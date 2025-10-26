# Quick Deployment Reference

## ✅ Works in Both Environments

Your setup automatically handles:
- **Local Development**: Uses `http://localhost:8080`
- **Production (Vercel)**: Uses `https://api.yourdomain.com`

## 🚀 Quick Deploy Steps

### 1. Deploy Frontend to Vercel (5 minutes)

```bash
cd packages/issuer/web/iumicert-issuer
vercel
```

Or use Vercel Dashboard: https://vercel.com/new

**Environment Variables to Set:**
```
NEXT_PUBLIC_API_URL=https://api.yourdomain.com
NEXT_PUBLIC_IUMICERT_CONTRACT_ADDRESS=0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60
NEXT_PUBLIC_ANALYTICS_DISABLED=true
```

### 2. Update Backend on VM

```bash
# SSH to your VM
ssh user@your-vm

# Add to .env file
echo "FRONTEND_URL=https://your-app.vercel.app" >> .env

# Restart
docker-compose restart issuer-backend
```

### 3. Done! 🎉

Visit: `https://your-app.vercel.app`

## 💰 Cost

- **Frontend (Vercel)**: $0/month (free tier)
- **Backend (VM)**: $5-10/month
- **Total**: ~$5-10/month

## 🧪 Testing Locally

```bash
# Terminal 1: Backend
./dev.sh

# Terminal 2: Frontend
cd web/iumicert-issuer
npm run dev
```

Visit: `http://localhost:3000`

## 🔧 Update Deployments

**Frontend**: Just push to Git
```bash
git push origin main  # Auto-deploys on Vercel
```

**Backend**: SSH and restart
```bash
cd /path/to/iumicert
git pull
docker-compose restart issuer-backend
```

## 📚 Full Guide

See `VERCEL-DEPLOYMENT.md` for detailed instructions.
