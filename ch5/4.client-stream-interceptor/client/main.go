package main

import (
    m "client/msg"
    "context"
    "log"
    "google.golang.org/grpc"
    "github.com/golang/protobuf/ptypes/wrappers"
    "time"
)

func clientStreamInterceptor (
    ctx context.Context,
    desc *grpc.StreamDesc,
    cc *grpc.ClientConn,
    method string,
    streamer grpc.Streamer,
    opts ...grpc.CallOption,
) (grpc.ClientStream, error) {
    log.Println("========[Client Interceptor] ", method) //1. Preprocessing phase has access to the RPC request prior to sending it out to the server
    s, err := streamer(ctx, desc, cc, method, opts...) //2. Calling the passed-in streamer to get a ClientStream
    if err != nil {
        return nil, err
    }

    return newWrappedStream(s), nil //3. Wrapping around the ClientStream, overloading its methods with intercepting logic, and returning it to the client application
}

type wrappedStream struct { //4. Wrapper stream  of grpc.ClientStream
    grpc.ClientStream
}

func (w *wrappedStream) RecvMsg(m interface{}) error { //5. Function to intercept messages received from streaming RPC
    log.Printf(
        "======== [Client Stream Interceptor] Receive a message (Type: %T) at %v", 
        m,
        time.Now().Format(time.RFC3339),
    )
    return w.ClientStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) error { //6. Function to intercept messages sent from streaming RPC
    log.Printf(
        "======= [Client Stream Interceptor] Send a message (Type: %T) at %v",
        m,
        time.Now().Format(time.RFC3339),
    )
    return w.ClientStream.SendMsg(m)
}

func newWrappedStream(s grpc.ClientStream) grpc.ClientStream {
    return &wrappedStream{s}
}

const address = "localhost:50051"

func main() {
    conn, err := grpc.Dial(
        address,
        grpc.WithInsecure(),
        grpc.WithStreamInterceptor(clientStreamInterceptor), //7. Registering a streaming interceptor
    )
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()

    client := m.NewMsgServiceClient(conn)
    stream, err := client.SendMsg(context.Background()) 
    if err != nil {
        log.Fatalf("%v.SendMsg(_) = _, %v", client, err)
    }

    msg1 := wrappers.StringValue{Value: "Hello 1"}
    if err := stream.Send(&msg1); err != nil {
        log.Fatalf("%v.Send(%v) = %v", stream, msg1, err)
    }

    msg2 := wrappers.StringValue{Value: "Hello 2"}
    if err := stream.Send(&msg2); err != nil {
        log.Fatalf("%v.Send(%v) = %v", stream, msg2, err)
    }

    res, err := stream.CloseAndRecv()
    if err != nil {
        log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
    }

    log.Printf("Response: %s", res)
}
