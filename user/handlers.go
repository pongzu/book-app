package user

import (
	"net/http"

	"github.com/satori/go.uuid"

	"trip_app/config"
)

var dbSessions = make(map[string]string)

func Top(w http.ResponseWriter, r *http.Request) {
	config.TPL.ExecuteTemplate(w, "top.gohtml", nil)
}

func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	u, err := CreateUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// create session
	var c *http.Cookie
	uuid, err := uuid.NewV4()
	c = &http.Cookie{
		Name:  "session",
		Value: uuid.String(),
	}
	http.SetCookie(w, c)
	//set dbSession
	dbSessions[c.Value] = u.Un
	// Redirect
	http.Redirect(w, r, "/books", http.StatusSeeOther)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(r) {
		http.Redirect(w, r, "/index", 303)
	}
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	u, err := GetUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var c *http.Cookie
	c, err = r.Cookie("session")
	if err != nil {
		uuid, err := uuid.NewV4()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		c = &http.Cookie{
			Name:  "session",
			Value: uuid.String(),
		}
		http.SetCookie(w, c)
	}
	dbSessions[c.Value] = u.Un

	http.Redirect(w, r, "/books", http.StatusSeeOther)
}
