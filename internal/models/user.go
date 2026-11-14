package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username    string    `json:"username"`
	Password    string    `json:"password,omitempty"`
	Phone       string    `gorm:"unique;not null" json:"phone"`
	DateOfBirth time.Time `json:"dob"`
	Gender      string    `json:"gender"`
	Email       string    `gorm:"unique;not null" json:"email"`
	Occupation  string    `json:"occupation"`
	Role        string    `default:"user" json:"role"`
	ProvinceID  string    `json:"province_id"`
	CityID      string    `json:"city_id"`
	Store       *Store    `json:"store" gorm:"foreignKey:UserID"`

	Address     []Address     `json:"address,omitempty"`
	Transaction []Transaction `json:"transaction,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserResponse struct {
	UserID       uint   `json:"user_id"`
	Username     string `json:"username"`
	Role         string `json:"role"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
