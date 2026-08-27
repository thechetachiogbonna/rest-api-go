package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rest-api-go/internal/config"
	"rest-api-go/internal/database"
	"rest-api-go/internal/routes"
)

type ApiResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func main() {
	pool, err := database.ConnectToDB(config.Env.DATABASE_URL)

	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	defer pool.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(ApiResponse{
			Message: "Hello from my first API in GO!",
			Status:  "OK",
		}); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
	})

	mux.Handle("/auth/", http.StripPrefix("/auth", routes.AuthRoutes(pool)))

	port := config.Env.PORT

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Println("Server is running on http://localhost:" + port)
	log.Fatal(server.ListenAndServe())
}
