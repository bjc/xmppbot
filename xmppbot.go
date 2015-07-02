package main

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/inconshreveable/log15.v2"

	"github.com/ThomsonReutersEikon/open-nitro/src/bots"
	"github.com/ThomsonReutersEikon/open-nitro/src/bots/xmppclient" // Register "xmpp" and "xmpp-bosh" bot schemes
	"github.com/bjc/goctl"
)

const sockPath = "/tmp/xmppbot"
const timeout = 1 * time.Second

var (
	gc       goctl.Goctl
	xb       *xmppclient.Bot
	stopChan chan bool
)

func emailFromJID(jid string) string {
	i := strings.Index(jid, "@")
	return strings.Replace(jid[:i], ".", "@", 1)
}

func dialHandler(args []string) string {
	if xb != nil {
		return "ERROR: bot is already connected."
	}

	if len(args) != 3 {
		return "ERROR: dial requires JID, password, and host[:port]"
	}

	url := fmt.Sprintf("xmpp-bosh://%s,%s:%s@%s", args[0], emailFromJID(args[0]), args[1], args[2])

	var err error
	b, err := bots.Dial(url, timeout)
	if err != nil {
		return fmt.Sprintf("ERROR: %s.", err)
	}

	var ok bool
	xb, ok = b.(*xmppclient.Bot)
	if !ok {
		return "ERROR: bot wasn't xmpp client."
	}
	return "ok"
}

func loginHandler(args []string) string {
	if xb == nil {
		return "ERROR: bot is not connected."
	}

	if err := xb.Login(); err != nil {
		return fmt.Sprintf("ERROR: couldn't login: %s.", err)
	}
	return "ok"
}

func bindHandler(args []string) string {
	if xb == nil {
		return "ERROR: bot is not connected."
	}

	if err := xb.Bind(); err != nil {
		return fmt.Sprintf("ERROR: couldn't bind resource: %s.", err)
	}
	return "ok"
}

func stopHandler(args []string) string {
	if xb == nil {
		return "ERROR: bot is not connected."
	}
	xb.Shutdown()
	xb = nil

	stopChan <- true
	return "Stopping"
}

func presenceHandler(args []string) string {
	if xb == nil {
		return "ERROR: bot is not connected."
	}
	xb.Sendf(`<presence/>`)
	return "ok"
}

func pingHandler(args []string) string {
	if xb == nil {
		return "ERROR: bot is not connected."
	}

	xb.Sendf(`<iq type='get' to='%s' id='ping'><ping xmlns='urn:xmpp:ping'/></iq>`, xb.JID().Domain())
	return "ok"
}

func main() {
	stopChan = make(chan bool, 1)

	goctl.Logger.SetHandler(log15.StdoutHandler)

	gc = goctl.NewGoctl(sockPath)
	gc.AddHandlers([]*goctl.Handler{
		{"dial", dialHandler},
		{"login", loginHandler},
		{"bind", bindHandler},
		{"stop", stopHandler},
		{"presence", presenceHandler},
		{"ping", pingHandler},
	})
	if err := gc.Start(); err != nil {
		log15.Crit("Coudln't start command listener.", "error", err)
		return
	}
	defer gc.Stop()

	<-stopChan
}