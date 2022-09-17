package books

type bookService struct {
	bookRepository BookRepository
}

func NewBookService(r BookRepository) BookService {
	return &bookService{
		bookRepository: r,
	}
}

func (s *bookService) CreateBook(book *Book) (err error) {

	return
}

func (s *bookService) UpdateBook(bookId uint, book *Book) (err error) {

	return
}

func (s *bookService) DeleteBook(bookId uint) (err error) {

	return
}

func (s *bookService) GetBook(bookId uint) (book *Book, err error) {

	return
}

func (s *bookService) GetBooks(title, author string) (books *[]Book, err error) {

	return
}
