# IU-MiCert Issuer System

Privacy-preserving academic credential management using Verkle trees, zero-knowledge proofs, and Ethereum blockchain integration.

## 🎯 What It Does

Universities can issue tamper-proof, verifiable academic credentials that students control:

- **🔐 Cryptographic Proofs**: 32-byte Verkle proofs for each course completion using Ethereum's go-verkle library
- **🎓 Privacy-Preserving**: Students share only relevant achievements without revealing full transcripts
- **⛓️ Blockchain-Anchored**: Term roots published to Ethereum for independent verification
- **✨ Zero-Knowledge**: Verify credentials without contacting the institution
- **🖥️ Modern Interface**: CLI tools + Next.js web dashboard for complete workflow

## 🚀 Quick Start

```bash
# Complete setup (reset + generate data + start server)
./reset.sh && ./generate.sh && ./dev.sh

# Or manual step-by-step:
./micert generate-data           # Create test data (5 students, 6 terms)
./micert batch-process           # Build Verkle trees
./micert generate-all-receipts   # Create student receipts with proofs
export ISSUER_PRIVATE_KEY="..."  # Set your Sepolia private key
./micert publish-roots           # Publish to blockchain

# Start web interface (in separate terminal)
cd web/iumicert-issuer && npm install && npm run dev
```

**Access**: Backend API at `http://localhost:8080` | Dashboard at `http://localhost:3001`

## 🏗️ Architecture

**Single Verkle Tree per Term** - Each academic term (e.g., Semester_1_2023) becomes one Verkle tree containing all student course completions.

### Core Components

1. **Verkle Tree System** (`crypto/verkle/`)
   - Uses Ethereum's `ethereum/go-verkle` for real cryptographic operations
   - Generates 32-byte IPA (Inner Product Argument) proofs per course
   - Enables selective disclosure: prove specific courses without revealing all

2. **Data Pipeline** (`cmd/`)
   - Generate realistic academic data (Vietnamese students, IU course codes)
   - Convert LMS data to Verkle-compatible format
   - Batch process multiple terms automatically

3. **Blockchain Integration** (`blockchain_integration/`)
   - Smart contracts deployed on Ethereum Sepolia testnet
   - IUMiCertRegistry: `0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60`
   - Publishes term roots for decentralized verification

4. **Interfaces**
   - **CLI** (`cmd/main.go`): 15+ commands for all operations
   - **REST API** (`cmd/api_server.go`): Backend services on port 8080
   - **Web Dashboard** (`web/iumicert-issuer/`): Next.js interface with wallet integration

## 📁 Project Structure

```
packages/issuer/
├── cmd/                      # CLI application & API server
│   ├── main.go              # CLI entry point (15+ commands)
│   ├── api_server.go        # REST API server
│   └── *.go                 # Data generation, conversion, processing
├── crypto/verkle/           # Verkle tree implementation
│   └── term_aggregation.go # Core Verkle operations
├── blockchain_integration/  # Ethereum smart contract interaction
├── data/                    # Generated academic data
│   ├── student_journeys/   # Student academic progression
│   └── verkle_terms/       # Verkle-compatible format
├── publish_ready/           # Output files
│   ├── receipts/           # Student receipts with 32-byte proofs
│   ├── roots/              # Term root commitments
│   └── transactions/       # Blockchain transaction records
├── web/iumicert-issuer/    # Next.js dashboard
│   ├── src/app/            # App router pages
│   ├── src/components/     # UI components (shadcn/ui)
│   └── src/lib/            # Blockchain & API client
├── scripts/                 # Automation
│   ├── reset.sh            # Clean all data
│   ├── generate.sh         # Generate complete dataset
│   └── dev.sh              # Start API server
└── .env                     # Configuration (Sepolia key, contracts)
```

## 🛠️ Prerequisites

- **Go 1.21+** - Backend CLI and API
- **Docker & Docker Compose** - Database container
- **Node.js 18+** - Web dashboard
- **Ethereum Wallet** - Sepolia testnet ETH for publishing

## 📦 Installation

```bash
# Clone and navigate
git clone https://github.com/nikola0x0/iumicert.git
cd iumicert/packages/issuer

# Build CLI
go build -o micert cmd/*.go

# Configure environment
cp .env.example .env
# Edit .env with your ISSUER_PRIVATE_KEY

# Install web dependencies
cd web/iumicert-issuer && npm install && cd ../..
```

## 🎬 Automated Workflows

### Reset System
```bash
./reset.sh  # Clears all data, prepares fresh start
```

### Generate Complete Dataset
```bash
./generate.sh  # Creates 5 students × 6 terms with Verkle trees & receipts
```

### Start Development Server
```bash
./dev.sh  # Runs API server on :8080 with CORS
```

## 🛠️ CLI Commands

| Command | Purpose | Example |
|---------|---------|---------|
| `init` | Initialize system | `./micert init "IU-VNUHCM"` |
| `generate-data` | Create test data | `./micert generate-data --students=5 --terms="Semester_1_2023,Semester_2_2023"` |
| `convert-data` | Convert to Verkle format | `./micert convert-data Semester_1_2023` |
| `add-term` | Add term with Verkle tree | `./micert add-term Semester_1_2023 data.json` |
| `batch-process` | Process all terms | `./micert batch-process` |
| `generate-receipt` | Create receipt | `./micert generate-receipt ITITIU00001` |
| `generate-all-receipts` | Create all receipts | `./micert generate-all-receipts` |
| `generate-addon-term` | Generate single additional term | `./micert generate-addon-term Test_Term` |
| `publish-roots` | Publish to blockchain | `./micert publish-roots Semester_1_2023` |
| `verify-local` | Verify receipt locally | `./micert verify-local receipt.json` |
| `test-verify` | Full cryptographic verification | `./micert test-verify receipt.json` |
| `display-receipt` | Show receipt details | `./micert display-receipt receipt.json` |
| `verification-guide` | Show verification guide | `./micert verification-guide` |
| `serve` | Start API server | `./micert serve --port 8080 --cors` |
| `migrate` | Run database migrations | `./micert migrate` |
| `db-import` | Import data to database | `./micert db-import` |

Run `./micert --help` or `./micert <command> --help` for details.


## ⚙️ Configuration

### Environment Variables

Edit `.env` file:

```bash
# Blockchain
ISSUER_PRIVATE_KEY=your_sepolia_private_key_here
IUMICERT_CONTRACT_ADDRESS=0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60
NETWORK=sepolia
SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161

# Database (PostgreSQL)
DATABASE_URL=postgresql://iumicert:iumicert_secret@localhost:5432/iumicert?sslmode=disable

# Optional
DEFAULT_GAS_LIMIT=500000
MAX_GAS_PRICE=20000000000  # 20 gwei
```

### Database Setup

**PostgreSQL** with GORM ORM via Docker:
```bash
# Start PostgreSQL container
docker-compose up -d postgres

# Verify database is running
docker ps | grep iumicert-postgres

# Migrations run automatically on API server start
./micert serve

# Optional: Access pgAdmin at http://localhost:5050
docker-compose up -d pgadmin
```

### Deployed Contracts

- **IUMiCertRegistry**: `0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60` - Stores term roots
- **ReceiptRevocationRegistry**: `0x8814ae511d54Dc10C088143d86110B9036B3aa92` - Revocation

**Get Sepolia ETH** (~0.01 ETH needed):
- https://cloud.google.com/application/web3/faucet/ethereum/sepolia
- https://sepolia-faucet.pk910.de

## 🎯 Usage

### Recommended: Web UI (Easiest)
```bash
# 1. Start database
docker-compose up -d postgres

# 2. Start backend API
./dev.sh

# 3. Start web dashboard (new terminal)
cd web/iumicert-issuer && npm run dev

# 4. Open http://localhost:3001
# - Generate test data in Demo Data page
# - Publish terms in Publish Terms page (connect wallet)
# - Verify receipts in Verifier page
```

### CLI Alternative
```bash
# Generate data and publish
./reset.sh && ./generate.sh
export ISSUER_PRIVATE_KEY="..."
./micert publish-roots

# Generate and verify receipts
./micert generate-receipt ITITIU00001
./micert verify-local publish_ready/receipts/receipt_ITITIU00001_all_terms.json
```

## 🌐 API Reference

REST API server on port 8080 with endpoint separation:

### Key Endpoints

**Issuer Operations** (`/api/issuer/*`)
- `GET /terms` - List all academic terms
- `GET /terms/{term_id}/roots` - Get term Verkle root
- `POST /receipts` - Generate student receipt
- `POST /blockchain/publish` - Publish term roots
- `GET /students/{student_id}/journey` - Student academic history

**Verifier Operations** (`/api/verifier/*`)
- `POST /receipt` - Verify receipt proofs
- `POST /course` - Verify specific course
- `GET /blockchain/transaction/{tx_hash}` - Check transaction

**System**
- `GET /health` - Health check
- `GET /status` - System status

### Test Data

Generated dataset includes:
- **Students**: ITITIU00001-ITITIU00005 (Vietnamese names)
- **Terms**: 6 semesters (Semester_1_2023 through Summer_2024)
- **Courses**: Real IU Vietnam codes (IT013IU, IT153IU, PE008IU, MA001IU, etc.)

## 🌐 Web Dashboard

Modern Next.js interface at `http://localhost:3001` - **Full-featured alternative to CLI**

### Setup
```bash
cd web/iumicert-issuer
npm install && npm run dev
```

### Features (All CLI Operations Available)

**1. Publish Terms (`/`)**
- Select and publish term roots to Sepolia blockchain
- Connect MetaMask wallet, select term, approve transaction
- Real-time transaction status with Etherscan links
- Automatic database updates after publishing

**2. Verifier (`/verifier`)**
- Upload and verify student receipt JSON files
- View complete academic journey with courses, grades, and cryptographic proofs
- Blockchain verification status for each term
- Expandable course details with 32-byte IPA proofs

**3. Demo Data (`/demo-data`)** - *Replaces CLI data generation*
- **Generate Test Data**: 1-100 students × customizable terms
  - Calls backend to generate student journeys
  - Automatically converts to Verkle format
  - Builds Verkle trees with cryptographic proofs
  - Generates receipts for all students
- **Add Custom Terms**: Create new academic terms dynamically
- **System Reset**: Clean all data (equivalent to `./reset.sh`)

## 🔐 Security Features

- **32-byte Verkle Proofs**: Real cryptographic proofs using Ethereum's go-verkle library
- **Zero-Knowledge Verification**: Prove credentials without revealing full transcript
- **Selective Disclosure**: Students receive full receipts and can create filtered versions showing only specific courses/terms while maintaining cryptographic verifiability
- **Blockchain Anchoring**: Immutable term roots on Ethereum for independent verification
- **Tamper-Proof**: Cryptographically impossible to forge credentials

### How Selective Disclosure Works

The issuer provides students with **complete receipts** containing all courses and cryptographic proofs. Students can then create **filtered versions** for specific purposes:

**Workflow:**
1. **Issuer generates full receipt** - Contains all terms and courses with individual Verkle proofs
2. **Student receives full receipt** - `ITITIU00001_full_journey.json`
3. **Student creates selective receipt** - Removes unwanted courses/terms from the JSON
4. **Verifier validates selective receipt** - Each remaining course proof still verifies against the blockchain root

**Why it's secure:**
- Each course has an **independent 32-byte cryptographic proof**
- Removing courses from the receipt doesn't affect verification of remaining courses
- The **Verkle root remains unchanged** (published on blockchain)
- Verifiers confirm revealed courses are genuine without seeing hidden ones
- **Impossible to add fake courses** - each proof must cryptographically verify against the published root

**Example:** Student has 30 courses across 6 terms but only wants to show IT courses from 2023:
```json
// Student manually filters their full receipt JSON:
// 1. Remove unwanted terms (e.g., keep only Semester_1_2023, Semester_2_2023)
// 2. Remove unwanted courses from "revealed_courses" array
// 3. Remove corresponding entries from "course_proofs" object
// 4. Keep "verkle_root" unchanged

// Result: Selective receipt with only 5 IT courses that still verifies!
```

**Note:** Currently, filtering is done manually by editing the JSON. A student/verifier portal with UI-based filtering is planned for easier selective disclosure.

## 🎯 Use Cases

This is the **Issuer System** for universities and academic institutions to:

- **Generate Academic Credentials**: Create verifiable receipts with cryptographic proofs for students
- **Publish Term Roots**: Anchor academic term data to Ethereum blockchain for independent verification
- **Manage Student Data**: Generate and process student academic journeys with Verkle trees
- **Issue Portable Credentials**: Students receive JSON receipts they can share with anyone

**Note**: The Verifier route (`/verifier`) in the web UI is for testing purposes only. In production, third-party verifiers (employers, other institutions) would verify credentials independently using the CLI verification commands or their own verification systems, without needing access to the issuer's interface.

## 🔧 Troubleshooting

### Common Issues

**Port already in use**
```bash
# Option 1: Use pkill (Linux/macOS)
pkill -f "micert serve"

# Option 2: Find and kill by port
# Linux/macOS:
lsof -ti:8080 | xargs kill -9
# Windows (PowerShell):
# Get-Process -Id (Get-NetTCPConnection -LocalPort 8080).OwningProcess | Stop-Process

# Then restart
./dev.sh
```

**Blockchain publishing fails**
- Check `ISSUER_PRIVATE_KEY` is set in `.env`
- Verify Sepolia ETH balance (~0.01 ETH needed)
- Ensure wallet is on Sepolia network

**Frontend can't connect to API**
- Verify API server is running on port 8080
- Check CORS is enabled: `./micert serve --cors`

**Receipt verification fails**
- Ensure term has been published to blockchain first
- Check receipt JSON file is valid

### Development Tips

- Rebuild CLI after code changes: `go build -o micert ./cmd`
- Test API directly: `curl localhost:8080/api/health`
- View logs in `./dev.sh` terminal for debugging


## 📚 Additional Resources

- Run `./micert verification-guide` for third-party verifier instructions
- View smart contract on Etherscan: [0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60](https://sepolia.etherscan.io/address/0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60)
- Contract ABIs available in `blockchain_integration/abi/`

## 🙏 Credits

Built with [Ethereum's go-verkle](https://github.com/ethereum/go-verkle), Cobra CLI, and Next.js

---

**IU-MiCert Issuer System** - Privacy-preserving academic credentials with blockchain technology
