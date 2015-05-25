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
			conv12 := conn1.StartConversation(acct2)
			msgBody := []byte("hallo")
			conv12.Send(msgBody)

			conv21 := <-conn2.AcceptConversations()
			ups := conv21.AwaitUpdates() // get this *before* checking or you haz a race
			msgs := conv21.GetMessages(0, -1)
			for len(msgs) < 1 {
				Printf("Snoozen\n")
				<-ups
				msgs = conv21.GetMessages(0, -1)
				ups = conv21.AwaitUpdates()
			}
			So(msgs[0].Body, ShouldResemble, msgBody)
		})
	})
}
