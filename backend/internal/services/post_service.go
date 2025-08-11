package services

import (
	"github.com/mehulambastha/wrink-ai/backend/internal/database"
	"github.com/mehulambastha/wrink-ai/backend/internal/models"
)

func CreateSuggestion(input models.CreateSuggestionDto) (*models.SuggestionResponseDto, error) {

	newSuggestion := models.Suggestion{
		Topic:          input.Topic,
		SearchKeywords: input.SearchKeywords,
		Search:         input.Search,
		UserID:         input.UserID,
	}

	result := database.DB.Create(&newSuggestion)

	if result.Error != nil {
		return nil, result.Error
	}

	var suggestionCreatedResponse models.SuggestionResponseDto
	suggestionCreatedResponse.ID = newSuggestion.ID
	suggestionCreatedResponse.Topic = newSuggestion.Topic
	suggestionCreatedResponse.Search = newSuggestion.Search
	suggestionCreatedResponse.SearchKeywords = newSuggestion.SearchKeywords
	suggestionCreatedResponse.UserID = newSuggestion.UserID
	suggestionCreatedResponse.PostID = newSuggestion.Post.ID
	suggestionCreatedResponse.SearchResultID = newSuggestion.SearchResult.ID

	return &suggestionCreatedResponse, nil
}
