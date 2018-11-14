package user

import (
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

// sessions: uuid for key, session for value
// session holds user name,time.Now() for hodling current user imformation
type session struct {
	un           string
	lastActicity time.Time
}

var sessions = make(map[string]session)

func (u *User) setSession(w http.ResponseWriter) {
	c := getCockie()
	http.SetCookie(w, c)

	sessions[c.Value] = session{
		un:           u.Un,
		lastActicity: time.Now(),
	}
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

// delite sessions once a hour
func clearnSession(sessions map[string]session) {
	for cookie, session := range sessions {
		if time.Now().Sub(session.lastActicity) > (time.Minute * 60) {
			delete(sessions, cookie)
		}
	}
}
