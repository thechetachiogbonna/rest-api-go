package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(struct {
			Message string `json:"message"`
			Status  string `json:"status"`
		}{
			Message: "Hello from my first API in GO!",
			Status:  "success",
		}); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
	})

	port := os.Getenv("PORT")

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Println("Server is running on http://localhost:" + port)
	log.Fatal(server.ListenAndServe())
}
