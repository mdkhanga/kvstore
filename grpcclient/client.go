package grpclient

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/mdkhanga/kvstore/kvmessages"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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
