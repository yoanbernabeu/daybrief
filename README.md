# DayBrief

An open-source Go CLI tool that aggregates content from RSS feeds, YouTube channels, and podcasts, uses the Gemini API to summarize and analyze each source, and sends an automated HTML newsletter by email.

## Features

- **Multi-source aggregation**: RSS feeds, YouTube channels, podcasts
- **AI-powered analysis**: Two-pass Gemini integration (summarize each source, then synthesize a newsletter)
- **Automated delivery**: HTML email via SMTP
- **Incremental updates**: Only processes new content since last execution
- **CI/CD ready**: Designed to run in GitHub Actions via cron

## Installation

Download the latest binary from [GitHub Releases](https://github.com/yoanbernabeu/daybrief/releases):

```bash
curl -sL https://github.com/yoanbernabeu/daybrief/releases/latest/download/daybrief-linux-amd64 -o daybrief
chmod +x daybrief
```

Or build from source:

```bash
git clone https://github.com/yoanbernabeu/daybrief.git
cd daybrief
make build
```

## Configuration

### `config.yaml`

```yaml
gemini:
  model: "gemini-3-flash-preview"

newsletter:
  language: "fr"
  max_highlights: 5
  editorial_prompt: "A casual, tech-savvy tone with a focus on practical insights."

mail:
  subject_prefix: "[DayBrief]"

sources:
  rss:
    - url: "https://blog.golang.org/feed.atom"
      name: "Go Blog"
  youtube:
    - channel_id: "UCxxxx"
      name: "My Channel"
  podcasts:
    - url: "https://example.com/podcast.xml"
      name: "My Podcast"
```

### Environment Variables

Create a `.env` file (see `.env.example`):

```env
GEMINI_API_KEY=your-api-key
YOUTUBE_API_KEY=your-youtube-api-key
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=user
SMTP_PASSWORD=pass
MAIL_FROM_NAME=DayBrief
MAIL_FROM_EMAIL=newsletter@example.com
DAYBRIEF_RECIPIENTS=user1@example.com,user2@example.com
```

## Usage

### Run the full pipeline

```bash
daybrief run --config config.yaml
```

### Preview in browser

```bash
daybrief preview --config config.yaml
```

### Check source health

```bash
daybrief sources --config config.yaml
```

## GitHub Actions

DayBrief is available as a GitHub Action. Add a workflow to your repository:

```yaml
name: DayBrief Newsletter

on:
  workflow_dispatch:
  schedule:
    - cron: "0 7 * * 1" # Every Monday at 7:00 UTC

jobs:
  newsletter:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: yoanbernabeu/daybrief@v1
        with:
          config: config.yaml
        env:
          GEMINI_API_KEY: ${{ secrets.GEMINI_API_KEY }}
          YOUTUBE_API_KEY: ${{ secrets.YOUTUBE_API_KEY }}
          SMTP_HOST: ${{ secrets.SMTP_HOST }}
          SMTP_PORT: ${{ secrets.SMTP_PORT }}
          SMTP_USERNAME: ${{ secrets.SMTP_USERNAME }}
          SMTP_PASSWORD: ${{ secrets.SMTP_PASSWORD }}
          MAIL_FROM_NAME: ${{ secrets.MAIL_FROM_NAME }}
          MAIL_FROM_EMAIL: ${{ secrets.MAIL_FROM_EMAIL }}
          DAYBRIEF_RECIPIENTS: ${{ secrets.DAYBRIEF_RECIPIENTS }}
```

Configure the required secrets in your repository settings under **Settings > Secrets and variables > Actions**.

## License

MIT - see [LICENSE](LICENSE) for details.
