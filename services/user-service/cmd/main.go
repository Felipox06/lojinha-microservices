package main

import (
	"encoding/json"
	"fmt" // Debug
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/Felipox06/lojinha-microservices/services/user-service/internal/handlers"
    "github.com/Felipox06/lojinha-microservices/services/user-service/internal/services"
)

func main(){
    // Primeira etapa: inicializar camadas (Dependency Injection)
	userService := services.NewUserService()
	userHandler := handlers.NewUserHandler(userService)

	// Segunda etapa: Configurar rotas
	router := mux.NewRouter()

	// Middleware
	router.Use(loggingMiddleware)
	router.Use(corsMiddleware)

	v1 := router.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/health", healthHandler).Methods("GET")


	v1.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	v1.HandleFunc("/users", userHandler.ListUsers).Methods("GET")
	v1.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	v1.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	v1.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	// Terceira etapa: Configurar servidor
	srv := &http.Server{
		Handler: router,
		Addr: ":8001",
		ReadTimeout: 15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout: 60 * time.Second,
	}

	fmt.Println("User Service v1 starting on :8001")
	log.Fatal(srv.ListenAndServe())
}

func loggingMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Printf(
			"[%s] %s %s %v",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			time.Since(start),
		)
	})
}

// CORS - permite requests do browser
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control_Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

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