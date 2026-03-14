package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/genai"
)

var newsletterSchema = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"subject":   map[string]any{"type": "string"},
		"editorial": map[string]any{"type": "string"},
		"highlights": map[string]any{
			"type": "array",
			"items": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"title":         map[string]any{"type": "string"},
					"source_name":   map[string]any{"type": "string"},
					"source_url":    map[string]any{"type": "string"},
					"thumbnail_url": map[string]any{"type": "string"},
					"analysis":      map[string]any{"type": "string"},
				},
				"required": []string{"title", "source_name", "source_url", "analysis"},
			},
		},
		"resources": map[string]any{
			"type": "array",
			"items": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"title":       map[string]any{"type": "string"},
					"source_name": map[string]any{"type": "string"},
					"source_url":  map[string]any{"type": "string"},
					"summary":     map[string]any{"type": "string"},
				},
				"required": []string{"title", "source_name", "source_url", "summary"},
			},
		},
	},
	"required": []string{"subject", "editorial", "highlights", "resources"},
}

func (c *Client) SynthesizeNewsletter(ctx context.Context, summaries []SourceSummary) (*Newsletter, error) {
	c.logger.Info("synthesizing newsletter", "summaries", len(summaries))

	summariesJSON, err := json.MarshalIndent(summaries, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshaling summaries: %w", err)
	}

	prompt := fmt.Sprintf(`You are a senior tech newsletter editor. Based on the following source summaries, create a newsletter.

Source summaries:
%s

Instructions:

EDITORIAL (field "editorial"):
Write a substantial editorial analysis (4-6 paragraphs). This is the centerpiece of the newsletter.
- Paragraph 1: Hook the reader with the most striking trend or insight emerging from the sources
- Paragraphs 2-4: Develop a real analysis — identify connections between sources, emerging trends, what this means for practitioners. Be opinionated and insightful, not just descriptive.
- Final paragraph: Wrap up with a forward-looking perspective or actionable takeaway
- Use "\n\n" to separate paragraphs in the JSON string
- Do NOT use markdown formatting (no **, no #, no bullet points). Write in plain prose.
- Editorial tone: %s

HIGHLIGHTS (field "highlights"):
- Select the top %d most interesting/relevant items
- For each, write a detailed "analysis" (3-5 sentences): explain why it matters, what the implications are, and what the reader should take away. Go beyond summarizing — add context and perspective.

RESOURCES (field "resources"):
- List all remaining items with a one-line summary each

SUBJECT:
- Generate a catchy, specific subject line that reflects the editorial angle (not generic)

Produce your output in %s.`, string(summariesJSON), c.editorialPrompt, c.maxHighlights, c.language)

	contents := []*genai.Content{
		{
			Role:  "user",
			Parts: []*genai.Part{{Text: prompt}},
		},
	}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType:   "application/json",
		ResponseJsonSchema: newsletterSchema,
	}

	newsletter, err := withRetry(func() (*Newsletter, error) {
		resp, err := c.client.Models.GenerateContent(ctx, c.model, contents, config)
		if err != nil {
			return nil, fmt.Errorf("gemini API call: %w", err)
		}

		if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
			return nil, fmt.Errorf("empty response from Gemini")
		}

		text := resp.Candidates[0].Content.Parts[0].Text

		var nl Newsletter
		if err := json.Unmarshal([]byte(text), &nl); err != nil {
			return nil, fmt.Errorf("parsing newsletter JSON: %w", err)
		}

		return &nl, nil
	}, c.logger)

	if err != nil {
		return nil, err
	}

	newsletter.GeneratedAt = time.Now().UTC().Format(time.RFC3339)

	return newsletter, nil
}
