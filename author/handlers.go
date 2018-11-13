package author

import (
	"book_app/commonStruct"
	"book_app/config"
	"log"
	"net/http"
	"strconv"
)

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
