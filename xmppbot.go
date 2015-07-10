package main

import (
	"flag"

	"github.com/bjc/goctl"
	"github.com/bjc/xmppbot/handlers"
	"gopkg.in/inconshreveable/log15.v2"
)

func main() {
	var sockPath string
	flag.StringVar(&sockPath, "f", "/tmp/xmppbot", "path to UNIX control socket")
	flag.Parse()

	stopChan := make(chan bool, 1)

	goctl.Logger.SetHandler(log15.StdoutHandler)

	gc := goctl.NewGoctl(sockPath)
	if err := gc.AddHandlers([]goctl.Handler{
		handlers.Dial{},
		handlers.Login{},
		handlers.Bind{},
		handlers.Stop{C: stopChan},
		handlers.Presence{},
		handlers.Raw{},
	}); err != nil {
		log15.Crit("Couldn't set up command handlers.", "error", err)
		return
	}
	if err := gc.Start(); err != nil {
		log15.Crit("Coudln't start command listener.", "error", err)
		return
	}
	defer gc.Stop()

	<-stopChan
}
