package gophermart

import (
	"context"
	"go-musthave-diploma-tpl/internal/models"
)

func (g Gophermart) GetNewAndProcessingOrders() []models.Order {
	return g.UserOrderRepo.GetNewAndProcessingOrders(context.Background())
}
func (g Gophermart) GetAccruals(orderID string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	orders := g.GetNewAndProcessingOrders()
	if orders != nil {
		for _, order := range orders {
			response := g.accrualProvider.GetAccrual(ctx, order.OrderID)
			print(response)
		}
	}
}
