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
	// Get(ctx context.Context, id string) (*Basket, error)
	// GetByCustomerId(ctx context.Context, customerId string) (*Basket, error)
	Create(c *gin.Context) (*models.Item, error)
	Delete(c *gin.Context) error
	CheckProduct(c *gin.Context) (bool, error)
	Update(c *gin.Context) error
	CalculatePrice(c *gin.Context) (float32, error)
	ClearCart(c *gin.Context) error
	Order(c *gin.Context) error
	// Delete(ctx context.Context, id string) (*Basket, error)

	// UpdateItem(ctx context.Context, basketId, itemId string, quantity int) error
	// AddItem(ctx context.Context, basketId, sku string, quantity int, price float64) (string, error)
	// DeleteItem(ctx context.Context, basketId, itemId string) error
}

func NewItemService(repo Repository, productRepo product.ProductRepository) Service {
	if repo == nil {
		return nil
	}

	return &ItemService{itemRepo: repo,
		productRepo: productRepo}
}

func (is *ItemService) Create(c *gin.Context) (*models.Item, error) {
	sku := c.Param("sku")
	quantity := c.Param("quantity")
	product, err := is.productRepo.GetBySKU(sku)
	zap.L().Debug("itemservice.create.productID", zap.Reflect("productID", product.ID))
	if err != nil {
		return nil, errors.New("Failed to get product")
	}
	zap.L().Debug("itemservice.create.parsequantity", zap.Reflect("quantity", quantity))
	quantityInt, err := strconv.Atoi(quantity)
	if err != nil {
		return nil, errors.New("cannot parse quantity")
	}
	zap.L().Debug("itemservice.create.quantity parsed")
	totalPrice := product.Price * float32(quantityInt)
	// TODO burayı ayrı bir cart servisi yapabilirsin
	cartID, ok := c.Get("cartID")
	if !ok {
		// response.RespondWithError(c, errors.New("Cart data not found"))
		return nil, errors.New("Cart data not found")
	}

	// TODO yukarıda getbysku ile çağırma olayını aşağıdaki fonksiyon içinde de yapabilirsin
	_, err = is.productRepo.CheckStock(sku, uint(quantityInt))
	if err != nil {
		return nil, err
	}
	// createdItem:=
	parsedCartId, err := uuid.Parse(fmt.Sprintf("%v", cartID))
	item := models.Item{
		ProductID:  product.ID,
		Product:    *product,
		Quantity:   uint(quantityInt),
		TotalPrice: totalPrice,
		CartID:     parsedCartId,
	}
	itemDb, err := is.itemRepo.create(&item)
	if err != nil {
		return nil, errors.New("cannot create item")
	}
	return itemDb, nil
}

// TODO total price'ı değiştir
func (is *ItemService) Delete(c *gin.Context) error {

	sku := c.Param("sku")
	// TODO burada cart id'ye gerek var mı?
	cartID, ok := c.Get("cartID")
	if !ok {
		// response.RespondWithError(c, errors.New("Cart data not found"))
		return errors.New("Cart data not found")
	}
	parsedCartId, err := uuid.Parse(fmt.Sprintf("%v", cartID))
	if err != nil {
		return err
	}
	id, err := is.productRepo.GetIDBySKU(sku)
	if err != nil {
		return err
	}
	// item, err := is.itemRepo.getItemWithProductSKU(sku, parsedCartId)

	item, err := is.itemRepo.getItemWithProductID(id, parsedCartId)
	if err != nil {
		return err
	}
	err = is.itemRepo.delete(item)
	if err != nil {
		return err
	}
	return nil
}

// TODO bulursa false dönüyo ya bu mantık değiştirilebilir ya fonksiyon ismi değiştiilebilir
func (is *ItemService) CheckProduct(c *gin.Context) (bool, error) {

	sku := c.Param("sku")
	zap.L().Debug("itemservice.checkProduct", zap.Reflect("productsku", sku))
	// TODO burada cart id'ye gerek var mı?
	cartID, ok := c.Get("cartID")
	if !ok {
		// response.RespondWithError(c, errors.New("Cart data not found"))
		return false, errors.New("Cart data not found")
	}
	parsedCartId, err := uuid.Parse(fmt.Sprintf("%v", cartID))
	if err != nil {
		return false, err
	}
	items, err := is.itemRepo.getItemsInCart(parsedCartId)
	zap.L().Debug("itemservice.checkProduc.getItemsInCart", zap.Reflect("items", items))
	if err != nil {
		return false, err
	}
	for _, v := range *items {
		zap.L().Debug("itemservice.checkProduct.for", zap.Reflect("item", v))
		if v.Product.Stock.SKU == sku {
			return false, nil
		}
	}
	return true, nil
}

func (is *ItemService) Update(c *gin.Context) error {
	sku := c.Param("sku")
	quantity := c.Param("quantity")

	quantityInt, err := strconv.Atoi(quantity)
	if err != nil {
		return errors.New("cannot parse quantity")
	}
	zap.L().Debug("itemservice.create.quantity parsed")
	// TODO burada cart id'ye gerek var mı?
	cartID, ok := c.Get("cartID")
	if !ok {
		// response.RespondWithError(c, errors.New("Cart data not found"))
		return errors.New("Cart data not found")
	}
	parsedCartId, err := uuid.Parse(fmt.Sprintf("%v", cartID))
	if err != nil {
		return err
	}
	id, err := is.productRepo.GetIDBySKU(sku)
	if err != nil {
		return err
	}
	// item, err := is.itemRepo.getItemWithProductSKU(sku, parsedCartId)

	// TODO yukarıda getbysku ile çağırma olayını aşağıdaki fonksiyon içinde de yapabilirsin
	_, err = is.productRepo.CheckStock(sku, uint(quantityInt))
	if err != nil {
		return err
	}
	// createdItem:=
	item, err := is.itemRepo.getItemWithProductID(id, parsedCartId)
	if err != nil {
		return err
	}
	item.Quantity = uint(quantityInt)

	// TODO aşağısı after update'de yapılabilir.
	totalPrice := item.Product.Price * float32(item.Quantity)
	item.TotalPrice = totalPrice

	zap.L().Debug("itemservice.create.quantity parsed", zap.Reflect("item", item))
	err = is.itemRepo.update(item)
	if err != nil {
		return err
	}
	return nil

}

func (is *ItemService) CalculatePrice(c *gin.Context) (float32, error) {
	cartID, ok := c.Get("cartID")
	if !ok {
		// response.RespondWithError(c, errors.New("Cart data not found"))
		return -1, errors.New("Cart data not found")
	}
	parsedCartId, err := uuid.Parse(fmt.Sprintf("%v", cartID))
	if err != nil {
		return -1, err
	}
	items, err := is.itemRepo.getItemsInCart(parsedCartId)
	zap.L().Debug("itemservice.checkProduc.getItemsInCart", zap.Reflect("items", items))
	if err != nil {
		return -1, err
	}
	var totalPrice float32
	for _, v := range *items {
		totalPrice += v.TotalPrice
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
