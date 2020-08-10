package main

import (
    "google.golang.org/grpc"
    "log"
    "context"

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
    updateStream, err := client.UpdateOrders(context.Background()) //1. Invoking UpdateOrders remote method
    if err != nil { //2. Handling errors related to UpdateOrders
        log.Fatalf("%v.UpdatedOrders(_) = _, %v", client, err)
    }

    //Updating order 1
    updOrder1 := pb.Order{
        Id: "1",
        Items: []string{"a", "b", "c"},
        Description: "description order 1",
        Price: 1.1,
        Destination: "destination 1",
    }

    if err := updateStream.Send(&updOrder1); err != nil { //3. Sending order update via client stream
        log.Fatalf("%v.Send(%v) = %v", updateStream, updOrder1, err) //4. Handling errors when sending messages to stream
    }

    //Updating order 2
    updOrder2 := pb.Order{
        Id: "2",
        Items: []string{"a2", "b2", "c2"},
        Description: "description order 2",
        Price: 2.2,
        Destination: "destination 2",
    }
    if err := updateStream.Send(&updOrder2); err != nil {
        log.Fatalf("%v.Send(%v) = %v", updateStream, updOrder2, err)
    }

    //Updating order 3
    updOrder3 := pb.Order{
        Id: "3",
        Items: []string{"a3", "b3", "c3"},
        Description: "description order 3",
        Price: 3.3,
        Destination: "destination 3",
    }
    if err := updateStream.Send(&updOrder3); err != nil {
        log.Fatalf("%v.Send(%v) = %v", updateStream, updOrder3, err)
    }

    updateRes, err := updateStream.CloseAndRecv() //5. Closing the stream and receiving the response
    if err != nil {
        log.Fatalf("%v.CloseAndRecv() got error %v, want %v", updateStream, err, nil)
    }

    log.Printf("Update Orders Res: %s", updateRes)
}
