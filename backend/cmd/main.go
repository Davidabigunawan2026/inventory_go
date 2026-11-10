package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"inventory_go/internal/database"
	"inventory_go/internal/handlers"
	"inventory_go/repositories"
)

func main() {

	// 1. CONNECT POSTGRESQL
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	log.Println("PostgreSQL Connected!")

	// 2. CREATE ROUTER
	mux := http.NewServeMux()

	// 3. PRODUCT REPOSITORY & HANDLER
	productRepo := repositories.NewProductRepository(db)
	productHandler := handlers.NewProductHandler(productRepo)

	// 4. CATEGORY REPOSITORY & HANDLER
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)

	// 5. WAREHOUSE REPOSITORY & HANDLER
	warehouseRepo := repositories.NewWarehouseRepository(db)
	warehouseHandler := handlers.NewWarehouseHandler(warehouseRepo)

	// STOCK IN REPOSITORY & HANDLER
	stockInRepo := repositories.NewStockInRepository(db)
	stockInHandler := handlers.NewStockInHandler(stockInRepo)

	// 6. STOCK OUT REPOSITORY & HANDLER
	stockOutRepo := repositories.NewStockOutRepository(db)
	stockOutHandler := handlers.NewStockOutHandler(stockOutRepo)

	// DASHBOARD REPOSITORY & HANDLER
	dashboardRepo := repositories.NewDashboardRepository(db)
	dashboardHandler := handlers.NewDashboardHandler(dashboardRepo)




	// 6. PRODUCT API ROUTES
	mux.HandleFunc("GET /api/products",
		productHandler.GetAllProducts)
	mux.HandleFunc("GET /api/products/{id}",
		productHandler.GetProductByID)
	mux.HandleFunc("POST /api/products",
		productHandler.CreateProduct)
	mux.HandleFunc("PUT /api/products/{id}",
		productHandler.UpdateProduct)
	mux.HandleFunc("DELETE /api/products/{id}",
		productHandler.DeleteProduct)

	// 7. CATEGORY API ROUTES
	mux.HandleFunc("GET /api/categories",
		categoryHandler.GetAllCategories)
	mux.HandleFunc("GET /api/categories/{id}",
		categoryHandler.GetCategoryByID)
	mux.HandleFunc("POST /api/categories",
		categoryHandler.CreateCategory)
	mux.HandleFunc("PUT /api/categories/{id}",
		categoryHandler.UpdateCategory)
	mux.HandleFunc("DELETE /api/categories/{id}",
		categoryHandler.DeleteCategory)

	// 8. WAREHOUSE API ROUTES
	mux.HandleFunc("GET /api/warehouses",
		warehouseHandler.GetAllWarehouses)
	mux.HandleFunc("GET /api/warehouses/{id}",
		warehouseHandler.GetWarehouseByID)
	mux.HandleFunc("POST /api/warehouses",
		warehouseHandler.CreateWarehouse)
	mux.HandleFunc("PUT /api/warehouses/{id}",
		warehouseHandler.UpdateWarehouse)
	mux.HandleFunc("DELETE /api/warehouses/{id}",
		warehouseHandler.DeleteWarehouse)

	// STOCK OUT API ROUTES
	mux.HandleFunc(
		"GET /api/stock-out",
		stockOutHandler.GetAllStockOut,
	)

	// STOCK IN API ROUTES
	mux.HandleFunc(
		"GET /api/stock-in",
		stockInHandler.GetAllStockIn,
	)

	// DASHBOARD API ROUTE
	mux.HandleFunc(
		"GET /api/dashboard",
		dashboardHandler.GetDashboard,
	)


	

	// 9. HEALTH CHECK
	mux.HandleFunc("GET /api/health",
		func(w http.ResponseWriter, r *http.Request) {

			ctx, cancel := context.WithTimeout(
				r.Context(), 3*time.Second,
			)
			defer cancel()

			if err := db.Ping(ctx); err != nil {
				http.Error(
					w,
					"Database unavailable",
					http.StatusServiceUnavailable,
				)
				return
			}

			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			json.NewEncoder(w).Encode(map[string]string{
				"status":   "success",
				"message":  "Inventory Go API is running!",
				"database": "PostgreSQL Connected",
			})
		})

	// 10. HTTP SERVER
	server := &http.Server{
		Addr:         ":8080",
		Handler:      corsMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// 11. START SERVER
	log.Println("================================")
	log.Println(" INVENTORY GO SYSTEM")
	log.Println(" PostgreSQL Connected!")
	log.Println(" Server running on port 8080")
	log.Println("================================")

	if err := server.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		log.Fatal("HTTP Server failed:", err)
	}
}

// CORS MIDDLEWARE
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set(
			"Access-Control-Allow-Origin",
			"http://localhost:5173",
		)

		w.Header().Set(
			"Access-Control-Allow-Methods",
			"GET, POST, PUT, DELETE, OPTIONS",
		)

		w.Header().Set(
			"Access-Control-Allow-Headers",
			"Content-Type, Authorization",
		)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
