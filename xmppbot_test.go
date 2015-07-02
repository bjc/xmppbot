package main

import (
	"testing"
	"time"
)

const TIMEOUT = 10 * time.Second

const JID = "xmppload0.reuters.com@array12.msgtst.reuters.com"
const EMAIL_ADDR = "xmppload0@reuters.com"
const PASSWORD = "Welcome1"
const HOST = "localhost:5222"

func TestConnect(t *testing.T) {
	bot := NewBot(JID, PASSWORD, HOST, TIMEOUT)
	if err := bot.Dial(); err != nil {
		t.Fatalf("Couldn't connect: %s.", err)
	}

	if err := bot.Login(); err != nil {
		t.Fatalf("Couldn't login bot: %s.", err)
	}

	if err := bot.Bind(); err != nil {
		t.Fatalf("Couldn't bind to resource: %s.", err)
	}
}
