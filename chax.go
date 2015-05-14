package main

import (
	"fmt"
	"net"
	"time"

	"polydawn.net/chax/chaxui"

	"github.com/agl/xmpp"
	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/math"
	"github.com/google/gxui/themes/dark"
	"github.com/inconshreveable/log15"
)

func main() {
	gl.StartDriver(appMain)
}

var Log = log15.New()

func appMain(driver gxui.Driver) {
	theme := dark.CreateTheme(driver)

	window := theme.CreateWindow(800, 600, "Hallo mundai")
	window.OnClose(driver.Terminate)
	window.SetPadding(math.Spacing{L: 10, T: 10, R: 10, B: 10})

	splitter := theme.CreateSplitterLayout()
	splitter.SetOrientation(gxui.Horizontal)
	window.AddChild(splitter)

	idkPlaceholder := theme.CreateTextBox()
	splitter.AddChild(idkPlaceholder)

	consoleLogCtrl := &chaxui.LinesAdapter{}
	consoleLog := theme.CreateList()
	consoleLog.SetAdapter(consoleLogCtrl)
	splitter.AddChild(consoleLog)

	consoleLogCtrl.Set([]string{"hai", "hay"})

	//hello()
}

func hello() {
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
		SkipTLS:   true, // sigh, current test server is EOF'ing me.  presumably suite mismatches; no indications; libpurple is fine with it.
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

	time.Sleep(5 * time.Minute)

	conn.SignalPresence("")
}
