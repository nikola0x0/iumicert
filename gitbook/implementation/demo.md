# Live Demonstration

## Deployed System

IU-MiCert is fully deployed and operational on Ethereum Sepolia testnet.

### Web Applications

**Student/Verifier Portal**
🌐 [https://iu-micert.vercel.app](https://iu-micert.vercel.app)

- Upload receipt JSON
- Verify cryptographic proofs
- View course details
- Check blockchain anchoring

**Issuer Dashboard**
🌐 [https://iumicert-issuer.vercel.app](https://iumicert-issuer.vercel.app)

- Generate demo data
- Build Verkle trees
- Publish roots to blockchain
- Create student receipts

### Smart Contracts

**IUMiCertRegistry**
📜 [0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60](https://sepolia.etherscan.io/address/0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60)

## Test Data

The system includes realistic test data:

- **5 students**: ITITIU00001 through ITITIU00005
- **6 terms**: Semester_1_2023 through Summer_2024
- **~30 courses per student**: Real IU Vietnam course codes
- **Vietnamese names**: Authentic student profiles

## Demo Workflow

### 1. Verify a Receipt

1. Visit [student portal](https://iu-micert.vercel.app)
2. Upload sample receipt (available in issuer dashboard)
3. System automatically:
   - Validates Verkle proof
   - Checks blockchain for term root
   - Verifies timestamps
   - Displays result

### 2. Generate New Receipts

1. Visit [issuer dashboard](https://iumicert-issuer.vercel.app)
2. Navigate to "Demo Data" page
3. Generate new student journeys
4. Build Verkle trees
5. Create verifiable receipts
6. Download and test

### 3. Explore Blockchain

1. Visit contract on [Etherscan](https://sepolia.etherscan.io/address/0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60)
2. View `RootPublished` events
3. See term commitments
4. Verify transaction timestamps

## Sample Use Cases

### Use Case 1: Job Application

**Scenario**: Software engineering position requires data structures knowledge

1. Student downloads receipt for "Data Structures" course
2. Uploads to employer's verification system
3. Employer sees:
   - ✅ Verified completion
   - ✅ Grade: A
   - ✅ Blockchain-anchored
   - ✅ Issued: Semester_1_2023

### Use Case 2: Graduate School Application

**Scenario**: Master's program requires prerequisite courses

1. Student exports receipts for required courses
2. Submits to admissions office
3. Office verifies instantly:
   - No need to contact undergraduate institution
   - Cryptographically proven authenticity
   - Complete timeline of prerequisite completion

### Use Case 3: Skills Portfolio

**Scenario**: Building verifiable skills portfolio

1. Student collects micro-credentials across years
2. Creates portfolio of verified achievements
3. Shares selectively based on job requirements
4. Each credential independently verifiable

## Performance Metrics

From deployed system:

- **Proof size**: 32 bytes per course
- **Verification time**: <100ms
- **API response**: ~200ms average
- **Blockchain query**: <500ms
- **Storage per term**: ~32 bytes on-chain

## Try It Yourself

1. **Visit**: [https://iu-micert.vercel.app](https://iu-micert.vercel.app)
2. **Get sample receipt**: Download from issuer dashboard
3. **Verify**: Upload and see instant verification
4. **Explore**: Check blockchain transaction on Etherscan

## System Status

- ✅ Smart contracts deployed and operational
- ✅ Issuer system generating receipts
- ✅ Verification portal live
- ✅ Test data available
- ✅ Blockchain anchoring working
- 🔄 Revocation feature in development

---

**Next**: Security analysis and evaluation.
