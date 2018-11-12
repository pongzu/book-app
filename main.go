package main

import (
	"book_app/books"
	"net/http"
)

func main() {
	http.HandleFunc("/", books.Index)
	http.HandleFunc("/books/create", books.Create)
	http.HandleFunc("/books/create/process", books.CreateProcess)
	http.HandleFunc("/books/show", books.Show)
	http.HandleFunc("/books/update", books.Update)
	http.HandleFunc("/books/update/process", books.UpdateProcess)
	http.HandleFunc("/books/delete/process", books.DeleteProcess)
	http.ListenAndServe(":8080", nil)
}
