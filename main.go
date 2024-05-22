package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	stripe "github.com/stripe/stripe-go"
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
	mux.HandleFunc("/api/stripe/webhook", webhookHandler)
	mux.HandleFunc("/signup", signupHandler)
	mux.HandleFunc("/checkout", checkoutHandler)
	mux.HandleFunc("/", indexHandler)

	http.ListenAndServe("127.0.0.1:8080", mux)

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	tmpl = template.Must(template.ParseFiles("static/html/index.html"))
	tmpl.Execute(w, nil)
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	tmpl = template.Must(template.ParseFiles("static/html/signup.html"))
	tmpl.Execute(w, nil)
}

func checkoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	tmpl = template.Must(template.ParseFiles("static/html/checkout.html"))
	tmpl.Execute(w, nil)
}
