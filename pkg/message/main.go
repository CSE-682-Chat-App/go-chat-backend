package message

import (
	"encoding/json"
	"fmt"
	ws "github.com/CSE-682-Chat-App/go-chat-backend/internal/pkg/websocket"
	auth "github.com/CSE-682-Chat-App/go-chat-backend/pkg/authorization"
	"log"
	"time"
)

type ChatMessage struct {
	Message   string     `json:"message"`
	Sender    *auth.User `json:"sender"`
	Timestamp time.Time  `json:"timestamp"`
}

func (m *ChatMessage) ToString() string {
	str, err := json.Marshal(m)

	if err != nil {
		return ""
	}
	return string(str)
}

var channelMessages = map[string][]*ChatMessage{}
var channelUsers = map[string][]*auth.User{}

func JoinChannel(m *ws.Message, h *ws.Hub, c *ws.Client) {
	user := c.User.(*auth.User)
	channel := m.Get("channel")

	channelUsers[channel] = append(channelUsers[channel], user)

	message := &ChatMessage{
		Message: fmt.Sprintf("%s Joined %s", user.Name, channel),
		Sender: &auth.User{
			Name: "System",
			UUID: "0",
		},
		Timestamp: time.Now(),
	}

	response := &ws.Message{
		Path: fmt.Sprintf("/message/%s", channel),
		Data: map[string]string{
			"message": message.ToString(),
		},
	}
	response.SetSender(c)

	h.Broadcast <- response
}
func LeaveChannel(m *ws.Message, h *ws.Hub, c *ws.Client) {
	user := c.User.(*auth.User)
	channel := m.Get("channel")

	for i, u := range channelUsers[channel] {
		if u == user {
			channelUsers[channel] = append(channelUsers[channel][:i], channelUsers[channel][i+1:]...)
			break
		}
	}

	message := &ChatMessage{
		Message: fmt.Sprintf("%s Left %s", user.Name, channel),
		Sender: &auth.User{
			Name: "System",
			UUID: "0",
		},
		Timestamp: time.Now(),
	}

	response := &ws.Message{
		Path: fmt.Sprintf("/message/%s", channel),
		Data: map[string]string{
			"message": message.ToString(),
		},
	}
	response.SetSender(c)

	h.Broadcast <- response
}

func OnMessage(m *ws.Message, h *ws.Hub, c *ws.Client) {
	channel := m.Get("channel")
	if channel == "" {
		channel = "general"
	}

	user := c.User.(*auth.User)

	message := &ChatMessage{
		Message:   m.Get("message"),
		Sender:    user,
		Timestamp: time.Now(),
	}

	channelMessages[channel] = append(channelMessages[channel], message)

	response := &ws.Message{
		Path: fmt.Sprintf("/message/%s", channel),
		Data: map[string]string{
			"message": message.ToString(),
		},
	}

	response.SetSender(c)

	log.Println(*response)

	h.Broadcast <- response
}
