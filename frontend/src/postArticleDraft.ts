export interface LinkedArticleDraft {
  url: string;
  title: string;
  source: string;
  summary: string;
  image_url?: string;
  published_at?: string;
}

export const postArticleDraftStorageKey = 'news-social-post-article-draft';

export const savePostArticleDraft = (article: LinkedArticleDraft) => {
  if (typeof window === 'undefined') {
    return;
  }

  window.sessionStorage.setItem(
    postArticleDraftStorageKey,
    JSON.stringify(article)
  );
};

export const readPostArticleDraft = (): LinkedArticleDraft | null => {
  if (typeof window === 'undefined') {
    return null;
  }

  const rawDraft = window.sessionStorage.getItem(postArticleDraftStorageKey);
  if (!rawDraft) {
    return null;
  }

  try {
    const parsedDraft = JSON.parse(rawDraft) as LinkedArticleDraft;
    if (!parsedDraft?.url || !parsedDraft?.title) {
      return null;
    }
    return parsedDraft;
  } catch (_error) {
    return null;
  }
};

export const clearPostArticleDraft = () => {
  if (typeof window === 'undefined') {
    return;
  }

  window.sessionStorage.removeItem(postArticleDraftStorageKey);
};
