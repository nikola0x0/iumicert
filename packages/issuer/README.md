# IU-MiCert Issuer System

A comprehensive academic credential management system for educational institutions using hybrid Merkle-Verkle tree architecture with blockchain integration.

## 🎓 Overview

The IU-MiCert Issuer System enables universities to:

- Process and validate academic course completions
- Generate cryptographic proofs for student achievements
- Issue verifiable academic journey receipts
- Publish term credentials to blockchain networks
- Visualize student learning progress

## 🏗️ Architecture

### Core Components

1. **Academic Data Processing Pipeline**

   - Converts LMS data into cryptographic structures
   - Builds student-level Merkle trees for each term
   - Aggregates into term-level Verkle trees

2. **Cryptographic Proof System**

   - Student-term Merkle trees for course completions
   - Term-level Verkle trees for efficient aggregation
   - Selective disclosure support for privacy

3. **Blockchain Integration**

   - Publishes term root commitments to Ethereum (Sepolia testnet)
   - Enables on-chain verification of credentials
   - Transaction monitoring and confirmation

4. **Web Interface**
   - Student journey visualization
   - Batch term processing
   - Receipt generation and verification

## 📁 Project Structure

```
packages/issuer/
├── cmd/                           # CLI application
│   ├── main.go                   # Main CLI entry point
│   ├── api_server.go            # REST API server
│   ├── data_generator.go        # Realistic test data generator
│   ├── data_converter.go        # LMS data converter
│   └── batch_processor.go       # Batch processing pipeline
├── data/                         # Data storage
│   ├── generated_student_data/   # Realistic academic data
│   │   ├── students/            # Individual student journeys
│   │   └── terms/               # Term summaries
│   ├── merkle_trees/            # Student-term Merkle trees
│   └── verkle_trees/            # Term-level Verkle trees
├── blockchain_ready/             # Blockchain integration
│   ├── receipts/                # Student academic receipts
│   ├── roots/                   # Term root commitments
│   └── transactions/            # Blockchain transactions
├── config/                       # System configuration
└── web/                         # Frontend application
    └── iumicert-issuer/         # Next.js web interface
```

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+
- Git

### Installation

1. **Clone and setup**

   ```bash
   git clone https://github.com/Niko1444/iumicert.git
   cd iumicert/packages/issuer
   ```

2. **Initialize the system**

   ```bash
   go run cmd/*.go init "IU-VNUHCM"
   ```

3. **Generate test data**

   ```bash
   go run cmd/*.go generate-data
   ```

4. **Process all terms into credential system**

   ```bash
   go run cmd/*.go batch-process
   ```

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

### Academic Data Management

- **Realistic Test Data Generation**: Creates 100 Vietnamese students with 6 terms of academic data using actual IU Vietnam course codes
- **LMS Integration**: Converts academic records into cryptographic format
- **Batch Processing**: Efficiently processes multiple terms simultaneously

### Cryptographic Security

- **Merkle Trees**: Individual student-term proof structures
- **Verkle Trees**: Efficient term-level aggregation
- **Selective Disclosure**: Privacy-preserving credential sharing
- **Tamper-proof Records**: Cryptographically secured academic data

### Blockchain Integration

- **Ethereum Compatibility**: Sepolia testnet integration
- **Root Publishing**: Term commitment anchoring
- **Transaction Monitoring**: Real-time blockchain status
- **Verification Support**: On-chain proof validation

## 🛠️ CLI Commands

### System Management

```bash
# Initialize repository
go run cmd/*.go init <institution-id>

# Check system status
go run cmd/*.go version

# Start API server
go run cmd/*.go serve [--port 8080] [--cors]
```

### Data Processing

```bash
# Generate test data
go run cmd/*.go generate-data

# Convert student data format
go run cmd/*.go convert-data <term-id>

# Process single term
go run cmd/*.go add-term <term-id> <data-file>

# Batch process all terms
go run cmd/*.go batch-process
```

### Receipt Generation

```bash
# Generate student receipt
go run cmd/*.go generate-receipt <student-id> <output-file>

# Verify receipt locally
go run cmd/*.go verify-local <receipt-file>
```

### Blockchain Operations

```bash
# Publish term roots
go run cmd/*.go publish-roots <term-id> [--network sepolia] [--gas-limit 500000]
```

## 🌐 API Endpoints

### System

- `GET /api/status` - System status and configuration
- `GET /api/health` - Health check

### Terms

- `GET /api/terms` - List all processed terms
- `POST /api/terms` - Process new academic term
- `GET /api/terms/{term_id}/receipts` - Get term receipts
- `GET /api/terms/{term_id}/roots` - Get term root commitment

### Students

- `GET /api/students` - List all students
- `GET /api/students/{student_id}/journey` - Get complete academic journey
- `GET /api/students/{student_id}/terms` - Get student terms

### Receipts

- `POST /api/receipts` - Generate student receipt
- `POST /api/receipts/verify` - Verify receipt
- `GET /api/receipts` - List all receipts

### Blockchain

- `POST /api/blockchain/publish` - Publish term roots
- `GET /api/blockchain/transactions` - List transactions
- `GET /api/blockchain/transactions/{tx_hash}` - Get transaction details

## 📈 Sample Data

The system includes realistic IU Vietnam academic data:

- **100 Students**: ITITIU00001-ITITIU00100 format
- **6 Academic Terms**: Semester 1/2 2023-2024, Summer 2023-2024
- **78 Real Courses**: Actual IU Vietnam course codes (IT064IU, MA001IU, etc.)
- **2,700 Completions**: 3-6 courses per student per term
- **Vietnamese Names**: Realistic student name generation

## 🔧 Configuration

System configuration is stored in `config/micert.json`:

```json
{
  "institution_id": "IU-Vietnam",
  "version": "1.0.0",
  "blockchain": {
    "default_network": "sepolia",
    "gas_limit": 500000,
    "confirmation_blocks": 3
  },
  "output_paths": {
    "receipts": "blockchain_ready/receipts",
    "roots": "blockchain_ready/roots",
    "transactions": "blockchain_ready/transactions"
  }
}
```

## 📱 Web Interface

The web interface provides:

- **Dashboard**: System overview and statistics
- **Student Journey Viewer**: Interactive academic timeline
- **Term Management**: Batch processing controls
- **Receipt Generator**: Credential creation tools
- **Blockchain Monitor**: Transaction tracking

Access at `http://localhost:3000` after starting the development server.

## 🔐 Security Features

- **Cryptographic Proofs**: Merkle/Verkle tree verification
- **Immutable Records**: Blockchain anchoring
- **Privacy Protection**: Selective disclosure
- **Tamper Detection**: Hash-based integrity
- **Audit Trail**: Complete transaction history

## 🚀 Production Deployment

For production use:

1. **Environment Setup**

   - Configure Ethereum mainnet/L2 network
   - Set up proper key management
   - Enable SSL/TLS for API server

2. **Data Integration**

   - Connect to actual LMS/SIS systems
   - Implement data validation pipelines
   - Set up automated term processing

3. **Monitoring**
   - Add logging and metrics
   - Configure alerting systems
   - Monitor blockchain transaction costs

## 📄 License

MIT License - see LICENSE file for details.

## 🤝 Contributing

1. Fork the repository
2. Create feature branch
3. Make changes with tests
4. Submit pull request

## 📞 Support

For technical support or questions:

- Create GitHub issue
- Contact: [Your contact information]

---

**IU-MiCert Issuer System** - Securing Academic Excellence with Blockchain Technology
