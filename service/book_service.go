package service

import (
	"errors"
	"pojok-baca-api/dto"
	"pojok-baca-api/model"
	"pojok-baca-api/repository"
)

type BookService interface {
	GetBooks() ([]model.Book, error)
	Create(book model.Book) (model.Book, error)
	GetBookByID(id uint) (model.Book, error)
	DeleteBookByID(id uint) error
	UpdateBookByID(req dto.UpdateBookRequest, id uint) (model.Book, error)
}

type bookService struct {
	repo repository.BookRepository
}

func NewBookService(r repository.BookRepository) BookService {
	return &bookService{repo: r}
}

func (s *bookService) GetBooks() ([]model.Book, error) {
	return s.repo.GetAll()
}

func (s *bookService) Create(book model.Book) (model.Book, error) {

	if book.Name == "" || book.Stok == 0 || book.Category == "" || book.RentalCost == 0 {
		return model.Book{}, errors.New("name, stok, category, rental cost required")
	}

	return s.repo.Create(book)
}

func (s *bookService) GetBookByID(id uint) (model.Book, error) {
	return s.repo.GetByID(id)
}

func (s *bookService) DeleteBookByID(id uint) error {
	return s.repo.Delete(id)
}

func (s *bookService) UpdateBookByID(req dto.UpdateBookRequest, id uint) (model.Book, error) {
	book := model.Book{
		Name:       req.Name,
		Stok:       req.Stok,
		RentalCost: req.RentalCost,
		Category:   req.Category,
	}
	return s.repo.Update(book, id)
}
