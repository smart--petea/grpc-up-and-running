package main

import (
    "log"
    "net"

    "google.golang.org/grpc"
)

import (
    "context"

    "github.com/gofrs/uuid"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
    pb "productinfo/service/ecommerce" 
)

type Server struct { //2. The server struct is an abstraction of the server. it allows attaching service methods to the server
    productMap map[string]*pb.Product
}

func (s *Server) AddProduct(ctx context.Context, in *pb.Product) (*pb.ProductID, error) {
    //3. The AddProduct method takes Product as a parameter and returns a ProductID. Product and ProductID structs are
    //defined in the product_info.pb.go, which is autogenerate from the product_info.proto definition

    //5. Both methods also have a Context parameter. A context object contains metadatada such as the identity of the
    //end user authoriztion token and the requests's deadline, and it will exist during the lifetime of the request

    //6. Both methods return an aerror in addition to the return value of the remove method (methods have multiple
    //return types). These errors are propagated to the consumers and can be used for error handling at the consumer
    //side

    out, err := uuid.NewV4()
    if err != nil {
        return nil, status.Errorf(codes.Internal, "Error while generating Product ID", err)
    }
    in.Id = out.String()
    if s.productMap == nil {
        s.productMap = make(map[string]*pb.Product)
    }

    s.productMap[in.Id] = in
    return &pb.ProductID{Value: in.Id}, status.New(codes.OK, "").Err()
}

func (s *Server) GetProduct(ctx context.Context, in *pb.ProductID) (*pb.Product, error) {
    value, exists := s.productMap[in.Value]
    if exists {
        return value, status.New(codes.OK, "").Err()
    }

    return nil, status.Errorf(codes.NotFound, "Product does not exist.", in.Value)
}

const (
    port = ":50051"
)

func main() {
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer()
    pb.RegisterProductInfoServer(s, &Server{})

    log.Printf("Starting gRPC listener on port " + port)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
