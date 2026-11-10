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

type ProductHandler struct {
	Repo *repositories.ProductRepository
}

func NewProductHandler(repo *repositories.ProductRepository) *ProductHandler {
	return &ProductHandler{
		Repo: repo,
	}
}

// JSON RESPONSE
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// ERROR RESPONSE
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{
		"status":  "error",
		"message": message,
	})
}

// DATABASE ERROR HANDLING
func handleDatabaseError(w http.ResponseWriter, err error) {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			respondError(w, http.StatusConflict, "SKU already exists")
			return
		case "23514", "23502":
			respondError(w, http.StatusBadRequest, "Invalid product data")
			return
		}
	}

	if errors.Is(err, pgx.ErrNoRows) {
		respondError(w, http.StatusNotFound, "Product not found")
		return
	}

	respondError(w, http.StatusInternalServerError, "Database operation failed")
}

// DECODE JSON
func decodeProduct(w http.ResponseWriter, r *http.Request) (models.Product, error) {
	var p models.Product

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&p); err != nil {
		return p, err
	}

	var extra interface{}

	if err := decoder.Decode(&extra); err != io.EOF {
		return p, errors.New("request must contain only one JSON object")
	}

	p.SKU = strings.TrimSpace(p.SKU)
	p.Name = strings.TrimSpace(p.Name)
	p.Unit = strings.TrimSpace(p.Unit)

	if p.SKU == "" {
		return p, errors.New("SKU is required")
	}

	if p.Name == "" {
		return p, errors.New("product name is required")
	}

	if len(p.SKU) > 50 || len(p.Name) > 150 {
		return p, errors.New("SKU or product name is too long")
	}

	if p.Unit == "" {
		p.Unit = "PCS"
	}

	if p.MinStock < 0 {
		return p, errors.New("minimum stock cannot be negative")
	}

	return p, nil
}

// GET ALL PRODUCTS
func (h *ProductHandler) GetAllProducts(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	products, err := h.Repo.GetAllProducts(r.Context())

	if err != nil {
		handleDatabaseError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   products,
	})
}

// GET PRODUCT BY ID
func (h *ProductHandler) GetProductByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := h.Repo.GetProductByID(r.Context(), id)

	if err != nil {
		handleDatabaseError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   product,
	})
}

// CREATE PRODUCT
func (h *ProductHandler) CreateProduct(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	p, err := decodeProduct(w, r)

	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	product, err := h.Repo.CreateProduct(r.Context(), &p)

	if err != nil {
		handleDatabaseError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"status":  "success",
		"message": "Product created successfully",
		"data":    product,
	})
}

// UPDATE PRODUCT
func (h *ProductHandler) UpdateProduct(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPut {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	p, err := decodeProduct(w, r)

	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	p.ID = id

	product, err := h.Repo.UpdateProduct(r.Context(), &p)

	if err != nil {
		handleDatabaseError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Product updated successfully",
		"data":    product,
	})
}

// DELETE PRODUCT
func (h *ProductHandler) DeleteProduct(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	err = h.Repo.DeleteProduct(r.Context(), id)

	if err != nil {
		handleDatabaseError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Product deleted successfully",
	})
}
