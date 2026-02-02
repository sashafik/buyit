package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

// Service URLs
var (
	AuthServiceURL    = "http://localhost:8081"
	ProductServiceURL = "http://localhost:8082"
	OrderServiceURL   = "http://localhost:8083"
)

func main() {
	if url := os.Getenv("AUTH_SERVICE_URL"); url != "" {
		AuthServiceURL = url
	}
	if url := os.Getenv("PRODUCT_SERVICE_URL"); url != "" {
		ProductServiceURL = url
	}
	if url := os.Getenv("ORDER_SERVICE_URL"); url != "" {
		OrderServiceURL = url
	}
	mux := http.NewServeMux()

	// Auth Service Routes
	mux.Handle("/auth/", http.StripPrefix("/auth", newProxy(AuthServiceURL)))

	// Product Service Routes (Public)
	mux.Handle("/products", newProxy(ProductServiceURL))
	mux.Handle("/products/", newProxy(ProductServiceURL))

	// Order Service Routes (Protected)
	mux.Handle("/orders", authMiddleware(newProxy(OrderServiceURL)))

	handler := corsMiddleware(mux)

	fmt.Println("API Gateway running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func newProxy(target string) http.Handler {
	url, _ := url.Parse(target)
	return httputil.NewSingleHostReverseProxy(url)
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Validate token with Auth Service
		req, err := http.NewRequest("GET", AuthServiceURL+"/validate", nil)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		req.Header.Set("Authorization", token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Auth Service unreachable", http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract User ID from response
		var user struct {
			ID string `json:"id"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			http.Error(w, "Failed to parse auth response", http.StatusInternalServerError)
			return
		}

		// Add User ID to headers for downstream services
		r.Header.Set("X-User-ID", user.ID)

		next.ServeHTTP(w, r)
	})
}
