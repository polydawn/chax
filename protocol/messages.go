package protocol

import (
	"time"
)

type Message struct {
	ID           int64  // Local ID.  Used as a key.
	Uid          string // Globally shared ID, if exists.  May be used for message acknowledgement.
	Author       Account
	TimeSender   time.Time // Sender's claim of their time (if the protocol supports it).
	TimeReceived time.Time // Local time when we received the message.
	Body         []byte
	// other?
	//  'recipients' is probably better managed by whatever contains pointers to this / its uid
	//    unless we go row-style in our data storage, which is an option
	// might make sense to have an recievedBy field for our own client ID... e.g. so you could merge logs see by two different JID/POP
}
