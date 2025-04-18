package core

import "errors"

type OrderService interface {
	CreateOrder(order Order) error
}

type OrderServiceImpl struct {
	repo OrderRepository
}

func NewOrderService(repo OrderRepository) OrderService {
	return &OrderServiceImpl{repo: repo}
}

func (s *OrderServiceImpl) CreateOrder(order Order) error {
	// Busisness Logic
	if order.Total <= 0 {
		return errors.New("total must be positive")
	}

	if err := s.repo.Save(order); err != nil {
		return err
	}

	return nil
}
