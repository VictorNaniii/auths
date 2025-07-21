package domain

import "github.com/google/uuid"

type Book struct {
	ID          uuid.UUID ` gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Title       string    `gorm:"size:255"`
	Author      string    `gorm:"size:255"`
	Description string    `gorm:"size:255"`
}
