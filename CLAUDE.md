# IU-MiCert Project - AI Assistant Instructions

## 🎯 Project Overview

Academic credential management system using Verkle trees, blockchain integration, and zero-knowledge proofs for privacy-preserving verification of student achievements.

**Key Concept**: Each course is treated as a micro-credential. The system enables term-by-term verification and maintains a complete, tamper-proof academic provenance timeline that cannot be backdated or manipulated.

## 🏗️ System Architecture

### Technology Stack

- **Backend**: Go 1.21+ (packages/issuer/)
- **Cryptography**: Ethereum's go-verkle library
- **Blockchain**: Ethereum Sepolia testnet
- **Frontend**: Next.js (issuer dashboard + student/verifier portal)
- **Smart Contracts**: Solidity (IUMiCertRegistry, ReceiptRevocationRegistry)

### Core Components

1. **Verkle Trees**: One tree per academic term containing all student course completions
2. **CLI Application**: `micert` binary with 15+ commands (see issuer README)
3. **REST API**: Server on port 8080 with CORS support
4. **Data Pipeline**: LMS data → Verkle format → Trees → Receipts → Blockchain

## 📁 Project Structure

```
packages/
├── issuer/          # Go backend + CLI + API + admin dashboard
│   ├── cmd/         # CLI commands & API server
│   ├── crypto/      # Verkle tree implementation
│   ├── data/        # Test data generation
│   └── web/         # Issuer dashboard (Next.js)
├── client/          # Student/Verifier portal (Next.js)
└── contracts/       # Smart contracts (Solidity + Foundry)
docs/                # Technical documentation
```

## 🚀 Quick Reference

### Common Commands

```bash
# Complete reset and regenerate
./reset.sh && ./generate.sh

# Start API server
./dev.sh

# Individual operations
./micert generate-data
./micert batch-process
./micert generate-all-receipts
./micert publish-roots  # Requires ISSUER_PRIVATE_KEY
```

### Key Directories in packages/issuer/

- `data/student_journeys/` - Generated test data
- `data/verkle_terms/` - Converted Verkle format
- `publish_ready/receipts/` - Student receipts with proofs
- `publish_ready/roots/` - Term root commitments
- `publish_ready/transactions/` - Blockchain records

## 🔑 Key Concepts

### Academic Provenance

- Each **course** = independent verifiable micro-credential
- Each **term** = one Verkle tree with all course completions
- **Timeline integrity**: Cannot backdate or manipulate achievement order
- **Selective disclosure**: Students can share specific courses without revealing full transcript

### Verkle Trees

- **32-byte proofs**: Compact cryptographic proofs per course
- **Term-based**: One tree per academic term (e.g., Semester_1_2023)
- **Blockchain anchoring**: Term roots published to Ethereum for independent verification

## 📊 Test Data Model

- **Students**: ITITIU00001-ITITIU00005 (Vietnamese names, IU course codes)
- **Terms**: 6 semesters (Semester_1_2023 through Summer_2024)
- **Courses**: Real IU Vietnam codes (IT013IU, IT153IU, PE008IU, etc.)

## 🌐 Deployed System

**Smart Contracts (Sepolia):**

- **IUMiCertRegistry**: `0x2452F0063c600BcFc232cC9daFc48B7372472f79` ✅ **ACTIVE** (v2 with versioning)
  - Enhanced with term root versioning for revocation support
  - Deployed: 2025-11-27
  - Etherscan: https://sepolia.etherscan.io/address/0x2452f0063c600bcfc232cc9dafc48b7372472f79
- ~~IUMiCertRegistry (v1)~~: `0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60` (Legacy)
- ~~ReceiptRevocationRegistry~~: `0x8814ae511d54Dc10C088143d86110B9036B3aa92` (Deprecated - superseded by versioning approach)

**Web Applications:**

- Student/Verifier Portal: https://iu-micert.vercel.app
- Issuer Dashboard: https://iumicert-issuer.vercel.app

## 📚 Documentation

- **Issuer System Details**: `packages/issuer/README.md` (comprehensive guide)
- **Technical Docs**: `docs/` directory
  - IPA_VERIFICATION_IMPLEMENTATION.md - Cryptographic implementation
  - VERKLE_MEMBERSHIP_PROOFS.md - Historical context
  - THESIS_DEFENSE_SCRIPT.md - Presentation guide
- **Client Portal**: `packages/client/iumicert-client/README.md`
- **Smart Contracts**: `packages/contracts/README.md`

## ⚙️ Configuration

- **Environment**: `.env` files in each package (see `.env.example` templates)
- **Private Keys**: Required for blockchain publishing (ISSUER_PRIVATE_KEY)
- **API Server**: Runs on localhost:8080 by default
- **CORS**: Configured for both local development and Vercel deployments

## 🎯 Current State

- ✅ All CLI commands implemented and tested
- ✅ Verkle tree architecture: single tree per term (simplified design)
- ✅ Smart contracts deployed on Sepolia testnet
- ✅ Web interfaces deployed (issuer + student/verifier portals)
- ✅ Full IPA verification with cryptographic binding
- 🔄 **Future enhancements**: Receipt revocation functionality

## 💡 Important Notes for AI Assistants

1. **Read issuer README first** for detailed operations: `packages/issuer/README.md`
2. **This is a thesis project** - focus on research contributions (provenance, temporal integrity)
3. **Test data is realistic** - uses actual IU Vietnam course codes and Vietnamese names
4. **Revocation is planned but not critical** - system is functional without it
5. **Documentation cleanup** - archive old deployment-specific docs, keep technical references

---
