package cart

import (
	"github.com/cagrikilicoglu/shopping-basket/internal/api"
	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/product"
	"go.uber.org/zap"
)

func cartToResponse(c *models.Cart) *api.Cart {
	zap.L().Debug("Cart.serializer.cartToResponse", zap.Reflect("userID", c.UserID))
	userIDstr := c.UserID.String()
	apiItems := make([]*api.Item, 0)

	// TODO aşağıyı itemsin içerisinde bir serializerla yapmak daha doğru
	for i := range c.Items {
		apiItems = append(apiItems, itemToResponse(&c.Items[i]))
	}

	return &api.Cart{
		UserID:     &userIDstr,
		Items:      apiItems,
		TotalPrice: &c.TotalPrice,
	}

}

func itemToResponse(i *models.Item) *api.Item {
	quantity := uint32(i.Quantity)
	return &api.Item{
		Product:    product.ProductToResponse(&i.Product),
		Quantity:   &quantity,
		TotalPrice: &i.TotalPrice,
	}
}
