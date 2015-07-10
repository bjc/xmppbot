package handlers

import (
	"fmt"

	"github.com/bjc/goctl"
)

type Bind struct{}

func (bh Bind) Name() string {
	return "bind"
}

func (bh Bind) Help() string {
	return "binds XMPP resource"
}

func (bh Bind) Run(_ *goctl.Goctl, args []string) string {
	if xb == nil {
		return "ERROR: bot is not connected."
	}

	if err := xb.Bind(); err != nil {
		return fmt.Sprintf("ERROR: couldn't bind resource: %s.", err)
	}
	return "ok"
}
