package main

import (
    pb "service/ecommerce"
    "log"
    "net"
    "google.golang.org/grpc"
    //"github.com/golang/protobuf/ptypes/wrappers"
    "io"
)

type Server struct {
    orderMap map[string]pb.Order
}

func NewServer() *Server {
    return &Server {
        orderMap: make(map[string]pb.Order),
    }
}

func (s *Server) ProcessOrders(stream pb.OrderManagement_ProcessOrdersServer) error {
    for {
        orderId, err := stream.Recv() //1. Read order IDs from the incoming stream
        if err == io.EOF { //2. Keep reading until the end of the stream is found
            return nil //4. Server-side end of the stream is marked by returning nil
        }
        if err != nil {
            return err
        }

        order := pb.Order{
            Id: orderId.Value,
            Items: []string{"a", "b", "c", "d"},
            Description: "description",
            Price: 4.4,
            Destination: "Italy",
        }

        comb := pb.CombinedShipment{
            Id: "1",
            Status: "Done",
            OrdersList: []*pb.Order{&order, &order, &order},
        }
        stream.Send(&comb)
    }
}

const port = ":50051"

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
