import React from 'react';

export interface NewsCardProps {
  headline: string;
  summary: string;
  category: string;
  source: string;
  timeAgo: string;
  articleUrl: string;
}

const categoryColors: Record<string, string> = {
  Technology: 'bg-blue-100 text-blue-700',
  Politics: 'bg-red-100 text-red-700',
  Sports: 'bg-green-100 text-green-700',
  Business: 'bg-orange-100 text-orange-700',
  Science: 'bg-yellow-100 text-yellow-700',
  Health: 'bg-teal-100 text-teal-700',
  World: 'bg-purple-100 text-purple-700',
};

const NewsCard: React.FC<NewsCardProps> = ({
  headline,
  summary,
  category,
  source,
  timeAgo,
  articleUrl,
}) => {
  const badge = categoryColors[category] || 'bg-gray-100 text-gray-700';

  return (
    <div className="h-full w-full flex flex-col">
      {/* Top half — Placeholder graphic */}
      <div className="relative h-1/2 w-full flex-shrink-0 bg-gray-200 flex items-center justify-center">
        {/* Simple doughnut / ring shape */}
        <div className="w-32 h-32 rounded-full border-[16px] border-gray-400 bg-gray-200" />
        <span
          className={`absolute top-4 left-4 text-xs font-semibold px-3 py-1 rounded-full ${badge}`}
        >
          {category}
        </span>
      </div>

      {/* Bottom half — Content */}
      <div className="flex-1 flex flex-col justify-center px-6 sm:px-10 md:px-16 py-6 bg-white">
        <h2 className="text-xl sm:text-2xl md:text-3xl font-bold text-gray-900 leading-tight">
          {headline}
        </h2>
        <p className="mt-4 text-base sm:text-lg text-gray-600 leading-relaxed">
          {summary}
        </p>
        <div className="mt-6 flex items-center justify-between">
          <p className="text-sm text-gray-400">
            {source} &middot; {timeAgo}
          </p>
          <a
            href={articleUrl}
            target="_blank"
            rel="noopener noreferrer"
            className="text-sm font-medium text-blue-600 hover:text-blue-800 transition-colors"
          >
            View full article &rarr;
          </a>
        </div>
      </div>
    </div>
  );
};

export default NewsCard;
