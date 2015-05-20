package xmpp

// misc parsing and data massaging functions

import (
	"polydawn.net/chax/protocol"

	"github.com/agl/xmpp"
)

func parseAccountFromString(str string) protocol.Account {
	from := xmpp.RemoveResourceFromJid(str)
	// TODO care
	_ = from
	return protocol.Account{}
}
