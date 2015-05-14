package main

import (
	"github.com/agl/xmpp"
)

type Account struct {
	Username string
	Domain   string
	Password string
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
