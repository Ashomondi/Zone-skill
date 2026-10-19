# Authentication Bypass Setup Guide

This guide explains how to disable authentication and allow free access to your Zone-skill application.

## Quick Setup

### Option 1: Environment Variable (Recommended)

Add this environment variable to your **Render frontend service**:

```bash
VITE_DISABLE_AUTH=true
```

That's it! Redeploy and authentication will be bypassed.

### Option 2: Local Testing

Create a `.env.local` file in the `frontend/` directory:

```bash
VITE_API_BASE_URL=https://zone-skill-1.onrender.com/api
VITE_DISABLE_AUTH=true
```

Then run the frontend:

```bash
cd frontend
npm run dev
```

## Deployment Steps (Render)

### 1. Update Frontend Environment Variables

Go to your frontend service on Render (`https://zone-skill.onrender.com`) and add:

```bash
VITE_DISABLE_AUTH=true
VITE_API_BASE_URL=https://zone-skill-1.onrender.com/api
```

### 2. Redeploy Frontend

Click "Manual Deploy" → "Deploy latest commit"

### 3. Test

Visit `https://zone-skill.onrender.com`

You should:
- Be automatically logged in as "Guest User"
- Have full access to all features
- Not see login/signup pages

## How It Works

1. **Frontend**: Auto-logs in with a guest user and adds `X-Disable-Auth: true` header to all API requests
2. **Backend**: Accepts `X-Disable-Auth: true` header and injects guest user ID into requests
3. **Result**: No authentication required, everyone uses the same guest account

## What This Changes

- ✅ No login/signup required
- ✅ Instant access to dashboard
- ✅ All features work normally
- ⚠️ Everyone shares the same guest account
- ⚠️ No data isolation between users

## Re-enabling Authentication

To turn authentication back on:

1. Go to Render frontend environment variables
2. Set `VITE_DISABLE_AUTH=false` or remove it entirely
3. Redeploy the frontend

## Troubleshooting

### Still seeing login page?
- Check that `VITE_DISABLE_AUTH=true` is set in Render environment variables
- Redeploy the frontend after changing environment variables
- Clear browser cache or try incognito mode

### API errors?
- Verify frontend can reach backend at `https://zone-skill-1.onrender.com/api`
- Check CORS settings in backend `.env`:
  ```bash
  CORS_ALLOWED_ORIGIN=https://zone-skill.onrender.com
  ```

### Backend rejecting requests?
- The backend middleware already supports `X-Disable-Auth` header
- No backend changes needed!

## Security Note

⚠️ **This is for demos and testing only!**

Do not use authentication bypass with:
- Real user data
- Production environments
- Any sensitive information

When bypass is enabled, all users share the same guest account.
