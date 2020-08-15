package main

import (
    h "client/hello"
    "log"

    "google.golang.org/grpc"
    "github.com/golang/protobuf/ptypes/wrappers"
    "context"
)

const address = "localhost:50051"

func main() {
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }

    c := h.NewHelloManagementClient(conn)
    response, err := c.SayHello(context.Background(), &wrappers.StringValue{Value: "Client hello"})
    if err != nil {
        log.Fatal("Could not get the response to my hello: ")
    }

    log.Printf("Got the response %s", response.Value)
}
