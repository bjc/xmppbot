package handlers

import (
	"fmt"

	"github.com/ThomsonReutersEikon/nitro/src/bots/xmppclient"
	"github.com/bjc/goctl"
)

type Raw struct{}

func (rh Raw) Name() string {
	return "raw"
}

func (rh Raw) Help() string {
	return "send raw data"
}

func (rh Raw) Run(_ *goctl.Goctl, args []string) string {
	if xb == nil {
		return "ERROR: bot is not connected."
	}

	if b, ok := xb.(*xmppclient.Bot); !ok {
		return "ERROR: can only send raw data on XMPP bot."
	} else if err := b.Sendf(args[0]); err != nil {
		return fmt.Sprintf("ERROR: sending raw data: %s", err)
	}

	return "ok"
}
