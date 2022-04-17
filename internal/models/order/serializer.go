package order

import (
	"github.com/cagrikilicoglu/shopping-basket/internal/api"
	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"github.com/cagrikilicoglu/shopping-basket/internal/models/item"
	"github.com/go-openapi/strfmt"
	"go.uber.org/zap"
)

// orderToResponse converts order database model to response model
func orderToResponse(o *models.Order) *api.Order {
	zap.L().Debug("Order.serializer.orderToResponse", zap.Reflect("order", o))

	apiItems := make([]*api.Item, 0)

	orderDate := strfmt.Date(o.CreatedAt)
	idStr := o.ID.String()

	for i := range o.Items {
		apiItems = append(apiItems, item.ItemToResponse(&o.Items[i]))
	}
	return &api.Order{
		ID:         &idStr,
		Items:      apiItems,
		TotalPrice: &o.TotalPrice,
		Status:     &o.Status,
		Date:       &orderDate,
	}

}

// ordersToResponse converts order database model to response model as a batch
func ordersToResponse(os *[]models.Order) []*api.Order {
	zap.L().Debug("Order.serializer.ordersToResponse", zap.Reflect("orders", os))
	orders := make([]*api.Order, 0)
	for i := range *os {
		osDeref := *os
		orders = append(orders, orderToResponse(&osDeref[i]))
	}
	return orders
}
