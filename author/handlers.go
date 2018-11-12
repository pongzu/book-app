package author

import (
	"book_app/commonStruct"
	"book_app/config"
	"errors"
	"log"
	"net/http"
	"strconv"
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

func Create(name string) error {
	_, err := config.DB.Exec("INSERT INTO authors (name) VALUES ($1)", name)
	if err != nil {
		return errors.New("500. Internal Server Error")
	}
	return nil
}

func Show(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))

	author, err := GetAuthor(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	rows, err := config.DB.Query("SELECT * FROM books WHERE author_id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var bks []commonStruct.Book
	for rows.Next() {
		bk := commonStruct.Book{}

		if err := rows.Scan(&bk.Id, &bk.Title, &bk.Author, &bk.Price, &bk.Author_id); err != nil {
			log.Println(err)
		}
		bks = append(bks, bk)
	}
	if err = rows.Err(); err != nil {
		log.Println(err)
	}
	log.Println(bks)

	author.Books = bks

	config.TPL.ExecuteTemplate(w, "author.gohtml", bks)

}
