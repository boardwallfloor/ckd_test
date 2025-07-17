package main

import (
	"boardwallfloor/ckd/internal/controller"
	"boardwallfloor/ckd/internal/db"
	"boardwallfloor/ckd/internal/middleware"
	"boardwallfloor/ckd/internal/service"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	dbConn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer dbConn.Close()

	queries := db.New(dbConn)

	authService := service.NewAuthService(queries)
	txService := service.NewTransactionService(queries)

	authController := controller.NewAuthController(authService)
	txController := controller.NewTransactionController(txService)

	r := mux.NewRouter()

	r.HandleFunc("/register", authController.Register).Methods("POST")
	r.HandleFunc("/login", authController.Login).Methods("POST")

	txRouter := r.PathPrefix("/transaction").Subrouter()
	txRouter.Use(middleware.AuthMiddleware)
	txRouter.HandleFunc("", txController.GetTransactions).Methods("GET")
	txRouter.HandleFunc("/process", txController.ProcessTransaction).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), r); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
