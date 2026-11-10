
package repositories

import (
	"context"

	"inventory_go/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardRepository struct {
	DB *pgxpool.Pool
}

func NewDashboardRepository(db *pgxpool.Pool) *DashboardRepository {
	return &DashboardRepository{DB: db}
}

// GetDashboardData retrieves all dashboard information.
func (r *DashboardRepository) GetDashboardData(
	ctx context.Context,
) (*models.DashboardData, error) {

	data := &models.DashboardData{
		MonthlyMovement:    []models.MonthlyMovement{},
		CategoryValues:     []models.CategoryValue{},
		LowStockItems:      []models.LowStockItem{},
		RecentTransactions: []models.RecentTransaction{},
	}

	// 1. TOTAL PRODUCTS
	err := r.DB.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM products
		WHERE is_active = TRUE
	`).Scan(&data.TotalProducts)
	if err != nil {
		return nil, err
	}

	// 2. TOTAL WAREHOUSES
	err = r.DB.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM warehouses
		WHERE is_active = TRUE
	`).Scan(&data.TotalWarehouses)
	if err != nil {
		return nil, err
	}

	// 3. TOTAL STOCK IN
	err = r.DB.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(d.quantity), 0)::double precision,
			COALESCE(
				SUM(d.quantity * d.unit_cost), 0
			)::double precision
		FROM stock_in_details d
		JOIN stock_in si ON si.id = d.stock_in_id
	`).Scan(
		&data.TotalStockInQuantity,
		&data.TotalStockInAmount,
	)
	if err != nil {
		return nil, err
	}

	// 4. TOTAL STOCK OUT
	err = r.DB.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(d.quantity), 0)::double precision,
			COALESCE(
				SUM(d.quantity * COALESCE(d.unit_cost, 0)), 0
			)::double precision
		FROM stock_out_details d
		JOIN stock_out so ON so.id = d.stock_out_id
	`).Scan(
		&data.TotalStockOutQuantity,
		&data.TotalStockOutAmount,
	)
	if err != nil {
		return nil, err
	}

	// 5. INVENTORY VALUE
	// Uses weighted average inbound unit cost per product.
	err = r.DB.QueryRow(ctx, `
		WITH inbound AS (
			SELECT
				product_id,
				COALESCE(SUM(quantity), 0) AS in_qty,
				COALESCE(SUM(quantity * unit_cost), 0) AS in_value
			FROM stock_in_details
			GROUP BY product_id
		),
		outbound AS (
			SELECT
				product_id,
				COALESCE(SUM(quantity), 0) AS out_qty
			FROM stock_out_details
			GROUP BY product_id
		),
		balance AS (
			SELECT
				COALESCE(i.product_id, o.product_id) AS product_id,
				COALESCE(i.in_qty, 0) -
					COALESCE(o.out_qty, 0) AS stock_balance,
				CASE
					WHEN COALESCE(i.in_qty, 0) > 0
					THEN i.in_value / i.in_qty
					ELSE 0
				END AS average_cost
			FROM inbound i
			FULL OUTER JOIN outbound o
				ON o.product_id = i.product_id
		)
		SELECT COALESCE(
			SUM(stock_balance * average_cost), 0
		)::double precision
		FROM balance
	`).Scan(&data.TotalInventoryValue)
	if err != nil {
		return nil, err
	}

	// 6. MONTHLY STOCK MOVEMENT
	// Displays the current month and previous five months.
	monthlyQuery := `
		WITH months AS (
			SELECT generate_series(
				date_trunc('month', CURRENT_DATE) - INTERVAL '5 months',
				date_trunc('month', CURRENT_DATE),
				INTERVAL '1 month'
			) AS month_start
		),
		stock_in_monthly AS (
			SELECT
				date_trunc('month', si.transaction_date) AS month_start,
				SUM(d.quantity * d.unit_cost) AS amount
			FROM stock_in si
			JOIN stock_in_details d
				ON d.stock_in_id = si.id
			GROUP BY 1
		),
		stock_out_monthly AS (
			SELECT
				date_trunc('month', so.transaction_date) AS month_start,
				SUM(d.quantity * COALESCE(d.unit_cost, 0)) AS amount
			FROM stock_out so
			JOIN stock_out_details d
				ON d.stock_out_id = so.id
			GROUP BY 1
		)
		SELECT
			TO_CHAR(m.month_start, 'Mon YYYY'),
			COALESCE(si.amount, 0)::double precision,
			COALESCE(so.amount, 0)::double precision
		FROM months m
		LEFT JOIN stock_in_monthly si
			ON si.month_start = m.month_start
		LEFT JOIN stock_out_monthly so
			ON so.month_start = m.month_start
		ORDER BY m.month_start
	`

	rows, err := r.DB.Query(ctx, monthlyQuery)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var item models.MonthlyMovement

		if err := rows.Scan(
			&item.Month,
			&item.StockInAmount,
			&item.StockOutAmount,
		); err != nil {
			rows.Close()
			return nil, err
		}

		data.MonthlyMovement = append(data.MonthlyMovement, item)
	}

	err = rows.Err()
	rows.Close()
	if err != nil {
		return nil, err
	}

	// 7. STOCK VALUE BY CATEGORY
	categoryQuery := `
		WITH inbound AS (
			SELECT
				product_id,
				SUM(quantity) AS in_qty,
				SUM(quantity * unit_cost) AS in_value
			FROM stock_in_details
			GROUP BY product_id
		),
		outbound AS (
			SELECT
				product_id,
				SUM(quantity) AS out_qty
			FROM stock_out_details
			GROUP BY product_id
		)
		SELECT
			COALESCE(NULLIF(p.category, ''), 'Uncategorized'),
			COALESCE(
				SUM(
					(
						COALESCE(i.in_qty, 0) -
						COALESCE(o.out_qty, 0)
					) *
					CASE
						WHEN COALESCE(i.in_qty, 0) > 0
						THEN i.in_value / i.in_qty
						ELSE 0
					END
				), 0
			)::double precision
		FROM products p
		LEFT JOIN inbound i ON i.product_id = p.id
		LEFT JOIN outbound o ON o.product_id = p.id
		WHERE p.is_active = TRUE
		GROUP BY p.category
		ORDER BY p.category
	`

	rows, err = r.DB.Query(ctx, categoryQuery)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var item models.CategoryValue

		if err := rows.Scan(
			&item.Category,
			&item.StockValue,
		); err != nil {
			rows.Close()
			return nil, err
		}

		data.CategoryValues = append(data.CategoryValues, item)
	}

	err = rows.Err()
	rows.Close()
	if err != nil {
		return nil, err
	}

	// 8. LOW STOCK ITEMS
	lowStockQuery := `
		SELECT
			p.id,
			p.name,
			p.sku,
			(
				COALESCE(si.total_in, 0) -
				COALESCE(so.total_out, 0)
			)::double precision
		FROM products p
		LEFT JOIN (
			SELECT product_id, SUM(quantity) AS total_in
			FROM stock_in_details
			GROUP BY product_id
		) si ON si.product_id = p.id
		LEFT JOIN (
			SELECT product_id, SUM(quantity) AS total_out
			FROM stock_out_details
			GROUP BY product_id
		) so ON so.product_id = p.id
		WHERE p.is_active = TRUE
		  AND (
			COALESCE(si.total_in, 0) -
			COALESCE(so.total_out, 0)
		  ) <= p.min_stock
		ORDER BY
			(
				COALESCE(si.total_in, 0) -
				COALESCE(so.total_out, 0)
			) ASC
	`

	rows, err = r.DB.Query(ctx, lowStockQuery)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var item models.LowStockItem

		if err := rows.Scan(
			&item.ProductID,
			&item.ProductName,
			&item.SKU,
			&item.StockBalance,
		); err != nil {
			rows.Close()
			return nil, err
		}

		data.LowStockItems = append(data.LowStockItems, item)
	}

	err = rows.Err()
	rows.Close()
	if err != nil {
		return nil, err
	}

	// 9. RECENT TRANSACTIONS
	recentQuery := `
		SELECT
			t.transaction_no,
			t.type,
			t.warehouse_name,
			t.transaction_date,
			t.total_amount
		FROM (
			SELECT
				si.transaction_no,
				'Stock In' AS type,
				w.warehouse_name,
				si.transaction_date,
				COALESCE(
					SUM(d.quantity * d.unit_cost), 0
				)::double precision AS total_amount
			FROM stock_in si
			JOIN warehouses w ON w.id = si.warehouse_id
			LEFT JOIN stock_in_details d
				ON d.stock_in_id = si.id
			GROUP BY
				si.id, si.transaction_no,
				w.warehouse_name, si.transaction_date

			UNION ALL

			SELECT
				so.transaction_no,
				'Stock Out' AS type,
				w.warehouse_name,
				so.transaction_date,
				COALESCE(
					SUM(d.quantity * COALESCE(d.unit_cost, 0)), 0
				)::double precision AS total_amount
			FROM stock_out so
			JOIN warehouses w ON w.id = so.warehouse_id
			LEFT JOIN stock_out_details d
				ON d.stock_out_id = so.id
			GROUP BY
				so.id, so.transaction_no,
				w.warehouse_name, so.transaction_date
		) t
		ORDER BY t.transaction_date DESC
		LIMIT 10
	`

	rows, err = r.DB.Query(ctx, recentQuery)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var item models.RecentTransaction

		if err := rows.Scan(
			&item.TransactionNo,
			&item.Type,
			&item.WarehouseName,
			&item.TransactionDate,
			&item.TotalAmount,
		); err != nil {
			rows.Close()
			return nil, err
		}

		data.RecentTransactions = append(
			data.RecentTransactions,
			item,
		)
	}

	err = rows.Err()
	rows.Close()
	if err != nil {
		return nil, err
	}

	return data, nil
}
