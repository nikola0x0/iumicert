# IU-MiCert Vercel Deployments

## 🌐 Production URLs

### Student/Verifier Portal (Public)
**URL**: https://iu-micert.vercel.app
**Purpose**: Public portal for students and verifiers
**Features**:
- Upload and verify academic receipts
- View student journeys
- IPA verification
- Search students

### Issuer Portal (Admin)
**URL**: https://iumicert-issuer.vercel.app
**Purpose**: Administrative dashboard for institution staff
**Features**:
- Generate test data
- Process academic terms
- Generate receipts
- Publish roots to blockchain
- Reset system (dangerous)

---

## 🔧 GitHub Secrets to Add

Add these **exact values** in GitHub Repository Settings:

**Settings → Secrets and variables → Actions → New repository secret**

```
Name: STUDENT_PORTAL_URL
Value: https://iu-micert.vercel.app
```

```
Name: ISSUER_PORTAL_URL
Value: https://iumicert-issuer.vercel.app
```

---

## 🎯 Vercel Environment Variables

### For Student Portal (iu-micert)

**Project**: iu-micert
**Settings → Environment Variables**

```
NEXT_PUBLIC_API_URL = http://YOUR_VM_IP:8080
```

Replace `YOUR_VM_IP` with your GCP VM external IP after creating it.

---

### For Issuer Portal (iumicert-issuer)

**Project**: iumicert-issuer (or whatever you name it)
**Settings → Environment Variables**

```
NEXT_PUBLIC_API_URL = http://YOUR_VM_IP:8080
NEXT_PUBLIC_IUMICERT_CONTRACT_ADDRESS = 0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60
NEXT_PUBLIC_ISSUER_USERNAME = admin
NEXT_PUBLIC_ISSUER_PASSWORD = your_secure_password
NEXT_PUBLIC_ANALYTICS_DISABLED = true
```

**⚠️ Security**: Use a strong password for `NEXT_PUBLIC_ISSUER_PASSWORD` in production!

---

## ✅ CORS Configuration

Your backend is now configured to allow requests from both Vercel URLs.

**In production** (on GCP VM), the backend will read these from environment variables:
- `STUDENT_PORTAL_URL=https://iu-micert.vercel.app`
- `ISSUER_PORTAL_URL=https://iumicert-issuer.vercel.app`

**For local development**, it uses:
- `STUDENT_PORTAL_URL=http://localhost:3000`
- `ISSUER_PORTAL_URL=http://localhost:3001`

---

## 🚀 Deployment Checklist

### 1. Student Portal (Already Deployed)
- [x] Deployed at iu-micert.vercel.app
- [ ] Update `NEXT_PUBLIC_API_URL` with VM IP
- [ ] Test verification works

### 2. Issuer Portal (To Deploy)
- [ ] Deploy to Vercel
- [ ] Should get: iumicert-issuer.vercel.app
- [ ] Add all environment variables
- [ ] Enable Vercel Password Protection or Authentication
- [ ] Test login works

### 3. Backend API
- [ ] Add GitHub secrets (STUDENT_PORTAL_URL, ISSUER_PORTAL_URL)
- [ ] Create GCP VM
- [ ] Deploy via GitHub Actions
- [ ] Test CORS from both portals

### 4. Final Testing
- [ ] Student portal can reach API
- [ ] Issuer portal can reach API
- [ ] No CORS errors in browser console
- [ ] Login works on issuer portal
- [ ] Verification works on student portal

---

## 📝 Quick Commands

### Deploy Issuer Portal to Vercel

```bash
cd packages/issuer/web/iumicert-issuer

# Login to Vercel
vercel login

# Deploy
vercel --prod

# It should automatically suggest: iumicert-issuer
# If not, you can set it manually in Vercel dashboard
```

### Get GCP VM External IP

```bash
gcloud compute instances describe YOUR_VM_NAME \
  --zone=YOUR_VM_ZONE \
  --format='get(networkInterfaces[0].accessConfigs[0].natIP)'
```

### Update Vercel Environment Variables

```bash
# Option 1: Via Dashboard (Recommended)
# Go to project → Settings → Environment Variables

# Option 2: Via CLI
vercel env add NEXT_PUBLIC_API_URL production
# Then enter: http://YOUR_VM_IP:8080
```

---

## 🔍 Testing CORS

After deployment, test from browser console on each portal:

```javascript
// On iu-micert.vercel.app
fetch('http://YOUR_VM_IP:8080/api/health')
  .then(r => r.json())
  .then(d => console.log('✅ Student Portal CORS works:', d))
  .catch(e => console.error('❌ CORS Error:', e));

// On iumicert-issuer.vercel.app
fetch('http://YOUR_VM_IP:8080/api/health')
  .then(r => r.json())
  .then(d => console.log('✅ Issuer Portal CORS works:', d))
  .catch(e => console.error('❌ CORS Error:', e));
```

---

## 📚 Related Files

- **Backend CORS**: `packages/issuer/cmd/api_server.go` (lines 204-224)
- **Backend .env**: `packages/issuer/.env` (lines 34-40)
- **Production config**: `packages/issuer/.env.production` (lines 25-28)
- **GitHub Actions**: `.github/workflows/deploy-issuer.yml`

---

**Your Vercel URLs are now configured everywhere!** ✅
