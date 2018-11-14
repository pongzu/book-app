package user

import (
	"net/http"

	uuid "github.com/satori/go.uuid"
)

// dbSeiions: uuid for key, user name for value db
// session holds user name for hodling current user imformation
var dbSessions = make(map[string]string)

func (u *User) setSession(w http.ResponseWriter) {
	c := getCockie()
	http.SetCookie(w, c)
	dbSessions[c.Value] = u.Un
}

// returns *http.Cockie
func getCockie() *http.Cookie {
	var c *http.Cookie
	uuid, _ := uuid.NewV4()
	c = &http.Cookie{
		Name:  "session",
		Value: uuid.String(),
	}
	return c
}
