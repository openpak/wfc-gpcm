package gpcm

import (
	"context"
	"os"
	"strconv"
	"strings"
	"testing"
	"unicode/utf16"

	"openpak/wfc-common/common"
	"openpak/wfc-common/database"
)

// The link gate against a real Postgres (WFC_TEST_DSN, a wwfc database with the users table):
// what an unlinked console is sent, per title and unit, and that a linked one goes on.
func TestLinkGateReplies(t *testing.T) {
	dsn := os.Getenv("WFC_TEST_DSN")
	if dsn == "" {
		t.Skip("WFC_TEST_DSN not set")
	}
	t.Setenv("WFC_LINK_FIRST", "1")
	cfg := common.Config{}
	// database.Start wants the pieces; parse them from the DSN postgres://u:p@host/db.
	rest := strings.TrimPrefix(dsn, "postgres://")
	cred, hostdb, _ := strings.Cut(rest, "@")
	cfg.Username, cfg.Password, _ = strings.Cut(cred, ":")
	cfg.DatabaseAddress, cfg.DatabaseName, _ = strings.Cut(hostdb, "/")
	db = database.Start(cfg)
	defer db.Close()
	if err := db.EnsurePairing(); err != nil {
		t.Fatal(err)
	}

	var sent []string
	beforeSend, beforeClose := sendPacket, closeConnection
	sendPacket = func(server string, idx uint64, b []byte) error { sent = append(sent, string(b)); return nil }
	closeConnection = func(server string, idx uint64) error { return nil }
	defer func() { sendPacket, closeConnection = beforeSend, beforeClose }()

	ctx := context.Background()
	const wii, ds = uint64(5550001), uint64(5550002)
	cleanup := func() {
		for _, u := range []uint64{wii, ds} {
			_ = db.ForgetPairing(ctx, u)
			_, _ = db.EraseUsers(ctx, []int64{int64(u)})
		}
	}
	cleanup()
	defer cleanup()

	mk := func(game string, unit byte) *GameSpySession {
		return &GameSpySession{ModuleName: "test", GameName: game, UnitCode: unit, LoginInfoSet: true, InGameName: "Tobagin", RemoteAddr: "198.51.100.9:1234", Language: 1}
	}

	// Mario Kart Wii: text with the code and the number.
	sent = nil
	if ok, _ := mk("mariokartwii", UnitCodeWii).linkGate(wii, "RMCP", 0, false); ok || len(sent) != 1 {
		t.Fatalf("unlinked Wii proceeded or sent %d packets", len(sent))
	}
	if !strings.Contains(sent[0], `\wl:err\22100`) || !strings.Contains(sent[0], `\wl:errmsg\`) {
		t.Fatalf("Mario Kart Wii reply: %q", sent[0])
	}
	text := errMsgText(t, sent[0])
	target, err := db.LookupPairingCode(ctx, digitsIn(text), database.NumberProof{Name: "Tobagin"})
	if err != nil || target.UserID != wii || target.GameCode != "RMCP" || target.Name != "Tobagin" {
		t.Fatalf("the shown number does not resolve: %q -> %+v %v", text, target, err)
	}
	if !strings.Contains(text, database.FormatPairingCode(target.Code)) {
		t.Fatalf("text lacks the code: %q", text)
	}

	// Another Wii title: the number is the error code, no text.
	sent = nil
	mk("smashbrawlx", UnitCodeWii).linkGate(wii, "RSBE", 0, false)
	if strings.Contains(sent[0], "wl:errmsg") || !strings.Contains(sent[0], `\wl:err\`+strconv.Itoa(numberOf(t, sent[0]))) {
		t.Fatalf("Brawl reply: %q", sent[0])
	}
	if n := numberOf(t, sent[0]); n < 24000 || n > 29999 {
		t.Fatalf("number %d", n)
	}

	// Wii on the exploit path: let in, flagged, no reply.
	sent = nil
	if ok, unlinked := mk("smashbrawlx", UnitCodeWii).linkGate(wii, "RSBE", 0, true); !ok || !unlinked || len(sent) != 0 {
		t.Fatalf("exploit path: ok %v unlinked %v sent %d", ok, unlinked, len(sent))
	}

	// DS with a configured err value.
	linkDSErr = 263
	defer func() { linkDSErr = 0 }()
	sent = nil
	mk("mariokartds", UnitCodeDS).linkGate(ds, "AMCE", 0, false)
	if !strings.Contains(sent[0], `\err\263\`) {
		t.Fatalf("DS reply: %q", sent[0])
	}

	// Linked: goes on, sends nothing.
	tg, _, err := db.BeginApproval(ctx, strconv.Itoa(numberOf(t, sent[0])), "acct-test", database.NumberProof{Name: "Tobagin"})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.FinishApproval(ctx, tg, "acct-test", "link-test"); err != nil {
		t.Fatal(err)
	}
	sent = nil
	if ok, unlinked := mk("mariokartds", UnitCodeDS).linkGate(ds, "AMCE", 0, false); !ok || unlinked || len(sent) != 0 {
		t.Fatalf("linked DS: ok %v unlinked %v sent %v", ok, unlinked, sent)
	}

	// Gate off: everyone goes on.
	t.Setenv("WFC_LINK_FIRST", "")
	if ok, _ := mk("mariokartwii", UnitCodeWii).linkGate(wii, "RMCP", 0, false); !ok {
		t.Fatal("gate off refused")
	}
}

func errMsgText(t *testing.T, msg string) string {
	t.Helper()
	i := strings.Index(msg, `\wl:errmsg\`)
	v := msg[i+len(`\wl:errmsg\`):]
	v = v[:strings.Index(v, `\`)]
	b, err := common.Base64DwcEncoding.DecodeString(v)
	if err != nil {
		t.Fatal(err)
	}
	u := make([]uint16, len(b)/2)
	for k := range u {
		u[k] = uint16(b[2*k])<<8 | uint16(b[2*k+1])
	}
	return string(utf16.Decode(u))
}

// digitsIn finds the 6-digit link number in a message.
func digitsIn(s string) string {
	run := ""
	for _, r := range s + " " {
		if r >= '0' && r <= '9' {
			run += string(r)
			continue
		}
		if len(run) == 5 && run != "22100" {
			return run
		}
		run = ""
	}
	return ""
}

func numberOf(t *testing.T, msg string) int {
	t.Helper()
	i := strings.Index(msg, `\wl:err\`)
	if i < 0 {
		t.Fatalf("no wl:err in %q", msg)
	}
	v := msg[i+len(`\wl:err\`):]
	v = v[:strings.Index(v, `\`)]
	n, _ := strconv.Atoi(v)
	return n
}
