# Disabling Authentication

For development, testing, or demo purposes, you can temporarily disable authentication to allow free access to the application.

## Frontend Setup

Add the following to your frontend `.env.local` file:

```bash
VITE_DISABLE_AUTH=true
```

## Backend Setup

The backend will automatically accept requests with the `X-Disable-Auth: true` header when authentication is bypassed on the frontend.

## How It Works

### Frontend
- When `VITE_DISABLE_AUTH=true`, the `useAuth` hook automatically logs in with a guest user
- All API requests include the `X-Disable-Auth: true` header instead of a JWT token
- Users can access all protected routes without signing in

### Backend
- The authentication middleware checks for the `X-Disable-Auth: true` header
- If present, it injects a guest user ID (`guest-user`) into the request context
- All endpoints work normally with the guest user

## Deployment

### Render Environment Variables

**Frontend (`https://zone-skill.onrender.com`)**
```bash
VITE_DISABLE_AUTH=true
VITE_API_BASE_URL=https://zone-skill-1.onrender.com/api
```

**Backend (`https://zone-skill-1.onrender.com`)**
```bash
# No changes needed - the backend already supports the X-Disable-Auth header
CORS_ALLOWED_ORIGIN=https://zone-skill.onrender.com
```

## Re-enabling Authentication

To re-enable authentication:

1. Remove `VITE_DISABLE_AUTH` from your environment variables or set it to `false`
2. Redeploy both frontend and backend
3. Users will be required to sign in again

## Security Warning

⚠️ **Do not use this in production with real user data!** This bypass is intended for:
- Development and testing
- Public demos
- Prototype presentations

When authentication is disabled, all users share the same guest account and can see each other's data.
