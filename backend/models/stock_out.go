
package models

import "time"

type StockOutTransaction struct {
	ID              int64     `json:"id"`
	TransactionNo   string    `json:"transaction_no"`
	TransactionDate time.Time `json:"transaction_date"`
	WarehouseID     int64     `json:"warehouse_id"`
	WarehouseName   string    `json:"warehouse_name"`
	ReferenceNo     *string   `json:"reference_no"`
	Destination     *string   `json:"destination"`
	Notes           *string   `json:"notes"`
	TotalQuantity   float64   `json:"total_quantity"`
	TotalAmount     float64   `json:"total_amount"`

	Details []StockOutDetail `json:"details,omitempty"`
}

type StockOutDetail struct {
	ID          int64     `json:"id"`
	StockOutID  int64     `json:"stock_out_id"`
	ProductID   int64     `json:"product_id"`
	ProductName string    `json:"product_name,omitempty"`
	Quantity    float64   `json:"quantity"`
	UnitCost    float64   `json:"unit_cost"`
	TotalCost   float64   `json:"total_cost"`
	CreatedAt   time.Time `json:"created_at"`
}
