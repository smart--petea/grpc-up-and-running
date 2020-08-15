package main

import (
    "context"
    "github.com/golang/protobuf/ptypes/wrappers"

    h "service/hello"
    "log"
    "google.golang.org/grpc"
    "net"
)

type Server struct {}

func (server *Server) SayHello(context.Context, *wrappers.StringValue) (*wrappers.StringValue, error) {
    return &wrappers.StringValue{Value: "Hi from server"}, nil
}

const (
    port = ":50051"
)

func orderUnaryServerInterceptor(ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    //Preprocessing logic
    //Gets info about the current RPC call by examining the args passed in
    log.Println("======[Server Interceptor] ", info.FullMethod) //1. Preprocessing phase: this is where you can intercept the message prior to invoking the respective RPC

    //Invoking the handler to complete the normal execution of a unary RPC
    m, err := handler(ctx, req) //2. Invoking the RPC method via UnaryHandler

    //Post processing logic
    log.Printf(" Post Proc Message : %s", m) //3. Postprocessing phase: you can process the response from the RPC invocation

    return m, err //4. Sending back the RPC response
}

func main() {
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer(
        grpc.UnaryInterceptor(orderUnaryServerInterceptor), //5. Registering the unary interceptor with the gRPC server
    )
    h.RegisterHelloManagementServer(s, &Server{})

    log.Printf("Starting gRPC listener on port " + port)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
