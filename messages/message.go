package messages

import (
	"bytes"
	"encoding/binary"
)

type PingMessage struct {
	Type MessageType
	data *[]byte
}

type MessageType int16

type Message interface {
	Serialize() ([]byte, error)
}

const (
	HELLO    MessageType = 0
	RESPONSE MessageType = 1
	PING     MessageType = 2
)

func (message *PingMessage) serialize() ([]byte, error) {

	buf := new(bytes.Buffer)

	// Write the int32 value
	if err := binary.Write(buf, binary.LittleEndian, message.Type); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
