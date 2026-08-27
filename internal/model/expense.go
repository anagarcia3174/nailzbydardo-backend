package model

import "time"

type Expense struct {
	ID string `json:"id"`
	ExpenseName string `json:"expense_name"`
	Price int64 `json:"price"`
	DatePurchased time.Time `json:"date_purchased"`
	ReceiptURL *string `json:"receipt_url"`
	CreatedAt time.Time `json:"created_at"`
}