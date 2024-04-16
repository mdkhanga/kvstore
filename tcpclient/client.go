package tcpclient

import (
	"fmt"
	"net"
)

func Connect(hostport string) (net.Conn, error) {

	fmt.Println("Connecting to " + hostport)

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return conn, nil

}
