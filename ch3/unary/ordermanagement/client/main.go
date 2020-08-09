package main

import (
    "log"

    pb "ordermanagement/client/ecommerce"
    "context"
    "google.golang.org/grpc"
    "github.com/golang/protobuf/ptypes/wrappers"
    "os"
)

const address = "localhost:50051"

func main() {
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()

    c := pb.NewOrderManagementClient(conn)
    order, err := c.GetOrder(context.Background(), &wrappers.StringValue{Value: "106"})
    if err != nil {
        log.Fatalf("Could not get the order: %v", err)
    }
    log.Print("GetOrder response: ", order)

    order, err = c.GetOrder(context.Background(), &wrappers.StringValue{Value: "107"})
    if err != nil {
        log.Fatalf("Could not get the order: %v", err)
    }
    log.Print("GetOrder response: ", order)
}
