
package models

import "time"

// DashboardData contains all dashboard information.
type DashboardData struct {
	TotalProducts       int64   `json:"total_products"`
	TotalWarehouses     int64   `json:"total_warehouses"`
	TotalStockInQuantity  float64 `json:"total_stock_in_quantity"`
	TotalStockOutQuantity float64 `json:"total_stock_out_quantity"`

	TotalInventoryValue float64 `json:"total_inventory_value"`
	TotalStockInAmount  float64 `json:"total_stock_in_amount"`
	TotalStockOutAmount float64 `json:"total_stock_out_amount"`

	MonthlyMovement []MonthlyMovement `json:"monthly_movement"`
	CategoryValues  []CategoryValue   `json:"category_values"`
	LowStockItems   []LowStockItem    `json:"low_stock_items"`
	RecentTransactions []RecentTransaction `json:"recent_transactions"`
}

// Monthly stock movement for chart.
type MonthlyMovement struct {
	Month          string  `json:"month"`
	StockInAmount  float64 `json:"stock_in_amount"`
	StockOutAmount float64 `json:"stock_out_amount"`
}

// Inventory value grouped by category.
type CategoryValue struct {
	Category   string  `json:"category"`
	StockValue float64 `json:"stock_value"`
}

// Products with low stock.
type LowStockItem struct {
	ProductID   int64   `json:"product_id"`
	ProductName string  `json:"product_name"`
	SKU         string  `json:"sku"`
	StockBalance float64 `json:"stock_balance"`
}

// Latest stock transactions.
type RecentTransaction struct {
	TransactionNo   string    `json:"transaction_no"`
	Type            string    `json:"type"`
	WarehouseName   string    `json:"warehouse_name"`
	TransactionDate time.Time `json:"transaction_date"`
	TotalAmount     float64   `json:"total_amount"`
}
