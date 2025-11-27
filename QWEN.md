# IU-MiCert: Blockchain-Based Verifiable Academic Micro-Credential Provenance System

## Project Overview

IU-MiCert is a blockchain-based system designed to address critical limitations in current academic credential verification by implementing verifiable micro-credential provenance. The system leverages **Verkle tree technology** as an improvement over traditional Merkle trees to provide enhanced credential verification, efficient storage, and seamless integration with existing credential management systems.

### Core Purpose
- Enhance credential verification through verifiable academic micro-credential provenance using Verkle trees
- Provide enhanced anti-forgery mechanisms through temporal verification
- Enable efficient storage and verification of granular learning achievements
- Maintain comprehensive audit trails of learning achievements with verifiable timestamps

### Key Features
- **Verkle Tree Implementation**: Compact proofs with efficient verification
- **Term-by-Term Verification**: Verify specific academic periods independently
- **Academic Provenance**: Complete, tamper-proof timeline of achievements
- **Micro-Credential Tracking**: Each course becomes an independently verifiable credential
- **Selective Disclosure**: Share specific achievements without revealing full transcript
- **Temporal Integrity**: Prevents backdating and timeline manipulation

### Technical Stack
- **Blockchain**: Ethereum (with potential for other EVM chains)
- **Smart Contracts**: Solidity with Foundry framework
- **Cryptography**: Verkle tree implementation using Ethereum's `go-verkle` library
- **Backend**: Go (for issuer system)
- **Frontend**: Next.js with React (for demo interfaces)
- **Backend API**: Node.js (for institutional integration)
- **Styling**: Tailwind CSS and shadcn/ui components

## Repository Structure

```
iumicert/
├── packages/
│   ├── issuer/                  # Issuer System (Go + CLI + REST API)
│   │   ├── cmd/                 # CLI commands & API server
│   │   ├── crypto/verkle/       # Verkle tree implementation
│   │   ├── data/                # Academic data & test generation
│   │   ├── publish_ready/       # Receipts, roots, blockchain records
│   │   └── web/iumicert-issuer/ # Admin dashboard (Next.js)
│   ├── client/                  # Verification Portal (Next.js)
│   │   └── iumicert-client/     # Public verification interface
│   ├── contracts/               # Smart contracts (Solidity + Foundry)
│   ├── crypto/                  # Cryptographic implementations (Go)
│   └── data/                    # Academic data structures
├── docs/                        # Technical documentation
└── README.md                    # Main project documentation
```

## Building and Running

### Prerequisites
- Go 1.21+
- Node.js 18+
- Docker & Docker Compose
- Ethereum wallet for Sepolia testnet

### Issuer System Setup
1. Navigate to the issuer directory: `cd packages/issuer`
2. Build the CLI: `go build -o micert cmd/*.go`
3. Configure environment: `cp .env.example .env` and edit with your private key
4. Install web dashboard dependencies: `cd web/iumicert-issuer && npm install`

### Running the System
```bash
# Complete setup (reset + generate data + start server)
./reset.sh && ./generate.sh && ./dev.sh

# Start web dashboard in separate terminal
cd web/iumicert-issuer && npm run dev
```

### Client Portal Setup
1. Navigate to client: `cd packages/client/iumicert-client`
2. Install dependencies: `npm install`
3. Configure environment: Create `.env.local` with `NEXT_PUBLIC_API_URL=http://localhost:8080`
4. Start development server: `npm run dev`

### Smart Contracts
The system includes two main contracts:
- **IUMiCertRegistry**: Stores term roots (`0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60`)
- **ReceiptRevocationRegistry**: Handles revocation (`0x8814ae511d54Dc10C088143d86110B9036B3aa92`)

Build with Foundry: `forge build`
Test with Foundry: `forge test`

## Development Conventions

### Code Structure
- Go backend follows standard Go module structure with Cobra CLI framework
- Next.js frontend uses App Router with React Server Components
- Solidity contracts follow Foundry project structure
- Cryptographic operations separated in dedicated crypto package

### Security Practices
- Uses 32-byte Verkle proofs for each course completion
- Implements zero-knowledge verification without contacting institutions
- Enables selective disclosure with cryptographic proofs
- Blockchain anchoring for immutable verification

### Testing
- Foundry tests for smart contracts
- Go tests for cryptographic implementations
- Integration tests through CLI and API commands
- End-to-end testing with the web interfaces

## Key Components

### Verkle Tree Manager
- Efficient storage/verification of micro-credentials using Ethereum's `go-verkle` library
- Generates 32-byte IPA (Inner Product Argument) proofs per course
- Enables selective disclosure: prove specific courses without revealing all

### Smart Contracts
- Automated credential issuance with term-based cycles
- Term roots published to Ethereum for independent verification

### Commitment Engine
- Minimizes on-chain storage while maintaining provenance
- Reduces gas costs for credential publishing

### Verification Protocols
- Efficient proof validation at scale
- Both local and blockchain verification capabilities

## Deployment

### Live Demo & Deployed Contracts
**Ethereum Sepolia Testnet:**
- **IUMiCertRegistry**: `0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60`
- **ReceiptRevocationRegistry**: `0x8814ae511d54Dc10C088143d86110B9036B3aa92`

**Web Applications:**
- **Student/Verifier Portal**: [https://iu-micert.vercel.app](https://iu-micert.vercel.app)
- **Issuer Dashboard**: [https://iumicert-issuer.vercel.app](https://iumicert-issuer.vercel.app)

### Client Portal
Deployed on Vercel with environment variables:
- `NEXT_PUBLIC_API_URL` for backend API connection

## CLI Commands

The issuer system includes 15+ commands:
- `init` - Initialize system
- `generate-data` - Create test data
- `convert-data` - Convert to Verkle format
- `add-term` - Add term with Verkle tree
- `batch-process` - Process all terms
- `generate-receipt` - Create receipt
- `publish-roots` - Publish to blockchain
- `verify-local` - Verify receipt locally
- `serve` - Start API server
- And more...

## Security Features

- **32-byte Verkle Proofs**: Real cryptographic proofs using Ethereum's go-verkle library
- **Zero-Knowledge Verification**: Prove credentials without revealing full transcript
- **Selective Disclosure**: Students receive full receipts and can create filtered versions showing only specific courses/terms while maintaining cryptographic verifiability
- **Blockchain Anchoring**: Immutable term roots on Ethereum for independent verification
- **Tamper-Proof**: Cryptographically impossible to forge credentials

## Academic Data Structure

The system processes student data through multiple phases:
1. **Student Journeys**: Academic progression data in `packages/data/student_journeys/`
2. **Verkle Terms**: Converted format for Verkle tree processing in `packages/data/verkle_terms/`
3. **Verkle Trees**: Cryptographic structures in `packages/data/verkle_trees/`
4. **Publish Ready**: Final output with receipts and blockchain records

## Research Context

This system was developed as part of a Bachelor's thesis titled "IU-MiCert: Blockchain-Based Verifiable Academic Micro-Credential Provenance System" at International University - Vietnam National University HCM, representing a comprehensive approach to solving academic credential verification challenges using blockchain technology and advanced cryptographic structures.