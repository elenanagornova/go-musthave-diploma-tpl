package models

import "time"

type OrderStatus string

const (
	NEW        = "NEW"        //"заказ загружен в систему, но не попал в обработку"
	PROCESSING = "PROCESSING" //"вознаграждение за заказ рассчитывается"
	INVALID    = "INVALID"    //"система расчёта вознаграждений отказала в расчёте"
	PROCESSED  = "PROCESSED"  //"данные по заказу проверены и информация о расчёте"
)

// Order заказ, загружаемый пользователем, за который могут быть начислены баллы лояльности
type Order struct {
	UID        string    `json:",omitempty"`
	OrderID    string    `json:"number"`
	UploadedAt time.Time `json:"uploaded_at"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual"`

	// RetryCount количество попыток получить начисления во внешнем сервисе
	RetryCount int `json:",omitempty"`
}

type OrderOption func(*Order)

func NewOrder(userID string, orderID string, orderOpt ...OrderOption) Order {
	order := Order{
		UID:        userID,
		OrderID:    orderID,
		UploadedAt: time.Now(),
		Status:     NEW,
		Accrual:    0,
		RetryCount: 0,
	}
	for _, opt := range orderOpt {
		opt(&order)
	}
	return order
}
