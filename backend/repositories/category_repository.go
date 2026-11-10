package repositories

import (
	"context"
	"fmt"

	"inventory_go/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository struct {
	DB *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

func (r *CategoryRepository) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	query := `
		SELECT id, name, description, is_active, created_at, updated_at
		FROM categories
		ORDER BY id ASC
	`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get categories: %w", err)
	}
	defer rows.Close()

	categories := make([]models.Category, 0)

	for rows.Next() {
		var c models.Category

		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Description,
			&c.IsActive,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan category: %w", err)
		}

		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate categories: %w", err)
	}

	return categories, nil
}

func (r *CategoryRepository) GetCategoryByID(ctx context.Context, id int64) (*models.Category, error) {
	query := `
		SELECT id, name, description, is_active, created_at, updated_at
		FROM categories
		WHERE id = $1
	`

	var c models.Category

	err := r.DB.QueryRow(ctx, query, id).Scan(
		&c.ID,
		&c.Name,
		&c.Description,
		&c.IsActive,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("get category by ID: %w", err)
	}

	return &c, nil
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, c *models.Category) (*models.Category, error) {
	query := `
		INSERT INTO categories (name, description, is_active)
		VALUES ($1, $2, $3)
		RETURNING id, name, description, is_active, created_at, updated_at
	`

	var created models.Category

	err := r.DB.QueryRow(
		ctx, query, c.Name, c.Description, c.IsActive,
	).Scan(
		&created.ID,
		&created.Name,
		&created.Description,
		&created.IsActive,
		&created.CreatedAt,
		&created.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("create category: %w", err)
	}

	return &created, nil
}

func (r *CategoryRepository) UpdateCategory(ctx context.Context, c *models.Category) (*models.Category, error) {
	query := `
		UPDATE categories
		SET name = $1,
			description = $2,
			is_active = $3,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING id, name, description, is_active, created_at, updated_at
	`

	var updated models.Category

	err := r.DB.QueryRow(
		ctx, query, c.Name, c.Description, c.IsActive, c.ID,
	).Scan(
		&updated.ID,
		&updated.Name,
		&updated.Description,
		&updated.IsActive,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("update category: %w", err)
	}

	return &updated, nil
}

func (r *CategoryRepository) DeleteCategory(ctx context.Context, id int64) error {
	query := `DELETE FROM categories WHERE id = $1`

	result, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("category with ID %d not found: %w", id, pgx.ErrNoRows)
	}

	return nil
}
