package item

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/product"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ItemService struct {
	itemRepo    Repository
	productRepo product.ProductRepository
}

type Service interface {
	// Get(ctx context.Context, id string) (*Basket, error)
	// GetByCustomerId(ctx context.Context, customerId string) (*Basket, error)
	Create(c *gin.Context) (*models.Item, error)
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
	if err != nil {
		return nil, errors.New("Failed to get product")
	}
	quantityInt, err := strconv.Atoi(quantity)
	if err != nil {
		return nil, errors.New("cannot parse quantity")
	}
	totalPrice := product.Price * float32(quantityInt)
	// TODO burayı ayrı bir cart servisi yapabilirsin
	cartID, ok := c.Get("cartID")
	if !ok {
		// response.RespondWithError(c, errors.New("Cart data not found"))
		return nil, errors.New("Cart data not found")
	}

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
