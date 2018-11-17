package books

import (
	"book_app/author"
	"book_app/commonStruct"
	"book_app/config"
	"errors"
	"log"
	"net/http"
	"strconv"
)

func AllBooks() ([]commonStruct.Book, error) {
	rows, err := config.DB.Query("SELECT * FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bks []commonStruct.Book
	for rows.Next() {
		bk := commonStruct.Book{}
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

func OneBook(r *http.Request) (commonStruct.Book, error) {
	var bk = commonStruct.Book{}
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

func PutBook(r *http.Request) (commonStruct.Book, error) {
	var bk = commonStruct.Book{}

	bk.Title = r.FormValue("title")
	bk.Author = r.FormValue("author")
	p := r.FormValue("price")

	/// cretate author_id
	if author.Find(bk.Author) {
		id, err := author.GetId(bk.Author)
		if err != nil {
			log.Println(err)
			return bk, err
		}
		bk.Author_id = id
	} else {
		if err := author.Create(bk.Author); err != nil {
			log.Println(err)
			return bk, err
		}
		id, err := author.GetId(bk.Author)
		if err != nil {
			log.Println(err)
			return bk, err
		}
		bk.Author_id = id
	}

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
	// to get serial id
	row := config.DB.QueryRow("select books.id from books where title = $1", bk.Title)
	if err := row.Scan(&bk.Id); err != nil {
		return bk, err
	}

	return bk, nil
}

// putAuthor(r){
// }

func PutComment(r *http.Request, id int) (commonStruct.Comment, error) {
	bookId := id
	userId, _ := strconv.Atoi(r.FormValue("current_user_id"))
	comment := r.FormValue("comment")

	var c commonStruct.Comment

	_, err := config.DB.Exec("INSERT INTO comments (comment, book_id, user_id) VALUES ($1, $2, $3)", comment, bookId, userId)
	if err != nil {
		return c, err
	}
	// to get serial id
	var i int
	row := config.DB.QueryRow("select comments.id from comments where comment = $1", comment)
	if err := row.Scan(&i); err != nil {
		return c, err
	}
	// create structre just for passing

	c.Id = i
	c.Comment = comment
	c.BookId = bookId
	c.UserId = userId

	return c, nil
}

func PutRelation(c commonStruct.Comment) error {
	_, err := config.DB.Exec("INSERT INTO relations (book_id, user_id, comment_id) VALUES ($1, $2, $3)", c.BookId, c.UserId, c.Id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateBook(r *http.Request) (commonStruct.Book, error) {
	var bk = commonStruct.Book{}
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
