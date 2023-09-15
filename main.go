package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"cloud.google.com/go/pubsub"
)

var (

	// Messages received by this instance.
	messagesMu sync.Mutex
	messages   []string
)

const maxMessages = 10

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

func pushHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("Call to push endpoint received: %s", r.Body)

	msg := &pubsub.Message{}
	if err := json.NewDecoder(r.Body).Decode(msg); err != nil {
		http.Error(w, fmt.Sprintf("Could not decode body: %v", err), http.StatusBadRequest)
		return
	}

	messagesMu.Lock()
	defer messagesMu.Unlock()
	// Limit to ten.
	message := string(msg.Data)
	messages = append(messages, message)
	if len(messages) > maxMessages {
		messages = messages[len(messages)-maxMessages:]
	}

	log.Printf("Message received: %s", message)
}