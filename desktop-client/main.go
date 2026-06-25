package main

import (
	"errors"
	"flag"

	"voice-bridge-client/internal/singleinstance"
	"voice-bridge-client/internal/ui"
)

func main() {
	defaultServer := flag.String("server", "wss://srf.cccz.cc/ws", "server websocket URL")
	flag.Parse()

	lock, err := singleinstance.Acquire("cccz-srf-voice-bridge-client")
	if errors.Is(err, singleinstance.ErrAlreadyRunning) {
		singleinstance.ShowAlreadyRunning()
		return
	}
	if err != nil {
		ui.WriteStartupError(err)
		return
	}
	defer lock.Release()

	if err := ui.Run(*defaultServer); err != nil {
		ui.WriteStartupError(err)
	}
}
