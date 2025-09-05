package models

import (
	"time"
)

type RefreshToken struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Token     string    `json:"-" gorm:"type:text;uniqueIndex;not null"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsRevoked bool      `json:"is_revoked" gorm:"default:false"`

	User User `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
