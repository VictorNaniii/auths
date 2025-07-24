package repository

import (
	"auth/internal/auth"
	"auth/internal/entity"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"

	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db}
}

func (r *AuthRepository) Register(data auth.RegisterUser) error {
	user := entity.User{
		//ID:        uuid.New(),
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Password:  data.Password,
	}

	result := r.db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (c *AuthRepository) CheckCredentials(data auth.LoginUser) (bool, error) {
	isExist := c.db.Where("email = ? and password = ? ", data.Email, data.Password).First(&entity.User{})

	if isExist.Error != nil {
		return false, isExist.Error
	}

	return true, nil
}
func (c *AuthRepository) CheckEmail(email string) (bool, error) {
	isEmail := c.db.Where("email = ? ", email).First(&entity.User{})
	if isEmail.Error != nil {
		return false, isEmail.Error
	}
	return true, nil
}
func (c *AuthRepository) ChekPasswordHas(email string, password string) (bool, error) {
	var user entity.User
	isPassword := c.db.Where("email = ? ", email).First(&user)
	if isPassword.Error != nil {
		return false, isPassword.Error
	}

	if CheckPasswordHash(password, user.Password) {
		return true, nil
	}

	return false, nil
}

func (c *AuthRepository) GetUserId(data auth.LoginUser) (uuid.UUID, error) {
	var user entity.User
	isExist := c.db.Where("email = ? ", data.Email).First(&user)
	if isExist.Error != nil {
		return uuid.Nil, isExist.Error
	}

	return user.ID, nil
}

func (c *AuthRepository) StoreRefreshToken(userId uuid.UUID, refreshToken string, expirationDate time.Time) error {
	var saveAuthToken = entity.AuthToken{}

	saveAuthToken.UserId = userId
	saveAuthToken.Token = refreshToken
	saveAuthToken.ExpireDate = expirationDate

	result := c.db.Create(&saveAuthToken)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (c *AuthRepository) FindRefreshToken(refreshToken string) (*entity.AuthToken, error) {
	var authToken entity.AuthToken
	isExist := c.db.Where("token = ?", refreshToken).First(&authToken)
	if isExist.Error != nil {
		return nil, isExist.Error
	}

	return &authToken, nil
}

func (c *AuthRepository) DeleteRefresh(userId uuid.UUID, token string) error {
	result := c.db.Where("token = ? AND user_id = ?", token, userId).Delete(&entity.AuthToken{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
