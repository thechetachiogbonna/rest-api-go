package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rest-api-go/internal/database"
	"rest-api-go/internal/utils"
)

func main() {
	pool, err := database.ConnectToDB(utils.Env.DATABASE_URL)

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

		fmt.Println(r.Header.Get("user-agent"))
		fmt.Println(r.RemoteAddr)

		w.WriteHeader(http.StatusOK)
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

	port := utils.Env.PORT

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Println("Server is running on http://localhost:" + port)
	log.Fatal(server.ListenAndServe())
}
