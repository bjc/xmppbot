package handlers

import (
	"fmt"

	"github.com/bjc/goctl"
)

type Presence struct{}

func (ph Presence) Name() string {
	return "presence"
}

func (ph Presence) Help() string {
	return "sends available presence broadcast"
}

func (ph Presence) Run(_ *goctl.Goctl, args []string) string {
	if xb == nil {
		return "ERROR: bot is not connected."
	}

	if err := xb.Online(); err != nil {
		return fmt.Sprintf("ERROR: sending available presence: %s", err)
	}
	return "ok"
}
