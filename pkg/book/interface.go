package book

type BookRepository interface {
	GetBooks(title, author string) (*[]Book, error)
	GetBook(bookID uint) (*Book, error)
	CreateBook(book *Book) error
	UpdateBook(bookID uint, book *Book) error
	DeleteBook(bookID uint) error
}

type BookService interface {
	GetBooks(title, author string) (*[]Book, error)
	GetBook(bookID uint) (*Book, error)
	CreateBook(book *Book) error
	UpdateBook(bookID uint, book *Book) error
	DeleteBook(bookID uint) error
}
