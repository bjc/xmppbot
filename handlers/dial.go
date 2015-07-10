package handlers

import (
	"fmt"
	"strings"
	"time"

	_ "github.com/ThomsonReutersEikon/nitro/src/sipbot" // "sip" scheme
	"github.com/ThomsonReutersEikon/open-nitro/src/bots"
	_ "github.com/ThomsonReutersEikon/open-nitro/src/bots/xmppclient" // Register "xmpp" and "xmpp-bosh" bot schemes
	"github.com/bjc/goctl"
)

const timeout = 1 * time.Second

var xb bots.Bot

type Dial struct{}

func (dh Dial) Name() string {
	return "dial"
}

func (dh Dial) Help() string {
	return "scheme jid password host[:port]"
}

func (dh Dial) Run(_ *goctl.Goctl, args []string) string {
	if xb != nil {
		return "ERROR: bot is already connected."
	}

	if len(args) != 4 {
		return "ERROR: dial requires scheme, JID, password, and host[:port]"
	}

	url := fmt.Sprintf("%s://%s,%s:%s@%s", args[0], args[1], emailFromJID(args[1]), args[2], args[3])

	var err error
	xb, err = bots.Dial(url, timeout)
	if err != nil {
		return fmt.Sprintf("ERROR: %s.", err)
	}
	return "ok"
}

func emailFromJID(jid string) string {
	i := strings.Index(jid, "@")
	return strings.Replace(jid[:i], ".", "@", 1)
}
