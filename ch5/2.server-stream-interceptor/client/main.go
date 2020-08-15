package main;

import (
    "google.golang.org/grpc"
    "log"
    "context"

    "github.com/golang/protobuf/ptypes/wrappers"
    h "client/hello"
    "io"
)

const address = "localhost:50052"

func main() {
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()

    c := h.NewHelloManagementClient(conn)
    stream, err := c.SayHello(
        context.Background(),
        &wrappers.StringValue{Value: "Hello"},
    )
    if err != nil {
        log.Fatal(err)
    }

    if stream == nil {
        log.Fatal("stream is nil")
    }

    for {
        msg, err := stream.Recv()
        if err == io.EOF {
            break;
        }

        if err != nil {
            log.Fatal(err)
        }

        log.Println(msg.Value)
    }
}
