# IU-MiCert Project Structure

**A Transparent and Granular Blockchain System for Verifiable Academic Micro-credential Provenance**

## 📁 Project Organization

```
iumicert/
├── contracts/                          # Smart Contract Layer
│   ├── src/
│   │   ├── IUMiCertRegistry.sol        # Term-level Verkle root storage
│   │   ├── IUMiCertVerifier.sol        # On-chain verification utilities
│   │   └── libraries/
│   │       └── VerkleProofLib.sol      # Verkle proof verification library
│   ├── deploy/                         # Deployment scripts
│   │   ├── 01-deploy-registry.ts
│   │   └── 02-deploy-verifier.ts
│   ├── test/                           # Contract tests
│   ├── hardhat.config.ts               # Hardhat configuration
│   └── package.json
│
├── packages/
│   ├── blockchain/                     # Blockchain Integration Layer
│   │   ├── deploy/                     # Deployment utilities
│   │   │   ├── registry-deployer.ts
│   │   │   └── network-config.ts
│   │   ├── scripts/                    # Blockchain interaction scripts
│   │   │   ├── publish-term-roots.ts
│   │   │   └── batch-verify.ts
│   │   ├── abis/                       # Contract ABIs
│   │   │   ├── IUMiCertRegistry.json
│   │   │   └── IUMiCertVerifier.json
│   │   └── package.json
│   │
│   ├── crypto/                         # Cryptographic Operations Layer
│   │   ├── verkle/                     # Verkle tree implementation
│   │   │   ├── tree.go                 # Core Verkle tree operations
│   │   │   ├── proofs.go               # Proof generation/verification
│   │   │   └── commitment.go           # Cryptographic commitments
│   │   ├── merkle/                     # Student-term Merkle trees
│   │   │   ├── student_term.go         # Student-level Merkle implementation
│   │   │   └── course_leaf.go          # Course completion leaf structure
│   │   ├── go.mod
│   │   └── go.sum
│   │
│   ├── issuer/                         # CLI Tools for Credential Issuance
│   │   ├── cmd/                        # CLI commands
│   │   │   ├── init.go                 # micert init
│   │   │   ├── add_term.go             # micert add-term
│   │   │   ├── generate_receipt.go     # micert generate-receipt
│   │   │   ├── verify_local.go         # micert verify-local
│   │   │   └── publish_roots.go        # micert publish-roots
│   │   ├── internal/                   # Internal packages
│   │   │   ├── term/                   # Term management
│   │   │   ├── receipt/                # Receipt generation
│   │   │   └── blockchain/             # Blockchain interactions
│   │   ├── scripts/                    # Utility scripts
│   │   │   ├── generate_test_data.js   # Test data generation
│   │   │   └── backup_proofs.js        # Proof backup utilities
│   │   ├── data/                       # Test and example data
│   │   │   ├── academic_records_export.json
│   │   │   ├── multi_semester_academic_export.json
│   │   │   └── enhanced_credential_proofs/
│   │   │       ├── journey_STU001.json
│   │   │       ├── term_STU001_Fall_2021.json
│   │   │       └── ...
│   │   ├── backup_scripts/             # Backup and maintenance
│   │   ├── hybrid_credential_system.go # Hybrid Merkle-Verkle implementation
│   │   ├── enhanced_credential_system.js # Enhanced system (legacy)
│   │   ├── go.mod
│   │   └── go.sum
│   │
│   ├── verifier/                       # Verification Services
│   │   ├── api/                        # REST API for verification
│   │   │   ├── handlers/
│   │   │   │   ├── verify_term.go
│   │   │   │   ├── verify_journey.go
│   │   │   │   └── batch_verify.go
│   │   │   └── server.go
│   │   ├── lib/                        # Verification libraries
│   │   │   ├── proof_validator.go
│   │   │   ├── temporal_checker.go
│   │   │   └── fraud_detector.go
│   │   ├── go.mod
│   │   └── go.sum
│   │
│   └── client/                         # Frontend Application
│       └── iumicert-client/            # Next.js Web Application
│           ├── src/
│           │   ├── app/                # App router (Next.js 13+)
│           │   │   ├── components/     # React components
│           │   │   │   ├── AnimatedBackground.tsx
│           │   │   │   ├── FileUploaderWrapper.tsx
│           │   │   │   ├── Header.tsx
│           │   │   │   ├── Footer.tsx
│           │   │   │   ├── VerificationInterface.tsx
│           │   │   │   ├── JourneyVisualization.tsx
│           │   │   │   └── SelectiveDisclosure.tsx
│           │   │   ├── verify/         # Verification pages
│           │   │   │   ├── page.tsx
│           │   │   │   ├── term/
│           │   │   │   │   └── page.tsx
│           │   │   │   └── journey/
│           │   │   │       └── page.tsx
│           │   │   ├── revoke/         # Revocation interface
│           │   │   │   └── page.tsx
│           │   │   ├── layout.tsx
│           │   │   ├── page.tsx        # Landing page
│           │   │   └── globals.css
│           │   ├── lib/                # Utility libraries
│           │   │   ├── blockchain.ts   # Blockchain interactions
│           │   │   ├── verification.ts # Verification logic
│           │   │   └── receipt-parser.ts # Receipt parsing
│           │   └── types/              # TypeScript definitions
│           │       ├── receipt.ts
│           │       └── blockchain.ts
│           ├── public/                 # Static assets
│           │   ├── logo.svg
│           │   ├── horizontal-certificate.svg
│           │   └── ...
│           ├── next.config.ts
│           ├── tailwind.config.ts
│           ├── package.json
│           └── README.md
│
├── docs/                               # Documentation
│   ├── api/                           # API documentation
│   ├── deployment/                    # Deployment guides
│   ├── user-guides/                   # User documentation
│   └── architecture.md               # System architecture
│
├── tests/                             # Integration tests
│   ├── e2e/                          # End-to-end tests
│   ├── performance/                  # Performance benchmarks
│   └── security/                     # Security tests
│
├── scripts/                          # Build and deployment scripts
│   ├── build-all.sh
│   ├── deploy-contracts.sh
│   └── generate-test-data.sh
│
├── .github/                          # GitHub workflows
│   └── workflows/
│       ├── ci.yml
│       └── deploy.yml
│
├── LICENSE                           # MIT License
├── README.md                         # Project overview
└── PROJECT_STRUCTURE.md             # This file
```

## 🎯 Component Responsibilities

### 🔗 Smart Contracts (`contracts/`)
- **IUMiCertRegistry.sol**: Stores term-level Verkle tree roots with temporal verification
- **IUMiCertVerifier.sol**: On-chain verification utilities and batch operations
- **VerkleProofLib.sol**: Cryptographic proof verification library

### ⚡ Crypto Layer (`packages/crypto/`)
- **Verkle Trees**: Per-term aggregation with constant-size proofs
- **Merkle Trees**: Student-term course completion tracking
- **Hybrid System**: Combines both approaches for optimal efficiency

### 🛠️ CLI Tools (`packages/issuer/`)
- **Term Management**: Add academic terms and course completions
- **Receipt Generation**: Create journey receipts with selective disclosure
- **Blockchain Integration**: Publish term roots and verify proofs
- **Data Processing**: Handle LMS/SIS integration

### 🌐 Frontend (`packages/client/`)
- **Multi-type Verification**: Single term and academic journey receipts
- **Selective Disclosure**: Granular control over revealed credentials
- **Timeline Visualization**: Academic progression tracking
- **Responsive Design**: Cross-platform compatibility

### 🔍 Verification Services (`packages/verifier/`)
- **API Layer**: REST endpoints for verification operations
- **Proof Validation**: Cryptographic proof verification
- **Fraud Detection**: Timeline manipulation and forgery detection
- **Temporal Checking**: Academic progression validation

### ⛓️ Blockchain Integration (`packages/blockchain/`)
- **Deployment Scripts**: Smart contract deployment automation
- **Network Configuration**: Multi-network support (Sepolia, Mainnet)
- **ABI Management**: Contract interface definitions
- **Interaction Scripts**: Blockchain operation utilities

## 🔄 Data Flow

1. **Credential Issuance**:
   ```
   LMS/SIS Data → CLI Tools → Student-Term Merkle Trees → 
   Term Verkle Trees → Smart Contract Storage
   ```

2. **Receipt Generation**:
   ```
   CLI Tools → Academic Journey Receipts → 
   Students (Selective Disclosure) → Verification Interface
   ```

3. **Verification Process**:
   ```
   Frontend Upload → Proof Parsing → Verkle Verification → 
   Blockchain Validation → Result Display
   ```

## 🏗️ Architecture Principles

- **Hybrid Cryptographic Design**: Merkle trees for student-level, Verkle trees for aggregation
- **Per-Term Deployment**: Independent Verkle trees for each academic term
- **Selective Disclosure**: Privacy-preserving credential revelation
- **Timeline Integrity**: Anti-forgery through temporal verification
- **Modular Components**: Independent, reusable system modules
- **Scalable Design**: Efficient handling of large credential datasets

## 📊 Key Features

- ✅ **Granular Micro-credential Tracking**
- ✅ **Verifiable Academic Provenance** 
- ✅ **Anti-forgery Timeline Protection**
- ✅ **Constant-size Verkle Proofs**
- ✅ **Selective Disclosure Privacy**
- ✅ **Multi-term Journey Verification**
- ✅ **Blockchain Transparency**
- ✅ **CLI and Web Interfaces**