package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

type InventoryRequest struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

var (
	products      = make(map[string]Product)
	productsMutex sync.RWMutex
)

func main() {
	// Seed some data
	seedData()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /products", listProductsHandler)
	mux.HandleFunc("GET /products/", getProductHandler) // Handles /products/{id}
	mux.HandleFunc("POST /products", createProductHandler)
	mux.HandleFunc("POST /inventory/decrement", decrementInventoryHandler)

	fmt.Println("Product Service running on port 8082")
	log.Fatal(http.ListenAndServe(":8082", mux))
}

func seedData() {
	p1 := Product{ID: "1", Name: "Laptop", Description: "High performance laptop", Price: 1200.00, Stock: 10}
	p2 := Product{ID: "2", Name: "Phone", Description: "Smartphone with good camera", Price: 800.00, Stock: 20}
	products["1"] = p1
	products["2"] = p2
}

func listProductsHandler(w http.ResponseWriter, r *http.Request) {
	productsMutex.RLock()
	defer productsMutex.RUnlock()

	productList := make([]Product, 0, len(products))
	for _, p := range products {
		productList = append(productList, p)
	}

	json.NewEncoder(w).Encode(productList)
}

func getProductHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/products/"):]
	if id == "" {
		http.Error(w, "Missing product ID", http.StatusBadRequest)
		return
	}

	productsMutex.RLock()
	product, exists := products[id]
	productsMutex.RUnlock()

	if !exists {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
}

func createProductHandler(w http.ResponseWriter, r *http.Request) {
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	p.ID = fmt.Sprintf("prod-%d", time.Now().UnixNano())
	
	productsMutex.Lock()
	products[p.ID] = p
	productsMutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func decrementInventoryHandler(w http.ResponseWriter, r *http.Request) {
	var req InventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	productsMutex.Lock()
	defer productsMutex.Unlock()

	product, exists := products[req.ProductID]
	if !exists {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	if product.Stock < req.Quantity {
		http.Error(w, "Insufficient stock", http.StatusConflict)
		return
	}

	product.Stock -= req.Quantity
	products[req.ProductID] = product

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}
