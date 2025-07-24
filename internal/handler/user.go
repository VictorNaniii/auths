package handler

import (
	"auth/internal/auth"
	"auth/internal/service"
)

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type UserHandler struct {
	//service service.IUserService
	authService service.AuthService
}

func NewUserHandler(authService service.AuthService) *UserHandler {
	return &UserHandler{authService: authService}
}

func (h *UserHandler) Login(data auth.LoginUser) (string, error) {

	loginUser, token, err := h.authService.Login(data)

	if err != nil {
		return "error", err
	}

	if !loginUser {
		return "Wrong credentials", nil
	}

	return token, nil
}

func (h *UserHandler) Register(data auth.RegisterUser) (string, error) {

	registerUser, err := h.authService.Register(data)
	if err != nil {
		return "error", err
	}

	if !registerUser {
		return "User alredy exist exist", nil
	}

	return "Register success", nil
}
