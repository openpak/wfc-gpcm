package gpcm

import (
	"context"
	"os"
	"time"
	"unicode/utf16"

	"openpak/wfc-common/logging"
	"openpak/wfc-common/mkwboard"
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

// MessageWithHandle is the message of the day with "OpenPak: @handle" appended for a linked
// console, so a player notices if someone else claimed their console. An empty handle leaves the
// message alone. The handle line is kept whole: if both would not fit the game's buffer, the
// base message is cut.
func MessageWithHandle(base, handle string) string {
	if handle == "" {
		return base
	}
	tail := "\n\nOpenPak: @" + handle
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

// linkedHandle is the handle of the account this console is linked to, or "" (not linked, no
// core, or no answer in time: a login never waits on it).
func (g *GameSpySession) linkedHandle() string {
	links, ok := linkCore.(*mkwboard.CoreLinks)
	if !ok || links == nil {
		return ""
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	l, linked, err := db.LinkedConsoleOf(ctx, g.User.UserId)
	if err != nil || !linked {
		return ""
	}
	h, err := links.Handle(ctx, l.AccountID)
	if err != nil {
		logging.Warn(g.ModuleName, "handle:", err)
		return ""
	}
	return h
}
