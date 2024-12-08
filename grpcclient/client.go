package grpclient

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	pb "github.com/mdkhanga/kvstore/kvmessages"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	status "google.golang.org/grpc/status"
)

// Queue to hold incoming messages
type MessageQueue struct {
	messages []*pb.ServerMessage
	mu       sync.Mutex
}

// Enqueue adds a message to the queue
func (q *MessageQueue) Enqueue(msg *pb.ServerMessage) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.messages = append(q.messages, msg)
}

// Dequeue removes and returns the oldest message from the queue
func (q *MessageQueue) Dequeue() *pb.ServerMessage {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.messages) == 0 {
		return nil
	}
	msg := q.messages[0]
	q.messages = q.messages[1:]
	return msg
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
		return
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

func CallGrpcServerv2(hostport string) {

	for {

		fmt.Println(" Calling grpc server")

		conn, err := grpc.NewClient(hostport, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("did not connect: %v", err)
			log.Println("Sleep for 5 sec and try again")
			time.Sleep(5 * time.Second)
			continue
		}
		defer conn.Close()

		c := pb.NewKVSeviceClient(conn)
		ctx := context.Background()
		// defer cancel()

		fmt.Println("Create KVclient")

		stream, err := c.Communicate(ctx)
		if err != nil {
			fmt.Println("Error getting bidirectinal strem")
			conn.Close()
			log.Println("Sleep for 5 sec and try again")
			time.Sleep(5 * time.Second)
			continue

		}

		sendMessageQueue := &MessageQueue{}
		receiveMessageQueue := &MessageQueue{}

		stopChan := make(chan struct{})

		go sendLoop(stream, sendMessageQueue, stopChan)

		go receiveLoop(stream, receiveMessageQueue, stopChan)

		<-stopChan
		log.Println("Stopping message processing due to stream error")
		stream.CloseSend()
		conn.Close()
		log.Println("Sleep for 5 sec and try again")
		time.Sleep(5 * time.Second)

	}

	/* for true {

		fmt.Println("Sending ping")

		sendMessageQueue.Enqueue(&pb.ServerMessage{
			Type: pb.MessageType_PING,
			Content: &pb.ServerMessage_Ping{
				Ping: &pb.PingRequest{Hello: 1},
			},
		})

		log.Printf("Server Queue length %d", len(sendMessageQueue.messages))

		time.Sleep(5000 * time.Millisecond)

	} */

}

func sendLoop(stream pb.KVSevice_CommunicateClient, messageQueue *MessageQueue, stopChan chan struct{}) {

	for {

		select {

		case <-stopChan:
			log.Println("Stopping send goroutine ..")
			return

		default:
			/* msg := messageQueue.Dequeue()
			if msg == nil {
				time.Sleep(1 * time.Second) // Wait before checking again
				continue
			} */

			msg := &pb.ServerMessage{
				Type: pb.MessageType_PING,
				Content: &pb.ServerMessage_Ping{
					Ping: &pb.PingRequest{Hello: 1},
				},
			}

			log.Printf("Dequed Sending message of type: %v", msg.Type)
			err := stream.Send(msg)
			if err != nil {
				log.Printf("Error sending message: %v", err)
				close(stopChan)
				return
			}

		}

		time.Sleep(5 * time.Second)

	}
}

func receiveLoop(stream pb.KVSevice_CommunicateClient, messageQueue *MessageQueue, stopChan chan struct{}) {

	for {
		msg, err := stream.Recv()
		if err != nil {

			code := status.Code(err)

			if code == codes.Unavailable || code == codes.Canceled || code == codes.DeadlineExceeded {

				log.Println("Unable to read from the stream. server seems unavailable")
				close(stopChan)
				return

			}

		}
		log.Printf("Received message of type: %v", msg.Type)

		if msg.Type == pb.MessageType_PING_RESPONSE {
			fmt.Println("Received Ping message from the stream ", msg.GetPingResponse().Hello)
		}

		// For now do nothing with the msg
		// messageQueue.Enqueue(msg)
	}

}

/*
func CallGrpcServerv3(hostport string) {

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

	stream, err := c.Communicate(ctx)
	if err != nil {
		fmt.Println("Error getting bidirectinal strem")
	}

	sendMessageQueue := &MessageQueue{}
	receiveMessageQueue := &MessageQueue{}

	go sendLoop(stream, sendMessageQueue)

	go receiveLoop(stream, receiveMessageQueue)

	for true {

		fmt.Println("Sending ping")

		sendMessageQueue.Enqueue(&pb.ServerMessage{
			Type: pb.MessageType_PING,
			Content: &pb.ServerMessage_Ping{
				Ping: &pb.PingRequest{Hello: 1},
			},
		})

		log.Printf("Server Queue length %d", len(sendMessageQueue.messages))

		time.Sleep(5000 * time.Millisecond)

	}

}
*/
