package main

import (
    "net"
    "log"
    "google.golang.org/grpc"
	empty "github.com/golang/protobuf/ptypes/empty"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
    "context"

    s "grpc/proto"
)

const (
    port = ":10000"
)

type Server struct {
    Articles []s.Article
}

func (server *Server) ReturnAllArticles(emt *empty.Empty, stream s.ServiceA_ReturnAllArticlesServer) error {
    log.Println("ReturnAllArticles function")
    for _, article := range server.Articles {
        log.Printf("\nsending article %v", article)
        err := stream.Send(&article)
        if err != nil {
            log.Printf("error %w", err)
        }
    }

    return nil
}

func (server *Server) HomePage(ctx context.Context, empt *empty.Empty) (*wrappers.StringValue, error) {
    return &wrappers.StringValue{Value: "Welcome to the HomePage!"}, nil
}

func NewServer() *Server {
    articles := []s.Article {
        s.Article{
            Title: "Hello",
            Desc: "Article Description",
            Content: "Article Content",
        },
        s.Article{
            Title: "Hello 1",
            Desc: "Article Description 1",
            Content: "Article Content 1",
        },
    }

    return &Server{
        Articles: articles,
    }
}

func main() {
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    gServer := grpc.NewServer()
    s.RegisterServiceAServer(gServer, NewServer()) 

    log.Printf("Starting gRPC listener on port " + port)
    if err := gServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
