package main

import (
	"log"
	"net/http"
	"strconv"
)

func main() {
	hub := newHub()
	go hub.run()

	http.HandleFunc("/socket", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	hub.On("/message", func(m *Message, h *Hub, c *Client) {
		if !c.User.IsAuthorized() {
			return
		}

		log.Println("Message Received", m.Path, m.Data)

		//Broadcast the message to everyone but sender
		response := &Message{
			Path: "/message",
			Data: map[string]string{
				"message": m.Data["message"],
			},
			// sender: c, //Uncomment to exclude the sender
		}

		h.broadcast <- response
	})

	hub.On("/signin", func(m *Message, h *Hub, c *Client) {
		c.User.CheckCredentials(m.Data["username"], m.Data["password"])

		response := &Message{
			Path: "/authStatusChanged",
			Data: map[string]string{
				"isAuthorized": strconv.FormatBool(c.User.IsAuthorized()),
			},
		}

		//Send the response just to the recipient
		c.send <- response.ToByte()
	})

	log.Println("Serving at localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
