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
	zap.L().Debug("User.repo.create", zap.Reflect("Cart", c))
	if err := cr.db.Create(c).Error; err != nil {
		zap.L().Error("User.repo.Create failed to create User", zap.Error(err))
		return nil, err
	}
	return c, nil
}
