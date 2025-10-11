⏺ IU-MiCert Issuer System Demo Guide

🎯 Demo Overview

This demo showcases the complete lifecycle of academic credential management
using the IU-MiCert system - from publishing new terms to generating verifiable
receipts and blockchain integration.

📋 Prerequisites

- Go 1.21+ installed
- Terminal/CLI access
- Navigate to: packages/issuer/

🎪 Demo Script

Step 1: Check System Status

# Verify the system is initialized

go run cmd/\*.go version

# Check existing terms

ls blockchain_ready/roots/

What to highlight:

- System is ready with existing 6 terms processed
- Blockchain-ready architecture in place

---

Step 2: Publish New Academic Term 🎓

# Add a new academic term with course completions

go run cmd/\*.go add-term "Fall_2025"
"data/converted_terms/Semester_1_2023_completions.json"

Expected Output:
📚 Adding academic term: Fall_2025
📖 Processing data from: data/converted_terms/Semester_1_2023_completions.json
📊 Loaded 21 course completions
🌳 Building student-term Merkle trees...
✓ Built Merkle tree for did:example:ITITIU00001: 3 courses
✓ Built Merkle tree for did:example:ITITIU00002: 4 courses
[... more students ...]
🔗 Preparing Verkle tree aggregation...
✅ Verkle root:
20b04358e69318369690a87a504e179acdcb7fc04d3be7c04d865b26bbb37f45
✅ Term added successfully!

Key Points:

- ✅ Processed 21 course completions for 5 students
- ✅ Built individual Merkle trees per student
- ✅ Generated term-level Verkle tree commitment
- ✅ Created blockchain-ready root data

---

Step 3: Generate Updated Student Receipt 📜

# Generate comprehensive academic journey receipt

go run cmd/\*.go generate-receipt "ITITIU00001" "receipts/demo_receipt.json"

Expected Output:
👤 Generating receipt for student: ITITIU00001
📚 Auto-discovered terms: [Fall_2025 Semester_1_2023 Semester_1_2024 ...]
🔐 Generating academic journey receipt...
✓ Generated receipt for term Fall_2025 (3 courses)
✓ Generated receipt for term Semester_1_2023 (3 courses)
[... more terms ...]
💾 Receipt saved to: receipts/demo_receipt.json
✅ Receipt generated successfully!

Key Points:

- ✅ Auto-discovered all student terms (now 7 total)
- ✅ Includes the newly published Fall_2025 term
- ✅ Generated cryptographic proofs for verification
- ✅ Saved to both custom and blockchain_ready locations

---

Step 4: Display Verifiable Academic Journey 🎓

# Display the receipt in human-readable format

go run cmd/\*.go display-receipt "receipts/demo_receipt.json" --blockchain

Expected Output:
📋 ACADEMIC JOURNEY RECEIPT
=============================================================
👤 Student: ITITIU00001
📅 Generated: 2025-08-26T08:23:00+07:00
📖 Type: Complete Academic Journey

📚 ACADEMIC TIMELINE:

---

[1] 📖 Fall_2025
⛓️ Blockchain Root: 20b04358e6931836...
📋 Courses (3 completed): 1. IT153IU - Discrete Mathematics [A-] (3 credits) 2. PH013IU - Physics 1 [B-] (2 credits) 3. PH013IU - Physics 1 [C] (2 credits)
📊 Term Summary: 3 courses, 7 credits, GPA: 2.80

[... 6 more terms ...]

==============================================================
📊 JOURNEY SUMMARY:
🎓 Total Courses: 21
📚 Total Credits: 67
📈 Overall GPA: 2.60

⛓️ BLOCKCHAIN VERIFICATION:

---

📋 Terms Published: 7
🔗 Blockchain Network: Sepolia Testnet
✅ All terms cryptographically anchored on blockchain

Key Points:

- ✅ Complete academic timeline with new term included
- ✅ Blockchain anchors for each term
- ✅ Comprehensive GPA and credit calculations
- ✅ Verification instructions for third parties

---

Step 5: Publish to Blockchain 🔗

# Prepare blockchain transaction for the new term

go run cmd/\*.go publish-roots "Fall_2025"

Expected Output:
⛓️ Publishing roots for term: Fall*2025
🌐 Target network: sepolia
🌳 Loading Verkle tree commitment...
✓ Verkle root:
20b04358e69318369690a87a504e179acdcb7fc04d3be7c04d865b26bbb37f45
📡 Preparing blockchain transaction...
💰 Estimating gas costs...
📡 [SIMULATION] Connecting to blockchain...
📨 [SIMULATION] Broadcasting transaction...
✅ Term roots prepared for blockchain publishing!
📄 Transaction data saved: blockchain_ready/transactions/tx_Fall_2025*\*.json
🔗 [SIMULATION] Transaction hash: 0x68ad0c89

Key Points:

- ✅ Prepared Ethereum transaction for Sepolia testnet
- ✅ Estimated gas costs and transaction parameters
- ✅ Generated transaction data for blockchain deployment
- ✅ Ready for production blockchain integration

---

Step 6: Verify Receipt Locally 🔍

# Demonstrate local verification without blockchain

go run cmd/\*.go verify-local "receipts/demo_receipt.json"

Expected Output:
🔍 Verifying receipt: receipts/demo_receipt.json
📋 Verifying receipt for student: ITITIU00001
🔐 Validating Verkle proofs...
✓ Term Fall_2025: Verkle root 20b04358e6931836...
✓ Term Semester_1_2023: Verkle root 20b04358e6931836...
[... more terms ...]
⏰ Checking temporal consistency...
✅ Local verification successful!

---

Step 7: Selective Disclosure Demo 🔒

# Generate privacy-preserving receipt with specific courses

go run cmd/\*.go generate-receipt "ITITIU00001" "receipts/selective_demo.json" \
 --selective \
 --courses "IT153IU,PH013IU" \
 --terms "Fall_2025"

# Display the selective receipt

go run cmd/\*.go display-receipt "receipts/selective_demo.json"

Key Points:

- ✅ Privacy-preserving credential sharing
- ✅ Selective disclosure of specific courses
- ✅ Maintains cryptographic integrity
- ✅ Enables targeted verification

---

🎯 Demo Key Messages

For Academic Institutions:

- Seamless Integration: Easy term publishing workflow
- Automated Processing: Bulk student processing capabilities
- Cryptographic Security: Tamper-proof academic records
- Blockchain Ready: Global verification infrastructure

For Students:

- Complete Journey: Comprehensive academic timeline
- Privacy Control: Selective disclosure options
- Instant Verification: Real-time receipt generation
- Global Recognition: Blockchain-anchored credentials

For Verifiers:

- Cryptographic Proofs: Mathematical verification
- Blockchain Anchored: Immutable record verification
- Easy Validation: Simple CLI verification commands
- Third-party Tools: Standard verification instructions

📊 Demo Statistics

After completing the demo:

- Terms Processed: 7 academic terms
- Students: 5 active students
- Course Completions: 21+ individual records
- Blockchain Roots: 7 Verkle tree commitments
- Receipts Generated: Multiple format options
- Verification Methods: Local + Blockchain ready

---

🚀 Next Steps Discussion

1. Production Deployment: Real blockchain integration
2. Institution Integration: LMS/SIS system connections
3. Student Portal: Web interface for receipt access
4. Verifier Tools: Third-party verification APIs
5. Mobile Apps: Student credential wallet applications

This demo showcases a complete academic credential lifecycle with cutting-edge
cryptographic security and blockchain integration! 🎉
