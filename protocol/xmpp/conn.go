package xmpp

import (
	"time"

	"polydawn.net/chax/protocol"

	"github.com/agl/xmpp"
	"github.com/inconshreveable/log15"
)

var _ protocol.Conn = &Conn{}

type Conn struct {
	raw         *xmpp.Conn // agl's connection type
	server      protocol.ServerDesc
	account     protocol.Account
	log         log15.Logger
	commandChan chan interface{}
}

/*
	Actor main method.  Call me once, in a goroutine you don't expect to return.
*/
func (conn *Conn) run() {
	// spawn another worker for slurping messages from the wire,
	//  so that we can just select on whole messages in the main loop.
	stanzaChan := make(chan xmpp.Stanza)
	go func() {
		defer close(stanzaChan)
		for {
			stanza, err := conn.raw.Next()
			if err != nil {
				panic(err.Error())
			}
			stanzaChan <- stanza
		}
	}()
	// expect to do heartbeats and maintenance periodically
	heartbeatTicker := time.NewTicker(5 * time.Second)
	// main loop
	for {
		select {
		case <-heartbeatTicker.C:
			conn.raw.SignalPresence("")
		case rawStanza := <-stanzaChan:
			// TODO ... many branches of processing
			conn.log.Debug("unhandled stanza", "stanza", rawStanza)
		case cmd := <-conn.commandChan:
			switch cmd := cmd.(type) {
			case func():
				cmd()
			}
		}
	}
}

func (conn *Conn) Send(recipient protocol.Account, body []byte) protocol.Message {
	// this is kind of cheating/easy since there's no actual ack to wait for :I
	conn.commandChan <- func() {
		conn.raw.Send(recipient.FQAN(), string(body))
	}
	// TODO we'd like a UUID of some kind from the ack, and a "universal" time if possible
	return protocol.Message{
		Author: conn.account,
		Body:   body,
	}
}
