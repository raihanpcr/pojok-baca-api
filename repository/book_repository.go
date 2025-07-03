package repository

import (
	"gorm.io/gorm"
	"pojok-baca-api/model"
)

type BookRepository interface {
	GetAll() ([]model.Book, error)
	Create(book model.Book) (model.Book, error)
	GetByID(id uint) (model.Book, error)
	Delete(id uint) error
	Update(book model.Book, id uint) (model.Book, error)
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db}
}

func (r *bookRepository) GetAll() ([]model.Book, error) {
	var books []model.Book
	err := r.db.Find(&books).Error
	return books, err
}

func (r *bookRepository) Create(book model.Book) (model.Book, error) {
	err := r.db.Create(&book).Error
	return book, err
}
func (r *bookRepository) GetByID(id uint) (model.Book, error) {
	var book model.Book
	err := r.db.Where("id = ?", id).First(&book).Error
	return book, err
}

func (r *bookRepository) Delete(id uint) error {
	err := r.db.Where("id = ?", id).Delete(&model.Book{}).Error
	if err != nil {
		return err
	}
	return r.db.Delete(&model.Book{}, id).Error
}

func (r *bookRepository) Update(book model.Book, id uint) (model.Book, error) {
	var b model.Book
	if err := r.db.First(&b, id).Error; err != nil {
		return b, err
	}

	b.Name = book.Name
	b.Stok = book.Stok
	b.RentalCost = book.RentalCost
	b.Category = book.Category
	
	if err := r.db.Save(&b).Error; err != nil {
		return b, err
	}

	return b, nil
}
