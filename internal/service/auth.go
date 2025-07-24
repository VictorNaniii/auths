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

func (r *AuthService) Login(data auth.LoginUser) (bool, string, string, error) {
	isValid := data.Validate()

	if isValid != nil {
		return false, "", "", fmt.Errorf("Invalid data user: %w", isValid)
	}

	isRightPassword, _ := r.repo.ChekPasswordHas(data.Email, data.Password)

	if !isRightPassword {
		return false, "", "", errors.New("Invalid User")
	}
	getUserId, err := r.repo.GetUserId(data)

	accessToken, err := GenerateAccessToken(getUserId)
	if err != nil {
		return false, "", "", err
	}
	rawRefreshToken, hashedRefreshToken := CreateRefreshToken()
	expirationDate := time.Now().Add(config.ExpireRefreshToken)
	if err := r.repo.StoreRefreshToken(getUserId, hashedRefreshToken, expirationDate); err != nil {
		return false, "", "", fmt.Errorf("Store Refresh Token: %w", err)
	}

	return true, accessToken, rawRefreshToken, nil
}

func (s *AuthService) Refresh(oldToken string) (access, refresh string, err error) {
	// Hash the incoming raw token to find it in the database
	hash := sha256.Sum256([]byte(oldToken + config.SaltForRefreshToken))
	hashedToken := base64.URLEncoding.EncodeToString(hash[:])

	rt, err := s.repo.FindRefreshToken(hashedToken)
	if err != nil {
		return "", "", errors.New("Refresh token is invalid")
	}

	// Check if token is expired
	if time.Now().After(rt.ExpireDate) {
		// Clean up expired token
		if err := s.repo.DeleteRefresh(rt.UserId, hashedToken); err != nil {
			// Log error but continue with the main error response
		}
		return "", "", errors.New("Refresh token has expired")
	}

	// Generate new access token
	access, err = GenerateAccessToken(rt.UserId)
	if err != nil {
		return "", "", err
	}

	// Generate new refresh token
	rawRefreshToken, newHashedRefreshToken := CreateRefreshToken()
	newExpirationDate := time.Now().Add(config.ExpireRefreshToken)

	// Delete old refresh token
	if err := s.repo.DeleteRefresh(rt.UserId, hashedToken); err != nil {
		return "", "", err
	}

	// Store new refresh token
	if err := s.repo.StoreRefreshToken(rt.UserId, newHashedRefreshToken, newExpirationDate); err != nil {
		return "", "", err
	}

	return access, rawRefreshToken, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	// Hash the incoming raw token to find it in the database
	hash := sha256.Sum256([]byte(refreshToken + config.SaltForRefreshToken))
	hashedToken := base64.URLEncoding.EncodeToString(hash[:])

	rt, err := s.repo.FindRefreshToken(hashedToken)
	if err != nil {
		return errors.New("Invalid refresh token")
	}

	// Delete the refresh token
	if err := s.repo.DeleteRefresh(rt.UserId, hashedToken); err != nil {
		return err
	}

	return nil
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
