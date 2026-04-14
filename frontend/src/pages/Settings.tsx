import React, { useState } from 'react';
import {
  BellAlertIcon,
  ChatBubbleLeftRightIcon,
  CheckCircleIcon,
  NewspaperIcon,
  ShieldCheckIcon,
  UserCircleIcon,
} from '@heroicons/react/24/outline';
import {
  AuthSession,
  getSessionDisplayName,
  isVerifiedAuthSession,
} from '../auth';

interface SettingsProps {
  authSession?: AuthSession | null;
}

type FeedMode = 'for-you' | 'latest' | 'following';
type SummaryDepth = 'quick' | 'balanced' | 'deep';
type AudienceMode = 'everyone' | 'followers' | 'mutuals';
type DmAccess = 'everyone' | 'following' | 'nobody';

interface ToggleRowProps {
  label: string;
  description: string;
  enabled: boolean;
  onToggle: () => void;
  dataCy: string;
}

const ToggleRow: React.FC<ToggleRowProps> = ({
  label,
  description,
  enabled,
  onToggle,
  dataCy,
}) => (
  <div className="flex items-start justify-between gap-4 rounded-[24px] border border-slate-200 bg-slate-50 px-4 py-4">
    <div>
      <p className="text-sm font-semibold text-slate-900">{label}</p>
      <p className="mt-1 text-sm leading-6 text-slate-500">{description}</p>
    </div>
    <button
      type="button"
      onClick={onToggle}
      aria-pressed={enabled}
      data-cy={dataCy}
      className={`relative mt-1 inline-flex h-7 w-12 shrink-0 items-center rounded-full transition ${
        enabled ? 'bg-slate-900' : 'bg-slate-300'
      }`}
    >
      <span
        className={`inline-block h-5 w-5 rounded-full bg-white transition ${
          enabled ? 'translate-x-6' : 'translate-x-1'
        }`}
      />
    </button>
  </div>
);

const feedModeOptions: Array<{ id: FeedMode; label: string }> = [
  { id: 'for-you', label: 'For you' },
  { id: 'latest', label: 'Latest' },
  { id: 'following', label: 'Following' },
];

const audienceOptions: Array<{ id: AudienceMode; label: string }> = [
  { id: 'everyone', label: 'Everyone' },
  { id: 'followers', label: 'Followers' },
  { id: 'mutuals', label: 'Mutuals' },
];

const dmOptions: Array<{ id: DmAccess; label: string }> = [
  { id: 'everyone', label: 'Everyone' },
  { id: 'following', label: 'Following only' },
  { id: 'nobody', label: 'Nobody' },
];

const summaryLabels: Record<SummaryDepth, string> = {
  quick: 'Quick read',
  balanced: 'Balanced',
  deep: 'Deep dive',
};

const quickLinks = [
  { label: 'Go to home', href: '#/', dataCy: 'settings-quick-link-home' },
  { label: 'Go to posts', href: '#/posts', dataCy: 'settings-quick-link-posts' },
  { label: 'Go to chat', href: '#/chat', dataCy: 'settings-quick-link-chat' },
  {
    label: 'Go to profile',
    href: '#/profile',
    dataCy: 'settings-quick-link-profile',
  },
];

const Settings: React.FC<SettingsProps> = ({ authSession = null }) => {
  const hasVerifiedSession = isVerifiedAuthSession(authSession);
  const [feedMode, setFeedMode] = useState<FeedMode>('for-you');
  const [summaryDepth, setSummaryDepth] = useState<SummaryDepth>('balanced');
  const [favoriteTopics, setFavoriteTopics] = useState(
    'technology, politics, markets'
  );
  const [smartSummaries, setSmartSummaries] = useState(true);
  const [followedTopicsFirst, setFollowedTopicsFirst] = useState(true);
  const [showCommunityTakes, setShowCommunityTakes] = useState(true);
  const [defaultAudience, setDefaultAudience] =
    useState<AudienceMode>('everyone');
  const [allowQuotePosts, setAllowQuotePosts] = useState(true);
  const [showLikeCounts, setShowLikeCounts] = useState(true);
  const [openComments, setOpenComments] = useState(true);
  const [dmAccess, setDmAccess] = useState<DmAccess>('following');
  const [readReceipts, setReadReceipts] = useState(false);
  const [activityStatus, setActivityStatus] = useState(true);
  const [sharedArticleChats, setSharedArticleChats] = useState(true);
  const [breakingNewsAlerts, setBreakingNewsAlerts] = useState(true);
  const [postActivityAlerts, setPostActivityAlerts] = useState(true);
  const [messageAlerts, setMessageAlerts] = useState(true);
  const [trendAlerts, setTrendAlerts] = useState(false);
  const [blurSensitiveMedia, setBlurSensitiveMedia] = useState(true);
  const [filterLowQualityReplies, setFilterLowQualityReplies] = useState(true);
  const [hideNewAccountRequests, setHideNewAccountRequests] = useState(false);
  const [mutedKeywords, setMutedKeywords] = useState('spoilers, giveaways');

  return (
    <main
      data-cy="settings-page"
      className="mx-auto min-h-[calc(100vh-56px)] max-w-6xl px-4 py-6 sm:px-6 lg:py-8"
    >
      <div className="grid gap-6 lg:grid-cols-[minmax(0,1.4fr)_340px]">
        <section className="space-y-6">
          <section className="rounded-[32px] bg-white p-6 shadow-sm sm:p-7">
            <div>
              <span className="inline-flex rounded-full bg-slate-100 px-3 py-1 text-xs font-semibold uppercase tracking-[0.18em] text-slate-500">
                Frontend placeholders
              </span>
              <h1 className="mt-4 text-3xl font-bold tracking-tight text-slate-900">
                Settings
              </h1>
              <p className="mt-3 max-w-2xl text-sm leading-6 text-slate-500">
                Tune how news summaries appear on home, how posts behave in the
                social feed, who can message you, and which alerts you want to
                receive. These controls are frontend-only for now.
              </p>
            </div>

          </section>

          <div className="grid gap-6 md:grid-cols-2">
            <section className="rounded-[32px] bg-white p-6 shadow-sm md:col-span-2">
              <div className="flex items-center gap-3">
                <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-slate-100">
                  <NewspaperIcon className="h-6 w-6 text-slate-700" />
                </div>
                <div>
                  <h2 className="text-xl font-bold text-slate-900">
                    Home feed and summaries
                  </h2>
                  <p className="text-sm text-slate-500">
                    Control how stories, summaries, and opinions show up together.
                  </p>
                </div>
              </div>

              <div className="mt-6">
                <p className="text-sm font-medium text-slate-700">
                  Default home feed
                </p>
                <div className="mt-3 flex flex-wrap gap-3">
                  {feedModeOptions.map((option) => {
                    const isActive = option.id === feedMode;

                    return (
                      <button
                        key={option.id}
                        type="button"
                        onClick={() => setFeedMode(option.id)}
                        data-cy={`settings-feed-mode-${option.id}`}
                        className={`rounded-full px-4 py-2 text-sm font-semibold transition ${
                          isActive
                            ? 'bg-slate-900 text-white'
                            : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
                        }`}
                      >
                        {option.label}
                      </button>
                    );
                  })}
                </div>
              </div>

              <div className="mt-6 grid gap-4 md:grid-cols-2">
                <label className="block">
                  <span className="text-sm font-medium text-slate-700">
                    Summary style
                  </span>
                  <select
                    value={summaryDepth}
                    onChange={(event) =>
                      setSummaryDepth(event.target.value as SummaryDepth)
                    }
                    data-cy="settings-summary-style"
                    className="mt-2 w-full rounded-[20px] border border-slate-200 bg-white px-4 py-3 text-sm text-slate-700 outline-none transition focus:border-slate-400"
                  >
                    <option value="quick">Quick read</option>
                    <option value="balanced">Balanced</option>
                    <option value="deep">Deep dive</option>
                  </select>
                </label>

                <label className="block">
                  <span className="text-sm font-medium text-slate-700">
                    Favorite topics
                  </span>
                  <input
                    type="text"
                    value={favoriteTopics}
                    onChange={(event) => setFavoriteTopics(event.target.value)}
                    data-cy="settings-favorite-topics"
                    className="mt-2 w-full rounded-[20px] border border-slate-200 px-4 py-3 text-sm text-slate-700 outline-none transition focus:border-slate-400"
                  />
                </label>
              </div>

              <div className="mt-6 space-y-4">
                <ToggleRow
                  label="Smart news summaries"
                  description="Show compact story summaries before opening the full article."
                  enabled={smartSummaries}
                  onToggle={() => setSmartSummaries((current) => !current)}
                  dataCy="settings-smart-summaries"
                />
                <ToggleRow
                  label="Prioritize followed topics"
                  description="Push stories related to followed topics higher in the home feed."
                  enabled={followedTopicsFirst}
                  onToggle={() => setFollowedTopicsFirst((current) => !current)}
                  dataCy="settings-followed-topics-first"
                />
                <ToggleRow
                  label="Show community takes under stories"
                  description="Surface opinions and reactions directly below news cards."
                  enabled={showCommunityTakes}
                  onToggle={() => setShowCommunityTakes((current) => !current)}
                  dataCy="settings-community-takes"
                />
              </div>
            </section>

            <section className="rounded-[32px] bg-white p-6 shadow-sm">
              <div className="flex items-center gap-3">
                <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-sky-50">
                  <UserCircleIcon className="h-6 w-6 text-sky-600" />
                </div>
                <div>
                  <h2 className="text-xl font-bold text-slate-900">
                    Posting and replies
                  </h2>
                  <p className="text-sm text-slate-500">
                    Set defaults for opinions, comments, and reactions.
                  </p>
                </div>
              </div>

              <div className="mt-6">
                <p className="text-sm font-medium text-slate-700">
                  Default audience
                </p>
                <div className="mt-3 flex flex-wrap gap-3">
                  {audienceOptions.map((option) => {
                    const isActive = option.id === defaultAudience;

                    return (
                      <button
                        key={option.id}
                        type="button"
                        onClick={() => setDefaultAudience(option.id)}
                        data-cy={`settings-audience-${option.id}`}
                        className={`rounded-full px-4 py-2 text-sm font-semibold transition ${
                          isActive
                            ? 'bg-slate-900 text-white'
                            : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
                        }`}
                      >
                        {option.label}
                      </button>
                    );
                  })}
                </div>
              </div>

              <div className="mt-6 space-y-4">
                <ToggleRow
                  label="Allow quote reposts"
                  description="Let other people share your post with their own take."
                  enabled={allowQuotePosts}
                  onToggle={() => setAllowQuotePosts((current) => !current)}
                  dataCy="settings-quote-posts"
                />
                <ToggleRow
                  label="Show like counts on my posts"
                  description="Keep public reaction counts visible on your content."
                  enabled={showLikeCounts}
                  onToggle={() => setShowLikeCounts((current) => !current)}
                  dataCy="settings-like-counts"
                />
                <ToggleRow
                  label="Open comments by default"
                  description="New opinion posts allow replies unless you switch this off."
                  enabled={openComments}
                  onToggle={() => setOpenComments((current) => !current)}
                  dataCy="settings-open-comments"
                />
              </div>
            </section>

            <section className="rounded-[32px] bg-white p-6 shadow-sm">
              <div className="flex items-center gap-3">
                <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-amber-50">
                  <BellAlertIcon className="h-6 w-6 text-amber-600" />
                </div>
                <div>
                  <h2 className="text-xl font-bold text-slate-900">
                    Notification preferences
                  </h2>
                  <p className="text-sm text-slate-500">
                    Decide which activity should pull you back into the app.
                  </p>
                </div>
              </div>

              <div className="mt-6 space-y-4">
                <ToggleRow
                  label="Breaking news alerts"
                  description="Push urgent story updates when major topics move fast."
                  enabled={breakingNewsAlerts}
                  onToggle={() =>
                    setBreakingNewsAlerts((current) => !current)
                  }
                  dataCy="settings-breaking-news-alerts"
                />
                <ToggleRow
                  label="Likes, comments, and reposts"
                  description="Notify you when people engage with your opinions and posts."
                  enabled={postActivityAlerts}
                  onToggle={() => setPostActivityAlerts((current) => !current)}
                  dataCy="settings-post-activity-alerts"
                />
                <ToggleRow
                  label="Direct messages"
                  description="Alert you when someone starts or replies in a chat."
                  enabled={messageAlerts}
                  onToggle={() => setMessageAlerts((current) => !current)}
                  dataCy="settings-message-alerts"
                />
                <ToggleRow
                  label="Trending conversations"
                  description="Recommend fast-moving discussions connected to the news feed."
                  enabled={trendAlerts}
                  onToggle={() => setTrendAlerts((current) => !current)}
                  dataCy="settings-trend-alerts"
                />
              </div>
            </section>

            <section className="rounded-[32px] bg-white p-6 shadow-sm">
              <div className="flex items-center gap-3">
                <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-indigo-50">
                  <ChatBubbleLeftRightIcon className="h-6 w-6 text-indigo-600" />
                </div>
                <div>
                  <h2 className="text-xl font-bold text-slate-900">
                    Messages and privacy
                  </h2>
                  <p className="text-sm text-slate-500">
                    Control who can reach you and how private chat feels.
                  </p>
                </div>
              </div>

              <div className="mt-6">
                <p className="text-sm font-medium text-slate-700">
                  Who can message me
                </p>
                <div className="mt-3 flex flex-wrap gap-3">
                  {dmOptions.map((option) => {
                    const isActive = option.id === dmAccess;

                    return (
                      <button
                        key={option.id}
                        type="button"
                        onClick={() => setDmAccess(option.id)}
                        data-cy={`settings-dm-access-${option.id}`}
                        className={`rounded-full px-4 py-2 text-sm font-semibold transition ${
                          isActive
                            ? 'bg-slate-900 text-white'
                            : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
                        }`}
                      >
                        {option.label}
                      </button>
                    );
                  })}
                </div>
              </div>

              <div className="mt-6 space-y-4">
                <ToggleRow
                  label="Read receipts"
                  description="Let chat partners know when you have seen their message."
                  enabled={readReceipts}
                  onToggle={() => setReadReceipts((current) => !current)}
                  dataCy="settings-read-receipts"
                />
                <ToggleRow
                  label="Show active status"
                  description="Display when you are currently online and available."
                  enabled={activityStatus}
                  onToggle={() => setActivityStatus((current) => !current)}
                  dataCy="settings-activity-status"
                />
                <ToggleRow
                  label="Allow chats from shared articles"
                  description="Open a one-to-one conversation directly from a story share."
                  enabled={sharedArticleChats}
                  onToggle={() => setSharedArticleChats((current) => !current)}
                  dataCy="settings-shared-article-chats"
                />
              </div>
            </section>

            <section className="rounded-[32px] bg-white p-6 shadow-sm">
              <div className="flex items-center gap-3">
                <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-emerald-50">
                  <ShieldCheckIcon className="h-6 w-6 text-emerald-600" />
                </div>
                <div>
                  <h2 className="text-xl font-bold text-slate-900">
                    Safety filters
                  </h2>
                  <p className="text-sm text-slate-500">
                    Placeholder controls for keeping replies and chats manageable.
                  </p>
                </div>
              </div>

              <label className="mt-6 block">
                <span className="text-sm font-medium text-slate-700">
                  Muted keywords
                </span>
                <textarea
                  value={mutedKeywords}
                  onChange={(event) => setMutedKeywords(event.target.value)}
                  data-cy="settings-muted-keywords"
                  rows={4}
                  className="mt-2 w-full rounded-[20px] border border-slate-200 px-4 py-3 text-sm leading-6 text-slate-700 outline-none transition focus:border-slate-400"
                />
              </label>

              <div className="mt-6 space-y-4">
                <ToggleRow
                  label="Blur sensitive media"
                  description="Hide potentially sensitive images until you choose to open them."
                  enabled={blurSensitiveMedia}
                  onToggle={() => setBlurSensitiveMedia((current) => !current)}
                  dataCy="settings-blur-sensitive-media"
                />
                <ToggleRow
                  label="Filter low-quality replies"
                  description="Down-rank noisy replies beneath your posts and news comments."
                  enabled={filterLowQualityReplies}
                  onToggle={() =>
                    setFilterLowQualityReplies((current) => !current)
                  }
                  dataCy="settings-filter-low-quality-replies"
                />
                <ToggleRow
                  label="Hide requests from brand new accounts"
                  description="Keep first-time message requests quieter until you review them."
                  enabled={hideNewAccountRequests}
                  onToggle={() =>
                    setHideNewAccountRequests((current) => !current)
                  }
                  dataCy="settings-hide-new-account-requests"
                />
              </div>
            </section>
          </div>
        </section>

        <aside className="space-y-6">
          <section className="rounded-[32px] bg-white p-6 shadow-sm">
            <p className="text-xs font-semibold uppercase tracking-[0.16em] text-slate-400">
              Account snapshot
            </p>
            <h2 className="mt-3 text-xl font-bold text-slate-900">
              {getSessionDisplayName(authSession, 'NewsHub user')}
            </h2>
            <p className="mt-1 text-sm text-slate-500">
              {authSession?.user?.email ||
                'Sign in to tie these settings to your profile later.'}
            </p>
            <div className="mt-4 inline-flex rounded-full bg-slate-100 px-3 py-1 text-xs font-semibold uppercase tracking-[0.14em] text-slate-600">
              {hasVerifiedSession ? 'Verified profile' : 'Guest placeholder'}
            </div>

            <div className="mt-5 space-y-3 rounded-[24px] border border-slate-200 bg-slate-50 p-4">
              <div
                data-cy="settings-snapshot-feed-mode"
                className="flex items-center justify-between text-sm"
              >
                <span className="text-slate-500">Home feed</span>
                <span className="font-semibold text-slate-900">
                  {feedModeOptions.find((option) => option.id === feedMode)?.label}
                </span>
              </div>
              <div
                data-cy="settings-snapshot-summary-style"
                className="flex items-center justify-between text-sm"
              >
                <span className="text-slate-500">Summary style</span>
                <span className="font-semibold text-slate-900">
                  {summaryLabels[summaryDepth]}
                </span>
              </div>
              <div
                data-cy="settings-snapshot-dm-access"
                className="flex items-center justify-between text-sm"
              >
                <span className="text-slate-500">DM access</span>
                <span className="font-semibold text-slate-900">
                  {dmOptions.find((option) => option.id === dmAccess)?.label}
                </span>
              </div>
            </div>
          </section>

          <section className="rounded-[32px] bg-white p-6 shadow-sm">
            <div className="flex items-center gap-3">
              <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-rose-50">
                <CheckCircleIcon className="h-6 w-6 text-rose-600" />
              </div>
              <div>
                <h2 className="text-xl font-bold text-slate-900">
                  Reader blend
                </h2>
                <p className="text-sm text-slate-500">
                  Static placeholders for what shapes your timeline.
                </p>
              </div>
            </div>

            <div className="mt-6 flex flex-wrap gap-2">
              {['Top stories', 'Following', 'Opinions', 'Local news'].map(
                (item) => (
                  <span
                    key={item}
                    className="rounded-full bg-slate-100 px-3 py-1.5 text-sm font-medium text-slate-600"
                  >
                    {item}
                  </span>
                )
              )}
            </div>

            <div className="mt-6 space-y-3">
              {[
                ['Saved stories', '18'],
                ['Draft posts', '4'],
                ['Unread chats', '7'],
              ].map(([label, value]) => (
                <div
                  key={label}
                  className="flex items-center justify-between rounded-[22px] border border-slate-200 px-4 py-3"
                >
                  <span className="text-sm text-slate-500">{label}</span>
                  <span className="text-sm font-semibold text-slate-900">
                    {value}
                  </span>
                </div>
              ))}
            </div>
          </section>

          <section className="rounded-[32px] bg-white p-6 shadow-sm">
            <h2 className="text-xl font-bold text-slate-900">Quick links</h2>
            <p className="mt-2 text-sm leading-6 text-slate-500">
              Jump back to the surfaces these settings affect.
            </p>

            <div data-cy="settings-quick-links" className="mt-6 space-y-3">
              {quickLinks.map(({ label, href, dataCy }) => (
                <a
                  key={label}
                  href={href}
                  data-cy={dataCy}
                  className="block rounded-[22px] border border-slate-200 px-4 py-3 text-sm font-semibold text-slate-700 transition hover:border-slate-300 hover:bg-slate-50"
                >
                  {label}
                </a>
              ))}
            </div>
          </section>
        </aside>
      </div>
    </main>
  );
};

export default Settings;
