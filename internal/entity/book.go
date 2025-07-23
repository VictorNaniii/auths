package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	ID          uuid.UUID ` gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Title       string    `gorm:"size:255"`
	Author      string    `gorm:"size:255"`
	Description string    `gorm:"size:255"`
}

func (Book) TableName() string {
	return "book_res"
}
