# DayBrief

An open-source GitHub Action that aggregates content from RSS feeds, YouTube channels, and podcasts, uses the Gemini API to summarize and analyze each source, and sends an automated HTML newsletter by email.

## Features

- **Multi-source aggregation**: RSS feeds, YouTube channels, podcasts
- **AI-powered analysis**: Two-pass Gemini integration (summarize each source, then synthesize a newsletter)
- **Automated delivery**: HTML email via SMTP
- **Incremental updates**: Only processes new content since last execution
- **Zero infrastructure**: Runs entirely in GitHub Actions, no server needed

## Quick Start

### 1. Create a new repository

Create a GitHub repository for your newsletter. This repo will hold your configuration and newsletter history.

### 2. Add `config.yaml`

Create a `config.yaml` at the root of your repository to define your sources and newsletter preferences:

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

| Option | Description |
|---|---|
| `gemini.model` | Gemini model to use (default: `gemini-3-flash-preview`) |
| `newsletter.language` | Newsletter language (default: `en`) |
| `newsletter.max_highlights` | Number of highlights in the newsletter (default: `5`) |
| `newsletter.default_lookback` | Time window for first run (default: `48h`) |
| `newsletter.editorial_prompt` | Tone and style instructions for the AI |
| `mail.subject_prefix` | Prefix added to email subjects |

### 3. Configure secrets

In your repository, go to **Settings > Secrets and variables > Actions** and add:

| Secret | Required | Description |
|---|---|---|
| `GEMINI_API_KEY` | Yes | [Google Gemini API key](https://ai.google.dev/) |
| `YOUTUBE_API_KEY` | If using YouTube sources | [YouTube Data API key](https://console.cloud.google.com/) |
| `SMTP_HOST` | Yes | SMTP server host (e.g. `smtp.gmail.com`) |
| `SMTP_PORT` | No | SMTP port (default: `587`) |
| `SMTP_USERNAME` | Yes | SMTP username |
| `SMTP_PASSWORD` | Yes | SMTP password |
| `MAIL_FROM_NAME` | No | Sender name (default: `DayBrief`) |
| `MAIL_FROM_EMAIL` | Yes | Sender email address |
| `DAYBRIEF_RECIPIENTS` | Yes | Comma-separated list of recipient emails |

### 4. Add the workflow

Create `.github/workflows/daybrief.yml` in your repository:

```yaml
name: DayBrief Newsletter

on:
  workflow_dispatch:
  schedule:
    - cron: "0 7 * * 1" # Every Monday at 7:00 UTC

permissions:
  contents: write

jobs:
  newsletter:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: yoanbernabeu/daybrief@v0.1.0
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

Adjust the cron schedule to your needs. The workflow can also be triggered manually via `workflow_dispatch`.

### 5. Run it

Go to the **Actions** tab in your repository, select "DayBrief Newsletter", and click **Run workflow** to test it. Once confirmed, the cron schedule will take care of the rest.

The action automatically commits newsletter output files to `output/` in your repository, which are used to track what content has already been processed (incremental updates).

## Web App

DayBrief includes a static web app (in `web/`) built with Astro 6, deployed on GitHub Pages:

- **Landing page** — Project presentation at [yoanbernabeu.github.io/daybrief](https://yoanbernabeu.github.io/daybrief)
- **Dashboard** — Web UI to edit `config.yaml` visually, manage sources, preview newsletters, and get your shareable URL
- **Public newsletter page** — Shareable page to browse newsletter archives (e.g. `yoanbernabeu.github.io/daybrief/owner/repo`)
- **Setup guide** — Step-by-step documentation with Gemini API setup and free email provider recommendations

```bash
cd web
npm install
npm run dev      # Dev server
npm run build    # Production build → web/dist/
```

## CLI Usage

DayBrief can also be used as a standalone CLI tool.

Download the binary from [GitHub Releases](https://github.com/yoanbernabeu/daybrief/releases):

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

Available commands:

```bash
daybrief run --config config.yaml       # Run the full newsletter pipeline
daybrief preview --config config.yaml   # Generate and preview in browser
daybrief sources --config config.yaml   # Check source accessibility
```

When running locally, create a `.env` file with the same variables as the GitHub secrets (see `.env.example`).

## License

MIT - see [LICENSE](LICENSE) for details.
