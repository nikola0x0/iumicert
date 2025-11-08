# Related Work

## Existing Blockchain Credential Systems

### Blockcerts (MIT Media Lab)

- **Approach**: Issues credentials as blockchain transactions
- **Limitation**: Focuses on whole credentials, not course-level
- **Technology**: Bitcoin/Ethereum with Merkle proofs

### Open Badges / Verifiable Credentials (W3C)

- **Approach**: Decentralized identifiers with JSON-LD credentials
- **Limitation**: No temporal provenance, requires trust in issuer timestamps
- **Technology**: DIDs, Verifiable Credentials specification

### Ethereum Attestation Service (EAS)

- **Approach**: On-chain attestations for various claims
- **Limitation**: High gas costs for granular credentials
- **Technology**: Ethereum smart contracts

## Academic Credential Research

- **Traditional PKI**: Certificate authorities, revocation lists
- **Smart Diploma**: Ethereum-based diploma verification
- **CredenceLedger**: Private blockchain for academic records

**Gap**: None of these systems combine:

1. Course-level granularity
2. Temporal provenance
3. Efficient verification
4. Selective disclosure

## Cryptographic Data Structures

### Merkle Trees

- **Advantages**: Proven security, wide adoption
- **Limitations**: O(log n) proof size, path disclosure reveals tree structure

### Verkle Trees

- **Innovation**: Constant-size proofs (~32 bytes)
- **Efficiency**: Better for large datasets
- **Adoption**: Ethereum Verkle tree upgrade proposal (EIP-6800)

**Why Verkle**: Our system manages thousands of course records per student - Verkle trees provide efficiency while maintaining security.

## Our Contribution

IU-MiCert uniquely combines:

- **Verkle trees** for efficient academic record management
- **Term-based architecture** for provenance timeline
- **Course-level credentials** as independent micro-credentials
- **Production deployment** on Ethereum testnet

The next section explains Verkle trees in accessible terms.
