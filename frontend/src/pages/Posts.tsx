import React, {
  FormEvent,
  useCallback,
  useEffect,
  useMemo,
  useState,
} from 'react';
import {
  ArrowUpOnSquareIcon,
  ChatBubbleLeftRightIcon,
  HeartIcon,
  NewspaperIcon,
  UserMinusIcon,
  UserPlusIcon,
} from '@heroicons/react/24/outline';
import { HeartIcon as HeartSolidIcon } from '@heroicons/react/24/solid';
import {
  LinkedArticleDraft,
  clearPostArticleDraft,
  readPostArticleDraft,
} from '../postArticleDraft';
import {
  AuthSession,
  getSessionDisplayName,
  getSessionHandle,
  isVerifiedAuthSession,
} from '../auth';

interface NewsReference {
  id: string;
  articleId?: string;
  headline: string;
  source: string;
  category: string;
  summary: string;
  articleUrl: string;
  imageUrl?: string;
  publishedAt?: string;
}

interface OpinionPost {
  id: string;
  author: string;
  handle: string;
  userEmail?: string;
  postedAt: string;
  body: string;
  newsId: string;
  linkedNews?: NewsReference;
  likeCount: number;
  shareCount: number;
  isLiked: boolean;
  commentCount?: number;
  comments: PostComment[];
}

interface PostComment {
  id: string;
  author: string;
  body: string;
}

interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
  feed?: T;
  post?: T;
  like?: T;
  comment?: T;
  comments?: T;
  follow?: T;
  message?: string;
}

interface BackendPost {
  id: string;
  user_email: string;
  username?: string;
  display_name?: string;
  avatar_url?: string;
  caption?: string;
  created_at: string;
  like_count: number;
  comment_count: number;
  liked_by_me: boolean;
  article: BackendFeedArticle;
}

interface BackendFeedArticle {
  id: string;
  title: string;
  description?: string;
  content?: string;
  summary?: string;
  author?: string;
  source_name?: string;
  source_id?: string;
  url: string;
  image_url?: string;
  published_at?: string;
  created_at: string;
}

interface BackendPostComment {
  id: string;
  post_id: string;
  user_email: string;
  username?: string;
  display_name?: string;
  avatar_url?: string;
  content: string;
  created_at: string;
}

interface PostsProps {
  authSession?: AuthSession | null;
}

const initialPosts: OpinionPost[] = [];

const apiBaseUrl = (process.env.REACT_APP_API_BASE_URL || '').replace(
  /\/$/,
  ''
);

const getStableArticleId = (url: string) =>
  `article-${url.replace(/[^a-zA-Z0-9]+/g, '-').slice(0, 48) || 'linked'}`;

const toNewsReference = (article: LinkedArticleDraft): NewsReference => ({
  id: article.id || getStableArticleId(article.url),
  articleId: article.id,
  headline: article.title,
  source: article.source || 'Unknown source',
  category: 'News',
  summary: article.summary || 'Summary unavailable.',
  articleUrl: article.url,
  imageUrl: article.image_url,
  publishedAt: article.published_at,
});

const createLinkedNewsFromDraft = () => {
  const draft = readPostArticleDraft();
  return draft ? toNewsReference(draft) : null;
};

const authHeaders = (authSession: AuthSession | null) => {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };

  if (authSession?.idToken) {
    headers.Authorization = `Bearer ${authSession.idToken}`;
  }

  return headers;
};

const readApiData = async <T,>(
  path: string,
  payloadKey: keyof ApiResponse<T>,
  init?: RequestInit
): Promise<T> => {
  const response = await fetch(`${apiBaseUrl}${path}`, init);
  const body = (await response.json().catch(() => null)) as
    | ApiResponse<T>
    | T
    | null;

  if (!response.ok) {
    const errorMessage =
      body && typeof body === 'object' && 'error' in body
        ? String((body as ApiResponse<T>).error)
        : `request failed with status ${response.status}`;
    throw new Error(errorMessage);
  }

  if (body && typeof body === 'object' && 'success' in body) {
    const apiBody = body as ApiResponse<T>;
    if (!apiBody.success) {
      throw new Error(apiBody.error || 'request failed');
    }
    return (apiBody[payloadKey] ?? apiBody.data) as T;
  }

  return body as T;
};

const formatPostTime = (createdAt?: string) => {
  if (!createdAt) {
    return 'Just now';
  }

  const createdMs = new Date(createdAt).getTime();
  if (Number.isNaN(createdMs)) {
    return 'Recently';
  }

  const diffMinutes = Math.max(0, Math.floor((Date.now() - createdMs) / 60000));
  if (diffMinutes < 1) {
    return 'Just now';
  }
  if (diffMinutes < 60) {
    return `${diffMinutes}m ago`;
  }

  const diffHours = Math.floor(diffMinutes / 60);
  if (diffHours < 24) {
    return `${diffHours}h ago`;
  }

  return `${Math.floor(diffHours / 24)}d ago`;
};

const mapBackendPost = (post: BackendPost): OpinionPost => {
  const linkedNews = toNewsReference({
    id: post.article.id,
    url: post.article.url,
    title: post.article.title,
    source: post.article.source_name || 'Unknown source',
    summary:
      post.article.summary ||
      post.article.description ||
      'Summary unavailable.',
    image_url: post.article.image_url,
    published_at: post.article.published_at,
  });
  const authorName = post.display_name || post.username || post.user_email;
  const handle = post.username
    ? `@${post.username.replace(/^@+/, '')}`
    : `@${post.user_email.split('@')[0] || 'newshub'}`;

  return {
    id: post.id,
    author: authorName,
    handle,
    userEmail: post.user_email,
    postedAt: formatPostTime(post.created_at),
    body: post.caption || '',
    newsId: linkedNews.id,
    linkedNews,
    likeCount: post.like_count,
    shareCount: 0,
    isLiked: post.liked_by_me,
    commentCount: post.comment_count,
    comments: [],
  };
};

const mapBackendComment = (
  comment: BackendPostComment,
  fallbackAuthor = ''
): PostComment => {
  const authorName =
    comment.display_name ||
    comment.username ||
    fallbackAuthor ||
    comment.user_email;

  return {
    id: comment.id,
    author: authorName,
    body: comment.content,
  };
};

const createPostId = () =>
  `post-${Date.now()}-${Math.random().toString(36).slice(2, 7)}`;

const createCommentId = () =>
  `comment-${Date.now()}-${Math.random().toString(36).slice(2, 7)}`;

const Posts: React.FC<PostsProps> = ({ authSession = null }) => {
  const [linkedNews] = useState<NewsReference | null>(() =>
    createLinkedNewsFromDraft()
  );
  const [availableNews, setAvailableNews] = useState<NewsReference[]>(
    linkedNews ? [linkedNews] : []
  );
  const [posts, setPosts] = useState<OpinionPost[]>(initialPosts);
  const [draft, setDraft] = useState('');
  const [selectedNewsId, setSelectedNewsId] = useState(
    () => linkedNews?.id || ''
  );
  const [followedHandles, setFollowedHandles] = useState<string[]>([]);
  const [commentDrafts, setCommentDrafts] = useState<Record<string, string>>(
    {}
  );
  const hasVerifiedSession = isVerifiedAuthSession(authSession);
  const sessionEmail = authSession?.user?.email?.trim().toLowerCase() || '';
  const canUseBackend = hasVerifiedSession && Boolean(sessionEmail);
  const currentAuthorName = getSessionDisplayName(authSession, 'You');
  const currentHandle = getSessionHandle(authSession, '@you');

  const selectedNews = useMemo(
    () =>
      availableNews.find((item) => item.id === selectedNewsId) ??
      availableNews[0],
    [availableNews, selectedNewsId]
  );

  const loadNews = useCallback(async () => {
    try {
      const response = await fetch(`${apiBaseUrl}/news?q=latest`);
      const data = await response.json();
      if (data.articles) {
        const mappedNews: NewsReference[] = data.articles.map((article: any) => ({
          id: article.id || getStableArticleId(article.url),
          articleId: article.id,
          headline: article.title || 'Untitled',
          source: article.source?.name || 'Unknown source',
          category: 'News',
          summary: article.summary || article.description || 'Summary unavailable.',
          articleUrl: article.url || 'https://example.com',
          imageUrl: article.urlToImage,
          publishedAt: article.publishedAt,
        }));
        setAvailableNews((prev) => {
          const newItems = mappedNews.filter(
            (n) => !prev.some((p) => p.id === n.id)
          );
          return [...prev, ...newItems];
        });
        if (!selectedNewsId && mappedNews.length > 0) {
          setSelectedNewsId(mappedNews[0].id);
        }
      }
    } catch (error) {
      // Ignore errors for now
    }
  }, [selectedNewsId]);

  const loadPosts = useCallback(async () => {
    if (!canUseBackend) {
      setPosts(initialPosts);
      return;
    }

    try {
      const backendPosts = await readApiData<BackendPost[]>(
        `/feed?user_email=${encodeURIComponent(sessionEmail)}`,
        'feed',
        { headers: authHeaders(authSession) }
      );

      if (Array.isArray(backendPosts)) {
        const mappedPosts = backendPosts.map(mapBackendPost);
        const hydratedPosts = await Promise.all(
          mappedPosts.map(async (post) => {
            try {
              const comments = await readApiData<BackendPostComment[]>(
                `/post-comments?post_id=${encodeURIComponent(post.id)}`,
                'comments',
                { headers: authHeaders(authSession) }
              );

              return {
                ...post,
                comments: Array.isArray(comments)
                  ? comments.map((comment) => mapBackendComment(comment))
                  : post.comments,
                commentCount: Array.isArray(comments)
                  ? comments.length
                  : post.commentCount,
              };
            } catch (_error) {
              return post;
            }
          })
        );

        setPosts(hydratedPosts);
      }
    } catch (_error) {
      // Keep the local starter posts when the backend API is unavailable.
    }
  }, [authSession, canUseBackend, sessionEmail]);

  useEffect(() => {
    void loadNews();
    void loadPosts();
  }, [loadNews, loadPosts]);

  useEffect(() => {
    const hash = window.location.hash;
    const match = hash.match(/postId=([^&]+)/);
    if (match && match[1] && posts.length > 0) {
      setTimeout(() => {
        const el = document.getElementById(`post-${match[1]}`);
        if (el) {
          el.scrollIntoView({ behavior: 'smooth', block: 'center' });
          el.classList.add('ring-4', 'ring-blue-300', 'transition-all');
          setTimeout(() => el.classList.remove('ring-4', 'ring-blue-300'), 2000);
        }
      }, 500);
    }
  }, [posts.length]);

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const trimmedDraft = draft.trim();
    if (!trimmedDraft) {
      return;
    }

    if (canUseBackend && selectedNews?.articleId) {
      try {
        await readApiData<{
          id: string;
          user_email: string;
          article_id: string;
          caption?: string;
          created_at: string;
        }>('/posts', 'post', {
          method: 'POST',
          headers: authHeaders(authSession),
          body: JSON.stringify({
            user_email: sessionEmail,
            article_id: selectedNews.articleId,
            caption: trimmedDraft,
          }),
        });

        await loadPosts();
        setDraft('');
        clearPostArticleDraft();
        return;
      } catch (_error) {
        // Local fallback keeps the composer usable if backend rejects the post.
      }
    }

    const nextPost: OpinionPost = {
      id: createPostId(),
      author: currentAuthorName,
      handle: currentHandle,
      userEmail: sessionEmail || undefined,
      postedAt: 'Just now',
      body: trimmedDraft,
      newsId: selectedNews.id,
      linkedNews: selectedNews,
      likeCount: 0,
      shareCount: 0,
      isLiked: false,
      comments: [],
    };

    setPosts((previousPosts) => [nextPost, ...previousPosts]);
    setDraft('');
    clearPostArticleDraft();
  };

  const toggleLike = (postId: string) => {
    const targetPost = posts.find((post) => post.id === postId);
    if (!targetPost) {
      return;
    }

    const nextLiked = !targetPost.isLiked;

    setPosts((previousPosts) =>
      previousPosts.map((post) => {
        if (post.id !== postId) {
          return post;
        }

        return {
          ...post,
          isLiked: nextLiked,
          likeCount: Math.max(0, post.likeCount + (nextLiked ? 1 : -1)),
        };
      })
    );

    if (!canUseBackend || !targetPost.userEmail) {
      return;
    }

    void readApiData<unknown>('/post-likes', nextLiked ? 'like' : 'data', {
      method: nextLiked ? 'POST' : 'DELETE',
      headers: authHeaders(authSession),
      body: JSON.stringify({
        user_email: sessionEmail,
        post_id: targetPost.id,
      }),
    }).catch(() => {
      setPosts((previousPosts) =>
        previousPosts.map((post) =>
          post.id === postId
            ? {
                ...post,
                isLiked: targetPost.isLiked,
                likeCount: targetPost.likeCount,
              }
            : post
        )
      );
    });
  };

  const sharePost = (postId: string) => {
    setPosts((previousPosts) =>
      previousPosts.map((post) =>
        post.id === postId ? { ...post, shareCount: post.shareCount + 1 } : post
      )
    );
  };

  const toggleFollow = (post: OpinionPost) => {
    const isCurrentlyFollowing = followedHandles.includes(post.handle);

    setFollowedHandles((previousHandles) =>
      isCurrentlyFollowing
        ? previousHandles.filter(
            (currentHandle) => currentHandle !== post.handle
          )
        : [...previousHandles, post.handle]
    );

    if (!canUseBackend || !post.userEmail || post.userEmail === sessionEmail) {
      return;
    }

    void readApiData<unknown>(
      '/following',
      isCurrentlyFollowing ? 'data' : 'follow',
      {
        method: isCurrentlyFollowing ? 'DELETE' : 'POST',
        headers: authHeaders(authSession),
        body: JSON.stringify({
          follower_email: sessionEmail,
          following_email: post.userEmail,
        }),
      }
    ).catch(() => {
      setFollowedHandles((previousHandles) =>
        isCurrentlyFollowing
          ? [...previousHandles, post.handle]
          : previousHandles.filter(
              (currentHandle) => currentHandle !== post.handle
            )
      );
    });
  };

  const addComment = (postId: string) => {
    const trimmedComment = commentDrafts[postId]?.trim();
    if (!trimmedComment) {
      return;
    }

    const targetPost = posts.find((post) => post.id === postId);
    const nextComment: PostComment = {
      id: createCommentId(),
      author: currentAuthorName,
      body: trimmedComment,
    };

    setPosts((previousPosts) =>
      previousPosts.map((post) =>
        post.id === postId
          ? {
              ...post,
              comments: [...post.comments, nextComment],
              commentCount: (post.commentCount ?? post.comments.length) + 1,
            }
          : post
      )
    );

    setCommentDrafts((previousDrafts) => ({
      ...previousDrafts,
      [postId]: '',
    }));

    if (!canUseBackend || !targetPost?.userEmail) {
      return;
    }

    void readApiData<BackendPostComment>('/post-comments', 'comment', {
      method: 'POST',
      headers: authHeaders(authSession),
      body: JSON.stringify({
        post_id: targetPost.id,
        user_email: sessionEmail,
        content: trimmedComment,
      }),
    })
      .then((createdComment) => {
        const persistedComment = mapBackendComment(
          createdComment,
          currentAuthorName
        );
        setPosts((previousPosts) =>
          previousPosts.map((post) =>
            post.id === postId
              ? {
                  ...post,
                  comments: post.comments.map((comment) =>
                    comment.id === nextComment.id ? persistedComment : comment
                  ),
                }
              : post
          )
        );
      })
      .catch(() => undefined);
  };

  const deletePost = async (postId: string) => {
    if (!window.confirm("Are you sure you want to delete this post?")) return;

    if (canUseBackend) {
      try {
        await readApiData(`/posts?id=${postId}`, 'data', {
          method: 'DELETE',
          headers: authHeaders(authSession),
        });
      } catch (e) {
        console.error("Failed to delete post", e);
        return;
      }
    }

    setPosts((prev) => prev.filter(p => p.id !== postId));
  };

  return (
    <main data-cy="posts-page" className="mx-auto max-w-6xl px-4 py-6 sm:px-6">
      {/* <section className="rounded-[32px] bg-gradient-to-br from-slate-900 via-slate-800 to-blue-900 px-6 py-8 text-white shadow-sm sm:px-8">
        <div className="flex flex-col gap-6 lg:flex-row lg:items-end lg:justify-between">
          <div className="max-w-2xl">
            <div className="inline-flex items-center gap-2 rounded-full bg-white/10 px-3 py-1 text-xs font-semibold uppercase tracking-[0.2em] text-blue-100">
              <SparklesIcon className="h-4 w-4" />
              Frontend Only Opinions
            </div>
            <h1 className="mt-4 text-3xl font-bold sm:text-4xl">
              Post your take on the news
            </h1>
            <p className="mt-3 text-sm leading-7 text-slate-200 sm:text-base">
              This is a local frontend demo. Posts, reactions, and the attached
              summary cards stay in the browser and use placeholder news items.
            </p>
          </div>

        </div>
      </section> */}

      <div className="mt-6 grid gap-6 lg:grid-cols-[360px_minmax(0,1fr)]">
        <aside className="space-y-6">
          <section className="rounded-[28px] bg-white p-5 shadow-sm">
            <div className="flex items-center gap-3">
              <div className="rounded-2xl bg-blue-50 p-3 text-blue-600">
                <NewspaperIcon className="h-6 w-6" />
              </div>
              <div>
                <h2 className="text-lg font-bold text-slate-900">
                  Attach a news summary
                </h2>
                <p className="text-sm text-slate-500">
                  Pick the story your post is reacting to.
                </p>
              </div>
            </div>

            <div className="mt-5 space-y-3">
              {availableNews.map((newsItem) => {
                const isSelected = newsItem.id === selectedNewsId;

                return (
                  <button
                    key={newsItem.id}
                    type="button"
                    onClick={() => setSelectedNewsId(newsItem.id)}
                    className={`w-full rounded-3xl border px-4 py-4 text-left transition ${
                      isSelected
                        ? 'border-blue-500 bg-blue-50'
                        : 'border-slate-200 hover:border-slate-300 hover:bg-slate-50'
                    }`}
                  >
                    <p className="text-xs font-semibold uppercase tracking-[0.18em] text-slate-400">
                      {newsItem.category} | {newsItem.source}
                    </p>
                    <p className="mt-2 text-sm font-semibold leading-6 text-slate-900">
                      {newsItem.headline}
                    </p>
                  </button>
                );
              })}
            </div>
          </section>

          {selectedNews && (
            <section className="rounded-[28px] border border-blue-100 bg-blue-50 p-5">
              <p className="text-xs font-semibold uppercase tracking-[0.18em] text-blue-600">
                Attached Summary Preview
              </p>
              <h3 className="mt-3 text-lg font-bold text-slate-900">
                {selectedNews.headline}
              </h3>
              <p className="mt-2 text-xs font-medium uppercase tracking-[0.16em] text-slate-500">
                {selectedNews.source}
              </p>
              <p className="mt-3 text-sm leading-6 text-slate-600">
                {selectedNews.summary}
              </p>
            </section>
          )}
        </aside>

        <section className="space-y-6">
          <form
            onSubmit={handleSubmit}
            className="rounded-[28px] bg-white p-6 shadow-sm"
          >
            <div className="flex items-center justify-between gap-3">
              <div>
                <h2 className="text-xl font-bold text-slate-900">
                  Share an opinion
                </h2>
                <p className="mt-1 text-sm text-slate-500">
                  Draft a quick take, reaction, or viewpoint.
                </p>
              </div>
              <span className="rounded-full bg-slate-100 px-3 py-1 text-xs font-semibold uppercase tracking-[0.18em] text-slate-500">
                {draft.trim().length}/280
              </span>
            </div>

            <label className="mt-5 block">
              <span className="sr-only">Write your opinion</span>
              <textarea
                data-cy="posts-draft"
                rows={5}
                maxLength={280}
                value={draft}
                onChange={(event) => setDraft(event.target.value)}
                placeholder="Write what you think about this story..."
                className="w-full resize-none rounded-[24px] border border-slate-200 bg-slate-50 px-4 py-4 text-sm leading-7 text-slate-700 outline-none transition placeholder:text-slate-400 focus:border-blue-400 focus:bg-white"
              />
            </label>

            {selectedNews && (
              <div className="mt-4 rounded-[24px] border border-slate-200 bg-slate-50 px-4 py-4">
                <p className="text-xs font-semibold uppercase tracking-[0.16em] text-slate-400">
                  News summary that will appear below your post
                </p>
                <p className="mt-2 text-sm font-semibold text-slate-900">
                  {selectedNews.headline}
                </p>
                <p className="mt-2 text-sm leading-6 text-slate-600">
                  {selectedNews.summary}
                </p>
              </div>
            )}

            <div className="mt-5 flex items-center justify-between gap-4">
              <p className="text-xs text-slate-400">
                {canUseBackend
                  ? selectedNews?.articleId
                    ? `Posting as ${currentHandle}. This will save to the database.`
                    : 'This story needs a saved article id before it can persist.'
                  : 'Sign in with a verified account to save posts.'}
              </p>
              <button
                type="submit"
                data-cy="posts-submit"
                disabled={!draft.trim()}
                className="rounded-2xl bg-blue-600 px-5 py-3 text-sm font-semibold text-white transition hover:bg-blue-700 disabled:cursor-not-allowed disabled:bg-slate-300"
              >
                Post opinion
              </button>
            </div>
          </form>

          <div className="space-y-5">
            {posts.map((post) => {
              const attachedNews =
                post.linkedNews ??
                availableNews.find((newsItem) => newsItem.id === post.newsId) ??
                availableNews[0];
              const isFollowing = followedHandles.includes(post.handle);
              const isOwnPost = Boolean(
                sessionEmail && post.userEmail === sessionEmail
              );
              const displayedCommentCount =
                post.commentCount ?? post.comments.length;

              return (
                <article
                  key={post.id}
                  id={`post-${post.id}`}
                  data-cy="post-card"
                  className="rounded-[28px] bg-white p-6 shadow-sm transition-all"
                >
                  <div className="flex items-start justify-between gap-4">
                    <div>
                      <h3 className="text-base font-bold text-slate-900">
                        {post.author}
                      </h3>
                      <p className="mt-1 text-sm text-slate-500">
                        {post.handle} | {post.postedAt}
                      </p>
                    </div>
                    <div className="flex items-center gap-2">
                      {!isOwnPost && (
                        <button
                          type="button"
                          data-cy={`follow-${post.id}`}
                          onClick={() => toggleFollow(post)}
                          className={`inline-flex items-center gap-2 rounded-full px-3 py-2 text-xs font-semibold uppercase tracking-[0.14em] transition ${
                            isFollowing
                              ? 'bg-slate-900 text-white hover:bg-slate-800'
                              : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
                          }`}
                        >
                          {isFollowing ? (
                            <UserMinusIcon className="h-4 w-4" />
                          ) : (
                            <UserPlusIcon className="h-4 w-4" />
                          )}
                          {isFollowing ? 'Following' : 'Follow'}
                        </button>
                      )}
                      {isOwnPost && (
                        <button
                          type="button"
                          onClick={() => deletePost(post.id)}
                          className="inline-flex items-center gap-2 rounded-full px-3 py-2 text-xs font-semibold uppercase tracking-[0.14em] transition bg-rose-100 text-rose-600 hover:bg-rose-200"
                        >
                          Delete
                        </button>
                      )}
                      <span className="rounded-full bg-slate-100 px-3 py-1 text-xs font-semibold uppercase tracking-[0.16em] text-slate-500">
                        Opinion
                      </span>
                    </div>
                  </div>

                  <p className="mt-5 text-[15px] leading-7 text-slate-700">
                    {post.body}
                  </p>

                  <div className="mt-5 flex flex-wrap items-center gap-3 border-y border-slate-100 py-4">
                    <button
                      type="button"
                      data-cy={`like-${post.id}`}
                      onClick={() => toggleLike(post.id)}
                      className={`inline-flex items-center gap-2 rounded-full px-4 py-2 text-sm font-medium transition ${
                        post.isLiked
                          ? 'bg-rose-50 text-rose-600 hover:bg-rose-100'
                          : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
                      }`}
                    >
                      {post.isLiked ? (
                        <HeartSolidIcon className="h-5 w-5" />
                      ) : (
                        <HeartIcon className="h-5 w-5" />
                      )}
                      Like {post.likeCount}
                    </button>

                    <button
                      type="button"
                      data-cy={`comment-trigger-${post.id}`}
                      className="inline-flex items-center gap-2 rounded-full bg-slate-100 px-4 py-2 text-sm font-medium text-slate-600 transition hover:bg-slate-200"
                    >
                      <ChatBubbleLeftRightIcon className="h-5 w-5" />
                      Comment {displayedCommentCount}
                    </button>

                  </div>

                  <div className="mt-5 rounded-[24px] border border-slate-200 bg-slate-50 p-4">
                    <div className="flex flex-col gap-3 sm:flex-row sm:items-end">
                      <label className="flex-1">
                        <span className="sr-only">Add a comment</span>
                        <textarea
                          rows={2}
                          value={commentDrafts[post.id] || ''}
                          onChange={(event) =>
                            setCommentDrafts((previousDrafts) => ({
                              ...previousDrafts,
                              [post.id]: event.target.value,
                            }))
                          }
                          data-cy={`comment-draft-${post.id}`}
                          placeholder="Add a quick comment..."
                          className="w-full resize-none rounded-[18px] border border-slate-200 bg-white px-4 py-3 text-sm leading-6 text-slate-700 outline-none transition placeholder:text-slate-400 focus:border-blue-400"
                        />
                      </label>
                      <button
                        type="button"
                        data-cy={`comment-submit-${post.id}`}
                        onClick={() => addComment(post.id)}
                        disabled={!commentDrafts[post.id]?.trim()}
                        className="rounded-2xl bg-blue-600 px-4 py-3 text-sm font-semibold text-white transition hover:bg-blue-700 disabled:cursor-not-allowed disabled:bg-slate-300"
                      >
                        Add comment
                      </button>
                    </div>

                    {post.comments.length > 0 && (
                      <div className="mt-4 space-y-3">
                        {post.comments.map((comment) => (
                          <div
                            key={comment.id}
                            className="rounded-[18px] bg-white px-4 py-3"
                          >
                            <p className="text-xs font-semibold uppercase tracking-[0.14em] text-slate-400">
                              {comment.author}
                            </p>
                            <p className="mt-1 text-sm leading-6 text-slate-600">
                              {comment.body}
                            </p>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>

                  {attachedNews && (
                    <a href={attachedNews.articleUrl || '#'} target="_blank" rel="noopener noreferrer" className="block mt-5 rounded-[24px] border border-blue-100 bg-blue-50 p-4 hover:bg-blue-100 transition">
                      <p className="text-xs font-semibold uppercase tracking-[0.16em] text-blue-600">
                        Attached News Summary (Click to View Full Article)
                      </p>
                      <h4 className="mt-2 text-sm font-bold text-slate-900">
                        {attachedNews.headline}
                      </h4>
                      <p className="mt-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                        {attachedNews.source} | {attachedNews.category}
                      </p>
                      <p className="mt-3 text-sm leading-6 text-slate-600">
                        {attachedNews.summary}
                      </p>
                    </a>
                  )}
                </article>
              );
            })}
          </div>
        </section>
      </div>
    </main>
  );
};

export default Posts;
