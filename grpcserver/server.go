package grpcserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/mdkhanga/kvstore/kvmessages"
	"google.golang.org/grpc"
	peer "google.golang.org/grpc/peer"
)

// server is used to implement helloworld.GreeterServer.
type Server struct {
	pb.UnimplementedKVSeviceServer
}

// SayHello implements helloworld.GreeterServer
func (s *Server) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	fmt.Println("Received", in.GetHello())
	peerInfo, ok := peer.FromContext(ctx)
	if ok {
		fmt.Println("Client Address", peerInfo.Addr.String())
		fmt.Println("Client Address", peerInfo.LocalAddr.String())
	}
	return &pb.PingResponse{Hello: in.GetHello()}, nil
}

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

func (s *Server) Communicate(stream pb.KVSevice_CommunicateServer) error {

	messageQueue := &MessageQueue{}

	go func() {
		for {
			in, err := stream.Recv()
			if err != nil {
				log.Printf("Error receiving message: %v", err)
				return
			}
			log.Printf("Received message of type: %v", in.Type)
			if in.Type == pb.MessageType_PING {
				fmt.Println("Received Ping message from the stream ", in.GetPing().Hello)
			}
			messageQueue.Enqueue(in)
		}
	}()

	// Goroutine to process the message queue and send responses as needed
	go func() {
		for {
			msg := messageQueue.Dequeue()
			if msg == nil {
				time.Sleep(1 * time.Second) // Wait before checking again
				continue
			}

			// Process each message type and decide what to send
			var response *pb.ServerMessage
			switch msg.Type {
			case pb.MessageType_PING:
				response = &pb.ServerMessage{
					Type: pb.MessageType_PING_RESPONSE,
					Content: &pb.ServerMessage_PingResponse{
						PingResponse: &pb.PingResponse{Hello: 2},
					},
				}
			case pb.MessageType_KEY_VALUE:
				log.Printf("Processing KeyValueMessage")
				// Handle KeyValueMessage
			default:
				log.Printf("Unknown message type received")
			}

			// Send the response if it was generated
			if response != nil {
				if err := stream.Send(response); err != nil {
					log.Printf("Error sending message: %v", err)
					return
				}
			}
		}
	}()

	// Block main goroutine to keep stream open
	select {}
}

func StartGrpcServer(portPtr *string) {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *portPtr))
	if err != nil {
		fmt.Println("failed to listen:", err)
	}

	s := grpc.NewServer()
	pb.RegisterKVSeviceServer(s, &Server{})
	fmt.Println("GRPC server listening at ", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
