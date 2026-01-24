package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/internal/model"
	"net/http"
	"strconv"
	"strings"
)

var lastCategoryID = 2
var categories = []model.Category{
	{ID: 1, Name: "Food", Description: "Snack and Meals"},
	{ID: 2, Name: "Beverage", Description: "Various drinks"},
}

var lastProductID = 3
var products = []model.Product{
	{ID: 1, Name: "Indomie Godog", Price: 3500, Stock: 10, CategoryId: 1},
	{ID: 2, Name: "Vit 1000ml", Price: 3000, Stock: 40, CategoryId: 2},
	{ID: 3, Name: "Kecap", Price: 12000, Stock: 50, CategoryId: 1},
}

func isCategoryExist(id int) bool {
	for _, c := range categories {
		if c.ID == id {
			return true
		}
	}
	return false
}

// ==========================================
//         PRODUCT HANDLER
// ==========================================

func productHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Case 1: /api/products
	if path == "/api/products" {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(products)
			return
		}

		if r.Method == "POST" {
			var newProduct model.Product
			if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}

			if !isCategoryExist(newProduct.CategoryId) {
				http.Error(w, "Category ID not found", http.StatusBadRequest)
				return
			}

			lastProductID++
			newProduct.ID = lastProductID
			products = append(products, newProduct)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newProduct)
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Case 2: /api/products/{id}
	if strings.HasPrefix(path, "/api/products/") {
		idStr := strings.TrimPrefix(path, "/api/products/")

		if idStr == "" {
			http.Error(w, "Product ID required", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		// GET by ID
		if r.Method == "GET" {
			for _, p := range products {
				if p.ID == id {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(p)
					return
				}
			}
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		// PUT (update)
		if r.Method == "PUT" {
			var inputData model.Product
			if err := json.NewDecoder(r.Body).Decode(&inputData); err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}

			for i := range products {
				if products[i].ID == id {
					if inputData.CategoryId != 0 {
						if !isCategoryExist(inputData.CategoryId) {
							http.Error(w, "Category ID not found", http.StatusBadRequest)
							return
						}
						products[i].CategoryId = inputData.CategoryId
					}

					if inputData.Name != "" {
						products[i].Name = inputData.Name
					}
					if inputData.Price != 0 {
						products[i].Price = inputData.Price
					}
					if inputData.Stock != 0 {
						products[i].Stock = inputData.Stock
					}

					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(products[i])
					return
				}
			}
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		// DELETE
		if r.Method == "DELETE" {
			for i, p := range products {
				if p.ID == id {
					products = append(products[:i], products[i+1:]...)
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted successfully"})
					return
				}
			}
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.NotFound(w, r)
}

// ==========================================
//        CATEGORY HANDLER
// ==========================================

func categoryHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Case 1: /api/categories
	if path == "/api/categories" {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(categories)
			return
		}

		if r.Method == "POST" {
			var newCategory model.Category
			if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}
			lastCategoryID++
			newCategory.ID = lastCategoryID
			categories = append(categories, newCategory)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newCategory)
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Case 2: /api/categories/{id}
	if strings.HasPrefix(path, "/api/categories/") {
		idStr := strings.TrimPrefix(path, "/api/categories/")

		if idStr == "" {
			http.Error(w, "Category ID required", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid Category ID", http.StatusBadRequest)
			return
		}

		// GET by ID
		if r.Method == "GET" {
			for _, c := range categories {
				if c.ID == id {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(c)
					return
				}
			}
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}

		// PUT (update)
		if r.Method == "PUT" {
			var inputData model.Category
			if err := json.NewDecoder(r.Body).Decode(&inputData); err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}

			for i, c := range categories {
				if c.ID == id {
					if inputData.Name != "" {
						categories[i].Name = inputData.Name
					}
					if inputData.Description != "" {
						categories[i].Description = inputData.Description
					}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(categories[i])
					return
				}
			}
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}

		// DELETE
		if r.Method == "DELETE" {
			for i, c := range categories {
				if c.ID == id {
					categories = append(categories[:i], categories[i+1:]...)
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(map[string]string{"message": "Category deleted successfully"})
					return
				}
			}
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.NotFound(w, r)
}

// ==========================================
//             MAIN FUNCTION
// ==========================================

func main() {
	// Product routes
	http.HandleFunc("/api/products", productHandler)
	http.HandleFunc("/api/products/", productHandler)

	// Category routes
	http.HandleFunc("/api/categories", categoryHandler)
	http.HandleFunc("/api/categories/", categoryHandler)

	// Health Check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	fmt.Println("Server running at http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
