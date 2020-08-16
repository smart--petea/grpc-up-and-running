package main

import (
    m "service/msg"
    "log"
    "net"
    "google.golang.org/grpc"
    "github.com/golang/protobuf/ptypes/wrappers"
    "io"
)

type Server struct {}


func (*Server) SendMsg(stream m.MsgService_SendMsgServer) error {
    response := ""
    for {
        msg, err := stream.Recv()
        if err == io.EOF {
            return stream.SendAndClose(&wrappers.StringValue{Value: response})
        }

        response = response + " " + msg.Value
        log.Println(msg.Value)
    }

    return nil
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
