package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/pubsub"
)

// var (
// 	topic *pubsub.Topic

// 	// Messages received by this instance.
// 	messagesMu sync.Mutex
// 	messages   []string
// )

// const maxMessages = 10

func main() {

    http.HandleFunc("/", hello)
	http.HandleFunc("/push", pushHandler)
    http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Hello from hello-app!")
}

// type pushRequest struct {
// 	now string
// }

func pushHandler(w http.ResponseWriter, r *http.Request) {

	msg := &pubsub.Message{}
	if err := json.NewDecoder(r.Body).Decode(msg); err != nil {
		http.Error(w, fmt.Sprintf("Could not decode body: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("Message received: %s", msg.Data)

	// messagesMu.Lock()
	// defer messagesMu.Unlock()
	// // Limit to ten.
	// messages = append(messages, string(msg.Message.Data))
	// if len(messages) > maxMessages {
	// 	messages = messages[len(messages)-maxMessages:]
	// }
}