package grpc

import (
	"context"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes"
	"github.com/initlevel5/microservices-example/app"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	svc app.ProductService
	*grpc.Server
	*log.Logger
}

func NewGrpcServer(svc app.ProductService, logger *log.Logger) *grpcServer {
	gs := grpc.NewServer()
	s := &grpcServer{svc: svc, Server: gs, Logger: logger}

	RegisterProductServiceServer(gs, s)
	reflection.Register(gs)

	return s
}

func (s *grpcServer) Serve(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		s.Printf("Serve(): Listen(): %v", err)
		return err
	}

	if err := s.Server.Serve(l); err != nil {
		s.Printf("Serve(): Serve(): %v", err)
		return err
	}
	return nil
}

func (s *grpcServer) CreateProduct(ctx context.Context, in *CreateProductRequest) (*CreateProductReply, error) {
	id, err := s.svc.CreateProduct(ctx, in.Title, in.Price, in.Manufacturer, in.Description)
	if err != nil {
		return &CreateProductReply{}, err
	}
	return &CreateProductReply{Id: id}, nil
}

func (s *grpcServer) DeleteProduct(ctx context.Context, in *DeleteProductRequest) (*DeleteProductReply, error) {
	id, err := s.svc.DeleteProduct(ctx, in.Id)
	if err != nil {
		return &DeleteProductReply{}, err
	}
	return &DeleteProductReply{Id: id}, nil
}

func (s *grpcServer) Product(ctx context.Context, in *ProductRequest) (*ProductReply, error) {
	p, err := s.svc.Product(ctx, in.Id)
	if err != nil {
		return &ProductReply{}, err
	}
	created, _ := ptypes.TimestampProto(p.Created)

	return &ProductReply{
		Id:           p.ID,
		Title:        p.Title,
		Price:        p.Price,
		Manufacturer: p.Manufacturer,
		Created:      created,
	}, nil
}

func (s *grpcServer) SearchProduct(ctx context.Context, in *SearchProductRequest) (*SearchProductReply, error) {
	id, err := s.svc.SearchProduct(ctx, in.Title)
	if err != nil {
		return &SearchProductReply{}, err
	}
	return &SearchProductReply{Id: id}, nil
}
