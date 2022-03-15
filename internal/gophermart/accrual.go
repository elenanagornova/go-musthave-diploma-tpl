package gophermart

import (
	"context"
	"go-musthave-diploma-tpl/internal/models"
	"log"
	"time"
)

const (
	workers      = 100
	retryTimeout = 10 * time.Second
)

var validStatuses = map[string]bool{
	"INVALID": true, "PROCESSED": true,
}

func (g Gophermart) GetNewAndProcessingOrders() []models.Order {
	return g.UserOrderRepo.GetNewAndProcessingOrders(context.Background())
}

func (g Gophermart) UpdateOrders(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(5 * time.Second):
			orders := g.UserOrderRepo.GetNewAndProcessingOrders(ctx)
			log.Println(orders)
			if len(orders) > 0 {
				or := newOrderRequest(g.AccrualAddr)
				toWriteOrders := or.run(ctx, orders)
				g.UserOrderRepo.UpdateOrdersStateFromAccrual(ctx, toWriteOrders)
			}
		}
	}
}
