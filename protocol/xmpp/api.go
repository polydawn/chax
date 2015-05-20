package xmpp

// this file is full of all the publicly exported methods that satsify the `protocol.Conn` interface.

import (
	"polydawn.net/chax/protocol"
)

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
