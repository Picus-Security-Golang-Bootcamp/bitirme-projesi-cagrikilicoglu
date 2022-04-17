package item

import (
	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Repository encapsulates the logic to access cart items from the data source.
type Repository interface {
	create(i *models.Item) (*models.Item, error)
	getItemsInCart(cartID uuid.UUID) (*[]models.Item, error)
	updateItemWithProductID(id, cartID uuid.UUID, quantity int, price float32) error
	removeFromCart(i *models.Item) error
	order(i *models.Item, orderID uuid.UUID) error
	deleteItemWithProductID(id, cartID uuid.UUID) error
	getItemWithProductSKU(sku string, cartID uuid.UUID) (*models.Item, error)
	getItemWithProductID(id, cartID uuid.UUID) (*models.Item, error)
}

type ItemRepository struct {
	db *gorm.DB
}

func (ir *ItemRepository) Migration() {
	ir.db.AutoMigrate(&models.Item{})
}

func NewItemRepository(db *gorm.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

//create creates an item in the database
func (ir *ItemRepository) create(i *models.Item) (*models.Item, error) {
	zap.L().Debug("item.repo.create", zap.Reflect("item", i))

	if err := ir.db.Preload("Product").Create(i).Error; err != nil {
		zap.L().Error("item.repo.Create failed to create item", zap.Error(err))
		return nil, err
	}
	return i, nil
}

//getItemsInCart fetches all items in the cart by cartID
func (ir *ItemRepository) getItemsInCart(cartID uuid.UUID) (*[]models.Item, error) {
	zap.L().Debug("item.repo.GetItemsInCart", zap.Reflect("cartID", cartID))
	var items *[]models.Item

	result := ir.db.Order("created_at").Where("is_ordered", false).Where(&models.Item{CartID: cartID}).Preload("Product").Find(&items)
	if result.Error != nil {
		zap.L().Error("item.repo.GetItemsInCart failed to get items", zap.Error(result.Error))
		return nil, result.Error
	}
	return items, nil
}

//getItemWithProductID fetches an item by the productID from the database
func (ir *ItemRepository) getItemWithProductID(id, cartID uuid.UUID) (*models.Item, error) {
	zap.L().Debug("item.repo.GetItemByProductID", zap.Reflect("ID", id), zap.Reflect("cartID", cartID))

	var item *models.Item
	if err := ir.db.Preload("Product").Where(&models.Item{CartID: cartID, ProductID: id}).Where("is_ordered = ?", false).First(&item).Error; err != nil {
		zap.L().Error("item.repo.getItemWithProductID failed to get item", zap.Error(err))
		return nil, err
	}
	return item, nil
}

//getItemWithProductSKU fetches an item by the productSKU from the database
func (ir *ItemRepository) getItemWithProductSKU(sku string, cartID uuid.UUID) (*models.Item, error) {

	zap.L().Debug("item.repo.GetItemByProductSKU", zap.Reflect("SKU", sku))

	var item *models.Item

	result := ir.db.Preload("Product").Where(&models.Item{CartID: cartID}).Joins("left join Product on item.Product_id = Product.id").Where("Product.sku = ?", sku).First(&item)

	zap.L().Debug("item.repo.GetItemByProductSKU.itemcheck", zap.Reflect("item", item))
	if result.Error != nil {
		zap.L().Error("item.repo.GetItemByProductSKU failed to get item", zap.Error(result.Error))
		return nil, result.Error
	}
	return item, nil
}

//updateItemWithProductID updates an item in the database with quantity and price inputs
func (ir *ItemRepository) updateItemWithProductID(id, cartID uuid.UUID, quantity int, price float32) error {
	zap.L().Debug("item.repo.updateItemWithProductID", zap.Reflect("ID", id), zap.Reflect("cartID", cartID))

	result := ir.db.Model(&models.Item{}).Preload("Product").Where(&models.Item{CartID: cartID, ProductID: id}).Where("is_ordered = ?", false).Select("quantity", "total_price").Updates(map[string]interface{}{"quantity": quantity, "total_price": price})

	if err := result.Error; err != nil {
		zap.L().Error("item.repo.updateItemWithProductID failed to update item", zap.Error(err))
		return err
	}
	return nil
}

//removeFromCart updates status of an item to isOrdered
func (ir *ItemRepository) removeFromCart(i *models.Item) error {

	zap.L().Debug("item.repo.removeFromCart", zap.Reflect("items", i))

	if err := ir.db.Model(&i).Preload("Product").Select("is_ordered").Update("is_ordered", true).Error; err != nil {
		zap.L().Error("item.repo.removeFromCart remove the item", zap.Error(err))
		return err
	}
	return nil
}

//order sets an orderID to an item
func (ir *ItemRepository) order(i *models.Item, orderID uuid.UUID) error {

	zap.L().Debug("item.repo.order", zap.Reflect("orderID", orderID))

	if err := ir.db.Model(&i).Preload("Product").Select("order_id").Update("order_id", orderID).Error; err != nil {
		zap.L().Error("item.repo.order", zap.Error(err))
		return err
	}
	return nil
}

//deleteItemWithProductID deletes an item by the productID from the database
func (ir *ItemRepository) deleteItemWithProductID(id, cartID uuid.UUID) error {
	zap.L().Debug("item.repo.deleteItemWithProductID", zap.Reflect("ID", id), zap.Reflect("cartID", cartID))
	result := ir.db.Preload("Product").Where(&models.Item{CartID: cartID, ProductID: id}).Where("is_ordered = ?", false).Delete(&models.Item{})
	if result.Error != nil {
		return result.Error
	}

	return nil

}
