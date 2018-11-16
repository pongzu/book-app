package user

import (
	"book_app/config"
	"errors"
	"fmt"
	"log"
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

// check already logged in
func alreadyLoggedIn(r *http.Request) bool {
	c, err := r.Cookie("session")
	if err != nil {
		return false
	}
	fmt.Println(sessions[c.Value])

	_, ok := sessions[c.Value]
	return ok
}

//get user that has already logged in
func GetCurrentUser(r *http.Request) (User, error) {
	// initialized user struct and put values from data base
	var u User

	c, err := r.Cookie("session")
	if err != nil {
		return u, errors.New("400. Bad Request. can not get cookie value")
	}
	// get userName from session value
	s := sessions[c.Value]
	row := config.DB.QueryRow("SELECT * FROM users WHERE username = $1", s.un)

	// just for scanning
	var p string
	if err := row.Scan(&u.Id, &u.Un, &u.Email, &p); err != nil {
		return u, errors.New("500. Internal Server error")
	}
	u.Password = []byte(p)
	log.Println("ここおおっっこっここここ")
	log.Println(u)

	return u, nil
}
