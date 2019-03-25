package authorization

import (
	"encoding/json"
	"fmt"
	ws "github.com/CSE-682-Chat-App/go-chat-backend/internal/pkg/websocket"
	uuid "github.com/google/uuid"
	"log"
	// "reflect"
)

//User is an authorized user
type User struct {
	UUID       string `json:"uuid"`
	Name       string `json:"name"`
	password   string
	authorized bool
}

func (u *User) ToString() string {
	str, err := json.Marshal(u)

	if err != nil {
		return ""
	}
	return string(str)
}

var users = map[string]*User{
	"test": &User{
		UUID:     uuid.New().String(),
		Name:     "test",
		password: "test",
	},
	"test2": &User{
		UUID:     uuid.New().String(),
		Name:     "test2",
		password: "test",
	},
}

//CheckCredentials checks whether the given credentials is authorized
func (u *User) CheckCredentials(password string) bool {
	u.authorized = u.password == password
	return u.authorized
}

func GetUser(m *ws.Message, h *ws.Hub, c *ws.Client) {
	user := c.User.(*User)
	response := &ws.Message{
		Path: "/user_response",
		Data: map[string]string{
			"user": user.ToString(),
		},
	}

	c.Send <- response.ToByte()
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
		log.Println(fmt.Sprintf("User %s signup up", username))
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

func AuthMiddleware(handler ws.OnCallback) ws.OnCallback {
	return func(m *ws.Message, h *ws.Hub, c *ws.Client) {
		if c.User == nil {
			return
		}
		user := c.User.(*User)

		if !user.authorized {
			response := &ws.Message{
				Path: "/auth_error",
				Data: map[string]string{
					"message":    "User is Unauthorized",
					"authorized": "false",
				},
			}

			c.Send <- response.ToByte()
			return
		}

		handler(m, h, c)
	}
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
