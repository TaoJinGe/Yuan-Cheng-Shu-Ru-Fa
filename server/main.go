package main

import (
	"log"
	"net/http"

	"voice-bridge-server/internal/auth"
	"voice-bridge-server/internal/config"
	"voice-bridge-server/internal/hub"
	"voice-bridge-server/internal/records"
	"voice-bridge-server/internal/web"
)

func main() {
	cfg := config.Load()

	sessions, err := auth.NewStore(cfg)
	if err != nil {
		log.Fatal(err)
	}
	recorder := records.New(cfg.RecordsDir)
	messageHub := hub.New(sessions, recorder)
	handler := web.NewHandler(cfg, sessions, messageHub)

	addr := ":" + cfg.Port
	log.Printf("voice bridge server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, handler.Routes()))
}
