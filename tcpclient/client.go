package tcpclient

import (
	"fmt"
	"net"
)

func connect(hostport string) (net.Conn, error) {

	fmt.Println("Connecting to " + hostport)

	return nil, nil

}
