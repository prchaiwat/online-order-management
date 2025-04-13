package service

import (
	"context"
	"sync"
	"time"

	"order-management-system/internal/config"
	"order-management-system/internal/model"
	"order-management-system/internal/repository"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order model.Order) (int, error)
	CreateOrderWithGoroutine(order model.Order) (int, error)
	GetOrder(ctx context.Context, id int) (model.Order, error)
	GetOrders(ctx context.Context, page, size int) (model.OrdersResponse, error)
	UpdateOrderStatus(ctx context.Context, id int, status string) error
}

type orderService struct {
	repo           repository.OrderRepository
	config         *config.Config
	orderSemaphore chan struct{}
	waitGroup      sync.WaitGroup
}

func NewOrderService(repo repository.OrderRepository, config *config.Config) OrderService {
	return &orderService{
		repo:           repo,
		config:         config,
		orderSemaphore: make(chan struct{}, config.MaxConcurrentOrders),
	}
}

func (s *orderService) CreateOrder(ctx context.Context, order model.Order) (int, error) {
	return s.repo.CreateOrder(ctx, order)
}

func (s *orderService) CreateOrderWithGoroutine(order model.Order) (int, error) {
	resultChannel := make(chan struct {
		orderID int
		err     error
	}, 1)

	s.waitGroup.Add(1)

	go func() {

		s.orderSemaphore <- struct{}{}

		defer func() { <-s.orderSemaphore }()
		defer s.waitGroup.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.RequestTimeoutSec)*time.Second)
		defer cancel()

		orderID, err := s.repo.CreateOrder(ctx, order)

		resultChannel <- struct {
			orderID int
			err     error
		}{orderID, err}
	}()

	result := <-resultChannel
	return result.orderID, result.err
}

func (s *orderService) GetOrder(ctx context.Context, id int) (model.Order, error) {
	return s.repo.GetOrderByID(ctx, id)
}

func (s *orderService) GetOrders(ctx context.Context, page, size int) (model.OrdersResponse, error) {
	orders, total, err := s.repo.GetOrders(ctx, page, size)
	if err != nil {
		return model.OrdersResponse{}, err
	}
	return model.OrdersResponse{
		Orders: orders,
		Total:  total,
	}, nil
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, id int, status string) error {
	return s.repo.UpdateOrderStatus(ctx, id, status)
}
