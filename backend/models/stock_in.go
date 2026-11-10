
package models

import "time"

// StockInTransaction represents stock in header
type StockInTransaction struct {
	ID              int64     `json:"id"`
	TransactionNo   string    `json:"transaction_no"`
	WarehouseID     int64     `json:"warehouse_id"`
	WarehouseName   string    `json:"warehouse_name,omitempty"`
	TransactionDate time.Time `json:"transaction_date"`
	ReferenceNo     *string   `json:"reference_no"`
	SupplierName    *string   `json:"supplier_name"`
	Notes           *string   `json:"notes"`
	TotalQuantity   float64   `json:"total_quantity"`
	TotalAmount     float64   `json:"total_amount"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	Details []StockInDetail `json:"details,omitempty"`
}

// StockInDetail represents stock in item details
type StockInDetail struct {
	ID        int64     `json:"id"`
	StockInID int64     `json:"stock_in_id"`
	ProductID int64     `json:"product_id"`
	ProductName string  `json:"product_name,omitempty"`
	Quantity  float64   `json:"quantity"`
	UnitCost  float64   `json:"unit_cost"`
	TotalCost float64   `json:"total_cost"`
	CreatedAt time.Time `json:"created_at"`
}
