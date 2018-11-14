package user

import (
	"fmt"
	"net/http"

	"book_app/config"
)

func Top(w http.ResponseWriter, r *http.Request) {
	config.TPL.ExecuteTemplate(w, "top.gohtml", nil)
}

func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	/// insert values from input to DB
	u, err := CreateUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	u.setSession(w)
	// Redirect
	http.Redirect(w, r, "/books", http.StatusSeeOther)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	u, err := GetUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// create session
	_, err = r.Cookie("session")
	if err != nil {
		u.setSession(w)
	}

	fmt.Println(sessions)
	http.Redirect(w, r, "/books", http.StatusSeeOther)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	if !alreadyLoggedIn(r) {
		return
	}

	c, err := r.Cookie("session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	delete(sessions, c.Value)

	// set cookie that has no values
	c = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
