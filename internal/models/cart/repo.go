package cart

import (
	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func (cr *CartRepository) Migration() {
	cr.db.AutoMigrate(&models.Cart{})
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

// GetByUserID fetches cart data with its items by userID input
func (cr *CartRepository) GetByUserID(id string) (*models.Cart, error) {
	var c *models.Cart
	zap.L().Debug("Cart.repo.getByUserID", zap.Reflect("id", id))
	if err := cr.db.Preload("Items.Product").Preload("Items", "is_ordered = ?", false).Where("user_id = ?", id).First(&c).Error; err != nil {
		zap.L().Error("Cart.repo.getByUserID failed to get Cart", zap.Error(err))
		return nil, err
	}
	return c, nil
}

// GetByCartID fetches cart data with its items by cartID input
func (cr *CartRepository) GetByCartID(id string) (*models.Cart, error) {
	var c *models.Cart
	zap.L().Debug("Cart.repo.GetByCartID", zap.Reflect("id", id))
	if err := cr.db.Preload("Items.Product").Preload("Items", "is_ordered = ?", false).Where("id = ?", id).First(&c).Error; err != nil {
		zap.L().Error("Cart.repo.GetByCartID failed to get Cart", zap.Error(err))
		return nil, err
	}
	return c, nil
}

// UpdateTotalPrice updates totalPrice of the cart
func (cr *CartRepository) UpdateTotalPrice(c *models.Cart, totalPrice float32) error {
	zap.L().Debug("cart.update.updateTotalPrice", zap.Reflect("cart", c), zap.Reflect("totalPrice", totalPrice))

	if result := cr.db.Model(&c).Select("TotalPrice").Update("total_price", totalPrice); result.Error != nil {
		zap.L().Error("cart.update.updateTotalPrice failed to get update total price", zap.Error(result.Error))
		return result.Error
	}
	return nil
}
