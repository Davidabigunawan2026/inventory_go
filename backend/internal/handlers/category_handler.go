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

type CategoryHandler struct {
	Repo *repositories.CategoryRepository
}

func NewCategoryHandler(repo *repositories.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{Repo: repo}
}

// DECODE CATEGORY JSON
func decodeCategory(w http.ResponseWriter, r *http.Request) (models.Category, error) {
	var c models.Category

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&c); err != nil {
		return c, err
	}

	var extra interface{}
	if err := decoder.Decode(&extra); err != io.EOF {
		return c, errors.New("request must contain only one JSON object")
	}

	c.Name = strings.TrimSpace(c.Name)

	if c.Name == "" {
		return c, errors.New("category name is required")
	}

	if len(c.Name) > 100 {
		return c, errors.New("category name is too long")
	}

	if c.Description != nil {
		description := strings.TrimSpace(*c.Description)
		c.Description = &description

		if description == "" {
			c.Description = nil
		}
	}

	return c, nil
}

// CATEGORY DATABASE ERROR
func handleCategoryDatabaseError(w http.ResponseWriter, err error) {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			respondError(w, http.StatusConflict, "Category name already exists")
			return
		case "23514", "23502":
			respondError(w, http.StatusBadRequest, "Invalid category data")
			return
		case "23503":
			respondError(w, http.StatusConflict, "Category is being used by another record")
			return
		}
	}

	if errors.Is(err, pgx.ErrNoRows) {
		respondError(w, http.StatusNotFound, "Category not found")
		return
	}

	respondError(w, http.StatusInternalServerError, "Database operation failed")
}

// GET ALL CATEGORIES
func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.Repo.GetAllCategories(r.Context())
	if err != nil {
		handleCategoryDatabaseError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   categories,
	})
}

// GET CATEGORY BY ID
func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	category, err := h.Repo.GetCategoryByID(r.Context(), id)
	if err != nil {
		handleCategoryDatabaseError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   category,
	})
}

// CREATE CATEGORY
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	c, err := decodeCategory(w, r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// New category defaults to active.
	c.IsActive = true

	category, err := h.Repo.CreateCategory(r.Context(), &c)
	if err != nil {
		handleCategoryDatabaseError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"status":  "success",
		"message": "Category created successfully",
		"data":    category,
	})
}

// UPDATE CATEGORY
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	c, err := decodeCategory(w, r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	c.ID = id

	category, err := h.Repo.UpdateCategory(r.Context(), &c)
	if err != nil {
		handleCategoryDatabaseError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Category updated successfully",
		"data":    category,
	})
}

// DELETE CATEGORY
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	err = h.Repo.DeleteCategory(r.Context(), id)
	if err != nil {
		handleCategoryDatabaseError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Category deleted successfully",
	})
}
