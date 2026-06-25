package main

import (
	"flag"

	"voice-bridge-client/internal/ui"
)

func main() {
	defaultServer := flag.String("server", "wss://srf.cccz.cc/ws", "server websocket URL")
	flag.Parse()

	if err := ui.Run(*defaultServer); err != nil {
		ui.WriteStartupError(err)
	}
}
