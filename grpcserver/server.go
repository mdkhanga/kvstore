package grpcserver

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/mdkhanga/kvstore/kvmessages"
	"google.golang.org/grpc"
)

// server is used to implement helloworld.GreeterServer.
type Server struct {
	pb.UnimplementedKVSeviceServer
}

// SayHello implements helloworld.GreeterServer
func (s *Server) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	fmt.Println("Received: %v", in.GetHello())
	return &pb.PingResponse{Hello: in.GetHello()}, nil
}

func StartGrpcServer(portPtr *string) {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *portPtr))
	if err != nil {
		fmt.Println("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterKVSeviceServer(s, &Server{})
	fmt.Println("GRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
