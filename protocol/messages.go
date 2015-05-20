package protocol

import (
	"time"
)

type Message struct {
	Uid          string
	Author       Account
	TimeSender   time.Time // Sender's claim of their time (if the protocol supports it).
	TimeReceived time.Time // Local time when we received the message.
	Body         []byte
	// other?
	//  'recipients' is probably better managed by whatever contains pointers to this / its uid
}
