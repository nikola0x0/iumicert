# Issuer Portal Authentication Guide

This guide explains how to secure your issuer portal with simple login authentication.

## Overview

The issuer portal now has **built-in login authentication** that protects all admin features. This is **much simpler** than API key management and provides a better user experience.

## How It Works

```
┌────────────────────────────────────────────────┐
│  User visits issuer portal                     │
│              ↓                                  │
│  Middleware checks auth cookie                 │
│              ↓                                  │
│  Not authenticated? → Redirect to /login       │
│              ↓                                  │
│  User enters username & password               │
│              ↓                                  │
│  Credentials valid? → Set auth cookie          │
│              ↓                                  │
│  Access granted to admin features              │
└────────────────────────────────────────────────┘
```

## Files Added

1. **`src/app/login/page.tsx`** - Login page with form
2. **`src/middleware.ts`** - Route protection
3. **`src/components/LogoutButton.tsx`** - Logout functionality
4. **`.env.example`** - Added auth credentials

## Setup Instructions

### Step 1: Set Environment Variables

For **local development**, create/update `.env.local`:

```bash
cd packages/issuer/web/iumicert-issuer

# Copy example
cp .env.example .env.local

# Edit .env.local and set:
NEXT_PUBLIC_ISSUER_USERNAME=your_username
NEXT_PUBLIC_ISSUER_PASSWORD=your_secure_password
```

For **Vercel production**, add environment variables in dashboard:

1. Go to Vercel Dashboard → Your Issuer Portal Project
2. Settings → Environment Variables
3. Add:
   - **Name**: `NEXT_PUBLIC_ISSUER_USERNAME`
   - **Value**: `admin` (or your chosen username)
   - Click "Add"

4. Add:
   - **Name**: `NEXT_PUBLIC_ISSUER_PASSWORD`
   - **Value**: `your_secure_password` (use a strong password!)
   - Click "Add"

5. Redeploy: Deployments → Click ⋮ on latest → Redeploy

### Step 2: Add Logout Button to Layout

Update your main layout to include the logout button:

```tsx
// src/app/layout.tsx
import LogoutButton from "@/components/LogoutButton";

export default function RootLayout({ children }) {
  return (
    <html>
      <body>
        <nav className="flex justify-between items-center p-4">
          <h1>IU-MiCert Issuer Portal</h1>
          <LogoutButton />
        </nav>
        {children}
      </body>
    </html>
  );
}
```

### Step 3: Test Locally

```bash
cd packages/issuer/web/iumicert-issuer

# Install dependencies
npm install

# Run development server
npm run dev

# Visit http://localhost:3000
# You'll be redirected to /login
# Enter username and password from .env.local
```

### Step 4: Deploy to Vercel

```bash
# From issuer portal directory
vercel --prod

# Or if already linked
vercel link
vercel --prod
```

## Security Features

### ✅ What's Protected
- All admin pages (data generation, processing, publishing)
- All admin API calls
- Session expires after 24 hours
- Logout clears session immediately

### ✅ What's NOT Protected (Student Portal)
The student/verifier portal at **iu-micert.vercel.app** remains **public** and accessible to everyone. Only the issuer portal requires login.

## Configuration Options

### Change Session Duration

Edit `src/app/login/page.tsx`:

```typescript
// Change from 24 hours to 7 days
document.cookie = `issuer_auth=true; path=/; max-age=604800`; // 7 days
```

### Multiple Users

For multiple staff members, you can:

**Option A: Single shared password** (current implementation)
- Simple, works for small teams
- Everyone uses same credentials

**Option B: Multiple credentials** (requires code update)
```typescript
// src/app/login/page.tsx
const validUsers = [
  { username: "admin1", password: "pass1" },
  { username: "admin2", password: "pass2" },
];

const isValid = validUsers.some(
  user => user.username === username && user.password === password
);
```

**Option C: Database-backed auth** (for larger teams)
- Store users in PostgreSQL
- Hash passwords with bcrypt
- More complex but more secure

## Comparison: Simple Login vs API Key

| Feature | Simple Login (✅ Implemented) | API Key |
|---------|------------------------------|---------|
| User Experience | Friendly login page | Technical header management |
| Setup Complexity | Easy (5 minutes) | Medium (add middleware) |
| Multiple Users | Easy to add | Same key for all |
| Session Management | Built-in (cookies) | Manual tracking |
| Logout | One click | Clear storage manually |
| Forgot Password | Can add reset flow | Share key again |
| **Recommended** | **YES** | For API-to-API only |

## Backend Simplification

With login auth in the issuer portal, you can **simplify the backend**:

### Before (with API key middleware):
```go
// Complex - need API key validation
router.HandleFunc("/api/admin/generate-data", RequireAPIKey(GenerateDataHandler))
```

### After (with portal login):
```go
// Simple - only issuer portal (logged-in users) can access
router.HandleFunc("/api/admin/generate-data", GenerateDataHandler)

// Optional: Add origin check for extra security
if origin != "https://your-issuer-portal.vercel.app" {
    http.Error(w, "Forbidden", http.StatusForbidden)
    return
}
```

## Where to Update API URL

### Student Portal (iu-micert.vercel.app)
1. Vercel Dashboard → iu-micert project
2. Settings → Environment Variables
3. Update `NEXT_PUBLIC_API_URL` = `http://YOUR_VM_IP:8080`
4. Redeploy

### Issuer Portal
1. Vercel Dashboard → Issuer Portal project
2. Settings → Environment Variables
3. Add/Update:
   - `NEXT_PUBLIC_API_URL` = `http://YOUR_VM_IP:8080`
   - `NEXT_PUBLIC_ISSUER_USERNAME` = `admin`
   - `NEXT_PUBLIC_ISSUER_PASSWORD` = `your_password`
4. Redeploy

## Testing Checklist

- [ ] Login page displays at /login
- [ ] Invalid credentials show error message
- [ ] Valid credentials redirect to dashboard
- [ ] Logout button appears on all pages
- [ ] Logout clears session and redirects to login
- [ ] Direct access to protected pages redirects to login
- [ ] Session persists across page refreshes
- [ ] Session expires after configured time

## Troubleshooting

### Issue: Infinite redirect loop
**Cause**: Middleware configuration error

**Fix**: Check `src/middleware.ts` matcher config excludes static files

### Issue: Login successful but still redirects
**Cause**: Cookie not being set

**Fix**: Check browser developer tools → Application → Cookies
- Should see `issuer_auth=true`
- If missing, check cookie path and domain settings

### Issue: Works locally but not on Vercel
**Cause**: Environment variables not set

**Fix**:
1. Vercel Dashboard → Environment Variables
2. Ensure `NEXT_PUBLIC_ISSUER_USERNAME` and `NEXT_PUBLIC_ISSUER_PASSWORD` are set
3. Redeploy

## Production Security Tips

1. **Use strong passwords**: Minimum 16 characters, mix of letters/numbers/symbols
2. **Change default credentials**: Never use `admin/admin123` in production
3. **Enable HTTPS**: Vercel provides this automatically
4. **Rotate passwords**: Update every 3-6 months
5. **Monitor login attempts**: Check Vercel logs for suspicious activity
6. **Add rate limiting**: Prevent brute force attacks (future enhancement)

## Future Enhancements

Possible improvements for larger deployments:

1. **Password Reset Flow**
   - Forgot password link
   - Email-based reset

2. **Multi-Factor Authentication (MFA)**
   - TOTP codes (Google Authenticator)
   - SMS codes

3. **Audit Logging**
   - Track who logged in when
   - Log admin actions

4. **Role-Based Access**
   - Admin vs. Read-only users
   - Different permission levels

5. **Database-Backed Users**
   - Store users in PostgreSQL
   - Password hashing with bcrypt
   - User management UI

---

**Recommendation**: Start with the simple login implementation provided. It's secure enough for most use cases and much easier to manage than API keys. You can always enhance it later if needed.
