package item

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/product"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ItemService struct {
	itemRepo    Repository
	productRepo product.ProductRepository
}

type Service interface {
	Create(c *gin.Context) (*models.Item, error)
	Delete(c *gin.Context) (float32, error)
	CheckProduct(c *gin.Context) (bool, error)
	Update(c *gin.Context) (float32, error)
	CalculatePrice(c *gin.Context) (float32, error)
	ClearCart(c *gin.Context) error
	Order(c *gin.Context) error
	getItemsFromCartID(c *gin.Context) (*[]models.Item, error)
	parsedCartIdFromCtx(c *gin.Context) (uuid.UUID, error)
	AddItem(c *gin.Context) (float32, error)
}

func NewItemService(repo Repository, productRepo product.ProductRepository) Service {
	if repo == nil {
		return nil
	}

	return &ItemService{itemRepo: repo,
		productRepo: productRepo}
}

func (is *ItemService) AddItem(c *gin.Context) (float32, error) {
	ok, err := is.CheckProduct(c)
	if !ok {
		if err != nil {
			return -1, err
		} else {
			return -1, errors.New("Product with given sku is already in the cart, please update the quantity")
		}
	}
	_, err = is.Create(c)
	if err != nil {
		return -1, err
	}
	totalPrice, err := is.CalculatePrice(c)
	if err != nil {
		return -1, err
	}
	return totalPrice, nil

}
func (is *ItemService) parsedCartIdFromCtx(c *gin.Context) (uuid.UUID, error) {
	cartID, ok := c.Get("cartID")
	zap.L().Debug("itemservice.parsedCartIdFromCtx", zap.Reflect("cartID", cartID))
	if !ok {
		zap.L().Error("itemservice.parsedCartIdFromCtx failed to fetch cartID", zap.Error(errors.New("cartID can not be fetched from context")))
		return uuid.Nil, errors.New("Cart data not found")
	}

	parsedCartId, err := uuid.Parse(fmt.Sprintf("%v", cartID))
	if err != nil {
		zap.L().Error("itemservice.parsedCartIdFromCtx failed to parse cartID", zap.Error(errors.New("cartID can not be parsed")))
		return uuid.Nil, err
	}
	return parsedCartId, nil
}

func (is *ItemService) getItemsFromCartID(c *gin.Context) (*[]models.Item, error) {

	cartID, err := is.parsedCartIdFromCtx(c)
	zap.L().Debug("itemservice.getItemsFromCartID", zap.Reflect("cartID", cartID))
	if err != nil {
		return nil, err
	}
	items, err := is.itemRepo.getItemsInCart(cartID)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (is *ItemService) CalculatePrice(c *gin.Context) (float32, error) {
	zap.L().Debug("itemservice.CalculatePrice")
	items, err := is.getItemsFromCartID(c)
	if err != nil {
		return -1, err
	}
	var totalPrice float32
	for _, v := range *items {
		totalPrice += v.TotalPrice
	}
	return totalPrice, nil
}

// CheckProduct checks if an item with the given product is existed in the cart. Note that function returns true if NOT EXIST.
func (is *ItemService) CheckProduct(c *gin.Context) (bool, error) {

	sku := c.Param("sku")
	zap.L().Debug("itemservice.CheckProduct", zap.Reflect("productsku", sku))

	items, err := is.getItemsFromCartID(c)
	if err != nil {
		return false, err
	}

	for _, v := range *items {
		zap.L().Debug("itemservice.CheckProduct.for", zap.Reflect("item", v))
		if v.Product.Stock.SKU == sku {
			return false, nil
		}
	}
	return true, nil
}

func (is *ItemService) Create(c *gin.Context) (*models.Item, error) {
	sku := c.Param("sku")
	quantity := c.Param("quantity")
	zap.L().Debug("itemservice.Create", zap.Reflect("sku", sku), zap.Reflect("quantity", quantity))

	product, err := is.productRepo.GetBySKU(sku)
	if err != nil {
		return nil, err
	}

	quantityInt, err := strconv.Atoi(quantity)
	if err != nil {
		return nil, errors.New("cannot parse quantity")
	}

	if product.Stock.Number < uint(quantityInt) {
		return nil, fmt.Errorf("Not enough %s in the stock,please request less than %d", *product.Name, (product.Stock.Number + 1))
	}

	// TODO aşğaıyı sil
	// _, err = is.productRepo.CheckStock(sku, uint(quantityInt))
	// if err != nil {
	// 	return nil, err
	// }

	totalPrice := product.Price * float32(quantityInt)
	parsedCartId, err := is.parsedCartIdFromCtx(c)
	if err != nil {
		return nil, err
	}

	itemToCreate := models.Item{
		ProductID:  product.ID,
		Product:    *product,
		Quantity:   uint(quantityInt),
		TotalPrice: totalPrice,
		CartID:     parsedCartId,
	}
	item, err := is.itemRepo.create(&itemToCreate)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (is *ItemService) Delete(c *gin.Context) (float32, error) {

	sku := c.Param("sku")
	zap.L().Debug("itemservice.Delete", zap.Reflect("sku", sku))

	parsedCartId, err := is.parsedCartIdFromCtx(c)
	if err != nil {
		return -1, err
	}

	product, err := is.productRepo.GetBySKU(sku)
	if err != nil {
		return -1, err
	}

	err = is.itemRepo.deleteItemWithProductID(product.ID, parsedCartId)
	if err != nil {
		return -1, err
	}
	totalPrice, err := is.CalculatePrice(c)
	if err != nil {
		return -1, err
	}
	return totalPrice, nil
}

func (is *ItemService) Update(c *gin.Context) (float32, error) {
	sku := c.Param("sku")
	quantity := c.Param("quantity")
	zap.L().Debug("itemservice.Update", zap.Reflect("sku", sku), zap.Reflect("quantity", quantity))

	ok, err := is.CheckProduct(c)
	if ok {
		return -1, errors.New("Product with given sku is not in the cart, please add the product")
	} else {
		if err != nil {
			return -1, err
		}
	}

	product, err := is.productRepo.GetBySKU(sku)

	quantityInt, err := strconv.Atoi(quantity)
	if err != nil {
		return -1, errors.New("cannot parse quantity")
	}

	if product.Stock.Number < uint(quantityInt) {
		return -1, fmt.Errorf("Not enough %s in the stock, please request less than %d", *product.Name, (product.Stock.Number + 1))
	}

	parsedCartId, err := is.parsedCartIdFromCtx(c)
	if err != nil {
		return -1, err
	}

	itemPrice := product.Price * float32(quantityInt)

	err = is.itemRepo.updateItemWithProductID(product.ID, parsedCartId, quantityInt, itemPrice)
	if err != nil {
		return -1, err
	}
	totalPrice, err := is.CalculatePrice(c)
	if err != nil {
		return -1, err
	}
	return totalPrice, nil

}

func (is *ItemService) ClearCart(c *gin.Context) error {

	cartID, ok := c.Get("cartID")
	if !ok {
		// response.RespondWithError(c, errors.New("Cart data not found"))
		return errors.New("Cart data not found")
	}
	parsedCartId, err := uuid.Parse(fmt.Sprintf("%v", cartID))
	if err != nil {
		return err
	}
	items, err := is.itemRepo.getItemsInCart(parsedCartId)
	if err != nil {
		return err
	}
	for i := range *items {
		itemsDeref := *items
		err := is.itemRepo.removeFromCart(&itemsDeref[i])
		if err != nil {
			return err
		}

	}
	return nil
}

func (is *ItemService) Order(c *gin.Context) error {

	// zap.L().Debug("item.order", zap.Reflect("item", orderID))
	orderID, ok := c.Get("orderID")

	zap.L().Debug("item.order.head", zap.Reflect("item", orderID))
	if !ok {
		// response.RespondWithError(c, errors.New("Cart data not found"))
		return errors.New("Order data not found")
	}
	parsedOrderId, err := uuid.Parse(fmt.Sprintf("%v", orderID))
	if err != nil {
		return err
	}
	cartID, ok := c.Get("cartID")
	if !ok {
		// response.RespondWithError(c, errors.New("Cart data not found"))
		return errors.New("Cart data not found")
	}
	parsedCartId, err := uuid.Parse(fmt.Sprintf("%v", cartID))
	if err != nil {
		return err
	}
	items, err := is.itemRepo.getItemsInCart(parsedCartId)
	if err != nil {
		return err
	}

	// TODO order serializerındaki gibi daha iyi handle edilebilir.
	for i := range *items {
		itemsDeref := *items

		zap.L().Debug("item.order", zap.Reflect("item", itemsDeref))

		productSKU := &itemsDeref[i].Product.Stock.SKU
		quantity := &itemsDeref[i].Quantity
		zap.L().Debug("item.order.updateStock", zap.Reflect("productSKU", productSKU), zap.Reflect("quantity", quantity))
		err := is.productRepo.UpdateStock(*productSKU, *quantity)
		if err != nil {
			zap.L().Error("order.service.UpdateStock failed to update product", zap.Error(err))
			return err
		}

		err = is.itemRepo.order(&itemsDeref[i], parsedOrderId)
		if err != nil {
			return err
		}

	}
	return nil
}
