package mock

import (
	"context"
	"errors"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/initlevel5/microservices-example/app"
)

var (
	errAlreadyExists = errors.New("mock: already exists")
	errNotFound      = errors.New("mock: not found")
)

var (
	defaultProducts = []*app.Product{
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

// productService represents a mock implementation of app.ProductService interface
type productService struct {
	mu *sync.RWMutex
	db map[string]*app.Product
	*log.Logger
}

func NewProductService(logger *log.Logger) *productService {
	s := &productService{
		mu:     &sync.RWMutex{},
		db:     make(map[string]*app.Product),
		Logger: logger,
	}

	now := time.Now()

	for _, p := range defaultProducts {
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

	s.mu.Lock()
	s.db[p.ID] = p
	s.mu.Unlock()

	return p.ID, nil
}

func (s *productService) DeleteProduct(ctx context.Context, id string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.db[id]; ok {
		delete(s.db, id)
		return id, nil
	}
	return "", errNotFound
}

func (s *productService) Product(ctx context.Context, id string) (*app.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if p, ok := s.db[id]; ok {
		return &app.Product{
			ID:           p.ID,
			Title:        p.Title,
			Price:        p.Price,
			Manufacturer: p.Manufacturer,
			Description:  p.Description,
			Created:	  p.Created,
		}, nil
	}
	return nil, errNotFound
}

func (s *productService) SearchProduct(ctx context.Context, title string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.db {
		if p.Title == title {
			return p.ID, nil
		}
	}
	return "", errNotFound
}
