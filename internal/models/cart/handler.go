package cart

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cagrikilicoglu/shopping-basket/internal/models/item"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/response"
	"github.com/cagrikilicoglu/shopping-basket/pkg/config"
	"github.com/cagrikilicoglu/shopping-basket/pkg/middleware"
	"github.com/gin-gonic/gin"
)

type cartHandler struct {
	repo        *CartRepository
	itemService item.Service
}

func NewCartHandler(r *gin.RouterGroup, repo *CartRepository, is item.Service, cfg *config.Config) {
	h := &cartHandler{repo: repo,
		itemService: is}

	r.GET("/cart", middleware.UserAuthMiddleware(cfg.JWTConfig.SecretKey), h.getCart)
	r.POST("/cart/add/sku/{sku}/quantity/{quanity}")

}

func (cr *cartHandler) getCart(c *gin.Context) {
	currentUserId, ok := c.Get("userID")

	if !ok {
		//TODO erroru farklı şekilde handle et
		response.RespondWithError(c, errors.New("User data not found"))
		return
	}

	cart, err := cr.repo.GetByUserID(fmt.Sprintf("%v", currentUserId))
	if err != nil {
		response.RespondWithError(c, errors.New("User not found"))
	}
	response.RespondWithJson(c, http.StatusOK, cartToResponse(cart))
}

func (cr *cartHandler) AddItem(c *gin.Context) {

	currentUserId, ok := c.Get("userID")

	if !ok {
		//TODO erroru farklı şekilde handle et
		response.RespondWithError(c, errors.New("User data not found"))
		return
	}

	cart, err := cr.repo.GetByUserID(fmt.Sprintf("%v", currentUserId))
	if err != nil {
		response.RespondWithError(c, errors.New("User not found"))
		return
	}
	c.Set("cartID", cart.ID)
	item, err := cr.itemService.Create(c)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}
	cart.Items = append(cart.Items, *item)
	cart.TotalPrice += item.TotalPrice
	updatedCart, err := cr.repo.AddItem(cart)
	if err != nil {
		response.RespondWithError(c, err)
		return
	}
	response.RespondWithJson(c, http.StatusOK, cartToResponse(updatedCart))

}
