package service

import (
	"auth/config"
	"auth/internal/auth"
	"auth/internal/repository"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
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

func (r *AuthService) Login(data auth.LoginUser) (bool, string, error) {
	isValid := data.Validate()

	if isValid != nil {
		return false, "", fmt.Errorf("Invalid data user: %w", isValid)
	}

	isRightPassword, _ := r.repo.ChekPasswordHas(data.Email, data.Password)

	if !isRightPassword {
		return false, "", errors.New("Invalid User")
	}
	getUserId, err := r.repo.GetUserId(data)

	accessToken, err := GenerateAccessToken(getUserId)
	if err != nil {
		return false, "", err
	}
	rawToken, refreshToken := CreateRefreshToken()

	if err := r.repo.StoreRefreshToken(accessToken, rawToken, refreshToken); err != nil {
		return false, "", fmt.Errorf("Store Refresh Token: %w", err)
	}

	return true, accessToken, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CreateRefreshToken() (string, string) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		panic(err)
	}

	rawString := base64.URLEncoding.EncodeToString(raw)
	hash := sha256.Sum256([]byte(rawString + config.SaltForRefreshToken))
	refreshToken := base64.URLEncoding.EncodeToString(hash[:])

	return rawString, refreshToken
}

func GenerateAccessToken(userId uuid.UUID) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"subj": userId,
		"iss":  "auth-app",
		"exp":  time.Now().Add(config.ExpJWT).Unix(),
		"iat":  time.Now().Unix(),
		"jti":  uuid.New().String(),
	})

	tokenString, err := claims.SignedString(config.SecretKey)
	if err != nil {
		return "", err
	}
	//fmt.Println("Create Token: ", tokenString)

	return tokenString, nil
}
