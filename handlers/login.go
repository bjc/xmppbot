package handlers

import (
	"fmt"

	"github.com/bjc/goctl"
)

type Login struct{}

func (lh Login) Name() string {
	return "login"
}

func (lh Login) Help() string {
	return "authenticates to server"
}

func (lh Login) Run(_ *goctl.Goctl, args []string) string {
	if xb == nil {
		return "ERROR: bot is not connected."
	}

	if err := xb.Login(); err != nil {
		return fmt.Sprintf("ERROR: couldn't login: %s.", err)
	}
	return "ok"
}
