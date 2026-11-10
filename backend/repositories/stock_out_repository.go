package repositories


import (
    "context"

    "inventory_go/models"

    "github.com/jackc/pgx/v5/pgxpool"
)


type StockOutRepository struct {
	DB *pgxpool.Pool
}

func NewStockOutRepository(db *pgxpool.Pool) *StockOutRepository {
	return &StockOutRepository{DB: db}
}

func (r *StockOutRepository) GetAllStockOut(
	ctx context.Context,
) ([]models.StockOutTransaction, error) {

	query := `
		SELECT
			so.id,
			so.transaction_no,
			so.transaction_date,
			so.warehouse_id,
			w.warehouse_name,
			so.reference_no,
			so.destination,
			so.notes,
			COALESCE(SUM(d.quantity), 0)::double precision
		FROM stock_out so
		JOIN warehouses w
			ON w.id = so.warehouse_id
		LEFT JOIN stock_out_details d
			ON d.stock_out_id = so.id
		GROUP BY
			so.id,
			so.transaction_no,
			so.transaction_date,
			so.warehouse_id,
			w.warehouse_name,
			so.reference_no,
			so.destination,
			so.notes
		ORDER BY so.transaction_date DESC, so.id DESC
	`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := []models.StockOutTransaction{}

	for rows.Next() {
		var item models.StockOutTransaction

		err := rows.Scan(
			&item.ID,
			&item.TransactionNo,
			&item.TransactionDate,
			&item.WarehouseID,
			&item.WarehouseName,
			&item.ReferenceNo,
			&item.Destination,
			&item.Notes,
			&item.TotalQuantity,
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
