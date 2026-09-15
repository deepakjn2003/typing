package content

import (
	"context"
	"fmt"
)

type ExtractedContent struct {
	Title      string
	Body       string
	SourceType string
	Source     string
}

type ContentExtractor interface {
	Extract(ctx context.Context, source string) (ExtractedContent, error)
}

type TextNormalizer interface {
	Normalize(text string) string
}

type Summarizer interface {
	Summarize(ctx context.Context, text string, opts SummaryOptions) (string, error)
}

type SummaryOptions struct {
	MinWords       int
	MaxWords       int
	WeakCharacters []string
}

func DefaultSummaryOptions() SummaryOptions {
	return SummaryOptions{
		MinWords: 150,
		MaxWords: 250,
	}
}

type Pipeline struct {
	extractor  ContentExtractor
	normalizer TextNormalizer
	summarizer Summarizer
}

func NewPipeline(extractor ContentExtractor, normalizer TextNormalizer, summarizer Summarizer) *Pipeline {
	return &Pipeline{
		extractor:  extractor,
		normalizer: normalizer,
		summarizer: summarizer,
	}
}

type ProcessResult struct {
	Title      string
	Passage    string
	Source     string
	SourceType string
}

func (p *Pipeline) Process(ctx context.Context, source, sourceType string, opts SummaryOptions) (ProcessResult, error) {
	extracted, err := p.extractor.Extract(ctx, source)
	if err != nil {
		return ProcessResult{}, fmt.Errorf("extraction failed: %w", err)
	}
	if extracted.Body == "" {
		return ProcessResult{}, fmt.Errorf("no content extracted from source")
	}

	normalized := p.normalizer.Normalize(extracted.Body)
	if normalized == "" {
		return ProcessResult{}, fmt.Errorf("content empty after normalization")
	}

	if opts.MinWords <= 0 {
		opts = DefaultSummaryOptions()
	}
	passage, err := p.summarizer.Summarize(ctx, normalized, opts)
	if err != nil {
		return ProcessResult{}, fmt.Errorf("summarization failed: %w", err)
	}

	title := extracted.Title
	if title == "" {
		runes := []rune(passage)
		if len(runes) > 50 {
			title = string(runes[:50]) + "..."
		} else {
			title = passage
		}
	}

	return ProcessResult{
		Title:      title,
		Passage:    passage,
		Source:     source,
		SourceType: sourceType,
	}, nil
}
