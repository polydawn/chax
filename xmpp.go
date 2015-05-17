package main

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"github.com/agl/xmpp"
)

type Account struct {
	Username string
	Domain   string
	Password string
}

func (a Account) FQAN() string {
	return fmt.Sprintf("%s@%s", a.Username, a.Domain)
}

type ServerDesc struct {
	Host string
	Port uint16
}

func Resolve(domain string) (sd ServerDesc) {
	var err error
	sd.Host, sd.Port, err = xmpp.Resolve(domain)
	if err != nil {
		panic("Failed to resolve XMPP server: " + err.Error())
	}
	return
}

func awaitVersionReply(ch <-chan xmpp.Stanza, user string) {
	stanza, ok := <-ch
	if !ok {
		panic("Version request to " + user + " timed out")
	}
	reply, ok := stanza.Value.(*xmpp.ClientIQ)
	if !ok {
		panic("Version request to " + user + " resulted in bad reply type")
	}

	if reply.Type == "error" {
		panic("Version request to " + user + " resulted in XMPP error")
	} else if reply.Type != "result" {
		panic("Version request to " + user + " resulted in response with unknown type: " + reply.Type)
	}

	buf := bytes.NewBuffer(reply.Query)
	var versionReply xmpp.VersionReply
	if err := xml.NewDecoder(buf).Decode(&versionReply); err != nil {
		panic("Failed to parse version reply from " + user + ": " + err.Error())
	}

	Log.Info(fmt.Sprintf("Version reply from %s: %#v", user, versionReply))
}
