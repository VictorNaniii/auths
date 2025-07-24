package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	FirstName string `gorm:"size:255"`
	LastName  string `gorm:"size:255"`
	Email     string `gorm:"size:255;unique"`
	Password  string `gorm:"size:255"`

	// Adding relation to Books
	Books []Book `gorm:"foreignKey:UserId"`
}

func (User) TableName() string {
	return "users"
}
