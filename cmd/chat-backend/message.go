package main

import (
	"encoding/json"
)

//Message is a message from the channels
type Message struct {
	Path string            `json:"path"`
	Data map[string]string `json:"data"`

	recepients []*Client
	sender     *Client
}

//NewMessage returns a new Message Struct
func NewMessage() *Message {
	return &Message{
		recepients: make([]*Client, 0),
		Data:       make(map[string]string),
	}
}

//SetSender sets the sender of the message
func (m *Message) SetSender(c *Client) {
	m.sender = c
}

//WasSentBy returns whether a message was sent by a client
func (m *Message) WasSentBy(c *Client) bool {
	return m.sender == c
}

//IsRecipient returns whether the client is a tagged recipient
func (m *Message) IsRecipient(c *Client) bool {
	//If there are no listed recipients then all clients are recipients
	if len(m.recepients) == 0 {
		return true
	}

	for _, r := range m.recepients {
		if r == c {
			return true
		}
	}
	return false
}

//ToString returns the Message json encoded as a string
func (m *Message) ToString() string {
	b := m.ToByte()
	return string(b)
}

//ToByte returns the message json encoded as a byte array
func (m *Message) ToByte() []byte {
	str, err := json.Marshal(m)

	if err != nil {
		return []byte{}
	}
	return str
}

//AddRecipient adds a potential recipient for a message
func (m *Message) AddRecipient(c *Client) {
	m.recepients = append(m.recepients, c)
}
