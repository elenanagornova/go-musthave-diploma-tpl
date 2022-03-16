package gophermart

import (
	"context"
	"encoding/json"
	"fmt"
	"go-musthave-diploma-tpl/internal/models"
	"io/ioutil"
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
		close(or.tasks)
	}()
	wg := sync.WaitGroup{}
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			or.worker()
		}()
	}
	ret := make([]Result, 0, len(orders))
	go func() {
		for result := range or.out {
			ret = append(ret, result)
		}
	}()
	log.Println("waiting complete")
	wg.Wait()
	log.Println("complete")
	return ret
}

func (or *orderRequest) worker() {
	for task := range or.tasks {
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

func (or *orderRequest) queryOrder(order models.Order) (*OrderResponse, error) {
	url := fmt.Sprintf("%s/api/orders/%s", or.accrualAddr, order.OrderID)
	resp, err := or.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed get %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected status 200, got %d (%s)", resp.StatusCode, resp.Status)
	}

	var orderResp OrderResponse
	b, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(b))
	if err := json.Unmarshal(b, orderResp); err != nil {
		return nil, fmt.Errorf("failed deserialize %s: %w", url, err)
	}
	return &orderResp, nil
}

type Result struct {
	OrderID    int
	Status     string
	Accrual    float32
	RetryCount int
}

type OrderResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual,omitempty"`
}
