package grpc

import (
	"log"

	"google.golang.org/grpc"
)

type GrpcClient struct {
	ProductServiceClient
	*grpc.ClientConn
	*log.Logger
}

func NewGrpcClient(addr string, logger *log.Logger) (*GrpcClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		logger.Printf("NewClient(): Dial(): %v", err)
		return nil, err
	}

	return &GrpcClient{
		ProductServiceClient: NewProductServiceClient(conn),
		ClientConn:           conn,
		Logger:               logger,
	}, nil
}

func (c *GrpcClient) Close() {
	c.ClientConn.Close()
}
