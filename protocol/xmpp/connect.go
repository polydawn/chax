package xmpp

// this file is full of connection bootstrapping stuff

import (
	"crypto/tls"
	"net"
	"time"

	"polydawn.net/chax/lib/log"
	"polydawn.net/chax/protocol"

	"github.com/agl/xmpp"
	"github.com/inconshreveable/log15"
)

func Dial(server protocol.ServerDesc, account protocol.Account, Log log15.Logger) protocol.Conn {
	xmppConfig := &xmpp.Config{
		InLog:          log.Writer{Log, log15.LvlDebug, " <- RECV <-"},
		OutLog:         log.Writer{Log, log15.LvlDebug, " -> SENT ->"},
		Log:            log.Writer{Log, log15.LvlDebug, " :: NOTE ::"},
		TrustedAddress: true, // current test account has cert for 'www.xmpp.pro', which is not our account domain; handle better later
		//ServerCertificateSHA256: certSHA256, // pinning todo
		TLSConfig: TLSConfig,
		SkipTLS:   true, // sigh, current test server is EOF'ing me.  presumably suite mismatches; no indications; libpurple is fine with it.
	}

	// do our own dial, because the default timeouts are... insane.  like, minutes.  plural.
	addr := server.Addr()
	sock, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		panic("Failed to connect to XMPP server: " + err.Error())
	}
	xmppConfig.Conn = sock

	// shake
	xmppConn, err := xmpp.Dial(addr, account.Username, account.Domain, account.Password, xmppConfig)
	if err != nil {
		panic("Failed to connect to XMPP server: " + err.Error())
	}

	// shrinkwrap for sale and launch actor
	conn := &Conn{
		raw:         xmppConn,
		server:      server,
		account:     account,
		commandChan: make(chan interface{}),
	}
	go conn.run()
	return conn
}

var TLSConfig = &tls.Config{
	MinVersion:   tls.VersionTLS10,
	CipherSuites: CipherSuites,
}

var CipherSuites = []uint16{
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
}
