package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type OrderItem struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type Order struct {
	ID        string      `json:"id"`
	UserID    string      `json:"userId"`
	Items     []OrderItem `json:"items"`
	Status    string      `json:"status"`
	CreatedAt time.Time   `json:"createdAt"`
}

type InventoryRequest struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

var (
	orders      = make(map[string]Order)
	ordersMutex sync.RWMutex
)

// Simplified: We assume Product Service is running at localhost:8082
const ProductServiceURL = "http://localhost:8082"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", createOrderHandler)
	mux.HandleFunc("GET /orders", listOrdersHandler)

	fmt.Println("Order Service running on port 8083")
	log.Fatal(http.ListenAndServe(":8083", mux))
}

func createOrderHandler(w http.ResponseWriter, r *http.Request) {
	// In a real app, we would extract UserID from context (set by middleware/gateway)
	// Here, we'll expect it in the body or header for simplicity, or just Mock it.
	// Let's assume the Gateway passes "X-User-ID" header.
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		// Fallback for direct testing
		userID = "guest"
	}

	var items []OrderItem
	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if len(items) == 0 {
		http.Error(w, "Cart is empty", http.StatusBadRequest)
		return
	}

	// Process each item
	for _, item := range items {
		if err := decrementInventory(item.ProductID, item.Quantity); err != nil {
			http.Error(w, fmt.Sprintf("Failed to process item %s: %v", item.ProductID, err), http.StatusConflict)
			// In a real system, we would need to rollback previous successful decrements here!
			return
		}
	}

	orderID := fmt.Sprintf("order-%d", time.Now().UnixNano())
	newOrder := Order{
		ID:        orderID,
		UserID:    userID,
		Items:     items,
		Status:    "Confirmed",
		CreatedAt: time.Now(),
	}

	ordersMutex.Lock()
	orders[orderID] = newOrder
	ordersMutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newOrder)
}

func decrementInventory(productID string, quantity int) error {
	reqBody, _ := json.Marshal(InventoryRequest{
		ProductID: productID,
		Quantity:  quantity,
	})

	resp, err := http.Post(ProductServiceURL+"/inventory/decrement", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("inventory check failed with status: %d", resp.StatusCode)
	}

	return nil
}

func listOrdersHandler(w http.ResponseWriter, r *http.Request) {
	ordersMutex.RLock()
	defer ordersMutex.RUnlock()

	// Filter by User ID if provided?
	userID := r.Header.Get("X-User-ID")

	userOrders := make([]Order, 0)
	for _, o := range orders {
		if userID == "" || o.UserID == userID {
			userOrders = append(userOrders, o)
		}
	}

	json.NewEncoder(w).Encode(userOrders)
}
