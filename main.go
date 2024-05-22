package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
)

var tmpl *template.Template

func main() {

	stripe.Key = os.Getenv("STRIPE_KEY")
	if stripe.Key == "" {
		log.Fatal("There is no stripe key")
	}

	// Create a webserver.
	mux := http.NewServeMux()

	// Redirect file requests to static dir.
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/api/stripe/create-payment-intent", paymentIntentHandler)
	mux.HandleFunc("/", indexHandler)

	http.ListenAndServe("127.0.0.1:8080", mux)

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	tmpl = template.Must(template.ParseFiles("static/html/index.html"))
	tmpl.Execute(w, nil)
}

func paymentIntentHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Parse the incoming request.
	err := r.ParseForm()
	if err != nil {
		log.Println("failed to parse form:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	products := strings.Split(r.FormValue("products"), ",")

	params := stripe.PaymentIntentParams{
		Amount:   stripe.Int64(calcAmount(products)),
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
	}{
		ClientSecret: pi.ClientSecret,
	})
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

func calcAmount(product []string) int64 {
	return 1500
}
