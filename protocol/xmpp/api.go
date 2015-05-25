package xmpp

// this file is full of all the publicly exported methods that satsify the `protocol.Conn` interface.

import (
	"polydawn.net/chax/protocol"
	"sync"
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

	sync.Mutex
	messages []protocol.Message
	awaiters chan chan struct{}
}

func (conv *conversation) Send(body []byte) protocol.Message {
	// this is kind of cheating/easy since there's no actual ack to wait for :I
	conv.conn.commandChan <- func() {
		conv.conn.raw.Send(conv.recipient.FQAN(), string(body))
	}
	// TODO we'd like a UUID of some kind from the ack, and a "universal" time if possible
	msg := protocol.Message{
		Author: conv.conn.account,
		Body:   body,
	}
	conv.pushMessage(msg)
	return msg
}

// locks the message slice, appends, and notifies.
func (conv *conversation) pushMessage(msg protocol.Message) {
	conv.Lock()
	defer conv.Unlock()
	conv.conn.log.Debug("message pushed", "msg", msg, "cardinality", len(conv.messages))
	conv.messages = append(conv.messages, msg)
	select {
	case awaiter := <-conv.awaiters:
		close(awaiter)
	default:
		// when we would block, just leave.
	}
}

func (conv *conversation) GetMessages(start, end int) []protocol.Message {
	conv.Lock()
	defer conv.Unlock()
	if end == -1 {
		end = len(conv.messages)
	}
	return conv.messages[start:end]
}

func (conv *conversation) AwaitUpdates() <-chan struct{} {
	awaiter := make(chan struct{})
	conv.awaiters <- awaiter
	return awaiter
}
