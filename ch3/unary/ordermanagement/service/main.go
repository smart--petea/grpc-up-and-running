package main

import (
    pb "ordermanagement/service/ecommerce"
    "log"
    "net"
    "google.golang.org/grpc"
    "context"
    "github.com/golang/protobuf/ptypes/wrappers"
    "errors"
)

type Server struct {
    orderMap map[string]pb.Order
}

func NewServer() *Server {
    var server Server
    server.orderMap = make(map[string]pb.Order)
    server.orderMap["106"] = pb.Order{
        Id: "106",
        Items: []string{"a item", "b item", "c item"},
        Description: "Very nice product",
        Price: 5.66,
        Destination: "Abroad",
    }

    return &server
}

func (server *Server) GetOrder(ctxt context.Context, orderId *wrappers.StringValue) (*pb.Order, error) {
    ord, exist := server.orderMap[(*orderId).Value]
    if !exist {
        return nil, errors.New("Order does not exist")
    }
    return &ord, nil
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
    pb.RegisterOrderManagementServer(s, NewServer())

    log.Printf("Starting gRPC listener on port " + port)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
