package auth

import (
	"auth/config"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

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

func AuthentificateMiddleware(c *gin.Context) {
	tokenFromCokie, err := c.Cookie("token")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	ok := VerifyToken(tokenFromCokie)

	if ok != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	c.Next()
}

func VerifyToken(token string) error {
	tokenVerify, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return config.SecretKey, nil
	})
	if err != nil {
		return err
	}

	if !tokenVerify.Valid {
		return err
	}

	return nil

}
