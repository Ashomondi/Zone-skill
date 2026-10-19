# Google Sign-In Setup

The app uses Google Identity Services (GIS) ID tokens. The frontend obtains a
signed ID token from Google, and the backend verifies it before issuing the
application's own JWT.

## 1. Create an OAuth client

1. Open [Google Cloud Console](https://console.cloud.google.com/).
2. Create or select a project.
3. Configure the OAuth consent screen if prompted.
4. Go to **APIs & Services → Credentials → Create credentials → OAuth client ID**.
5. Choose **Web application**.
6. Add the frontend dev origins under **Authorized JavaScript origins**:
   - `http://localhost:3000`
   - `http://localhost:5173` (optional, if you run Vite on its default port)
7. No redirect URI is needed for the ID-token flow.

## 2. Configure the backend

Add the client ID to `backend/.env`:

```env
GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
```

`GOOGLE_CLIENT_SECRET` is not needed for this flow.

Restart the backend after changing `.env`.

## 3. Configure the frontend

Create `frontend/.env.local`:

```env
VITE_GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
```

Restart Vite after changing the env file. Vite only exposes variables prefixed
with `VITE_`.

## 4. How the flow works

1. The user clicks **Continue with Google**.
2. Google returns an ID token to the frontend.
3. The frontend sends it to `POST /api/auth/google`:

   ```json
   { "credential": "<google-id-token>" }
   ```

4. The backend verifies the token's signature, issuer, audience, and expiry
   with Google's public keys.
5. The backend requires `email_verified: true`.
6. If a local account already uses that email, the Google identity is linked to
   it. Otherwise a new Google-only user is created.
7. The backend returns the same `{ token, user }` response as password login.

## Notes

- The client ID is public; the client secret must never be placed in frontend
  code.
- If Google sign-in is not configured, the button shows a clear configuration
  message instead of silently doing nothing.
- Existing email/password accounts are linked automatically only when Google
  reports the email as verified.
- Google-only accounts have no password hash. Password login for those accounts
  returns the normal "invalid credentials" response.
