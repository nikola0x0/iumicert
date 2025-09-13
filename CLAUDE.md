# IU-MiCert Issuer System - Project Overview for Claude

## 🎯 Project Purpose
Academic credential management system for universities using Verkle trees with blockchain integration and zero-knowledge proofs. Enables privacy-preserving verification of student achievements.

## 🏗️ Architecture Overview

### Core Technology Stack
- **Backend**: Go 1.21+ with Cobra CLI framework
- **Cryptography**: Ethereum's go-verkle library for 32-byte proofs
- **Blockchain**: Ethereum Sepolia testnet for root commitments
- **Frontend**: Next.js web interface (in `web/iumicert-issuer/`)
- **Smart Contracts**: IUMiCertRegistry (0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60)

### Key Components
1. **Verkle Trees**: One tree per academic term containing all student course completions
2. **CLI Application**: `micert` binary with comprehensive commands
3. **REST API**: Server on port 8080 with CORS support
4. **Data Pipeline**: LMS → Verkle format → Trees → Receipts → Blockchain

## 📁 Critical Directories
```
packages/issuer/
├── cmd/                    # CLI commands (main.go, api_server.go, etc.)
├── data/                   # Academic data storage
│   ├── student_journeys/   # Generated test data
│   └── verkle_terms/       # Converted Verkle format
├── publish_ready/          # Blockchain-ready outputs
│   ├── receipts/          # Student receipts with proofs
│   ├── roots/             # Term root commitments
│   └── transactions/      # Blockchain records
├── crypto/verkle/         # Verkle tree implementation
│   └── term_aggregation.go # Core Verkle operations
└── blockchain_integration/ # Ethereum integration
```

## 🚀 Quick Commands

### Setup & Generation
```bash
# Complete reset and regenerate
./reset.sh && ./generate.sh

# Start development server
./dev.sh

# Manual process
./micert generate-data        # Create test data
./micert batch-process        # Process all terms
./micert generate-all-receipts # Create receipts
```

### Blockchain Operations
```bash
# Requires ISSUER_PRIVATE_KEY env var
export ISSUER_PRIVATE_KEY="your_sepolia_private_key"
./micert publish-roots        # Publish to blockchain
```

### Verification
```bash
./micert verify-local <receipt.json>  # Local verification
./micert display-receipt <receipt.json> # Human-readable display
```

## 🔑 Key Concepts

### Verkle Trees
- **Single tree per term**: Each academic term (e.g., Semester_1_2023) becomes one Verkle tree
- **32-byte proofs**: Compact cryptographic proofs for each course completion
- **Selective disclosure**: Students can prove specific courses without revealing full transcript

### Academic Data Flow
1. **Generate**: Create realistic student journeys (5 students, 6 terms)
2. **Convert**: Transform to Verkle-compatible format
3. **Process**: Build Verkle trees with cryptographic commitments
4. **Receipt**: Generate verifiable credentials with proofs
5. **Publish**: Anchor term roots to Ethereum blockchain

### Privacy Features
- Zero-knowledge proofs for credential verification
- Selective course/term disclosure
- No need to contact institution for verification
- Tamper-proof cryptographic integrity

## 🛠️ Development Tips

### Environment Setup
```bash
# Copy and configure
cp .env.example .env
# Add Sepolia private key to .env
```

### Testing Workflow
1. `./reset.sh` - Clean all data
2. `./generate.sh` - Create complete dataset
3. `./dev.sh` - Start API server
4. Test specific operations with `./micert <command>`

### Common Issues
- **No private key**: Set `ISSUER_PRIVATE_KEY` environment variable
- **Gas issues**: Adjust gas limits in `.env` or use `--gas-limit` flag
- **Data missing**: Run `./generate.sh` to create test data

## 📊 Data Model

### Students
- IDs: ITITIU00001-ITITIU00005 (expandable)
- Vietnamese names and IU course codes
- 6 terms of progression (2023-2024 academic year)

### Terms
- Semester_1_2023, Semester_2_2023, Summer_2023
- Semester_1_2024, Semester_2_2024, Summer_2024

### Courses
- Real IU Vietnam codes: IT013IU, IT153IU, PE008IU, MA001IU, etc.
- Natural prerequisite chains and academic progression

## 🌐 API Endpoints

Key endpoints (all prefixed with `/api/`):
- `GET /terms` - List all terms
- `POST /receipts` - Generate receipt
- `POST /receipts/verify` - Verify receipt
- `POST /blockchain/publish` - Publish to blockchain
- `GET /students/{id}/journey` - Student academic history

## 📝 Important Files

- `cmd/main.go` - CLI entry point with all commands
- `crypto/verkle/term_aggregation.go` - Core Verkle implementation
- `blockchain_integration/integration.go` - Ethereum publishing
- `cmd/api_server.go` - REST API implementation
- `config/micert.json` - System configuration

## 🔐 Security Notes

- Private keys should NEVER be committed
- Use environment variables for sensitive data
- Verkle proofs are cryptographically secure (32-byte)
- Blockchain anchoring provides immutability

## 🎯 Current Status

- Branch: `issuer_v2` 
- Architecture: Single Verkle tree per term (simplified from earlier multi-tree design)
- All CLI commands fully implemented and tested
- Deployed contracts on Sepolia testnet
- Web interface available in `web/iumicert-issuer/`

## 🚨 Quick Troubleshooting

```bash
# Check system status
./micert version

# View help for any command
./micert <command> --help

# Check if data exists
ls data/student_journeys/students/

# Verify API is running
curl http://localhost:8080/api/health
```

---
**Last Updated**: Project uses single Verkle tree architecture with full CLI implementation.

## Task Master AI Instructions
**Import Task Master's development workflow commands and guidelines, treat as if import is in the main CLAUDE.md file.**
@./.taskmaster/CLAUDE.md
