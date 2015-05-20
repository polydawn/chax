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
	heartbeatTicker := time.NewTicker(10 * time.Second)
	// main loop
	for {
		select {
		case <-heartbeatTicker.C:
			conn.raw.SignalPresence("")
		case rawStanza := <-stanzaChan:
			switch stanza := rawStanza.Value.(type) {
			case *xmpp.ClientMessage:
				conn.log.Debug("unhandled message stanza", "stanza", stanza)
			case *xmpp.ClientPresence:
				conn.log.Debug("unhandled presence stanza", "stanza", stanza)
			case *xmpp.ClientIQ:
				conn.log.Debug("unhandled IQ request stanza", "stanza", stanza)
				if stanza.Type != "get" && stanza.Type != "set" {
					continue
				}
				// TODO control the spiraling inanity of this parsing
				// reply := s.processIQ(stanza)
				var reply interface{}
				// TODO also this does not seem to me like it should necessarily be done while blocking the main event loop; review
				if reply == nil {
					reply = xmpp.ErrorReply{
						Type:  "cancel",
						Error: xmpp.ErrorBadRequest{},
					}
				}
				if err := conn.raw.SendIQReply(stanza.From, "result", stanza.Id, reply); err != nil {
					conn.log.Warn("failed to send IQ message", "reply", reply, "error", err)
				}
			case *xmpp.StreamError:
				var text string
				if len(stanza.Text) > 0 {
					text = stanza.Text
				} else {
					text = fmt.Sprintf("%s", stanza.Any)
				}
				conn.log.Warn("error from server", "stanza", text)
			default:
				conn.log.Warn("unrecognized stanza", "name", rawStanza.Name, "value", rawStanza.Value)
			}
		case cmd := <-conn.commandChan:
			switch cmd := cmd.(type) {
			case func():
				cmd()
			}
		}
	}
}
