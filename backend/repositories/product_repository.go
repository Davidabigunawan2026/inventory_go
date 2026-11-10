package repositories

import (
	"context"
	"fmt"

	"inventory_go/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	DB *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{
		DB: db,
	}
}

// GET ALL PRODUCTS
func (r *ProductRepository) GetAllProducts(
	ctx context.Context,
) ([]models.Product, error) {

	query := `
		SELECT
			id, sku, name, description, category,
			unit, min_stock, is_active, created_at, updated_at
		FROM products
		ORDER BY id ASC
	`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get products: %w", err)
	}
	defer rows.Close()

	products := make([]models.Product, 0)

	for rows.Next() {
		var p models.Product

		err := rows.Scan(
			&p.ID,
			&p.SKU,
			&p.Name,
			&p.Description,
			&p.Category,
			&p.Unit,
			&p.MinStock,
			&p.IsActive,
			&p.CreatedAt,
			&p.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}

		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate products: %w", err)
	}

	return products, nil
}

// GET PRODUCT BY ID
func (r *ProductRepository) GetProductByID(
	ctx context.Context,
	id int64,
) (*models.Product, error) {

	query := `
		SELECT
			id, sku, name, description, category,
			unit, min_stock, is_active, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	var p models.Product

	err := r.DB.QueryRow(ctx, query, id).Scan(
		&p.ID,
		&p.SKU,
		&p.Name,
		&p.Description,
		&p.Category,
		&p.Unit,
		&p.MinStock,
		&p.IsActive,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("get product by ID: %w", err)
	}

	return &p, nil
}

// CREATE PRODUCT
func (r *ProductRepository) CreateProduct(
	ctx context.Context,
	p *models.Product,
) (*models.Product, error) {

	query := `
		INSERT INTO products (
			sku, name, description, category,
			unit, min_stock, is_active
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING
			id, sku, name, description, category,
			unit, min_stock, is_active, created_at, updated_at
	`

	var created models.Product

	err := r.DB.QueryRow(
		ctx,
		query,
		p.SKU,
		p.Name,
		p.Description,
		p.Category,
		p.Unit,
		p.MinStock,
		p.IsActive,
	).Scan(
		&created.ID,
		&created.SKU,
		&created.Name,
		&created.Description,
		&created.Category,
		&created.Unit,
		&created.MinStock,
		&created.IsActive,
		&created.CreatedAt,
		&created.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("create product: %w", err)
	}

	return &created, nil
}

// UPDATE PRODUCT
func (r *ProductRepository) UpdateProduct(
	ctx context.Context,
	p *models.Product,
) (*models.Product, error) {

	query := `
		UPDATE products
		SET
			sku = $1,
			name = $2,
			description = $3,
			category = $4,
			unit = $5,
			min_stock = $6,
			is_active = $7,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
		RETURNING
			id, sku, name, description, category,
			unit, min_stock, is_active, created_at, updated_at
	`

	var updated models.Product

	err := r.DB.QueryRow(
		ctx,
		query,
		p.SKU,
		p.Name,
		p.Description,
		p.Category,
		p.Unit,
		p.MinStock,
		p.IsActive,
		p.ID,
	).Scan(
		&updated.ID,
		&updated.SKU,
		&updated.Name,
		&updated.Description,
		&updated.Category,
		&updated.Unit,
		&updated.MinStock,
		&updated.IsActive,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("update product: %w", err)
	}

	return &updated, nil
}

// DELETE PRODUCT
func (r *ProductRepository) DeleteProduct(
	ctx context.Context,
	id int64,
) error {

	query := `
		DELETE FROM products
		WHERE id = $1
	`

	result, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("product with ID %d not found: %w", id, pgx.ErrNoRows)
	}

	return nil
}
