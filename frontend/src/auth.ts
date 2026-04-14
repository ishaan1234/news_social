export interface AuthUser {
  uid: string;
  email: string;
  display_name?: string;
  email_verified?: boolean;
}

export interface AuthApiResponse {
  success: boolean;
  message?: string;
  user?: AuthUser;
  id_token?: string;
  refresh_token?: string;
  expires_in?: string;
}

export interface AuthSession {
  user: AuthUser | null;
  idToken: string;
  refreshToken: string;
  expiresIn: string;
  message: string;
  authenticatedAt: string;
}

const authStorageKey = 'news-social-auth-session';

export const readStoredAuthSession = (): AuthSession | null => {
  if (typeof window === 'undefined') {
    return null;
  }

  const rawSession = window.localStorage.getItem(authStorageKey);
  if (!rawSession) {
    return null;
  }

  try {
    const parsedSession = JSON.parse(rawSession) as AuthSession;
    if (!parsedSession || typeof parsedSession !== 'object') {
      return null;
    }

    return parsedSession;
  } catch (_error) {
    return null;
  }
};

export const storeAuthSession = (session: AuthSession) => {
  if (typeof window === 'undefined') {
    return;
  }

  window.localStorage.setItem(authStorageKey, JSON.stringify(session));
};

export const clearStoredAuthSession = () => {
  if (typeof window === 'undefined') {
    return;
  }

  window.localStorage.removeItem(authStorageKey);
};

export const isVerifiedAuthSession = (session: AuthSession | null) =>
  Boolean(session?.user?.email_verified);

export const buildAuthSession = (response: AuthApiResponse): AuthSession => ({
  user: response.user ?? null,
  idToken: response.id_token?.trim() || '',
  refreshToken: response.refresh_token?.trim() || '',
  expiresIn: response.expires_in?.trim() || '',
  message: response.message?.trim() || '',
  authenticatedAt: new Date().toISOString(),
});

const titleCase = (value: string) =>
  value
    .split(' ')
    .filter(Boolean)
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1).toLowerCase())
    .join(' ');

export const getSessionDisplayName = (
  session: AuthSession | null,
  fallback = 'NewsHub User'
) => {
  const explicitName = session?.user?.display_name?.trim();
  if (explicitName) {
    return explicitName;
  }

  const emailLocalPart = session?.user?.email?.split('@')[0]?.trim();
  if (!emailLocalPart) {
    return fallback;
  }

  return titleCase(emailLocalPart.replace(/[._-]+/g, ' ')) || fallback;
};

export const getSessionHandle = (
  session: AuthSession | null,
  fallback = '@newshub'
) => {
  const emailLocalPart = session?.user?.email?.split('@')[0]?.trim();
  if (!emailLocalPart) {
    return fallback;
  }

  return `@${emailLocalPart.replace(/[^a-zA-Z0-9._-]/g, '') || 'newshub'}`;
};

export const getInitials = (value: string) => {
  const parts = value.trim().split(/\s+/).filter(Boolean);
  if (parts.length === 0) {
    return 'NH';
  }

  if (parts.length === 1) {
    return parts[0].slice(0, 2).toUpperCase();
  }

  return `${parts[0][0]}${parts[1][0]}`.toUpperCase();
};
