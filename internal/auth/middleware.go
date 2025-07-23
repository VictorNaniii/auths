package auth

import "github.com/go-playground/validator/v10"

type RegisterUser struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Password  string `json:"password" validate:"required,min=8,max=32"`
	Email     string `json:"email" validate:"required,email"`
}

type LoginUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type AuthUser interface {
	Register(data RegisterUser) (string, error)
	Login(data LoginUser) (string, error)
}

var validate = validator.New()

func (i RegisterUser) Validate() error {
	err := validate.Struct(i)
	if err != nil {
		return err
	}
	return nil
}

func (i LoginUser) Validate() error {
	err := validate.Struct(i)
	if err != nil {
		return err
	}
	return nil
}
