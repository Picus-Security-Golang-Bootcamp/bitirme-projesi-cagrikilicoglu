package item

import (
	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"github.com/google/uuid"
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

	getItemWithProductSKU(sku string, cartID uuid.UUID) (*models.Item, error)
	delete(i *models.Item) error
	getItemsInCart(cartID uuid.UUID) (*[]models.Item, error)
	getItemWithProductID(id, cartID uuid.UUID) (*models.Item, error)
	update(i *models.Item) error
	removeFromCart(i *models.Item) error
	order(i *models.Item, orderID uuid.UUID) error
	deleteItemWithProductID(id, cartID uuid.UUID) error
	updateItemWithProductID(id, cartID uuid.UUID, quantity int, price float32) error
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

func (ir *ItemRepository) getItemsInCart(cartID uuid.UUID) (*[]models.Item, error) {
	zap.L().Debug("item.repo.GetItemsInCart", zap.Reflect("cartID", cartID))
	var items *[]models.Item

	result := ir.db.Where("is_ordered", false).Where(&models.Item{CartID: cartID}).Preload("Product").Find(&items)
	if result.Error != nil {
		zap.L().Error("item.repo.GetItemsInCart failed to get items", zap.Error(result.Error))
		return nil, result.Error
	}
	return items, nil
}

func (ir *ItemRepository) create(i *models.Item) (*models.Item, error) {
	zap.L().Debug("item.repo.create", zap.Reflect("item", i))

	if err := ir.db.Preload("Product").Create(i).Error; err != nil {
		zap.L().Error("item.repo.Create failed to create item", zap.Error(err))
		return nil, err
	}
	return i, nil
}
func (ir *ItemRepository) deleteItemWithProductID(id, cartID uuid.UUID) error {
	zap.L().Debug("item.repo.deleteItemWithProductID", zap.Reflect("ID", id), zap.Reflect("cartID", cartID))
	result := ir.db.Preload("Product").Where(&models.Item{CartID: cartID, ProductID: id}).Where("is_ordered = ?", false).Delete(&models.Item{})
	if result.Error != nil {
		return result.Error
	}

	return nil

}

func (ir *ItemRepository) updateItemWithProductID(id, cartID uuid.UUID, quantity int, price float32) error {
	zap.L().Debug("item.repo.updateItemWithProductID", zap.Reflect("ID", id), zap.Reflect("cartID", cartID))

	// result := ir.db.Model(&models.Item{}).Preload("Product").Where(&models.Item{CartID: cartID, ProductID: id}).Where("is_ordered = ?", false).Select("quantity").Update("quantity", quantity)

	result := ir.db.Model(&models.Item{}).Preload("Product").Where(&models.Item{CartID: cartID, ProductID: id}).Where("is_ordered = ?", false).Select("quantity", "total_price").Updates(map[string]interface{}{"quantity": quantity, "total_price": price})

	if err := result.Error; err != nil {
		zap.L().Error("item.repo.updateItemWithProductID failed to update item", zap.Error(err))
		return err
	}
	return nil
}
func (ir *ItemRepository) getItemWithProductID(id, cartID uuid.UUID) (*models.Item, error) {
	zap.L().Debug("item.repo.GetItemByProductID", zap.Reflect("ID", id), zap.Reflect("cartID", cartID))

	var item *models.Item
	if err := ir.db.Preload("Product").Where(&models.Item{CartID: cartID, ProductID: id}).Where("is_ordered = ?", false).First(&item).Error; err != nil {
		zap.L().Error("item.repo.getItemWithProductID failed to get item", zap.Error(err))
		return nil, err
	}
	return item, nil
}

// func (ir *ItemRepository) update(i *models.Item) (*models.Item, error) {
// 	zap.L().Debug("item.repo.update", zap.Reflect("item", i))

// 	if err := ir.db.Preload("Products").Save(i).Error; err != nil {
// 		zap.L().Error("item.repo.Save failed to save item", zap.Error(err))
// 		return nil, err
// 	}
// 	return i, nil
// }

func (ir *ItemRepository) getItemWithProductSKU(sku string, cartID uuid.UUID) (*models.Item, error) {

	zap.L().Debug("item.repo.GetItemByProductSKU", zap.Reflect("SKU", sku))

	var item *models.Item
	// id, err := is.pro
	// result := ir.db.Preload(clause.Associations).Where(&models.Item{Product: models.Product{Stock: models.Stock{SKU: sku}}, CartID: cartID}).First(&item)
	result := ir.db.Preload("Product").Where(&models.Item{CartID: cartID}).Joins("left join Product on item.Product_id = Product.id").Where("Product.sku = ?", sku).First(&item)

	zap.L().Debug("item.repo.GetItemByProductSKU.itemcheck", zap.Reflect("item", item))
	if result.Error != nil {
		zap.L().Error("item.repo.GetItemByProductSKU failed to get item", zap.Error(result.Error))
		return nil, result.Error
	}
	return item, nil
}
func (ir *ItemRepository) delete(i *models.Item) error {
	result := ir.db.Delete(i)
	if result.Error != nil {
		return result.Error
	}

	return nil

}

func (ir *ItemRepository) update(i *models.Item) error {
	zap.L().Debug("item.repo.update.item", zap.Reflect("item", i))
	// if err := ir.db.Preload("Product").Select("OrderID").Save(i).Error; err != nil {

	result := ir.db.Model(&i).Preload("Product").Select("quantity").Update("quantity", int(i.Quantity))

	zap.L().Debug("item.repo.update.item.result", zap.Reflect("result", i))
	if err := result.Error; err != nil {
		zap.L().Error("item.repo.update failed to update item", zap.Error(err))
		return err
	}
	return nil
}

func (ir *ItemRepository) removeFromCart(i *models.Item) error {

	zap.L().Debug("itemservice.repo.removefromcart", zap.Reflect("items", i))

	if err := ir.db.Model(&i).Preload("Product").Select("is_ordered").Update("is_ordered", true).Error; err != nil {
		zap.L().Error("item.repo.update failed to update item", zap.Error(err))
		return err
	}
	return nil
}
func (ir *ItemRepository) order(i *models.Item, orderID uuid.UUID) error {

	zap.L().Debug("item.repo.order", zap.Reflect("orderID", orderID))

	if err := ir.db.Model(&i).Preload("Product").Select("order_id").Update("order_id", orderID).Error; err != nil {
		zap.L().Error("item.repo.update failed to update item", zap.Error(err))
		return err
	}
	return nil
}
