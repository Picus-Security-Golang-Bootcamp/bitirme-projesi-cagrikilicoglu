package order

import (
	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"github.com/google/uuid"
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

func (or *OrderRepository) delete(o *models.Order) error {
	result := or.db.Delete(o)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (or *OrderRepository) Create(o *models.Order) error {

	// var o *models.Order
	// c.User = u
	zap.L().Debug("Order.repo.create", zap.Reflect("Order", o))
	if err := or.db.Create(o).Error; err != nil {
		zap.L().Error("Cart.repo.Create failed to create Cart", zap.Error(err))
		return err
	}
	return nil
}
func (or *OrderRepository) getWithID(id uuid.UUID) (*models.Order, error) {
	var o *models.Order
	if err := or.db.Preload("Items.Product").Preload("Items").Where("id", id).First(&o).Error; err != nil {
		zap.L().Error("order.repo.getWithID failed to create Cart", zap.Error(err))
		return nil, err
	}
	return o, nil
}
func (or *OrderRepository) getWithUserID(id uuid.UUID) (*[]models.Order, error) {

	var orders *[]models.Order
	if err := or.db.Unscoped().Preload("Items.Product").Preload("Items").Where("user_id", id).Find(&orders).Error; err != nil {
		zap.L().Error("order.repo.getWithID failed to create Cart", zap.Error(err))
		return nil, err
	}
	return orders, nil
}
