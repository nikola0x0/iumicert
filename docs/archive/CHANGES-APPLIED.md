# Changes Applied for CI/CD Setup

## ✅ All Updates Complete

### 1. CORS Configuration (api_server.go)
**File**: `packages/issuer/cmd/api_server.go`

**Added support for both Vercel portals**:
- `STUDENT_PORTAL_URL` - For iu-micert.vercel.app
- `ISSUER_PORTAL_URL` - For your admin dashboard
- Kept legacy `FRONTEND_URL` for backward compatibility

**Lines 210-224**: Now reads three environment variables and adds all to allowed origins list.

---

### 2. Docker Compose Environment Variables
**File**: `packages/issuer/docker-compose.yml`

**Added lines 56-57**:
```yaml
STUDENT_PORTAL_URL: ${STUDENT_PORTAL_URL:-}
ISSUER_PORTAL_URL: ${ISSUER_PORTAL_URL:-}
```

This ensures the backend container receives these URLs from the `.env` file.

---

### 3. GitHub Actions Workflow Updates
**File**: `.github/workflows/deploy-issuer.yml`

**Updated in 4 locations**:

1. **Line 91-92**: Added secrets to workflow environment
   ```yaml
   STUDENT_PORTAL_URL: ${{ secrets.STUDENT_PORTAL_URL }}
   ISSUER_PORTAL_URL: ${{ secrets.ISSUER_PORTAL_URL }}
   ```

2. **Line 132-133**: Added to `.env` file generation on VM
   ```bash
   STUDENT_PORTAL_URL=${STUDENT_PORTAL_URL}
   ISSUER_PORTAL_URL=${ISSUER_PORTAL_URL}
   ```

3. **Line 182-183**: Added to docker-compose.prod.yml environment
   ```yaml
   STUDENT_PORTAL_URL: \${STUDENT_PORTAL_URL}
   ISSUER_PORTAL_URL: \${ISSUER_PORTAL_URL}
   ```

4. **Line 269**: Added to SSH command that executes deployment
   ```bash
   STUDENT_PORTAL_URL='${STUDENT_PORTAL_URL}' ISSUER_PORTAL_URL='${ISSUER_PORTAL_URL}'
   ```

---

## 🔑 GitHub Secrets You Need to Add

Before deploying, add these **2 NEW secrets** in GitHub:

**Settings → Secrets and variables → Actions → New repository secret**

| Secret Name | Value | Example |
|-------------|-------|---------|
| `STUDENT_PORTAL_URL` | Student portal Vercel URL | `https://iu-micert.vercel.app` |
| `ISSUER_PORTAL_URL` | Issuer portal Vercel URL | `https://your-issuer.vercel.app` |

### All Required Secrets (Complete List)

✅ Already Should Have:
- `GCP_SA_KEY`
- `GCP_VM_NAME`
- `GCP_VM_ZONE`
- `POSTGRES_PASSWORD`
- `ISSUER_PRIVATE_KEY`
- `SEPOLIA_RPC_URL`
- `IUMICERT_CONTRACT_ADDRESS`

🆕 Need to Add:
- `STUDENT_PORTAL_URL`
- `ISSUER_PORTAL_URL`

---

## 📝 Local .env Files

### Backend (.env)
Add these lines to `packages/issuer/.env`:
```bash
STUDENT_PORTAL_URL=http://localhost:3000
ISSUER_PORTAL_URL=http://localhost:3001
```

For production, these will be set via GitHub secrets.

---

## 🧪 Testing

### Local Testing
```bash
cd packages/issuer

# Set environment variables
export STUDENT_PORTAL_URL=http://localhost:3000
export ISSUER_PORTAL_URL=http://localhost:3001

# Start with docker-compose
docker-compose up

# Check CORS in logs
docker-compose logs issuer-backend | grep "CORS"
```

### After Deployment
1. Check backend logs for CORS configuration
2. Test from browser console on both portals:
   ```javascript
   fetch('http://YOUR_VM_IP:8080/api/health')
     .then(r => r.json())
     .then(d => console.log('✅ Success:', d))
     .catch(e => console.log('❌ CORS Error:', e))
   ```

---

## 🚀 Ready to Deploy!

All code changes are complete. You can now:

1. **Commit these changes**:
   ```bash
   git add .
   git commit -m "Add CORS support for both Vercel portals"
   git push origin main
   ```

2. **Add GitHub secrets** (2 new ones listed above)

3. **Create GCP VM** (follow PRE-DEPLOYMENT-CHECKLIST.md)

4. **Trigger deployment** (push will auto-deploy or trigger manually)

---

## 📚 Related Documentation

- **Complete Setup Guide**: `PRE-DEPLOYMENT-CHECKLIST.md`
- **CI/CD Details**: `packages/issuer/CI-CD-SETUP.md`
- **Authentication**: `packages/issuer/AUTHENTICATION-SETUP.md`
