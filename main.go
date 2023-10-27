package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"slices"
	"time"

	"cloud.google.com/go/pubsub"
)

var (
	topic *pubsub.Topic
)

func init() {

	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

func main() {

	http.HandleFunc("/", hello)
	http.HandleFunc("/publish", publishHandler)
	http.HandleFunc("/push", pushHandler)
	http.HandleFunc("/chart", chartHandler)

	// setup pubsub
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

	// open up for business
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
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

type jsonResult struct {
	Message string `json:"message"`
}

func chartHandler(w http.ResponseWriter, r *http.Request) {

	origin := r.Header.Get("Origin")
	allowed := []string{"http://localhost:3000", "https://storage.googleapis.com"}
	if slices.Contains(allowed, origin) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Content-Type", "application/json")

	m := make(map[string]int)
	m["a"] = 1

	json.NewEncoder(w).Encode(m)
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

	origin := r.Header.Get("Origin")
	allowed := []string{"http://localhost:3000", "https://storage.googleapis.com"}
	if slices.Contains(allowed, origin) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	name := r.URL.Query().Get("name")

	n := 1
	requests := r.URL.Query().Get("requests")
	fmt.Sscan(requests, &n)

	ctx := context.Background()

	for i := 1; i <= n; i++ {

		time.Sleep(randomDuration(2, 5))

		msg := &pubsub.Message{
			Data: []byte(name),
		}

		if _, err := topic.Publish(ctx, msg).Get(ctx); err != nil {
			http.Error(w, fmt.Sprintf("Could not publish message: %v", err), http.StatusBadRequest)
			return
		}

		log.Printf("Message published: " + string(msg.Data))
	}

	msg := &jsonResult{
		Message: fmt.Sprintf("Message(s) published: %v", n),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

func randomDuration(min int, max int) time.Duration {
	return time.Duration(rand.Intn(max-min)+min) * time.Second
}
