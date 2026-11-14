package models

import "time"

type Store struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Photo     *string   `json:"photo"`
	UserID    uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
