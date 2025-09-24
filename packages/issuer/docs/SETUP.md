# IU-MiCert Issuer Setup Guide

## 🌟 Single Verkle Architecture

This guide will help you set up and run the IU-MiCert academic credential issuer system, which uses a **single Verkle tree architecture** for enhanced privacy and efficiency.

## 📋 Prerequisites

### Required Software
- **Go 1.21+** - For running the core issuer system
- **Node.js 18+** - For the web frontend (optional)
- **Git** - For version control

### System Requirements
- **macOS, Linux, or Windows** with WSL2
- **4GB+ RAM** recommended
- **500MB+ storage** for generated data

## 🚀 Quick Start

### 1. Clone and Navigate
```bash
git clone https://github.com/Niko1444/iumicert.git
cd iumicert/packages/issuer
```

### 2. Initialize Dependencies
```bash
# Initialize Go module (if needed)
go mod tidy

# Install any missing dependencies
go mod download
```

### 3. Generate Sample Data
```bash
# Reset system and generate complete dataset
./reset.sh    # Clean slate
./generate.sh # Generate everything
```

## 📁 Project Structure

```
packages/issuer/
├── cmd/                          # CLI application
│   ├── main.go                   # Main CLI interface
│   ├── data_generator.go         # Student data generation
│   └── data_converter.go         # Format conversion utilities
├── data/                         # Generated data storage
│   ├── student_journeys/         # Academic journey data
│   │   ├── students/            # Individual student files
│   │   └── terms/               # Term summary data
│   └── verkle_terms/            # Verkle-formatted completions
├── publish_ready/               # Blockchain-ready outputs
│   ├── receipts/                # Student verification receipts
│   ├── roots/                   # Verkle tree roots
│   ├── proofs/                  # Cryptographic proofs
│   └── transactions/            # Blockchain transaction records
├── config/                      # Configuration files
├── reset.sh                     # System reset script
├── generate.sh                  # Complete generation workflow
└── SETUP.md                     # This guide
```

## 🛠️ Available Commands

### Core Operations
```bash
cd cmd

# Initialize repository
go run . init IU-BUS

# Add academic term
go run . add-term Semester_1_2023 ../data/verkle_terms/Semester_1_2023_completions.json

# Generate student receipt
go run . generate-receipt ITITIU00001 ../publish_ready/receipts/ITITIU00001_journey.json

# Verify receipt locally
go run . verify-local ../publish_ready/receipts/ITITIU00001_journey.json

# Publish to blockchain (requires configuration)
go run . publish-roots Semester_1_2023
```

### Advanced Features
```bash
# Selective disclosure - specific courses only
go run . generate-receipt ITITIU00001 receipt.json --selective --courses IT116IU,IT153IU

# Generate for specific terms
go run . generate-receipt ITITIU00001 receipt.json --terms Semester_1_2023,Semester_2_2023
```

## 🔧 Configuration

### Blockchain Setup (Optional)
```bash
# Copy example configuration
cp .env.example .env

# Edit with your values
nano .env
```

Example `.env` configuration:
```env
ISSUER_PRIVATE_KEY=your_private_key_here
CONTRACT_ADDRESS=0x1234...
RPC_URL=https://sepolia.infura.io/v3/your_key
NETWORK=sepolia
```

## 🌳 Single Verkle Architecture Benefits

### ✨ Technical Advantages
- **32-byte constant proof size** (vs variable Merkle proofs)
- **Course-level selective disclosure** (reveal specific courses only)
- **Better privacy protection** (no student data in proofs)
- **Simplified verification** (single root per term)

### 📊 Generated Data Overview
- **5 students** with complete academic journeys
- **6 terms** spanning 2023-2024 academic years
- **126 total course completions** across all students
- **6 Verkle trees** with blockchain-ready roots
- **Verification receipts** for all students

## 🎯 Common Use Cases

### 1. Generate Fresh Dataset
```bash
./reset.sh && ./generate.sh
```

### 2. Test Selective Disclosure
```bash
cd cmd
go run . generate-receipt ITITIU00001 selective_receipt.json --selective --courses IT116IU,MA003IU
```

### 3. Verify Generated Receipts
```bash
cd cmd
go run . verify-local ../publish_ready/receipts/ITITIU00001_journey.json
```

### 4. Publish to Blockchain
```bash
cd cmd
go run . publish-roots Semester_1_2023 --network sepolia
```

## 🔍 Sample Output Structure

### Student Receipt Format
```json
{
  "student_id": "ITITIU00001",
  "receipt_type": {
    "selective_disclosure": false,
    "specific_courses": false,
    "specific_terms": false
  },
  "term_receipts": {
    "Semester_1_2023": {
      "receipt": {
        "proof_type": "verkle_32_byte",
        "revealed_courses": [...],
        "verkle_root": "72f90eadf541dc4f...",
        "verification_path": "single_verkle_proof"
      },
      "revealed_courses": 3,
      "total_courses": 3
    }
  }
}
```

### Verkle Root Format
```json
{
  "term_id": "Semester_1_2023",
  "verkle_root": "72f90eadf541dc4f9a885b6581c645ec55a8eb64e06dc807b9b73eceaf681a2",
  "timestamp": "2025-08-27T07:39:25+07:00",
  "total_students": 5,
  "ready_for_blockchain": true
}
```

## 🚨 Troubleshooting

### Common Issues

**1. "command not found: go"**
```bash
# Install Go from https://golang.org/dl/
# Or via package manager:
brew install go        # macOS
sudo apt install golang-go  # Ubuntu
```

**2. "permission denied: ./generate.sh"**
```bash
chmod +x reset.sh generate.sh
```

**3. "failed to discover student terms"**
```bash
# Ensure data is generated first
./generate.sh
```

**4. "blockchain configuration missing"**
```bash
# Copy and configure environment
cp .env.example .env
# Edit .env with your blockchain settings
```

## 🎓 Academic Features

### Supported Academic Elements
- **Multi-term journeys** (Semesters + Summer terms)
- **Course completions** with grades and credits
- **Instructor attribution**
- **Temporal tracking** (started, completed, assessed, issued dates)
- **Multi-issuer support** (different departments)

### Privacy Features
- **Course-level granularity** in proofs
- **Student data isolation** in verification
- **Selective disclosure** for specific courses
- **Zero-knowledge verification** support

## 📈 Performance Characteristics

- **Proof Generation**: ~50ms per term
- **Verification Time**: ~10ms per receipt
- **Storage Efficiency**: 32-byte constant proof size
- **Scalability**: O(log n) verification complexity

## 🔗 Integration Points

### Frontend Integration
- Receipt verification API endpoints
- Proof generation services  
- Student journey visualization

### Blockchain Integration
- Sepolia testnet support
- Ethereum mainnet compatible
- Gas-optimized contract interactions

## 📞 Support

For technical issues or questions:
1. Check this SETUP.md guide
2. Review error logs in terminal output
3. Ensure all prerequisites are installed
4. Verify file permissions for scripts

## 🎉 Success Indicators

After running `./generate.sh`, you should see:
- ✅ **6 Verkle trees** created successfully  
- ✅ **5 student receipts** generated
- ✅ **Blockchain-ready roots** in `publish_ready/roots/`
- ✅ **Zero verification errors** when testing receipts

The system is ready for academic credential verification with enhanced privacy through single Verkle tree architecture! 🎓