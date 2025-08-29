package pipeline

import (
	"context"
	"log"
	"time"

	"github.com/mehulambastha/wrink-ai/backend/internal/models"
	"gorm.io/gorm"
)

type ContentPipeline struct {
	db              *gorm.DB
	searchService   SearchService
	llmService      LLMService
	linkedinService LinkedinService
	steps           []PipelineStep
}

func NewContentPipeline(db *gorm.DB, searchSvc SearchService, llmSvc LLMService, linkedinSvc LinkedinService) *ContentPipeline {
	pipeline := &ContentPipeline{
		db:              db,
		searchService:   searchSvc,
		llmService:      llmSvc,
		linkedinService: linkedinSvc,
	}

	pipeline.steps = []PipelineStep{
		NewSearchInternetStep(searchSvc),
		NewGenerateContentStep(llmSvc),
		NewPublishToLinkedinStep(linkedinSvc),
	}

	return pipeline
}

func (p *ContentPipeline) execute(ctx context.Context, workflowID string) {
	log.Printf("Executing pipieline for workflow %s", workflowID)

	// 1. will create a new instance for this, therefore fetching the model first
	var instance models.WorkflowInstance

	if err := p.db.Preload("Suggestion").First(&instance, "id = ?", workflowID).Error; err != nil {
		log.Fatalf("[FATAl]: Could not fetch workflow instance %s: %v", workflowID, err)
		return
	}

	// 2. TODO: stateful step execution, by looping through p.steps and calling their execute methods, and updating the instance status

	for _, step := range p.steps {
		stepRecord := models.WorkflowStep{
			WorkflowID: instance.ID,
			StepName:   step.Name(),
			Status:     "running",
			StartedAt:  &[]time.Time{time.Now()}[0],
		}

		p.db.Create(&stepRecord)

		instance.CurrentStep = step.Name()
		instance.Status = "running"
		p.db.Save(&instance)

		outputData, err := step.Execute(ctx, &instance)

		if err != nil {
			instance.Status = "failed"
			instance.ErrorMessage = err.Error()
			p.db.Save(&instance)

			stepRecord.Status = "failed"
			stepRecord.ErrorMessage = err.Error()
			stepRecord.CompletedAt = &[]time.Time{time.Now()}[0]
			p.db.Save(&stepRecord)

			log.Printf("PIPELINE FAILED: Workflow %s failed at step %s", workflowID, step.Name())
			return
		}

		stepRecord.Status = "completed"
		stepRecord.OutputData = outputData
		stepRecord.CompletedAt = &[]time.Time{time.Now()}[0]
		p.db.Save(&stepRecord)

	}

	instance.Status = "completed"
	instance.CurrentStep = ""
	p.db.Save(&instance)

	log.Printf("PIPELINE COMPLETED: Workflow %s completed successfully", workflowID)
}

func (p *ContentPipeline) ExecuteAsync(workflowID string) {
	log.Printf("Starting pipeline for workflowID: %s in the background", workflowID)

	go p.execute(context.Background(), workflowID)
}
