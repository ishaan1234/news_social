import React from 'react';

interface PlaceholderPageProps {
  title: string;
  description: string;
}

const PlaceholderPage: React.FC<PlaceholderPageProps> = ({
  title,
  description,
}) => {
  return (
    <main className="mx-auto flex min-h-[calc(100vh-56px)] max-w-4xl items-center px-4 py-10 sm:px-6">
      <section className="w-full rounded-[32px] bg-white p-8 shadow-sm sm:p-12">
        <span className="inline-flex rounded-full bg-slate-100 px-3 py-1 text-xs font-semibold uppercase tracking-[0.2em] text-slate-500">
          Frontend Placeholder
        </span>
        <h1 className="mt-6 text-3xl font-bold text-slate-900 sm:text-4xl">
          {title}
        </h1>
        <p className="mt-4 max-w-2xl text-base leading-7 text-slate-600">
          {description}
        </p>
      </section>
    </main>
  );
};

export default PlaceholderPage;
