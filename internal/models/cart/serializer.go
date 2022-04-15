package cart

import (
	"github.com/cagrikilicoglu/shopping-basket/internal/api"
	"github.com/cagrikilicoglu/shopping-basket/internal/models"

	"github.com/cagrikilicoglu/shopping-basket/internal/models/item"

	"go.uber.org/zap"
)

func cartToResponse(c *models.Cart) *api.Cart {
	zap.L().Debug("Cart.serializer.cartToResponse", zap.Reflect("cart", c))
	userIDstr := c.UserID.String()
	apiItems := make([]*api.Item, 0)

	for i := range c.Items {
		apiItems = append(apiItems, item.ItemToResponse(&c.Items[i]))
	}

	return &api.Cart{
		UserID:     &userIDstr,
		Items:      apiItems,
		TotalPrice: &c.TotalPrice,
	}

}
