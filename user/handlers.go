package user

import (
	"fmt"
	"net/http"

	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"

	"trip_app/config"
)

var dbSessions = make(map[string]string)

type User struct {
	Un       string
	Email    string
	Password []byte
}

// for trash for scanning
type TrashScanner struct{}

func (TrashScanner) Scan(interface{}) error {
	return nil
}
func Top(w http.ResponseWriter, r *http.Request) {
	config.TPL.ExecuteTemplate(w, "top.gohtml", nil)
}

func Signup(w http.ResponseWriter, r *http.Request) {
	postValidate(w, r)

	u := User{}
	u.Un = r.FormValue("username")
	u.Email = r.FormValue("email")
	p := r.FormValue("password")

	blackchech(w, u.Un, u.Email, p)
	// create hash password
	bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	u.Password = bs
	_, err = config.DB.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", u.Un, u.Email, string(u.Password))
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
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
	// template
	config.TPL.ExecuteTemplate(w, "index.gohtml", u)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(r) {
		http.Redirect(w, r, "/index", 303)
	}

	postValidate(w, r)

	u := User{}
	email := r.FormValue("email")
	password := r.FormValue("password")

	blackchech(w, email, password)

	row := config.DB.QueryRow("SELECT * FROM users WHERE email = $1", email)

	var p string

	err := row.Scan(TrashScanner{}, &u.Un, &u.Email, &p)
	if err != nil {
		http.Error(w, "InternalSevrError", 500)
		return
	}
	u.Password = []byte(p)

	err = bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	if err != nil {
		http.Error(w, "wrong password got here", 400)
		fmt.Println(u.Password)
		return
	}

	var c *http.Cookie
	c, err = r.Cookie("session")
	if err != nil {
		uuid, err := uuid.NewV4()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		c = &http.Cookie{
			Name:  "session",
			Value: uuid.String(),
		}
		http.SetCookie(w, c)
	}
	dbSessions[c.Value] = u.Un

	config.TPL.ExecuteTemplate(w, "index.gohtml", u)
}
