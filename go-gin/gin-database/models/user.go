package models

import (
	"time"
)

type UserRole string

const (
	AdminRole       UserRole = "admin"
	RegularUserRole UserRole = "user"
)

type User struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name" gorm:"size:100;not null"`
	Email     string     `json:"email" gorm:"size:120;uniqueIndex;not null"`
	Password  string     `json:"-" gorm:"size:255;not null"`
	Role      UserRole   `json:"role" gorm:"type:enum('admin','user');default:'user';not null"`
	Birthday  *time.Time `json:"birthday" gorm:"type:date"`
	Phone     *string    `json:"phone" gorm:"size:30"`
	IsActive  bool       `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
