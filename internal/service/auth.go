package service

import (
	"auth/internal/auth"
	"auth/internal/repository"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (r *AuthService) Register(data auth.RegisterUser) (bool, error) {
	isValid := data.Validate()

	if isValid != nil {
		return false, isValid
	}

	pssword, ok := HashPassword(data.Password)

	if ok != nil {
		return false, isValid
	}

	data.Password = pssword

	err := r.repo.Register(data)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *AuthService) Login(data auth.LoginUser) (bool, error) {
	isValid := data.Validate()

	if isValid != nil {
		return false, errors.New("Invalid User")
	}

	isRightPassword, _ := r.repo.ChekPasswordHas(data.Email, data.Password)

	if !isRightPassword {
		return false, errors.New("Invalid User")
	}

	return true, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
