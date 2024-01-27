package main

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:embed template/*.html
var fs embed.FS

var messages []string

func GetIndex(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFS(fs, "template/index.html")
	t.Execute(w, map[string]any{"Messages": messages})
}

func PostMessage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	message := r.FormValue("message")
	messages = append(messages, message)
	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	r := chi.NewRouter()
	r.Get("/", GetIndex)
	r.Post("/messages", PostMessage)
	http.ListenAndServe("127.0.0.1:8888", r)
}
