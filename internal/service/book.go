package service

import (
	"auth/internal/entity"
	"auth/internal/model"
	"auth/internal/repository"
	"github.com/google/uuid"
)

type IBookService interface {
	AddBook(book model.BookRes, userId uuid.UUID) (string, error)
	GetAllBooks() ([]entity.Book, error)
	EditBook(id string, book model.ChangeData) (entity.Book, error)
}

type BookService struct {
	repo repository.BooksRepository
}

func NewBookService(repo repository.BooksRepository) IBookService {
	return &BookService{repo: repo}
}

func (s *BookService) AddBook(book model.BookRes, userId uuid.UUID) (string, error) {
	return s.repo.AddBook(book, userId)
}

func (s *BookService) GetAllBooks() ([]entity.Book, error) {
	result, err := s.repo.GetAllBooks()
	if err != nil {
		return []entity.Book{}, err
	}
	return result, nil
}

func (s *BookService) EditBook(id string, book model.ChangeData) (entity.Book, error) {
	itemToUpdate, err := s.repo.EditBook(id, book)
	if err != nil {
		return entity.Book{}, err
	}
	return itemToUpdate, nil
}
