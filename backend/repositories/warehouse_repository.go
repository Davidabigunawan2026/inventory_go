package repositories

import (
	"context"
	"fmt"

	"inventory_go/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WarehouseRepository struct {
	DB *pgxpool.Pool
}

func NewWarehouseRepository(db *pgxpool.Pool) *WarehouseRepository {
	return &WarehouseRepository{
		DB: db,
	}
}

// GET ALL WAREHOUSES
func (r *WarehouseRepository) GetAllWarehouses(
	ctx context.Context,
) ([]models.Warehouse, error) {

	query := `
		SELECT
			id,
			warehouse_code,
			warehouse_name,
			address,
			city,
			contact_person,
			phone,
			email,
			is_active,
			created_at,
			updated_at
		FROM warehouses
		ORDER BY id ASC
	`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	warehouses := []models.Warehouse{}

	for rows.Next() {
		var w models.Warehouse

		err := rows.Scan(
			&w.ID,
			&w.WarehouseCode,
			&w.WarehouseName,
			&w.Address,
			&w.City,
			&w.ContactPerson,
			&w.Phone,
			&w.Email,
			&w.IsActive,
			&w.CreatedAt,
			&w.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		warehouses = append(warehouses, w)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return warehouses, nil
}

// GET WAREHOUSE BY ID
func (r *WarehouseRepository) GetWarehouseByID(
	ctx context.Context,
	id int64,
) (models.Warehouse, error) {

	var w models.Warehouse

	query := `
		SELECT
			id,
			warehouse_code,
			warehouse_name,
			address,
			city,
			contact_person,
			phone,
			email,
			is_active,
			created_at,
			updated_at
		FROM warehouses
		WHERE id = $1
	`

	err := r.DB.QueryRow(ctx, query, id).Scan(
		&w.ID,
		&w.WarehouseCode,
		&w.WarehouseName,
		&w.Address,
		&w.City,
		&w.ContactPerson,
		&w.Phone,
		&w.Email,
		&w.IsActive,
		&w.CreatedAt,
		&w.UpdatedAt,
	)

	if err != nil {
		return w, err
	}

	return w, nil
}

// CREATE WAREHOUSE
func (r *WarehouseRepository) CreateWarehouse(
	ctx context.Context,
	w models.Warehouse,
) (models.Warehouse, error) {

	query := `
		INSERT INTO warehouses (
			warehouse_code,
			warehouse_name,
			address,
			city,
			contact_person,
			phone,
			email,
			is_active
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING
			id,
			warehouse_code,
			warehouse_name,
			address,
			city,
			contact_person,
			phone,
			email,
			is_active,
			created_at,
			updated_at
	`

	err := r.DB.QueryRow(
		ctx,
		query,
		w.WarehouseCode,
		w.WarehouseName,
		w.Address,
		w.City,
		w.ContactPerson,
		w.Phone,
		w.Email,
		w.IsActive,
	).Scan(
		&w.ID,
		&w.WarehouseCode,
		&w.WarehouseName,
		&w.Address,
		&w.City,
		&w.ContactPerson,
		&w.Phone,
		&w.Email,
		&w.IsActive,
		&w.CreatedAt,
		&w.UpdatedAt,
	)

	if err != nil {
		return w, err
	}

	return w, nil
}

// UPDATE WAREHOUSE
func (r *WarehouseRepository) UpdateWarehouse(
	ctx context.Context,
	w models.Warehouse,
) (models.Warehouse, error) {

	query := `
		UPDATE warehouses
		SET
			warehouse_code = $1,
			warehouse_name = $2,
			address = $3,
			city = $4,
			contact_person = $5,
			phone = $6,
			email = $7,
			is_active = $8,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $9
		RETURNING
			id,
			warehouse_code,
			warehouse_name,
			address,
			city,
			contact_person,
			phone,
			email,
			is_active,
			created_at,
			updated_at
	`

	err := r.DB.QueryRow(
		ctx,
		query,
		w.WarehouseCode,
		w.WarehouseName,
		w.Address,
		w.City,
		w.ContactPerson,
		w.Phone,
		w.Email,
		w.IsActive,
		w.ID,
	).Scan(
		&w.ID,
		&w.WarehouseCode,
		&w.WarehouseName,
		&w.Address,
		&w.City,
		&w.ContactPerson,
		&w.Phone,
		&w.Email,
		&w.IsActive,
		&w.CreatedAt,
		&w.UpdatedAt,
	)

	if err != nil {
		return w, err
	}

	return w, nil
}

// DELETE WAREHOUSE
func (r *WarehouseRepository) DeleteWarehouse(
	ctx context.Context,
	id int64,
) error {

	query := `
		DELETE FROM warehouses
		WHERE id = $1
	`

	result, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf(
			"warehouse not found: %w",
			pgx.ErrNoRows,
		)
	}

	return nil
}
