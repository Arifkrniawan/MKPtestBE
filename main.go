package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

var products = []Product{
	{ID: 1, Name: "Product A", Price: 100, Stock: 10},
	{ID: 2, Name: "Product B", Price: 200, Stock: 20},
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	storedPassword, ok := users[user.Username]
	if !ok || storedPassword != user.Password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	authToken := "dummy_token"

	response := map[string]string{"token": authToken}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}

func productHandler(w http.ResponseWriter, r *http.Request) {

	authToken := r.Header.Get("Authorization")
	if authToken != "dummy_token" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	jsonResponse, err := json.Marshal(products)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func main() {

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/products", productHandler)

	fmt.Println("Server started at localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
