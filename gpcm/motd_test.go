package gpcm

import (
	"strings"
	"testing"
	"unicode/utf16"
)

func TestMessageWithHandle(t *testing.T) {
	base := "Welcome to\nOpenPak Wi-Fi Connection!"
	if got := MessageWithHandle(base, ""); got != base {
		t.Fatalf("no handle changed the message: %q", got)
	}
	got := MessageWithHandle(base, "tobagin")
	if !strings.HasPrefix(got, base) || !strings.HasSuffix(got, "OpenPak: @tobagin") {
		t.Fatalf("unexpected message %q", got)
	}
	if got := MessageWithHandle("", "tobagin"); strings.HasPrefix(got, "\n") {
		t.Fatalf("empty base leaves leading blank lines: %q", got)
	}
}

func TestMessageWithHandleFitsTheGameBuffer(t *testing.T) {
	got := MessageWithHandle(strings.Repeat("x", 1000), "tobagin")
	if n := len(utf16.Encode([]rune(got))); n > motdLimit {
		t.Fatalf("%d units, over the limit %d", n, motdLimit)
	}
	if !strings.HasSuffix(got, "@tobagin") {
		t.Fatalf("the handle was cut: %q", got)
	}
	if b64 := ((len(utf16.Encode([]rune(got)))*2 + 2) / 3) * 4; b64 >= 512 {
		t.Fatalf("base64 length %d does not fit 512", b64)
	}
}
