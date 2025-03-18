package main

import (
	"log"
	"net/http"
	"time"

	"github.com/azhar.firdaus/simple-messenger/client"
	"github.com/azhar.firdaus/simple-messenger/config"
	"github.com/azhar.firdaus/simple-messenger/dao"
	"github.com/azhar.firdaus/simple-messenger/messaging"
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

	client.GlobalClient, err = client.NewClient(config.GlobalConfig)
	if err != nil {
		log.Fatalf("failed to build client %v", err)
		return
	}

	go consumeCreateMessageTopic(config.GlobalConfig.KafkaBroker, client.GlobalClient)

	router := mux.NewRouter()
	router.HandleFunc("/message", routes.CreateMessage).Methods("POST")

	// Start the server
	log.Printf("Server started on :%v", *config.GlobalConfig.Port)
	port := ":" + *config.GlobalConfig.Port
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Could not start server: %v", err)
		return
	}
}

func consumeCreateMessageTopic(broker *string, globalClient *client.Client) {
	kafkaClient := messaging.NewKafkaMessageQueueClientImpl(*broker, "create_message")
	chatDAO := globalClient.ChatDAO
	consumerGroupHandler := messaging.ConsumerGroupHandler{
		Handler: func(topic *string, partition *int32, offset *int64, key, value *[]byte) {
			log.Printf("Received message: %s (topic: %s, partition: %d, offset: %d)\n", string(*value), *topic, *partition, *offset)
			data := string(*value)
			now := time.Now()
			chat := dao.Chat{
				Data: []*dao.Message{
					{
						Data:      &data,
						CreatedAt: &now,
					},
				},
				CreatedAt: &now,
			}
			chatDAO.InsertOne(&chat)
		},
	}
	kafkaClient.Consume(&consumerGroupHandler)
}
