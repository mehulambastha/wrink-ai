package models

import "gorm.io/gorm"

type Suggestion struct {
	gorm.Model
	Topic          string `gorm:"not null"`
	SearchKeywords string `gorm:"not null"`
	Search         bool   `gorm:"default:true"`

	UserID uint
	User   User `gorm:"foreignKey:UserID"`

	Post         Post
	SearchResult SearchResult
}

type Post struct {
	gorm.Model

	content string

	UserID uint
	User   User `gorm:"foreignKey:UserID"`

	SuggestionID uint
}

type SearchResult struct {
	gorm.Model

	SearchKeyword string

	Phrases []string `gorm:"type:text[]"`

	UserID uint
	User   User `gorm:"foreignKey:UserID"`

	SuggestionID uint
}

type CreateSuggestionDto struct {
	Topic          string `json:"topic" binding:"required"`
	SearchKeywords string `json:"searchKeywords" binding:"required"`
	Search         bool   `json:"search" binding:"required"`
}

type CreatePostDto struct {
	Content string `json:"content" binding:"required"`
}

type CreateSearchResultDto struct {
	SearchKeyword string   `json:"searchKeyword" binding:"required"`
	Phrases       []string `json:"phrases" binding:"required"`
}
