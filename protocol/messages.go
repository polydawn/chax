package protocol

import (
	"time"
)

type Message struct {
	Uid    string
	Author Account
	Time   time.Time
	Body   string
	// other?
	//  'recipients' is probably better managed by whatever contains pointers to this / its uid
}
