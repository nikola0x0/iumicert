# Architecture Overview

## System Design

IU-MiCert follows a **term-centric architecture** where each academic period (semester) is independently managed and verifiable.

## High-Level Flow

```
1. LMS Data → 2. Verkle Trees → 3. Blockchain → 4. Verifiable Receipts
```

### 1. Data Ingestion

- Academic records exported from Learning Management System (LMS)
- Organized by academic term (e.g., Semester_1_2023)
- Each course includes: student ID, course code, grade, timestamps

### 2. Verkle Tree Construction

- **One tree per term** containing all student course completions
- Each course = one leaf in the tree
- Tree commitment (root) represents entire term

### 3. Blockchain Anchoring

- Term roots published to Ethereum smart contract
- Immutable timestamp of credential issuance
- Enables independent verification without trusting issuer

### 4. Receipt Generation

Students receive **verifiable receipts** containing:

- Course details (selectively disclosed)
- Cryptographic proof (32 bytes)
- Blockchain reference (term root)

## Core Components

### Issuer System

**Role**: Educational institution
**Responsibilities**:

- Process LMS data
- Build Verkle trees per term
- Publish roots to blockchain
- Generate student receipts

**Technology**: Go backend with CLI interface

### Smart Contracts

**Functions**:

- `publishRoot(termId, root)`: Anchor term commitment
- `getRoot(termId)`: Retrieve published root
- Revocation registry (future feature)

**Deployment**: Ethereum Sepolia testnet

### Verification Portal

**Role**: Students, employers, verifiers
**Capabilities**:

- Upload receipt JSON
- Verify cryptographic proofs
- Check blockchain anchoring
- View course details

**Technology**: Next.js web application

## Data Flow Example

**Scenario**: Verify student completed "Data Structures" in Fall 2023

1. **Student** uploads receipt to verification portal
2. **System** extracts Verkle proof and term root
3. **Verification**:
   - Validate cryptographic proof against term root
   - Check term root exists on blockchain
   - Confirm timestamps match
4. **Result**: ✅ Verified or ❌ Invalid

## Security Model

### Threat Protection

| Threat | Mitigation |
|--------|-----------|
| Fake credentials | Cryptographic proofs verifiable against blockchain |
| Backdating | Term-based architecture prevents timeline manipulation |
| Tampering | Receipt modifications invalidate cryptographic proofs |
| Impersonation | Student identity tied to Verkle tree structure |

### Trust Assumptions

- **Blockchain integrity**: Ethereum network security
- **Issuer authenticity**: Institution's smart contract address
- **Cryptography**: Verkle tree and IPA security proofs

## Scalability

- **Per-term trees**: Manageable tree sizes (~1000-5000 students × ~6 courses = ~30,000 leaves)
- **Constant proofs**: 32 bytes regardless of transcript length
- **Parallel verification**: Each receipt independently verifiable

---

**Next**: Learn about how micro-credentials work in IU-MiCert.
