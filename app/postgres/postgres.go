package postgres

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

var (
	errNotImplemented = errors.New("postgres: not implemented")
)

type productService struct {
	db *sql.DB
	*log.Logger
}

func NewProductService(logger *log.Logger) *productService {
	return &productService{Logger: logger}
}

func (s *productService) CreateProduct(ctx context.Context, title string, price float64, manufacturer string, description string) (string, error) {
	return "", errNotImplemented
}

func (s *productService) DeleteProduct(ctx context.Context, id string) (string, error) {
	return "", errNotImplemented
}

func (s *productService) Product(ctx context.Context, id string) (*Product, error) {
	return nil, errNotImplemented
}

func (s *productService) SearchProduct(ctx context.Context, title string) (string, error) {
	return "", errNotImplemented
}
