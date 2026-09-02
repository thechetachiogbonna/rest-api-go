package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"rest-api-go/internal/config"
	"rest-api-go/internal/database"
	"rest-api-go/internal/routes"
	"rest-api-go/internal/utils"

	"github.com/rs/cors"
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
			utils.WriteError(w, http.StatusMethodNotAllowed, errors.New("Method Not Allowed"))
			return
		}

		if r.URL.Path != "/" {
			utils.WriteError(w, http.StatusNotFound, errors.New("Not Found"))
			return
		}

		if err := utils.WriteJson(w, http.StatusOK, ApiResponse{
			Message: "Hello from my first API in GO!",
			Status:  "OK",
		}); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
	})

	mux.Handle("/api/auth/", http.StripPrefix("/api/auth", routes.AuthRoutes(pool)))

	port := config.Env.PORT

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{config.Env.FRONTEND_URL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            true,
	})

	handler := c.Handler(mux)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	fmt.Println("Server is running on http://localhost:" + port)
	log.Fatal(server.ListenAndServe())
}
