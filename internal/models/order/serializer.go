package order

import (
	"github.com/cagrikilicoglu/shopping-basket/internal/api"
	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/product"
	"github.com/go-openapi/strfmt"
	"go.uber.org/zap"
)

func orderToResponse(o *models.Order) *api.Order {
	zap.L().Debug("Cart.serializer.cartToResponse", zap.Reflect("userID", o.Status))

	// date, _ := time.Parse(strfmt.RFC3339FullDate, o.CreatedAt)
	orderDate := strfmt.Date(o.CreatedAt)

	idStr := o.ID.String()
	apiItems := make([]*api.Item, 0)

	// TODO aşağıyı itemsin içerisinde bir serializerla yapmak daha doğru
	for i := range o.Items {
		apiItems = append(apiItems, itemToResponse(&o.Items[i]))
	}

	return &api.Order{
		ID:         &idStr,
		Items:      apiItems,
		TotalPrice: &o.TotalPrice,
		Status:     &o.Status,
		Date:       &orderDate,
	}

}

func ordersToResponse(os *[]models.Order) []*api.Order {
	orders := make([]*api.Order, 0)
	for _, o := range *os {
		orders = append(orders, orderToResponse(&o))
	}
	return orders
}

// TODO bu fonksiyonu başka bir yerde handle etmek daha iyi olur
func itemToResponse(i *models.Item) *api.Item {
	quantity := uint32(i.Quantity)
	return &api.Item{
		Product:    product.ProductToResponse(&i.Product),
		Quantity:   &quantity,
		TotalPrice: &i.TotalPrice,
	}
}
