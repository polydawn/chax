package xmpp

import (
	"time"

	"polydawn.net/chax/protocol"

	"github.com/agl/xmpp"
)

var _ protocol.Conn = &Conn{}

type Conn struct {
	raw         *xmpp.Conn // agl's connection type
	server      protocol.ServerDesc
	account     protocol.Account
	commandChan chan interface{}
}

/*
	Actor main method.  Call me once, in a goroutine you don't expect to return.
*/
func (conn *Conn) run() {
	heartbeatTicker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-heartbeatTicker.C:
			conn.raw.SignalPresence("")
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
