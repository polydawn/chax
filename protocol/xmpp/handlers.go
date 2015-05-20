package xmpp

import (
	"github.com/agl/xmpp"
)

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
