package author

import (
	"book_app/commonStruct"
	"book_app/config"
	"errors"
)

func GetAuthor(id int) (commonStruct.Author, error) {
	var author = commonStruct.Author{}
	row := config.DB.QueryRow("SELECT * FROM authors where id = $1", id)
	if err := row.Scan(&author.Id, &author.Name); err != nil {
		return author, errors.New("500. Internal Server Error")
	}
	return author, nil
}
