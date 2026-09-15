package content

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	readability "codeberg.org/readeck/go-readability/v2"
)

const (
	maxResponseSize = 5 * 1024 * 1024
	defaultTimeout  = 30 * time.Second
)

type PasteExtractor struct{}

func (e *PasteExtractor) Extract(_ context.Context, source string) (ExtractedContent, error) {
	if strings.TrimSpace(source) == "" {
		return ExtractedContent{}, fmt.Errorf("pasted content is empty")
	}

	return ExtractedContent{
		Title:      deriveTitle(source),
		Body:       source,
		SourceType: "paste",
		Source:     source,
	}, nil
}

type URLExtractor struct {
	client *http.Client
}

func NewURLExtractor() *URLExtractor {
	return &URLExtractor{
		client: &http.Client{
			Timeout: defaultTimeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 10 {
					return fmt.Errorf("too many redirects")
				}
				return nil
			},
		},
	}
}

func (e *URLExtractor) Extract(ctx context.Context, source string) (ExtractedContent, error) {
	parsed, err := url.ParseRequestURI(source)
	if err != nil {
		return ExtractedContent{}, fmt.Errorf("invalid URL: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return ExtractedContent{}, fmt.Errorf("URL must use http or https scheme, got %q", parsed.Scheme)
	}
	if parsed.Host == "" {
		return ExtractedContent{}, fmt.Errorf("URL has no host")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, source, nil)
	if err != nil {
		return ExtractedContent{}, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", "typing/1.0 (article reader)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")

	resp, err := e.client.Do(req)
	if err != nil {
		return ExtractedContent{}, fmt.Errorf("fetching URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ExtractedContent{}, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	article, err := readability.FromReader(io.LimitReader(resp.Body, maxResponseSize), parsed)
	if err != nil {
		return ExtractedContent{}, fmt.Errorf("extracting article content: %w", err)
	}

	var buf bytes.Buffer
	if err := article.RenderText(&buf); err != nil {
		return ExtractedContent{}, fmt.Errorf("rendering article text: %w", err)
	}

	body := strings.TrimSpace(buf.String())
	if body == "" {
		return ExtractedContent{}, fmt.Errorf("no article content found at URL")
	}

	title := strings.TrimSpace(article.Title())
	if title == "" {
		title = parsed.Host + parsed.Path
	}

	return ExtractedContent{
		Title:      title,
		Body:       body,
		SourceType: "url",
		Source:     source,
	}, nil
}

func deriveTitle(text string) string {
	lines := strings.SplitN(text, "\n", 2)
	firstLine := strings.TrimSpace(lines[0])
	if len(firstLine) > 0 && len(firstLine) <= 60 {
		return firstLine
	}

	runes := []rune(strings.TrimSpace(text))
	if len(runes) > 50 {
		cut := string(runes[:50])
		if lastSpace := strings.LastIndex(cut, " "); lastSpace > 20 {
			return cut[:lastSpace] + "..."
		}
		return cut + "..."
	}
	return string(runes)
}
