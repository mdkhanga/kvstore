package messages

type Message struct {
	id        int
	operation int
	data      *[]byte
}

type MessageType int16

const (
	HELLO    MessageType = 0
	RESPONSE MessageType = 1
	PING     MessageType = 2
)
