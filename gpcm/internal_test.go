package gpcm

import (
	"net/http"
	"net/http/httptest"
	"openpak/wfc-common/service"
	"openpak/wfc-common/wfcapi"
	"strings"
	"testing"
)

func TestSearchVerificationIsScopedToCurrentSessionAndGame(t *testing.T) {
	publishedMu.Lock()
	published = map[uint32]wfcapi.Session{7: {Login: wfcapi.Login{ProfileID: 7, Generation: "new"}, SessionKey: 123, GameName: "game", UniqueNick: "native-name"}}
	publishedMu.Unlock()
	t.Cleanup(func() { publishedMu.Lock(); published = map[uint32]wfcapi.Session{}; publishedMu.Unlock() })
	key := strings.Repeat("k", 32)
	t.Setenv("WFC_INTERNAL_KEY", key)
	mux := http.NewServeMux()
	RegisterInternal(mux)
	h := service.Guard(mux)
	for _, tc := range []struct {
		body  string
		valid bool
	}{{`{"ProfileID":7,"SessionKey":123,"GameName":"game"}`, true}, {`{"ProfileID":7,"SessionKey":124,"GameName":"game"}`, false}, {`{"ProfileID":7,"SessionKey":123,"GameName":"other"}`, false}} {
		r := httptest.NewRequest("POST", "/internal/v1/gpcm/verify", strings.NewReader(tc.body))
		r.Header.Set("X-WFC-Internal-Key", key)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		want := `"Valid":false`
		if tc.valid {
			want = `"Valid":true`
		}
		if w.Code != 200 || !strings.Contains(w.Body.String(), want) {
			t.Fatal(w.Code, w.Body)
		}
	}
	unpublish(7, "old")
	if _, ok := published[7]; !ok {
		t.Fatal("old connection removed new session")
	}
	unpublish(7, "new")
	if _, ok := published[7]; ok {
		t.Fatal("logout left a searchable session")
	}
}
func TestOldRoomCallbackCannotKickNewConnection(t *testing.T) {
	mutex.Lock()
	sessions[7] = &GameSpySession{Generation: "new"}
	mutex.Unlock()
	defer func() { mutex.Lock(); delete(sessions, 7); mutex.Unlock() }()
	mux := http.NewServeMux()
	RegisterInternal(mux)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest("POST", "/internal/v1/gpcm/callback", strings.NewReader(`{"ProfileID":7,"Generation":"old","Reason":"network_error"}`)))
	if w.Code != 200 {
		t.Fatal(w.Code)
	}
}
