package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Book struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Title       string `gorm:"size:255"`
	Author      string `gorm:"size:255"`
	Description string `gorm:"size:255"`

	UserId uuid.UUID `gorm:"type:uuid;not null"`
	User   User      `gorm:"foreignKey:UserId;references:ID"`
}

func (Book) TableName() string {
	return "book_res"
}
