package main;

import (
    h "service/hello"
    "log"
    "net"
    "google.golang.org/grpc"
    "github.com/golang/protobuf/ptypes/wrappers"
    "time"
)

type Server struct {//1. Wrapper stram of the grpc.ServerStream
    grpc.ServerStream
}

func (w *Server) RecvMsg(m interface{}) error { //2. Implementing the RecvMsg function of the wrapper to processess messages received with stream RPC
    log.Printf("======[Server Stream Interceptor Wrapper] Receive a message (Type: %T) at %s", m, time.Now().Format(time.RFC3339))
    return w.ServerStream.RecvMsg(m)
}

func (w *Server) SendMsg(m interface{}) error { //3. Implementing the SendMsg function of the wrapper to process messages sent with stream RPC
    log.Printf("=====[Server Stream Interceptor Wrapper] Send a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
    return w.ServerStream.SendMsg(m)
}

func newServer(s grpc.ServerStream) grpc.ServerStream { //4. Creating an instance of the new wrapper stream
    return &Server{s}
}

func (*Server) SayHello(msg *wrappers.StringValue, stream h.HelloManagement_SayHelloServer) error {
    log.Printf("Got %s", msg.Value)

    stream.Send(&wrappers.StringValue{Value: "Nice"})
    stream.Send(&wrappers.StringValue{Value: "to"})
    stream.Send(&wrappers.StringValue{Value: "meet"})
    stream.Send(&wrappers.StringValue{Value: "yout"})

    return nil
}

func helloServerStreamInterceptor( //5. Streamin interprocessor implementation
    srv interface{},
    ss grpc.ServerStream,
    info *grpc.StreamServerInfo,
    handler grpc.StreamHandler,
) error {
    log.Println("=======[Server Stream Interceptor] ", info.FullMethod) //6. Preprocessor phase
    err := handler(srv, newServer(ss)) //7.Invoking the streaming RPC with the wrapper stream.
    if err != nil {
        log.Printf("RPC failed with error %v", err)
    }
    return err
}

const port = ":50052"

func main() {
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer(
        grpc.StreamInterceptor(helloServerStreamInterceptor),
    )
    h.RegisterHelloManagementServer(s, &Server{})

    log.Printf("Starting gRPC listener on port %s",  port)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
