package main

import (
	"github.com/CSE-682-Chat-App/go-chat-backend/internal/pkg/websocket"
	auth "github.com/CSE-682-Chat-App/go-chat-backend/pkg/authorization"
	message "github.com/CSE-682-Chat-App/go-chat-backend/pkg/message"
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
	hub.On("/user", auth.AuthMiddleware(auth.GetUser))
	hub.On("/join", auth.AuthMiddleware(message.JoinChannel))
	hub.On("/leave", auth.AuthMiddleware(message.LeaveChannel))

	hub.On("/message", auth.AuthMiddleware(message.OnMessage))

	log.Println("Serving at localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
