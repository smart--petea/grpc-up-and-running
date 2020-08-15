package main

import (
    m "service/msg"
    "net"
    "log"
    "context"
    "github.com/golang/protobuf/ptypes/wrappers"
    "google.golang.org/grpc"
)

type Server struct {}

func (*Server) SendMsg(ctx context.Context, request *wrappers.StringValue) (*wrappers.StringValue, error) {
    return &wrappers.StringValue{Value: "Hi from server"}, nil
}

const port = ":50051"

func main() {
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer()
    m.RegisterMsgServiceServer(s, &Server{})

    log.Printf("Starting gRPC listener on port " + port)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
