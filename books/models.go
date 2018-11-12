package books

import (
	"book_app/author"
	"book_app/config"
	"errors"
	"net/http"
	"strconv"
)

type Book struct {
	Id        int
	Title     string
	Author    string
	Price     float64
	Author_id int
}

func AllBooks() ([]Book, error) {
	rows, err := config.DB.Query("SELECT * FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bks []Book
	for rows.Next() {
		bk := Book{}
		if err := rows.Scan(&bk.Id, &bk.Title, &bk.Author, &bk.Price, &bk.Author_id); err != nil {
			return nil, err
		}
		bks = append(bks, bk)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return bks, nil
}

func PutBook(r *http.Request) (Book, error) {
	var bk = Book{}

	bk.Title = r.FormValue("title")
	bk.Author = r.FormValue("author")
	p := r.FormValue("price")
	//bk.Author_id = 1 // ひとまず

	author.FindAuthor(bk.Author)

	if bk.Title == "" || bk.Author == "" || p == "" {
		return bk, errors.New("400. Bad request. All fields must be complete.")
	}

	f64, err := strconv.ParseFloat(p, 32)
	if err != nil {
		return bk, errors.New("406. Not Acceptable. Price must be a number.")
	}
	bk.Price = f64

	_, err = config.DB.Exec("INSERT INTO books (title, author, price, author_id) VALUES ($1, $2, $3, $4)", bk.Title, bk.Author, bk.Price, bk.Author_id)
	if err != nil {
		return bk, err
	}
	return bk, nil
}

func OneBook(r *http.Request) (Book, error) {
	bk := Book{}
	id, _ := strconv.Atoi(r.FormValue("id"))

	if id == 0 {
		return bk, errors.New("400. Bad Reqest")
	}

	row := config.DB.QueryRow("SELECT * FROM books WHERE id = $1", id)
	if err := row.Scan(&bk.Id, &bk.Title, &bk.Author, &bk.Price, &bk.Author_id); err != nil {
		return bk, err
	}
	return bk, nil
}

func UpdateBook(r *http.Request) (Book, error) {
	bk := Book{}
	bk.Id, _ = strconv.Atoi(r.FormValue("id"))

	bk.Title = r.FormValue("title")
	bk.Author = r.FormValue("author")
	p := r.FormValue("price")
	bk.Author_id = 1

	if bk.Id == 0 || bk.Title == "" || bk.Author == "" || p == "" {
		return bk, errors.New("400. Bad Request. Fields can't be empty.")
	}

	// convert form values
	f64, err := strconv.ParseFloat(p, 32)
	if err != nil {
		return bk, errors.New("406. Not Acceptable. Enter number for price.")
	}
	bk.Price = f64

	// insert values
	_, err = config.DB.Exec("UPDATE books SET id = $1, title=$2, author=$3, price=$4, author_id = $5  WHERE id=$1;", bk.Id, bk.Title, bk.Author, bk.Price, bk.Author_id)
	if err != nil {
		return bk, err
	}
	return bk, nil
}

func DeleteBook(r *http.Request) error {
	id, _ := strconv.Atoi(r.FormValue("id"))
	if id == 0 {
		return errors.New("400. Bad Request.")
	}

	_, err := config.DB.Exec("DELETE FROM books WHERE id=$1;", id)
	if err != nil {
		return errors.New("500. Internal Server Error")
	}
	return nil
}
