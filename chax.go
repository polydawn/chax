package main

import (
	"fmt"
	"net"
	"time"

	"polydawn.net/chax/chaxui"
	"polydawn.net/chax/protocol"
	xmpp2 "polydawn.net/chax/protocol/xmpp" // todo get agl's code outta sight of main and mostly hidden behind protocol agnosticness

	"github.com/agl/xmpp"
	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/math"
	"github.com/google/gxui/themes/dark"
	"github.com/inconshreveable/log15"
)

var Log = log15.New()

func main() {
	gl.StartDriver(appMain)
}

func getTitle() string {
	// TODO(lfaraone): include username
	return "Chax"
}
func appMain(driver gxui.Driver) {
	theme := dark.CreateTheme(driver)

	window := theme.CreateWindow(800, 600, getTitle())
	window.OnClose(driver.Terminate)
	window.SetPadding(math.Spacing{L: 10, T: 10, R: 10, B: 10})

	window.AddChild(chaxui.NewChatWindowContents(theme))

	go hello()
}

func hello() {
	// this one's fun!  it doesn't respond to... any... of my IQs, apparently.
	account := protocol.Account{
		Username: "testpilot",
		Domain:   "crypt.mn",
		Password: "asdf",
	}
	serverDesc := xmpp2.ResolveServer("crypt.mn")

	// this one's fun!  "PLAIN authentication is not an option"
	//	account := protocol.Account{
	//		Username: "testpilot",
	//		Domain:   "im.koderoot.net",
	//		Password: "asdf",
	//	}
	//	serverDesc := xmpp2.ResolveServer("im.koderoot.net")

	addr := fmt.Sprintf("%s:%d", serverDesc.Host, serverDesc.Port)

	xmppConfig := &xmpp.Config{
		InLog:          LogWriter{Log, log15.LvlDebug, " <- RECV <-"},
		OutLog:         LogWriter{Log, log15.LvlDebug, " -> SENT ->"},
		Log:            LogWriter{Log, log15.LvlDebug, " :: NOTE ::"},
		TrustedAddress: true, // current test account has cert for 'www.xmpp.pro', which is not our account domain; handle better later
		//ServerCertificateSHA256: certSHA256, // pinning todo
		TLSConfig: TLSConfig,
		SkipTLS:   true, // sigh, current test server is EOF'ing me.  presumably suite mismatches; no indications; libpurple is fine with it.
	}

	Log.Info("connecting", "account", account, "server", serverDesc)

	// do our own dial, because the default timeouts are... insane.  like, minutes.  plural.
	sock, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		panic("Failed to connect to XMPP server: " + err.Error())
	}
	xmppConfig.Conn = sock

	// shake
	conn, err := xmpp.Dial(addr, account.Username, account.Domain, account.Password, xmppConfig)
	if err != nil {
		panic("Failed to connect to XMPP server: " + err.Error())
	}

	// announce that we're alive
	conn.SignalPresence("")

	// talk to ourselves
	// there's no such thing as message acknowledgement, apparently
	conn.Send(account.FQAN(), string("hallomsg"))

	// fire an IQ.  we'll try to collect on it later.
	rosterReply, _, err := conn.RequestRoster()
	if err != nil {
		panic("Failed to request roster: " + err.Error())
	}

	// fire an IQ.  wait for it, bceause in theory there's no earthy reason you wouldn't get a version response immediately, right?
	//	replyChan, _, err := conn.SendIQ(account.FQAN(), "get", xmpp.VersionQuery{})
	//	if err != nil {
	//		panic("Error sending version request: " + err.Error())
	//	}
	//	awaitVersionReply(replyChan, account.FQAN())

	heartbeatTicker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-heartbeatTicker.C:
			conn.SignalPresence("keks")
		case rosterStanza, ok := <-rosterReply:
			if !ok {
				panic("Failed to read roster: " + err.Error())
			}
			var roster []xmpp.RosterEntry
			if roster, err = xmpp.ParseRoster(rosterStanza); err != nil {
				panic("Failed to parse roster: " + err.Error())
			}
			Log.Info("Roster received", "roster", roster)
		}
	}
}
