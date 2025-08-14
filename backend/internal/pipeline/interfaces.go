package pipeline

import (
	"context"
	"github.com/mehulambastha/wrink-ai/backend/internal/models"
	"gorm.io/datatypes"
)

type PipelineStep interface {
	Execute(ctx context.Context, instance *models.WorkflowInstance) (datatypes.JSON, error)
	Name() string
}

type SearchService interface {
	Search(ctx context.Context, keywords string, topic string) (string, error)
}

type LLMService interface {
	GeneratePost(ctx context.Context, searchResults string) (string, error)
}

type LinkedinService interface {
	PublishPost(ctx context.Context, content string) (string, error)
}
