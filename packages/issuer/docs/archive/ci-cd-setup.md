# CI/CD Setup Guide for IU-MiCert Complete System

This guide explains how to set up automated deployments for the complete IU-MiCert system with **two web portals** and a shared backend API.

## System Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    IU-MiCert System                         │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────────────┐      ┌──────────────────────┐    │
│  │  Issuer Portal      │      │  Student/Verifier    │    │
│  │  (Admin Dashboard)  │      │  Portal              │    │
│  │                     │      │                      │    │
│  │  - Generate Data    │      │  - View Receipts     │    │
│  │  - Process Terms    │      │  - Verify Credentials│    │
│  │  - Publish to Chain │      │  - Search Students   │    │
│  │  - Reset System     │      │  - IPA Verification  │    │
│  │                     │      │                      │    │
│  │  Vercel Deploy      │      │  Vercel Deploy       │    │
│  │  (Restricted)       │      │  iu-micert.vercel.app│    │
│  └──────────┬──────────┘      └──────────┬───────────┘    │
│             │                            │                 │
│             └────────────┬───────────────┘                 │
│                          │                                 │
│                ┌─────────▼─────────┐                       │
│                │   Backend API     │                       │
│                │   (Go + Docker)   │                       │
│                │                   │                       │
│                │  Port 8080        │                       │
│                │  GCP VM           │                       │
│                │  + PostgreSQL     │                       │
│                │  + Ethereum       │                       │
│                └───────────────────┘                       │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

## Components

### 1. Backend API (GCP VM)
- **Location**: GCP VM with Docker
- **Purpose**: Shared API for both portals
- **Deployment**: GitHub Actions → Docker → GCP VM
- **Authorization**: API key-based for admin operations

### 2. Issuer Portal (Admin)
- **Location**: `packages/issuer/web/iumicert-issuer/`
- **Purpose**: Administrative dashboard for institution staff
- **Deployment**: Vercel (restricted access)
- **Features**: Generate data, process terms, publish roots, reset system

### 3. Student/Verifier Portal
- **Location**: `packages/client/iumicert-client/`
- **Purpose**: Public portal for students and verifiers
- **Deployment**: Vercel at **iu-micert.vercel.app**
- **Features**: View receipts, verify credentials, search students

## Prerequisites

- GitHub repository with the code
- GCP VM (AMD64) with Docker installed
- GCP Service Account with VM access
- Ethereum Sepolia private key and RPC URL
- Vercel account for frontend deployments

---

## Part 1: Backend API Deployment (GCP VM)

### Step 1: Prepare Your GCP VM

#### Option A: Automated Setup (Recommended)

```bash
# Copy the setup script to your VM
# From your local machine:
scp packages/issuer/scripts/vm-setup.sh YOUR_VM_USER@YOUR_VM_IP:~/

# SSH into VM
ssh YOUR_VM_USER@YOUR_VM_IP

# Run setup script
bash ~/vm-setup.sh

# Log out and back in for Docker group changes
exit
```

#### Option B: Manual Setup

```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Install Docker Compose
sudo apt-get install -y docker-compose-plugin

# Create project directory
mkdir -p ~/iumicert
cd ~/iumicert

# Configure firewall (CRITICAL: Allow both portals)
sudo ufw allow 22/tcp
sudo ufw allow 8080/tcp
sudo ufw enable

# Log out and back in
exit
```

### Step 2: Configure API Authorization

The API needs to differentiate between:
- **Public endpoints** (student/verifier portal) - No auth required
- **Admin endpoints** (issuer portal) - Require API key

Create a strong API key:

```bash
# Generate a secure API key
openssl rand -hex 32

# Save this as ISSUER_API_KEY secret
```

### Step 3: Create GCP Service Account

1. **Go to GCP Console** → IAM & Admin → Service Accounts
2. **Create Service Account** with name: `github-actions-deployer`
3. **Grant roles**:
   - Compute Instance Admin (v1) - `roles/compute.instanceAdmin.v1`
   - Service Account User - `roles/iam.serviceAccountUser`
4. **Create JSON key**:
   - Click on the service account
   - Keys → Add Key → Create New Key → JSON
   - Download and save the JSON file securely

### Step 4: Configure GitHub Secrets

Go to your GitHub repository → Settings → Secrets and variables → Actions

#### Backend Secrets

| Secret Name | Description | Example |
|------------|-------------|---------|
| `GCP_SA_KEY` | Service account JSON key | `{"type": "service_account",...}` |
| `GCP_VM_NAME` | Your VM instance name | `iumicert-issuer-vm` |
| `GCP_VM_ZONE` | Your VM zone | `us-central1-a` |
| `ISSUER_PRIVATE_KEY` | Ethereum private key (without 0x) | `abc123...` |
| `SEPOLIA_RPC_URL` | Infura/Alchemy RPC URL | `https://sepolia.infura.io/v3/YOUR_KEY` |
| `POSTGRES_PASSWORD` | Database password | `your_secure_password` |
| `IUMICERT_CONTRACT_ADDRESS` | Deployed contract | `0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60` |
| `ISSUER_API_KEY` | Admin API key (generated above) | `abc123...` |

### Step 5: Update Backend for CORS

The backend needs to allow requests from **both Vercel deployments**:

Update your API server configuration:

```go
// cmd/api_server.go
allowedOrigins := []string{
    "http://localhost:3000",
    "http://localhost:3001",
    "https://iu-micert.vercel.app",                    // Student portal
    "https://issuer-portal.vercel.app",                 // Issuer portal (your domain)
    "https://*.vercel.app",                             // Preview deployments
}
```

### Step 6: Add API Key Middleware

Protected endpoints need authorization checking:

```go
// middleware/auth.go
func RequireAPIKey(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        apiKey := r.Header.Get("X-API-Key")
        expectedKey := os.Getenv("ISSUER_API_KEY")

        if apiKey == "" || apiKey != expectedKey {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
    }
}

// Apply to admin routes
router.HandleFunc("/api/admin/generate-data", RequireAPIKey(GenerateDataHandler))
router.HandleFunc("/api/admin/batch-process", RequireAPIKey(BatchProcessHandler))
router.HandleFunc("/api/admin/reset", RequireAPIKey(ResetHandler))
// etc.
```

### Step 7: Deploy Backend with GitHub Actions

Create `.github/workflows/deploy-issuer.yml`:

```yaml
name: Deploy Issuer Backend to GCP

on:
  push:
    branches: [main, issuer_v2]
    paths:
      - 'packages/issuer/**'
      - 'packages/crypto/**'
      - '.github/workflows/deploy-issuer.yml'
  workflow_dispatch:

jobs:
  deploy:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Log in to GitHub Container Registry
        run: echo "${{ secrets.GITHUB_TOKEN }}" | docker login ghcr.io -u ${{ github.actor }} --password-stdin

      - name: Build and push Docker image
        run: |
          docker build -f packages/issuer/Dockerfile \
            --platform linux/amd64 \
            -t ghcr.io/${{ github.repository_owner }}/iumicert-issuer:latest \
            .
          docker push ghcr.io/${{ github.repository_owner }}/iumicert-issuer:latest

      - name: Deploy to GCP VM
        uses: google-github-actions/ssh-compute@v1
        with:
          instance_name: ${{ secrets.GCP_VM_NAME }}
          zone: ${{ secrets.GCP_VM_ZONE }}
          ssh_private_key: ${{ secrets.GCP_SSH_KEY }}
          command: |
            cd ~/iumicert/packages/issuer
            docker compose -f docker-compose.prod.yml pull
            docker compose -f docker-compose.prod.yml up -d

      - name: Verify deployment
        run: |
          sleep 10
          curl -f "http://${{ secrets.GCP_VM_EXTERNAL_IP }}:8080/api/health" || exit 1
```

---

## Part 2: Student/Verifier Portal Deployment (Vercel)

### Current Status
- **Already deployed** at: **iu-micert.vercel.app**
- **Repository**: `packages/client/iumicert-client/`

### Environment Variables

In Vercel dashboard for student portal, set:

| Variable | Value |
|----------|-------|
| `NEXT_PUBLIC_API_URL` | `http://YOUR_GCP_VM_IP:8080` |

### API Endpoints Used by Student Portal

Public endpoints (no auth required):
- `GET /api/health` - Health check
- `POST /api/verifier/ipa-verify` - Verify receipts
- `GET /api/students` - List students
- `GET /api/students/{id}/journey` - Get student journey

### Update Vercel Deployment

```bash
# From local machine
cd packages/client/iumicert-client

# Install Vercel CLI if needed
npm i -g vercel

# Deploy
vercel --prod

# Or link to existing project
vercel link
vercel --prod
```

---

## Part 3: Issuer Portal Deployment (Vercel)

### Repository Location
- **Path**: `packages/issuer/web/iumicert-issuer/`

### Environment Variables

In Vercel dashboard for issuer portal, set:

| Variable | Value | Description |
|----------|-------|-------------|
| `NEXT_PUBLIC_API_URL` | `http://YOUR_GCP_VM_IP:8080` | Backend API URL |
| `NEXT_PUBLIC_ISSUER_API_KEY` | `your_api_key_from_github_secrets` | Admin API key |

### Protected Features

The issuer portal has access to admin endpoints:
- `POST /api/admin/generate-data` - Generate test data
- `POST /api/admin/batch-process` - Process all terms
- `POST /api/admin/generate-receipts` - Create receipts
- `POST /api/admin/publish-roots` - Publish to blockchain
- `POST /api/admin/reset` - Reset system (DANGEROUS)

### Deploy Issuer Portal

```bash
cd packages/issuer/web/iumicert-issuer

# Install dependencies
npm install

# Deploy to Vercel
vercel --prod

# Configure environment variables in Vercel dashboard
# Settings → Environment Variables → Add NEXT_PUBLIC_API_URL and NEXT_PUBLIC_ISSUER_API_KEY
```

### Restrict Access

Since this is an admin portal, you should:

1. **Use Vercel Authentication** (recommended):
   - Vercel → Settings → Authentication → Enable
   - Add your institution email addresses

2. **Or use Vercel Password Protection**:
   - Vercel → Settings → Password Protection → Enable
   - Share password only with authorized staff

---

## API Endpoint Authorization Reference

### Public Endpoints (No Auth)
Used by student/verifier portal:
```
GET  /api/health
GET  /api/terms
GET  /api/students
GET  /api/students/:id/journey
POST /api/verifier/ipa-verify
POST /api/receipts/verify
```

### Protected Endpoints (API Key Required)
Used by issuer portal only:
```
POST /api/admin/generate-data
POST /api/admin/batch-process
POST /api/admin/generate-receipts
POST /api/admin/publish-roots
POST /api/admin/reset
POST /api/blockchain/publish
```

### How Protection Works

Issuer portal sends API key in header:
```typescript
// Issuer portal API client
const response = await fetch(`${API_BASE_URL}/api/admin/generate-data`, {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'X-API-Key': process.env.NEXT_PUBLIC_ISSUER_API_KEY,
  },
});
```

Backend validates API key:
```go
func RequireAPIKey(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        apiKey := r.Header.Get("X-API-Key")
        if apiKey != os.Getenv("ISSUER_API_KEY") {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    }
}
```

---

## Complete Deployment Checklist

### Backend (GCP VM)
- [ ] VM created and Docker installed
- [ ] GitHub secrets configured
- [ ] CORS allows both Vercel domains
- [ ] API key middleware added
- [ ] GitHub Actions workflow created
- [ ] Firewall allows port 8080
- [ ] Health endpoint accessible

### Student Portal (iu-micert.vercel.app)
- [ ] Already deployed ✓
- [ ] Environment variable `NEXT_PUBLIC_API_URL` set
- [ ] Public endpoints accessible
- [ ] Verification working

### Issuer Portal
- [ ] Deployed to Vercel
- [ ] Environment variables set (API_URL + API_KEY)
- [ ] Authentication/password protection enabled
- [ ] Admin endpoints working
- [ ] Only accessible to authorized users

---

## Testing the Complete System

### 1. Test Backend API

```bash
# Health check
curl http://YOUR_VM_IP:8080/api/health

# Public endpoint (no auth)
curl http://YOUR_VM_IP:8080/api/students

# Admin endpoint (should fail without key)
curl http://YOUR_VM_IP:8080/api/admin/generate-data
# Expected: 401 Unauthorized

# Admin endpoint (with key)
curl -H "X-API-Key: YOUR_API_KEY" \
  -X POST http://YOUR_VM_IP:8080/api/admin/generate-data
# Expected: Success
```

### 2. Test Student Portal

Visit: **https://iu-micert.vercel.app**

- Upload and verify a receipt
- Search for students
- View academic journeys
- Check that all features work

### 3. Test Issuer Portal

Visit: Your issuer portal URL

- Login with credentials
- Generate test data
- Process terms
- Generate receipts
- Check that admin features work

---

## Monitoring & Maintenance

### Monitor Backend

```bash
# SSH into VM
ssh YOUR_VM_USER@YOUR_VM_IP

# Check health
~/check-health.sh

# View logs
cd ~/iumicert/packages/issuer
docker compose -f docker-compose.prod.yml logs -f

# Check API key authentication
docker compose logs issuer-backend | grep "Unauthorized"
```

### Monitor Vercel Deployments

1. Vercel Dashboard → Your Projects
2. Check deployment status
3. View function logs
4. Monitor bandwidth usage

### Update Deployments

**Backend**: Push to GitHub → Auto-deploys
```bash
git push origin main
```

**Frontend**: Use Vercel CLI or push to GitHub
```bash
vercel --prod
```

---

## Security Best Practices

1. **Never commit API keys** - Use environment variables
2. **Rotate API keys regularly** - Update every 3-6 months
3. **Use HTTPS in production** - Set up SSL/TLS on VM
4. **Enable Vercel authentication** - For issuer portal
5. **Monitor API usage** - Check for unusual patterns
6. **Restrict VM firewall** - Only allow necessary ports
7. **Regular backups** - Automated database backups
8. **Audit logs** - Track who uses admin endpoints

---

## Troubleshooting

### Issue: Student portal can't reach API

**Check**:
1. VM firewall allows port 8080
2. `NEXT_PUBLIC_API_URL` is correct in Vercel
3. CORS allows Vercel domain
4. Backend is running: `docker compose ps`

### Issue: Issuer portal gets 401 Unauthorized

**Check**:
1. API key matches between Vercel and VM
2. API key sent in header: `X-API-Key`
3. Middleware is applied to protected routes
4. Check backend logs for auth failures

### Issue: CORS errors

**Fix**: Update allowed origins in backend
```go
allowedOrigins := []string{
    "https://iu-micert.vercel.app",
    "https://your-issuer-portal.vercel.app",
    "https://*.vercel.app",  // Preview deployments
}
```

### Issue: Vercel deployment fails

**Check**:
1. `package.json` has correct scripts
2. Environment variables set in Vercel
3. Build command: `npm run build`
4. Output directory: `.next` or `out`

---

## Cost Estimation

| Component | Provider | Estimated Cost |
|-----------|----------|----------------|
| Backend VM (e2-medium) | GCP | ~$25-35/month |
| PostgreSQL storage | GCP | ~$5/month |
| Student Portal | Vercel | FREE (Hobby plan) |
| Issuer Portal | Vercel | FREE (Hobby plan) |
| Ethereum gas fees | Sepolia | FREE (testnet) |
| **Total** | | **~$30-40/month** |

---

## Quick Reference Commands

```bash
# Deploy backend
git push origin main

# Deploy student portal
cd packages/client/iumicert-client && vercel --prod

# Deploy issuer portal
cd packages/issuer/web/iumicert-issuer && vercel --prod

# Check backend health
curl http://YOUR_VM_IP:8080/api/health

# View backend logs
ssh VM && docker compose logs -f

# Test admin endpoint
curl -H "X-API-Key: KEY" -X POST http://VM_IP:8080/api/admin/generate-data
```

---

**System Status**:
- ✅ Student Portal: **iu-micert.vercel.app** (deployed)
- ⏳ Issuer Portal: Pending Vercel deployment
- ⏳ Backend API: Pending GCP VM deployment with authorization

**Next Steps**:
1. Deploy backend with API key middleware
2. Deploy issuer portal to Vercel
3. Configure all environment variables
4. Test complete end-to-end flow
