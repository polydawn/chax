package protocol

// we wrap most things in blocking interfaces, even if they're using state and ack-based
// protocols in the backend (like xmpp's IQs).

type Conn interface {
	StartConversation(recipient Account) Conversation
	AcceptConversations() <-chan Conversation // Handle new incoming conversations.
}

/*
	A `Conversation` is a glorified list of messages, gathered by
	virtue of being addressed to the same recipients.

	The recipients may be a single user or may be a group (sometimes aka MUC).

	The list of messages may also contain metadata events, such as join
	and part events.
*/
type Conversation interface {
	Send(msg []byte) Message
	GetMessages(start, end int) []Message // 1 is the beginning.  '-1' for end will get "the rest".
	UpdateSignaller() <-chan struct{}     // Returned chan will be closed after new messages are logged.
}
