package main

import (
	"log"
	"net/http"

	"github.com/azhar.firdaus/simple-messenger/routes"
)

func main() {
	// Register handlers
	http.HandleFunc("/hello", routes.HelloHandler)

	// Start the server
	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
