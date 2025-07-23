package handler

import (
	"auth/internal/entity"
	"auth/internal/model"
	"auth/internal/service"
)

type BookHandler struct {
	service service.IBookService
}

func NewBookHandler(svc service.IBookService) *BookHandler {
	return &BookHandler{service: svc}
}

func (h *BookHandler) AddBook(book model.BookRes) (string, error) {
	return h.service.AddBook(book)
}

func (h *BookHandler) GetAllBook() ([]entity.Book, error) {
	result, err := h.service.GetAllBooks()
	if err != nil {
		return result, err
	}
	return result, nil
}

func (h *BookHandler) EditBook(id string, data model.ChangeData) (entity.Book, error) {
	result, err := h.service.EditBook(id, data)

	if err != nil {
		return result, err
	}

	return result, nil
}
