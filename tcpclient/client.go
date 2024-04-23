package tcpclient

import (
	"fmt"
	"net"
	"time"
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

	for true {

		message := "Hello, server!"
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}

		// Receive and print the echoed message from the server
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error receiving message:", err)
			return
		}

		fmt.Printf("Received from server: %s\n", buffer[:n])

		time.Sleep(200 * time.Millisecond)
	}

}
