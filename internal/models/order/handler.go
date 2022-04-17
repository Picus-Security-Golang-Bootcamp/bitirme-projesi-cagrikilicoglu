package order

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/cart"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/item"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/response"
	"github.com/cagrikilicoglu/shopping-basket/pkg/config"
	"github.com/cagrikilicoglu/shopping-basket/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// TODO aşağıdakini başka bir yere taşıyabilir miyiz? config gibi
var maxAllowedCancelDay = 14

type orderHandler struct {
	orderRepo   *OrderRepository
	cartRepo    cart.CartRepo
	itemService item.Service
}

func NewOrderHandler(r *gin.RouterGroup, orderRepo *OrderRepository, cartRepo cart.CartRepo, is item.Service, cfg *config.Config) {
	h := &orderHandler{orderRepo: orderRepo,
		cartRepo:    cartRepo,
		itemService: is}

	r.POST("/order", middleware.UserAuthMiddleware(cfg.JWTConfig.SecretKey), h.placeOrder)
	r.DELETE("/order/id/:id/cancel", middleware.UserAuthMiddleware(cfg.JWTConfig.SecretKey), h.cancelOrder)
	r.GET("/order/history", middleware.UserAuthMiddleware(cfg.JWTConfig.SecretKey), h.getOrders)

}

func (oh *orderHandler) placeOrder(c *gin.Context) {

	cart, err := oh.getCartFromUserID(c)
	zap.L().Debug("order.handler.placeOrder", zap.Reflect("cart", cart))
	if err != nil {
		response.RespondWithError(c, err)
		return
	}
	if len(cart.Items) == 0 {
		response.RespondWithError(c, errors.New("Your cart is empty"))
		return
	}

	order := createOrderFromCart(cart)
	err = oh.orderRepo.Create(order)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}

	c.Set("orderID", order.ID)

	err = oh.itemService.Order(c)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}

	orderPlaced, err := oh.orderRepo.getWithID(order.ID)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}
	response.RespondWithJson(c, http.StatusCreated, orderToResponse(orderPlaced))

}

func (oh *orderHandler) cancelOrder(c *gin.Context) {

	id := c.Param("id")
	zap.L().Debug("order.handler.cancelOrder", zap.Reflect("id", id))

	orderIDParsed, err := uuid.Parse(fmt.Sprintf("%v", id))
	if err != nil {
		response.RespondWithError(c, err)
	}

	order, err := oh.orderRepo.getWithID(orderIDParsed)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}

	allowedCancelDeadline := order.CreatedAt.AddDate(0, 0, maxAllowedCancelDay)
	// TODO remove line below
	// allowedCancelDeadline := order.CreatedAt.Add(time.Minute * 1)
	if !time.Now().Before(allowedCancelDeadline) {
		response.RespondWithError(c, errors.New("Order cannot be canceled after 14 days :("))
		return
	}
	err = oh.orderRepo.delete(order)

	if err != nil {
		response.RespondWithError(c, err)
	}
	response.RespondWithJson(c, http.StatusOK, "Order successfully canceled")

}

func (oh *orderHandler) getOrders(c *gin.Context) {

	userID, ok := c.Get("userID")
	zap.L().Debug("order.handler.cancelOrder", zap.Reflect("userid", userID))
	if !ok {
		response.RespondWithError(c, errors.New("User data not found"))
		return
	}
	userIDParsed, err := uuid.Parse(fmt.Sprintf("%v", userID))
	if err != nil {
		response.RespondWithError(c, err)
	}
	orders, err := oh.orderRepo.getWithUserID(userIDParsed)

	response.RespondWithJson(c, http.StatusOK, ordersToResponse(orders))

}

func createOrderFromCart(c *models.Cart) *models.Order {
	return &models.Order{
		UserID:     c.UserID,
		TotalPrice: c.TotalPrice,
	}
}

func (oh *orderHandler) getCartFromUserID(c *gin.Context) (*models.Cart, error) {

	currentUserId, ok := c.Get("userID")
	zap.L().Debug("cart.handler.getCartFromUserID", zap.Reflect("currentUserId", currentUserId))
	if !ok {
		zap.L().Error("cart.handler.getCartFromUserID failed to fetch userID", zap.Error(errors.New("UserID can not be fetched from context")))
		return nil, errors.New("User data not found")
	}

	cart, err := oh.cartRepo.GetByUserID(currentUserId.(string))
	if err != nil {
		return nil, err
	}
	c.Set("cartID", cart.ID)
	return cart, nil
}
