# IU-MiCert Data Flow: Generation to Publication

## 🏗️ System Architecture Overview

The IU-MiCert system uses a **single Verkle tree architecture** for academic credential management, providing privacy-preserving verification with constant-size proofs. This document outlines the complete data flow from initial generation to blockchain publication.

## 📊 Core Components

- **Backend**: Go-based CLI with Cobra framework
- **Cryptography**: Real Verkle trees using `ethereum/go-verkle` library
- **Blockchain**: Ethereum Sepolia testnet integration
- **Web Interface**: Next.js with wagmi/viem for MetaMask integration
- **API**: REST API server for web interface communication

## 🔄 Complete Data Flow Pipeline

### Phase 1: Raw Data Generation

**Command**: `./micert generate-data`

**Input**: System configuration  
**Output**: Academic journey data

```
📊 Creates:
├── data/student_journeys/students/journey_ITITIU00001.json
├── data/student_journeys/students/journey_ITITIU00002.json
├── data/student_journeys/students/journey_ITITIU00003.json
├── data/student_journeys/students/journey_ITITIU00004.json
├── data/student_journeys/students/journey_ITITIU00005.json
└── data/student_journeys/terms/summary_*.json
```

**Data Characteristics**:

- 5 students with realistic IU Vietnam course codes
- 6 academic terms (2 semesters + summer for 2 years)
- Complete academic progressions with grades and timestamps
- Realistic course names, instructors, and credit hours

### Phase 2: Term Processing (Per Academic Term)

#### Step 2a: Data Conversion

**Command**: `./micert convert-data <term-id>`

**Input**: `data/student_journeys/students/*.json`  
**Output**: `data/verkle_terms/<term-id>_completions.json`

```
🔄 Transformation:
Student Journey Format → Verkle Tree Format

Example: Semester_1_2023_completions.json
[
  {
    "student_id": "ITITIU00001",
    "course_id": "IT153IU",
    "course_name": "Discrete Mathematics",
    "grade": "A+",
    "credits": 3,
    "term_id": "Semester_1_2023"
  }
  // ... more completions
]
```

#### Step 2b: Verkle Tree Creation

**Command**: `./micert add-term <term-id> <completions-file>`

**Input**: `data/verkle_terms/<term-id>_completions.json`  
**Output**:

- `data/verkle_trees/<term-id>/verkle_tree.json`
- `publish_ready/roots/root_<term-id>.json`

```
🌳 Verkle Tree Structure:
├── verkle_tree.json (complete tree with all student data)
└── root_<term-id>.json (blockchain-ready root commitment)

Root Format:
{
  "ready_for_blockchain": true,
  "term_id": "Semester_1_2023",
  "timestamp": "2025-08-27T07:39:25+07:00",
  "total_students": 5,
  "verkle_root": "72f90eadf541dc4fc9a885b6581c645ec55a8eb64e06dc807b9b73eceaf681a2"
}
```

### Phase 3: Receipt Generation

**Command**: `./micert generate-receipt <student-id>`

**Input**:

- `data/verkle_trees/*/verkle_tree.json`
- `data/student_journeys/students/journey_<student-id>.json`

**Output**: `publish_ready/receipts/<student-id>_journey.json`

```
🎓 Receipt Structure:
{
  "student_id": "ITITIU00001",
  "student_name": "Student ITITIU00001",
  "courses": [
    {
      "course_id": "IT153IU",
      "course_name": "Discrete Mathematics",
      "grade": "A+",
      "credits": 3,
      "term_id": "Semester_1_2023",
      "verkle_proof": "32-byte-proof-data"
    }
  ],
  "created_at": "2025-08-27T07:39:28+07:00",
  "blockchain_anchor": "verkle-root-reference"
}
```

### Phase 4: Blockchain Publication

**Interface**: Web UI with MetaMask integration  
**Command**: Via web interface → `publishTermRoot()`

**Input**: `publish_ready/roots/root_<term-id>.json`  
**Output**:

- Blockchain transaction on Sepolia
- `publish_ready/transactions/<tx-hash>.json`

```
⛓️ Blockchain Publication:
1. User selects term in web interface
2. MetaMask popup for transaction signature
3. Smart contract call: publishTermRoot(verkleRoot, termId, totalStudents)
4. Transaction confirmation and storage
```

## 📁 Directory Structure

### Working Data (`data/`)

```
data/
├── student_journeys/          # 📊 Raw Academic Data
│   ├── students/             # Individual student journey files
│   │   ├── journey_ITITIU00001.json
│   │   ├── journey_ITITIU00002.json
│   │   └── ...
│   └── terms/                # Term summary files
│       ├── summary_Semester_1_2023.json
│       └── ...
├── verkle_terms/             # 🔄 Converted Verkle Format
│   ├── Semester_1_2023_completions.json
│   ├── Semester_2_2023_completions.json
│   └── ...
└── verkle_trees/             # 🌳 Built Verkle Structures
    ├── Semester_1_2023/
    │   └── verkle_tree.json
    └── ...
```

### Publication Ready (`publish_ready/`)

```
publish_ready/
├── receipts/                 # 🎓 Student Verification Receipts
│   ├── ITITIU00001_journey.json
│   ├── ITITIU00002_journey.json
│   └── ...
├── roots/                    # 🔗 Blockchain-Ready Roots
│   ├── root_Semester_1_2023.json
│   ├── root_Semester_2_2023.json
│   └── ...
├── proofs/                   # 🔐 Individual Proof Files (future)
└── transactions/             # ⛓️ Blockchain Transaction Records
    └── tx_<hash>.json
```

## ⚡ Automated Pipeline

### Full Generation Script (`./generate.sh`)

```bash
#!/bin/bash
# Complete end-to-end pipeline

# 1. Generate student academic journeys
./micert generate-data

# 2. Process each term (6 terms total)
for term in "Semester_1_2023" "Semester_2_2023" "Summer_2023" \
           "Semester_1_2024" "Semester_2_2024" "Summer_2024"; do
    ./micert convert-data $term
    ./micert add-term $term data/verkle_terms/${term}_completions.json
done

# 3. Generate receipts for all students
for student in "ITITIU00001" "ITITIU00002" "ITITIU00003" \
               "ITITIU00004" "ITITIU00005"; do
    ./micert generate-receipt $student
done
```

### Quick Commands

```bash
# Reset and regenerate everything
./reset.sh && ./generate.sh

# Start development environment
./dev.sh  # Starts API server on port 8080

# Manual term processing
./micert batch-process  # Process all terms automatically

# Manual receipt generation
./micert generate-all-receipts  # Generate all student receipts
```

## 🔍 Key Features

### Single Verkle Tree Architecture

- **One Verkle tree per academic term** (not hybrid Merkle-Verkle)
- **32-byte constant proof size** regardless of tree size
- **Course-level selective disclosure** for privacy
- **Real cryptographic implementation** using `ethereum/go-verkle v0.2.2`

### Privacy-Preserving Verification

- Students can prove specific course completions without revealing other courses
- Proofs don't expose sensitive academic information
- Zero-knowledge verification of academic achievements

### Blockchain Integration

- **Sepolia testnet** deployment for testing
- **MetaMask integration** for secure transaction signing
- **Smart contract verification** of academic credentials
- **Immutable audit trail** of published term roots

## 🌐 Web Interface Integration

### API Endpoints

```
GET  /api/terms                    # List all academic terms
GET  /api/terms/{id}/roots         # Get Verkle root for term
GET  /api/terms/{id}/receipts      # Get student receipts for term
POST /api/receipts                 # Generate new receipt
POST /api/blockchain/publish       # Publish to blockchain
```

### Frontend Flow

1. **Connect MetaMask wallet**
2. **Browse academic terms** and student data
3. **Generate receipts** with selective disclosure
4. **Publish term roots** to blockchain via MetaMask
5. **Verify receipts** locally and on-chain

## 🎯 Current System Status

### Generated Data

- **5 students** with complete academic journeys
- **6 academic terms** with realistic IU course progressions
- **6 Verkle trees** with cryptographic proofs
- **6 blockchain-ready roots** for term publication
- **Multiple student receipts** with 32-byte proofs

### Deployed Components

- **Smart Contracts**: Deployed on Sepolia testnet
- **API Server**: REST endpoints for web integration
- **Web Interface**: Next.js application with MetaMask
- **CLI Tools**: Complete command-line management suite

## 🚀 Next Steps

### For Thesis Demonstration

1. **Performance Benchmarking**: Compare Verkle vs Merkle proof sizes
2. **Selective Disclosure Demo**: Show course-level privacy features
3. **Blockchain Verification**: Demonstrate end-to-end verification flow
4. **Scalability Testing**: Test with larger student/course datasets

### For Production Deployment

1. **Mainnet Migration**: Deploy contracts to Ethereum mainnet
2. **LMS Integration**: Connect with university information systems
3. **Enhanced Security**: Implement HSM for key management
4. **Monitoring & Alerts**: Add comprehensive system monitoring

## 📚 Technical References

- **Verkle Trees**: [ethereum/go-verkle](https://github.com/ethereum/go-verkle)
- **Smart Contracts**: Solidity contracts with OpenZeppelin
- **Web3 Integration**: wagmi + viem for Ethereum interaction
- **Academic Standards**: IU Vietnam course code structure

---

_This system demonstrates practical implementation of Verkle trees for academic credential management, providing privacy-preserving verification with blockchain integration._
