package author

import (
	"book_app/books"
	"book_app/config"
	"errors"
)

type Author struct {
	Id    int
	Name  string
	Books []books.Book
}

func FindAuthor(name string) bool {
	var author = Author{}

	row := config.DB.QueryRow("SELECT * FROM authors WHERE name = $1", name)
	if err := row.Scan(&author.Id, &author.Name); err != nil {
		return false
	}
	return name == author.Name
}

func CreateProcess(name string) {
	_, err = config.DB.Exec("INSERT INTO authors (name) VALUES ($1)", name)
	if err != nil {
		return errors.New("500. Internal Server Error")
	}
}
