package mockdb

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/initlevel5/microservices-example/app"
)

var (
	errAlreadyExists = errors.New("mockdb: already exists")
	errNotFound      = errors.New("mockdb: not found")
)

var (
	products = []*app.Product{
		{
			ID:           "1",
			Title:        "Socks",
			Price:        2.95,
			Manufacturer: "Adidas",
			Description:  "Best socks",
		},
		{
			ID:           "2",
			Title:        "Jeans",
			Price:        20.99,
			Manufacturer: "Levi's",
			Description:  "Good jeans",
		},
		{
			ID:           "3",
			Title:        "T-shirt",
			Price:        7.45,
			Manufacturer: "Ostin",
		},
	}
)

type productService struct {
	db map[string]*app.Product
	*log.Logger
}

func NewProductService(logger *log.Logger) *productService {
	s := &productService{
		db:     make(map[string]*app.Product),
		Logger: logger,
	}

	now := time.Now()

	for _, p := range products {
		p.Created = now
		s.db[p.ID] = p
	}
	return s
}

func (s *productService) CreateProduct(ctx context.Context, title string, price float64, manufacturer string, description string) (string, error) {
	_, err := s.SearchProduct(ctx, title)
	if err == nil {
		err = errAlreadyExists
	}
	if err != nil && err != errNotFound {
		return "", err
	}

	p := &app.Product{
		ID:           strconv.FormatInt(int64(len(s.db)+1), 10),
		Title:        title,
		Price:        price,
		Manufacturer: manufacturer,
		Description:  description,
		Created:      time.Now(),
	}
	s.db[p.ID] = p
	return p.ID, nil
}

func (s *productService) DeleteProduct(ctx context.Context, id string) (string, error) {
	_, ok := s.db[id]
	if ok {
		delete(s.db, id)
		return id, nil
	}
	return "", errNotFound
}

func (s *productService) Product(ctx context.Context, id string) (*app.Product, error) {
	if p, ok := s.db[id]; ok {
		return p, nil
	}
	return nil, errNotFound
}

func (s *productService) SearchProduct(ctx context.Context, title string) (string, error) {
	for _, p := range s.db {
		if p.Title == title {
			return p.ID, nil
		}
	}
	return "", errNotFound
}
