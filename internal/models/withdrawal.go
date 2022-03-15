package models

import "time"

type Withdrawal struct {
	OrderNum    string    `json:"order_num"`
	Sum         float32   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

func NewWithdrawal(OrderNum string, Sum float32) Withdrawal {
	return Withdrawal{
		OrderNum:    OrderNum,
		Sum:         Sum,
		ProcessedAt: time.Now(),
	}

}
