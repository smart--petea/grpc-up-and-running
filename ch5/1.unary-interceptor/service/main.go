package main

import (
    "context"
    "github.com/golang/protobuf/ptypes/wrappers"

    h "service/hello"
    "log"
    "google.golang.org/grpc"
    "net"
)

type Server struct {}

func (server *Server) SayHello(context.Context, *wrappers.StringValue) (*wrappers.StringValue, error) {
    return &wrappers.StringValue{Value: "Hi from server"}, nil
}

const (
    port = ":50051"
)

func main() {
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer()
    h.RegisterHelloManagementServer(s, &Server{})

    log.Printf("Starting gRPC listener on port " + port)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
