package main

import (
    pb "service/ecommerce"
    "log"
    "net"
    "google.golang.org/grpc"
    "github.com/golang/protobuf/ptypes/wrappers"
)

type Server struct {}

func (*Server) SearchOrders(searchQuery *wrappers.StringValue, stream pb.OrderManagement_SearchOrdersServer) error {
    order := pb.Order{
        Id: "106",
        Items: []string{"a item", "b item", "c item"},
        Description: "Very nice product",
        Price: 5.66,
        Destination: "Abroad",
    }
    stream.Send(&order)

    order = pb.Order{
        Id: "106",
        Items: []string{"a item", "b item", "c item"},
        Description: "Very nice product",
        Price: 5.66,
        Destination: "Abroad",
    }
    stream.Send(&order)

    order = pb.Order{
        Id: "106",
        Items: []string{"a item", "b item", "c item"},
        Description: "Very nice product",
        Price: 5.66,
        Destination: "Abroad",
    }
    stream.Send(&order)
    return nil
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
    pb.RegisterOrderManagementServer(s, &Server{})

    log.Printf("Starting gRPC listener on port " + port)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
