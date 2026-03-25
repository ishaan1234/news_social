import React, { FormEvent, useEffect, useRef, useState } from 'react';
import {
  ChatBubbleLeftRightIcon,
  EllipsisHorizontalIcon,
  MagnifyingGlassIcon,
  PaperAirplaneIcon,
} from '@heroicons/react/24/outline';

type MessageStatus = 'Sent' | 'Delivered' | 'Read';

interface Message {
  id: string;
  sender: 'me' | 'friend';
  text: string;
  sentAt: string;
  status?: MessageStatus;
}

interface Conversation {
  id: string;
  name: string;
  username: string;
  role: string;
  avatar: string;
  accent: string;
  lastActive: string;
  about: string;
  unreadCount: number;
  messages: Message[];
}

const storageKey = 'news-social-chat-demo';

const seedConversations: Conversation[] = [
  {
    id: 'maya-chen',
    name: 'Maya Chen',
    username: '@maya',
    role: 'Technology Desk',
    avatar: 'MC',
    accent: 'from-cyan-500 to-blue-500',
    lastActive: 'Online now',
    about: 'Watching AI chip policy, cloud earnings, and developer tools.',
    unreadCount: 2,
    messages: [
      {
        id: 'maya-1',
        sender: 'friend',
        text: 'The chip export story has changed three times today. Are you covering it?',
        sentAt: '9:18 AM',
      },
      {
        id: 'maya-2',
        sender: 'me',
        text: 'Yes. I am comparing the policy note against the market reaction before I post anything.',
        sentAt: '9:20 AM',
        status: 'Read',
      },
      {
        id: 'maya-3',
        sender: 'friend',
        text: 'Perfect. If you find a clean source, send it over. I want to reference it in my summary.',
        sentAt: '9:22 AM',
      },
    ],
  },
  {
    id: 'jordan-lee',
    name: 'Jordan Lee',
    username: '@jord',
    role: 'World News',
    avatar: 'JL',
    accent: 'from-orange-400 to-rose-500',
    lastActive: 'Active 12m ago',
    about: 'Focuses on geopolitics, elections, and emerging policy shifts.',
    unreadCount: 0,
    messages: [
      {
        id: 'jordan-1',
        sender: 'me',
        text: 'Do you think the headline is too broad for the push notification?',
        sentAt: 'Yesterday',
        status: 'Read',
      },
      {
        id: 'jordan-2',
        sender: 'friend',
        text: 'A little. The body is strong, but the opener needs one concrete detail.',
        sentAt: 'Yesterday',
      },
      {
        id: 'jordan-3',
        sender: 'friend',
        text: 'Try anchoring it on the treaty timeline instead of the reaction quotes.',
        sentAt: 'Yesterday',
      },
    ],
  },
  {
    id: 'nina-patel',
    name: 'Nina Patel',
    username: '@nina',
    role: 'Business & Markets',
    avatar: 'NP',
    accent: 'from-emerald-400 to-teal-600',
    lastActive: 'Active 1h ago',
    about: 'Tracks earnings, macro shifts, and how headlines move sentiment.',
    unreadCount: 1,
    messages: [
      {
        id: 'nina-1',
        sender: 'friend',
        text: 'I put together a quick market take. Want me to drop it into the thread?',
        sentAt: '8:04 AM',
      },
      {
        id: 'nina-2',
        sender: 'me',
        text: 'Yes. Keep it short and tie it back to rates so readers understand the move.',
        sentAt: '8:09 AM',
        status: 'Delivered',
      },
    ],
  },
];

const quickReplies = [
  'Send me the source link.',
  'I can turn that into a short post.',
  'Let us compare headlines before publishing.',
];

const readStoredConversations = (): Conversation[] => {
  if (typeof window === 'undefined') {
    return seedConversations;
  }

  const storedConversations = window.localStorage.getItem(storageKey);

  if (!storedConversations) {
    return seedConversations;
  }

  try {
    const parsedConversations = JSON.parse(storedConversations);

    if (Array.isArray(parsedConversations) && parsedConversations.length > 0) {
      return parsedConversations as Conversation[];
    }
  } catch (_error) {
    return seedConversations;
  }

  return seedConversations;
};

const createMessageId = () =>
  `message-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`;

const formatCurrentTime = () =>
  new Date().toLocaleTimeString([], {
    hour: 'numeric',
    minute: '2-digit',
  });

const Chat: React.FC = () => {
  const [conversations, setConversations] = useState<Conversation[]>(
    readStoredConversations
  );
  const [selectedId, setSelectedId] = useState<string>(
    () => readStoredConversations()[0]?.id ?? ''
  );
  const [searchTerm, setSearchTerm] = useState('');
  const [draftMessage, setDraftMessage] = useState('');
  const messageEndRef = useRef<HTMLDivElement>(null);

  const visibleConversations = conversations.filter((conversation) => {
    const normalizedSearch = searchTerm.trim().toLowerCase();

    if (!normalizedSearch) {
      return true;
    }

    return [
      conversation.name,
      conversation.username,
      conversation.role,
      conversation.about,
    ].some((value) => value.toLowerCase().includes(normalizedSearch));
  });

  useEffect(() => {
    if (typeof window !== 'undefined') {
      window.localStorage.setItem(storageKey, JSON.stringify(conversations));
    }
  }, [conversations]);

  useEffect(() => {
    if (!visibleConversations.some((conversation) => conversation.id === selectedId)) {
      setSelectedId(visibleConversations[0]?.id ?? conversations[0]?.id ?? '');
    }
  }, [conversations, selectedId, visibleConversations]);

  const activeConversation =
    conversations.find((conversation) => conversation.id === selectedId) ??
    conversations[0];

  useEffect(() => {
    messageEndRef.current?.scrollIntoView?.({
      behavior: 'smooth',
      block: 'end',
    });
  }, [activeConversation?.id, activeConversation?.messages.length]);

  const selectConversation = (conversationId: string) => {
    setSelectedId(conversationId);
    setConversations((previousConversations) =>
      previousConversations.map((conversation) =>
        conversation.id === conversationId
          ? { ...conversation, unreadCount: 0 }
          : conversation
      )
    );
  };

  const sendMessage = (text: string) => {
    if (!activeConversation) {
      return;
    }

    const trimmedMessage = text.trim();

    if (!trimmedMessage) {
      return;
    }

    const newMessage: Message = {
      id: createMessageId(),
      sender: 'me',
      text: trimmedMessage,
      sentAt: formatCurrentTime(),
      status: 'Sent',
    };

    setConversations((previousConversations) =>
      previousConversations.map((conversation) =>
        conversation.id === activeConversation.id
          ? {
              ...conversation,
              messages: [...conversation.messages, newMessage],
              unreadCount: 0,
            }
          : conversation
      )
    );
    setDraftMessage('');
  };

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    sendMessage(draftMessage);
  };

  const activeMessages = activeConversation?.messages ?? [];

  return (
    <main
      data-cy="chat-page"
      className="mx-auto flex min-h-[calc(100vh-56px)] max-w-6xl flex-col gap-4 px-4 py-5 sm:px-6 lg:flex-row lg:py-6"
    >
      <aside className="w-full shrink-0 overflow-hidden rounded-[28px] bg-white shadow-sm lg:w-[340px]">
        <div className="border-b border-slate-100 px-5 py-5">
          <div className="flex items-start justify-between gap-3">
            <div>
              <p className="text-xs font-semibold uppercase tracking-[0.25em] text-blue-600">
                Direct Messages
              </p>
              <h1 className="mt-2 text-2xl font-bold text-slate-900">
                Chat with your people
              </h1>
            </div>
            <div className="rounded-2xl bg-blue-50 p-3 text-blue-600">
              <ChatBubbleLeftRightIcon className="h-6 w-6" />
            </div>
          </div>
          <p className="mt-3 text-sm leading-6 text-slate-500">
            Frontend-only demo. Conversations are mocked and new messages stay in
            this browser.
          </p>

          <label className="mt-5 flex items-center gap-3 rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
            <MagnifyingGlassIcon className="h-5 w-5 text-slate-400" />
            <input
              type="text"
              value={searchTerm}
              onChange={(event) => setSearchTerm(event.target.value)}
              placeholder="Search chats"
              data-cy="chat-search"
              className="w-full border-none bg-transparent text-sm text-slate-700 outline-none placeholder:text-slate-400"
            />
          </label>
        </div>

        <div className="max-h-[320px] overflow-y-auto lg:max-h-[calc(100vh-220px)]">
          {visibleConversations.length === 0 ? (
            <div className="px-5 py-8 text-sm text-slate-500">
              No chats match that search yet.
            </div>
          ) : (
            visibleConversations.map((conversation) => {
              const lastMessage = conversation.messages[conversation.messages.length - 1];
              const isActive = conversation.id === activeConversation?.id;

              return (
                <button
                  key={conversation.id}
                  type="button"
                  onClick={() => selectConversation(conversation.id)}
                  data-cy={`conversation-${conversation.id}`}
                  className={`flex w-full items-start gap-3 border-b border-slate-100 px-5 py-4 text-left transition ${
                    isActive ? 'bg-slate-50' : 'hover:bg-slate-50/70'
                  }`}
                >
                  <div
                    className={`flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl bg-gradient-to-br ${conversation.accent} text-sm font-semibold text-white`}
                  >
                    {conversation.avatar}
                  </div>
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center justify-between gap-3">
                      <p className="truncate text-sm font-semibold text-slate-900">
                        {conversation.name}
                      </p>
                      <span className="shrink-0 text-xs text-slate-400">
                        {lastMessage?.sentAt}
                      </span>
                    </div>
                    <p className="mt-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                      {conversation.role}
                    </p>
                    <p className="mt-2 truncate text-sm text-slate-500">
                      {lastMessage?.sender === 'me' ? 'You: ' : ''}
                      {lastMessage?.text}
                    </p>
                  </div>
                  {conversation.unreadCount > 0 && (
                    <span className="mt-1 inline-flex h-6 min-w-[1.5rem] items-center justify-center rounded-full bg-blue-600 px-2 text-xs font-semibold text-white">
                      {conversation.unreadCount}
                    </span>
                  )}
                </button>
              );
            })
          )}
        </div>
      </aside>

      <section className="flex min-h-[65vh] flex-1 flex-col overflow-hidden rounded-[28px] bg-white shadow-sm">
        {activeConversation ? (
          <>
            <header className="border-b border-slate-100 px-5 py-5 sm:px-6">
              <div className="flex items-start justify-between gap-4">
                <div className="flex items-center gap-4">
                  <div
                    className={`flex h-14 w-14 items-center justify-center rounded-2xl bg-gradient-to-br ${activeConversation.accent} text-base font-semibold text-white`}
                  >
                    {activeConversation.avatar}
                  </div>
                  <div>
                    <div className="flex items-center gap-3">
                      <h2 className="text-xl font-bold text-slate-900">
                        {activeConversation.name}
                      </h2>
                      <span className="rounded-full bg-emerald-50 px-3 py-1 text-xs font-semibold text-emerald-700">
                        {activeConversation.lastActive}
                      </span>
                    </div>
                    <p className="mt-1 text-sm text-slate-500">
                      {activeConversation.username} | {activeConversation.role}
                    </p>
                    <p className="mt-2 max-w-2xl text-sm leading-6 text-slate-500">
                      {activeConversation.about}
                    </p>
                  </div>
                </div>
                <button
                  type="button"
                  className="rounded-2xl border border-slate-200 p-3 text-slate-400 transition hover:border-slate-300 hover:text-slate-600"
                >
                  <EllipsisHorizontalIcon className="h-6 w-6" />
                </button>
              </div>
            </header>

            <div
              data-cy="chat-messages"
              className="flex-1 overflow-y-auto bg-slate-50/70 px-4 py-6 sm:px-6"
            >
              <div className="mx-auto flex max-w-3xl flex-col gap-3">
                {activeMessages.map((message) => {
                  const isOwnMessage = message.sender === 'me';

                  return (
                    <div
                      key={message.id}
                      className={`flex ${
                        isOwnMessage ? 'justify-end' : 'justify-start'
                      }`}
                    >
                      <div
                        className={`max-w-[85%] rounded-[24px] px-4 py-3 shadow-sm sm:max-w-[70%] ${
                          isOwnMessage
                            ? 'rounded-br-md bg-blue-600 text-white'
                            : 'rounded-bl-md bg-white text-slate-800'
                        }`}
                      >
                        <p className="text-sm leading-6">{message.text}</p>
                        <div
                          className={`mt-2 flex items-center justify-end gap-2 text-xs ${
                            isOwnMessage ? 'text-blue-100' : 'text-slate-400'
                          }`}
                        >
                          <span>{message.sentAt}</span>
                          {message.status && <span>{message.status}</span>}
                        </div>
                      </div>
                    </div>
                  );
                })}
                <div ref={messageEndRef} />
              </div>
            </div>

            <div className="border-t border-slate-100 px-4 py-4 sm:px-6">
              <div className="mb-3 flex flex-wrap gap-2">
                {quickReplies.map((reply) => (
                  <button
                    key={reply}
                    type="button"
                    onClick={() => setDraftMessage(reply)}
                    className="rounded-full border border-slate-200 px-3 py-2 text-xs font-medium text-slate-600 transition hover:border-slate-300 hover:bg-slate-50"
                  >
                    {reply}
                  </button>
                ))}
              </div>

              <form
                onSubmit={handleSubmit}
                className="flex flex-col gap-3 rounded-[24px] border border-slate-200 bg-white p-3 sm:flex-row sm:items-end"
              >
                <label className="flex-1">
                  <span className="sr-only">Type your message</span>
                  <textarea
                    rows={2}
                    value={draftMessage}
                    onChange={(event) => setDraftMessage(event.target.value)}
                    placeholder="Type a message..."
                    data-cy="chat-draft"
                    className="w-full resize-none border-none bg-transparent px-1 py-2 text-sm leading-6 text-slate-700 outline-none placeholder:text-slate-400"
                  />
                </label>
                <button
                  type="submit"
                  disabled={!draftMessage.trim()}
                  data-cy="chat-send"
                  className="inline-flex items-center justify-center gap-2 rounded-2xl bg-blue-600 px-4 py-3 text-sm font-semibold text-white transition hover:bg-blue-700 disabled:cursor-not-allowed disabled:bg-slate-300"
                >
                  <PaperAirplaneIcon className="h-5 w-5" />
                  Send
                </button>
              </form>
            </div>
          </>
        ) : (
          <div className="flex flex-1 items-center justify-center px-6 text-sm text-slate-500">
            Pick a conversation to start chatting.
          </div>
        )}
      </section>
    </main>
  );
};

export default Chat;
