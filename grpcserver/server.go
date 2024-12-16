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
	"google.golang.org/grpc/codes"
	peer "google.golang.org/grpc/peer"
	status "google.golang.org/grpc/status"
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

	sendMessageQueue := &MessageQueue{}
	receiveMessageQueue := &MessageQueue{}

	stopChan := make(chan struct{})

	var once sync.Once

	// Function to safely close the stopChan
	closeStopChan := func() {
		once.Do(func() {
			close(stopChan)
		})
	}

	go receiveLoop(stream, receiveMessageQueue, stopChan, closeStopChan)

	go processMessageLoop(receiveMessageQueue, sendMessageQueue, stopChan, closeStopChan)

	go sendLoop(stream, sendMessageQueue, stopChan, closeStopChan)

	<-stopChan
	log.Println("Stopping message processing due to stream error")

	return nil

	// Block main goroutine to keep stream open
	// select {}
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

func receiveLoop(stream pb.KVSevice_CommunicateServer, messageQueue *MessageQueue, stopChan chan struct{}, closeStopChan func()) {

	ctx := stream.Context()

	for {

		select {

		case <-ctx.Done():
			log.Println("Client disconnected or context canceled (receiver)")
			// close(stopChan)
			closeStopChan()
			return

		default:
			in, err := stream.Recv()
			if err != nil {

				code := status.Code(err)

				if code == codes.Unavailable || code == codes.Canceled || code == codes.DeadlineExceeded {

					log.Println("Unable to read from the stream. server seems unavailable")
					close(stopChan)
					return
				}
			}

			log.Printf("Received message of type: %v", in.Type)
			if in.Type == pb.MessageType_PING {
				log.Printf("Received Ping message from the stream %d %s %d", in.GetPing().Hello, in.GetPing().Hostname, in.GetPing().Port)

				messageQueue.Enqueue(in)
				log.Printf("Server Queue length %d", len(messageQueue.messages))
			}

		}

	}

}

func sendLoop(stream pb.KVSevice_CommunicateServer, messageQueue *MessageQueue, stopChan chan struct{}, closeStopChan func()) {

	ctx := stream.Context()

	for {
		select {
		case <-ctx.Done(): // Client disconnected or context canceled
			log.Println("Client disconnected or context canceled (sender)")
			// close(stopChan)
			closeStopChan()
			return
		case <-stopChan: // Stop signal received
			log.Println("Stop signal received for sender goroutine")
			return
		default:
			// Send a message to the client (dummy example message)

			msg := messageQueue.Dequeue()
			if msg == nil {
				time.Sleep(1 * time.Second) // Wait before checking again
				continue
			}

			// Process each message type and decide what to send
			/* var response *pb.ServerMessage
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
			} */

			// Send the response if it was generated
			// if response != nil {
			if err := stream.Send(msg); err != nil {
				log.Printf("Error sending message: %v", err)

				closeStopChan()
				return
			}
			// }
		}
	}
}

func processMessageLoop(receiveMessageQueue *MessageQueue, sendMessageQueue *MessageQueue, stopChan chan struct{}, closeStopChan func()) {

	for {

		select {

		case <-stopChan:
			log.Println("Stop signal received for processing goroutine")
			return

		default:

			msg := receiveMessageQueue.Dequeue()
			if msg == nil {
				time.Sleep(1 * time.Second) // Wait before checking again
				continue
			}

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

			sendMessageQueue.Enqueue(response)

		}

	}

}
