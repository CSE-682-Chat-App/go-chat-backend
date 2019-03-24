package authorization

import (
	ws "github.com/CSE-682-Chat-App/go-chat-backend/internal/pkg/websocket"
	uuid "github.com/google/uuid"
	// "log"
)

//User is an authorized user
type User struct {
	UUID       string
	Name       string
	password   string
	authorized bool
}

var users = map[string]*User{}

//CheckCredentials checks whether the given credentials is authorized
func (u *User) CheckCredentials(password string) bool {
	u.authorized = u.password == password
	return u.authorized
}

func Signup(m *ws.Message, h *ws.Hub, c *ws.Client) {
	username := m.Get("username")
	password := m.Get("password")

	response := &ws.Message{
		Path: "/signup_response",
		Data: map[string]string{
			"message": "",
			"success": "false",
		},
	}

	if _, ok := users[username]; !ok {
		users[username] = &User{
			Name:     username,
			password: password,
			UUID:     uuid.New().String(),
		}

		response.Data["message"] = "User created successfully!"
		response.Data["success"] = "true"
	} else {
		response.Data["message"] = "User already exists!"
	}

	c.Send <- response.ToByte()
}

//Authenticate handles authenticating the user
func Authenticate(m *ws.Message, h *ws.Hub, c *ws.Client) {
	username := m.Get("username")
	password := m.Get("password")

	response := &ws.Message{
		Path: "/auth_response",
		Data: map[string]string{
			"message": "",
			"success": "false",
		},
	}

	if _, ok := users[username]; ok && users[username].CheckCredentials(password) {
		c.User = users[username]
		response.Data["message"] = "Success"
		response.Data["success"] = "true"
	} else {
		response.Data["message"] = "Incorrect username or password"
		response.Data["success"] = "false"
	}

	c.Send <- response.ToByte()
}

// hub.On("/message", func(m *Message, h *Hub, c *Client) {
// 	if !c.User.IsAuthorized() {
// 		return
// 	}
//
// 	log.Println("Message Received", m.Path, m.Data)
//
// 	//Broadcast the message to everyone but sender
// 	response := &Message{
// 		Path: "/message",
// 		Data: map[string]string{
// 			"message": m.Data["message"],
// 		},
// 		// sender: c, //Uncomment to exclude the sender
// 	}
//
// 	h.broadcast <- response
// })
//
// hub.On("/signin", func(m *Message, h *Hub, c *Client) {
// 	c.User.CheckCredentials(m.Data["username"], m.Data["password"])
//
// 	response := &Message{
// 		Path: "/authStatusChanged",
// 		Data: map[string]string{
// 			"isAuthorized": strconv.FormatBool(c.User.IsAuthorized()),
// 		},
// 	}
//
// 	//Send the response just to the recipient
// 	c.send <- response.ToByte()
// })
