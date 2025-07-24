package repository

import (
	"auth/internal/entity"
	"auth/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//type IBookRepository interface {
//	AddBook(book model.BookRes) (string, error)
//	GetAllBook() ([]entity.Book, error)
//	//DeleteBook(id string) error
//	EditBook(id string, book model.BookRes) (entity.Book, error)
//}

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) AddBook(book model.BookRes) (string, error) {
	//var createBook entity.Book
	createBook := entity.Book{
		Title:       book.Title,
		Author:      book.Author,
		Description: book.Description,
	}
	if err := r.db.Create(&createBook).Error; err != nil {
		return "error", err
	}
	return "Book was successfully added: " + book.Title, nil
}

func (r *BookRepository) GetAllBook() ([]entity.Book, error) {
	var book []entity.Book
	r.db.Raw("SELECT * FROM book_res").Scan(&book)
	return book, nil
}

func (r *BookRepository) DeleteBook(id string) error {
	uuidChange, err := uuid.FromBytes([]byte(id))
	if err != nil {
		return err
	}

	var result = entity.Book{
		ID: uuidChange,
	}

	if err := r.db.Delete(&result).Error; err != nil {
		return err
	}

	return nil
}

func (r *BookRepository) EditBook(id string, changeData model.ChangeData) (entity.Book, error) {
	uuidChange, err := uuid.FromBytes([]byte(id))
	if err != nil {
		return entity.Book{}, err
	}

	var book entity.Book
	if err := r.db.First(&book, "id = ?", uuidChange).Error; err != nil {
		return entity.Book{}, err
	}

	if changeData.Title != nil {
		book.Title = *changeData.Title
	}
	if changeData.Author != nil {
		book.Author = *changeData.Author
	}
	if changeData.Description != nil {
		book.Description = *changeData.Description
	}

	if err := r.db.Save(&book).Error; err != nil {
		return entity.Book{}, err
	}

	return book, nil
}

func (r *BookRepository) GetBook(id string) (entity.Book, error) {
	var book entity.Book
	if err := r.db.First(&book, "id = ?", id).Error; err != nil {
		return entity.Book{}, err
	}
	return book, nil
}

func (r *BookRepository) GetAllBooks() ([]entity.Book, error) {
	var books []entity.Book
	if err := r.db.Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}
