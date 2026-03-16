package ai

import (
	"context"

	"github.com/yoanbernabeu/daybrief/internal/gemini"
	"github.com/yoanbernabeu/daybrief/internal/sources"
)

type Provider interface {
	SummarizeSource(ctx context.Context, item sources.SourceItem) (*gemini.SourceSummary, error)
	SynthesizeNewsletter(ctx context.Context, summaries []gemini.SourceSummary) (*gemini.Newsletter, error)
}
