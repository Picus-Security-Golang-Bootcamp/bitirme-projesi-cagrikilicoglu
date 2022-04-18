package cart

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/item"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/response"
	"github.com/cagrikilicoglu/shopping-basket/pkg/config"
	"github.com/cagrikilicoglu/shopping-basket/pkg/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var maxItemsForCart = 20

type cartHandler struct {
	repo        *CartRepository
	itemService item.Service
}

func NewCartHandler(r *gin.RouterGroup, repo *CartRepository, is item.Service, cfg *config.Config) {
	h := &cartHandler{repo: repo,
		itemService: is}

	r.GET("/", middleware.UserAuthMiddleware(cfg.JWTConfig.SecretKey), h.getCart)
	r.POST("/add/sku/:sku/quantity/:quantity", middleware.UserAuthMiddleware(cfg.JWTConfig.SecretKey), h.addItem)
	r.DELETE("/delete/sku/:sku", middleware.UserAuthMiddleware(cfg.JWTConfig.SecretKey), h.deleteItem)
	r.PUT("/update/sku/:sku/quantity/:quantity", middleware.UserAuthMiddleware(cfg.JWTConfig.SecretKey), h.updateItem)
}

// getCart fetches cart data from user id
func (cr *cartHandler) getCart(c *gin.Context) {

	cart, err := cr.getCartFromUserID(c)
	zap.L().Debug("cart.handler.getCart", zap.Reflect("cart", cart))

	if err != nil {
		response.RespondWithError(c, err)
		return
	}

	totalPrice, err := cr.itemService.CalculatePrice(c)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}

	err = cr.repo.UpdateTotalPrice(cart, totalPrice)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}
	response.RespondWithJson(c, http.StatusOK, cartToResponse(cart))
}

// addItem adds a product to the cart and returns updated cart
func (cr *cartHandler) addItem(c *gin.Context) {

	cart, err := cr.getCartFromUserID(c)
	zap.L().Debug("cart.handler.addItem", zap.Reflect("cart", cart))

	if err != nil {
		response.RespondWithError(c, err)
		return
	}

	err = checkItemNumber(cart)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}

	totalPrice, err := cr.itemService.AddItem(c)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}

	err = cr.repo.UpdateTotalPrice(cart, totalPrice)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}

	updatedCart, err := cr.repo.GetByCartID(fmt.Sprintf("%v", cart.ID))
	if err != nil {
		response.RespondWithError(c, err)
		return
	}

	response.RespondWithJson(c, http.StatusOK, cartToResponse(updatedCart))
}

// deleteItem deletes a product from the cart
func (cr *cartHandler) deleteItem(c *gin.Context) {

	cart, err := cr.getCartFromUserID(c)
	zap.L().Debug("cart.handler.deleteItem", zap.Reflect("cart", cart))

	totalPrice, err := cr.itemService.Delete(c)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}

	err = cr.repo.UpdateTotalPrice(cart, totalPrice)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}
	response.RespondWithJson(c, http.StatusOK, "Item successfully deleted from the cart")

}

// updateItem updates quantity of a product that is already in the cart
func (cr *cartHandler) updateItem(c *gin.Context) {

	cart, err := cr.getCartFromUserID(c)
	zap.L().Debug("cart.handler.updateItem", zap.Reflect("cart", cart))

	totalPrice, err := cr.itemService.Update(c)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}

	err = cr.repo.UpdateTotalPrice(cart, totalPrice)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}

	updatedCart, err := cr.repo.GetByCartID(fmt.Sprintf("%v", cart.ID))
	if err != nil {
		response.RespondWithError(c, err)
		return
	}
	response.RespondWithJson(c, http.StatusOK, cartToResponse(updatedCart))

}

// getCartFromUserID fetches cart of the user by ID
func (cr *cartHandler) getCartFromUserID(c *gin.Context) (*models.Cart, error) {

	currentUserId, ok := c.Get("userID")
	zap.L().Debug("cart.handler.getCartFromUserID", zap.Reflect("currentUserId", currentUserId))
	if !ok {
		zap.L().Error("cart.handler.getCartFromUserID failed to fetch userID", zap.Error(errors.New("UserID can not be fetched from context")))
		return nil, errors.New("User data not found")
	}

	cart, err := cr.repo.GetByUserID(currentUserId.(string))
	if err != nil {
		return nil, err
	}
	c.Set("cartID", cart.ID)
	return cart, nil

}

//checkItemNumber checks if item number in the cart is below maximum
func checkItemNumber(c *models.Cart) error {
	if len(c.Items) >= maxItemsForCart {
		return errors.New("You exceed maximum number of items")
	}
	return nil
}
