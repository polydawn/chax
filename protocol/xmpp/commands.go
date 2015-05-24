package xmpp

import (
	"polydawn.net/chax/protocol"
)

type newConversationCmd struct {
	recipient protocol.Account
	ack       chan *conversation
}
