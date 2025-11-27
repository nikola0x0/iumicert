# Smart Contracts

## Contract Architecture

IU-MiCert uses two main smart contracts deployed on Ethereum Sepolia testnet.

### IUMiCertRegistry

**Purpose**: Store term roots and manage credential lifecycle

**Key Functions**:

```solidity
// Publish a term's Verkle tree root
function publishRoot(string termId, bytes32 root) external onlyIssuer

// Retrieve a term's root
function getRoot(string termId) external view returns (bytes32)

// Check if term exists
function termExists(string termId) external view returns (bool)
```

**Deployment**: `0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60`

**View on Etherscan**: [Sepolia Explorer](https://sepolia.etherscan.io/address/0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60)

### ReceiptRevocationRegistry

**Purpose**: Future feature for credential revocation

**Deployment**: `0x8814ae511d54Dc10C088143d86110B9036B3aa92`

**Status**: Implemented but not yet integrated

## Contract Design Principles

### Gas Efficiency

- Minimal storage (only term roots, ~32 bytes each)
- Efficient data structures (mappings)
- Batch operations where possible

### Access Control

- **Issuer role**: Only authorized institution can publish roots
- **Public reading**: Anyone can verify roots
- **Upgradeability**: Proxy pattern for future enhancements

### Data Integrity

- **Immutable roots**: Once published, cannot be changed
- **Event emissions**: Transparent logging of all operations
- **Timestamp tracking**: Automatic blockchain timestamp for each root

## Root Publication Flow

1. **Institution**: Builds term Verkle tree
2. **Commitment**: Calculates tree root (32 bytes)
3. **Transaction**: Calls `publishRoot(termId, root)`
4. **Blockchain**: Permanently stores root + timestamp
5. **Event**: `RootPublished` event emitted

## Verification Flow

1. **Student**: Presents receipt with proof
2. **Verifier**: Extracts term ID from receipt
3. **Contract Query**: Calls `getRoot(termId)`
4. **Verification**: Validates proof against retrieved root
5. **Result**: ✅ Valid or ❌ Invalid

## Security Considerations

### Attack Prevention

| Attack Vector | Mitigation |
|---------------|------------|
| Fake roots | Only issuer can publish (access control) |
| Root modification | Immutable once published |
| Replay attacks | Each term has unique ID |
| Gas griefing | Rate limiting at API layer |

### Trust Model

- **Trust issuer address**: Institution controls private key
- **Trust blockchain**: Ethereum consensus guarantees
- **Don't trust**: Individual receipts (must verify against blockchain)

## Gas Costs

Approximate costs on Ethereum Sepolia:

- `publishRoot()`: ~50,000 gas (~$0.10 at 20 gwei)
- `getRoot()`: 0 gas (view function)
- **Per term**: One-time cost for publishing root
- **Per student**: Zero on-chain cost (verification off-chain)

## Example Transaction

Publishing Semester_1_2023 root:

```
Function: publishRoot(string termId, bytes32 root)
Parameters:
  - termId: "Semester_1_2023"
  - root: 0x1a2b3c...def (32 bytes)
Gas Used: 48,234
Timestamp: 2024-01-15 10:23:45 UTC
Tx Hash: 0xabc123...
```

## Future Enhancements

- **Batch publishing**: Multiple terms in one transaction
- **Revocation**: Mark specific credentials as revoked
- **Multi-issuer**: Support consortium of institutions
- **Layer 2**: Deploy to L2 for lower costs

---

**Next**: Live demonstration and use cases.
