package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type WorkflowInstance struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;"`
	UserID       uint
	Status       string `gorm:"type:varchar(20);not null"`
	CurrentStep  string `gorm:"ype:varchar(20)"`
	ErrorMessage string
	RetryCount   int `gorm:"default:0`

	SuggestionID uint
	Suggestion   Suggestion `gorm:"foreignKey:SuggestionID"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Steps []WorkflowStep `gorm:"foreignKey:WorkflowID"`
}

type WorkflowStep struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;"`
	WorkflowID   uuid.UUID `gorm:"type:uuid;not null"`
	StepName     string    `gorm:"type:varchar(50);not null"`
	Status       string    `gorm:"type:varchar(20);not null"`
	InputData    datatypes.JSON
	OutputData   datatypes.JSON
	StartedAt    *time.Time
	CompletedAt  *time.Time
	ErrorMessage string
	RetryCount   int `gorm:"default:0"`
}

func (instance *WorkflowInstance) BeforeCreate(tx *gorm.DB) (err error) {
	instance.ID = uuid.New()
	return
}

func (step *WorkflowStep) BeforeCreate(tx *gorm.DB) (err error) {
	step.ID = uuid.New()
	return
}
