package messages

import (
	"bytes"
	"encoding/binary"
)

type PingMessage struct {
	Type MessageType
}

type MessageType int16

type Message interface {
	Serialize() ([]byte, error)
	Deserialize([]byte) (Message, error)
}

const (
	HELLO    MessageType = 0
	RESPONSE MessageType = 1
	PING     MessageType = 2
)

func (message *PingMessage) Serialize() ([]byte, error) {

	buf := new(bytes.Buffer)

	// Write the int32 value
	if err := binary.Write(buf, binary.LittleEndian, message.Type); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (message *PingMessage) Deserialize(data []byte) error {

	buf := bytes.NewReader(data)

	if err := binary.Read(buf, binary.LittleEndian, message.Type); err != nil {
		return err
	}

	return nil
}
