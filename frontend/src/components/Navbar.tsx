import React, { useState } from 'react';
import {
  MagnifyingGlassIcon,
  BellIcon,
  Bars3Icon,
  XMarkIcon,
} from '@heroicons/react/24/outline';

const navLinks = [
  { name: 'News', path: '/' },
  { name: 'Posts', path: '/posts' },
  { name: 'Chat', path: '/chat' },
  { name: 'Profile', path: '/profile' },
  { name: 'Settings', path: '/settings' },
];

const Navbar: React.FC = () => {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  return (
    <nav className="bg-white shadow-sm sticky top-0 z-50">
      <div className="max-w-6xl mx-auto px-4 sm:px-6">
        <div className="flex items-center justify-between h-14">

          {/* Logo */}
          <div className="flex items-center gap-2">
            <a href="/" className="flex items-center gap-2">
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
            {navLinks.map((link) => (
              <a
                key={link.name}
                href={link.path}
                className={`text-sm ${link.name === 'News'
                  ? 'text-blue-600 font-medium'
                  : 'text-gray-500 hover:text-gray-900'
                  }`}
              >
                {link.name}
              </a>
            ))}
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
          {navLinks.map((link) => (
            <a
              key={link.name}
              href={link.path}
              className={`block text-sm py-1.5 ${link.name === 'News'
                ? 'text-blue-600 font-medium'
                : 'text-gray-500 hover:text-gray-900'
                }`}
              onClick={() => setMobileMenuOpen(false)}
            >
              {link.name}
            </a>
          ))}
        </div>
      )}
    </nav>
  );
};

export default Navbar;