package model

import "time"

type Service struct {
	ID string `json:"id"`
	ServiceName string `json:"service_name"`
	Price int64 `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}