package main

import (
	"book_app/author"
	"book_app/books"
	"book_app/user"
	"net/http"
)

func main() {
	http.HandleFunc("/", user.Top)
	http.HandleFunc("/signup", user.Signup)
	http.HandleFunc("/login", user.Login)
	http.HandleFunc("/books", books.Index)
	http.HandleFunc("/books/create", books.Create)
	http.HandleFunc("/books/create/process", books.CreateProcess)
	http.HandleFunc("/books/show", books.Show)
	http.HandleFunc("/books/update", books.Update)
	http.HandleFunc("/books/update/process", books.UpdateProcess)
	http.HandleFunc("/books/delete/process", books.DeleteProcess)
	http.HandleFunc("/author/show", author.Show)
	http.ListenAndServe(":8080", nil)
}
