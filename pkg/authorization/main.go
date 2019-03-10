package authorization

//User is an authorized user
type User struct {
	UUID       string
	Name       string
	token      string
	authorized bool
}

//CheckCredentials checks whether the given credentials is authorized
func (u *User) CheckCredentials(user string, password string) bool {
	if user == "test" && password == "test" {
		u.UUID = "1234"
		u.authorized = true
	}

	return u.authorized
}

//IsAuthorized sets whether the user is authorized
func (u *User) IsAuthorized() bool {
	return u.authorized
}
