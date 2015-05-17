package protocol

// we wrap most things in blocking interfaces, even if they're using state and ack-based
// protocols in the backend (like xmpp's IQs).

type Conn interface {
	Send(recepient Account, message Message)
}
