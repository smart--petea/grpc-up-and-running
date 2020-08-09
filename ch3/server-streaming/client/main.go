package main

import (
    "google.golang.org/grpc"
    "log"
    "context"

    "github.com/golang/protobuf/ptypes/wrappers"
    pb "client/ecommerce"
    "io"
    "fmt"
)

const address = "localhost:50051"

func main() {
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()

    c := pb.NewOrderManagementClient(conn)
    searchStream, _ := c.SearchOrders(context.Background(),
        &wrappers.StringValue{Value: "Google"}) //1. the SearchOrders function returns a client stream of OrderManagement_SearchOrdersClient, which has a Recv method

    if searchStream == nil {
        log.Fatalf("searchStream is nil")
    }


    for {
        searchOrder, err := searchStream.Recv() //2. Calling the client stream's Recv() method to retrieve Order response one by one
        if err == io.EOF { //3. When the end of the stream is found Recv returns an io.EOF
            break
        }

        fmt.Printf("got the order %v\n", searchOrder)
    }
}
