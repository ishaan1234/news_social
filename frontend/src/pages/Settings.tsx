import React from 'react';
import { AuthSession } from '../auth';

interface SettingsProps {
  authSession?: AuthSession | null;
}

const Settings: React.FC<SettingsProps> = ({ authSession = null }) => {
  const handleDeleteAccount = () => {
    if (window.confirm("Are you sure you want to delete your account? This action cannot be undone.")) {
      alert("Account deletion feature is coming soon!");
      // Here you could make an API call to DELETE /profile or DELETE /account
    }
  };

  return (
    <main
      data-cy="settings-page"
      className="mx-auto min-h-[calc(100vh-56px)] max-w-2xl px-4 py-6 sm:px-6 lg:py-8"
    >
      <section className="rounded-[32px] bg-white p-6 shadow-sm sm:p-7">
        <h1 className="text-3xl font-bold tracking-tight text-slate-900">
          Settings
        </h1>
        <p className="mt-3 text-sm leading-6 text-slate-500">
          Manage your account preferences here.
        </p>
      </section>

      <section className="mt-6 rounded-[32px] bg-white p-6 shadow-sm sm:p-7">
        <h2 className="text-xl font-bold text-red-600">Danger Zone</h2>
        <p className="mt-2 text-sm text-slate-500">
          Permanently delete your account and all associated data.
        </p>
        <button
          onClick={handleDeleteAccount}
          className="mt-4 rounded-full bg-red-100 px-4 py-2 text-sm font-semibold text-red-600 hover:bg-red-200 transition"
        >
          Delete Account
        </button>
      </section>
    </main>
  );
};

export default Settings;
