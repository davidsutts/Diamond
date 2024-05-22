package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
)

type Item struct {
	ID string `json:"id"`
}

type Order struct {
	Items []Item `json:"items"`
}

func paymentIntentHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("failed to read body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var order Order

	// Unmarshal the JSON data into the struct
	err = json.Unmarshal([]byte(body), &order)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	log.Printf("%+v", order)

	cost := calcAmount(order.Items[0].ID)
	if cost == 0 {
		log.Println("invalid product")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	params := stripe.PaymentIntentParams{
		Amount:   stripe.Int64(cost),
		Currency: stripe.String(string(stripe.CurrencyAUD)),
	}

	pi, err := paymentintent.New(&params)
	if err != nil {
		log.Println("failed to create new payment intent:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("New payment intent created:", pi.ClientSecret)

	// Write client secret to response writer.
	writeJSON(w, struct {
		ClientSecret string `json:"clientSecret"`
		OrderTotal   int64  `json:"orderTotal"`
	}{
		ClientSecret: pi.ClientSecret,
		OrderTotal:   cost,
	})
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	event := stripe.Event{}
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	err = json.Unmarshal(payload, &event)
	if err != nil {
		log.Println("error unmarshalling payload:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("%+v", event)
	w.WriteHeader(http.StatusAccepted)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("json.NewEncoder.Encode: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := io.Copy(w, &buf); err != nil {
		log.Printf("io.Copy: %v", err)
		return
	}
}

const premiumSubscription = "prem"
const superSubscription = "super"

func calcAmount(product string) int64 {
	switch product {
	case premiumSubscription:
		return 599 // $5.99
	case superSubscription:
		return 1099 // $10.99
	default:
		return 0
	}
}
