package tcpserver

import (
	"fmt"
	"net"
)

func Listen() {

	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on port 8081")

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

		fmt.Printf("Received: %s", buffer[:n])

		// Echo back to the client
		conn.Write(buffer[:n])
	}
}
