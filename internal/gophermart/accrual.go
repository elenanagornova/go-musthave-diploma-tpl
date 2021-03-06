package gophermart

import (
	"context"
	"go-musthave-diploma-tpl/internal/models"
	"log"
	"time"
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
		case <-time.After(1 * time.Second):
			orders := g.UserOrderRepo.GetNewAndProcessingOrders(ctx)
			if len(orders) > 0 {
				log.Println("Updating", len(orders))
				or := newOrderRequest(g.AccrualAddr)
				toWriteOrders := or.run(ctx, orders)
				g.UpdateStates(ctx, toWriteOrders)
			}
		}
	}
}
