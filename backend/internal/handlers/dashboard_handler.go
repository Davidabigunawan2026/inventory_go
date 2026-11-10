
package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"inventory_go/repositories"
)

type DashboardHandler struct {
	Repository *repositories.DashboardRepository
}

func NewDashboardHandler(
	repo *repositories.DashboardRepository,
) *DashboardHandler {
	return &DashboardHandler{
		Repository: repo,
	}
}

// GET /api/dashboard
func (h *DashboardHandler) GetDashboard(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx, cancel := context.WithTimeout(
		r.Context(),
		10*time.Second,
	)
	defer cancel()

	data, err := h.Repository.GetDashboardData(ctx)
	if err != nil {
		http.Error(
			w,
			"Failed to retrieve dashboard data: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "success",
		"data":      data,
		"timestamp": time.Now(),
	})
}
