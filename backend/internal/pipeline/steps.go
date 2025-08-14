package pipeline

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/mehulambastha/wrink-ai/backend/internal/database"
	"github.com/mehulambastha/wrink-ai/backend/internal/models"
	"gorm.io/datatypes"
)

type SearchInternetStep struct {
	searchSvc SearchService
}

func NewSearchInternetStep(svc SearchService) *SearchInternetStep {
	return &SearchInternetStep{searchSvc: svc}
}

func (s *SearchInternetStep) Name() string {
	return "search_internet"
}

func (s *SearchInternetStep) Execute(ctx context.Context, instance *models.WorkflowInstance) (datatypes.JSON, error) {
	log.Printf("EXECUTING STEP: %s for workflwo %s", s.Name(), instance.ID)

	keywords := instance.Suggestion.SearchKeywords
	topic := instance.Suggestion.Topic

	seachResults, err := s.searchSvc.Search(ctx, keywords, topic)

	if err != nil {
		return nil, err
	}

	outputData := map[string]string{
		"search_results_text": seachResults,
	}

	outputJson, err := json.Marshal(outputData)

	if err != nil {
		return nil, err
	}

	log.Printf("Step %s completed for workflow %s", s.Name(), instance.ID)

	return outputJson, nil
}

// Content generation, this is MOCK right now.
type GenerateContentStep struct {
	llmSvc LLMService
}

func NewGenerateContentStep(svc LLMService) *GenerateContentStep {
	return &GenerateContentStep{llmSvc: svc}
}

func (g *GenerateContentStep) Name() string {
	return "generate_content"
}

func (g *GenerateContentStep) Execute(ctx context.Context, instance *models.WorkflowInstance) (datatypes.JSON, error) {
	log.Printf("EXECUTING STEP: %s for workflow %s", g.Name(), instance.ID)

	var previousStep models.WorkflowStep
	err := database.DB.Where("workflow_id = ? AND step_name = ?", instance.ID, "search_internet").First(&previousStep).Error

	if err != nil {
		return nil, errors.New("could not find previous step 'search_internet' to get input from")
	}

	var searchOutput map[string]string

	if err := json.Unmarshal(previousStep.OutputData, &searchOutput); err != nil {
		return nil, err
	}

	searchResultText := searchOutput["search_results_text"]

	// FINALLY CALL LLM LLMService
	generatedContent, err := g.llmSvc.GeneratePost(ctx, searchResultText)

	if err != nil {
		return nil, err
	}

	// return some output for this
	outputData := map[string]string{"generated_content": generatedContent}
	outputJson, err := json.Marshal(outputData)

	if err != nil {
		return nil, err
	}

	log.Printf("Step %s completed for wokflow %s", g.Name(), instance.ID)

	return outputJson, nil
}

// THe PublishtOlinkeding Step

type PublishToLinkedinStep struct {
	linkedinSvc LinkedinService
}

func NewPublishToLinkedinStep(svc LinkedinService) *PublishToLinkedinStep {
	return &PublishToLinkedinStep{linkedinSvc: svc}
}

func (p *PublishToLinkedinStep) Name() string {
	return "publish_to_linkedin"
}

func (p *PublishToLinkedinStep) Execute(ctx context.Context, instance *models.WorkflowInstance) (datatypes.JSON, error) {
	log.Printf("EXECUTING STEP: %s for workflow %s", p.Name(), instance.ID)

	var previousStep models.WorkflowStep
	err := database.DB.Where("workflow_id = ? AND step_name = ?", instance.ID, "generate_content").First(&previousStep).Error

	if err != nil {
		return nil, errors.New("could not find previous step 'generate_content' to get input from")
	}

	var contentOutput map[string]string
	if err := json.Unmarshal(previousStep.OutputData, &contentOutput); err != nil {
		return nil, err
	}

	contentToPublish := contentOutput["generated_content"]

	postURL, err := p.linkedinSvc.PublishPost(ctx, contentToPublish)

	if err != nil {
		return nil, err
	}

	outputData := map[string]string{"linkedin_post_url": postURL}

	outputJson, err := json.Marshal(outputData)

	if err != nil {
		return nil, err
	}

	log.Printf("Step %s completed for workflow %s", p.Name(), instance.ID)

	return outputJson, nil
}
