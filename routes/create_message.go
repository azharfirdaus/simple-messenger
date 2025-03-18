package routes

import (
	"encoding/json"
	"net/http"

	"github.com/azhar.firdaus/simple-messenger/config"
	msg "github.com/azhar.firdaus/simple-messenger/messaging"
)

type CreatedMessageRequest struct {
	ID      *uint64 `json:"ID"`
	Message *string `json:"message"`
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	var err error

	xUserID := r.Header.Get("X-UserId")
	if xUserID == "" {
		http.Error(w, "X-userId is empty", http.StatusBadRequest)
		return
	}

	var request CreatedMessageRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Body message is not recognized", http.StatusBadRequest)
		return
	}

	if request.ID == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "ID is required"})
	}

	if request.Message == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "message is required"})
	} else if len(*request.Message) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "message is empty"})
	}

	kafkaClient := msg.NewKafkaMessageQueueClientImpl(*config.GlobalConfig.KafkaBroker, "create_message")
	err = kafkaClient.Produce(request.ID, request.Message)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}
