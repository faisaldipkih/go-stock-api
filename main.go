package main

import (
    "log"
    "net/http"

    "go-stock-api/db"
    "go-stock-api/handlers"
    "go-stock-api/middleware"

    "github.com/gorilla/mux"
)

func main() {
    db.Init()

    r := mux.NewRouter()

    // Pasang logging middleware
    r.Use(middleware.LoggingMiddleware)

    // Define routes
    r.HandleFunc("/stocks", handlers.GetStocks).Methods("GET")
    r.HandleFunc("/stocks/{id}", handlers.GetStock).Methods("GET")
    r.HandleFunc("/stocks", handlers.CreateStock).Methods("POST")
    r.HandleFunc("/stocks/{id}", handlers.UpdateStock).Methods("PUT")
    r.HandleFunc("/stocks/{id}", handlers.DeleteStock).Methods("DELETE")

    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
