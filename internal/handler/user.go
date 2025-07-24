package handler

import (
	"auth/internal/auth"
	"auth/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
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

func (h *UserHandler) Login(data auth.LoginUser) (string, string, error) {

	loginUser, accessToken, refreshToken, err := h.authService.Login(data)

	if err != nil {
		return "", "", err
	}

	if !loginUser {
		return "", "", errors.New("Wrong credentials")
	}

	return accessToken, refreshToken, nil
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

func (h *UserHandler) RefreshToken(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	access, refresh, err := h.authService.Refresh(body.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (h *UserHandler) Logout(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authService.Logout(body.RefreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
