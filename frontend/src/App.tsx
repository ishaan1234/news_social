import React, { useEffect, useState } from 'react';
import {
  AuthSession,
  clearStoredAuthSession,
  readStoredAuthSession,
  storeAuthSession,
} from './auth';
import Navbar from './components/Navbar';
import Auth from './pages/Auth';
import Home from './pages/Home';
import Chat from './pages/Chat';
import Posts from './pages/Posts';
import Profile from './pages/Profile';
import Settings from './pages/Settings';

export type AppRoute =
  | '/'
  | '/posts'
  | '/chat'
  | '/profile'
  | '/settings'
  | '/auth';

const validRoutes: AppRoute[] = [
  '/',
  '/posts',
  '/chat',
  '/profile',
  '/settings',
  '/auth',
];
const validRouteSet = new Set<AppRoute>(validRoutes);

const getCurrentRoute = (): AppRoute => {
  if (typeof window === 'undefined') {
    return '/';
  }

  const candidate = window.location.hash
    ? window.location.hash.replace(/^#/, '') || '/'
    : window.location.pathname || '/';
  const routeOnly = candidate.split('?')[0] || '/';

  return validRouteSet.has(routeOnly as AppRoute)
    ? (routeOnly as AppRoute)
    : '/';
};

const renderRoute = (
  route: AppRoute,
  authSession: AuthSession | null,
  handleAuthSuccess: (session: AuthSession) => void,
  handleSignOut: () => void
) => {
  switch (route) {
    case '/chat':
      return <Chat />;
    case '/posts':
      return <Posts authSession={authSession} />;
    case '/profile':
      return <Profile authSession={authSession} />;
    case '/auth':
      return (
        <Auth
          authSession={authSession}
          onAuthSuccess={handleAuthSuccess}
          onSignOut={handleSignOut}
        />
      );
    case '/settings':
      return <Settings authSession={authSession} />;
    case '/':
    default:
      return <Home />;
  }
};

function App() {
  const [currentRoute, setCurrentRoute] = useState<AppRoute>(getCurrentRoute);
  const [authSession, setAuthSession] = useState<AuthSession | null>(
    readStoredAuthSession
  );

  const handleAuthSuccess = (session: AuthSession) => {
    setAuthSession(session);
    storeAuthSession(session);
  };

  const handleSignOut = () => {
    clearStoredAuthSession();
    setAuthSession(null);
  };

  useEffect(() => {
    const syncRoute = () => {
      setCurrentRoute(getCurrentRoute());
    };

    syncRoute();
    window.addEventListener('hashchange', syncRoute);
    window.addEventListener('popstate', syncRoute);

    return () => {
      window.removeEventListener('hashchange', syncRoute);
      window.removeEventListener('popstate', syncRoute);
    };
  }, []);

  return (
    <div className="min-h-screen bg-gray-100">
      <Navbar
        currentPath={currentRoute}
        authSession={authSession}
        onSignOut={handleSignOut}
      />
      {renderRoute(currentRoute, authSession, handleAuthSuccess, handleSignOut)}
    </div>
  );
}

export default App;
