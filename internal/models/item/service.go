package item

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

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

func (is *ItemService) parsedOrderIdFromCtx(c *gin.Context) (uuid.UUID, error) {
	orderID, ok := c.Get("orderID")
	zap.L().Debug("itemservice.parsedOrderIdFromCtx", zap.Reflect("orderID", orderID))

	if !ok {
		zap.L().Error("itemservice.parsedOrderIdFromCtx failed to fetch orderID", zap.Error(errors.New("orderID can not be fetched from context")))
		return uuid.Nil, errors.New("Order data not found")
	}
	parsedOrderId, err := uuid.Parse(fmt.Sprintf("%v", orderID))
	if err != nil {
		zap.L().Error("itemservice.parsedOrderIdFromCtx failed to parse orderID", zap.Error(errors.New("orderID can not be parsed")))
		return uuid.Nil, err
	}
	return parsedOrderId, nil
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

func (is *ItemService) Order(c *gin.Context) error {

	orderID, err := is.parsedOrderIdFromCtx(c)
	zap.L().Debug("itemservice.Order", zap.Reflect("orderID", orderID))
	if err != nil {
		return err
	}

	items, err := is.getItemsFromCartID(c)
	if err != nil {
		return err
	}

	for i := range *items {
		itemsDeref := *items

		sku := &itemsDeref[i].Product.Stock.SKU
		quantity := &itemsDeref[i].Quantity
		product, err := is.productRepo.GetBySKU(*sku)
		if err != nil {
			return err
		}
		if product.Stock.Number < *quantity {
			return fmt.Errorf("Not enough %s in the stock, please request less than %d", *product.Name, (product.Stock.Number + 1))
		}

		var mu sync.Mutex
		mu.Lock()
		err = is.productRepo.UpdateStock(*sku, *quantity)
		if err != nil {
			return err
		}
		mu.Unlock()

		err = is.itemRepo.order(&itemsDeref[i], orderID)
		if err != nil {
			return err
		}
		err = is.itemRepo.removeFromCart(&itemsDeref[i])
		if err != nil {
			return err
		}

	}
	return nil
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
