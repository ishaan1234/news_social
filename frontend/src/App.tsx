import React, { useEffect, useState } from 'react';
import Navbar from './components/Navbar';
import Home from './pages/Home';
import Chat from './pages/Chat';
import Posts from './pages/Posts';
import PlaceholderPage from './pages/PlaceholderPage';

export type AppRoute = '/' | '/posts' | '/chat' | '/profile' | '/settings';

const validRoutes: AppRoute[] = ['/', '/posts', '/chat', '/profile', '/settings'];
const validRouteSet = new Set<AppRoute>(validRoutes);

const getCurrentRoute = (): AppRoute => {
  if (typeof window === 'undefined') {
    return '/';
  }

  const candidate = window.location.hash
    ? window.location.hash.replace(/^#/, '') || '/'
    : window.location.pathname || '/';

  return validRouteSet.has(candidate as AppRoute) ? (candidate as AppRoute) : '/';
};

const renderRoute = (route: AppRoute) => {
  switch (route) {
    case '/chat':
      return <Chat />;
    case '/posts':
      return <Posts />;
    case '/profile':
      return (
        <PlaceholderPage
          title="Profile is still a placeholder."
          description="Chat is available now. Profile data and editing can plug into the same frontend routing pattern later."
        />
      );
    case '/settings':
      return (
        <PlaceholderPage
          title="Settings are still pending."
          description="This screen is intentionally static for now so the frontend can focus on the new chat flow."
        />
      );
    case '/':
    default:
      return <Home />;
  }
};

function App() {
  const [currentRoute, setCurrentRoute] = useState<AppRoute>(getCurrentRoute);

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
      <Navbar currentPath={currentRoute} />
      {renderRoute(currentRoute)}
    </div>
  );
}

export default App;
