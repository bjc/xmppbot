package handlers

import "github.com/bjc/goctl"

type Stop struct {
	C chan bool
}

func (sh Stop) Name() string {
	return "stop"
}

func (sh Stop) Help() string {
	return "stops this bot"
}

func (sh Stop) Run(_ *goctl.Goctl, args []string) string {
	if xb != nil {
		xb.Shutdown()
		xb = nil
	}

	sh.C <- true
	return "Stopping"
}
