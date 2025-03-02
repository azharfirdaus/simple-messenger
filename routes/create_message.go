package routes

import (
	"encoding/json"
	"net/http"

	msg "github.com/azhar.firdaus/simple-messenger/messaging"
)

type CreatedMessageRequest struct {
	Message *string `json:"message"`
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	var err error
	var request CreatedMessageRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Body message is not recognized", http.StatusBadRequest)
		return
	}

	if request.Message == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "message is required"})
	} else if len(*request.Message) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "message is empty"})
	}

	kafkaClient := msg.NewKafkaMessageQueueClientImpl("localhost:9092", "create_message")
	err = kafkaClient.Produce(*request.Message)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}
