# IU-MiCert Data Flow: From LMS to Blockchain Verification

This document explains the complete data flow in the IU-MiCert system, from raw LMS data to blockchain-anchored verifiable credentials.

---

## Overview

IU-MiCert transforms university Learning Management System (LMS) data into cryptographically verifiable micro-credentials using Verkle trees and blockchain anchoring. The system treats **each course as an independent micro-credential** and maintains a **tamper-proof academic provenance timeline**.

---

## Phase 1: LMS Data Input

### What We Assume from the LMS

**Real-World Scenario:**
Universities use Learning Management Systems (LMS) like:

- **IU Vietnam**: Edusoft, Blackboard
- **Other institutions**: Moodle, Canvas, Banner

**LMS Provides:**

- End-of-term student records
- Course completions with grades
- Timestamps (course start, completion, assessment, credential issuance)
- Student identifiers from LMS system

**In Our System:**

- We **simulate** this with a data generator
- In production, this would be **actual exports from university LMS**

### Generated Student Journey Format

**Location**: `data/student_journeys/students/journey_ITITIU00001.json`

```json
{
  "student_id": "did:example:ITITIU00001",
  "terms": {
    "Semester_1_2023": {
      "courses": [
        {
          "course_id": "CH011IU",
          "course_name": "Chemistry for Engineers",
          "grade": "A+",
          "credits": 3,
          "instructor": "Prof. White",
          "issuer_id": "IU-PHYS",
          "started_at": "2023-08-31T05:21:44Z",
          "completed_at": "2023-11-06T12:46:30Z",
          "assessed_at": "2023-11-09T12:46:30Z",
          "issued_at": "2023-11-23T12:46:30Z",
          "attempt_no": 1,
          "term_id": "Semester_1_2023"
        }
      ],
      "gpa": 3.5,
      "total_credits": 12
    }
  }
}
```

**Key Fields:**

- **student_id**: DID (Decentralized Identifier) - see [DID Explanation](#did-decentralized-identifiers)
- **Temporal data**: Four timestamps showing credential lifecycle
- **Course metadata**: ID, name, grade, credits, instructor
- **Issuer tracking**: Which department issued the credential

---

## DID (Decentralized Identifiers)

### Why DIDs Instead of Plain Student IDs?

**Student ID**: `ITITIU00001` (just an identifier)

**DID**: `did:example:ITITIU00001`

**Format**: `did:method:identifier`

- `did` - Prefix indicating this is a Decentralized Identifier
- `example` - Method/namespace (in production: `iu` or `iu-vietnam`)
- `ITITIU00001` - The actual student ID from LMS

**Advantages:**

1. **Global Uniqueness**: `did:iu:12345` vs `did:mit:12345` (different universities, same ID)
2. **Standardization**: W3C Verifiable Credentials standard compliance
3. **Interoperability**: Other systems recognize DIDs as structured identifiers
4. **Future-Proof**: Enables cross-institution credential verification
5. **Context Preservation**: The `did:iu:` prefix tells you this is from IU's system

**In Production:**

- Use the LMS's existing student ID
- Wrap it in DID format: `did:iu:STUDENT_ID`
- Maintains compatibility while adding global uniqueness

---

## Phase 2: Data Conversion to Verkle Format

### Command

```bash
./micert convert-data Semester_1_2023
```

### Input

`data/student_journeys/students/*.json` (all student journeys)

### Output

`data/verkle_terms/Semester_1_2023_completions.json`

### Transformation

Extracts all course completions for a specific term from all students:

```json
[
  {
    "student_id": "ITITIU00001",
    "course_id": "CH011IU",
    "course_name": "Chemistry for Engineers",
    "grade": "A+",
    "credits": 3,
    "term_id": "Semester_1_2023",
    "started_at": "2023-08-31T05:21:44Z",
    "completed_at": "2023-11-06T12:46:30Z",
    "assessed_at": "2023-11-09T12:46:30Z",
    "issued_at": "2023-11-23T12:46:30Z"
  },
  {
    "student_id": "ITITIU00002",
    "course_id": "CH011IU",
    "course_name": "Chemistry for Engineers",
    "grade": "B+",
    "credits": 3,
    "term_id": "Semester_1_2023",
    ...
  }
]
```

**Purpose**: Organize data by academic term for Verkle tree construction.

---

## Phase 3: Verkle Tree Construction

### Command

```bash
./micert add-term Semester_1_2023 data/verkle_terms/Semester_1_2023_completions.json
```

### Architecture

**One Verkle tree per academic term** containing all students' course completions for that term.

### Verkle Tree Key-Value Structure

For each course completion, we create:

#### Key Construction

```
Raw Key: "did:example:ITITIU00001:Semester_1_2023:CH011IU"
         └────-----─┬----------─┘ └────-─┬──────┘└─-─┬───┘
               Student DID            Term ID      Course

Separators are colons (:) to ensure uniqueness.

Hashed Key: SHA256(raw key) = 32-byte fixed-length key
            0x529dde14ab6fa8f3d16e8990a3eb3422790912f2c18f1537f8f8d6a0622bd3
```

**Why Hash?**

- **Verkle trees require 32-byte fixed-length keys** (for polynomial commitments)
- **Security**: Prevents key manipulation
- **Privacy**: Original key is hidden in the tree
- **Efficiency**: Uniform key size for cryptographic operations

#### Value Construction

```
Course Data: {
  "course_id": "CH011IU",
  "grade": "A+",
  "credits": 3,
  "started_at": "2023-08-31T05:21:44Z",
  "completed_at": "2023-11-06T12:46:30Z",
  "assessed_at": "2023-11-09T12:46:30Z",
  "issued_at": "2023-11-23T12:46:30Z",
  ...
}

Hashed Value: SHA256(JSON serialized data) = 32-byte commitment
              0x0510313122f10f7e812f8d851fd9ce316ac4808962b6921d2b023fb671216a4d
```

**Why Hash the Value?**

- **Integrity**: Any change to course data changes the hash
- **Privacy**: Actual data not stored in tree structure, only commitment
- **Fixed-length**: Required for Verkle tree cryptography

### Verkle Tree Structure

Each term's Verkle tree contains:

- **Leaf nodes**: (key_hash → value_hash) pairs for each course completion
- **Internal nodes**: Polynomial commitments computed from children
- **Root commitment**: 32-byte cryptographic commitment representing entire term

**Example**: Semester_1_2023 tree with 15 courses (3 students × 5 courses each):

```
Verkle Root: 0x425f67ea3a169e47e51c3d9057940dddc8793ced4dee4f5ccf93f5e87265160b
  ├─ Student ITITIU00001 courses (5 courses)
  ├─ Student ITITIU00002 courses (5 courses)
  └─ Student ITITIU00003 courses (5 courses)
```

### Outputs

**1. Complete Verkle Tree**
`data/verkle_trees/Semester_1_2023_verkle_tree.json`:

```json
{
  "term_id": "Semester_1_2023",
  "verkle_root": [66, 95, 103, 234, ...],  // 32-byte array
  "published_at": "2025-10-20T09:05:44+07:00",
  "version": "1.0.0",
  "course_entries": {
    "did:example:ITITIU00001:Semester_1_2023:CH011IU": {
      "course_id": "CH011IU",
      "student_id": "ITITIU00001",
      "grade": "A+",
      ...
    }
  }
}
```

**2. Blockchain-Ready Root**
`publish_ready/roots/root_Semester_1_2023.json`:

```json
{
  "ready_for_blockchain": true,
  "term_id": "Semester_1_2023",
  "verkle_root": "425f67ea3a169e47e51c3d9057940dddc8793ced4dee4f5ccf93f5e87265160b",
  "timestamp": "2025-10-20T09:05:44+07:00",
  "total_students": 5
}
```

---

## Phase 4: Receipt Generation

### Command

```bash
./micert generate-receipt ITITIU00001
```

### What Happens

For each course the student completed:

1. **Locate the course** in the Verkle tree
2. **Generate membership proof**: Cryptographic proof that this course exists in the tree
3. **Extract StateDiff**: The key-value pair being proven
4. **Package everything** into the receipt

### Receipt Structure

**Location**: `publish_ready/receipts/ITITIU00001_journey.json`

```json
{
  "student_id": "ITITIU00001",
  "blockchain_ready": true,
  "generation_timestamp": "2025-10-20T09:19:02+07:00",
  "receipt_type": {
    "selective_disclosure": false,
    "specific_courses": false,
    "specific_terms": false
  },
  "term_receipts": {
    "Semester_1_2023": {
      "receipt": {
        "verkle_root": "425f67ea3a169e47e51c3d9057940dddc8793ced4dee4f5ccf93f5e87265160b",
        "revealed_courses": [
          {
            "course_id": "CH011IU",
            "course_name": "Chemistry for Engineers",
            "grade": "A+",
            "credits": 3,
            "started_at": "2023-08-31T05:21:44Z",
            "completed_at": "2023-11-06T12:46:30Z",
            "assessed_at": "2023-11-09T12:46:30Z",
            "issued_at": "2023-11-23T12:46:30Z"
          }
        ],
        "course_proofs": {
          "CH011IU": {
            "course_key": "did:example:ITITIU00001:Semester_1_2023:CH011IU",
            "course_id": "CH011IU",
            "verkle_proof": {
              "commitmentsByPath": ["0x4ec4e765...", "0x618c0d97..."],
              "d": "0x161653d6...",
              "depthExtensionPresent": "0x0a",
              "ipaProof": {
                "cl": ["0x70506cef...", ...],  // 8 commitments
                "cr": ["0x111378d4...", ...],  // 8 commitments
                "finalEvaluation": "0x073a9765..."
              },
              "otherStems": []
            },
            "state_diff": [
              {
                "stem": "0x529dde14ab6fa8f3d16e8990a3eb3422790912f2c18f1537f8f8d6a0622bd3",
                "suffixDiffs": [
                  {
                    "suffix": 187,
                    "currentValue": "0x0510313122f10f7e812f8d851fd9ce316ac4808962b6921d2b023fb671216a4d",
                    "newValue": null
                  }
                ]
              }
            ]
          }
        },
        "metadata": {
          "verification_level": "full_ipa",
          "total_courses": 5,
          "revealed_courses": 5
        }
      }
    }
  }
}
```

### What the Receipt Contains

**Per Term:**

1. **verkle_root**: 32-byte Verkle tree root commitment (to be checked against blockchain)
2. **revealed_courses**: Actual course data (grades, timestamps, etc.)
3. **course_proofs**: Cryptographic proofs for each course
   - **verkle_proof**: IPA proof with commitments
   - **state_diff**: The key-value pair being proven (stem + suffix + value)
   - **course_key**: Original unhashed key for verification

**Why This Structure?**

- **Separates data from proofs**: Makes selective disclosure possible
- **Independent proofs**: Each course has its own 32-byte proof
- **Cryptographic binding**: StateDiff is bound to VerkleProof via IPA verification
- **Blockchain anchor**: verkle_root can be verified on-chain

---

## Phase 5: Blockchain Publication

### Command

Via web interface or `./micert publish-roots Semester_1_2023`

### What Gets Stored On-Chain

**Smart Contract**: `IUMiCertRegistry` at `0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60`

**Function Call**:

```solidity
publishTermRoot(
  bytes32 _verkleRoot,      // 32-byte Verkle root
  string memory _termId,    // "Semester_1_2023"
  uint256 _totalStudents    // 5
)
```

**Stored Data Structure**:

```solidity
struct TermRoot {
  string termId;           // "Semester_1_2023"
  uint256 totalStudents;   // 5
  uint256 publishedAt;     // block.timestamp
  bool exists;             // true
}

mapping(bytes32 => TermRoot) public termRoots;
// verkleRoot → TermRoot
```

**Critical Point**: Only the **Verkle root (32 bytes) + metadata** goes on-chain, NOT the actual student data or courses.

### Transaction Record

**Location**: `publish_ready/transactions/tx_<hash>.json`

```json
{
  "transaction_hash": "0x8a7e5c2...",
  "block_number": 12345678,
  "term_id": "Semester_1_2023",
  "verkle_root": "425f67ea3a169e47e51c3d9057940dddc8793ced4dee4f5ccf93f5e87265160b",
  "gas_used": 125000,
  "timestamp": "2025-10-20T09:30:15+07:00"
}
```

---

## Phase 6: Verification Process

### The Two-Part Verification

When a verifier receives a student's receipt, verification happens in two stages:

### Part 1: Blockchain Verification (On-Chain)

**What the Verifier Does:**

1. **Receives receipt JSON** (from student)

2. **Extracts verkle_root** from the receipt

3. **Queries smart contract**:

   ```solidity
   getTermRootInfo(verkleRoot)
   → returns (string termId, uint256 totalStudents, uint256 publishedAt, bool exists)
   ```

4. **Security Checks**:

   - Does this verkle_root exist on blockchain? (`exists == true`)
   - Was it officially published by the institution?
   - Does the term_id in receipt match blockchain? (prevents term swapping)
   - When was it published? (timestamp verification)

5. **Rejection conditions**:
   - Verkle root not found on blockchain
   - Term ID mismatch between receipt and blockchain record
   - Root not published by authorized institution

**Security Rationale:**

This step establishes the trust anchor by verifying the Verkle root was officially published by the institution. Without this check, an attacker could construct their own Verkle tree with fabricated courses and generate valid cryptographic proofs against it. The blockchain verification ensures only institution-published roots are accepted, combining authority verification with cryptographic integrity.

**Implementation**: `cmd/api_server.go` lines 1412-1483

### Part 2: Cryptographic Verification (Off-Chain)

**What the Verifier Does:**

1. **Now using blockchain-verified root**, for each course:

   - **Recomputes course key hash**:

     ```
     Raw key: "did:example:ITITIU00001:Semester_1_2023:CH011IU"
     Key hash: SHA256(raw key) = 0x529dde14...
     ```

   - **Recomputes course value hash**:

     ```
     Course data: {course_id, grade, credits, timestamps...}
     Value hash: SHA256(JSON) = 0x05103131...
     ```

   - **Checks these appear in proof's StateDiff**:
     - Finds matching stem and suffix in StateDiff
     - Verifies the value in StateDiff matches computed hash

2. **IPA Verification** (the critical cryptographic step):
   - Uses the `verkle_proof` to **reconstruct the tree root** from the StateDiff
   - Algorithm: `PreStateTreeFromProof(proof, stateDiff)`
   - Compares reconstructed root with `verkle_root` in the receipt
   - **If they match** → proof is cryptographically valid, data wasn't tampered

**Why This Works:**

- Verifier **doesn't trust the verkle_root blindly**
- **Cryptographic proof reconstructs the root** - tampering causes mismatch
- Uses go-verkle's IPA (Inner Product Argument) verification
- Impossible to forge a valid proof without the original tree

**Code**: `packages/crypto/verkle/membership_verifier.go` (`VerifyMembershipProof`)

### Two-Layer Verification Model

| Step  | Verification Type         | What It Proves                                     | Protects Against                                       |
| ----- | ------------------------- | -------------------------------------------------- | ------------------------------------------------------ |
| **1** | **Blockchain** (on-chain) | Root was officially published                      | Fake Verkle trees, unauthorized issuers, fake receipts |
| **2** | **Cryptographic** (IPA)   | Receipt data matches the blockchain-verified proof | Data tampering within legitimate receipts              |

**Security Design:**

The verification order is critical. Blockchain verification occurs first to establish a trust anchor, confirming the Verkle root was published by the authorized institution. Only then does cryptographic verification proceed, ensuring the receipt data matches the blockchain-verified root. This two-layer approach prevents both unauthorized credential creation and data tampering within legitimate receipts.

**Combined Result**: The receipt contains authentic data from a legitimate, authorized source.

---

## Selective Disclosure

### How It Works

**Full Receipt:**

- Student receives receipt with **all courses** from all terms
- Contains complete academic history

**Selective Disclosure:**

1. Student **manually edits the receipt JSON**
2. **Removes unwanted courses** from `revealed_courses` array
3. **Removes corresponding proofs** from `course_proofs` object
4. **Keeps the verkle_root unchanged**
5. Shares filtered receipt with verifier

**Example:**

```json
// Original: 30 courses across 6 terms
// Filtered: Only IT courses from 2023 (5 courses)

{
  "term_receipts": {
    "Semester_1_2023": {
      "receipt": {
        "verkle_root": "425f67ea...", // SAME ROOT
        "revealed_courses": [
          // Only IT courses, removed chemistry/math
        ],
        "course_proofs": {
          // Only proofs for IT courses
        }
      }
    }
    // Removed other terms entirely
  }
}
```

### Why Verkle Trees Enable This

**Merkle Trees**:

- Need sibling nodes for entire path to root
- Removing data breaks the verification path
- Can't selectively disclose without revealing structure

**Verkle Trees**:

- **Each course has independent 32-byte proof**
- Removing courses doesn't affect other course proofs
- All remaining proofs **still verify against the same root**
- Verifier only sees what student shares

**Security:**

- **Cannot add fake courses**: Each proof must cryptographically verify
- **Cannot change grades**: Would change value hash, proof fails
- **Cannot backdate**: Timestamps are in the verified course data
- Root remains unchanged and verified on blockchain

---

## Complete Data Flow Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│ Phase 1: LMS Data                                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  University LMS (Edusoft/Blackboard)                           │
│  │                                                              │
│  ├─→ End-of-term student records                               │
│  └─→ Course completions with timestamps                        │
│                                                                  │
│         ↓ (simulated by generate-data)                         │
│                                                                  │
│  Student Journey JSON                                          │
│  - DID: did:example:ITITIU00001                                │
│  - Terms → Courses → Grades + Timestamps                       │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────────┐
│ Phase 2: Conversion (Per Term)                                 │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  convert-data Semester_1_2023                                  │
│  │                                                              │
│  ├─→ Extract all course completions for this term              │
│  └─→ From all students                                         │
│                                                                  │
│         ↓                                                        │
│                                                                  │
│  Verkle Format: Semester_1_2023_completions.json               │
│  [course1, course2, course3, ...]                              │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────────┐
│ Phase 3: Verkle Tree Construction                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  add-term Semester_1_2023                                      │
│  │                                                              │
│  For each course completion:                                   │
│  ├─→ Key = SHA256("did:iu:ID:term:course")  → 32 bytes        │
│  ├─→ Value = SHA256(course JSON data)       → 32 bytes        │
│  └─→ Insert into Verkle tree                                   │
│                                                                  │
│         ↓                                                        │
│                                                                  │
│  Verkle Tree Built                                             │
│  ├─ Root: 0x425f67ea... (32-byte commitment)                  │
│  └─ All student courses cryptographically committed            │
│                                                                  │
│         ↓                                                        │
│                                                                  │
│  Output: root_Semester_1_2023.json (ready for blockchain)     │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────────┐
│ Phase 4: Receipt Generation (Per Student)                      │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  generate-receipt ITITIU00001                                  │
│  │                                                              │
│  For each of student's courses:                                │
│  ├─→ Generate 32-byte Verkle membership proof                  │
│  ├─→ Extract StateDiff (key-value pair)                        │
│  └─→ Package: proof + data + root                              │
│                                                                  │
│         ↓                                                        │
│                                                                  │
│  Student Receipt                                               │
│  ├─ revealed_courses: [course data]                            │
│  ├─ course_proofs: {courseId → proof + StateDiff}             │
│  └─ verkle_root: for each term                                 │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────────┐
│ Phase 5: Blockchain Publication                                │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Web UI: Connect MetaMask + Publish Term                       │
│  │                                                              │
│  ├─→ publishTermRoot(verkleRoot, termId, totalStudents)       │
│  └─→ Only 32-byte root + metadata goes on-chain                │
│                                                                  │
│         ↓                                                        │
│                                                                  │
│  Ethereum Blockchain (Sepolia)                                 │
│  ├─ Verkle Root: 0x425f67ea...                                │
│  ├─ Term ID: "Semester_1_2023"                                 │
│  ├─ Students: 5                                                │
│  └─ Published: timestamp                                       │
│                                                                  │
│  Storage: ~200 bytes per term (extremely efficient!)           │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────────┐
│ Phase 6: Verification                                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Student shares receipt → Verifier receives JSON               │
│                                                                  │
│  STEP 1: Blockchain Verification (Trust Anchor)                │
│  ├─→ Extract verkle_root from receipt                          │
│  ├─→ Query smart contract:                                     │
│  │   getTermRootInfo(verkle_root)                              │
│  │                                                              │
│  ├─→ Security Checks:                                          │
│  │   ├─ Does root exist on blockchain? (exists == true)        │
│  │   ├─ Does term_id match? (prevents term swapping)           │
│  │   └─ When was it published? (timestamp)                     │
│  │                                                              │
│  ├─→ If any check fails:                                       │
│  │   └─ Verification terminates with error                     │
│  │                                                              │
│  └─→ All checks pass:                                          │
│      Root confirmed as officially published by institution     │
│                                                                  │
│  Purpose of blockchain verification:                           │
│  - Establishes trust anchor before cryptographic verification  │
│  - Prevents use of attacker-generated Verkle trees             │
│  - Confirms institutional authority (owner-only publication)   │
│                                                                  │
│  STEP 2: Cryptographic Verification (Data Integrity)           │
│  ├─→ Now using blockchain-verified root:                       │
│  │   For each course:                                          │
│  │   ├─ Recompute key hash from courseKey                      │
│  │   ├─ Recompute value hash from course data                  │
│  │   ├─ Check these exist in proof's StateDiff                 │
│  │   └─ IPA verification:                                      │
│  │       ├─ Use proof to reconstruct tree root                 │
│  │       ├─ Compare with blockchain-verified root              │
│  │       └─ Match confirms data integrity                      │
│  │                                                              │
│  └─→ Result: All course proofs cryptographically verified      │
│                                                                  │
│  VERIFICATION COMPLETE (Both Layers Passed)                    │
│  - Layer 1: Root verified on blockchain (institutional source) │
│  - Layer 2: Data verified via IPA proofs (data integrity)      │
│  - Academic provenance confirmed                               │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Key Concepts Explained

### Academic Provenance

**Definition**: Complete, verifiable timeline of a student's learning achievements that cannot be backdated or manipulated.

**How IU-MiCert Achieves This:**

1. **Term-by-term organization**: Each academic period is a separate Verkle tree
2. **Temporal integrity**: Four timestamps per course (started, completed, assessed, issued)
3. **Cryptographic binding**: All data is cryptographically committed in the tree
4. **Blockchain anchoring**: Roots published at specific times, creating immutable timeline
5. **Cannot backdate**: Once published, timeline is locked on blockchain

### Micro-Credentials

**Traditional**: Verify entire degree as one credential

**IU-MiCert**: Each course is an **independent micro-credential**

- Separate 32-byte cryptographic proof per course
- Can verify individual achievements
- Enables granular verification (specific courses, not full transcript)
- Supports lifelong learning (add credentials over time)

### Why Hashing is Critical

**Problem**: Verkle trees use polynomial commitments requiring fixed-length inputs

**Solution**: Hash everything to 32 bytes

1. **Keys** (variable length):

   ```
   "did:example:ITITIU00001:Semester_1_2023:CH011IU" (52 chars)
   → SHA256 → 0x529dde14... (32 bytes)
   ```

2. **Values** (variable length JSON):
   ```
   {"course_id":"CH011IU","grade":"A+",...} (200+ chars)
   → SHA256 → 0x05103131... (32 bytes)
   ```

**Benefits:**

- • Fixed-length for polynomial commitments
- • Tamper-evident (any change breaks the hash)
- • Privacy (actual data hidden in tree)
- • Efficient verification

---

## Data Storage Architecture

### File System (Development & Backup)

```
packages/issuer/
│
├── data/
│   ├── student_journeys/          # Phase 1: Raw LMS data (JSON files)
│   │   ├── students/              # Individual student journeys
│   │   └── terms/                 # Term summaries
│   │
│   ├── verkle_terms/              # Phase 2: Converted data (JSON files)
│   │   └── Semester_*_completions.json
│   │
│   └── verkle_trees/              # Phase 3: Built Verkle trees (JSON files)
│       └── Semester_*/verkle_tree.json
│
└── publish_ready/
    ├── roots/                     # Phase 3: Blockchain-ready roots (JSON files)
    │   └── root_Semester_*.json
    │
    ├── receipts/                  # Phase 4: Student receipts (JSON files)
    │   └── ITITIU*_journey.json
    │
    └── transactions/              # Phase 5: Blockchain records (JSON files)
        └── tx_*.json
```

### Database (Production Storage)

**PostgreSQL Database** with the following tables:

**1. Students Table:**
- Student ID, DID, name, email
- Enrollment date, expected graduation
- Status (active, graduated, withdrawn)

**2. Terms Table:**
- Term ID (Semester_1_2023)
- Start/end dates
- **Verkle root** (hex string + binary)
- **Blockchain transaction hash** and block number
- Published timestamp

**3. TermReceipts Table** (Main Receipt Storage):
- Receipt ID (unique per student-term)
- Student ID, Term ID
- **VerkleProof** (JSONB - full proof structure with IPA data)
- **StateDiff** (JSONB - key-value pairs)
- **RevealedCourses** (JSONB - course completion array)
- Verkle root hex (for quick lookups)
- Blockchain verification status
- Transaction hash and block number

**4. AccumulatedReceipts Table:**
- Multi-term receipts (progress reports, diplomas)
- References to multiple term receipts
- Aggregated statistics (GPA, total credits)

**5. BlockchainTransactions Table:**
- Transaction hash, term ID
- Verkle root, block number
- Gas used, status (pending/confirmed)
- Timestamps

**6. VerificationLog Table:**
- Receipt ID being verified
- Verifier information
- Verification mode (local, blockchain, full_ipa)
- Success/failure with error messages
- Timestamp and IP tracking

### Dual Storage Strategy

**Why Both Files and Database?**

**JSON Files** (file system):
- • Portable receipts students can download/share
- • Easy backup and version control
- • Human-readable for debugging
- • Compatible with offline verification tools

**PostgreSQL Database**:
- • Fast queries (find student by ID, get all receipts for term)
- • JSONB for efficient proof storage and querying
- • Blockchain metadata tracking
- • Verification logging and analytics
- • Production-ready scalability

**Data Flow**: Generate in files → Import to database → Serve via API

---

## Automated Pipeline

### Complete Generation Script

```bash
./reset.sh && ./generate.sh
```

**What This Does:**

1. **Reset**: Clear all previous data
2. **Generate**: Create 5 students × 6 terms
3. **Convert**: Transform each term to Verkle format
4. **Build**: Create 6 Verkle trees with cryptographic commitments
5. **Receipts**: Generate verifiable credentials for all students

**Time**: ~10-15 seconds for complete dataset

### Individual Commands

```bash
# Step-by-step manual process
./micert generate-data              # Create student journeys
./micert batch-process              # Process all terms
./micert generate-all-receipts      # Create all receipts
./micert publish-roots              # Publish to blockchain (requires key)
```

---

## Security Properties

### Tamper-Proof Timeline

**Cannot modify past achievements** because:

1. Course data is **hashed into Verkle tree**
2. Root is **published on blockchain** with timestamp
3. **Any change** to course data changes the hash
4. Changed hash means **proof verification fails**
5. Blockchain record is **immutable**

### Cryptographic Integrity

**Cannot forge credentials** because of **two-layer security**:

**Layer 1 - Blockchain (Authority Check):**

1. Only institution can publish roots (smart contract owner-only function)
2. Verifier checks root exists on-chain BEFORE accepting any proofs
3. Prevents attacker from creating their own Verkle tree with fake courses
4. Even with valid cryptographic proofs, fake root is rejected

**Layer 2 - IPA Verification (Data Integrity Check):**

1. Given blockchain-verified root, verify course data matches proofs
2. Reconstruct root from proofs - must match blockchain root
3. Cannot tamper with course data without breaking cryptographic verification
4. Impossible to generate valid proof without original Verkle tree

**Attack Scenarios Prevented:**

- **Fake tree attack**: Creating a Verkle tree with fabricated courses
  - Prevented by: Blockchain verification (root not published by institution)
- **Data tampering**: Modifying grades in legitimate receipt
  - Prevented by: IPA verification (reconstructed root mismatch)
- **Course addition**: Adding extra courses to legitimate receipt
  - Prevented by: IPA verification (no valid proof for added course)
- **Backdating**: Altering timestamps in course data
  - Prevented by: IPA verification (value hash changes, proof invalid)

### Privacy Preservation

**Selective disclosure** without compromising security:

1. Each course has **independent proof**
2. Removing courses doesn't break other proofs
3. All remaining proofs verify against **same root**
4. Verifier only sees what student chooses to share

---

## System Advantages

### vs Traditional Credentials

- • **Granular verification**: Check specific courses, not whole degree
- • **No institution contact**: Verify independently via blockchain
- • **Tamper-proof**: Cryptographically impossible to forge
- • **Privacy**: Student controls what to share

### vs Merkle Trees

- • **Smaller proofs**: 32 bytes vs O(log n)
- • **Selective disclosure**: Independent proofs per course
- • **Efficient verification**: Constant-time, not logarithmic
- • **Better privacy**: No structural information leaked

### vs Other Blockchain Credentials

- • **Micro-credential support**: Course-level granularity
- • **Temporal integrity**: Complete provenance timeline
- • **Scalability**: Constant proof size regardless of transcript size
- • **Storage efficient**: Only 32-byte roots on-chain

---

## Technical Stack

- **Verkle Trees**: `ethereum/go-verkle` v0.2.2
- **Cryptography**: IPA (Inner Product Arguments) for polynomial commitments
- **Blockchain**: Ethereum Sepolia testnet
- **Smart Contract**: `IUMiCertRegistry` at `0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60`
- **Backend**: Go 1.21+ with Cobra CLI
- **API**: REST server on port 8080
- **Frontend**: Next.js with wagmi/viem for MetaMask

---

## Related Documentation

- **Main README**: `../README.md` - Complete system overview
- **Setup Guide**: `setup.md` - Local development setup
- **Deployment**: `deployment.md` - Production deployment
- **Demo Guide**: `demo-guide.md` - System demonstration
- **Implementation**: `../../../docs/implementation-guide.md` - Cryptographic details
- **Theory**: `../../../docs/mathematical-foundation.md` - Verkle tree mathematics

---

**System Status**: Production-ready with full IPA verification and blockchain-anchored security
