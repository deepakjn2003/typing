package summarization

import (
	"context"
	"fmt"
	"strings"

	"github.com/deepakjn2003/typing/internal/content"
	"google.golang.org/genai"
)

const DefaultModel = "gemini-2.5-flash"

type GeminiSummarizer struct {
	client *genai.Client
	model  string
}

func NewGeminiSummarizer(ctx context.Context, apiKey, model string) (*GeminiSummarizer, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("Gemini API key is required (set GEMINI_API_KEY)")
	}

	if model == "" {
		model = DefaultModel
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("creating Gemini client: %w", err)
	}

	return &GeminiSummarizer{
		client: client,
		model:  model,
	}, nil
}

func (s *GeminiSummarizer) Summarize(ctx context.Context, text string, opts content.SummaryOptions) (string, error) {
	if strings.TrimSpace(text) == "" {
		return "", fmt.Errorf("cannot summarize empty text")
	}

	minWords := opts.MinWords
	maxWords := opts.MaxWords
	if minWords <= 0 {
		minWords = 150
	}
	if maxWords <= 0 {
		maxWords = 250
	}

	systemPrompt, userPrompt := BuildSummarizationPrompt(text, minWords, maxWords, opts.WeakCharacters)

	temp := float32(0.3)
	result, err := s.client.Models.GenerateContent(ctx, s.model, []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{Text: userPrompt},
			},
		},
	}, &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{
				{Text: systemPrompt},
			},
		},
		Temperature:     &temp,
		MaxOutputTokens: 1024,
	})
	if err != nil {
		return "", fmt.Errorf("Gemini API call failed: %w", err)
	}

	if result == nil || len(result.Candidates) == 0 {
		return "", fmt.Errorf("Gemini returned no candidates")
	}

	candidate := result.Candidates[0]
	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		return "", fmt.Errorf("Gemini returned empty content")
	}

	var output strings.Builder
	for _, part := range candidate.Content.Parts {
		if part.Text != "" {
			output.WriteString(part.Text)
		}
	}

	passage := strings.TrimSpace(output.String())
	if passage == "" {
		return "", fmt.Errorf("Gemini returned empty passage")
	}

	return passage, nil
}
