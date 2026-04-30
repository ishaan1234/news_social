import React, { ChangeEvent, FormEvent, useEffect, useState } from 'react';
import {
  ArrowTrendingUpIcon,
  ArrowUpOnSquareIcon,
  CalendarDaysIcon,
  ChatBubbleLeftRightIcon,
  EllipsisHorizontalIcon,
  EnvelopeIcon,
  LinkIcon,
  MapPinIcon,
  PencilSquareIcon,
  SparklesIcon,
  UserPlusIcon,
} from '@heroicons/react/24/outline';
import { HeartIcon } from '@heroicons/react/24/solid';
import {
  AuthSession,
  getInitials,
  getSessionDisplayName,
  getSessionHandle,
  isVerifiedAuthSession,
} from '../auth';

const apiBaseUrl = (process.env.REACT_APP_API_BASE_URL || '').replace(
  /\/$/,
  ''
);

const authHeaders = (authSession: AuthSession | null) => {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };

  if (authSession?.idToken) {
    headers.Authorization = `Bearer ${authSession.idToken}`;
  }

  return headers;
};

interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
  feed?: T;
}

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
    throw new Error('Request failed');
  }

  if (body && typeof body === 'object' && 'success' in body) {
    const apiBody = body as ApiResponse<T>;
    if (!apiBody.success) {
      throw new Error(apiBody.error || 'Request failed');
    }
    return (apiBody[payloadKey] ?? apiBody.data) as T;
  }

  return body as T;
};

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
  article: {
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
  };
}

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
    return `${diffMinutes}m`;
  }

  const diffHours = Math.floor(diffMinutes / 60);
  if (diffHours < 24) {
    return `${diffHours}h`;
  }

  return `${Math.floor(diffHours / 24)}d`;
};

type TimelineTab = 'Posts' | 'Replies' | 'Media' | 'Likes';

interface ProfileDraft {
  name: string;
  handle: string;
  role: string;
  bio: string;
  location: string;
  website: string;
  email: string;
}

interface TimelineItem {
  id: string;
  tab: TimelineTab;
  eyebrow: string;
  body: string;
  timestamp: string;
  replies: number;
  reposts: number;
  likes: number;
  views: string;
  mediaLabel?: string;
}

interface TrendItem {
  id: string;
  category: string;
  title: string;
  posts: string;
}

interface SuggestedProfile {
  id: string;
  name: string;
  handle: string;
  bio: string;
  initials: string;
  accent: string;
}

interface ProfileProps {
  authSession?: AuthSession | null;
}

const initialProfile: ProfileDraft = {
  name: 'Avery Stone',
  handle: '@averystone',
  role: 'Tech & Policy Editor',
  bio:
    'Covering AI infrastructure, regulation, and market narratives with a bias toward clear explainers and sharp social packaging.',
  location: 'Brooklyn, New York',
  website: 'newshub.example/avery',
  email: 'avery@newshub.example',
};

const buildProfileFromSession = (authSession: AuthSession | null): ProfileDraft => {
  if (!authSession) {
    return initialProfile;
  }

  const name = getSessionDisplayName(authSession, initialProfile.name);
  const handle = getSessionHandle(authSession, initialProfile.handle);
  const websiteHandle = handle.replace(/^@/, '') || 'profile';

  return {
    ...initialProfile,
    name,
    handle,
    website: `newshub.example/${websiteHandle}`,
    email: authSession.user?.email || initialProfile.email,
  };
};

const timelineTabs: TimelineTab[] = ['Posts', 'Replies', 'Media', 'Likes'];

const initialTimelineItems: TimelineItem[] = [];

const trendItems: TrendItem[] = [
  {
    id: 'trend-1',
    category: 'Technology',
    title: 'AI infrastructure spending',
    posts: '4,281 posts',
  },
  {
    id: 'trend-2',
    category: 'Politics',
    title: 'Economic credibility debate',
    posts: '2,019 posts',
  },
  {
    id: 'trend-3',
    category: 'Business',
    title: 'EV pricing pressure',
    posts: '1,644 posts',
  },
];

const suggestedProfiles: SuggestedProfile[] = [
  {
    id: 'follow-1',
    name: 'Maya Chen',
    handle: '@maya',
    bio: 'Technology desk editor watching chips, cloud, and developer tools.',
    initials: 'MC',
    accent: 'from-cyan-500 to-blue-500',
  },
  {
    id: 'follow-2',
    name: 'Jordan Lee',
    handle: '@jord',
    bio: 'World news and policy angles with a clean headline instinct.',
    initials: 'JL',
    accent: 'from-orange-400 to-rose-500',
  },
  {
    id: 'follow-3',
    name: 'Nina Patel',
    handle: '@nina',
    bio: 'Markets, earnings, and sentiment shifts across breaking stories.',
    initials: 'NP',
    accent: 'from-emerald-400 to-teal-600',
  },
];

const Profile: React.FC<ProfileProps> = ({ authSession = null }) => {
  const [profile, setProfile] = useState<ProfileDraft>(() =>
    buildProfileFromSession(authSession)
  );
  const [draftProfile, setDraftProfile] = useState<ProfileDraft>(() =>
    buildProfileFromSession(authSession)
  );
  const [activeTab, setActiveTab] = useState<TimelineTab>('Posts');
  const [isEditing, setIsEditing] = useState(false);
  const [followedProfiles, setFollowedProfiles] = useState<string[]>([]);
  const [timelineItems, setTimelineItems] = useState<TimelineItem[]>(initialTimelineItems);
  const [followersCount, setFollowersCount] = useState<number>(0);
  const [followingCount, setFollowingCount] = useState<number>(0);
  const profileInitials = getInitials(profile.name);

  useEffect(() => {
    const fetchTimeline = async () => {
      if (!isVerifiedAuthSession(authSession) || !authSession?.user?.email) {
        setTimelineItems([]);
        return;
      }
      try {
        const email = authSession.user.email;
        const backendPosts = await readApiData<BackendPost[]>(
          `/feed?user_email=${encodeURIComponent(email)}`,
          'feed',
          { headers: authHeaders(authSession) }
        );
        if (Array.isArray(backendPosts)) {
          const mappedItems: TimelineItem[] = backendPosts
            .filter((post) => post.user_email === email)
            .map((post) => ({
              id: post.id,
              tab: 'Posts',
              eyebrow: post.article.title,
              body: post.caption || post.article.summary || post.article.description || 'Viewed article',
              timestamp: formatPostTime(post.created_at),
              replies: post.comment_count,
              reposts: 0,
              likes: post.like_count,
              views: '0',
            }));
          setTimelineItems(mappedItems);
        }
      } catch (e) {
        setTimelineItems([]);
      }

      try {
        const email = authSession.user.email;
        const stats = await readApiData<{ followers: number; following: number }>(
          `/following?email=${encodeURIComponent(email)}`,
          'data',
          { headers: authHeaders(authSession) }
        );
        setFollowersCount(stats.followers);
        setFollowingCount(stats.following);
      } catch (e) {
        setFollowersCount(0);
        setFollowingCount(0);
      }

      try {
        const email = authSession.user.email;
        const profileData = await readApiData<any>(
          `/profile?email=${encodeURIComponent(email)}`,
          'data',
          { headers: authHeaders(authSession) }
        );
        if (profileData) {
          setProfile((prev) => {
            const newProfile = {
              ...prev,
              name: profileData.display_name || prev.name,
              handle: profileData.username || prev.handle,
              role: profileData.role || prev.role,
              bio: profileData.bio || prev.bio,
              location: profileData.location || prev.location,
              website: profileData.website || prev.website,
            };
            setDraftProfile(newProfile);
            return newProfile;
          });
        }
      } catch (e) {
        // ignore
      }
    };
    fetchTimeline();
  }, [authSession]);

  useEffect(() => {
    const syncedProfile = buildProfileFromSession(authSession);
    setProfile(syncedProfile);
    setDraftProfile(syncedProfile);
  }, [authSession]);

  const visibleTimelineItems = timelineItems.filter(
    (item) => item.tab === activeTab
  );

  const handleFieldChange =
    (field: keyof ProfileDraft) =>
    (event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
      setDraftProfile((previousDraft) => ({
        ...previousDraft,
        [field]: event.target.value,
      }));
    };

  const openEditor = () => {
    setDraftProfile(profile);
    setIsEditing(true);
  };

  const closeEditor = () => {
    setDraftProfile(profile);
    setIsEditing(false);
  };

  const saveProfile = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    if (authSession?.user?.email) {
      try {
        await readApiData<any>('/profile', 'data', {
          method: 'PUT',
          headers: authHeaders(authSession),
          body: JSON.stringify({
            email: authSession.user.email,
            display_name: draftProfile.name,
            username: draftProfile.handle,
            role: draftProfile.role,
            bio: draftProfile.bio,
            location: draftProfile.location,
            website: draftProfile.website,
          }),
        });
      } catch (e) {
        console.error('Failed to save profile', e);
      }
    }

    setProfile(draftProfile);
    setIsEditing(false);
  };

  const toggleFollow = (profileId: string) => {
    setFollowedProfiles((previousProfiles) =>
      previousProfiles.includes(profileId)
        ? previousProfiles.filter((currentProfileId) => currentProfileId !== profileId)
        : [...previousProfiles, profileId]
    );
  };

  return (
    <main data-cy="profile-page" className="relative overflow-hidden">
      <div className="mx-auto max-w-6xl px-4 py-6 sm:px-6 lg:py-8">
        <section className="overflow-hidden rounded-[32px] border border-slate-200 bg-white shadow-sm">
          <div className="px-5 py-6 sm:px-6">
            <div className="flex items-start justify-between gap-4">
              <div className="flex h-28 w-28 items-center justify-center rounded-[28px] border-4 border-white bg-gradient-to-br from-amber-300 via-orange-300 to-sky-300 text-3xl font-bold text-slate-950 shadow-lg shadow-slate-300/60 sm:h-32 sm:w-32 sm:text-4xl">
                {profileInitials}
              </div>

              <button
                type="button"
                onClick={openEditor}
                data-cy="profile-edit"
                className="inline-flex items-center gap-2 rounded-full border border-slate-300 bg-white px-5 py-2.5 text-sm font-semibold text-slate-700 transition hover:bg-slate-50"
              >
                <PencilSquareIcon className="h-5 w-5" />
                Edit profile
              </button>
            </div>

            <div className="mt-4">
              <div className="flex flex-wrap items-center gap-2">
                <h1 className="text-3xl font-bold tracking-tight text-slate-900">
                  {profile.name}
                </h1>
              </div>
              <p className="mt-1 text-sm text-slate-500">{profile.handle}</p>
              <p className="mt-4 max-w-3xl text-[15px] leading-7 text-slate-700">
                {profile.bio}
              </p>

              <div className="mt-4 flex flex-wrap gap-4 text-sm text-slate-500">
                <span className="inline-flex items-center gap-2">
                  <MapPinIcon className="h-5 w-5 text-slate-400" />
                  {profile.location}
                </span>
                <span className="inline-flex items-center gap-2">
                  <LinkIcon className="h-5 w-5 text-slate-400" />
                  {profile.website}
                </span>
                <span className="inline-flex items-center gap-2">
                  <CalendarDaysIcon className="h-5 w-5 text-slate-400" />
                  Joined February 2026
                </span>
              </div>

              <div className="mt-4 flex flex-wrap items-center gap-4 text-sm">
                <span className="text-slate-500">
                  <span className="font-semibold text-slate-900">{followingCount}</span> Following
                </span>
                <span className="text-slate-500">
                  <span className="font-semibold text-slate-900">{followersCount}</span> Followers
                </span>
              </div>
            </div>

            <div className="mt-6 border-t border-slate-100">
              <div className="flex overflow-x-auto">
                {timelineTabs.map((tab) => {
                  const isActive = tab === activeTab;

                  return (
                    <button
                      key={tab}
                      type="button"
                      onClick={() => setActiveTab(tab)}
                      data-cy={`profile-tab-${tab.toLowerCase()}`}
                      className={`relative min-w-[110px] px-4 py-4 text-sm font-semibold transition ${
                        isActive
                          ? 'text-slate-900'
                          : 'text-slate-500 hover:bg-slate-50 hover:text-slate-700'
                      }`}
                    >
                      {tab}
                      {isActive && (
                        <span className="absolute inset-x-4 bottom-0 h-1 rounded-full bg-sky-500" />
                      )}
                    </button>
                  );
                })}
              </div>
            </div>
          </div>
        </section>

        <div className="mt-6 grid gap-6 xl:grid-cols-[minmax(0,1fr)_340px]">
          <section className="overflow-hidden rounded-[32px] border border-slate-200 bg-white shadow-sm">
            <div className="border-b border-slate-100 px-5 py-4">
              <p className="text-xs font-semibold uppercase tracking-[0.18em] text-slate-400">
                {activeTab}
              </p>
              <p className="mt-2 text-sm text-slate-500">
                Timeline-style cards for the selected profile tab.
              </p>
            </div>

            <div>
              {visibleTimelineItems.length === 0 ? (
                <div className="px-5 py-8 text-center">
                  <p className="text-sm text-slate-500">No posts to show.</p>
                </div>
              ) : (
                visibleTimelineItems.map((item, index) => (
                  <a href={`#/posts?postId=${item.id}`} key={item.id} className="block group">
                    <article
                      data-cy="profile-activity-card"
                      className={`px-5 py-5 transition group-hover:bg-slate-50 ${
                        index < visibleTimelineItems.length - 1
                          ? 'border-b border-slate-100'
                          : ''
                      }`}
                    >
                      <div className="flex items-start gap-4">
                    <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl bg-gradient-to-br from-amber-300 via-orange-300 to-sky-300 text-sm font-bold text-slate-950">
                      {profileInitials}
                    </div>

                    <div className="min-w-0 flex-1">
                      <div className="flex items-start justify-between gap-3">
                        <div className="min-w-0">
                          <div className="flex flex-wrap items-center gap-2">
                            <p className="font-semibold text-slate-900">
                              {profile.name}
                            </p>
                            <span className="text-sm text-slate-500">
                              {profile.handle}
                            </span>
                            <span className="text-sm text-slate-400">
                              | {item.timestamp}
                            </span>
                          </div>
                          <p className="mt-1 text-xs font-semibold uppercase tracking-[0.14em] text-sky-600">
                            {item.eyebrow}
                          </p>
                        </div>

                        <button
                          type="button"
                          className="rounded-full p-2 text-slate-400 transition hover:bg-slate-100 hover:text-slate-600"
                          aria-label="Post options"
                        >
                          <EllipsisHorizontalIcon className="h-5 w-5" />
                        </button>
                      </div>

                      <p className="mt-3 text-[15px] leading-7 text-slate-700">
                        {item.body}
                      </p>

                      {item.mediaLabel && (
                        <div className="mt-4 overflow-hidden rounded-[24px] border border-slate-200">
                          <div className="h-40 bg-gradient-to-br from-sky-400 via-blue-500 to-slate-950" />
                          <div className="border-t border-slate-200 bg-white px-4 py-3">
                            <p className="text-xs font-semibold uppercase tracking-[0.14em] text-slate-400">
                              Media
                            </p>
                            <p className="mt-1 text-sm font-semibold text-slate-900">
                              {item.mediaLabel}
                            </p>
                          </div>
                        </div>
                      )}

                      <div className="mt-5 flex flex-wrap items-center gap-4 text-sm text-slate-500">
                        <span className="inline-flex items-center gap-2">
                          <ChatBubbleLeftRightIcon className="h-5 w-5" />
                          {item.replies}
                        </span>
                        <span className="inline-flex items-center gap-2">
                          <ArrowUpOnSquareIcon className="h-5 w-5" />
                          {item.reposts}
                        </span>
                        <span className="inline-flex items-center gap-2">
                          <HeartIcon className="h-5 w-5 text-rose-500" />
                          {item.likes}
                        </span>
                        <span className="inline-flex items-center gap-2">
                          <ArrowTrendingUpIcon className="h-5 w-5" />
                          {item.views}
                        </span>
                      </div>
                      </div>
                    </div>
                  </article>
                </a>
                )))}
            </div>
          </section>

          <aside className="space-y-6">

            <section className="rounded-[28px] border border-slate-200 bg-white p-5 shadow-sm">
              <h2 className="text-xl font-bold text-slate-900">About</h2>
              <div className="mt-4 space-y-3 text-sm text-slate-600">
                <p>{profile.role}</p>
                <p className="inline-flex items-center gap-2">
                  <MapPinIcon className="h-5 w-5 text-slate-400" />
                  {profile.location}
                </p>
                <p className="inline-flex items-center gap-2">
                  <EnvelopeIcon className="h-5 w-5 text-slate-400" />
                  {profile.email}
                </p>
                <p>Covers AI policy, markets, media strategy, and startups.</p>
              </div>
            </section>

            <section className="rounded-[28px] border border-slate-200 bg-white p-5 shadow-sm">
              <div className="flex items-center justify-between gap-3">
                <h2 className="text-xl font-bold text-slate-900">
                  Trending in NewsHub
                </h2>
                <SparklesIcon className="h-6 w-6 text-sky-500" />
              </div>

              <div className="mt-4 space-y-4">
                {trendItems.map((trend) => (
                  <article
                    key={trend.id}
                    className="rounded-[22px] bg-slate-50 px-4 py-4"
                  >
                    <p className="text-xs font-semibold uppercase tracking-[0.14em] text-slate-400">
                      {trend.category}
                    </p>
                    <h3 className="mt-2 text-sm font-semibold text-slate-900">
                      {trend.title}
                    </h3>
                    <p className="mt-2 text-sm text-slate-500">{trend.posts}</p>
                  </article>
                ))}
              </div>
            </section>

            <section className="rounded-[28px] border border-slate-200 bg-white p-5 shadow-sm">
              <h2 className="text-xl font-bold text-slate-900">Who to follow</h2>

              <div className="mt-4 space-y-4">
                {suggestedProfiles.map((person) => (
                  <article key={person.id} className="flex items-start gap-3">
                    <div
                      className={`flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl bg-gradient-to-br ${person.accent} text-sm font-semibold text-white`}
                    >
                      {person.initials}
                    </div>
                    <div className="min-w-0 flex-1">
                      <div className="flex items-start justify-between gap-3">
                        <div>
                          <p className="font-semibold text-slate-900">
                            {person.name}
                          </p>
                          <p className="text-sm text-slate-500">
                            {person.handle}
                          </p>
                        </div>
                        <button
                          type="button"
                          data-cy={`profile-follow-${person.id}`}
                          onClick={() => toggleFollow(person.id)}
                          className={`inline-flex items-center gap-2 rounded-full px-3 py-2 text-xs font-semibold uppercase tracking-[0.14em] transition ${
                            followedProfiles.includes(person.id)
                              ? 'border border-slate-200 bg-white text-slate-700 hover:bg-slate-50'
                              : 'bg-slate-900 text-white hover:bg-slate-800'
                          }`}
                        >
                          <UserPlusIcon className="h-4 w-4" />
                          {followedProfiles.includes(person.id) ? 'Following' : 'Follow'}
                        </button>
                      </div>
                      <p className="mt-2 text-sm leading-6 text-slate-500">
                        {person.bio}
                      </p>
                    </div>
                  </article>
                ))}
              </div>
            </section>
          </aside>
        </div>

        {isEditing && (
          <div className="fixed inset-0 z-[60] overflow-y-auto bg-slate-950/40 p-4 sm:p-6">
            <div className="flex min-h-full items-start justify-center">
              <div
                data-cy="profile-edit-modal"
                className="flex max-h-[calc(100vh-2rem)] w-full max-w-2xl flex-col overflow-hidden rounded-[28px] bg-white shadow-2xl sm:max-h-[calc(100vh-3rem)]"
              >
                <div className="shrink-0 border-b border-slate-100 px-5 py-4 sm:px-6 sm:py-5">
                  <p className="text-xs font-semibold uppercase tracking-[0.18em] text-sky-600">
                    Edit profile
                  </p>
                  <h2 className="mt-2 text-2xl font-bold text-slate-900">
                    Tune the timeline identity
                  </h2>
                </div>

                <form
                  onSubmit={saveProfile}
                  className="flex min-h-0 flex-1 flex-col"
                >
                  <div className="min-h-0 flex-1 overflow-y-auto px-5 py-5 sm:px-6 sm:py-6">
                    <div className="grid gap-4 sm:grid-cols-2">
                      <label className="block">
                        <span className="text-sm font-medium text-slate-700">
                          Name
                        </span>
                        <input
                          type="text"
                          value={draftProfile.name}
                          onChange={handleFieldChange('name')}
                          data-cy="profile-name-input"
                          className="mt-2 w-full rounded-[18px] border border-slate-200 px-4 py-3 text-sm text-slate-700 outline-none transition focus:border-sky-400"
                        />
                      </label>

                      <label className="block">
                        <span className="text-sm font-medium text-slate-700">
                          Role
                        </span>
                        <input
                          type="text"
                          value={draftProfile.role}
                          onChange={handleFieldChange('role')}
                          className="mt-2 w-full rounded-[18px] border border-slate-200 px-4 py-3 text-sm text-slate-700 outline-none transition focus:border-sky-400"
                        />
                      </label>

                      <label className="block">
                        <span className="text-sm font-medium text-slate-700">
                          Handle
                        </span>
                        <input
                          type="text"
                          value={draftProfile.handle}
                          onChange={handleFieldChange('handle')}
                          className="mt-2 w-full rounded-[18px] border border-slate-200 px-4 py-3 text-sm text-slate-700 outline-none transition focus:border-sky-400"
                        />
                      </label>

                      <label className="block">
                        <span className="text-sm font-medium text-slate-700">
                          Location
                        </span>
                        <input
                          type="text"
                          value={draftProfile.location}
                          onChange={handleFieldChange('location')}
                          className="mt-2 w-full rounded-[18px] border border-slate-200 px-4 py-3 text-sm text-slate-700 outline-none transition focus:border-sky-400"
                        />
                      </label>

                      <label className="block">
                        <span className="text-sm font-medium text-slate-700">
                          Website
                        </span>
                        <input
                          type="text"
                          value={draftProfile.website}
                          onChange={handleFieldChange('website')}
                          className="mt-2 w-full rounded-[18px] border border-slate-200 px-4 py-3 text-sm text-slate-700 outline-none transition focus:border-sky-400"
                        />
                      </label>

                      <label className="block">
                        <span className="text-sm font-medium text-slate-700">
                          Email
                        </span>
                        <input
                          type="email"
                          value={draftProfile.email}
                          onChange={handleFieldChange('email')}
                          className="mt-2 w-full rounded-[18px] border border-slate-200 px-4 py-3 text-sm text-slate-700 outline-none transition focus:border-sky-400"
                        />
                      </label>
                    </div>

                    <label className="mt-4 block">
                      <span className="text-sm font-medium text-slate-700">
                        Bio
                      </span>
                      <textarea
                        rows={4}
                        value={draftProfile.bio}
                        onChange={handleFieldChange('bio')}
                        className="mt-2 w-full resize-none rounded-[22px] border border-slate-200 px-4 py-4 text-sm leading-7 text-slate-700 outline-none transition focus:border-sky-400"
                      />
                    </label>

                    <div className="mt-4 rounded-[24px] bg-slate-50 px-4 py-4">
                      <p className="text-sm leading-6 text-slate-500">
                        This dialog scrolls inside the window so it should fit
                        smaller screens without needing to zoom out.
                      </p>
                    </div>
                  </div>

                  <div className="shrink-0 border-t border-slate-100 bg-white px-5 py-4 sm:px-6">
                    <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                      <p className="text-sm text-slate-500">
                        Changes stay on the frontend for now.
                      </p>
                      <div className="flex gap-3">
                        <button
                          type="button"
                          onClick={closeEditor}
                          data-cy="profile-cancel"
                          className="rounded-full border border-slate-200 px-4 py-3 text-sm font-semibold text-slate-600 transition hover:bg-slate-50"
                        >
                          Cancel
                        </button>
                        <button
                          type="submit"
                          data-cy="profile-save"
                          className="rounded-full bg-sky-500 px-5 py-3 text-sm font-semibold text-white transition hover:bg-sky-600"
                        >
                          Save profile
                        </button>
                      </div>
                    </div>
                  </div>
                </form>
              </div>
            </div>
          </div>
        )}
      </div>
    </main>
  );
};

export default Profile;
