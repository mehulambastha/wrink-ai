package services

import (
	"log"

	"github.com/mehulambastha/wrink-ai/backend/internal/models"
	"github.com/mehulambastha/wrink-ai/backend/internal/pipeline"
	"gorm.io/gorm"
)

type SuggestionService struct {
	db       *gorm.DB
	pipeline *pipeline.ContentPipeline
}

func NewSuggestionService(d *gorm.DB, p *pipeline.ContentPipeline) *SuggestionService {
	return &SuggestionService{
		db:       d,
		pipeline: p,
	}
}

func (s *SuggestionService) CreateSuggestion(input models.CreateSuggestionDto) (*models.SuggestionResponseDto, error) {

	newSuggestion := models.Suggestion{
		Topic:          input.Topic,
		SearchKeywords: input.SearchKeywords,
		Search:         input.Search,
		UserID:         input.UserID,
	}

	var workflow *models.WorkflowInstance

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Created suggestion first
		if err := tx.Create(&newSuggestion).Error; err != nil {
			log.Fatalf("Failed to create the suggestion")
			return err
		}

		// Now creating the workflow instance
		workflow = &models.WorkflowInstance{
			UserID:       input.UserID,
			SuggestionID: newSuggestion.ID,
			Status:       "pending",
		}

		if err := tx.Create(&workflow).Error; err != nil {
			log.Fatalf("Failed to create workflow instance for this suggestion")
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Failed to trigger suggestion save and/or pipeline save")
		return nil, err
	}

	s.pipeline.ExecuteAsync(workflow.ID.String())

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
