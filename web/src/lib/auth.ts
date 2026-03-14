const TOKEN_KEY = 'daybrief_token';
const OWNER_KEY = 'daybrief_owner';
const REPO_KEY = 'daybrief_repo';

export function setAuth(token: string, owner: string, repo: string): void {
  localStorage.setItem(TOKEN_KEY, token);
  localStorage.setItem(OWNER_KEY, owner);
  localStorage.setItem(REPO_KEY, repo);
}

export function getToken(): string | null {
  if (typeof window === 'undefined') return null;
  return localStorage.getItem(TOKEN_KEY);
}

export function getOwner(): string | null {
  if (typeof window === 'undefined') return null;
  return localStorage.getItem(OWNER_KEY);
}

export function getRepo(): string | null {
  if (typeof window === 'undefined') return null;
  return localStorage.getItem(REPO_KEY);
}

export function clearAuth(): void {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(OWNER_KEY);
  localStorage.removeItem(REPO_KEY);
}

export function isAuthenticated(): boolean {
  return !!getToken() && !!getOwner() && !!getRepo();
}
