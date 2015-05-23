package xmpp

import (
	"testing"

	"polydawn.net/chax/lib/log"
	"polydawn.net/chax/protocol"
	"polydawn.net/chax/protocol/tests"

	"github.com/inconshreveable/log15"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCoreCompliance(t *testing.T) {
	Convey("Spec Compliance: XMPP Protocol", t, func(c C) {
		// setup logging
		lg := log15.New()
		lg.SetHandler(log.PrinterHandler(c.Print, log15.LogfmtFormat()))

		// set up accounts
		server := ResolveServer("crypt.mn")
		acct1 := protocol.Account{
			Username: "testpilot",
			Domain:   "crypt.mn",
			Password: "asdf",
		}
		acct2 := protocol.Account{
			Username: "testpilot2",
			Domain:   "crypt.mn",
			Password: "asdf",
		}

		// run compliance tests
		tests.CheckMessageBounce(
			acct1, func() protocol.Conn { return Dial(server, acct1, lg.New("stream", "acct1")) },
			acct2, func() protocol.Conn { return Dial(server, acct2, lg.New("stream", "acct2")) },
		)
	})
}
