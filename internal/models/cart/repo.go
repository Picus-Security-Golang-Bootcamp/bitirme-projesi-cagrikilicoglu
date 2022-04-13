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

func (cr *CartRepository) Create(u *models.User) (*models.Cart, error) {

	var c *models.Cart
	// c.User = u
	zap.L().Debug("Cart.repo.create", zap.Reflect("Cart", c))
	if err := cr.db.Create(c).Error; err != nil {
		zap.L().Error("Cart.repo.Create failed to create Cart", zap.Error(err))
		return nil, err
	}
	return c, nil
}

func (cr *CartRepository) GetByUserID(id string) (*models.Cart, error) {
	var c *models.Cart
	zap.L().Debug("Cart.repo.getByUserID", zap.Reflect("id", id))
	if err := cr.db.Preload("Items.Product").Preload("Items").Where("user_id = ?", id).First(&c).Error; err != nil {
		zap.L().Error("Cart.repo.Create failed to get Cart", zap.Error(err))
		return nil, err
	}
	return c, nil

}
func (cr *CartRepository) AddItem(c *models.Cart) (*models.Cart, error) {
	zap.L().Debug("product.cart.AddItem", zap.Reflect("cart", c))

	if result := cr.db.Preload("Items.Product").Preload("Items").Save(&c); result.Error != nil {
		return nil, result.Error
	}

	return c, nil

}

func (cr *CartRepository) DeleteItem(c *models.Cart) (*models.Cart, error) {
	zap.L().Debug("cart.delete.deleteItem", zap.Reflect("cart", c))

	if result := cr.db.Preload("Items.Product").Preload("Items").Save(&c); result.Error != nil {
		return nil, result.Error
	}

	return c, nil

}

// func(cr *CartRepository) CheckProduct(c *models.Cart,sku string) (bool){
// 	if result :=  cr.db.Preload("Items.Product").Preload("Items").Where(c.Items.Product)
// }
