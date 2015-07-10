package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"gopkg.in/inconshreveable/log15.v2"

	_ "github.com/ThomsonReutersEikon/nitro/src/sipbot" // "sip" scheme
	"github.com/ThomsonReutersEikon/open-nitro/src/bots"
	"github.com/ThomsonReutersEikon/open-nitro/src/bots/xmppclient" // Register "xmpp" and "xmpp-bosh" bot schemes
	"github.com/bjc/goctl"
)

const timeout = 1 * time.Second

var (
	gc       goctl.Goctl
	xb       bots.Bot
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

	if err := xb.Online(); err != nil {
		return fmt.Sprintf("ERROR: sending available presence: %s", err)
	}
	return "ok"
}

func rawHandler(args []string) string {
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

func main() {
	var sockPath string
	flag.StringVar(&sockPath, "f", "/tmp/xmppbot", "path to UNIX control socket")
	flag.Parse()

	stopChan = make(chan bool, 1)

	goctl.Logger.SetHandler(log15.StdoutHandler)

	gc = goctl.NewGoctl(sockPath)
	gc.AddHandlers([]*goctl.Handler{
		{"dial", dialHandler},
		{"login", loginHandler},
		{"bind", bindHandler},
		{"stop", stopHandler},
		{"presence", presenceHandler},
		{"raw", rawHandler},
	})
	if err := gc.Start(); err != nil {
		log15.Crit("Coudln't start command listener.", "error", err)
		return
	}
	defer gc.Stop()

	<-stopChan
}
