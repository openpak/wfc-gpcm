package main

import (
	"context"
	"net/http"
	"openpak/wfc-common/service"
	"openpak/wfc-gpcm/gpcm"
	"os"
	"time"
)

var httpServer *http.Server

func startHTTP(bool) {
	mux := http.NewServeMux()
	gpcm.RegisterInternal(mux)
	h := http.NewServeMux()
	h.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok\n")) })
	h.Handle("/", service.Guard(mux))
	addr := os.Getenv("WFC_HTTP_ADDR")
	if addr == "" {
		addr = ":20200"
	}
	httpServer = service.Server(addr, h)
	go func() {
		if e := httpServer.ListenAndServe(); e != nil && e != http.ErrServerClosed {
			panic(e)
		}
	}()
}
func stopHTTP() {
	if httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		httpServer.Shutdown(ctx)
	}
}
