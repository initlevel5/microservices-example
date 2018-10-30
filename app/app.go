package app

import (
	"context"
	"time"
)

type Product struct {
	ID           string
	Title        string
	Price        float64
	Manufacturer string
	Description  string
	Created      time.Time
}

type ProductService interface {
	CreateProduct(ctx context.Context, title string, price float64, manufacturer string, description string) (string, error)
	DeleteProduct(ctx context.Context, id string) (string, error)
	Product(ctx context.Context, id string) (*Product, error)
	SearchProduct(ctx context.Context, title string) (string, error)
}
