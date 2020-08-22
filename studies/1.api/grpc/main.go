package main

import (
    "net"
    "log"
    "google.golang.org/grpc"
	empty "github.com/golang/protobuf/ptypes/empty"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
    "context"

    s "grpc/proto"
)

const (
    port = ":10000"
)

type Server struct {}

func (server *Server) HomePage(ctx context.Context, empt *empty.Empty) (*wrappers.StringValue, error) {
    return &wrappers.StringValue{Value: "Welcome to the HomePage!"}, nil
}

func main() {
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    gServer := grpc.NewServer()
    s.RegisterServiceAServer(gServer, &Server{}) 

    log.Printf("Starting gRPC listener on port " + port)
    if err := gServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
