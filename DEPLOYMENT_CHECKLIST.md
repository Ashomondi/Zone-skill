# 🚀 Deployment Checklist - Authentication Bypass

Use this checklist to deploy Zone-skill with authentication disabled.

## ✅ Render Configuration

### Frontend Service (`zone-skill.onrender.com`)

**Environment Variables:**
```bash
VITE_DISABLE_AUTH=true
VITE_API_BASE_URL=https://zone-skill-1.onrender.com/api
```

**Steps:**
1. Go to https://dashboard.render.com
2. Select your frontend service
3. Go to Environment tab
4. Add `VITE_DISABLE_AUTH` = `true`
5. Confirm `VITE_API_BASE_URL` = `https://zone-skill-1.onrender.com/api`
6. Click "Save Changes"
7. Wait for automatic redeploy (or trigger manually)

### Backend Service (`zone-skill-1.onrender.com`)

**No changes needed!** The backend already supports the bypass.

**Existing Environment Variables (verify):**
```bash
CORS_ALLOWED_ORIGIN=https://zone-skill.onrender.com
DATABASE_URL=<your-postgres-url>
JWT_SECRET=<your-secret>
```

## 🧪 Testing After Deployment

### 1. Test Frontend Access
Visit: https://zone-skill.onrender.com

**Expected:**
- ✅ No login page shown
- ✅ Direct access to dashboard
- ✅ User shown as "Guest User" in navbar
- ✅ All navigation links work

### 2. Test API Endpoints
Open browser console and check network tab:

**Expected:**
- ✅ All requests include `X-Disable-Auth: true` header
- ✅ No `Authorization: Bearer` headers
- ✅ API returns 200 status codes
- ✅ No 401 Unauthorized errors

### 3. Test Features
Try these features:

- [ ] Dashboard loads
- [ ] Browse jobs (Discover page)
- [ ] Save a job
- [ ] View saved jobs
- [ ] Create an application
- [ ] View applications
- [ ] Upload resume (if storage is configured)
- [ ] Chat assistant (if Gemini is configured)

## 🐛 Common Issues

### Issue: Still showing login page
**Solution:**
- Verify `VITE_DISABLE_AUTH=true` is set
- Check the service redeployed (check deployment logs)
- Clear browser cache (`Ctrl+Shift+R` or `Cmd+Shift+R`)
- Try incognito/private window

### Issue: 401 Unauthorized errors
**Solution:**
- Check `X-Disable-Auth` header is being sent (network tab)
- Verify backend is running
- Check CORS settings on backend

### Issue: CORS errors
**Solution:**
- Verify `CORS_ALLOWED_ORIGIN` on backend includes frontend URL
- Check it's `https://zone-skill.onrender.com` (no trailing slash)

### Issue: Features not working
**Solution:**
- Check backend environment variables (DATABASE_URL, etc.)
- View backend logs in Render dashboard
- Ensure database migrations have run

## 📊 Monitoring

After deployment, monitor:

1. **Render Dashboard**: Check both services are "Live"
2. **Backend Logs**: Look for startup messages and errors
3. **Frontend Logs**: Check browser console for errors
4. **Network Tab**: Verify API calls succeed

## 🔄 Rollback Plan

If something goes wrong:

1. Go to Render frontend service
2. Remove or set `VITE_DISABLE_AUTH=false`
3. Redeploy
4. Authentication will be re-enabled

## 📝 Notes

- All users will share the same guest account (`guest-user`)
- No data isolation between users
- Suitable for demos, not production
- To re-enable auth, simply remove the environment variable

## 🎯 Success Criteria

Deployment is successful when:
- ✅ Frontend loads without login prompt
- ✅ Dashboard displays correctly
- ✅ Job browsing works
- ✅ Applications can be created
- ✅ No console errors
- ✅ All features accessible

---

**Last Updated:** Deploy this configuration for your demo!
