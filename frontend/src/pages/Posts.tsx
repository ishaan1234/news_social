import React, { FormEvent, useCallback, useEffect, useMemo, useState } from 'react';
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

interface NewsReference {
  id: string;
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
  backendId?: number;
  author: string;
  handle: string;
  postedAt: string;
  body: string;
  newsId: string;
  linkedNews?: NewsReference;
  likeCount: number;
  shareCount: number;
  isLiked: boolean;
  viewerVote?: number;
  commentCount?: number;
  comments: PostComment[];
}

interface PostComment {
  id: string;
  backendId?: number;
  author: string;
  body: string;
}

interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}

interface BackendPost {
  id: number;
  author_name?: string;
  author_handle?: string;
  body: string;
  article: LinkedArticleDraft;
  vote_score: number;
  viewer_vote?: number;
  comment_count: number;
  share_count: number;
  created_at: string;
}

interface BackendPostComment {
  id: number;
  post_id: number;
  author_name?: string;
  content: string;
}

interface BackendVoteSummary {
  post_id: number;
  vote_score: number;
  viewer_vote: number;
}

const placeholderNews: NewsReference[] = [
  {
    id: 'chip-policy',
    headline: 'Chip export rules tighten as AI demand keeps rising',
    source: 'Tech Daily',
    category: 'Technology',
    summary:
      'Regulators are tightening export restrictions on advanced chips while cloud and AI spending remains elevated. Companies now need to balance compliance risk, supply chain planning, and investor pressure around long-term growth.',
    articleUrl: 'https://example.com/chip-policy',
  },
  {
    id: 'election-debate',
    headline: 'Election debate shifts focus toward economic credibility',
    source: 'World Wire',
    category: 'Politics',
    summary:
      'The latest debate centered on inflation, wages, and public trust in economic leadership. Analysts say the exchange may matter less for headline moments and more for how undecided voters judge competence and stability.',
    articleUrl: 'https://example.com/election-debate',
  },
  {
    id: 'ev-market',
    headline: 'EV makers push expansion as pricing pressure intensifies',
    source: 'Market Brief',
    category: 'Business',
    summary:
      'Electric vehicle companies are expanding capacity and retail presence even as competition pushes margins lower. The main question is whether volume growth can offset price pressure quickly enough to preserve investor confidence.',
    articleUrl: 'https://example.com/ev-market',
  },
];

const initialPosts: OpinionPost[] = [
  {
    id: 'post-1',
    author: 'Maya Chen',
    handle: '@maya',
    postedAt: '12m ago',
    body:
      'Most coverage is treating this like a policy shock, but the bigger story is execution risk. If supply planning lags, the headline impact will outlast the announcement itself.',
    newsId: 'chip-policy',
    likeCount: 14,
    shareCount: 3,
    isLiked: false,
    comments: [
      {
        id: 'comment-1',
        author: 'Nina',
        body: 'That execution-risk angle is stronger than most of the headlines.',
      },
    ],
  },
  {
    id: 'post-2',
    author: 'Jordan Lee',
    handle: '@jord',
    postedAt: '28m ago',
    body:
      'The debate takeaway is not who had the sharpest line. It is which candidate sounded like they understood household economics in practical terms.',
    newsId: 'election-debate',
    likeCount: 9,
    shareCount: 2,
    isLiked: true,
    comments: [
      {
        id: 'comment-2',
        author: 'Maya',
        body: 'Agreed. Tone mattered less than whether the answer felt grounded.',
      },
    ],
  },
];

const apiBaseUrl = (process.env.REACT_APP_API_BASE_URL || '').replace(/\/$/, '');
const viewerStorageKey = 'news-social-post-viewer-id';

const getStableArticleId = (url: string) =>
  `article-${url.replace(/[^a-zA-Z0-9]+/g, '-').slice(0, 48) || 'linked'}`;

const toNewsReference = (article: LinkedArticleDraft): NewsReference => ({
  id: getStableArticleId(article.url),
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

const getViewerId = () => {
  if (typeof window === 'undefined') {
    return 'server-render-viewer';
  }

  const existingId = window.localStorage.getItem(viewerStorageKey);
  if (existingId) {
    return existingId;
  }

  const nextId = `viewer-${Date.now()}-${Math.random()
    .toString(36)
    .slice(2, 8)}`;
  window.localStorage.setItem(viewerStorageKey, nextId);
  return nextId;
};

const readApiData = async <T,>(path: string, init?: RequestInit): Promise<T> => {
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
    return apiBody.data as T;
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
  const linkedNews = toNewsReference(post.article);

  return {
    id: `api-post-${post.id}`,
    backendId: post.id,
    author: post.author_name || 'Anonymous',
    handle: post.author_handle || '@newshub',
    postedAt: formatPostTime(post.created_at),
    body: post.body,
    newsId: linkedNews.id,
    linkedNews,
    likeCount: post.vote_score,
    shareCount: post.share_count,
    isLiked: (post.viewer_vote || 0) > 0,
    viewerVote: post.viewer_vote || 0,
    commentCount: post.comment_count,
    comments: [],
  };
};

const createPostId = () =>
  `post-${Date.now()}-${Math.random().toString(36).slice(2, 7)}`;

const createCommentId = () =>
  `comment-${Date.now()}-${Math.random().toString(36).slice(2, 7)}`;

const Posts: React.FC = () => {
  const [viewerId] = useState(getViewerId);
  const [linkedNews] = useState<NewsReference | null>(() =>
    createLinkedNewsFromDraft()
  );
  const [posts, setPosts] = useState<OpinionPost[]>(initialPosts);
  const [draft, setDraft] = useState('');
  const [selectedNewsId, setSelectedNewsId] = useState(
    () => linkedNews?.id || placeholderNews[0].id
  );
  const [followedHandles, setFollowedHandles] = useState<string[]>(['@maya']);
  const [commentDrafts, setCommentDrafts] = useState<Record<string, string>>({});
  const availableNews = useMemo(
    () => (linkedNews ? [linkedNews, ...placeholderNews] : placeholderNews),
    [linkedNews]
  );

  const selectedNews = useMemo(
    () =>
      availableNews.find((item) => item.id === selectedNewsId) ??
      availableNews[0],
    [availableNews, selectedNewsId]
  );

  const loadPosts = useCallback(async () => {
    try {
      const backendPosts = await readApiData<BackendPost[]>(
        `/api/posts?viewer_id=${encodeURIComponent(viewerId)}`
      );

      if (Array.isArray(backendPosts)) {
        setPosts(backendPosts.map(mapBackendPost));
      }
    } catch (_error) {
      // Keep the local starter posts when the database-backed API is unavailable.
    }
  }, [viewerId]);

  useEffect(() => {
    void loadPosts();
  }, [loadPosts]);

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const trimmedDraft = draft.trim();
    if (!trimmedDraft) {
      return;
    }

    let nextPost: OpinionPost = {
      id: createPostId(),
      author: 'You',
      handle: '@you',
      postedAt: 'Just now',
      body: trimmedDraft,
      newsId: selectedNews.id,
      linkedNews: selectedNews,
      likeCount: 0,
      shareCount: 0,
      isLiked: false,
      comments: [],
    };

    try {
      const createdPost = await readApiData<BackendPost>('/api/posts', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          author_id: viewerId,
          author_name: 'You',
          author_handle: '@you',
          body: trimmedDraft,
          article: {
            url: selectedNews.articleUrl,
            title: selectedNews.headline,
            source: selectedNews.source,
            summary: selectedNews.summary,
            image_url: selectedNews.imageUrl,
            published_at: selectedNews.publishedAt,
          },
        }),
      });
      nextPost = mapBackendPost(createdPost);
    } catch (_error) {
      // Local fallback keeps the composer usable without a configured database.
    }

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
    const nextVote = nextLiked ? 1 : 0;

    setPosts((previousPosts) =>
      previousPosts.map((post) => {
        if (post.id !== postId) {
          return post;
        }

        return {
          ...post,
          isLiked: nextLiked,
          viewerVote: nextVote,
          likeCount: post.likeCount + (nextLiked ? 1 : -1),
        };
      })
    );

    if (!targetPost.backendId) {
      return;
    }

    void readApiData<BackendVoteSummary>('/api/posts/votes', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        post_id: targetPost.backendId,
        voter_id: viewerId,
        value: nextVote,
      }),
    })
      .then((summary) => {
        setPosts((previousPosts) =>
          previousPosts.map((post) =>
            post.id === postId
              ? {
                ...post,
                isLiked: summary.viewer_vote > 0,
                viewerVote: summary.viewer_vote,
                likeCount: summary.vote_score,
              }
              : post
          )
        );
      })
      .catch(() => undefined);
  };

  const sharePost = (postId: string) => {
    const targetPost = posts.find((post) => post.id === postId);

    setPosts((previousPosts) =>
      previousPosts.map((post) =>
        post.id === postId
          ? { ...post, shareCount: post.shareCount + 1 }
          : post
      )
    );

    if (!targetPost?.backendId) {
      return;
    }

    void readApiData<{ post_id: number; share_count: number }>(
      '/api/posts/share',
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ post_id: targetPost.backendId }),
      }
    )
      .then((result) => {
        setPosts((previousPosts) =>
          previousPosts.map((post) =>
            post.id === postId
              ? { ...post, shareCount: result.share_count }
              : post
          )
        );
      })
      .catch(() => undefined);
  };

  const toggleFollow = (handle: string) => {
    setFollowedHandles((previousHandles) =>
      previousHandles.includes(handle)
        ? previousHandles.filter((currentHandle) => currentHandle !== handle)
        : [...previousHandles, handle]
    );
  };

  const addComment = (postId: string) => {
    const trimmedComment = commentDrafts[postId]?.trim();
    if (!trimmedComment) {
      return;
    }

    const targetPost = posts.find((post) => post.id === postId);
    const nextComment: PostComment = {
      id: createCommentId(),
      author: 'You',
      body: trimmedComment,
    };

    setPosts((previousPosts) =>
      previousPosts.map((post) =>
        post.id === postId
          ? {
            ...post,
            comments: [
              ...post.comments,
              nextComment,
            ],
            commentCount: (post.commentCount ?? post.comments.length) + 1,
          }
          : post
      )
    );

    setCommentDrafts((previousDrafts) => ({
      ...previousDrafts,
      [postId]: '',
    }));

    if (!targetPost?.backendId) {
      return;
    }

    void readApiData<BackendPostComment>('/api/posts/comments', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        post_id: targetPost.backendId,
        author_id: viewerId,
        author_name: 'You',
        content: trimmedComment,
      }),
    }).catch(() => undefined);
  };

  return (
    <main
      data-cy="posts-page"
      className="mx-auto max-w-6xl px-4 py-6 sm:px-6"
    >
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
                    className={`w-full rounded-3xl border px-4 py-4 text-left transition ${isSelected
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

            <div className="mt-5 flex items-center justify-between gap-4">
              <p className="text-xs text-slate-400">
                Frontend demo only. Nothing is persisted to a server.
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
              const displayedCommentCount =
                post.commentCount ?? post.comments.length;

              return (
                <article
                  key={post.id}
                  data-cy="post-card"
                  className="rounded-[28px] bg-white p-6 shadow-sm"
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
                      <button
                        type="button"
                        data-cy={`follow-${post.id}`}
                        onClick={() => toggleFollow(post.handle)}
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

                    <button
                      type="button"
                      data-cy={`share-${post.id}`}
                      onClick={() => sharePost(post.id)}
                      className="inline-flex items-center gap-2 rounded-full bg-slate-100 px-4 py-2 text-sm font-medium text-slate-600 transition hover:bg-slate-200"
                    >
                      <ArrowUpOnSquareIcon className="h-5 w-5" />
                      Share {post.shareCount}
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

                  <div className="mt-5 rounded-[24px] border border-blue-100 bg-blue-50 p-4">
                    <p className="text-xs font-semibold uppercase tracking-[0.16em] text-blue-600">
                      Attached News Summary
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
                  </div>
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
