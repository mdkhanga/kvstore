package tcpserver

import (
	"bytes"
	"fmt"
	"net"

	"encoding/binary"
)

func Listen(port string) {

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on port " + port)

	for {
		// Accept a connection and handle it in a Goroutine
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Handle the incoming connection
	fmt.Println("Accepted connection from", conn.RemoteAddr())

	// Implement your logic for handling the connection here

	// For example, you can echo messages back to the client
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}

		var msglength int32
		binary.Read(bytes.NewReader(buffer[:4]), binary.BigEndian, &msglength)

		dataBytes := byteArray[4 : 4+length]

		fmtString := fmt.Sprintf("Received: %s", buffer[:n])
		fmt.Println(fmtString)

		// Echo back to the client
		conn.Write(buffer[:n])
	}
}
