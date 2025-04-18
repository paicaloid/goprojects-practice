package usecases

import (
	"go-clean/entities"
)

type OrderUseCase interface {
	CreateOrder(order entities.Order) error
	ReadOrders() ([]entities.Order, error)
}

type OrderService struct {
	repo OrderRepository
}

func NewOrderService(repo OrderRepository) OrderUseCase {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(order entities.Order) error {
	return s.repo.Save(order)
}

func (s *OrderService) ReadOrders() ([]entities.Order, error) {
	orders, err := s.repo.ReadAll()
	if err != nil {
		return nil, err
	}
	return orders, nil
}
