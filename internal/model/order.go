package model

import (
	"time"
)

type Order struct {
	ID           int         `json:"order_id"`
	CustomerName string      `json:"customer_name"`
	TotalAmount  float64     `json:"total_amount"`
	Status       string      `json:"status"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	Items        []OrderItem `json:"items,omitempty"`
}

type OrderItem struct {
	ID          int     `json:"id,omitempty"`
	OrderID     int     `json:"order_id,omitempty"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

type OrderStatus struct {
	Status string `json:"status"`
}

type OrdersResponse struct {
	Orders []Order `json:"orders"`
	Total  int     `json:"total"`
}

type OrderByIdResponse struct {
	OrderID      int         `json:"order_id"`
	CustomerName string      `json:"customer_name"`
	TotalAmount  float64     `json:"total_amount"`
	Items        []OrderItem `json:"items"`
}

type OrderStatusResponse struct {
	OrderID int    `json:"order_id"`
	Status  string `json:"status"`
}
