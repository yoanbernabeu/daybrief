package gemini

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/yoanbernabeu/daybrief/internal/sources"
	"google.golang.org/genai"
)

var sourceSummarySchema = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"title":         map[string]any{"type": "string"},
		"summary":       map[string]any{"type": "string"},
		"key_points":    map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
		"source_type":   map[string]any{"type": "string"},
		"source_url":    map[string]any{"type": "string"},
		"source_name":   map[string]any{"type": "string"},
		"thumbnail_url": map[string]any{"type": "string"},
	},
	"required": []string{"title", "summary", "key_points", "source_type", "source_url", "source_name"},
}

func (c *Client) SummarizeSource(ctx context.Context, item sources.SourceItem) (*SourceSummary, error) {
	c.logger.Debug("summarizing source", "title", item.Title, "type", item.SourceType)

	parts, tools := c.buildSummarizeRequest(item)
	prompt := c.buildSummarizePrompt(item)

	allParts := append([]*genai.Part{{Text: prompt}}, parts...)
	contents := []*genai.Content{
		{
			Role:  "user",
			Parts: allParts,
		},
	}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType:   "application/json",
		ResponseJsonSchema: sourceSummarySchema,
	}
	if len(tools) > 0 {
		config.Tools = tools
	}

	result, err := withRetry(func() (*SourceSummary, error) {
		resp, err := c.client.Models.GenerateContent(ctx, c.model, contents, config)
		if err != nil {
			return nil, fmt.Errorf("gemini API call: %w", err)
		}

		if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
			return nil, fmt.Errorf("empty response from Gemini")
		}

		text := resp.Candidates[0].Content.Parts[0].Text

		var summary SourceSummary
		if err := json.Unmarshal([]byte(text), &summary); err != nil {
			return nil, fmt.Errorf("parsing summary JSON: %w", err)
		}

		return &summary, nil
	}, c.logger)

	if err != nil {
		c.logger.Warn("failed to summarize source", "title", item.Title, "error", err)
		return nil, err
	}

	// Preserve source metadata
	result.SourceType = item.SourceType
	result.SourceURL = item.URL
	result.SourceName = item.SourceName
	result.ThumbnailURL = item.ThumbnailURL

	return result, nil
}

func (c *Client) buildSummarizePrompt(item sources.SourceItem) string {
	langDirective := fmt.Sprintf("\n\nProduce your output in %s.", c.language)

	switch item.SourceType {
	case "rss":
		return fmt.Sprintf(`You are a content analyst. Analyze the following web article and provide a structured summary.

Article title: %s
Article URL: %s
Source: %s

Read the article at the URL provided and produce:
- A concise title
- A 2-3 sentence summary capturing the main points
- 3-5 key points as bullet items%s`, item.Title, item.URL, item.SourceName, langDirective)

	case "youtube":
		return fmt.Sprintf(`You are a content analyst. Analyze the following YouTube video and provide a structured summary.

Video title: %s
Video URL: %s
Channel: %s

Watch the video and produce:
- A concise title
- A 2-3 sentence summary capturing the main points
- 3-5 key points as bullet items%s`, item.Title, item.URL, item.SourceName, langDirective)

	case "podcast":
		return fmt.Sprintf(`You are a content analyst. Analyze the following podcast episode and provide a structured summary.

Episode title: %s
Episode URL: %s
Podcast: %s

Listen to the episode and produce:
- A concise title
- A 2-3 sentence summary capturing the main points
- 3-5 key points as bullet items%s`, item.Title, item.URL, item.SourceName, langDirective)

	default:
		return fmt.Sprintf(`Analyze the following content and provide a structured summary.
Title: %s
URL: %s
Source: %s%s`, item.Title, item.URL, item.SourceName, langDirective)
	}
}

func (c *Client) buildSummarizeRequest(item sources.SourceItem) ([]*genai.Part, []*genai.Tool) {
	var parts []*genai.Part
	var tools []*genai.Tool

	switch item.SourceType {
	case "rss":
		tools = append(tools, &genai.Tool{
			URLContext: &genai.URLContext{},
		})
	case "youtube":
		parts = append(parts, &genai.Part{
			FileData: &genai.FileData{
				FileURI:  item.URL,
				MIMEType: "video/youtube",
			},
		})
	case "podcast":
		if item.AudioURL != "" {
			parts = append(parts, &genai.Part{
				FileData: &genai.FileData{
					FileURI:  item.AudioURL,
					MIMEType: "audio/mpeg",
				},
			})
		}
	}

	return parts, tools
}
