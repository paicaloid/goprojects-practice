package adapters

import (
	"go-clean/entities"
	"go-clean/usecases"

	"gorm.io/gorm"
)

type GormOrderRepository struct {
	db *gorm.DB
}

func NewGormOrderRepository(db *gorm.DB) usecases.OrderRepository {
	return &GormOrderRepository{db: db}
}

func (r *GormOrderRepository) Save(order entities.Order) error {
	if result := r.db.Create(&order); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormOrderRepository) ReadAll() ([]entities.Order, error) {
	var orders []entities.Order
	if result := r.db.Find(&orders); result.Error != nil {
		return nil, result.Error
	}

	return orders, nil
}
