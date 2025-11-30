# Credential Revocation System

## Overview

The IU-MiCert system supports credential revocation through a **term versioning** approach. When a credential needs to be revoked (e.g., due to grade correction, academic misconduct discovery, or administrative error), the affected term's Verkle tree is rebuilt without the revoked credential and a new version is published to the blockchain.

**Key Principle**: Each academic term has its own Verkle tree. Revocation only affects the specific term containing the credential - other terms remain unchanged.

## Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                    REVOCATION WORKFLOW                              │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  1. REGISTRAR validates student complaint (offline process)        │
│                          ↓                                          │
│  2. ADMIN creates revocation request via Dashboard                  │
│                          ↓                                          │
│  3. Request stored in DB with status "approved"                     │
│                          ↓                                          │
│  4. When ANY term is published via Dashboard:                       │
│      ├── System checks for approved revocations                     │
│      ├── Rebuilds affected term's Verkle tree                       │
│      ├── Publishes new version via SupersedeTerm()                  │
│      └── Marks revocation as "processed"                            │
│                          ↓                                          │
│  5. Old receipts become invalid (root mismatch)                     │
│     Students can download new receipts                              │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

## How It Works

### Per-Term Versioning

Each term maintains its own version history on the blockchain:

| Term | Version | Root Hash | Status |
|------|---------|-----------|--------|
| Semester_1_2023 | v1 | `0xabc123...` | Superseded |
| Semester_1_2023 | v2 | `0xdef456...` | Active |
| Semester_2_2023 | v1 | `0x789abc...` | Active |

### Automatic Processing

Revocations are processed **automatically** when any term is published via the Dashboard:

1. Admin publishes `Semester_2_2024` → Dashboard calls blockchain
2. Backend detects approved revocation for `Semester_1_2023`
3. Backend rebuilds `Semester_1_2023` tree (removes revoked credential)
4. Backend calls `SupersedeTerm("Semester_1_2023", newRoot, reason)`
5. `Semester_1_2023` is now v2, revocation marked as "processed"
6. `Semester_2_2024` is published as v1

This ensures revocations are processed in the background without manual intervention.

## Smart Contract

The `IUMiCertRegistry` contract supports versioning:

```solidity
// Publish new term (v1)
function publishTerm(bytes32 termId, bytes32 root, uint256 students)

// Supersede existing term (v2, v3, ...)
function supersedeTerm(bytes32 termId, bytes32 newRoot, uint256 students, string reason)

// Check root status
function checkRootStatus(bytes32 termId, bytes32 root) returns (bool isValid, bool isSuperseded, uint256 version)

// Get version history
function getTermHistory(bytes32 termId) returns (uint256[] versions, bytes32[] roots)
```

**Contract Address (Sepolia)**: `0x2452F0063c600BcFc232cC9daFc48B7372472f79`

## Database Schema

### revocation_requests
Stores revocation requests with audit trail:

| Column | Type | Description |
|--------|------|-------------|
| request_id | VARCHAR(255) | Unique ID (revoke_req_UUID) |
| student_id | VARCHAR(50) | e.g., ITITIU00001 |
| term_id | VARCHAR(50) | e.g., Semester_1_2023 |
| course_id | VARCHAR(50) | e.g., IT089IU |
| reason | TEXT | Why revoked |
| status | VARCHAR(50) | pending/approved/processed/rejected |
| processed_at | TIMESTAMP | When processed |
| processed_by_tx_hash | VARCHAR(66) | Blockchain tx |
| processed_in_version | INTEGER | Which version |

### term_root_versions
Tracks version history per term:

| Column | Type | Description |
|--------|------|-------------|
| term_id | VARCHAR(50) | Term identifier |
| version | INTEGER | 1, 2, 3... |
| root_hash | VARCHAR(66) | Verkle root (0x...) |
| is_superseded | BOOLEAN | If newer version exists |
| superseded_by | VARCHAR(66) | Next version's root |
| credentials_revoked | INTEGER | Count removed |
| tx_hash | VARCHAR(66) | Blockchain tx |

### revocation_batches
Records batch processing:

| Column | Type | Description |
|--------|------|-------------|
| batch_id | VARCHAR(255) | Unique batch ID |
| term_id | VARCHAR(50) | Affected term |
| old_version | INTEGER | Previous version |
| new_version | INTEGER | New version |
| request_count | INTEGER | Revocations in batch |
| tx_hash | VARCHAR(66) | Blockchain tx |

## API Endpoints

All under `/api/issuer/revocations`:

### Create Request
```http
POST /api/issuer/revocations
Content-Type: application/json

{
  "student_id": "ITITIU00001",
  "term_id": "Semester_1_2023",
  "course_id": "IT089IU",
  "reason": "Grade correction - incorrect entry",
  "requested_by": "admin",
  "notes": "Validated by registrar on 2025-11-30"
}
```

### List Requests
```http
GET /api/issuer/revocations?status=approved
GET /api/issuer/revocations?term_id=Semester_1_2023
```

### Get Statistics
```http
GET /api/issuer/revocations/stats
```

Response:
```json
{
  "approved_requests": 2,
  "processed_requests": 10,
  "pending_requests": 0,
  "total_batches": 5
}
```

### Process Manually (Optional)
```http
POST /api/issuer/revocations/process
```

Triggers immediate processing of all approved revocations.

### Delete Request
```http
DELETE /api/issuer/revocations/{request_id}
```

### Get Term Versions
```http
GET /api/issuer/terms/{term_id}/versions
```

## CLI Commands

### Process Revocations
```bash
./micert process-revocations
```

Processes all approved revocations across all terms.

### Publish with Revocation Check
```bash
./micert publish-roots Semester_1_2024
```

Automatically checks and processes approved revocations before publishing.

## Dashboard Usage

### Creating a Revocation Request

1. Navigate to **Revocation Management** tab
2. Fill in the form:
   - Student ID (e.g., `ITITIU00003`)
   - Term (select from dropdown)
   - Course ID (e.g., `IT013IU`)
   - Reason (description)
3. Click **Submit Request**
4. Request appears with "APPROVED" status

### Processing Revocations

Revocations are processed automatically when you:
1. Go to **Publish Term** tab
2. Select any term and publish to blockchain
3. System processes all pending revocations in background

### Monitoring

- **Refresh button** in Revocation Management to see updated status
- **Statistics dashboard** shows counts by status
- **Version History** endpoint shows term version progression

## Verification Impact

When a credential is revoked:

1. **Old receipts fail verification**
   - Root hash no longer matches active version
   - `checkRootStatus()` returns `isSuperseded: true`

2. **Students must download new receipts**
   - New receipts contain updated proofs
   - New root hash matches active version

3. **Verifiers see clear status**
   - API indicates if root is superseded
   - Provides reason and timestamp

## Example Workflow

### Scenario: Student ITITIU00003's IT013IU grade was incorrectly entered

**Step 1**: Registrar validates the complaint (offline)

**Step 2**: Admin creates revocation request
```bash
POST /api/issuer/revocations
{
  "student_id": "ITITIU00003",
  "term_id": "Semester_1_2023",
  "course_id": "IT013IU",
  "reason": "Grade entry error - validated by registrar",
  "requested_by": "admin"
}
```

**Step 3**: Admin publishes next term (e.g., Semester_2_2023)
- System detects approved revocation
- Rebuilds Semester_1_2023 tree (18 → 17 entries)
- Calls `SupersedeTerm("Semester_1_2023", newRoot, "Revoked 1 credential")`
- Marks revocation as "processed"

**Step 4**: Student downloads new receipt
- Old receipt for Semester_1_2023 is invalid
- New receipt has updated root and proofs
- IT013IU no longer appears in receipt

## Technical Details

### Tree Rebuilding Process

```
1. Load term's Verkle tree from disk
   └─ data/verkle_trees/Semester_1_2023_verkle_tree.json

2. Remove revoked credential from CourseEntries
   └─ delete("did:example:ITITIU00003:Semester_1_2023:IT013IU")

3. Rebuild Verkle tree
   └─ New root hash generated

4. Publish to blockchain
   └─ SupersedeTerm() with new root

5. Update database
   └─ TermRootVersion, RevocationRequest, RevocationBatch

6. Save updated tree to disk
   └─ Overwrites existing file
```

### Gas Costs (Sepolia)

| Operation | Gas Used |
|-----------|----------|
| publishTerm (v1) | ~180,000 |
| supersedeTerm (v2+) | ~320,000 |

### Configuration

Required in `.env`:
```
ISSUER_PRIVATE_KEY=0x...  # For signing supersedeTerm transactions
CONTRACT_ADDRESS=0x2452F0063c600BcFc232cC9daFc48B7372472f79
NETWORK=sepolia
```

## Limitations

1. **Batch processing only**: Individual credential revocation requires full tree rebuild
2. **No real-time revocation**: Processed during term publication
3. **Student notification**: Students must check for new receipts (no push notification)
4. **Gas costs**: Each supersedeTerm costs ~320k gas

## Future Enhancements

- Real-time revocation list (without tree rebuild)
- Student notification system
- Revocation certificate generation
- Audit log export
