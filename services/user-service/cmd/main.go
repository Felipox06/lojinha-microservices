package main

import (
	"encoding/json"
	"fmt" // Debug
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main(){
	router := mux.NewRouter()

	v1 := router.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/health", healthHandler).Methods("GET")
	v1.HandleFunc("/users", createUser).Methods("POST")
	v1.HandleFunc("/users/{id}", getUser).Methods("GET")

	router.Use(loggingMiddleware)

	srv := &http.Server{
		Handler: router,

		Addr: ":8001",
		ReadTimeout: 15 * time.Second,
		WriteTimeout: 15 * time.Second,

	}

	fmt.Println("User Service v1 starting on :8001")
	log.Fatal(srv.ListenAndServe())
}

func loggingMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"service": "user-service",
		"version": "1.0.0",
	})
}

func createUser(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User creation endpoint",
	})
}

func getUser(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	userID := vars["id"]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Getting user %s", userID),
	})
}

