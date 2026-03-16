package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/yoanbernabeu/daybrief/internal/gemini"
)

func (c *Client) SynthesizeNewsletter(ctx context.Context, summaries []gemini.SourceSummary) (*gemini.Newsletter, error) {
	c.logger.Info("synthesizing newsletter with openai", "summaries", len(summaries))

	summariesJSON, err := json.MarshalIndent(summaries, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshaling summaries: %w", err)
	}

	prompt := fmt.Sprintf(`You are a senior tech newsletter editor. Based on the source summaries below, return one JSON object with this exact shape:
{
  "subject": string,
  "editorial": string,
  "highlights": [
    {
      "title": string,
      "source_name": string,
      "source_url": string,
      "thumbnail_url": string,
      "analysis": string
    }
  ],
  "resources": [
    {
      "title": string,
      "source_name": string,
      "source_url": string,
      "summary": string
    }
  ]
}

Source summaries:
%s

Instructions:
- Output language: %s.
- Editorial tone: %s.
- Editorial must be 4-6 paragraphs separated with "\n\n".
- Select the top %d items for highlights, each with 3-5 sentences of analysis.
- Put all remaining items into resources with one-line summaries.
- Generate a specific and catchy subject line.
- Return JSON only.`, string(summariesJSON), c.language, c.editorialPrompt, c.maxHighlights)

	text, err := c.completeJSON(ctx, prompt)
	if err != nil {
		return nil, err
	}

	var nl gemini.Newsletter
	if err := json.Unmarshal([]byte(text), &nl); err != nil {
		return nil, fmt.Errorf("parsing newsletter JSON: %w", err)
	}

	nl.GeneratedAt = time.Now().UTC().Format(time.RFC3339)
	return &nl, nil
}
