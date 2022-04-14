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
	cartRepo    *cart.CartRepository
	itemService item.Service
}

func NewOrderHandler(r *gin.RouterGroup, orderRepo *OrderRepository, cartRepo *cart.CartRepository, is item.Service, cfg *config.Config) {
	h := &orderHandler{orderRepo: orderRepo,
		cartRepo:    cartRepo,
		itemService: is}

	r.POST("/order", middleware.UserAuthMiddleware(cfg.JWTConfig.SecretKey), h.placeOrder)
	r.DELETE("/order/id/:id/cancel", middleware.UserAuthMiddleware(cfg.JWTConfig.SecretKey), h.cancelOrder)
	r.GET("/order/history", middleware.UserAuthMiddleware(cfg.JWTConfig.SecretKey), h.getOrders)

}

func (oh *orderHandler) cancelOrder(c *gin.Context) {
	// currentUserId, ok := c.Get("userID")

	// if !ok {
	// 	//TODO erroru farklı şekilde handle et
	// 	response.RespondWithError(c, errors.New("User data not found"))
	// 	return
	// }
	// currentUserIDParsed, err := uuid.Parse(fmt.Sprintf("%v", currentUserId))
	// if err != nil {
	// 	response.RespondWithError(c, err)
	// }
	id := c.Param("id")
	orderIDParsed, err := uuid.Parse(fmt.Sprintf("%v", id))
	if err != nil {
		response.RespondWithError(c, err)
	}

	order, err := oh.orderRepo.getWithID(orderIDParsed)
	if err != nil {
		response.RespondWithError(c, errors.New("Order cannot be found"))
		return
	}

	allowedCancelDeadline := order.CreatedAt.AddDate(0, 0, maxAllowedCancelDay)
	// allowedCancelDeadline := order.CreatedAt.Add(time.Minute * 1)
	if !time.Now().Before(allowedCancelDeadline) {
		response.RespondWithError(c, errors.New("Order cannot be canceled after 14 days :("))
		return
	}
	err = oh.orderRepo.delete(order)

	if err != nil {
		response.RespondWithError(c, err)
	}
	response.RespondWithJson(c, http.StatusOK, fmt.Sprintf("Order successfully canceled from the cart"))

	// response.RespondWithJson()

}

func (oh *orderHandler) getOrders(c *gin.Context) {

	currentUserId, ok := c.Get("userID")
	if !ok {
		//TODO erroru farklı şekilde handle et
		response.RespondWithError(c, errors.New("User data not found"))
		return
	}
	currentUserIDParsed, err := uuid.Parse(fmt.Sprintf("%v", currentUserId))
	if err != nil {
		response.RespondWithError(c, err)
	}
	orders, err := oh.orderRepo.getWithUserID(currentUserIDParsed)
	zap.L().Debug("order.handler.getorders", zap.Reflect("orders", orders))
	responsed := ordersToResponse(orders)
	zap.L().Debug("order.handler.getorders", zap.Reflect("responsed", responsed))
	response.RespondWithJson(c, http.StatusOK, responsed)

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
	err = oh.orderRepo.Create(order)
	if err != nil {
		response.RespondWithError(c, errors.New("Order cannot be placed"))
		return
	}
	c.Set("orderID", order.ID)

	err = oh.itemService.Order(c)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}
	err = oh.itemService.ClearCart(c)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}
	orderPlaced, err := oh.orderRepo.getWithID(order.ID)
	if err != nil {
		response.RespondWithError(c, errors.New("Order cannot be placed"))
		return
	}
	response.RespondWithJson(c, http.StatusCreated, orderToResponse(orderPlaced))

}

// TODO aşağıdaki fonksiyon servise taşıanilbir
func createOrderFromCart(c *models.Cart) *models.Order {
	return &models.Order{
		UserID: c.UserID,
		// Items:      c.Items,
		TotalPrice: c.TotalPrice,
	}
}
