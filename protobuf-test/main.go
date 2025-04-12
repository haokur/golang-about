package main

import (
  "context"
  "log"
  "net"

  pb "protobuf_test/hello"

  "google.golang.org/grpc"
)

type server struct {
  pb.UnimplementedHelloServiceServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
  log.Printf("Received request from: %v", in.GetName())
  return &pb.HelloResponse{Message: "Hello, " + in.GetName() + "!"}, nil
}

func main() {
  lis, err := net.Listen("tcp", ":50051")
  if err != nil {
    log.Fatalf("failed to listen: %v", err)
  }

  s := grpc.NewServer()
  pb.RegisterHelloServiceServer(s, &server{})

  log.Println("gRPC server listening at :50051")
  if err := s.Serve(lis); err != nil {
    log.Fatalf("failed to serve: %v", err)
  }
}
