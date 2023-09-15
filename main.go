package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {

	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	// log.SetFlags(log.LstdFlags | log.Lmicroseconds)

    http.HandleFunc("/", hello)
	http.HandleFunc("/push", pushHandler)
    http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Hello from hello-app!")
}

type pushRequest struct {
	Message struct {
		Attributes map[string]string
		Data       []byte
		ID         string `json:"message_id"`
	}
	Subscription string
}

func pushHandler(w http.ResponseWriter, r *http.Request) {

	msg := &pushRequest{}
	if err := json.NewDecoder(r.Body).Decode(msg); err != nil {
		http.Error(w, fmt.Sprintf("Could not decode body: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("Message received: %s", string(msg.Message.Data))
}