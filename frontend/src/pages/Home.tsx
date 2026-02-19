import React, { useCallback, useEffect, useRef, useState } from 'react';
import NewsCard, { NewsCardProps } from '../components/NewsCard';
import { ChevronUpIcon, ChevronDownIcon } from '@heroicons/react/24/outline';

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

const Home: React.FC = () => {
  const [current, setCurrent] = useState(0);
  const containerRef = useRef<HTMLDivElement>(null);
  const isScrolling = useRef(false);
  const total = placeholderNews.length;

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
      className="relative w-full overflow-hidden"
      style={{ height: 'calc(100vh - 56px)' }}
    >
      {/* Sliding container */}
      <div
        className="h-full transition-transform duration-500 ease-in-out"
        style={{ transform: `translateY(-${current * 100}%)` }}
      >
        {placeholderNews.map((item, i) => (
          <div key={i} className="h-full w-full flex-shrink-0" style={{ height: 'calc(100vh - 56px)' }}>
            <NewsCard {...item} />
          </div>
        ))}
      </div>

      {/* Navigation arrows */}
      <div className="absolute right-4 bottom-6 flex flex-col gap-2 z-10">
        <button
          onClick={goPrev}
          disabled={current === 0}
          className="p-2 bg-white/80 backdrop-blur rounded-full shadow hover:bg-white disabled:opacity-30 transition"
        >
          <ChevronUpIcon className="w-5 h-5 text-gray-700" />
        </button>
        <button
          onClick={goNext}
          disabled={current === total - 1}
          className="p-2 bg-white/80 backdrop-blur rounded-full shadow hover:bg-white disabled:opacity-30 transition"
        >
          <ChevronDownIcon className="w-5 h-5 text-gray-700" />
        </button>
      </div>

      {/* Dots indicator */}
      <div className="absolute right-4 top-1/2 -translate-y-1/2 flex flex-col gap-1.5 z-10">
        {placeholderNews.map((_, i) => (
          <button
            key={i}
            onClick={() => goTo(i)}
            className={`w-2 h-2 rounded-full transition-all ${i === current ? 'bg-blue-600 scale-125' : 'bg-gray-300'
              }`}
          />
        ))}
      </div>
    </div>
  );
};

export default Home;
