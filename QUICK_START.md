# 🚀 Quick Start: Disable Authentication

**Goal:** Let users access your app at `https://zone-skill.onrender.com` without logging in.

## 1️⃣ Go to Render Dashboard

Visit: https://dashboard.render.com/

## 2️⃣ Select Your Frontend Service

Click on: `zone-skill` (or whatever your frontend service is named)

## 3️⃣ Add Environment Variable

1. Click on **Environment** in the left sidebar
2. Click **Add Environment Variable**
3. Enter:
   - **Key:** `VITE_DISABLE_AUTH`
   - **Value:** `true`
4. Click **Save Changes**

## 4️⃣ Wait for Deployment

Render will automatically redeploy your frontend. This takes ~2-3 minutes.

## 5️⃣ Test It!

Visit: https://zone-skill.onrender.com

You should:
- ✅ Skip the login page
- ✅ Go straight to the dashboard
- ✅ See "Guest User" in the top right
- ✅ Have full access to all features

## ✅ That's It!

No backend changes needed. The backend already supports this.

---

## 🔄 To Re-enable Login

1. Go back to Render → Frontend service → Environment
2. Find `VITE_DISABLE_AUTH`
3. Delete it or change to `false`
4. Save and redeploy

---

## 🐛 Troubleshooting

**Still seeing login page?**
- Wait for deployment to finish (check "Events" tab)
- Clear browser cache or use incognito mode
- Verify the variable is exactly: `VITE_DISABLE_AUTH` = `true`

**API errors?**
- Check both services are "Live" in Render
- Verify backend is at `https://zone-skill-1.onrender.com`
- Check backend logs for errors

---

## 📖 More Details

For complete documentation, see:
- `AUTH_BYPASS_SETUP.md` - Detailed setup guide
- `DEPLOYMENT_CHECKLIST.md` - Full deployment checklist
- `AUTH_CHANGES_SUMMARY.md` - Technical details

---

**Need help?** Check the other documentation files or review the code changes.
