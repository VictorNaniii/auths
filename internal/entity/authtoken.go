package entity

import (
	"github.com/google/uuid"
	"time"
)

type AuthToken struct {
	ID    string `gorm:"primary_key"`
	Token string `gorm:"gorm:type:text"`

	ExpireDate time.Time `gorm:"type:timestamp"`

	UserId uuid.UUID `gorm:"type:uuid;not null"`
	User   User      `gorm:"foreignKey:UserId;references:ID"`
}

func (AuthToken) TableName() string {
	return "auth_token"
}
