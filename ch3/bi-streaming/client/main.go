package main

import (
    "google.golang.org/grpc"
    "log"
    "context"
    "time"
    "io"
    "github.com/golang/protobuf/ptypes/wrappers"

    pb "client/ecommerce"
)

const address = "localhost:50051"

func main() {
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()

    client := pb.NewOrderManagementClient(conn)
    streamProcOrder, _ := client.ProcessOrders(context.Background()) //1. Invoke the remote method and obtain the stream reference for writing and reading from the client side
    if streamProcOrder == nil {
        log.Fatalf("streamProcOrder is nil");
    }

    if err := streamProcOrder.Send(&wrappers.StringValue{Value:"102"}); err != nil { //2. Send a message to the service
        log.Fatalf("%v.Send(%v) = %v", client, "102", err)
    }

    if err := streamProcOrder.Send(&wrappers.StringValue{Value:"103"}); err!= nil {
        log.Fatalf("%v.Send(%v) = %v", client, "103", err)
    }

    if err := streamProcOrder.Send(&wrappers.StringValue{Value:"104"}); err != nil {
        log.Fatalf("%v.Send(%v) = %v", client, "104", err)
    }

    channel := make(chan struct{}) //3. Create a channel to use for Goroutines
    go asncClientBidirectionalRPC(streamProcOrder, channel) //4. Invoke the function using Goroutines to read the messages in parallel from the service
    time.Sleep(time.Millisecond * 1000) //5. Mimic a delay when sending some messages to the service

    if err := streamProcOrder.Send(&wrappers.StringValue{Value: "101"}); err != nil {
        log.Fatalf("%v.Send(%v) = %v", client, "101", err)
    }

    if err := streamProcOrder.CloseSend(); err != nil { //6. Mark the end of stream for the client stream (order IDs)
        log.Fatal(err)
    }
    <-channel
}

func asncClientBidirectionalRPC (
    streamProcOrder pb.OrderManagement_ProcessOrdersClient,
    c chan struct{},
) {
    for {
        combinedShipment, errProcOrder := streamProcOrder.Recv() //7. Read service's messages on the client side
        if errProcOrder == io.EOF { //8. Condition to detect the end of the stream
            break
        }

        log.Printf("Combined shipment: %v", combinedShipment.OrdersList)
    }
    <-c
}
