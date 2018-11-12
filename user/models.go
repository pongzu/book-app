package user

import (
	"log"
	"net/http"
)

func postValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
}

func blackchech(w http.ResponseWriter, inputs ...string) {
	for _, v := range inputs {
		if v == "" {
			http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
			return
		}
	}
	log.Println("no balank for inputs")
}

func alreadyLoggedIn(r *http.Request) bool {
	c, err := r.Cookie("session")
	if err != nil {
		return false
	}
	_, ok := dbSessions[c.Value]
	return ok
}
