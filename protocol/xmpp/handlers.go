package xmpp

import (
	"time"

	"polydawn.net/chax/protocol"

	"github.com/agl/xmpp"
)

func (conn *Conn) dispatchClientMessage(stanza *xmpp.ClientMessage) {
	switch stanza.Type {
	case "error":
		conn.handleClientMessageError(stanza)
	default:
		conn.handleClientMessageRegular(stanza)
	}
	// TODO be advised that afaict, typing notifications are also under this heading
	// add better debugging formats... current log output doesn't contain unrecognized fields/tags :I
}

func (conn *Conn) handleClientMessageError(stanza *xmpp.ClientMessage) {
	conn.log.Warn("Error reported remote client", "from", stanza.From, "message", stanza.Body)
}

func (conn *Conn) handleClientMessageRegular(stanza *xmpp.ClientMessage) {
	msg := conn.parseMessage(stanza)
	conn.log.Info("MESSAGE", "msg", msg)
	conv := conn.getConversation(msg.Author)
	conv.pushMessage(msg) // sync's internally; this is where the conn actor hands off.
}

func (conn *Conn) parseMessage(stanza *xmpp.ClientMessage) protocol.Message {
	// `conn` ref only here for the logger, consider removing
	msg := protocol.Message{}

	// TODO want
	// msg.UID

	// parse claimed author
	msg.Author = parseAccountFromString(stanza.From)

	// parse time
	if stanza.Delay != nil && len(stanza.Delay.Stamp) > 0 {
		// An XEP-0203 Delayed Delivery <delay/> element exists for this message, meaning that someone sent it while we were offline.
		// Show the timestamp for when the message was sent, rather than time.Now().
		var err error
		msg.TimeReceived, err = time.Parse(time.RFC3339, stanza.Delay.Stamp)
		if err != nil {
			conn.log.Warn("Ridiculous time format from remote client", "from", stanza, "timestr", stanza.Delay.Stamp)
		}
	} else {
		// Local clocks are insane and we'd rather store both local perception and senders perception, but we'll take what we can get.
		// TODO We dearly wish for a universal coordinated time tag :/ If there's a XEP that allows this, plz, plz implement.
		msg.TimeReceived = time.Now()
	}

	// parse message body.  may be html in there -- deal with that downstream during presentation as necessary.
	msg.Body = []byte(stanza.Body)

	return msg
}

func (conn *Conn) dispatchClientIQ(reqStanza *xmpp.ClientIQ) {
	// TODO entire method is janky cargo cult, if possible fold it back into an nice flat switch statement in the main loop.
	//  that organization is preferable so we can see all the message types we support in one place.
	if reqStanza.Type != "get" && reqStanza.Type != "set" {
		return
	}
	// TODO consider the IQ subtype and come up with a reply // reply := s.processIQ(reqStanza)
	var reply interface{}
	// TODO also this does not seem to me like it should necessarily be done while blocking the main event loop; review
	if reply == nil {
		reply = xmpp.ErrorReply{
			Type:  "cancel",
			Error: xmpp.ErrorBadRequest{},
		}
	}
	if err := conn.raw.SendIQReply(reqStanza.From, "result", reqStanza.Id, reply); err != nil {
		conn.log.Warn("failed to send IQ message", "reply", reply, "error", err)
	}
}
