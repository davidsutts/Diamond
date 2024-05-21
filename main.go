package main

import (
	"html/template"
	"net/http"
)

var tmpl *template.Template

func main() {

	// Create a webserver.
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)

	http.ListenAndServe("127.0.0.1:8080", mux)

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl = template.Must(template.ParseFiles("static/html/index.html"))
	tmpl.Execute(w, nil)
}
