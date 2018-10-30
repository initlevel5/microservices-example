package graphql

import (
	"context"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/graph-gophers/graphql-go"

	"github.com/initlevel5/microservices-example/app"
	agrpc "github.com/initlevel5/microservices-example/app/grpc"
)

type resolver struct {
	c *agrpc.GrpcClient
	*log.Logger
}

func NewResolver(c *agrpc.GrpcClient, logger *log.Logger) *resolver {
	return &resolver{c: c, Logger: logger}
}

func (r *resolver) Product(args struct{ ID graphql.ID }) *productResolver {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	resp, err := r.c.Product(ctx, &agrpc.ProductRequest{Id: string(args.ID)})
	if err != nil {
		r.Printf("Product(): %v", err)
		return nil
	}

	created, _ := ptypes.Timestamp(resp.Created)

	return &productResolver{
		&app.Product{
			ID:           resp.Id,
			Title:        resp.Title,
			Price:        resp.Price,
			Manufacturer: resp.Manufacturer,
			Description:  resp.Description,
			Created:      created,
		},
	}
}

func (r *resolver) SearchProduct(args struct{ Title string }) *graphql.ID {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	resp, err := r.c.SearchProduct(ctx, &agrpc.SearchProductRequest{Title: args.Title})
	if err != nil {
		r.Printf("SearchProduct(): %v", err)
		return nil
	}

	id := new(graphql.ID)
	*id = graphql.ID(resp.Id)
	return id
}

func (r *resolver) CreateProduct(args struct {
	Title        string
	Price        float64
	Manufacturer string
	Description  *string
}) *graphql.ID {

	var descr string

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if args.Description != nil {
		descr = *args.Description
	}

	resp, err := r.c.CreateProduct(ctx, &agrpc.CreateProductRequest{
		Title:        args.Title,
		Price:        args.Price,
		Manufacturer: args.Manufacturer,
		Description:  descr,
	})
	if err != nil {
		r.Printf("CreateProduct(): %v", err)
		return nil
	}

	id := new(graphql.ID)
	*id = graphql.ID(resp.Id)
	return id
}

func (r *resolver) DeleteProduct(args struct{ ID graphql.ID }) *graphql.ID {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	resp, err := r.c.DeleteProduct(ctx, &agrpc.DeleteProductRequest{Id: string(args.ID)})
	if err != nil {
		r.Printf("DeleteProduct(): %v", err)
		return nil
	}
	id := new(graphql.ID)
	*id = graphql.ID(resp.Id)
	return id
}

type productResolver struct {
	p *app.Product
}

func (r *productResolver) ID() graphql.ID {
	return graphql.ID(r.p.ID)
}

func (r *productResolver) Title() string {
	return r.p.Title
}

func (r *productResolver) Price() float64 {
	return r.p.Price
}

func (r *productResolver) Manufacturer() string {
	return r.p.Manufacturer
}

func (r *productResolver) Description() *string {
	if r.p.Description == "" {
		return nil
	}
	return &r.p.Description
}

func (r *productResolver) Created() string {
	return r.p.Created.String()
}
