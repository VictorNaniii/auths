package entity

import (
	"github.com/google/uuid"
	"time"
)

type AuthToken struct {
	ID    uint   `gorm:"primaryKey;autoIncrement:true;unique"`
	Token string `gorm:"gorm:type:text"`

	ExpireDate time.Time `gorm:"type:timestamp"`

	UserId uuid.UUID `gorm:"type:uuid;not null"`
	User   User      `gorm:"foreignKey:UserId;references:ID"`
}

func (AuthToken) TableName() string {
	return "auth_token"
}
