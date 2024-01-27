package main

import (
	"embed"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

//go:embed template/*.html
var fs embed.FS

type Message struct {
	ID   string
	Body string
}

var messages []Message
var messageID int

func GetIndex(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.Header.Get("Accept"), "text/vnd.turbo-stream.html") {
		w.Header().Set("Content-Type", "text/vnd.turbo-stream.html")
		lastMessageID := r.FormValue("lastMessageId")
		messages := messages
		if lastMessageID != "" {
			for i, message := range messages {
				if message.ID == lastMessageID {
					messages = messages[i+1:]
					break
				}
			}
		}
		t, _ := template.ParseFS(fs, "template/messages.html")
		t.Execute(w, map[string]any{"Messages": messages})
		return
	}
	t, _ := template.ParseFS(fs, "template/index.html")
	t.Execute(w, map[string]any{"Messages": messages})
}

func PostMessage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	message := r.FormValue("message")
	lastMessageID := r.FormValue("lastMessageId")
	messageID++
	messages = append(messages, Message{
		ID:   "message-" + strconv.Itoa(messageID),
		Body: message,
	})
	redirectTo := "/"
	if lastMessageID != "" {
		redirectTo = "/?lastMessageId=" + lastMessageID
	}
	http.Redirect(w, r, redirectTo, http.StatusFound)
}

func main() {
	r := chi.NewRouter()
	r.Get("/", GetIndex)
	r.Post("/messages", PostMessage)
	http.ListenAndServe("127.0.0.1:8888", r)
}
