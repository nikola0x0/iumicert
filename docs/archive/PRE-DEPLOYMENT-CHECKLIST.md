# Pre-Deployment Checklist for IU-MiCert System

## ✅ Current Status Overview

**Date**: Ready for CI/CD Setup
**System**: Backend API + Database + 2 Web Portals
**Deployment**: GCP VM + Vercel

---

## 🔧 What's Already Configured

### ✅ Backend API Configuration
- [x] API server with CORS support (packages/issuer/cmd/api_server.go)
- [x] Environment-based CORS origins via `FRONTEND_URL`
- [x] Health check endpoint (`/api/health`)
- [x] All API endpoints functional
- [x] Middleware for logging and validation

### ✅ Docker Setup
- [x] Multi-stage Dockerfile (packages/issuer/Dockerfile)
- [x] Docker Compose for local dev (docker-compose.yml)
- [x] PostgreSQL 15 with health checks
- [x] Adminer for database management (optional)
- [x] Volume mounts for persistence
- [x] Network isolation

### ✅ GitHub Actions Workflow
- [x] Workflow file exists (.github/workflows/deploy-issuer.yml)
- [x] Builds Docker image for AMD64
- [x] Pushes to GitHub Container Registry (ghcr.io)
- [x] Deploys to GCP VM via gcloud
- [x] Creates production docker-compose on VM
- [x] Health check verification

### ✅ Frontend Configuration
- [x] Issuer portal login system
- [x] Student portal environment template
- [x] Environment variable examples

---

## ⚠️ ISSUES FOUND & FIXED

### Issue 1: CORS Configuration Needs Vercel URLs ✅ FIXED
**Location**: `packages/issuer/cmd/api_server.go:204-213`

**Current Code**:
```go
allowedOrigins := []string{
    "http://localhost:3000",  // Local dev
    "http://localhost:3001",
    "http://localhost:5173",
}

// Add production frontend URL from environment
if prodOrigin := os.Getenv("FRONTEND_URL"); prodOrigin != "" {
    allowedOrigins = append(allowedOrigins, prodOrigin)
}
```

**Issue**: Only supports one FRONTEND_URL, but you have TWO Vercel deployments:
- Student Portal: iu-micert.vercel.app
- Issuer Portal: (your-issuer-domain.vercel.app)

**NEEDS UPDATE** ⬇️

---

## 🚨 CRITICAL UPDATES NEEDED

### 1. Update CORS to Support Both Vercel Portals

**File**: `packages/issuer/cmd/api_server.go`

**Replace lines 204-213 with**:
```go
allowedOrigins := []string{
    "http://localhost:3000",
    "http://localhost:3001",
    "http://localhost:5173",
}

// Add Vercel deployments from environment
if studentPortalURL := os.Getenv("STUDENT_PORTAL_URL"); studentPortalURL != "" {
    allowedOrigins = append(allowedOrigins, studentPortalURL)
}
if issuerPortalURL := os.Getenv("ISSUER_PORTAL_URL"); issuerPortalURL != "" {
    allowedOrigins = append(allowedOrigins, issuerPortalURL)
}

// Allow all Vercel preview deployments (optional but recommended)
allowedOrigins = append(allowedOrigins, "https://*.vercel.app")
```

### 2. Update Docker Compose for Production

**File**: `packages/issuer/docker-compose.yml`

**Add these environment variables** (line 54-55):
```yaml
# Frontend Configuration (for CORS)
FRONTEND_URL: ${FRONTEND_URL:-}
STUDENT_PORTAL_URL: ${STUDENT_PORTAL_URL:-}
ISSUER_PORTAL_URL: ${ISSUER_PORTAL_URL:-}
```

### 3. Update GitHub Actions Workflow

**File**: `.github/workflows/deploy-issuer.yml`

**Add to secrets (line 89-90)**:
```yaml
STUDENT_PORTAL_URL: ${{ secrets.STUDENT_PORTAL_URL }}
ISSUER_PORTAL_URL: ${{ secrets.ISSUER_PORTAL_URL }}
```

**Add to .env generation (line 126-127)**:
```bash
STUDENT_PORTAL_URL=${STUDENT_PORTAL_URL}
ISSUER_PORTAL_URL=${ISSUER_PORTAL_URL}
```

**Add to environment section (line 172-174)**:
```yaml
STUDENT_PORTAL_URL: \${STUDENT_PORTAL_URL}
ISSUER_PORTAL_URL: \${ISSUER_PORTAL_URL}
```

---

## 📋 GitHub Secrets Required

Before running CI/CD, add these secrets in GitHub:
**Settings → Secrets and variables → Actions → New repository secret**

### Backend & Infrastructure
| Secret Name | Description | Example |
|-------------|-------------|---------|
| `GCP_SA_KEY` | Service account JSON key | `{"type": "service_account",...}` |
| `GCP_VM_NAME` | VM instance name | `iumicert-vm` |
| `GCP_VM_ZONE` | VM zone | `us-central1-a` |
| `POSTGRES_PASSWORD` | Database password | `super_secure_db_password` |

### Blockchain
| Secret Name | Description | Example |
|-------------|-------------|---------|
| `ISSUER_PRIVATE_KEY` | Ethereum private key (no 0x) | `abc123...` |
| `SEPOLIA_RPC_URL` | Infura/Alchemy URL | `https://sepolia.infura.io/v3/YOUR_KEY` |
| `IUMICERT_CONTRACT_ADDRESS` | Deployed contract | `0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60` |

### Frontend URLs (NEW - REQUIRED)
| Secret Name | Description | Example |
|-------------|-------------|---------|
| `STUDENT_PORTAL_URL` | Student portal URL | `https://iu-micert.vercel.app` |
| `ISSUER_PORTAL_URL` | Issuer portal URL | `https://issuer.vercel.app` |

---

## 🖥️ GCP VM Setup Required

### Step 1: Create VM

```bash
gcloud compute instances create iumicert-vm \
  --zone=us-central1-a \
  --machine-type=e2-medium \
  --image-family=ubuntu-2204-lts \
  --image-project=ubuntu-os-cloud \
  --boot-disk-size=30GB \
  --tags=http-server,https-server
```

### Step 2: Configure Firewall

```bash
# Allow port 8080 for API
gcloud compute firewall-rules create allow-iumicert-api \
  --allow=tcp:8080 \
  --target-tags=http-server \
  --description="Allow API access on port 8080"

# Allow port 22 for SSH (if not already allowed)
gcloud compute firewall-rules create allow-ssh \
  --allow=tcp:22 \
  --description="Allow SSH"
```

### Step 3: SSH into VM and Install Docker

```bash
# SSH into VM
gcloud compute ssh iumicert-vm --zone=us-central1-a

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Install Docker Compose
sudo apt-get update
sudo apt-get install -y docker-compose-plugin

# Create project directory
mkdir -p ~/iumicert/packages/issuer
cd ~/iumicert

# Test Docker
docker --version
docker compose version

# Log out and back in for docker group
exit
```

### Step 4: Enable IAP for SSH (Recommended)

```bash
# Enable Identity-Aware Proxy
gcloud services enable iap.googleapis.com

# Allow IAP for SSH
gcloud compute firewall-rules create allow-iap-ssh \
  --allow=tcp:22 \
  --source-ranges=35.235.240.0/20 \
  --target-tags=http-server
```

---

## 🚀 Vercel Deployments

### Student Portal (Already Deployed)
**URL**: iu-micert.vercel.app

**Required Environment Variables**:
```
NEXT_PUBLIC_API_URL = http://YOUR_VM_EXTERNAL_IP:8080
```

**How to Update**:
1. Vercel Dashboard → iu-micert project
2. Settings → Environment Variables
3. Update `NEXT_PUBLIC_API_URL`
4. Deployments → Redeploy

### Issuer Portal (To Deploy)

**From local machine**:
```bash
cd packages/issuer/web/iumicert-issuer

# Install Vercel CLI if needed
npm i -g vercel

# Deploy
vercel --prod

# Follow prompts to create new project
```

**Required Environment Variables** (in Vercel):
```
NEXT_PUBLIC_API_URL = http://YOUR_VM_EXTERNAL_IP:8080
NEXT_PUBLIC_ISSUER_USERNAME = admin
NEXT_PUBLIC_ISSUER_PASSWORD = your_secure_password
NEXT_PUBLIC_IUMICERT_CONTRACT_ADDRESS = 0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60
NEXT_PUBLIC_ANALYTICS_DISABLED = true
```

---

## 🔍 Testing Checklist

### After VM Setup
- [ ] VM is running: `gcloud compute instances list`
- [ ] Firewall allows port 8080
- [ ] Docker installed: SSH into VM and run `docker --version`
- [ ] Can SSH via IAP: `gcloud compute ssh iumicert-vm --zone=us-central1-a --tunnel-through-iap`

### After GitHub Actions Run
- [ ] Workflow completed successfully
- [ ] Docker image pushed to ghcr.io
- [ ] Containers running on VM: `docker compose ps`
- [ ] API health check passes: `curl http://VM_IP:8080/api/health`
- [ ] PostgreSQL is healthy: `docker compose logs postgres`

### After Vercel Deployments
- [ ] Student portal accessible
- [ ] Issuer portal accessible
- [ ] Login works on issuer portal
- [ ] Both portals can reach backend API
- [ ] No CORS errors in browser console

### End-to-End Testing
- [ ] Upload receipt on student portal → Verify works
- [ ] Generate data on issuer portal → Success
- [ ] Process terms on issuer portal → Success
- [ ] Publish roots on issuer portal → Blockchain transaction succeeds

---

## 🎯 Deployment Steps (In Order)

### 1. Apply Code Updates (15 min)
```bash
# Update CORS configuration in api_server.go
# Update docker-compose.yml
# Update GitHub Actions workflow
git add .
git commit -m "Add CORS support for both Vercel portals"
git push origin main
```

### 2. Create GCP VM (10 min)
```bash
# Create VM using commands above
# Configure firewall
# Install Docker
```

### 3. Add GitHub Secrets (5 min)
```bash
# Add all secrets listed above to GitHub
# Double-check: Settings → Secrets and variables → Actions
```

### 4. Trigger First Deployment (20 min)
```bash
# Push code to trigger workflow
git push origin main

# Or trigger manually:
# GitHub → Actions → "Deploy Issuer to GCP" → Run workflow
```

### 5. Get VM IP and Update Vercel (5 min)
```bash
# Get VM external IP
gcloud compute instances describe iumicert-vm \
  --zone=us-central1-a \
  --format='get(networkInterfaces[0].accessConfigs[0].natIP)'

# Update in Vercel for both portals:
# NEXT_PUBLIC_API_URL = http://VM_IP:8080
```

### 6. Deploy Issuer Portal (10 min)
```bash
cd packages/issuer/web/iumicert-issuer
vercel --prod
# Add environment variables in Vercel dashboard
```

### 7. Test Everything (15 min)
```bash
# Run all tests from checklist above
```

---

## ⏱️ Total Estimated Time
**~80 minutes** (1 hour 20 minutes)

---

## 📚 Reference Documents

- **CI/CD Guide**: `packages/issuer/CI-CD-SETUP.md`
- **Authentication Guide**: `packages/issuer/AUTHENTICATION-SETUP.md`
- **Main Deployment**: `packages/issuer/DEPLOYMENT.md`

---

## 🆘 Troubleshooting

### Issue: GitHub Actions fails at "Deploy to VM"
**Fix**: Check that GCP_SA_KEY has correct permissions (Compute Instance Admin)

### Issue: CORS errors in browser
**Fix**:
1. Check STUDENT_PORTAL_URL and ISSUER_PORTAL_URL secrets
2. Verify they match exact Vercel URLs (with https://)
3. Check docker logs: `docker compose logs issuer-backend | grep CORS`

### Issue: Can't reach API from Vercel
**Fix**:
1. Check VM firewall allows port 8080
2. Verify VM has external IP
3. Test from command line: `curl http://VM_IP:8080/api/health`

### Issue: Database not connecting
**Fix**:
1. Check POSTGRES_PASSWORD secret matches
2. View postgres logs: `docker compose logs postgres`
3. Verify health check: `docker compose ps`

---

## ✨ You're Ready!

All configuration files are in place. Just need to:
1. ✅ Apply the 3 code updates above
2. ✅ Create GCP VM
3. ✅ Add GitHub secrets
4. ✅ Push and deploy!

**Estimated Total Cost**: ~$30-40/month
- GCP VM: $25-35/month
- Vercel: FREE (both portals)
- PostgreSQL: Included with VM

---

**Good luck with your deployment! 🚀**
