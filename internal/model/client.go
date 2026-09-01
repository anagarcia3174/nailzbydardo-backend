 package model

import "time"

type Client struct {
    ID            string     `json:"id"`
    ClientName    string     `json:"client_name"`
    ContactMethod *string    `json:"contact_method"`
    Notes         *string    `json:"notes"`
    Birthday      *time.Time `json:"birthday"`
    CreatedAt     time.Time  `json:"created_at"`
    DeletedAt     *time.Time `json:"deleted_at"`
}

type ClientSummary struct {
    ID            string  `json:"id"`
    ClientName    string  `json:"client_name"`
    ContactMethod *string `json:"contact_method"`
}

type ClientSpent struct {
    TotalSpent int64 `json:"total_spent"`
    TotalTips int64 `json:"total_tips"`
}