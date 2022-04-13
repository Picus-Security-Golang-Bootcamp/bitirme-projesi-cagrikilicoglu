package order

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/cart"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/item"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/response"
	"github.com/cagrikilicoglu/shopping-basket/pkg/config"
	"github.com/cagrikilicoglu/shopping-basket/pkg/middleware"
	"github.com/gin-gonic/gin"
)

type orderHandler struct {
	orderRepo   *OrderRepository
	cartRepo    *cart.CartRepository
	itemService item.Service
}

func NewOrderHandler(r *gin.RouterGroup, orderRepo *OrderRepository, cartRepo *cart.CartRepository, is item.Service, cfg *config.Config) {
	h := &orderHandler{orderRepo: orderRepo,
		cartRepo:    cartRepo,
		itemService: is}

	// r.GET("/order", middleware.UserAuthMiddleware(cfg.JWTConfig.SecretKey), h.getCart)
	r.POST("/order", middleware.UserAuthMiddleware(cfg.JWTConfig.SecretKey), h.placeOrder)
	// r.DELETE("/cart/delete/sku/:sku", middleware.UserAuthMiddleware(cfg.JWTConfig.SecretKey), h.DeleteItem)
	// r.PUT("/cart/update/sku/:sku/quantity/:quantity", middleware.UserAuthMiddleware(cfg.JWTConfig.SecretKey), h.UpdateItem)

}

func (oh *orderHandler) placeOrder(c *gin.Context) {
	// TODO aşağısı pek çok fonksiyonda var burayı ayır
	currentUserId, ok := c.Get("userID")
	if !ok {
		//TODO erroru farklı şekilde handle et
		response.RespondWithError(c, errors.New("User data not found"))
		return
	}
	cart, err := oh.cartRepo.GetByUserID(fmt.Sprintf("%v", currentUserId))
	if err != nil {
		response.RespondWithError(c, errors.New("User not found"))
		return
	}
	c.Set("cartID", cart.ID)
	order := createOrderFromCart(cart)
	orderPlaced, err := oh.orderRepo.Create(order)
	if err != nil {
		response.RespondWithError(c, errors.New("Order cannot be placed"))
		return
	}
	oh.itemService.ClearCart(c)
	response.RespondWithJson(c, http.StatusCreated, orderToResponse(orderPlaced))

}

// TODO aşağıdaki fonksiyon servise taşıanilbir
func createOrderFromCart(c *models.Cart) *models.Order {
	return &models.Order{
		UserID:     c.UserID,
		Items:      c.Items,
		TotalPrice: c.TotalPrice,
	}
}
