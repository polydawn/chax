package main

import (
	"fmt"
	"net"
	"time"

	"github.com/agl/xmpp"
	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/themes/dark"
	"github.com/inconshreveable/log15"
)

func main() {
	gl.StartDriver(appMain)
}

var Log = log15.New()

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

func appMain(driver gxui.Driver) {
	theme := dark.CreateTheme(driver)

	window := theme.CreateWindow(800, 600, "Hallo mundai")
	window.OnClose(driver.Terminate)

	account := Account{
		Username: "testpilot",
		Domain:   "crypt.mn",
		Password: "asdf",
	}
	serverDesc := Resolve("crypt.mn")
	addr := fmt.Sprintf("%s:%d", serverDesc.Host, serverDesc.Port)

	xmppConfig := &xmpp.Config{
		InLog:          LogWriter{log15.Root(), log15.LvlDebug, " <- RECV <-"},
		OutLog:         LogWriter{log15.Root(), log15.LvlDebug, " -> SENT ->"},
		Log:            LogWriter{log15.Root(), log15.LvlDebug, " :: NOTE ::"},
		TrustedAddress: true, // current test account has cert for 'www.xmpp.pro', which is not our account domain; handle better later
		//ServerCertificateSHA256: certSHA256, // pinning todo
		TLSConfig: TLSConfig,
		SkipTLS: true, // sigh, current test server is EOF'ing me.  presumably suite mismatches; no indications; libpurple is fine with it.
	}

	Log.Info("connecting", "account", account, "server", serverDesc)

	// do our own dial, because the default timeouts are... insane.  like, minutes.  plural.
	sock, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		panic("Failed to connect to XMPP server: " + err.Error())
	}
	xmppConfig.Conn = sock

	// shake
	conn, err := xmpp.Dial(addr, account.Username, account.Domain, account.Password, xmppConfig)
	if err != nil {
		panic("Failed to connect to XMPP server: " + err.Error())
	}

	//	rosterReply, _, err := conn.RequestRoster()
	//	if err != nil {
	//		panic("Failed to request roster: "+err.Error())
	//		return
	//	}

	conn.SignalPresence("")
}
