package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"` // In real app, store hash
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

var (
	users       = make(map[string]User)
	tokens      = make(map[string]string) // token -> userID
	usersMutex  sync.RWMutex
	tokensMutex sync.RWMutex
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", registerHandler)
	mux.HandleFunc("POST /login", loginHandler)
	mux.HandleFunc("GET /validate", validateHandler)

	fmt.Println("Auth Service running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	usersMutex.Lock()
	defer usersMutex.Unlock()

	// Check if user exists
	for _, u := range users {
		if u.Username == creds.Username {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}
	}

	userID := fmt.Sprintf("user-%d", time.Now().UnixNano())
	newUser := User{
		ID:       userID,
		Username: creds.Username,
		Password: creds.Password, // Storing plain text for simplicity of demo
	}
	users[userID] = newUser

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	usersMutex.RLock()
	var foundUser *User
	for _, u := range users {
		if u.Username == creds.Username && u.Password == creds.Password {
			foundUser = &u
			break
		}
	}
	usersMutex.RUnlock()

	if foundUser == nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate simple token
	token := fmt.Sprintf("token-%s-%d", foundUser.ID, time.Now().UnixNano())

	tokensMutex.Lock()
	tokens[token] = foundUser.ID
	tokensMutex.Unlock()

	json.NewEncoder(w).Encode(AuthResponse{
		Token: token,
		User:  *foundUser,
	})
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	tokensMutex.RLock()
	userID, exists := tokens[token]
	tokensMutex.RUnlock()

	if !exists {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	usersMutex.RLock()
	user, userExists := users[userID]
	usersMutex.RUnlock()

	if !userExists {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(user)
}
