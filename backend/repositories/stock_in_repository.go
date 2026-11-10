
package repositories

import (
    "context"

    "inventory_go/models"

    "github.com/jackc/pgx/v5/pgxpool"
)


type StockInRepository struct {
	DB *pgxpool.Pool
}

func NewStockInRepository(db *pgxpool.Pool) *StockInRepository {
	return &StockInRepository{DB: db}
}

// GetAllStockIn retrieves all stock in transactions
func (r *StockInRepository) GetAllStockIn(
	ctx context.Context,
) ([]models.StockInTransaction, error) {

	query := `
		SELECT
			si.id,
			si.transaction_no,
			si.transaction_date,
			si.warehouse_id,
			w.warehouse_name,
			si.reference_no,
			si.supplier_name,
			si.notes,
			COALESCE(SUM(d.quantity), 0)::double precision,
			COALESCE(
				SUM(d.quantity * d.unit_cost), 0
			)::double precision,
			si.created_at,
			si.updated_at
		FROM stock_in si
		JOIN warehouses w
			ON w.id = si.warehouse_id
		LEFT JOIN stock_in_details d
			ON d.stock_in_id = si.id
		GROUP BY
			si.id,
			si.transaction_no,
			si.transaction_date,
			si.warehouse_id,
			w.warehouse_name,
			si.reference_no,
			si.supplier_name,
			si.notes,
			si.created_at,
			si.updated_at
		ORDER BY
			si.transaction_date DESC,
			si.id DESC
	`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := []models.StockInTransaction{}

	for rows.Next() {
		var item models.StockInTransaction

		err := rows.Scan(
			&item.ID,
			&item.TransactionNo,
			&item.TransactionDate,
			&item.WarehouseID,
			&item.WarehouseName,
			&item.ReferenceNo,
			&item.SupplierName,
			&item.Notes,
			&item.TotalQuantity,
			&item.TotalAmount,
			&item.CreatedAt,
			&item.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		transactions = append(transactions, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}



// GetStockInByID retrieves one stock in transaction with its details
func (r *StockInRepository) GetStockInByID(
	ctx context.Context,
	id int64,
) (*models.StockInTransaction, error) {

	query := `
		SELECT
			si.id,
			si.transaction_no,
			si.transaction_date,
			si.warehouse_id,
			w.warehouse_name,
			si.reference_no,
			si.supplier_name,
			si.notes,
			COALESCE(SUM(d.quantity), 0)::double precision,
			COALESCE(
				SUM(d.quantity * d.unit_cost), 0
			)::double precision,
			si.created_at,
			si.updated_at
		FROM stock_in si
		JOIN warehouses w
			ON w.id = si.warehouse_id
		LEFT JOIN stock_in_details d
			ON d.stock_in_id = si.id
		WHERE si.id = $1
		GROUP BY
			si.id,
			w.warehouse_name
	`

	var item models.StockInTransaction

	err := r.DB.QueryRow(ctx, query, id).Scan(
		&item.ID,
		&item.TransactionNo,
		&item.TransactionDate,
		&item.WarehouseID,
		&item.WarehouseName,
		&item.ReferenceNo,
		&item.SupplierName,
		&item.Notes,
		&item.TotalQuantity,
		&item.TotalAmount,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	detailQuery := `
		SELECT
			d.id,
			d.stock_in_id,
			d.product_id,
			p.product_name,
			d.quantity,
			d.unit_cost,
			(d.quantity * d.unit_cost)::double precision,
			d.created_at
		FROM stock_in_details d
		JOIN products p
			ON p.id = d.product_id
		WHERE d.stock_in_id = $1
		ORDER BY d.id ASC
	`

	rows, err := r.DB.Query(ctx, detailQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	item.Details = []models.StockInDetail{}

	for rows.Next() {
		var detail models.StockInDetail

		err := rows.Scan(
			&detail.ID,
			&detail.StockInID,
			&detail.ProductID,
			&detail.ProductName,
			&detail.Quantity,
			&detail.UnitCost,
			&detail.TotalCost,
			&detail.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		item.Details = append(item.Details, detail)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &item, nil
}


