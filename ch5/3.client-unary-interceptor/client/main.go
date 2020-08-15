package main;

import (
    m "client/msg"
    "log"

    "google.golang.org/grpc"
    "github.com/golang/protobuf/ptypes/wrappers"
    "context"
)

func msgUnaryClientInterceptor(
    ctx context.Context,
    method string, 
    req interface{},
    reply interface{},
    cc *grpc.ClientConn,
    invoker grpc.UnaryInvoker,
    opts ...grpc.CallOption,
) error {
    log.Println("Method: " + method) //1. Preprocessing phase has access to the RPC request prior to sending it out to the server

    //Invoking the remote method
    err := invoker(ctx, method, req, reply, cc, opts...) //2. Invoking the RPC method via UnaryInvoker

    //Postprocessor phase
    log.Println(reply) //3. Postprocessing phase where you can process the response or error results

    return err //4. Returning an error back to the gRPC client application along with a reply, which is passed as an argument.
}

const address = "localhost:50051"

func main() {
    conn, err := grpc.Dial(
        address,
        grpc.WithInsecure(),
        grpc.WithUnaryInterceptor(msgUnaryClientInterceptor), //5. Setting up a connection to the server by passing a unary interceptor as a dial option
    )
    if err != nil {
        log.Fatalf("did cont connect: %v", err)
    }

    c := m.NewMsgServiceClient(conn)
    response, err := c.SendMsg(
        context.Background(),
        &wrappers.StringValue{Value: "Hi from client"},
    )
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Got the response %s", response.Value)
}
