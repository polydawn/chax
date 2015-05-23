package tests

import (
	"polydawn.net/chax/protocol"

	. "github.com/smartystreets/goconvey/convey"
)

func CheckMessageBounce(
	acct1 protocol.Account, conn1dialer func() protocol.Conn,
	acct2 protocol.Account, conn2dialer func() protocol.Conn,
) {
	Convey("Dialing establishes a connection", func() {
		conn1 := conn1dialer()
		conn2 := conn2dialer()

		Convey("Sending a message from acct1 succeeds", func() {
			conn1.Send(acct2, []byte("hallo"))
			_ = conn2
		})
	})
}
