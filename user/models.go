package user

import (
	"book_app/config"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// for trash for scanning
type TrashScanner struct{}

func (TrashScanner) Scan(interface{}) error {
	return nil
}

type User struct {
	Id       int
	Un       string
	Email    string
	Password []byte
}

func CreateUser(w http.ResponseWriter, r *http.Request) (User, error) {
	u := User{}
	u.Un = r.FormValue("username")
	u.Email = r.FormValue("email")
	p := r.FormValue("password")

	if err := blackchech(u.Un, u.Email, p); err != nil {
		return u, err
	}

	// create hash password
	bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
	if err != nil {
		return u, err
	}
	u.Password = bs
	_, err = config.DB.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", u.Un, u.Email, string(u.Password))
	if err != nil {
		return u, err
	}

	return u, nil
}

func GetUser(w http.ResponseWriter, r *http.Request) (User, error) {
	u := User{}
	email := r.FormValue("email")
	password := r.FormValue("password")

	if err := blackchech(email, password); err != nil {
		return u, err
	}

	row := config.DB.QueryRow("SELECT * FROM users WHERE email = $1", email)

	// just for scanning
	var p string

	err := row.Scan(&u.Id, &u.Un, &u.Email, &p)
	if err != nil {
		return u, err
	}
	u.Password = []byte(p)

	err = bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	if err != nil {
		fmt.Println(u.Password)
		return u, errors.New("400. got wrong password")
	}
	return u, nil
}
