package auth

import (
	"auth/config"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"strings"
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
func AuthMiddleware(c *gin.Context) {
	authz := c.GetHeader("Authorization")
	parts := strings.SplitN(authz, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err := VerifyToken(parts[1])
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	//c.Set("userID", userID)
	c.Next()
}
func GetUserIDFromToken(c *gin.Context) {
	tokenStr, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token missing"})
		return
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKeyType
		}
		return config.SecretKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		subj, ok := claims["subj"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid subject"})
			return
		}

		userID, err := uuid.Parse(subj)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid UUID"})
			return
		}
		c.Set("userID", userID.String())
		//c.JSON(http.StatusOK, gin.H{"user_id": userID})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to extract claims"})
}
