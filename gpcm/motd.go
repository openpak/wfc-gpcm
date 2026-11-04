package gpcm

import (
	"context"
	"time"
	"os"
	"unicode/utf16"

	"openpak/wfc-common/database"
	"openpak/wfc-common/logging"
)

var motdFilepath = "./motd.txt"

func GetMessageOfTheDay() (string, error) {
	contents, err := os.ReadFile(motdFilepath)
	if err != nil {
		return "", err
	}

	return string(contents), nil
}

// motdLimit is how many UTF-16 units of message we send. The game decodes the base64 value
// into a 512-byte buffer (gpiValueForKey in the payload's mkwFeature.cpp), which holds 384
// bytes, 192 units; a little under that is kept as margin.
const motdLimit = 180

// MessageWithPairing is the message of the day with the console's pairing code appended: the
// only text an unlinked Mario Kart Wii player is guaranteed to read, which is how the code
// reaches them. An empty code (console already linked, or the database was unreachable)
// leaves the message alone. The code lines come first-class: if the whole thing would not fit
// the game's buffer, the base message is cut, never the code.
func MessageWithPairing(base, code string) string {
	if code == "" {
		return base
	}
	tail := "\n\nLink this Wii at openpak.org/link\nCode: " + database.FormatPairingCode(code)
	if base == "" {
		tail = tail[2:]
	}
	room := motdLimit - len(utf16.Encode([]rune(tail)))
	units := utf16.Encode([]rune(base))
	if room < 0 {
		return tail
	}
	if len(units) > room {
		units = units[:room]
	}
	return string(utf16.Decode(units)) + tail
}

// pairingCode is the code an unlinked console should be shown, or "" when it is linked or the
// database does not answer in time: a login never waits on, or fails because of, linking.
func (g *GameSpySession) pairingCode() string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	code, err := db.IssuePairingCode(ctx, g.User.UserId)
	if err != nil {
		logging.Warn(g.ModuleName, "pairing code:", err)
		return ""
	}
	return code
}
