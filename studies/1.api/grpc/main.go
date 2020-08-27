package main

import (
    "net"
    "log"
    "google.golang.org/grpc"
    "google.golang.org/grpc/status"
    "google.golang.org/grpc/codes"
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

func (server *Server) CreateNewArticle(ctx context.Context, article *s.Article) (*s.Article, error) {
    log.Printf("CreateNewArticle function arg=%v", article)

    existentArticle := server.GetArticleById(article.Id)
    if existentArticle != nil {
        return nil, status.Error(codes.AlreadyExists, "There already exists such an article")
    }

    server.Articles = append(server.Articles, *article)

    return article, nil
}

func (server *Server) DeleteArticle(ctx context.Context, id *wrappers.StringValue) (*s.Article, error) {
    log.Printf("DeleteArticle function arg=%v", id)

    article := server.GetArticleById(id.Value)
    if article == nil {
        return nil, status.Error(codes.NotFound, "The article is not found")
    }

    server.RemoveArticle(article)
    return article, nil
}

func (server *Server) RemoveArticle(article *s.Article) {
    if len(server.Articles) == 0 {
        return
    }

    var index int
    var art s.Article
    for index, art = range server.Articles {
        if art.Id == article.Id {
            break
        }
    }

    if index >= len(server.Articles) {
        return
    }

    server.Articles = append(server.Articles[:index], server.Articles[index + 1:]...)
}

func (server *Server) GetArticleById(id string) *s.Article {
    for _, article := range server.Articles {
        if article.Id == id {
            return &article
        }
    }
    return nil
}

func (server *Server) ReturnSingleArticle(ctx context.Context, id *wrappers.StringValue) (*s.Article, error) {
    log.Printf("ReturnSingleArticle function arg=%v", id)

    article := server.GetArticleById(id.Value)
    if article == nil {
        return nil, status.Error(codes.NotFound, "the article is not found")
    }

    return article, nil
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
            Id: "0",
            Title: "Hello",
            Desc: "Article Description",
            Content: "Article Content",
        },
        s.Article{
            Id: "1",
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
