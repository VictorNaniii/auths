package repository

import (
	"auth/internal/entity"
	"auth/internal/model"
	"github.com/google/uuid"
)

type BooksRepository interface {
	AddBook(book model.BookRes, userId uuid.UUID) (string, error)
	DeleteBook(id string) error
	GetBook(id string) (entity.Book, error)
	GetAllBooks() ([]entity.Book, error)
	EditBook(id string, book model.ChangeData) (entity.Book, error)
}

type UserRepository interface {
	Signup(user entity.User) (string, error)
	Login(user entity.User) (string, error)
}
