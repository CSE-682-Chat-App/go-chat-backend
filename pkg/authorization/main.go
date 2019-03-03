package authorization

type User struct {
	Name       string
	token      string
	authorized bool
}

//CheckCredentials checks whether the given credentials is authorized
func (u *User) CheckCredentials(user string, password string) bool {
	if user == "test" && password == "test" {
		u.authorized = true
	}

	return u.authorized
}

//IsAuthorized sets whether the user is authorized
func (u *User) IsAuthorized() bool {
	return u.authorized
}
