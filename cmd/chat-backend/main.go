package main

import (
	"github.com/CSE-682-Chat-App/go-chat-backend/internal/pkg/websocket"
	auth "github.com/CSE-682-Chat-App/go-chat-backend/pkg/authorization"
	"log"
	"net/http"
)

func main() {
	hub := websocket.New()
	go hub.Run()

	http.HandleFunc("/socket", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r)
	})

	hub.On("/auth", auth.Authenticate)
	hub.On("/signup", auth.Signup)

	log.Println("Serving at localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
