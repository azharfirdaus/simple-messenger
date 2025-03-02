package main

import (
	"log"
	"net/http"

	"github.com/azhar.firdaus/simple-messenger/config"
	"github.com/azhar.firdaus/simple-messenger/routes"
	"github.com/gorilla/mux"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var err error

	config.GlobalConfig, err = config.ReadConfig()
	if err != nil {
		return
	}

	router := mux.NewRouter()
	router.HandleFunc("/hello", routes.HelloHandler).Methods("GET")
	router.HandleFunc("/message", routes.CreateMessage).Methods("POST")

	// Start the server
	log.Printf("Server started on :%v", *config.GlobalConfig.Port)
	port := ":" + *config.GlobalConfig.Port
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Could not start server: %v", err)
		return
	}
}
