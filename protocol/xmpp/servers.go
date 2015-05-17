package xmpp

import (
	"polydawn.net/chax/protocol"

	"github.com/agl/xmpp"
)

/*
	Attempts to map a server domain name to a typical `ServerDesc` by use
	of SRV records, which is a standard that xmpp servers do often abide by.
*/
func ResolveServer(domain string) (sd protocol.ServerDesc) {
	var err error
	sd.Host, sd.Port, err = xmpp.Resolve(domain)
	if err != nil {
		panic("Failed to resolve XMPP server: " + err.Error())
	}
	return
}
