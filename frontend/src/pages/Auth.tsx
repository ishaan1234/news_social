import React, { FormEvent, useEffect, useState } from 'react';
import {
  ArrowPathIcon,
  CheckCircleIcon,
  EnvelopeIcon,
  KeyIcon,
  UserCircleIcon,
} from '@heroicons/react/24/outline';
import {
  AuthApiResponse,
  AuthSession,
  buildAuthSession,
  getSessionDisplayName,
  isVerifiedAuthSession,
} from '../auth';

type AuthMode = 'signup' | 'login';
type FeedbackTone = 'success' | 'error' | 'info';

interface AuthPageProps {
  authSession: AuthSession | null;
  onAuthSuccess: (session: AuthSession) => void;
  onSignOut: () => void;
}

interface AuthFeedback {
  tone: FeedbackTone;
  message: string;
}

const apiBaseUrl = (process.env.REACT_APP_API_BASE_URL || '').replace(/\/$/, '');

const modeCopy: Record<
  AuthMode,
  {
    title: string;
    description: string;
    submitLabel: string;
    endpoint: string;
  }
> = {
  signup: {
    title: 'Create account',
    description:
      'Set up your account and we will send a verification email to finish activation.',
    submitLabel: 'Create account',
    endpoint: '/auth/signup',
  },
  login: {
    title: 'Log in',
    description: 'Use your email and password to access your account.',
    submitLabel: 'Log in',
    endpoint: '/auth/login',
  },
};

const resendVerificationEndpoint = '/auth/verify-email/resend';

const requiresEmailVerification = (message: string) =>
  /not verified|verify your email|verify it first|verification email/i.test(
    message
  );

const maskEmail = (value: string) => {
  const email = value.trim();
  if (!email.includes('@')) {
    return email;
  }

  const [localPart, domain] = email.split('@');
  if (localPart.length <= 2) {
    return `${localPart[0] || ''}*@${domain}`;
  }

  return `${localPart.slice(0, 2)}${'*'.repeat(
    Math.max(localPart.length - 2, 1)
  )}@${domain}`;
};

const postAuthRequest = async (path: string, payload: Record<string, string>) => {
  const response = await fetch(`${apiBaseUrl}${path}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(payload),
  });

  const responseBody = (await response.json().catch(() => null)) as
    | (AuthApiResponse & { error?: string })
    | null;

  if (!response.ok) {
    throw new Error(
      responseBody?.error ||
        responseBody?.message ||
        `request failed with status ${response.status}`
    );
  }

  return responseBody as AuthApiResponse;
};

const Auth: React.FC<AuthPageProps> = ({
  authSession,
  onAuthSuccess,
  onSignOut,
}) => {
  const hasVerifiedSession = isVerifiedAuthSession(authSession);
  const [mode, setMode] = useState<AuthMode>(hasVerifiedSession ? 'login' : 'signup');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isResendingVerification, setIsResendingVerification] = useState(false);
  const [feedback, setFeedback] = useState<AuthFeedback | null>(null);
  const [showVerificationHelp, setShowVerificationHelp] = useState(
    Boolean(authSession && !authSession.user?.email_verified)
  );

  const [signupName, setSignupName] = useState('');
  const [signupEmail, setSignupEmail] = useState('');
  const [signupPassword, setSignupPassword] = useState('');

  const [loginEmail, setLoginEmail] = useState('');
  const [loginPassword, setLoginPassword] = useState('');
  const [verificationEmail, setVerificationEmail] = useState('');
  const [verificationPassword, setVerificationPassword] = useState('');

  useEffect(() => {
    const sessionEmail = authSession?.user?.email?.trim() || '';
    if (!sessionEmail) {
      return;
    }

    setLoginEmail((previousEmail) => previousEmail || sessionEmail);
    setSignupEmail((previousEmail) => previousEmail || sessionEmail);
    setVerificationEmail((previousEmail) => previousEmail || sessionEmail);
  }, [authSession?.user?.email]);

  useEffect(() => {
    if (!authSession) {
      return;
    }

    if (!authSession.user?.email_verified) {
      setShowVerificationHelp(true);
    }

    setFeedback({
      tone: authSession.user?.email_verified ? 'success' : 'info',
      message:
        authSession.message ||
        (authSession.user?.email_verified
          ? 'You are logged in.'
          : 'Your account is almost ready. Verify your email to finish signing in.'),
    });
  }, [authSession]);

  const handleReturnToLogin = () => {
    setShowVerificationHelp(false);
    setMode('login');
    setFeedback(null);
  };

  const handleUseDifferentEmail = () => {
    onSignOut();
    setShowVerificationHelp(false);
    setMode('signup');
    setFeedback(null);
    setSignupName('');
    setSignupEmail('');
    setSignupPassword('');
    setLoginEmail('');
    setLoginPassword('');
    setVerificationEmail('');
    setVerificationPassword('');
  };

  const handleResendVerification = async () => {
    setIsResendingVerification(true);
    setFeedback(null);

    try {
      const response = await postAuthRequest(resendVerificationEndpoint, {
        email: verificationEmail,
        password: verificationPassword,
      });

      setFeedback({
        tone: 'success',
        message: response.message || 'Verification email sent.',
      });
    } catch (error) {
      setFeedback({
        tone: 'error',
        message:
          error instanceof Error
            ? error.message
            : 'Unable to resend verification email.',
      });
    } finally {
      setIsResendingVerification(false);
    }
  };

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setIsSubmitting(true);
    setFeedback(null);

    try {
      let response: AuthApiResponse;
      const activeEmail =
        mode === 'signup' ? signupEmail.trim() : loginEmail.trim();
      const activePassword = mode === 'signup' ? signupPassword : loginPassword;

      switch (mode) {
        case 'signup':
          response = await postAuthRequest(modeCopy.signup.endpoint, {
            email: signupEmail,
            password: signupPassword,
            display_name: signupName,
          });
          break;
        case 'login':
          response = await postAuthRequest(modeCopy.login.endpoint, {
            email: loginEmail,
            password: loginPassword,
          });
          break;
        default:
          response = {
            success: false,
            message: 'unsupported auth mode',
          };
      }

      const session = buildAuthSession(response);
      onAuthSuccess(session);

      if (!session.user?.email_verified) {
        setVerificationEmail(activeEmail);
        setVerificationPassword(activePassword);
        setShowVerificationHelp(true);
        setFeedback({
          tone: 'info',
          message:
            response.message ||
            'Check your inbox and verify your email before logging in.',
        });
        return;
      }

      setShowVerificationHelp(false);
      setVerificationPassword('');
      setFeedback({
        tone: 'success',
        message: response.message || 'Authentication successful.',
      });

      if (typeof window !== 'undefined') {
        window.location.hash = '#/profile';
      }
    } catch (error) {
      const message =
        error instanceof Error ? error.message : 'Authentication failed.';
      const needsVerification = requiresEmailVerification(message);

      if (needsVerification) {
        setVerificationEmail(
          mode === 'signup' ? signupEmail.trim() : loginEmail.trim()
        );
        setVerificationPassword(
          mode === 'signup' ? signupPassword : loginPassword
        );
        setShowVerificationHelp(true);
      }

      setFeedback({
        tone: needsVerification ? 'info' : 'error',
        message,
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  const verificationDestination = verificationEmail || authSession?.user?.email || '';
  const isVerificationStep = showVerificationHelp && !hasVerifiedSession;

  return (
    <main
      data-cy="auth-page"
      className="mx-auto min-h-[calc(100vh-56px)] max-w-xl px-4 py-6 sm:px-6 lg:py-8"
    >
      <section className="rounded-[32px] bg-white p-6 shadow-sm sm:p-7">
        <h1 className="text-3xl font-bold tracking-tight text-slate-900">
          {isVerificationStep ? 'Verify your email' : 'Account access'}
        </h1>
        <p className="mt-3 text-sm leading-6 text-slate-500">
          {isVerificationStep
            ? verificationDestination
              ? `We sent a verification link to ${maskEmail(
                  verificationDestination
                )}. Open it to activate the account, or resend it below.`
              : 'Open the verification link we sent to your inbox, or resend it below.'
            : 'Sign up or log in to manage your business account.'}
        </p>

        {!isVerificationStep && (
          <>
            <div className="mt-5 flex flex-wrap gap-2">
              {(['signup', 'login'] as AuthMode[]).map((currentMode) => {
                const isActive = currentMode === mode;

                return (
                  <button
                    key={currentMode}
                    type="button"
                    onClick={() => {
                      setMode(currentMode);
                      setFeedback(null);
                      setShowVerificationHelp(false);
                    }}
                    data-cy={`auth-mode-${currentMode}`}
                    className={`rounded-full px-4 py-2 text-sm font-semibold transition ${
                      isActive
                        ? 'bg-slate-900 text-white'
                        : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
                    }`}
                  >
                    {currentMode === 'signup' ? 'Sign up' : 'Log in'}
                  </button>
                );
              })}
            </div>

            <h2 className="mt-5 text-2xl font-bold text-slate-900">
              {modeCopy[mode].title}
            </h2>
            <p className="mt-2 text-sm leading-6 text-slate-500">
              {modeCopy[mode].description}
            </p>
          </>
        )}

        {feedback && (
          <div
            data-cy="auth-feedback"
            className={`mt-5 rounded-[24px] px-4 py-4 text-sm leading-6 ${
              feedback.tone === 'success'
                ? 'bg-emerald-50 text-emerald-800'
                : feedback.tone === 'error'
                ? 'bg-rose-50 text-rose-700'
                : 'bg-blue-50 text-blue-800'
            }`}
          >
            {feedback.message}
          </div>
        )}

        {!isVerificationStep && (
          <form onSubmit={handleSubmit} className="mt-6 space-y-4">
            {mode === 'signup' && (
              <label className="block">
                <span className="text-sm font-medium text-slate-700">
                  Display name
                </span>
                <div className="mt-2 flex items-center gap-3 rounded-[22px] border border-slate-200 px-4 py-3">
                  <UserCircleIcon className="h-5 w-5 text-slate-400" />
                  <input
                    type="text"
                    value={signupName}
                    onChange={(event) => setSignupName(event.target.value)}
                    data-cy="auth-signup-name"
                    placeholder="Avery Stone"
                    className="w-full border-none bg-transparent text-sm text-slate-700 outline-none placeholder:text-slate-400"
                  />
                </div>
              </label>
            )}

            <label className="block">
              <span className="text-sm font-medium text-slate-700">Email</span>
              <div className="mt-2 flex items-center gap-3 rounded-[22px] border border-slate-200 px-4 py-3">
                <EnvelopeIcon className="h-5 w-5 text-slate-400" />
                <input
                  type="email"
                  value={mode === 'signup' ? signupEmail : loginEmail}
                  onChange={(event) => {
                    if (mode === 'signup') {
                      setSignupEmail(event.target.value);
                    } else {
                      setLoginEmail(event.target.value);
                    }
                  }}
                  data-cy="auth-email"
                  placeholder="you@example.com"
                  className="w-full border-none bg-transparent text-sm text-slate-700 outline-none placeholder:text-slate-400"
                />
              </div>
            </label>

            <label className="block">
              <span className="text-sm font-medium text-slate-700">Password</span>
              <div className="mt-2 flex items-center gap-3 rounded-[22px] border border-slate-200 px-4 py-3">
                <KeyIcon className="h-5 w-5 text-slate-400" />
                <input
                  type="password"
                  value={mode === 'signup' ? signupPassword : loginPassword}
                  onChange={(event) => {
                    if (mode === 'signup') {
                      setSignupPassword(event.target.value);
                    } else {
                      setLoginPassword(event.target.value);
                    }
                  }}
                  data-cy="auth-password"
                  placeholder="Minimum 6 characters"
                  className="w-full border-none bg-transparent text-sm text-slate-700 outline-none placeholder:text-slate-400"
                />
              </div>
            </label>

            <button
              type="submit"
              disabled={isSubmitting}
              data-cy="auth-submit"
              className="inline-flex w-full items-center justify-center gap-2 rounded-[22px] bg-blue-600 px-5 py-3 text-sm font-semibold text-white transition hover:bg-blue-700 disabled:cursor-not-allowed disabled:bg-slate-300"
            >
              {isSubmitting ? (
                <ArrowPathIcon className="h-5 w-5 animate-spin" />
              ) : mode === 'login' ? (
                <CheckCircleIcon className="h-5 w-5" />
              ) : (
                <EnvelopeIcon className="h-5 w-5" />
              )}
              {isSubmitting ? 'Working...' : modeCopy[mode].submitLabel}
            </button>
          </form>
        )}

        {isVerificationStep && (
          <section
            data-cy="auth-verification-help"
            className="mt-6 rounded-[24px] border border-blue-100 bg-blue-50 p-5"
          >
            <p className="text-xs font-semibold uppercase tracking-[0.16em] text-blue-500">
              Email verification
            </p>
            <h3 className="mt-3 text-lg font-bold text-slate-900">
              Need another email?
            </h3>
            <p className="mt-2 text-sm leading-6 text-slate-600">
              If the first message did not arrive, resend it here with the same
              email and password. Check spam or promotions first, then try again.
            </p>

            <div className="mt-4 space-y-4">
              <label className="block">
                <span className="text-sm font-medium text-slate-700">Email</span>
                <div className="mt-2 flex items-center gap-3 rounded-[22px] border border-blue-100 bg-white px-4 py-3">
                  <EnvelopeIcon className="h-5 w-5 text-slate-400" />
                  <input
                    type="email"
                    value={verificationEmail}
                    onChange={(event) => setVerificationEmail(event.target.value)}
                    data-cy="auth-resend-email"
                    placeholder="you@example.com"
                    className="w-full border-none bg-transparent text-sm text-slate-700 outline-none placeholder:text-slate-400"
                  />
                </div>
              </label>

              <label className="block">
                <span className="text-sm font-medium text-slate-700">
                  Password
                </span>
                <div className="mt-2 flex items-center gap-3 rounded-[22px] border border-blue-100 bg-white px-4 py-3">
                  <KeyIcon className="h-5 w-5 text-slate-400" />
                  <input
                    type="password"
                    value={verificationPassword}
                    onChange={(event) =>
                      setVerificationPassword(event.target.value)
                    }
                    data-cy="auth-resend-password"
                    placeholder="Password for this account"
                    className="w-full border-none bg-transparent text-sm text-slate-700 outline-none placeholder:text-slate-400"
                  />
                </div>
              </label>
            </div>

            <button
              type="button"
              onClick={handleResendVerification}
              disabled={isResendingVerification}
              data-cy="auth-resend-submit"
              className="mt-4 inline-flex w-full items-center justify-center gap-2 rounded-[22px] bg-slate-900 px-5 py-3 text-sm font-semibold text-white transition hover:bg-slate-800 disabled:cursor-not-allowed disabled:bg-slate-300"
            >
              {isResendingVerification ? (
                <ArrowPathIcon className="h-5 w-5 animate-spin" />
              ) : (
                <EnvelopeIcon className="h-5 w-5" />
              )}
              {isResendingVerification ? 'Sending...' : 'Resend verification email'}
            </button>

            <div className="mt-4 flex flex-wrap gap-3">
              <button
                type="button"
                onClick={handleReturnToLogin}
                data-cy="auth-back-to-login"
                className="rounded-full border border-slate-200 bg-white px-4 py-2 text-sm font-semibold text-slate-700 transition hover:border-slate-300 hover:bg-slate-50"
              >
                I verified, go to login
              </button>
              <button
                type="button"
                onClick={handleUseDifferentEmail}
                data-cy="auth-start-over"
                className="rounded-full border border-transparent px-4 py-2 text-sm font-semibold text-slate-500 transition hover:bg-white"
              >
                Use a different email
              </button>
            </div>
          </section>
        )}

        {authSession && hasVerifiedSession && (
          <section className="mt-6 rounded-[24px] border border-slate-200 bg-slate-50 p-4">
            <p className="text-xs font-semibold uppercase tracking-[0.16em] text-slate-400">
              Current session
            </p>
            <h3 className="mt-3 text-lg font-bold text-slate-900">
              {getSessionDisplayName(authSession)}
            </h3>
            <p className="mt-1 text-sm text-slate-500">
              {authSession.user?.email || 'No email returned'}
            </p>
            <p className="mt-3 text-sm leading-6 text-slate-600">
              {authSession.user?.email_verified
                ? 'This account is verified and ready to use.'
                : 'Email verification is still pending. You can resend the email above if needed.'}
            </p>

            <div className="mt-4 flex flex-wrap gap-3">
              <a
                href="#/profile"
                className="rounded-full bg-slate-900 px-4 py-2 text-sm font-semibold text-white transition hover:bg-slate-800"
              >
                Go to profile
              </a>
              <button
                type="button"
                onClick={onSignOut}
                data-cy="auth-signout"
                className="rounded-full border border-slate-200 px-4 py-2 text-sm font-semibold text-slate-600 transition hover:bg-white"
              >
                Sign out
              </button>
            </div>
          </section>
        )}
      </section>
    </main>
  );
};

export default Auth;
