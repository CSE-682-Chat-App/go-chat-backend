package main

import (
// "encoding/json"
// "log"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan *Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	messageCallbacks map[string][]onCallback
}

func newHub() *Hub {
	return &Hub{
		broadcast:        make(chan *Message),
		register:         make(chan *Client),
		unregister:       make(chan *Client),
		clients:          make(map[*Client]bool),
		messageCallbacks: make(map[string][]onCallback),
	}
}

type onCallback func(*Message, *Hub, *Client)

//On registers an event path callback
func (h *Hub) On(path string, cb onCallback) {
	h.messageCallbacks[path] = append(h.messageCallbacks[path], cb)
}

//Handle calls all the handlers for a path
func (h *Hub) Handle(m *Message, c *Client) {
	if cbs, ok := h.messageCallbacks[m.Path]; ok {
		for _, cb := range cbs {
			cb(m, h, c)
		}
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				if message.WasSentBy(client) {
					break
				}
				if !message.IsRecipient(client) {
					break
				}

				select {
				case client.send <- message.ToByte():
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
