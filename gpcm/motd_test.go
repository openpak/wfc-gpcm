package gpcm

import (
	"strings"
	"testing"
	"unicode/utf16"
)

func TestMessageWithPairing(t *testing.T) {
	base := "Welcome to\nWiiLink Wi-Fi Connection!"
	if got := MessageWithPairing(base, ""); got != base {
		t.Fatalf("linked console changed the message: %q", got)
	}
	got := MessageWithPairing(base, "ABCD2345")
	if !strings.HasPrefix(got, base) || !strings.HasSuffix(got, "Code: ABCD-2345") || !strings.Contains(got, "openpak.org/link") {
		t.Fatalf("unexpected message %q", got)
	}
	if got := MessageWithPairing("", "ABCD2345"); strings.HasPrefix(got, "\n") {
		t.Fatalf("empty base leaves leading blank lines: %q", got)
	}
}

func TestMessageWithPairingFitsTheGameBuffer(t *testing.T) {
	long := strings.Repeat("x", 1000)
	got := MessageWithPairing(long, "ABCD2345")
	if n := len(utf16.Encode([]rune(got))); n > motdLimit {
		t.Fatalf("%d units, over the limit %d", n, motdLimit)
	}
	if !strings.HasSuffix(got, "ABCD-2345") {
		t.Fatalf("the code was cut: %q", got)
	}
	// base64 of the UTF-16 bytes must fit the game's 512-byte value buffer.
	if b64 := ((len(utf16.Encode([]rune(got)))*2 + 2) / 3) * 4; b64 >= 512 {
		t.Fatalf("base64 length %d does not fit 512", b64)
	}
}
