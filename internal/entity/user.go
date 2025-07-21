package domain

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID `gorm:"primary_key" `
	FirstName string    `gorm:"size:255"`
	LastName  string    `gorm:"size:255"`
	Email     string    `gorm:"size:255"`
}
