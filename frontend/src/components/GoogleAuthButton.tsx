import React, { useEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { AlertCircle, Loader2 } from 'lucide-react';
import { CredentialResponse, GoogleLogin, GoogleOAuthProvider } from '@react-oauth/google';
import { useAuth } from '../hooks/useAuth';

const GoogleMark = () => (
  <svg viewBox="0 0 24 24" className="h-[18px] w-[18px]" aria-hidden="true">
    <path fill="#4285F4" d="M21.35 12.22c0-.71-.06-1.4-.18-2.05H12v3.88h5.24a4.48 4.48 0 0 1-1.94 2.94v2.52h3.14c1.84-1.69 2.91-4.19 2.91-7.29Z" />
    <path fill="#34A853" d="M12 21.72c2.62 0 4.82-.87 6.43-2.35l-3.14-2.52c-.87.58-1.99.92-3.29.92-2.53 0-4.67-1.71-5.44-4v2.6H3.32a9.72 9.72 0 0 0 8.68 5.35Z" />
    <path fill="#FBBC05" d="M6.56 13.77a5.84 5.84 0 0 1 0-3.54v-2.6H3.32a9.72 9.72 0 0 0 0 8.74l3.24-2.6Z" />
    <path fill="#EA4335" d="M12 6.23c1.43 0 2.71.49 3.72 1.44l2.79-2.79C16.82 3.3 14.62 2.28 12 2.28a9.72 9.72 0 0 0-8.68 5.35l3.24 2.6c.77-2.29 2.91-4 5.44-4Z" />
  </svg>
);

const ErrorBanner: React.FC<{ children: React.ReactNode }> = ({ children }) => (
  <p role="alert" className="mt-2 flex gap-2 rounded border border-[#B5573C] bg-[#B5573C]/10 p-3 text-sm text-[#3A2A1C]">
    <AlertCircle className="mt-0.5 h-4 w-4 shrink-0 text-[#B5573C]" />
    <span>{children}</span>
  </p>
);

export const GoogleAuthButton: React.FC = () => {
  const clientId = (import.meta.env.VITE_GOOGLE_CLIENT_ID as string | undefined)?.trim();
  const navigate = useNavigate();
  const { googleLogin } = useAuth();
  const containerRef = useRef<HTMLDivElement>(null);
  const [width, setWidth] = useState(400);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    const updateWidth = () => {
      const available = containerRef.current?.clientWidth;
      if (available) setWidth(Math.max(200, Math.min(400, Math.floor(available))));
    };
    updateWidth();
    window.addEventListener('resize', updateWidth);
    return () => window.removeEventListener('resize', updateWidth);
  }, []);

  const handleSuccess = async (response: CredentialResponse) => {
    if (!response.credential) {
      setError('Google did not return a credential. Please try again.');
      return;
    }
    setBusy(true);
    setError('');
    try {
      await googleLogin(response.credential);
      navigate('/dashboard');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Google sign-in failed. Please try again.');
    } finally {
      setBusy(false);
    }
  };

  if (!clientId) {
    return (
      <div ref={containerRef} className="w-full">
        <button
          type="button"
          onClick={() => setError('Google sign-in is not configured yet. Set VITE_GOOGLE_CLIENT_ID in frontend/.env.local and restart Vite.')}
          className="flex h-11 w-full items-center justify-center gap-3 rounded border border-[#3A2A1C]/40 text-[15px] font-medium text-[#3A2A1C]"
        >
          <GoogleMark />Continue with Google
        </button>
        {error ? <ErrorBanner>{error}</ErrorBanner> : null}
      </div>
    );
  }

  return (
    <div ref={containerRef} className="w-full">
      <GoogleOAuthProvider clientId={clientId}>
        {busy ? (
          <button type="button" disabled className="flex h-11 w-full items-center justify-center gap-2 rounded border border-[#3A2A1C]/40 text-[15px] font-medium text-[#3A2A1C] opacity-70">
            <Loader2 className="h-4 w-4 animate-spin" />Connecting to Google...
          </button>
        ) : (
          <GoogleLogin
            onSuccess={(response) => void handleSuccess(response)}
            onError={() => setError('Google sign-in was cancelled or failed. Please try again.')}
            useOneTap={false}
            type="standard"
            theme="outline"
            size="large"
            text="continue_with"
            shape="rectangular"
            logo_alignment="left"
            width={width}
          />
        )}
      </GoogleOAuthProvider>
      {error ? <ErrorBanner>{error}</ErrorBanner> : null}
    </div>
  );
};
