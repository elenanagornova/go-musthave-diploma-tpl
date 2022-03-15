package gophermart

import (
	"context"
	"encoding/json"
	"fmt"
	"go-musthave-diploma-tpl/internal/models"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type orderRequest struct {
	tasks       chan models.Order
	out         chan Result
	client      http.Client
	accrualAddr string
}

func newOrderRequest(accrualAddr string) *orderRequest {
	return &orderRequest{
		tasks:       make(chan models.Order, workers),
		out:         make(chan Result),
		client:      http.Client{Timeout: 30 * time.Second},
		accrualAddr: accrualAddr,
	}
}

func (or *orderRequest) run(ctx context.Context, orders []models.Order) []Result {
	go func() {
		for _, order := range orders {
			or.tasks <- order
		}
	}()
	wg := sync.WaitGroup{}
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			or.worker(ctx)
		}()
	}
	ret := make([]Result, 0, len(orders))
	go func() {
		for result := range or.out {
			ret = append(ret, result)
		}
	}()
	wg.Wait()
	return ret
}

func (or *orderRequest) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case task := <-or.tasks:
			status := task.Status
			accrual := task.Accrual
			resp, err := or.queryOrder(task)
			task.RetryCount++
			if err == nil && resp != nil {
				status = resp.Status
				accrual = resp.Accrual
			}
			oid, _ := strconv.Atoi(task.OrderID)

			or.out <- Result{
				OrderID:    oid,
				Status:     status,
				Accrual:    accrual,
				RetryCount: task.RetryCount,
			}

			<-time.After(retryTimeout)
		}
	}
}

func (or *orderRequest) queryOrder(order models.Order) (*orderResponse, error) {
	url := fmt.Sprintf("%s/api/orders/%s", or.accrualAddr, order.OrderID)
	log.Println("REQUEST to ACCRUAL:" + url)
	resp, err := or.client.Get(url)
	log.Println("RESPONSE FROM ACCRUAL:" + resp.Status + "!!!")
	if err != nil {
		return nil, fmt.Errorf("failed get %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected status 200, got %d (%s)", resp.StatusCode, resp.Status)
	}

	orderResp := new(orderResponse)
	if err := json.NewDecoder(resp.Body).Decode(orderResp); err != nil {
		return nil, fmt.Errorf("failed deserialize %s: %w", url, err)
	}
	return orderResp, err
}

type Result struct {
	OrderID    int
	Status     string
	Accrual    float32
	RetryCount int
}

type orderResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual"`
}
