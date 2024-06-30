package tcpclient

import (
	"fmt"
	"net"
	"time"

	"encoding/binary"

	"github.com/mdkhanga/kvstore/messages"
)

func Connect(hostport string) (net.Conn, error) {

	fmt.Println("Connecting to " + hostport)

	conn, err := net.Dial("tcp", hostport)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return conn, nil

}

func CallServer(hostport string) {

	conn, err := Connect(hostport)

	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}

	for true {

		/* message := "Hello, server!"
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		} */

		msg := messages.PingMessage{Type: messages.PING}
		data, err := msg.Serialize()

		// Calculate the length of the serialized data
		dataLength := len(data)

		fmt.Println("Length of data to write ", dataLength)

		// Write the length of the byte array to the socket
		if err := binary.Write(conn, binary.LittleEndian, int16(dataLength)); err != nil {
			fmt.Println("Error writing data length to socket:", err)
			return
		}

		// this might need be a loop. What is all data is not written
		n := dataLength

		for n > 0 {
			count, err := conn.Write(data)
			if err != nil {
				fmt.Println("Error writing data length to socket:", err)
			}
			fmt.Println("number of bytes written = ", count)
			n = n - count
		}

		// Receive and print the echoed message from the server
		buffer := make([]byte, 1024)
		n, err = conn.Read(buffer)
		if err != nil {
			fmt.Println("Error receiving message:", err)
			return
		}

		fmt.Printf("Received from server: %s\n", buffer[:n])

		time.Sleep(200 * time.Millisecond)
	}

}
