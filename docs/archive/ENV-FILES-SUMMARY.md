# Environment Files Summary

## 📁 All Environment Files in Project

### 1. **Root Level**
- `.env.example` - TaskMaster AI API keys (not related to deployment)

### 2. **Backend API** (`packages/issuer/`)
| File | Purpose | Status |
|------|---------|--------|
| `.env` | Local development | ✅ Updated with CORS URLs |
| `.env.example` | Template for developers | ✅ Updated with CORS URLs |
| `.env.production` | Production template | ✅ Updated with CORS URLs |

### 3. **Issuer Portal** (`packages/issuer/web/iumicert-issuer/`)
| File | Purpose | Status |
|------|---------|--------|
| `.env.local` | Local development | ✅ Updated with login credentials |
| `.env.example` | Template | ✅ Already has login credentials |

### 4. **Student Portal** (`packages/client/iumicert-client/`)
| File | Purpose | Status |
|------|---------|--------|
| `.env.local` | Local development | ✅ Already configured |
| `.env.example` | Template | ✅ Already configured |

### 5. **Smart Contracts** (`packages/contracts/`)
| File | Purpose | Status |
|------|---------|--------|
| `.env` | Deployment config | ✅ Has private key & RPC URL |

---

## 🔑 Complete Variable Reference

### Backend API Variables

#### **`.env` (Local Development)**
```bash
# Blockchain
ISSUER_PRIVATE_KEY=ed335d49f5b8f5017315a7bb55da31b10e7859d27864e22a08b7d9375024343e
IUMICERT_CONTRACT_ADDRESS=0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60
NETWORK=sepolia
SEPOLIA_RPC_URL=https://sepolia.drpc.org

# Gas
DEFAULT_GAS_LIMIT=500000
MAX_GAS_PRICE=20000000000

# App
DEBUG=false
LOG_LEVEL=info

# Database
POSTGRES_PASSWORD=iumicert_secret
DATABASE_URL=postgresql://iumicert:iumicert_secret@localhost:5432/iumicert?sslmode=disable
DB_MAX_CONNECTIONS=100
DB_MAX_IDLE=10

# CORS (Local)
STUDENT_PORTAL_URL=http://localhost:3000
ISSUER_PORTAL_URL=http://localhost:3001
```

#### **`.env.production` (Production Template)**
```bash
# Blockchain
ISSUER_PRIVATE_KEY=YOUR_SEPOLIA_PRIVATE_KEY_WITHOUT_0x_PREFIX
IUMICERT_CONTRACT_ADDRESS=0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60
NETWORK=sepolia
SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/YOUR_INFURA_API_KEY

# Gas
DEFAULT_GAS_LIMIT=500000
MAX_GAS_PRICE=20000000000

# Database
POSTGRES_PASSWORD=CHANGE_ME_STRONG_PASSWORD

# App
DEBUG=false
LOG_LEVEL=info

# CORS (Production)
STUDENT_PORTAL_URL=https://iu-micert.vercel.app
ISSUER_PORTAL_URL=https://your-issuer-portal.vercel.app
```

---

### Issuer Portal (Admin Dashboard) Variables

#### **`.env.local`**
```bash
# API
NEXT_PUBLIC_API_URL=http://localhost:8080

# Contract
NEXT_PUBLIC_IUMICERT_CONTRACT_ADDRESS=0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60

# Authentication
NEXT_PUBLIC_ISSUER_USERNAME=admin
NEXT_PUBLIC_ISSUER_PASSWORD=admin123

# Analytics
NEXT_PUBLIC_ANALYTICS_DISABLED=true
```

#### **For Vercel Production** (Set in Dashboard)
```bash
NEXT_PUBLIC_API_URL=http://YOUR_VM_IP:8080
NEXT_PUBLIC_IUMICERT_CONTRACT_ADDRESS=0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60
NEXT_PUBLIC_ISSUER_USERNAME=admin
NEXT_PUBLIC_ISSUER_PASSWORD=your_secure_password
NEXT_PUBLIC_ANALYTICS_DISABLED=true
```

---

### Student Portal Variables

#### **`.env.local`**
```bash
NEXT_PUBLIC_API_URL=http://localhost:8080
```

#### **For Vercel Production** (Set in Dashboard)
```bash
NEXT_PUBLIC_API_URL=http://YOUR_VM_IP:8080
```

---

### Smart Contracts Variables

#### **`.env`**
```bash
PRIVATE_KEY=0xed335d49f5b8f5017315a7bb55da31b10e7859d27864e22a08b7d9375024343e
SEPOLIA_RPC_URL=https://sepolia.drpc.org/
ETHERSCAN_API_KEY=
```

---

## 🎯 What Changed

### ✅ Updated Files

1. **`packages/issuer/.env`**
   - ➕ Added `STUDENT_PORTAL_URL=http://localhost:3000`
   - ➕ Added `ISSUER_PORTAL_URL=http://localhost:3001`

2. **`packages/issuer/.env.example`**
   - ➕ Added database configuration
   - ➕ Added `STUDENT_PORTAL_URL` and `ISSUER_PORTAL_URL`

3. **`packages/issuer/.env.production`**
   - ➕ Added `STUDENT_PORTAL_URL=https://iu-micert.vercel.app`
   - ➕ Added `ISSUER_PORTAL_URL=https://your-issuer-portal.vercel.app`

4. **`packages/issuer/web/iumicert-issuer/.env.local`**
   - ➕ Added `NEXT_PUBLIC_ISSUER_USERNAME=admin`
   - ➕ Added `NEXT_PUBLIC_ISSUER_PASSWORD=admin123`

---

## 📋 GitHub Secrets Needed

For CI/CD deployment, add these in **GitHub Repository Settings → Secrets**:

### Backend & Infrastructure
- `GCP_SA_KEY` - Service account JSON
- `GCP_VM_NAME` - VM instance name
- `GCP_VM_ZONE` - VM zone
- `POSTGRES_PASSWORD` - Database password
- `ISSUER_PRIVATE_KEY` - Ethereum private key
- `SEPOLIA_RPC_URL` - RPC endpoint URL
- `IUMICERT_CONTRACT_ADDRESS` - Contract address

### Frontend URLs (NEW)
- `STUDENT_PORTAL_URL` - `https://iu-micert.vercel.app`
- `ISSUER_PORTAL_URL` - Your issuer Vercel URL

---

## 🔍 Files NOT Modified

These files are fine as-is:
- ✅ `packages/client/iumicert-client/.env.local` - Already correct
- ✅ `packages/client/iumicert-client/.env.example` - Already correct
- ✅ `packages/contracts/.env` - Contract deployment only
- ✅ `packages/issuer/web/iumicert-issuer/.env.example` - Already has auth

---

## 🚀 Next Steps

1. **Review all variables** above for correctness
2. **Update production values** in `.env.production` if needed
3. **Add GitHub secrets** (2 new ones for portal URLs)
4. **Deploy issuer portal to Vercel** to get its URL
5. **Update `ISSUER_PORTAL_URL`** everywhere with actual Vercel URL
6. **Commit and push** to trigger deployment

---

## ⚠️ Security Notes

**Never commit these to git**:
- Actual private keys (use `.gitignore`)
- Real database passwords
- Production RPC URLs with API keys

**Safe to commit**:
- `.env.example` files (templates)
- Local development placeholder values
- Public contract addresses

---

**All environment files are now consistent and ready for deployment!** ✅
