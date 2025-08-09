package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name           string `gorm:"not null"`
	Email          string `gorm:"unique;not null"`
	HashedPassword string `gorm:"not null"`
	Schedule       string `gorm:"type:varchar(10);check:schedule IN ('DAILY','WEEKLY','BIWEEKLY','MINUTE');not null;default:'MINUTE'"`

	Posts []Post `gorm:"foreignKey:UserID"`

	Suggestions []Suggestion `gorm:"foreignKey:UserID"`

	SearchResults []SearchResult `gorm:"foreignKey:UserID"`
}

// CreateUserDto
type CreateUserDto struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginDto struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
