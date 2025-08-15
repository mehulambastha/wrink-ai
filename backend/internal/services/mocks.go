package services

import (
	"context"
	"fmt"
	"time"
)

type MockSearchService struct{}

func (m *MockSearchService) Search(ctx context.Context, keywords string, topic string) (string, error) {
	fmt.Printf("MOCK SEARCH: Searching internet for keywords: %s on topic: %s", keywords, topic)

	time.Sleep(2 * time.Second)

	searchResults := fmt.Sprintf("These are the scraped search results for %s.", keywords)

	return searchResults, nil
}

type MockLLMService struct{}

func (m *MockLLMService) GeneratePost(ctx context.Context, searchResults string) (string, error) {
	fmt.Printf("MOCK LLM: Gnerating content for the search results: %s", searchResults)

	time.Sleep(time.Second * 2)

	generatedPost := "## Mock Title : This is very good \n #Big Heading ### really small one. And over\n and out!"

	return generatedPost, nil
}

type MockLinkedinService struct{}

func (m *MockLinkedinService) PublishPost(ctx context.Context, content string) (string, error) {
	fmt.Printf("MOCK LINKEDIN: Posting to Linkedin with content: %s", content)

	time.Sleep(time.Second * 2)

	postURL := "https://www.likedin.com"

	return postURL, nil
}
