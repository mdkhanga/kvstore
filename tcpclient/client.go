package tcpclient

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"encoding/binary"

	"github.com/mdkhanga/kvstore/messages"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/mdkhanga/kvstore/kvmessages"
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
			n = n - count
		}

		// Receive and print the echoed message from the server
		buffer := make([]byte, 1024)
		n, err = conn.Read(buffer)
		if err != nil {
			fmt.Println("Error receiving message:", err)
			return
		}

		// fmt.Printf("Received from server: %s\n", buffer[:n])

		time.Sleep(1000 * time.Millisecond)
	}

}

func CallGrpcServer(hostport string) {

	fmt.Println(" Calling grpc server")

	conn, err := grpc.NewClient(hostport, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewKVSeviceClient(conn)
	ctx := context.Background()
	// defer cancel()

	fmt.Println("Create KVclient")

	var resp *pb.PingResponse
	resp, err = c.Ping(ctx, &pb.PingRequest{Hello: 1})

	if err != nil {
		fmt.Println("got error on ping: %v", err)
	}

	fmt.Println("called ping")

	if resp.Hello == 1 {
		fmt.Println("Get Ping Response")
	}

	stream, err := c.Communicate(ctx)
	if err != nil {
		fmt.Println("Error getting bidirectinal strem")
	}

	for true {

		fmt.Println("Sending ping")
		// resp, err = c.Ping(ctx, &pb.PingRequest{Hello: 1})
		err := stream.Send(&pb.ServerMessage{
			Type: pb.MessageType_PING,
			Content: &pb.ServerMessage_Ping{
				Ping: &pb.PingRequest{Hello: 1},
			},
		})
		if err != nil {
			fmt.Println("Error sending Ping message: %v", err)
			return
		}

		/* if err != nil {
			fmt.Println("got error on ping: %v", err)
		} else {

			fmt.Println("Got PingResponse", resp.Hello)
		} */

		in, err := stream.Recv()
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			return
		}
		log.Printf("Received message of type: %v", in.Type)
		if in.Type == pb.MessageType_PING_RESPONSE {
			fmt.Println("Received Ping message from the stream ", in.GetPingResponse().Hello)
		}

		time.Sleep(1000 * time.Millisecond)

	}

}
