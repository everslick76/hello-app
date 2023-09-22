package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/pubsub"
)

var (
	topic *pubsub.Topic
)

func main() {

	setupLogging()

	setupPubSub()

	setupRest()

	// open up for business
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
	
}

func setupLogging() {

	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

func setupPubSub() {

	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, "cloud-core-376009")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	topic = client.Topic("hello")

	// check if the topic exists
	exists, err := topic.Exists(ctx)
	if err != nil || !exists {
		log.Fatal(err)
	}
}

func setupRest() {

    http.HandleFunc("/", hello)
	http.HandleFunc("/push", pushHandler)
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