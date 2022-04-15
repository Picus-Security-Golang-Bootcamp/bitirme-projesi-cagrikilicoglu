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

type cartHandler struct {
	repo        *CartRepository
	itemService item.Service
}

func NewCartHandler(r *gin.RouterGroup, repo *CartRepository, is item.Service, cfg *config.Config) {
	h := &cartHandler{repo: repo,
		itemService: is}

	r.GET("/", middleware.UserAuthMiddleware(cfg.JWTConfig.SecretKey), h.getCart)
	r.POST("/add/sku/:sku/quantity/:quantity", middleware.UserAuthMiddleware(cfg.JWTConfig.SecretKey), h.AddItem)
	r.DELETE("/delete/sku/:sku", middleware.UserAuthMiddleware(cfg.JWTConfig.SecretKey), h.DeleteItem)
	r.PUT("/update/sku/:sku/quantity/:quantity", middleware.UserAuthMiddleware(cfg.JWTConfig.SecretKey), h.UpdateItem)

}

func (cr *cartHandler) getCart(c *gin.Context) {

	// scid
	// currentUserId, ok := c.Get("userID")
	// zap.L().Debug("cart.handler.getCart", zap.Reflect("currentUserId", currentUserId))
	// if !ok {
	// 	zap.L().Error("cart.handler.getCart failed to fetch userID", zap.Error(errors.New("UserID can not be fetched from context")))
	// 	response.RespondWithError(c, errors.New("User data not found"))
	// 	return
	// }

	// cart, err := cr.repo.GetByUserID(currentUserId.(string))
	// if err != nil {
	// 	response.RespondWithError(c, errors.New("User not found"))
	// 	return
	// }
	// c.Set("cartID", cart.ID)
	// scid

	cart, err := cr.getCartFromUserID(c)

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

func (cr *cartHandler) AddItem(c *gin.Context) {

	// // scid
	// currentUserId, ok := c.Get("userID")
	// zap.L().Debug("cart.handler.AddItem", zap.Reflect("currentUserId", currentUserId))

	// if !ok {
	// 	zap.L().Error("cart.handler.AddItem failed to fetch userID", zap.Error(errors.New("UserID can not be fetched from context")))
	// 	response.RespondWithError(c, errors.New("User data not found"))
	// 	return
	// }

	// cart, err := cr.repo.GetByUserID(currentUserId.(string))
	// if err != nil {
	// 	response.RespondWithError(c, errors.New("User not found"))
	// 	return
	// }
	// c.Set("cartID", cart.ID)
	// //scid

	cart, err := cr.getCartFromUserID(c)
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

func (cr *cartHandler) DeleteItem(c *gin.Context) {

	// //scid
	// currentUserId, ok := c.Get("userID")
	// zap.L().Debug("cart.handler.DeleteItem", zap.Reflect("currentUserId", currentUserId))

	// if !ok {
	// 	zap.L().Error("cart.handler.DeleteItem failed to fetch userID", zap.Error(errors.New("UserID can not be fetched from context")))
	// 	response.RespondWithError(c, errors.New("User data not found"))
	// 	return
	// }

	// cart, err := cr.repo.GetByUserID(fmt.Sprintf("%v", currentUserId))
	// if err != nil {
	// 	response.RespondWithError(c, errors.New("User not found"))
	// 	return
	// }
	// c.Set("cartID", cart.ID)
	// //scid

	cart, err := cr.getCartFromUserID(c)
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

func (cr *cartHandler) UpdateItem(c *gin.Context) {
	// // scid
	// currentUserId, ok := c.Get("userID")
	// zap.L().Debug("cart.handler.UpdateItem", zap.Reflect("currentUserId", currentUserId))

	// if !ok {
	// 	response.RespondWithError(c, errors.New("User data not found"))
	// 	return
	// }

	// cart, err := cr.repo.GetByUserID(fmt.Sprintf("%v", currentUserId))
	// if err != nil {
	// 	response.RespondWithError(c, errors.New("User not found"))
	// 	return
	// }
	// c.Set("cartID", cart.ID)
	// //scid

	cart, err := cr.getCartFromUserID(c)
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

func (cr *cartHandler) getCartFromUserID(c *gin.Context) (*models.Cart, error) {

	currentUserId, ok := c.Get("userID")
	zap.L().Debug("cart.handler.getCart", zap.Reflect("currentUserId", currentUserId))
	if !ok {
		zap.L().Error("cart.handler.getCart failed to fetch userID", zap.Error(errors.New("UserID can not be fetched from context")))
		return nil, errors.New("User data not found")
	}

	cart, err := cr.repo.GetByUserID(currentUserId.(string))
	if err != nil {
		return nil, err
	}
	c.Set("cartID", cart.ID)
	return cart, nil

}

// func (cr *cartHandler) UpdatePrice(c *gin.Context) {

// 	currentUserId, ok := c.Get("userID")

// 	if !ok {
// 		response.RespondWithError(c, errors.New("User data not found"))
// 		return
// 	}

// 	cart, err := cr.repo.GetByUserID(fmt.Sprintf("%v", currentUserId))
// 	if err != nil {
// 		response.RespondWithError(c, errors.New("User not found"))
// 		return
// 	}
// 	c.Set("cartID", cart.ID)

// 	ok, err = cr.itemService.CheckProduct(c)

// 	if ok {
// 		response.RespondWithError(c, errors.New("Product with given sku is not in the cart please add the product first"))
// 		return
// 	}

// 	totalPrice, err := cr.itemService.CalculatePrice(c)
// 	if err != nil {
// 		response.RespondWithError(c, err)
// 		return
// 	}
// 	err = cr.repo.UpdateTotalPrice(cart, totalPrice)
// 	if err != nil {
// 		response.RespondWithError(c, err)
// 		return
// 	}
// 	updatedCart, err := cr.repo.GetByUserID(fmt.Sprintf("%v", currentUserId))
// 	response.RespondWithJson(c, http.StatusOK, cartToResponse(updatedCart))

// }

// func (cr *cartHandler) SetCartIDToHeader(c *gin.Context) error {
// 	currentUserId, ok := c.Get("userID")
// 	zap.L().Debug("cart.handler.getCart", zap.Reflect("currentUserId", currentUserId))
// 	if !ok {
// 		zap.L().Error("cart.handler.getCart failed to fetch userID", zap.Error(errors.New("UserID can not be fetched from context")))
// 		return errors.New("User data not found")
// 	}

// 	cart, err := cr.repo.GetByUserID(currentUserId.(string))
// 	if err != nil {
// 		return err
// 	}
// 	c.Set("cartID", cart.ID)
// }

// func (cr *cartHandler) CheckProduct(c *gin.Context) {

// 	currentUserId, ok := c.Get("userID")

// 	if !ok {
// 		//TODO erroru farklı şekilde handle et
// 		response.RespondWithError(c, errors.New("User data not found"))
// 		return
// 	}

// 	cart, err := cr.repo.GetByUserID(fmt.Sprintf("%v", currentUserId))
// 	if err != nil {
// 		response.RespondWithError(c, errors.New("User not found"))
// 		return
// 	}
// 	c.Set("cartID", cart.ID)
// 	ok, err = cr.itemService.CheckProduct()
// 	if err != nil {
// 		response.RespondWithError(c, err)
// 		return
// 	}

// }
