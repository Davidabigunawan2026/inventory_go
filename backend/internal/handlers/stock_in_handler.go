
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"inventory_go/repositories"
)

type StockInHandler struct {
	Repo *repositories.StockInRepository
}

func NewStockInHandler(
	repo *repositories.StockInRepository,
) *StockInHandler {
	return &StockInHandler{Repo: repo}
}

func (h *StockInHandler) GetAllStockIn(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()

	transactions, err := h.Repo.GetAllStockIn(ctx)
	if err != nil {
		log.Println("Get Stock In Error:", err)

		http.Error(
			w,
			"Failed to retrieve stock in transactions",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"data":      transactions,
		"status":    "success",
		"timestamp": time.Now(),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("JSON Encoding Error:", err)
	}
}
