package item

import (
	"github.com/cagrikilicoglu/shopping-basket/internal/api"
	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/product"
	"go.uber.org/zap"
)

func ItemToResponse(i *models.Item) *api.Item {
	zap.L().Debug("item.serializer.itemToResponse", zap.Reflect("item", i))
	quantity := uint32(i.Quantity)
	return &api.Item{
		Product:    product.ProductToResponse(&i.Product),
		Quantity:   &quantity,
		TotalPrice: &i.TotalPrice,
	}
}
