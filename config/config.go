package config

import (
	"encoding/json"
	"log"
	"os"
)

var GlobalConfig *Config

type Config struct {
	Port        *string `json:"port"`
	KafkaBroker *string `json:"kafkaBroker"`
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
