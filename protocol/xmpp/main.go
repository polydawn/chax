package xmpp

import (
	"fmt"
	"time"

	"polydawn.net/chax/protocol"

	"github.com/agl/xmpp"
	"github.com/inconshreveable/log15"
)

var _ protocol.Conn = &Conn{}

type Conn struct {
	raw           *xmpp.Conn // agl's connection type
	server        protocol.ServerDesc
	account       protocol.Account
	log           log15.Logger
	commandChan   chan interface{}
	conversations map[protocol.Account]*conversation
	incoming      chan protocol.Conversation // buffered.  we dump new conversations here if we have to create them.
}

/*
	Actor main method.  Call me once, in a goroutine you don't expect to return.
*/
func (conn *Conn) run() {
	// the first thing we have to do upon a new connection is announce presence.
	//  if you don't announce presence, the server won't consider you online,
	//   nevermind the tcp dial and the whole auth shake, and won't send you messages.
	conn.raw.SignalPresence("")
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
	heartbeatTicker := time.NewTicker(10 * time.Second)
	// main loop
	for {
		select {
		case <-heartbeatTicker.C:
			conn.raw.SignalPresence("")
		case rawStanza := <-stanzaChan:
			switch stanza := rawStanza.Value.(type) {
			case *xmpp.ClientMessage:
				conn.dispatchClientMessage(stanza)
			case *xmpp.ClientPresence:
				conn.log.Debug("unhandled presence stanza", "stanza", stanza)
			case *xmpp.ClientIQ:
				conn.log.Debug("unhandled IQ request stanza", "stanza", stanza)
				conn.dispatchClientIQ(stanza)
			case *xmpp.StreamError:
				conn.log.Warn("error from server", "message", stringifyStreamError(stanza))
			default:
				conn.log.Warn("unrecognized stanza", "name", rawStanza.Name, "value", rawStanza.Value)
			}
		case cmd := <-conn.commandChan:
			switch cmd := cmd.(type) {
			case func():
				cmd()
			case *newConversationCmd:
				conv, ok := conn.conversations[cmd.recipient]
				if !ok {
					conv = conn.newConversation(cmd.recipient)
				}
				cmd.ack <- conv
			}
		}
	}
}

// get conversation; create if not exists.
// sends on `incoming` if it had to create
func (conn *Conn) getConversation(recipient protocol.Account) *conversation {
	conv, ok := conn.conversations[recipient]
	if !ok {
		conv = conn.newConversation(recipient)
		conn.incoming <- conv
	}
	return conv
}

func (conn *Conn) newConversation(recipient protocol.Account) *conversation {
	// must be called from the connection's master actor routine
	conv := &conversation{
		conn:      conn,
		recipient: recipient,
		messages:  make([]protocol.Message, 0, 20),
		awaiters:  make(chan chan struct{}, 999),
	}
	conn.conversations[recipient] = conv

	return conv
}

/*
	Stringify StreamError, preferring the human-readable text if any, and
	falling back to the xml blob if nothing better is around.
*/
func stringifyStreamError(stanza *xmpp.StreamError) string {
	if len(stanza.Text) > 0 {
		return stanza.Text
	}
	return fmt.Sprintf("%s", stanza.Any)
}
