package gpcm

// Link-first (prds repos/wfc-link-first-design.md, card prds#139). A console plays online only
// once it is linked to an OpenPak account. wfc-nas refuses every unlinked console at NAS login
// with its link number (shown as "Error Code: 2xxxx") except Mario Kart Wii, which comes here
// and is refused at GPCM login with text (the 8-symbol code and the same number, error 22100).
// This gate also catches any other console that reaches GPCM unlinked: a Wii title is sent the
// link number as wl:err (Brawl shows it, some titles show nothing), a DS its stock error.
//
// A Wii on the DNS-exploit path logs in without the payload; it is let in with no address kept
// and flagged unlinked so QR2 can deliver the payload and nothing else, and its exlogin (the
// first point the payload can show the code) is refused.
//
// WFC_LINK_FIRST=1 turns the gate on. Off, logins behave as before link-first.

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	"openpak/wfc-common/common"
	"openpak/wfc-common/database"
	"openpak/wfc-common/gperror"
	"openpak/wfc-common/logging"
	"openpak/wfc-common/mkwboard"
)

var (
	// linkCore heals a half-done approval (core link made, local cache not written). Nil
	// without WFC_CORE_ADDRESS: a half-done approval then waits for wfc-api's reaper.
	linkCore database.CoreChecker
	// linkDSErr is the GP err value a refused DS is sent (0: the generic login error).
	linkDSErr, _ = strconv.Atoi(os.Getenv("WFC_LINK_DS_ERR"))
)

func startLinkGate(cfg common.Config) {
	if !database.LinkFirst() {
		logging.Notice("GPCM", "link-first gate is off (WFC_LINK_FIRST unset)")
		return
	}
	if cfg.CoreAddress != "" {
		if links, _, err := mkwboard.DialCore(cfg.CoreAddress, cfg.CoreKey); err == nil {
			linkCore = links
		} else {
			logging.Error("GPCM", "link gate: account core dial:", err)
		}
	}
	logging.Notice("GPCM", "link-first gate is on: unlinked consoles are refused with their pairing code")
}

func (g *GameSpySession) remoteIP() string {
	ip := g.RemoteAddr
	if i := strings.LastIndex(ip, ":"); i >= 0 {
		ip = ip[:i]
	}
	return ip
}

// linkGate decides a login. It returns proceed=true to continue the login; unlinked=true
// when the login may only continue as an unlinked exploit session (allowUnlinked). On
// proceed=false the refusal has been sent.
func (g *GameSpySession) linkGate(userID uint64, gsbrcd string, deviceID uint32, allowUnlinked bool) (proceed, unlinked bool) {
	if !database.LinkFirst() {
		return true, false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	r, err := db.LinkGate(ctx, userID, linkCore)
	if r == database.GateLinked {
		return true, false
	}
	if r == database.GateUnavailable {
		logging.Warn(g.ModuleName, "link gate: refused user", userID, "reason=core_unavailable:", err)
		g.replyError(gperror.LinkUnavailableError())
		return false, false
	}

	// Unlinked. A banned console is told it is banned, never shown a code.
	if banned, reason, err := db.ConsoleBanned(ctx, userID, deviceID, g.remoteIP()); err != nil {
		logging.Warn(g.ModuleName, "link gate: ban check:", err)
		g.replyError(gperror.LinkUnavailableError())
		return false, false
	} else if banned {
		logging.Notice(g.ModuleName, "link gate: refused user", userID, "reason=console_banned")
		g.replyError(GPError{ErrorCode: ErrLogin.ErrorCode, ErrorString: "This console is banned.", Fatal: true,
			WWFCMessage: WWFCMsgProfileBannedTOS, Reason: reason})
		return false, false
	}

	if allowUnlinked {
		logging.Notice(g.ModuleName, "link gate: unlinked exploit login for user", userID, "game", g.GameCode, "(payload first, code at exlogin)")
		return true, true
	}
	g.refuseWithCode(ctx, userID, gsbrcd)
	return false, false
}

// refuseWithCode issues (or reuses) the console's pairing code and refuses the login with it.
func (g *GameSpySession) refuseWithCode(ctx context.Context, userID uint64, gsbrcd string) {
	unit := byte(database.UnitWii)
	if g.UnitCode == UnitCodeDS {
		unit = database.UnitDS
	}
	p, err := db.IssuePairingCode(ctx, database.PairingRequest{UserID: userID, Unit: unit, GameCode: gsbrcd, Name: g.InGameName,
		NetTag: database.NetTag(database.NetTagKey(), g.remoteIP(), time.Now())})
	if err != nil {
		logging.Warn(g.ModuleName, "link gate: refused user", userID, "reason=no_code:", err)
		g.replyError(gperror.LinkUnavailableError())
		return
	}
	dsErr := 0
	if g.UnitCode == UnitCodeDS {
		dsErr = linkDSErr
	}
	logging.Notice(g.ModuleName, "link gate: refused user", userID, "unit", unitName(g.UnitCode), "game", g.GameCode,
		"reason=not_linked code="+database.CodeHash(p.Code))
	g.replyError(gperror.NotLinkedError(g.GameName, database.FormatPairingCode(p.Code), database.DisplayNumber(p.Number), dsErr))
}

func unitName(u byte) string {
	if u == UnitCodeDS {
		return "ds"
	}
	return "wii"
}

// kickUser closes every session of a console: the account core unlinked it, deleted its
// account or banned it. Callers hold mutex.
func kickUser(userID uint64, reason string) int {
	n := 0
	for pid, s := range sessions {
		if s.User.UserId != userID || userID == 0 {
			continue
		}
		n++
		switch reason {
		case "banned":
			kickPlayer(pid, "banned")
		default:
			s.replyError(gperror.UnlinkedKickError())
		}
	}
	if n > 0 {
		logging.Notice("GPCM", "link gate: kicked user", userID, "reason="+reason, "sessions", n)
	}
	return n
}
