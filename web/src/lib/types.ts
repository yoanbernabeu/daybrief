// Mirror of Go structs from internal/config/config.go

export interface RSSSource {
  url: string;
  name: string;
}

export interface YouTubeSource {
  channel_id: string;
  name: string;
}

export interface PodcastSource {
  url: string;
  name: string;
}

export interface Sources {
  rss: RSSSource[];
  youtube: YouTubeSource[];
  podcasts: PodcastSource[];
}

export interface AIConfig {
  provider: "gemini" | "openai";
  model: string;
}

export interface NewsletterConfig {
  language: string;
  max_highlights: number;
  editorial_prompt: string;
  default_lookback: string;
}

export interface MailConfig {
  subject_prefix: string;
}

export interface DayBriefConfig {
  ai: AIConfig;
  newsletter: NewsletterConfig;
  mail: MailConfig;
  sources: Sources;
}

// Mirror of Go structs from internal/gemini/types.go

export interface Highlight {
  title: string;
  source_name: string;
  source_url: string;
  thumbnail_url?: string;
  analysis: string;
}

export interface Resource {
  title: string;
  source_name: string;
  source_url: string;
  summary: string;
}

export interface Newsletter {
  generated_at: string;
  subject: string;
  editorial: string;
  highlights: Highlight[];
  resources: Resource[];
}
