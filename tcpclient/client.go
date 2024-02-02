package tcpclient

import (
	"fmt"
	"net"
)

func Connect(hostport string) (net.Conn, error) {

	fmt.Println("Connecting to " + hostport)

	return nil, nil

}
