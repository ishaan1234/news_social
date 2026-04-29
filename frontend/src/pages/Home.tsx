import React, { useCallback, useEffect, useRef, useState } from 'react';
import NewsCard, { NewsCardProps } from '../components/NewsCard';
import { ChevronUpIcon, ChevronDownIcon } from '@heroicons/react/24/outline';
import { savePostArticleDraft } from '../postArticleDraft';

interface BackendArticle {
  id?: string;
  source?: {
    name?: string;
  };
  title?: string;
  description?: string;
  url?: string;
  urlToImage?: string;
  publishedAt?: string;
  summary?: string;
}

interface BackendNewsResponse {
  articles?: BackendArticle[];
}

const placeholderNews: NewsCardProps[] = [
  {
    headline: 'Headline placeholder',
    summary: 'This is random news summary of headline placeholder.',
    category: 'Technology',
    source: 'Tech Daily',
    timeAgo: '2h ago',
    articleUrl: 'https://example.com',
    imageUrl: 'https://picsum.photos/seed/news1/800/600',
  },
  {
    headline: 'Headline placeholder',
    summary: 'This is random news summary of headline placeholder.',
    category: 'World',
    source: 'World News',
    timeAgo: '4h ago',
    articleUrl: 'https://example.com',
  },
  {
    headline: 'Headline placeholder',
    summary: 'This is random news summary of headline placeholder.',
    category: 'Business',
    source: 'Finance Wire',
    timeAgo: '5h ago',
    articleUrl: 'https://example.com',
  },
  {
    headline: 'Headline placeholder',
    summary: 'This is random news summary of headline placeholder.',
    category: 'Sports',
    source: 'Sports Central',
    timeAgo: '6h ago',
    articleUrl: 'https://example.com',
  },
  {
    headline: 'Headline placeholder',
    summary: 'This is random news summary of headline placeholder.',
    category: 'Science',
    source: 'Space Report',
    timeAgo: '8h ago',
    articleUrl: 'https://example.com',
  },
  {
    headline: 'Headline placeholder',
    summary: 'This is random news summary of headline placeholder.',
    category: 'Health',
    source: 'Health Today',
    timeAgo: '10h ago',
    articleUrl: 'https://example.com',
  },
];

const backendQuery = 'tesla';
const apiBaseUrl = (process.env.REACT_APP_API_BASE_URL || '').replace(
  /\/$/,
  ''
);

const formatTimeAgo = (publishedAt?: string) => {
  if (!publishedAt) {
    return 'Recently';
  }

  const publishedMs = new Date(publishedAt).getTime();
  if (Number.isNaN(publishedMs)) {
    return 'Recently';
  }

  const diffMinutes = Math.max(
    0,
    Math.floor((Date.now() - publishedMs) / 60000)
  );

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

  const diffDays = Math.floor(diffHours / 24);
  return `${diffDays}d ago`;
};

const mapBackendArticle = (article: BackendArticle): NewsCardProps => {
  return {
    headline: article.title?.trim() || 'Untitled article',
    summary:
      article.summary?.trim() ||
      article.description?.trim() ||
      'Summary unavailable.',
    category: 'News',
    source: article.source?.name?.trim() || 'Unknown source',
    timeAgo: formatTimeAgo(article.publishedAt),
    articleUrl: article.url?.trim() || 'https://example.com',
    articleId: article.id?.trim() || undefined,
    imageUrl: article.urlToImage?.trim() || undefined,
  };
};

const Home: React.FC = () => {
  const [articles, setArticles] = useState<NewsCardProps[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [loadError, setLoadError] = useState('');
  const [current, setCurrent] = useState(0);
  const containerRef = useRef<HTMLDivElement>(null);
  const isScrolling = useRef(false);
  const newsItems = articles.length > 0 ? articles : placeholderNews;
  const total = newsItems.length;

  const loadNews = useCallback(async (signal?: AbortSignal) => {
    setIsLoading(true);
    setLoadError('');

    try {
      const response = await fetch(
        `${apiBaseUrl}/news?q=${encodeURIComponent(backendQuery)}`,
        { signal }
      );

      if (!response.ok) {
        throw new Error(`request failed with status ${response.status}`);
      }

      const data = (await response.json()) as BackendNewsResponse;
      const nextArticles = (data.articles || []).map(mapBackendArticle);

      if (nextArticles.length === 0) {
        throw new Error('no articles returned');
      }

      setArticles(nextArticles);
      setCurrent(0);
    } catch (error) {
      if (error instanceof DOMException && error.name === 'AbortError') {
        return;
      }

      setArticles([]);
      setLoadError(
        'Unable to load live news. Showing placeholder stories instead.'
      );
    } finally {
      if (!signal?.aborted) {
        setIsLoading(false);
      }
    }
  }, []);

  const goTo = useCallback(
    (index: number) => {
      if (index < 0 || index >= total || isScrolling.current) return;
      isScrolling.current = true;
      setCurrent(index);
      setTimeout(() => {
        isScrolling.current = false;
      }, 600);
    },
    [total]
  );

  const goNext = useCallback(() => goTo(current + 1), [current, goTo]);
  const goPrev = useCallback(() => goTo(current - 1), [current, goTo]);
  const handleCreatePost = useCallback((item: NewsCardProps) => {
    savePostArticleDraft({
      id: item.articleId,
      url: item.articleUrl,
      title: item.headline,
      source: item.source,
      summary: item.summary,
      image_url: item.imageUrl,
    });

    if (typeof window !== 'undefined') {
      window.location.hash = '#/posts?compose=1';
    }
  }, []);

  useEffect(() => {
    const controller = new AbortController();
    void loadNews(controller.signal);

    return () => controller.abort();
  }, [loadNews]);

  useEffect(() => {
    setCurrent((previousCurrent) =>
      Math.min(previousCurrent, Math.max(newsItems.length - 1, 0))
    );
  }, [newsItems.length]);

  // Keyboard navigation
  useEffect(() => {
    const handleKey = (e: KeyboardEvent) => {
      if (e.key === 'ArrowDown' || e.key === 'ArrowRight') {
        e.preventDefault();
        goNext();
      } else if (e.key === 'ArrowUp' || e.key === 'ArrowLeft') {
        e.preventDefault();
        goPrev();
      }
    };
    window.addEventListener('keydown', handleKey);
    return () => window.removeEventListener('keydown', handleKey);
  }, [goNext, goPrev]);

  // Wheel / scroll navigation
  useEffect(() => {
    const el = containerRef.current;
    if (!el) return;

    const handleWheel = (e: WheelEvent) => {
      e.preventDefault();
      if (Math.abs(e.deltaY) < 30) return;
      if (e.deltaY > 0) goNext();
      else goPrev();
    };

    el.addEventListener('wheel', handleWheel, { passive: false });
    return () => el.removeEventListener('wheel', handleWheel);
  }, [goNext, goPrev]);

  // Touch swipe navigation
  useEffect(() => {
    const el = containerRef.current;
    if (!el) return;
    let touchStartY = 0;

    const onTouchStart = (e: TouchEvent) => {
      touchStartY = e.touches[0].clientY;
    };
    const onTouchEnd = (e: TouchEvent) => {
      const diff = touchStartY - e.changedTouches[0].clientY;
      if (Math.abs(diff) < 50) return;
      if (diff > 0) goNext();
      else goPrev();
    };

    el.addEventListener('touchstart', onTouchStart);
    el.addEventListener('touchend', onTouchEnd);
    return () => {
      el.removeEventListener('touchstart', onTouchStart);
      el.removeEventListener('touchend', onTouchEnd);
    };
  }, [goNext, goPrev]);

  return (
    <div
      ref={containerRef}
      data-cy="home-feed"
      className="relative w-full overflow-hidden"
      style={{ height: 'calc(100vh - 56px)' }}
    >
      {isLoading && (
        <div className="absolute left-4 top-4 z-20 rounded-2xl bg-slate-900/80 px-4 py-2 text-sm font-medium text-white">
          Loading live news...
        </div>
      )}

      {!isLoading && loadError && (
        <div className="absolute left-4 top-4 z-20 flex max-w-sm items-center gap-3 rounded-2xl bg-white/95 px-4 py-3 shadow-lg">
          <p className="text-sm text-slate-600">{loadError}</p>
          <button
            type="button"
            onClick={() => void loadNews()}
            className="rounded-full bg-blue-600 px-3 py-1.5 text-xs font-semibold text-white transition hover:bg-blue-700"
          >
            Retry
          </button>
        </div>
      )}

      {/* Sliding container */}
      <div
        data-cy="news-track"
        className="h-full transition-transform duration-500 ease-in-out"
        style={{ transform: `translateY(-${current * 100}%)` }}
      >
        {newsItems.map((item, i) => (
          <div
            key={i}
            className="h-full w-full flex-shrink-0"
            style={{ height: 'calc(100vh - 56px)' }}
          >
            <NewsCard {...item} onCreatePost={() => handleCreatePost(item)} />
          </div>
        ))}
      </div>

      {/* Navigation arrows */}
      <div className="absolute right-4 bottom-6 flex flex-col gap-2 z-10">
        <button
          onClick={goPrev}
          disabled={current === 0}
          aria-label="Previous article"
          data-cy="news-prev"
          className="p-2 bg-white/80 backdrop-blur rounded-full shadow hover:bg-white disabled:opacity-30 transition"
        >
          <ChevronUpIcon className="w-5 h-5 text-gray-700" />
        </button>
        <button
          onClick={goNext}
          disabled={current === total - 1}
          aria-label="Next article"
          data-cy="news-next"
          className="p-2 bg-white/80 backdrop-blur rounded-full shadow hover:bg-white disabled:opacity-30 transition"
        >
          <ChevronDownIcon className="w-5 h-5 text-gray-700" />
        </button>
      </div>

      {/* Dots indicator */}
      <div className="absolute right-4 top-1/2 -translate-y-1/2 flex flex-col gap-1.5 z-10">
        {newsItems.map((_, i) => (
          <button
            key={i}
            onClick={() => goTo(i)}
            aria-label={`Go to article ${i + 1}`}
            data-cy={`news-dot-${i}`}
            className={`w-2 h-2 rounded-full transition-all ${
              i === current ? 'bg-blue-600 scale-125' : 'bg-gray-300'
            }`}
          />
        ))}
      </div>
    </div>
  );
};

export default Home;
