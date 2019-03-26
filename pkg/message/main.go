package message

import (
	"encoding/json"
	"fmt"
	ws "github.com/CSE-682-Chat-App/go-chat-backend/internal/pkg/websocket"
	auth "github.com/CSE-682-Chat-App/go-chat-backend/pkg/authorization"
	"log"
	"time"
)

type Timestamp struct {
	time.Time
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	log.Println(t)
	return []byte(fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))), nil
}
func (t Timestamp) UnmarshalJSON(b []byte) error {
	log.Println(t)
	return nil
}

type ChatMessage struct {
	Message   string     `json:"message"`
	Sender    *auth.User `json:"sender"`
	Timestamp Timestamp  `json:"timestamp"`
}

func (m *ChatMessage) ToString() string {
	str, err := json.Marshal(m)

	if err != nil {
		return ""
	}
	return string(str)
}

var channelMessages = map[string][]*ChatMessage{
	"general": []*ChatMessage{},
}
var channelUsers = map[string][]*auth.User{}

func messagesJson(channel string) string {
	str, err := json.Marshal(channelMessages[channel])
	if err != nil {
		log.Println(err)
		return "[]"
	}
	return string(str)
}

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
		Timestamp: Timestamp{time.Now()},
	}

	//Send the user the backlog of messages
	joinedResponse := &ws.Message{
		Path: fmt.Sprintf("/joined/%s", channel),
		Data: map[string]string{
			"messages": messagesJson(channel),
		},
	}
	c.Send <- joinedResponse.ToByte()

	//Broadcast that the user joined
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
		Timestamp: Timestamp{time.Now()},
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
		Timestamp: Timestamp{time.Now()},
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
