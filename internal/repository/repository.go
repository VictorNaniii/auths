package repository

import (
	"auth/internal/entity"
	"auth/internal/model"
)

type BooksRepository interface {
	AddBook(book model.BookRes) (string, error)
	DeleteBook(id string) error
	GetBook(id string) (entity.Book, error)
	GetAllBooks() ([]entity.Book, error)
	EditBook(id string, book model.ChangeData) (entity.Book, error)
}

type UserRepository interface {
	Signup(user entity.User) (string, error)
	Login(user entity.User) (string, error)
}
