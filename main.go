package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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
	http.HandleFunc("/publish", publishHandler)

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

func publishHandler(w http.ResponseWriter, r *http.Request) {

	n := 1
	requests := r.URL.Query().Get("requests")
	fmt.Sscan(requests, &n)

	ctx := context.Background()

	for i := 1; i <= n; i++ {

		currentTime := time.Now().String()

		msg := &pubsub.Message{
			Data: []byte(currentTime),
		}
	
		if _, err := topic.Publish(ctx, msg).Get(ctx); err != nil {
			http.Error(w, fmt.Sprintf("Could not publish message: %v", err), http.StatusBadRequest)
			return
		}
	
		log.Printf("Message published: " + string(msg.Data))
	}

	log.Printf("%s Message(s) published", requests)

	fmt.Fprint(w, "Message(s) published: " + requests)
}