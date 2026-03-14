import { clearAuth, getToken } from './auth';

const API_BASE = 'https://api.github.com';

export interface GitHubFile {
  name: string;
  path: string;
  sha: string;
  size: number;
  type: 'file' | 'dir';
  content?: string;
  encoding?: string;
}

export interface GitHubError {
  message: string;
  status: number;
  isRateLimit?: boolean;
  isAuth?: boolean;
}

function getHeaders(token?: string | null): HeadersInit {
  const headers: HeadersInit = {
    Accept: 'application/vnd.github.v3+json',
  };
  const t = token ?? getToken();
  if (t) {
    headers['Authorization'] = `Bearer ${t}`;
  }
  return headers;
}

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const error: GitHubError = {
      message: 'An error occurred',
      status: response.status,
    };

    if (response.status === 401) {
      error.message = 'Invalid token. Please log in again.';
      error.isAuth = true;
      clearAuth();
    } else if (response.status === 403) {
      const remaining = response.headers.get('X-RateLimit-Remaining');
      if (remaining === '0') {
        const resetTime = response.headers.get('X-RateLimit-Reset');
        const resetDate = resetTime ? new Date(parseInt(resetTime) * 1000) : new Date();
        error.message = `Rate limit reached. Try again after ${resetDate.toLocaleTimeString()}.`;
        error.isRateLimit = true;
      } else {
        error.message = 'Access denied. Check your token permissions.';
      }
    } else if (response.status === 404) {
      error.message = 'Repository or file not found.';
    }

    throw error;
  }

  if (response.status === 204) {
    return {} as T;
  }

  return response.json();
}

export async function validateAccess(
  owner: string,
  repo: string,
  token: string,
): Promise<{ valid: boolean; push: boolean }> {
  const response = await fetch(`${API_BASE}/repos/${owner}/${repo}`, {
    headers: getHeaders(token),
  });

  if (!response.ok) {
    return { valid: false, push: false };
  }

  const data = await response.json();
  return {
    valid: true,
    push: data.permissions?.push ?? false,
  };
}

export async function getFileContent(
  owner: string,
  repo: string,
  path: string,
  token?: string | null,
): Promise<{ content: string; sha: string }> {
  const response = await fetch(`${API_BASE}/repos/${owner}/${repo}/contents/${path}`, {
    headers: getHeaders(token),
  });

  const file = await handleResponse<GitHubFile>(response);

  if (!file.content || file.encoding !== 'base64') {
    throw { message: 'Cannot read file content', status: 400 } as GitHubError;
  }

  const content = decodeURIComponent(escape(atob(file.content)));
  return { content, sha: file.sha };
}

export async function updateFile(
  owner: string,
  repo: string,
  path: string,
  content: string,
  sha: string,
  message: string,
  token: string,
): Promise<void> {
  const response = await fetch(`${API_BASE}/repos/${owner}/${repo}/contents/${path}`, {
    method: 'PUT',
    headers: { ...getHeaders(token), 'Content-Type': 'application/json' },
    body: JSON.stringify({
      message,
      content: btoa(unescape(encodeURIComponent(content))),
      sha,
      branch: 'main',
    }),
  });

  await handleResponse(response);
}

export async function listDirectory(
  owner: string,
  repo: string,
  path: string,
  token?: string | null,
): Promise<GitHubFile[]> {
  const response = await fetch(`${API_BASE}/repos/${owner}/${repo}/contents/${path}`, {
    headers: getHeaders(token),
  });

  const result = await handleResponse<GitHubFile | GitHubFile[]>(response);
  return Array.isArray(result) ? result : [result];
}

export async function getNewsletters(
  owner: string,
  repo: string,
  token?: string | null,
): Promise<{ date: string; url: string }[]> {
  const files = await listDirectory(owner, repo, 'output', token);
  return files
    .filter((f) => f.name.endsWith('.json'))
    .map((f) => ({
      date: f.name.replace('.json', ''),
      url: `${API_BASE}/repos/${owner}/${repo}/contents/output/${f.name}`,
    }))
    .sort((a, b) => b.date.localeCompare(a.date));
}

export async function fetchNewsletterJSON(
  owner: string,
  repo: string,
  filename: string,
  token?: string | null,
): Promise<string> {
  const { content } = await getFileContent(owner, repo, `output/${filename}`, token);
  return content;
}
