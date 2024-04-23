package messages

type Message struct {
	Type MessageType
	data *[]byte
}

type MessageType int16

const (
	HELLO    MessageType = 0
	RESPONSE MessageType = 1
	PING     MessageType = 2
)

func createPingMessage() *[]byte {

}
