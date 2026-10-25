package gpcm

import (
	"net/http"
	"openpak/wfc-common/wfcapi"
	"sync"
	"time"
)

var publishedMu sync.RWMutex
var published = map[uint32]wfcapi.Session{}

func publish(g *GameSpySession) {
	g.PublishedRevision = time.Now().UnixNano()
	v := wfcapi.Session{Login: wfcapi.Login{ProfileID: g.User.ProfileId, Generation: g.Generation, Revision: g.PublishedRevision, GameCode: g.GameCode, InGameName: g.InGameName, ConsoleFriendCode: g.ConsoleFriendCode, FriendKeyGame: g.User.GsbrCode[:4], GPPublicIP: g.RemoteAddr, NeedsExploit: g.NeedsExploit, DeviceAuthenticated: g.DeviceAuthenticated, Restricted: g.User.Restricted}, SessionKey: g.SessionKey, GameName: g.GameName, UniqueNick: "000000000" + g.User.GsbrCode[:4] + "0000000"}
	publishedMu.Lock()
	published[g.User.ProfileId] = v
	publishedMu.Unlock()
}
func unpublish(pid uint32, generation string) {
	publishedMu.Lock()
	defer publishedMu.Unlock()
	if v, ok := published[pid]; ok && v.Login.Generation == generation {
		delete(published, pid)
	}
}
func RegisterInternal(mux *http.ServeMux) {
	mux.HandleFunc("POST /internal/v1/gpcm/{op}", func(w http.ResponseWriter, r *http.Request) {
		var in wfcapi.GPCMRequest
		if !wfcapi.Decode(w, r, &in) {
			return
		}
		out := wfcapi.GPCMResponse{}
		switch r.PathValue("op") {
		case "sessions":
			publishedMu.RLock()
			for _, v := range published {
				out.Sessions = append(out.Sessions, v)
			}
			publishedMu.RUnlock()
		case "verify":
			publishedMu.RLock()
			v, ok := published[in.ProfileID]
			publishedMu.RUnlock()
			if ok && v.SessionKey == in.SessionKey && v.GameName == in.GameName {
				out.Valid = true
				out.UniqueNick = v.UniqueNick
			}
		case "callback":
			if in.Generation == "" {
				http.Error(w, "missing generation", 400)
				return
			}
			mutex.Lock()
			if s, ok := sessions[in.ProfileID]; ok && s.Generation == in.Generation {
				kickPlayer(in.ProfileID, in.Reason)
			}
			mutex.Unlock()
		case "kick":
			if in.Custom == "restricted" {
				KickPlayerCustomMessage(in.ProfileID, in.Reason, WWFCMsgProfileRestrictedCustom)
			} else {
				KickPlayerCustomMessage(in.ProfileID, in.Reason, WWFCMsgKickedCustom)
			}
		default:
			http.NotFound(w, r)
			return
		}
		wfcapi.Reply(w, out)
	})
}
