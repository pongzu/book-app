package author

import (
	"book_app/commonStruct"
	"book_app/config"
	"errors"
)

func Find(name string) bool {
	var author = commonStruct.Author{}

	row := config.DB.QueryRow("SELECT * FROM authors WHERE name = $1", name)
	if err := row.Scan(&author.Id, &author.Name); err != nil {
		return false
	}
	return true
}

func GetId(name string) (int, error) {
	var author = commonStruct.Author{}
	row := config.DB.QueryRow("SELECT * FROM authors WHERE name = $1", name)
	if err := row.Scan(&author.Id, &author.Name); err != nil {
		return 0, errors.New("500. Internal Server Error")
	}
	return author.Id, nil
}

func GetAuthor(id int) (commonStruct.Author, error) {
	var author = commonStruct.Author{}
	row := config.DB.QueryRow("SELECT * FROM authors where id = $1", id)
	if err := row.Scan(&author.Id, &author.Name); err != nil {
		return author, errors.New("500. Internal Server Error")
	}
	return author, nil
}

func Create(name string) error {
	_, err := config.DB.Exec("INSERT INTO authors (name) VALUES ($1)", name)
	if err != nil {
		return errors.New("500. Internal Server Error")
	}
	return nil
}
