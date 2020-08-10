package main

import (
    pb "service/ecommerce"
    "log"
    "net"
    "google.golang.org/grpc"
    "github.com/golang/protobuf/ptypes/wrappers"
    "io"
)

type Server struct {
    orderMap map[string]pb.Order
}

func NewServer() *Server {
    return &Server{
        orderMap: make(map[string]pb.Order),
    }
}

func (s *Server) UpdateOrders(stream pb.OrderManagement_UpdateOrdersServer) error {
    ordersStr := "Updated Order IDs : "
    for {
        order, err := stream.Recv() //1. Read message from the client stream
        if err == io.EOF { //2. Check for end of stream
            //Finished reading the order stream.
            return stream.SendAndClose(&wrappers.StringValue{Value: "Orders processed " + ordersStr})
        }

        //Update order
        s.orderMap[order.Id] = *order
        log.Printf("Order ID %v: Updated", order)
        ordersStr += order.Id + ", "
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
