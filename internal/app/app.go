package app

import (
	"context"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/deepakjn2003/typing/internal/config"
	"github.com/deepakjn2003/typing/internal/content"
	"github.com/deepakjn2003/typing/internal/storage"
	"github.com/deepakjn2003/typing/internal/summarization"
	"github.com/deepakjn2003/typing/internal/tui"
)

func Run() error {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	db, err := storage.OpenDatabase(cfg.DBPath)
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}
	defer db.Close()

	repo := storage.NewSQLiteSessionRepository(db)
	normalizer := &content.DefaultNormalizer{}

	var pipeline *content.Pipeline
	if cfg.HasGeminiKey() {
		summarizer, err := summarization.NewGeminiSummarizer(ctx, cfg.GeminiAPIKey, cfg.GeminiModel)
		if err != nil {
			return fmt.Errorf("initializing summarizer: %w", err)
		}
		pipeline = content.NewPipeline(&autoExtractor{}, normalizer, summarizer)
	} else {
		pipeline = content.NewPipeline(&autoExtractor{}, normalizer, &passthroughSummarizer{})
	}

	model := tui.NewModel(pipeline, repo, cfg)
	p := tea.NewProgram(model)
	_, err = p.Run()
	return err
}

type autoExtractor struct {
	paste content.PasteExtractor
	url   *content.URLExtractor
}

func (e *autoExtractor) Extract(ctx context.Context, source string) (content.ExtractedContent, error) {
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		if e.url == nil {
			e.url = content.NewURLExtractor()
		}
		return e.url.Extract(ctx, source)
	}
	return e.paste.Extract(ctx, source)
}

type passthroughSummarizer struct{}

func (s *passthroughSummarizer) Summarize(_ context.Context, text string, _ content.SummaryOptions) (string, error) {
	return text, nil
}
