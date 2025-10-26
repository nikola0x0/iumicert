# Vercel Deployment Guide for IU-MiCert Web Portal

## Architecture Overview

```
┌─────────────────────────────────────────┐
│         Vercel (Global CDN)             │
│  web/iumicert-issuer (Next.js Frontend) │
│         https://your-app.vercel.app     │
└──────────────┬──────────────────────────┘
               │
               │ API Calls over HTTPS
               ▼
┌─────────────────────────────────────────┐
│     Your VM (Backend Infrastructure)    │
│  https://api.yourdomain.com             │
│  ┌───────────────────────────────────┐  │
│  │  Go Backend API (:8080)           │  │
│  │  - Verkle operations              │  │
│  │  - Blockchain integration         │  │
│  └──────────┬────────────────────────┘  │
│             │                            │
│  ┌──────────▼────────────────────────┐  │
│  │  PostgreSQL (:5432)               │  │
│  │  - Academic records               │  │
│  └───────────────────────────────────┘  │
└─────────────────────────────────────────┘
```

## Prerequisites

- Vercel account (free tier works)
- Your VM with backend deployed
- Domain or use Vercel's provided URL

## Step 1: Prepare Backend on VM

### 1.1 Set up SSL/TLS for your backend

Your backend API needs HTTPS. Options:

**Option A: Use Nginx as reverse proxy (Recommended)**
```bash
# Install Nginx
sudo apt update && sudo apt install nginx certbot python3-certbot-nginx

# Configure Nginx
sudo nano /etc/nginx/sites-available/iumicert-api

# Add configuration:
server {
    listen 80;
    server_name api.yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# Enable site
sudo ln -s /etc/nginx/sites-available/iumicert-api /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx

# Get SSL certificate
sudo certbot --nginx -d api.yourdomain.com
```

**Option B: Use Cloudflare Tunnel (No domain needed)**
```bash
# Install cloudflared
wget https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb
sudo dpkg -i cloudflared-linux-amd64.deb

# Login and create tunnel
cloudflared tunnel login
cloudflared tunnel create iumicert-api
cloudflared tunnel route dns iumicert-api api.yourdomain.com

# Run tunnel
cloudflared tunnel run iumicert-api --url http://localhost:8080
```

### 1.2 Update Backend Environment

Edit your `.env` file on the VM:
```bash
# Add your Vercel deployment URL
FRONTEND_URL=https://your-app.vercel.app
```

### 1.3 Restart Backend
```bash
docker-compose down
docker-compose up -d
```

## Step 2: Deploy to Vercel

### 2.1 Install Vercel CLI (Optional)
```bash
npm i -g vercel
```

### 2.2 Deploy via Vercel Dashboard (Easiest)

1. Go to https://vercel.com/new
2. Import your Git repository
3. Configure project:
   - **Framework Preset**: Next.js
   - **Root Directory**: `packages/issuer/web/iumicert-issuer`
   - **Build Command**: `npm run build` (default)
   - **Output Directory**: `.next` (default)

4. Add Environment Variables:
   ```
   NEXT_PUBLIC_IUMICERT_CONTRACT_ADDRESS=0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60
   NEXT_PUBLIC_API_URL=https://api.yourdomain.com
   NEXT_PUBLIC_ANALYTICS_DISABLED=true
   ```

5. Click "Deploy"

### 2.3 Deploy via CLI

```bash
# Navigate to Next.js app
cd packages/issuer/web/iumicert-issuer

# Login to Vercel
vercel login

# Deploy
vercel

# Follow prompts:
# - Link to existing project or create new one
# - Confirm settings
# - Deploy!

# For production deployment
vercel --prod
```

## Step 3: Configure Environment Variables

### On Vercel (Frontend):
```env
NEXT_PUBLIC_IUMICERT_CONTRACT_ADDRESS=0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60
NEXT_PUBLIC_API_URL=https://api.yourdomain.com
NEXT_PUBLIC_ANALYTICS_DISABLED=true
```

### On VM (Backend):
```env
FRONTEND_URL=https://your-app.vercel.app
```

## Step 4: Verify Deployment

### 4.1 Check Frontend
Visit: `https://your-app.vercel.app`

### 4.2 Test API Connection
Open browser console and check for CORS errors. Should see successful API calls.

### 4.3 Test Features
- Generate receipt
- Verify receipt
- Check blockchain integration

## Local Development Setup

Works exactly the same as before:

```bash
# Terminal 1: Start backend
cd packages/issuer
./dev.sh

# Terminal 2: Start frontend
cd packages/issuer/web/iumicert-issuer
npm run dev
```

Frontend uses `http://localhost:8080` automatically.

## Troubleshooting

### CORS Errors
```bash
# Check backend logs
docker logs iumicert-issuer

# Verify FRONTEND_URL is set correctly
docker exec iumicert-issuer env | grep FRONTEND_URL

# Restart backend
docker-compose restart issuer-backend
```

### API Connection Fails
- Ensure backend has HTTPS enabled
- Check firewall allows port 443
- Verify DNS resolves correctly: `nslookup api.yourdomain.com`
- Test backend directly: `curl https://api.yourdomain.com/api/health`

### Environment Variables Not Loading
- Redeploy on Vercel after changing env vars
- Check Vercel dashboard > Settings > Environment Variables
- Verify `NEXT_PUBLIC_` prefix for client-side variables

## Cost Breakdown

### Vercel (Frontend)
- **Free Tier**: 
  - 100GB bandwidth/month
  - Unlimited deployments
  - Automatic HTTPS
  - Global CDN
  - **Cost: $0/month**

### VM (Backend + Database)
- DigitalOcean/Linode/AWS: $5-10/month (basic droplet)
- Backend storage: Minimal (< 1GB)
- Database: PostgreSQL in Docker (no extra cost)
- **Total: $5-10/month**

### Total Monthly Cost: **~$5-10**

## Updating Deployments

### Frontend (Vercel)
Automatic! Push to Git and Vercel auto-deploys.

```bash
git push origin main
# Vercel automatically deploys
```

### Backend (VM)
```bash
# SSH to VM
ssh user@your-vm

# Pull latest changes
cd /path/to/iumicert
git pull

# Rebuild and restart
docker-compose down
docker-compose up -d --build
```

## Security Checklist

- [ ] Backend uses HTTPS (not HTTP)
- [ ] CORS configured with specific origins (not `*`)
- [ ] Private keys in environment variables (never in code)
- [ ] Firewall configured (only ports 80, 443, 22 open)
- [ ] Database not exposed publicly (only internal Docker network)
- [ ] Regular security updates on VM

## Custom Domain (Optional)

### For Vercel:
1. Go to Vercel Dashboard > Settings > Domains
2. Add custom domain (e.g., `issuer.iumicert.com`)
3. Update DNS records as instructed
4. Update backend `FRONTEND_URL` to new domain

### For Backend:
1. Point A record to your VM IP
2. Update Nginx/Cloudflare config with new domain
3. Get new SSL certificate: `sudo certbot --nginx -d api.iumicert.com`

## Monitoring

### Vercel Analytics (Free)
- Automatic performance monitoring
- Web Vitals tracking
- View in Vercel Dashboard

### Backend Monitoring
```bash
# Check logs
docker logs -f iumicert-issuer

# Check health
curl https://api.yourdomain.com/api/health

# Check resources
docker stats iumicert-issuer
```

## Support

- Vercel Documentation: https://vercel.com/docs
- IU-MiCert Issues: https://github.com/yourusername/iumicert/issues
