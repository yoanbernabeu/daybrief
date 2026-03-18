package openai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/yoanbernabeu/daybrief/internal/gemini"
	"github.com/yoanbernabeu/daybrief/internal/sources"
)

func (c *Client) SummarizeSource(ctx context.Context, item sources.SourceItem) (*gemini.SourceSummary, error) {
	c.logger.Debug("summarizing source with openai", "title", item.Title, "type", item.SourceType)

	prompt := c.buildSummarizePrompt(item)
	text, err := c.completeJSON(ctx, prompt)
	if err != nil {
		c.logger.Warn("failed to summarize source", "title", item.Title, "error", err)
		return nil, err
	}

	var summary gemini.SourceSummary
	if err := json.Unmarshal([]byte(text), &summary); err != nil {
		return nil, fmt.Errorf("parsing summary JSON: %w", err)
	}

	summary.SourceType = item.SourceType
	summary.SourceURL = item.URL
	summary.SourceName = item.SourceName
	summary.ThumbnailURL = item.ThumbnailURL

	return &summary, nil
}

func (c *Client) buildSummarizePrompt(item sources.SourceItem) string {
	return fmt.Sprintf(`You are a content analyst. Produce a JSON object with this exact shape:
{
  "title": string,
  "summary": string,
  "key_points": string[],
  "source_type": string,
  "source_url": string,
  "source_name": string,
  "thumbnail_url": string
}

Content metadata:
- title: %q
- source_type: %q
- source_name: %q
- source_url: %q
- thumbnail_url: %q
- audio_url: %q

Instructions:
- Write in %s.
- Provide a concise rewritten title.
- Provide a 2-3 sentence summary.
- Provide 3-5 key points.
- If full content is inaccessible from the provided metadata/URL, still provide the best possible summary and keep it factual.
- Return JSON only.`, item.Title, item.SourceType, item.SourceName, item.URL, item.ThumbnailURL, item.AudioURL, c.language)
}
