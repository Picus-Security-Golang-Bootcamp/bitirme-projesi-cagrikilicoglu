package order

import (
	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func (or *OrderRepository) Migration() {
	or.db.AutoMigrate(&models.Order{})
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (or *OrderRepository) Create(o *models.Order) (*models.Order, error) {

	// var o *models.Order
	// c.User = u
	zap.L().Debug("Order.repo.create", zap.Reflect("Order", o))
	if err := or.db.Create(o).Error; err != nil {
		zap.L().Error("Cart.repo.Create failed to create Cart", zap.Error(err))
		return nil, err
	}
	return o, nil
}
