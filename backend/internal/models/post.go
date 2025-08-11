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
	Search         bool   `json:"search"`
	UserID         uint   `json:"-"`
}

type CreatePostDto struct {
	Content      string `json:"content" binding:"required"`
	UserID       uint   `json:"userId" binding:"required"`
	SuggestionID uint   `json:"suggestionId" binding:"required"`
}

type CreateSearchResultDto struct {
	SearchKeyword string   `json:"searchKeyword" binding:"required"`
	Phrases       []string `json:"phrases" binding:"required"`
}

type SuggestionResponseDto struct {
	ID             uint   `json:"id"`
	Search         bool   `json:"search"`
	SearchKeywords string `json:"searchKeywords"`
	Topic          string `json:"topic"`
	UserID         uint   `json:"userId"`
	PostID         uint   `json:"postId"`
	SearchResultID uint   `json:"searchResultId"`
}
