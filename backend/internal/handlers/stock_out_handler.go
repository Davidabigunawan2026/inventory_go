package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"inventory_go/repositories"
)

type StockOutHandler struct {
	Repo *repositories.StockOutRepository
}

func NewStockOutHandler(
	repo *repositories.StockOutRepository,
) *StockOutHandler {
	return &StockOutHandler{Repo: repo}
}

func (h *StockOutHandler) GetAllStockOut(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()

	transactions, err := h.Repo.GetAllStockOut(ctx)
	if err != nil {
		log.Println("Get Stock Out Error:", err)

		http.Error(
			w,
			"Failed to retrieve stock out transactions",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"data":      transactions,
		"status":    "success",
		"timestamp": time.Now(),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("JSON Encoding Error:", err)
	}
}
