package books

import (
	"book_app/commonStruct"
	"book_app/config"
	"book_app/user"
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

var CurrentUser *user.User

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	bks, err := AllBooks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get current_user_id
	currentUser, err := user.GetCurrentUser(r)
	if err != nil {
		fmt.Printf("can not get current user get with this reason: %v", err)
	}

	// create data structure to pass to template
	data := struct {
		BOOKS       []commonStruct.Book
		CurrentUser user.User
	}{
		bks,
		currentUser,
	}

	config.TPL.ExecuteTemplate(w, "books.gohtml", data)
}

func Create(w http.ResponseWriter, r *http.Request) {
	u, err := user.GetCurrentUser(r)
	if err != nil {
		fmt.Printf("can not get current user get with this reason: %v", err)
	}
	config.TPL.ExecuteTemplate(w, "create.gohtml", u)
}

func CreateProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	bk, err := PutBook(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	c, err := PutComment(r, bk.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	if err := PutRelation(c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	http.Redirect(w, r, "/books", http.StatusSeeOther)
	return
}

func Show(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	bk, err := OneBook(r)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	config.TPL.ExecuteTemplate(w, "show.gohtml", bk)

	fmt.Print(bk)

}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	bk, err := OneBook(r)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	config.TPL.ExecuteTemplate(w, "update.gohtml", bk)
}

func UpdateProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	_, err := UpdateBook(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	http.Redirect(w, r, "/books", http.StatusSeeOther)
}

func DeleteProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	err := DeleteBook(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/books", http.StatusSeeOther)
}
