package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"inventory_go/models"
	"inventory_go/repositories"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type WarehouseHandler struct {
	Repo *repositories.WarehouseRepository
}

func NewWarehouseHandler(
	repo *repositories.WarehouseRepository,
) *WarehouseHandler {
	return &WarehouseHandler{
		Repo: repo,
	}
}

// GET ALL WAREHOUSES
func (h *WarehouseHandler) GetAllWarehouses(
	w http.ResponseWriter,
	r *http.Request,
) {
	warehouses, err := h.Repo.GetAllWarehouses(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to retrieve warehouses")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"data":   warehouses,
		"status": "success",
	})
}

// GET WAREHOUSE BY ID
func (h *WarehouseHandler) GetWarehouseByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid warehouse ID")
		return
	}

	warehouse, err := h.Repo.GetWarehouseByID(r.Context(), id)
	if err != nil {
		handleWarehouseError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"data":   warehouse,
		"status": "success",
	})
}

// CREATE WAREHOUSE
func (h *WarehouseHandler) CreateWarehouse(
	w http.ResponseWriter,
	r *http.Request,
) {
	var warehouse models.Warehouse

	if err := decodeWarehouse(w, r, &warehouse); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	warehouse.IsActive = true

	result, err := h.Repo.CreateWarehouse(r.Context(), warehouse)
	if err != nil {
		handleWarehouseError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"data":    result,
		"status":  "success",
		"message": "Warehouse created successfully",
	})
}

// UPDATE WAREHOUSE
func (h *WarehouseHandler) UpdateWarehouse(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid warehouse ID")
		return
	}

	var warehouse models.Warehouse

	if err := decodeWarehouse(w, r, &warehouse); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	warehouse.ID = id

	result, err := h.Repo.UpdateWarehouse(r.Context(), warehouse)
	if err != nil {
		handleWarehouseError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"data":    result,
		"status":  "success",
		"message": "Warehouse updated successfully",
	})
}

// DELETE WAREHOUSE
func (h *WarehouseHandler) DeleteWarehouse(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid warehouse ID")
		return
	}

	if err := h.Repo.DeleteWarehouse(r.Context(), id); err != nil {
		handleWarehouseError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Warehouse deleted successfully",
	})
}

// DECODE AND VALIDATE JSON
func decodeWarehouse(
	w http.ResponseWriter,
	r *http.Request,
	warehouse *models.Warehouse,
) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(warehouse); err != nil {
		return errors.New("invalid JSON request body")
	}

	var extra interface{}
	if err := decoder.Decode(&extra); err != io.EOF {
		return errors.New("request must contain one JSON object")
	}

	warehouse.WarehouseCode = strings.TrimSpace(warehouse.WarehouseCode)
	warehouse.WarehouseName = strings.TrimSpace(warehouse.WarehouseName)

	if warehouse.WarehouseCode == "" {
		return errors.New("warehouse code is required")
	}

	if warehouse.WarehouseName == "" {
		return errors.New("warehouse name is required")
	}

	if len(warehouse.WarehouseCode) > 30 {
		return errors.New("warehouse code must not exceed 30 characters")
	}

	if len(warehouse.WarehouseName) > 150 {
		return errors.New("warehouse name must not exceed 150 characters")
	}

	return nil
}

// DATABASE ERROR HANDLER
func handleWarehouseError(w http.ResponseWriter, err error) {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			respondError(w, http.StatusConflict, "Warehouse code already exists")
			return
		case "23503":
			respondError(w, http.StatusConflict, "Warehouse is referenced by another record")
			return
		case "23502", "23514":
			respondError(w, http.StatusBadRequest, "Invalid warehouse data")
			return
		}
	}

	if errors.Is(err, pgx.ErrNoRows) {
		respondError(w, http.StatusNotFound, "Warehouse not found")
		return
	}

	respondError(w, http.StatusInternalServerError, "Warehouse database operation failed")
}
