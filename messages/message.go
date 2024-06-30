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
	GetType() MessageType
	Serialize() ([]byte, error)
	Deserialize([]byte) (Message, error)
}

const (
	HELLO    MessageType = 0
	RESPONSE MessageType = 1
	PING     MessageType = 2
	UNKNOWN  MessageType = -1
)

func (message *PingMessage) GetType() MessageType {
	return message.Type
}

func (message *PingMessage) Serialize() ([]byte, error) {

	buf := new(bytes.Buffer)

	// Write the type value
	if err := binary.Write(buf, binary.LittleEndian, message.Type); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (message *PingMessage) Deserialize(data []byte) error {

	buf := bytes.NewReader(data)

	var typee int16

	if err := binary.Read(buf, binary.LittleEndian, &typee); err != nil {
		return err
	}

	if typee == 2 {
		message.Type = PING
	} else {
		message.Type = UNKNOWN
	}

	return nil
}

func Deserialize(data []byte) (Message, error) {

	return nil, nil
}
