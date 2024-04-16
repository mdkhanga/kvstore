package tcpclient

import (
	"fmt"
	"net"
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

	Connect(hostport)

}
