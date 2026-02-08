package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// set up database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("failed to initialize database:", err)
	}
	defer db.Close()

	// =====================
	// PRODUCT SETUP
	// =====================

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Product routes
	http.HandleFunc("/api/products", productHandler.HandleProducts)     // GET & POST
	http.HandleFunc("/api/products/", productHandler.HandleProductByID) // GET, PUT, DELETE

	// =====================
	// CATEGORY SETUP
	// =====================

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Category routes
	http.HandleFunc("/api/categories", categoryHandler.HandleCategories)    // GET & POST
	http.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID) // GET, PUT, DELETE

	// =====================
	// TRANSACTION SETUP
	// =====================
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	http.HandleFunc("/api/checkout", transactionHandler.HandleCheckout) // POST
	http.HandleFunc("/api/report/sales-summary", transactionHandler.HandleReport)

	// Health Check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	fmt.Println("Server running at http://localhost:" + config.Port)

	addr := "0.0.0.0:" + config.Port
	fmt.Println("Starting server at", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
