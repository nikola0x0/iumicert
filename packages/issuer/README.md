# IU-MiCert CLI - Academic Micro-credential Issuer

**A Transparent and Granular Blockchain System for Verifiable Academic Micro-credential Provenance**

The IU-MiCert CLI provides comprehensive tools for managing academic micro-credentials using a hybrid Merkle-Verkle tree architecture with blockchain integration.

## 🏗️ Architecture Overview

- **Student-Term Level**: Merkle trees for individual course completions with temporal verification
- **Term-Aggregation Level**: Verkle trees for efficient multi-student aggregation
- **Blockchain Integration**: Smart contracts store term roots with anti-forgery protection
- **Selective Disclosure**: Privacy-preserving credential revelation

## 🚀 Quick Start

### Prerequisites

- Go 1.21+ installed
- Git for repository access

### Installation

1. **Clone and navigate to the issuer package:**
   ```bash
   cd packages/issuer
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Build the CLI:**
   ```bash
   go build -o micert ./cmd/
   ```

4. **Verify installation:**
   ```bash
   ./micert --help
   ```

## 🧪 Testing the System

### Run Complete System Test

The fastest way to see the system in action:

```bash
./micert test
```

This comprehensive test will:
- ✅ Generate realistic academic test data (3 students, 4 courses each)
- ✅ Build student-term Merkle trees with course completions
- ✅ Aggregate into term-level Verkle trees
- ✅ Generate and verify full academic receipts
- ✅ Test selective disclosure (partial course revelation)
- ✅ Validate anti-forgery mechanisms
- ✅ Test error handling for invalid inputs

**Expected Output:**
```
🧪 Starting IU-MiCert Hybrid System Test
📊 Step 1: Generating test data...
✅ Generated 12 course completions for term Fall_2024
🌲 Step 3: Testing Student-Term Merkle Trees...
  ✓ Student did:example:STU001: 4 courses, root: 9c456b11...
  ✓ All Merkle proofs verified for did:example:STU001
🌳 Step 4: Testing Term-Level Verkle Tree...
✅ Term published with Verkle root: 6fd96498...
📄 Step 5: Testing Receipt Generation and Verification...
  ✅ Receipt verified successfully for did:example:STU001
🎉 All tests passed successfully!
```

## 📋 CLI Commands Reference

### Core Commands

#### `micert init [institution-id]`
Initialize a new credential repository for an institution.

```bash
./micert init "IU-Vietnam"
```

#### `micert add-term [term-id] [data-file]`
Add academic term data with course completions.

```bash
./micert add-term "Fall_2024" "./data/fall_2024_completions.json"
```

**Supported flags:**
- `--format`: Input data format (`json`, `csv`) [default: `json`]
- `--validate`: Validate input data [default: `true`]

#### `micert generate-receipt [student-id] [output-file]`
Generate academic journey receipt for verification.

```bash
# Full receipt (all courses)
./micert generate-receipt "did:example:student001" "./receipts/student001_full.json"

# Selective disclosure
./micert generate-receipt "did:example:student001" "./receipts/student001_selective.json" \
  --courses "CS101,MATH101" \
  --selective
```

**Supported flags:**
- `--terms`: Specific terms to include (comma-separated)
- `--courses`: Specific courses to include (comma-separated)
- `--selective`: Enable selective disclosure mode

#### `micert verify-local [receipt-file]`
Perform off-chain verification without blockchain interaction.

```bash
./micert verify-local "./receipts/student001_full.json"
```

#### `micert publish-roots [term-id]`
Publish term Verkle roots to blockchain (requires setup).

```bash
./micert publish-roots "Fall_2024" \
  --network sepolia \
  --private-key $PRIVATE_KEY \
  --gas-limit 500000
```

### Utility Commands

#### `micert version`
Display version and system information.

#### `micert test`
Run comprehensive system testing suite.

## 📊 Data Formats

### Course Completion JSON Format

```json
{
  "issuer_id": "IU-CS",
  "student_id": "STU001", 
  "term_id": "Fall_2024",
  "course_id": "CS101",
  "course_name": "Introduction to Programming",
  "attempt_no": 1,
  "started_at": "2024-09-01T00:00:00Z",
  "completed_at": "2024-12-15T00:00:00Z", 
  "assessed_at": "2024-12-20T00:00:00Z",
  "issued_at": "2024-12-25T00:00:00Z",
  "grade": "A",
  "credits": 3,
  "instructor": "Prof. Johnson"
}
```

### Verification Receipt Format

```json
{
  "term_id": "Fall_2024",
  "student_did": "did:example:student001",
  "student_term_root": "9c456b11546e7cc696d5dae24db...",
  "verkle_proof": "base64-encoded-proof-data",
  "verkle_root": "6fd9649863592b48eea506bd93...", 
  "published_at": "2024-12-31T23:59:59Z",
  "revealed_courses": [...],
  "merkle_proofs": {...},
  "raw_timestamps": {...},
  "metadata": {
    "generated_at": "2025-01-15T10:30:00Z",
    "total_courses": 4,
    "revealed_courses": 2,
    "verification_level": "selective"
  }
}
```

## 🔍 Understanding the Output

### Merkle Tree Verification
```
✓ Student did:example:STU001: 4 courses, root: 9c456b11546e7cc6
```
- Each student's courses form a Merkle tree
- Root hash represents all course completions for that term
- Individual course proofs can be verified against this root

### Verkle Tree Aggregation
```
✅ Term published with Verkle root: 6fd9649863592b48
```
- All student term roots aggregated into single Verkle commitment
- Enables efficient verification of any student's term without revealing others
- Constant-size proofs regardless of student population

### Verification Results
```
✅ Receipt verified successfully for did:example:STU001
```
- Cryptographic proof validation passed
- Timeline consistency verified (started < completed < assessed < issued)
- Anti-forgery checks passed (issued ≤ term publication time)

## 🛡️ Security Features

### Anti-forgery Protection
1. **Timeline Validation**: Course timestamps must be in logical order
2. **Publication Deadline**: All courses must be issued before term publication
3. **Cryptographic Proofs**: Merkle and Verkle proofs prevent tampering
4. **Hash Verification**: Course data integrity protected by cryptographic hashes

### Privacy Features
1. **Selective Disclosure**: Reveal only specific courses without exposing others
2. **Zero Knowledge**: Verkle proofs don't reveal student population or other students' data
3. **Granular Control**: Students choose exactly which achievements to share

## 🧩 System Architecture

```
Course Completions
       ↓
Student-Term Merkle Trees (per student, per term)
       ↓  
Term-Level Verkle Tree (all students in term)
       ↓
Blockchain Storage (Verkle roots only)
```

### Data Flow
1. **Academic records** → Course completion entries
2. **Student level** → Merkle tree of courses (with timestamps)
3. **Term level** → Verkle tree of student term roots
4. **Blockchain** → Term Verkle roots published on-chain
5. **Verification** → Students generate receipts, verifiers validate proofs

## 🔧 Advanced Usage

### Custom Configuration

Create `~/.micert.yaml`:
```yaml
institution:
  id: "IU-Vietnam"
  name: "International University"

blockchain:
  network: "sepolia"
  contract_address: "0x..."
  
verification:
  require_timestamps: true
  max_receipt_age_days: 30
```

### Batch Operations

Process multiple students/terms:
```bash
# Add multiple terms
for term in Fall_2024 Spring_2025; do
  ./micert add-term "$term" "./data/${term,,}_completions.json"
done

# Generate receipts for all students  
./micert generate-receipt "did:example:student001" "./receipts/student001.json"
./micert generate-receipt "did:example:student002" "./receipts/student002.json"
```

### Integration with Learning Management Systems

The CLI can process data from common LMS exports:

```bash
# From Canvas LMS export
./micert add-term "Fall_2024" "./canvas_export.json" --format json

# From CSV export (custom format)
./micert add-term "Fall_2024" "./gradebook.csv" --format csv
```

## 🐛 Troubleshooting

### Common Issues

**Build Errors:**
```bash
# Clean and rebuild
go clean -cache
go mod tidy
go build -o micert ./cmd/
```

**Missing Dependencies:**
```bash
# Reinstall modules
rm go.sum
go mod download
go mod tidy
```

**Test Failures:**
```bash
# Run with verbose output
./micert test --verbose
```

### Error Messages

| Error | Solution |
|-------|----------|
| `student not found in term` | Ensure student has course completions in specified term |
| `verkle tree not published` | Run `publish-term` before generating receipts |
| `invalid timestamp ordering` | Check course completion timeline (started < completed < assessed < issued) |
| `merkle proof verification failed` | Data may be corrupted, regenerate student term tree |

## 📚 Additional Resources

- **Thesis Document**: Complete technical details and methodology
- **Smart Contracts**: `../../contracts/src/` - On-chain verification logic  
- **Test Data**: `../crypto/testdata/` - Example academic records
- **Frontend**: `../client/` - Web interface for verification

## 🤝 Contributing

This is a research implementation for academic thesis work. For questions or issues, please refer to the thesis documentation or contact the author.

## 📄 License

MIT License - See LICENSE file for details.

---

**Author**: Le Tien Phat  
**Institution**: International University - VNU HCMC  
**Thesis**: "A Transparent and Granular Blockchain System for Verifiable Academic Micro-credential Provenance"