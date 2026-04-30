import React, { useState } from 'react';
import {
  MagnifyingGlassIcon,
  BellIcon,
  Bars3Icon,
  XMarkIcon,
} from '@heroicons/react/24/outline';
import {
  AuthSession,
  getInitials,
  getSessionDisplayName,
  isVerifiedAuthSession,
} from '../auth';

const navLinks = [
  { name: 'News', path: '/', href: '#/' },
  { name: 'Posts', path: '/posts', href: '#/posts' },
  { name: 'Profile', path: '/profile', href: '#/profile' },
  { name: 'Settings', path: '/settings', href: '#/settings' },
];

interface NavbarProps {
  currentPath: string;
  authSession: AuthSession | null;
  onSignOut: () => void;
}

const Navbar: React.FC<NavbarProps> = ({
  currentPath,
  authSession,
  onSignOut,
}) => {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const sessionName = getSessionDisplayName(authSession, 'NewsHub User');
  const sessionInitials = getInitials(sessionName);
  const hasVerifiedSession = isVerifiedAuthSession(authSession);

  return (
    <nav className="bg-white shadow-sm sticky top-0 z-50">
      <div className="max-w-6xl mx-auto px-4 sm:px-6">
        <div className="flex items-center justify-between h-14">

          {/* Logo */}
          <div className="flex items-center gap-2">
            <a href="#/" className="flex items-center gap-2">
              <div className="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center">
                <span className="text-white font-bold text-sm">N</span>
              </div>
              <span className="text-xl font-bold text-gray-900">
                NewsHub
              </span>
            </a>
          </div>

          {/* Desktop Links */}
          <div className="hidden md:flex items-center gap-6">
            {navLinks.map((link) => {
              const isActive = link.path === currentPath;

              return (
                <a
                  key={link.name}
                  href={link.href}
                  data-cy={`nav-${link.name.toLowerCase()}`}
                  aria-current={isActive ? 'page' : undefined}
                  className={`text-sm transition-colors ${isActive
                    ? 'text-blue-600 font-medium'
                    : 'text-gray-500 hover:text-gray-900'
                    }`}
                >
                  {link.name}
                </a>
              );
            })}
          </div>

          {/* Right Actions */}
          <div className="flex items-center gap-2">

            {/* Search */}
            <div className="hidden sm:flex items-center bg-gray-100 rounded-full px-3 py-1.5">
              <MagnifyingGlassIcon className="w-4 h-4 text-gray-400" />
              <input
                type="text"
                placeholder="Search..."
                className="bg-transparent border-none outline-none ml-2 text-sm text-gray-700 w-36"
              />
            </div>

            {/* Notifications */}
            <button className="p-2 text-gray-500 hover:text-gray-900 relative">
              <BellIcon className="w-5 h-5" />
              <span className="absolute top-1.5 right-1.5 w-2 h-2 bg-red-500 rounded-full" />
            </button>

            {hasVerifiedSession ? (
              <>
                <a
                  href="#/profile"
                  data-cy="nav-account"
                  className="hidden sm:flex items-center gap-3 rounded-full border border-slate-200 px-2 py-1.5 transition hover:border-slate-300 hover:bg-slate-50"
                >
                  <span className="flex h-8 w-8 items-center justify-center rounded-full bg-slate-900 text-xs font-semibold text-white">
                    {sessionInitials}
                  </span>
                  <span className="max-w-[140px] truncate pr-1 text-sm font-medium text-slate-700">
                    {sessionName}
                  </span>
                </a>
                <button
                  type="button"
                  onClick={onSignOut}
                  data-cy="nav-signout"
                  className="hidden rounded-full border border-slate-200 px-4 py-2 text-sm font-semibold text-slate-600 transition hover:border-slate-300 hover:bg-slate-50 sm:inline-flex"
                >
                  Sign out
                </button>
              </>
            ) : (
              <a
                href="#/auth"
                data-cy="nav-auth"
                aria-current={currentPath === '/auth' ? 'page' : undefined}
                className={`hidden rounded-full px-4 py-2 text-sm font-semibold transition sm:inline-flex ${
                  currentPath === '/auth'
                    ? 'bg-slate-900 text-white'
                    : 'border border-slate-200 text-slate-700 hover:border-slate-300 hover:bg-slate-50'
                }`}
              >
                Sign in
              </a>
            )}

            <button
              className="md:hidden p-2 text-gray-500"
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
            >
              {mobileMenuOpen
                ? <XMarkIcon className="w-5 h-5" />
                : <Bars3Icon className="w-5 h-5" />
              }
            </button>

          </div>
        </div>
      </div>

      {mobileMenuOpen && (
        <div className="md:hidden border-t border-gray-100 bg-white px-4 py-3 space-y-2">
          {navLinks.map((link) => {
            const isActive = link.path === currentPath;

            return (
              <a
                key={link.name}
                href={link.href}
                data-cy={`nav-${link.name.toLowerCase()}`}
                className={`block text-sm py-1.5 ${isActive
                  ? 'text-blue-600 font-medium'
                  : 'text-gray-500 hover:text-gray-900'
                  }`}
                onClick={() => setMobileMenuOpen(false)}
              >
                {link.name}
              </a>
            );
          })}

          {hasVerifiedSession ? (
            <>
              <a
                href="#/profile"
                data-cy="nav-account"
                className="block text-sm py-1.5 text-gray-500 hover:text-gray-900"
                onClick={() => setMobileMenuOpen(false)}
              >
                Account
              </a>
              <button
                type="button"
                data-cy="nav-signout"
                onClick={() => {
                  onSignOut();
                  setMobileMenuOpen(false);
                }}
                className="block w-full text-left text-sm py-1.5 text-gray-500 hover:text-gray-900"
              >
                Sign out
              </button>
            </>
          ) : (
            <a
              href="#/auth"
              data-cy="nav-auth"
              className={`block text-sm py-1.5 ${
                currentPath === '/auth'
                  ? 'text-blue-600 font-medium'
                  : 'text-gray-500 hover:text-gray-900'
              }`}
              onClick={() => setMobileMenuOpen(false)}
            >
              Sign in
            </a>
          )}
        </div>
      )}
    </nav>
  );
};

export default Navbar;
