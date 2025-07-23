package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	//ID        uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	FirstName string    `gorm:"size:255"`
	LastName  string    `gorm:"size:255"`
	Email     string    `gorm:"size:255; unique"`
	Password  string    `gorm:"size:255"`
}

func (User) TableName() string {
	return "users"
}
