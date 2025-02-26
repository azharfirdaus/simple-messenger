package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/azhar.firdaus/simple-messenger/routes"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config, err := ReadConfig()
	if err != nil {
		return
	}

	// Register handlers
	http.HandleFunc("/hello", routes.HelloHandler)

	// Start the server
	log.Printf("Server started on :%v", *config.Port)
	port := ":" + *config.Port
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Could not start server: %v", err)
		return
	}
}

func ReadConfig() (*Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
		return nil, err
	}

	return &config, nil
}
