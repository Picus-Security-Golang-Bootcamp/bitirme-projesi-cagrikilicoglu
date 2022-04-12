package item

import (
	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Repository encapsulates the logic to access basket from the data source.
type Repository interface {
	// // Get returns the basket with the specified basket Id.
	// Get(ctx context.Context, id string) *Basket
	// // GetByCustomerId returns the basket with the specified customer Id.
	// GetByCustomerId(ctx context.Context, customerId string) *Basket
	// Create saves a new basket in the storage.
	create(i *models.Item) (*models.Item, error)
	// // Update updates the basket with given Is in the storage.
	// Update(ctx context.Context, basket Basket) error
	// // Delete removes the basket with given Is from the storage.
	// Delete(ctx context.Context, basket Basket) error
}

type ItemRepository struct {
	db *gorm.DB
}

func (ir *ItemRepository) Migration() {
	zap.L().Debug("item.repo.migration")
	ir.db.AutoMigrate(&models.Item{})
}

func NewItemRepository(db *gorm.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (ir *ItemRepository) create(i *models.Item) (*models.Item, error) {
	zap.L().Debug("item.repo.create", zap.Reflect("item", i))

	if err := ir.db.Preload("Products").Create(i).Error; err != nil {
		zap.L().Error("item.repo.Create failed to create item", zap.Error(err))
		return nil, err
	}
	return i, nil
}
