package xmpp

// this file is full of all the publicly exported methods that satsify the `protocol.Conn` interface.

import (
	"polydawn.net/chax/protocol"
)

func (conn *Conn) StartConversation(recipient protocol.Account) protocol.Conversation {
	cmd := &newConversationCmd{
		recipient: recipient,
		ack:       make(chan *conversation, 1),
	}
	conn.commandChan <- cmd
	return <-cmd.ack
}

func (conn *Conn) AcceptConversations() <-chan protocol.Conversation {
	return conn.incoming
}

type conversation struct {
	conn      *Conn
	recipient protocol.Account

	messages []protocol.Message
	awaiters chan chan struct{}
}

func (conv *conversation) Send(msg []byte) protocol.Message {
	// this is kind of cheating/easy since there's no actual ack to wait for :I
	conv.conn.commandChan <- func() {
		conv.conn.raw.Send(conv.recipient.FQAN(), string(msg))
	}
	// TODO we'd like a UUID of some kind from the ack, and a "universal" time if possible
	return protocol.Message{
		Author: conv.conn.account,
		Body:   msg,
	}
}
func (conv *conversation) GetMessages(start, end int64) {
}
func (conv *conversation) AwaitUpdates() <-chan struct{} {
	awaiter := make(chan struct{})
	conv.awaiters <- awaiter
	return awaiter
}
