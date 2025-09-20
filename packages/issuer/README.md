# IU-MiCert Issuer System

A comprehensive academic credential management system for educational institutions using Verkle tree architecture with blockchain integration and zero-knowledge proofs.

## 📋 Current System Status

✅ **Fully Operational** - All core features implemented and tested  
✅ **API Restructured** - Issuer/Verifier endpoint separation complete  
✅ **Frontend Fixed** - Blockchain publishing shows success instead of 500 errors  
✅ **Duplicate Prevention** - Smart contract rejects already-published terms gracefully  
🔧 **Ready for Extensions** - Verifier interface and authentication middleware ready to implement  

**Last Updated**: September 2025 - System running on `issuer_v2` branch with single Verkle tree architecture

## 🎓 Overview

The IU-MiCert Issuer System enables universities to:

- **Process Academic Data**: Convert LMS data into cryptographic Verkle tree structures
- **Generate Zero-Knowledge Proofs**: Create 32-byte Verkle proofs for student achievements
- **Issue Verifiable Receipts**: Generate privacy-preserving academic journey receipts
- **Blockchain Integration**: Publish term root commitments to Ethereum networks
- **Selective Di**Note\*\*: All commands listed above are fully implemented and tested. This reflects the current state of the system using single Verkle tree architecture.

## 🚀 Quick Start Workflow

For new users, here's the typical workflow to get started:

### 1. Fresh Setup (Complete Pipeline)

```bash
# Reset everything and generate fresh data
./reset.sh && ./generate.sh

# Start development server
./dev.sh
```

### 2. Manual Step-by-Step Process

```bash
# 1. Generate academic data
./micert generate-data

# 2. Convert and process all terms
./micert batch-process

# 3. Generate receipts for all students
./micert generate-all-receipts

# 4. Publish term roots to blockchain (requires .env setup)
./micert publish-roots

# 5. Verify a receipt
./micert verify-local publish_ready/receipts/receipt_ITITIU00001_all_terms.json
```

### 3. Environment Configuration

```bash
# Copy and configure environment
cp .env.example .env
# Edit .env with your Sepolia private key

# Test blockchain connection
./micert publish-roots --help
```

## 📚 Documentationosure\*\*: Enable students to prove specific achievements without revealing all data

## 🏗️ Architecture

### Core Components

1. **Verkle Tree System**

   - **Single Verkle Trees per Term**: Each academic term becomes one Verkle tree
   - **Course-Level Proofs**: 32-byte proofs for individual course completions
   - **Ethereum Integration**: Uses `ethereum/go-verkle` library for real cryptographic operations
   - **Selective Disclosure**: Students can prove specific courses without revealing entire transcript

2. **Academic Data Pipeline**

   - **LMS Data Processing**: Converts university data into Verkle-compatible format
   - **Student Journey Generation**: Creates realistic academic progression data
   - **Term Aggregation**: Combines all course completions into single Verkle trees

3. **Blockchain Integration**

   - **Smart Contracts**: Deployed on Sepolia testnet for term root storage
   - **On-Chain Verification**: Enables anyone to verify academic credentials
   - **Transaction Monitoring**: Tracks blockchain confirmations and gas costs

4. **CLI & Web Interface**
   - **Command-Line Tools**: Complete CLI for data generation, processing, and verification
   - **Web Dashboard**: Next.js interface for student journey visualization
   - **REST API**: Backend services for receipt generation and verification

## 📁 Project Structure

```
packages/issuer/
├── cmd/                           # CLI application
│   ├── main.go                   # Main CLI entry point with all commands
│   ├── api_server.go            # REST API server
│   ├── data_generator.go        # Realistic test data generator
│   ├── data_converter.go        # LMS data converter
│   └── batch_processor.go       # Batch processing pipeline
├── data/                         # Academic data storage
│   ├── student_journeys/         # Generated student academic journeys
│   │   ├── students/            # Individual student progression files
│   │   └── terms/               # Term summary data
│   └── verkle_terms/            # Term data converted for Verkle trees
├── publish_ready/                # Blockchain-ready outputs
│   ├── receipts/                # Student academic receipts with proofs
│   ├── roots/                   # Term Verkle root commitments
│   ├── proofs/                  # 32-byte cryptographic proofs
│   └── transactions/            # Blockchain transaction records
├── blockchain_integration/       # Ethereum integration
│   ├── client.go                # Blockchain client setup
│   ├── contracts.go             # Smart contract interactions
│   └── integration.go           # Publishing and verification logic
├── crypto/verkle/                # Verkle tree implementation
│   └── term_aggregation.go      # Core Verkle tree operations
├── config/                       # System configuration
│   ├── env.go                   # Environment variable management
│   └── micert.json              # Application settings
├── web/iumicert-issuer/         # Next.js web interface
├── scripts/                     # Automation scripts
│   ├── reset.sh                 # Clean all data and reset system
│   ├── generate.sh              # Generate complete dataset
│   └── dev.sh                   # Start development server
└── .env                         # Environment configuration
```

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+** (for backend CLI and API)
- **Node.js 18+** (for frontend web interface)
- **Git** (for version control)
- **Ethereum Wallet** (for blockchain integration - Sepolia testnet)

### Installation & Setup

1. **Clone the repository**

   ```bash
   git clone https://github.com/Niko1444/iumicert.git
   cd iumicert/packages/issuer
   ```

2. **Build the CLI application**

   ```bash
   go build -o micert cmd/*.go
   ```

3. **Configure environment** (copy and edit .env file)
   ```bash
   cp .env.example .env
   # Edit .env with your Sepolia private key and settings
   ```

### Quick Demo (Automated Scripts)

**🧹 Reset System (Clean All Data)**

```bash
./reset.sh
```

This script:

- Clears all existing data directories
- Recreates clean directory structure
- Prepares system for fresh data generation

**🚀 Generate Complete Dataset**

```bash
./generate.sh
```

This script:

- Generates 5 student academic journeys across 6 terms
- Converts data to Verkle tree format
- Creates term-level Verkle trees with real cryptographic proofs
- Generates receipts for all students
- Outputs blockchain-ready term roots

**💻 Start Development Server**

```bash
./dev.sh
```

Starts the REST API server on port 8080 with CORS enabled for frontend development.

## 🛠️ CLI Commands Reference

The `micert` CLI provides comprehensive tools for managing academic credentials with Verkle trees and blockchain integration.

### Data Generation & Setup

**Initialize System**

```bash
./micert init "IU-VNUHCM"
```

- Creates initial configuration for the institution
- Sets up directory structure
- Configures default settings

**Generate Academic Data**

```bash
./micert generate-data
```

- Generates realistic student academic journeys (5 students, 6 terms)
- Creates course completion data with grades and timestamps
- Simulates natural academic progression across multiple semesters

**Convert Data to Verkle Format**

```bash
./micert convert-data <term-id>
./micert convert-data Semester_1_2023
```

- Converts LMS-style data into Verkle tree compatible format
- Prepares data for cryptographic tree construction
- Outputs structured JSON for term processing

### Term Processing & Verkle Trees

**Add Academic Term**

```bash
./micert add-term <term-id> <data-file>
./micert add-term Semester_1_2023 ./data/verkle_terms/Semester_1_2023_completions.json
```

- Creates single Verkle tree for entire academic term
- Includes all student course completions in one tree
- Generates cryptographic root commitment
- Uses real `ethereum/go-verkle` library for 32-byte proofs

**Batch Process All Terms**

```bash
./micert batch-process
```

- Automatically processes all available terms from generated data
- Converts data and creates Verkle trees for each term
- Equivalent to running convert-data and add-term for all terms

### Receipt Generation & Verification

**Generate Student Receipt**

```bash
./micert generate-receipt <student-id>
./micert generate-receipt ITITIU00001

# With specific term filter
./micert generate-receipt ITITIU00001 --terms "Semester_1_2023,Semester_2_2023"

# With course filter (selective disclosure)
./micert generate-receipt ITITIU00001 --courses "IT013IU,IT153IU"
```

- Generates cryptographic receipt with Verkle proofs
- Supports selective disclosure (specific courses/terms only)
- Creates 32-byte proofs for privacy-preserving verification
- Outputs blockchain-ready receipt files

**Display Receipt (Human-Readable)**

```bash
./micert display-receipt <receipt-file>
./micert display-receipt publish_ready/receipts/receipt_ITITIU00001_all_terms.json
```

- Shows receipt contents in readable format
- Displays course details, grades, and proof information
- Shows blockchain verification status
- Useful for debugging and demonstration

**Generate All Student Receipts**

```bash
./micert generate-all-receipts

# With output directory
./micert generate-all-receipts --output-dir ./custom_receipts

# With selective terms/courses
./micert generate-all-receipts --terms "Semester_1_2023,Semester_2_2023" --courses "IT013IU,IT153IU"

# Enable selective disclosure mode
./micert generate-all-receipts --selective
```

- Batch generates receipts for all students in system
- Creates comprehensive receipt set for testing
- Supports filtering by terms and courses
- Outputs individual files for each student

### Blockchain Integration

**Publish Term Roots**

```bash
# Publish all term roots to blockchain
./micert publish-roots

# Publish specific term
./micert publish-roots Test_Term_Fixed

# With custom gas settings and private key
./micert publish-roots --gas-limit 600000 --network sepolia --private-key YOUR_PRIVATE_KEY

# Using environment variable for private key (recommended)
export ISSUER_PRIVATE_KEY=your_private_key_here
./micert publish-roots Semester_1_2023
```

- Publishes Verkle tree root commitments to Ethereum
- Uses deployed smart contracts on Sepolia testnet
- Requires ISSUER_PRIVATE_KEY environment variable or --private-key flag
- Enables on-chain verification of academic credentials
- Supports custom gas limits and network selection

### Verification

**Local Verification**

```bash
./micert verify-local <receipt-file>
./micert verify-local publish_ready/receipts/receipt_ITITIU00001_all_terms.json
```

- Verifies receipt cryptographic proofs locally
- Checks Verkle proof validity without blockchain
- Fast verification for development and testing

**Verification Guide**

```bash
./micert verification-guide
```

- Shows step-by-step guide for third-party verifiers
- Explains how to verify receipts independently
- Includes both local and blockchain verification methods

### Development & Utilities

**Start API Server**

```bash
# Start API server with default settings
./micert serve

# Custom port and CORS settings
./micert serve --port 8080 --cors

# Disable CORS for production
./micert serve --port 3000 --cors=false
```

- Starts REST API server for web interface
- Default port 8080, CORS enabled by default
- Provides endpoints for receipt generation and verification
- Integrates with React frontend application

**Version Information**

```bash
./micert version
```

- Shows application version and build information
- Displays Go version and commit hash

**Help & Documentation**

```bash
./micert --help                    # General help
./micert <command> --help          # Command-specific help
./micert generate-receipt --help   # Example: receipt generation options
```

## 📋 Automation Scripts

### Reset Script (`reset.sh`)

```bash
./reset.sh
```

**Purpose**: Complete system reset and cleanup

- Removes all generated data (student journeys, Verkle terms, receipts, etc.)
- Clears blockchain-ready outputs (roots, proofs, transactions)
- Recreates clean directory structure
- Prepares system for fresh data generation

**Use Cases**:

- Starting over with clean slate
- Debugging data generation issues
- Preparing for new test scenarios

### Generation Script (`generate.sh`)

```bash
./generate.sh
```

**Purpose**: End-to-end data generation and processing

- **Step 1**: Generates 5 student academic journeys across 6 terms
- **Step 2**: Converts each term's data to Verkle format
- **Step 3**: Creates Verkle trees with real cryptographic proofs
- **Step 4**: Generates receipts for all students
- **Step 5**: Outputs blockchain-ready term roots

**What You Get**:

- Complete academic dataset (students: ITITIU00001-ITITIU00005)
- 6 term Verkle trees (Semester_1_2023 through Summer_2024)
- Individual student receipts with 32-byte proofs
- Term root commitments ready for blockchain publishing

### Development Script (`dev.sh`)

```bash
./dev.sh
```

**Purpose**: Start development environment

- Launches REST API server on port 8080
- Enables CORS for frontend development
- Provides backend services for web interface

5. **Start the API server**

   ```bash
   go run cmd/*.go serve --port 8080 --cors
   ```

6. **Launch web interface**
   ```bash
   cd web/iumicert-issuer
   npm install
   npm run dev
   ```

## 📊 Features

### Verkle Tree Architecture

- **Single Verkle Trees per Term**: Each academic term becomes one complete Verkle tree
- **32-byte Proofs**: Compact, efficient proofs using `ethereum/go-verkle` library
- **Course-Level Granularity**: Individual proofs for each course completion
- **Real Cryptographic Operations**: Uses Ethereum Foundation's official Verkle implementation
- **Selective Disclosure**: Students can prove specific achievements without revealing full transcript

### Academic Data Management

- **Realistic Test Data Generation**: Creates 5 Vietnamese students with 6 terms of academic progression
- **IU Vietnam Course Codes**: Uses actual course identifiers (IT013IU, IT153IU, PE008IU, etc.)
- **Natural Academic Flow**: Simulates realistic semester progression with prerequisites
- **LMS Data Conversion**: Transforms university records into Verkle-compatible format
- **Batch Processing**: Efficiently processes multiple terms simultaneously

### Privacy & Security

- **Zero-Knowledge Proofs**: Verify credentials without revealing unnecessary information
- **Tamper-Proof Records**: Cryptographically secured academic data with blockchain anchoring
- **Privacy-Preserving Verification**: Third parties can verify specific achievements only
- **Decentralized Trust**: No need to contact issuing institution for verification

### Blockchain Integration

- **Ethereum Sepolia Testnet**: Live deployment for testing and demonstration
- **Smart Contract Integration**: IUMiCertRegistry and ReceiptRevocationRegistry deployed
- **Root Publishing**: Automatic term commitment anchoring to blockchain
- **Gas Optimization**: Configurable gas limits and price settings
- **Transaction Monitoring**: Real-time blockchain confirmation tracking
- **On-Chain Verification**: Public verification without institutional dependency

## ⚙️ Environment Configuration

### Required Environment Variables

Create or edit `.env` file in the project root:

```bash
# Blockchain Configuration
ISSUER_PRIVATE_KEY=your_sepolia_private_key_here          # Your Ethereum wallet private key
IUMICERT_CONTRACT_ADDRESS=0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60  # Deployed contract
NETWORK=sepolia                                           # Blockchain network

# RPC Configuration
SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161  # Infura endpoint

# Gas Settings
DEFAULT_GAS_LIMIT=500000                                  # Transaction gas limit
MAX_GAS_PRICE=20000000000                                # Maximum gas price (20 gwei)

# Development Settings
DEBUG=true                                               # Enable debug logging
LOG_LEVEL=info                                          # Logging verbosity
```

### Deployed Contracts (Sepolia Testnet)

- **IUMiCertRegistry**: `0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60`

  - Stores term Verkle root commitments
  - Enables public verification of academic credentials
  - Provides term metadata and timestamps

- **ReceiptRevocationRegistry**: `0x8814ae511d54Dc10C088143d86110B9036B3aa92`
  - Manages credential revocation (if needed)
  - Provides revocation checking for verifiers

### Getting Sepolia ETH

For blockchain operations, you'll need Sepolia testnet ETH:

1. **Sepolia Faucet**: https://sepoliafaucet.com/
2. **Alchemy Faucet**: https://sepoliafaucet.com/
3. **Chainlink Faucet**: https://faucets.chain.link/sepolia

**💡 Tip**: You only need a small amount (~0.01 ETH) for publishing term roots.

## 🎯 Typical Workflow

### For New Setup (Complete Process)

1. **Reset & Setup**

   ```bash
   ./reset.sh                    # Clean slate
   ./micert init "IU-VNUHCM"    # Initialize system
   ```

2. **Generate & Process Data**

   ```bash
   ./generate.sh                 # Complete data generation pipeline
   # OR manually:
   # ./micert generate-data
   # ./micert batch-process
   ```

3. **Publish to Blockchain**

   ```bash
   export ISSUER_PRIVATE_KEY="your_private_key"
   ./micert publish-roots
   ```

4. **Generate Student Receipts**
   ```bash
   ./micert generate-receipt ITITIU00001
   ./micert display-receipt publish_ready/receipts/receipt_ITITIU00001_all_terms.json
   ```

### For Development & Testing

1. **Start Backend**

   ```bash
   ./dev.sh                      # API server on :8080
   ```

2. **Start Frontend**

   ```bash
   cd web/iumicert-issuer
   npm install && npm run dev    # Next.js on :3000
   ```

3. **Verify Locally**

   ```bash

   ```

## 🌐 API Endpoints

The REST API provides programmatic access with user-type separation:

### **Issuer Endpoints** (`/api/issuer/*`)
*For institutional dashboard and administrative operations*

**Terms Management**
- `GET /api/issuer/terms` - List all processed academic terms
- `POST /api/issuer/terms` - Process new academic term  
- `GET /api/issuer/terms/{term_id}/receipts` - Get all receipts for term
- `GET /api/issuer/terms/{term_id}/roots` - Get term Verkle root commitment

**Receipt Generation**
- `POST /api/issuer/receipts` - Generate new student receipt
- `GET /api/issuer/receipts` - List all generated receipts

**Student Management** 
- `GET /api/issuer/students` - List all students in system
- `GET /api/issuer/students/{student_id}/journey` - Complete academic journey
- `GET /api/issuer/students/{student_id}/terms` - Student's completed terms

**Blockchain Operations**
- `POST /api/issuer/blockchain/publish` - Publish term roots to blockchain
- `GET /api/issuer/blockchain/transactions` - List all blockchain transactions
- `GET /api/issuer/blockchain/transactions/{tx_hash}` - Get transaction details

### **Verifier Endpoints** (`/api/verifier/*`) 
*Public endpoints for students, employers, and third-party verifiers*

**Receipt Verification**
- `POST /api/verifier/receipt` - Verify receipt cryptographic proofs  
- `POST /api/verifier/course` - Verify specific course completion
- `GET /api/verifier/receipt/{receipt_id}` - Retrieve specific receipt by ID

**Student Data Access**
- `GET /api/verifier/journey/{student_id}` - Get student academic journey (if public)

**Blockchain Verification**
- `GET /api/verifier/blockchain/transaction/{tx_hash}` - Check transaction status

### **System Endpoints** (Public)

- `GET /api/health` - Health check endpoint  
- `GET /api/status` - System status and configuration

### **Legacy Endpoints** (Backward Compatibility)

- `GET /api/terms` - List terms (legacy)
- `GET /api/terms/{term_id}/roots` - Get term roots (legacy) 
- `POST /api/receipts/verify` - Verify receipt (legacy)
- `POST /api/receipts/verify-course` - Verify course (legacy)
- `POST /api/blockchain/publish` - Publish to blockchain (legacy)
- `GET /api/blockchain/transactions` - List transactions (legacy)

## 🎓 Academic Data Sample

The system includes realistic IU Vietnam academic progression:

- **5 Students**: ITITIU00001-ITITIU00005 (expandable to 100+)
- **6 Academic Terms**: Semester 1/2 2023-2024, Summer 2023-2024
- **Real Course Codes**: IT013IU, IT153IU, PE008IU, MA001IU, etc.
- **Natural Progression**: Simulates prerequisite chains and academic flow
- **Vietnamese Context**: Authentic names and institutional structure

## 🔧 System Configuration

Configuration is managed through `config/micert.json` and `.env` files:

```json
{
  "institution_id": "IU-VNUHCM",
  "version": "1.0.0",
  "blockchain": {
    "default_network": "sepolia",
    "gas_limit": 500000,
    "confirmation_blocks": 3
  },
  "output_paths": {
    "receipts": "publish_ready/receipts",
    "roots": "publish_ready/roots",
    "transactions": "publish_ready/transactions"
  }
}
```

## 🌐 Web Interface

The Next.js web interface provides:

- **Dashboard**: System overview and term statistics
- **Student Journey Viewer**: Interactive academic timeline with course details
- **Receipt Generator**: User-friendly credential creation tools
- **Verification Portal**: Local and blockchain-based proof verification
- **Blockchain Monitor**: Real-time transaction tracking and status

### Recent Improvements ✨

- **Fixed Frontend Error Handling**: Blockchain publishing now shows success messages instead of 500 errors
- **Duplicate Transaction Prevention**: System prevents republishing already-published terms
- **Graceful Error Recovery**: Backend API failures no longer interfere with successful blockchain transactions

Start with:

```bash
cd web/iumicert-issuer && npm install && npm run dev
```

Access at `http://localhost:3000`

### API Architecture Update 🏗️

The system now supports dual interfaces:

- **Issuer Dashboard** (Port 3000): For institution administrators - uses legacy and new `/api/issuer/*` endpoints
- **Verifier Interface** (Future): For students/employers - will use new `/api/verifier/*` endpoints

Both interfaces are served by a single backend API server with clear endpoint separation for security and scalability.

## 🔐 Security & Privacy

### Cryptographic Features

- **Real Verkle Proofs**: 32-byte proofs using `ethereum/go-verkle`
- **Zero-Knowledge Verification**: Prove achievements without revealing transcript
- **Selective Disclosure**: Students control what information to share
- **Tamper-Proof Records**: Cryptographic integrity protection

### Blockchain Security

- **Immutable Anchoring**: Term roots stored permanently on Ethereum
- **Decentralized Verification**: No need to contact issuing institution
- **Public Auditability**: Anyone can verify credential authenticity
- **Revocation Support**: Optional credential invalidation mechanism

## 🎯 Use Cases

### For Students

- **Privacy-Preserving Sharing**: Share only relevant course completions
- **Instant Verification**: Recipients can verify immediately without institution contact
- **Portable Credentials**: Works independently of university systems
- **Tamper-Evident**: Impossible to forge or modify achievements

### For Employers/Verifiers

- **Independent Verification**: Verify credentials without contacting university
- **Real-Time Validation**: Instant cryptographic proof checking
- **Fraud Prevention**: Cryptographically impossible to forge
- **Granular Information**: Verify specific courses/skills only

### For Institutions

- **Reduced Administrative Load**: Automated verification reduces support requests
- **Enhanced Trust**: Cryptographic proofs increase credential reliability
- **Modern Technology**: Demonstrates innovation in educational technology
- **Student Privacy**: Enables selective sharing without full transcript disclosure

## 🚀 Production Deployment

### Requirements

- **Ethereum Mainnet/L2**: Deploy contracts to production network
- **Key Management**: Secure private key storage (HSM recommended)
- **SSL/TLS**: Secure API endpoints
- **Database**: Persistent storage for large-scale operations

### Integration Steps

1. **LMS Connection**: Integrate with university information systems
2. **Data Pipeline**: Set up automated term processing workflows
3. **Monitoring**: Add comprehensive logging and alerting
4. **Scaling**: Configure load balancing for high-volume operations

## 🔧 Troubleshooting

### Common Issues & Recent Fixes

**❌ "Publishing failed: API Error: 500 Internal Server Error" (FIXED ✅)**
- **Issue**: Frontend showed error despite successful blockchain transactions
- **Solution**: Fixed frontend error handling to separate blockchain success from API failures
- **Status**: Resolved - now shows success messages correctly

**❌ Duplicate Term Publishing Errors (FIXED ✅)**  
- **Issue**: Smart contract rejected already-published terms causing 500 errors
- **Solution**: Added duplicate detection logic in both CLI and API
- **Status**: Resolved - gracefully handles republication attempts

**❌ Transaction History Not Clearing**
- **Issue**: `reset.sh` doesn't clear blockchain transaction records 
- **Solution**: Transaction files in `publish_ready/transactions/` persist intentionally
- **Workaround**: Manually delete files if needed: `rm publish_ready/transactions/tx_*.json`

### Development Tips

- **API Server**: Always rebuild binary after code changes: `go build -o micert ./cmd`
- **Frontend**: Uses legacy endpoints by default - new structured endpoints available  
- **Blockchain**: Set `ISSUER_PRIVATE_KEY` environment variable for publishing
- **Testing**: Use `curl` to test API endpoints directly: `curl localhost:8080/api/health`

### Port Conflicts  

If you see "port already in use" errors:
```bash
# Kill existing processes
pkill -f "micert serve"
lsof -ti:8080 | xargs kill -9

# Restart cleanly  
./dev.sh
```

## � Complete CLI Command Reference

This table shows all implemented commands in the current system:

| Command                 | Purpose                          | Example Usage                                 | Status         |
| ----------------------- | -------------------------------- | --------------------------------------------- | -------------- |
| `init`                  | Initialize credential repository | `./micert init "IU-VNUHCM"`                   | ✅ Implemented |
| `generate-data`         | Generate academic test data      | `./micert generate-data`                      | ✅ Implemented |
| `convert-data`          | Convert data to Verkle format    | `./micert convert-data Semester_1_2023`       | ✅ Implemented |
| `add-term`              | Add term with Verkle tree        | `./micert add-term Semester_1_2023 data.json` | ✅ Implemented |
| `batch-process`         | Process all terms automatically  | `./micert batch-process`                      | ✅ Implemented |
| `generate-receipt`      | Generate single student receipt  | `./micert generate-receipt ITITIU00001`       | ✅ Implemented |
| `generate-all-receipts` | Generate all student receipts    | `./micert generate-all-receipts`              | ✅ Implemented |
| `display-receipt`       | Show receipt details             | `./micert display-receipt receipt.json`       | ✅ Implemented |
| `verify-local`          | Verify receipt locally           | `./micert verify-local receipt.json`          | ✅ Implemented |
| `publish-roots`         | Publish roots to blockchain      | `./micert publish-roots`                      | ✅ Implemented |
| `serve`                 | Start API server                 | `./micert serve --port 8080`                  | ✅ Implemented |
| `verification-guide`    | Show verification steps          | `./micert verification-guide`                 | ✅ Implemented |
| `version`               | Show version info                | `./micert version`                            | ✅ Implemented |
| `completion`            | Generate shell completions       | `./micert completion bash`                    | ✅ Implemented |
| `help`                  | Command help                     | `./micert help <command>`                     | ✅ Implemented |

**Note**: All commands listed above are fully implemented and tested. This reflects the current state of the system using single Verkle tree architecture.

## �📚 Documentation

- **API Documentation**: OpenAPI/Swagger specs available
- **Verification Guide**: Step-by-step instructions for third parties
- **Smart Contract Documentation**: Contract ABI and interaction guides
- **Deployment Guide**: Production deployment best practices

## 🤝 Contributing

We welcome contributions! Please:

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

## � License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **Ethereum Foundation** - For the `ethereum/go-verkle` library
- **IU Vietnam** - For academic structure and course code references
- **Cobra CLI** - For command-line interface framework
- **Next.js** - For the web interface framework

---

**🎓 IU-MiCert Issuer System** - _Securing Academic Excellence with Blockchain Technology_

_Empowering students with verifiable, privacy-preserving academic credentials._
