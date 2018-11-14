package user

import (
	"errors"
	"log"
	"net/http"
)

// return err if any blancs from input are found
func blackchech(inputs ...string) error {
	for _, v := range inputs {
		if v == "" {
			return errors.New(http.StatusText(405))
		}
	}
	log.Println("no balank for inputs")
	return nil
}

// check already logged in
func alreadyLoggedIn(r *http.Request) bool {
	c, err := r.Cookie("session")
	if err != nil {
		return false
	}
	_, ok := dbSessions[c.Value]
	return ok
}
