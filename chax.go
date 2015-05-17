package main

import (
	"polydawn.net/chax/chaxui"
	"polydawn.net/chax/protocol"
	"polydawn.net/chax/protocol/xmpp"

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

	ui := chaxui.New(theme)
	window.AddChild(ui.BaseLayout)

	go hello()
}

func hello() {
	account, serverDesc := testpartner()

	Log.Info("connecting", "account", account, "server", serverDesc)
	conn := xmpp.Dial(serverDesc, account, Log)
	Log.Info("connected!", "account", account, "server", serverDesc)

	// talk to ourselves
	// there's no such thing as message acknowledgement, apparently
	Log.Info("sending hello", "account", account, "server", serverDesc)
	conn.Send(account, []byte("hallomsg"))
	Log.Info("hello sent!", "account", account, "server", serverDesc)
}

func testpartner() (protocol.Account, protocol.ServerDesc) {
	switch 1 {
	case 1:
		// this one's fun!  it doesn't respond to... any... of my IQs, apparently.
		return protocol.Account{
				Username: "testpilot",
				Domain:   "crypt.mn",
				Password: "asdf",
			},
			xmpp.ResolveServer("crypt.mn")
	case 2:
		// this one's fun!  "PLAIN authentication is not an option"
		return protocol.Account{
				Username: "testpilot",
				Domain:   "im.koderoot.net",
				Password: "asdf",
			},
			xmpp.ResolveServer("im.koderoot.net")
	default:
		panic("silly")
	}
}
