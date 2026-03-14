import jsYaml from 'js-yaml';
import type { DayBriefConfig } from './types';

export function parseConfig(yamlString: string): DayBriefConfig {
  const raw = jsYaml.load(yamlString) as Partial<DayBriefConfig> | null;

  const config: DayBriefConfig = {
    gemini: {
      model: raw?.gemini?.model || 'gemini-3-flash-preview',
    },
    newsletter: {
      language: raw?.newsletter?.language || 'en',
      max_highlights: raw?.newsletter?.max_highlights || 5,
      editorial_prompt: raw?.newsletter?.editorial_prompt || '',
      default_lookback: raw?.newsletter?.default_lookback || '48h',
    },
    mail: {
      subject_prefix: raw?.mail?.subject_prefix || '',
    },
    sources: {
      rss: raw?.sources?.rss || [],
      youtube: raw?.sources?.youtube || [],
      podcasts: raw?.sources?.podcasts || [],
    },
  };

  return config;
}

export function stringifyConfig(config: DayBriefConfig): string {
  return jsYaml.dump(config, {
    indent: 2,
    lineWidth: -1,
    noRefs: true,
    quotingType: '"',
    forceQuotes: false,
  });
}
