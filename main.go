package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/seanjh/bosh-chat/message"
)

var templates = template.Must(template.ParseGlob("templates/**.html"))

func renderTemplate(w http.ResponseWriter, tmpl string) {
	err := templates.ExecuteTemplate(w, tmpl+".html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		renderTemplate(w, "index")
	} else {
		http.NotFound(w, r)
	}
}

func main() {
	fmt.Println("Starting server.")

	http.HandleFunc("/", index)
	http.Handle(
		"/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	http.HandleFunc("/messages/", message.HandleMessages)
	message.StartWriter()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
