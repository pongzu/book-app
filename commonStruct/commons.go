package commonStruct

type Book struct {
	Id        int
	Title     string
	Author    string
	Price     float64
	Author_id int
}

type Author struct {
	Id    int
	Name  string
	Books []Book
}

type Comment struct {
	Id      int
	Comment string
	BookId  int
	UserId  int
}
