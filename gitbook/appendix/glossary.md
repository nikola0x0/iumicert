# Glossary

## Core Concepts

**Academic Provenance**
: The verifiable chronological record of a student's learning achievements, including when each credential was earned and in what context.

**Micro-Credential**
: A small, focused credential representing a specific skill or achievement, typically smaller in scope than a full degree (e.g., a single course completion).

**Selective Disclosure**
: The ability to reveal only specific parts of one's academic record without exposing the entire transcript.

**Temporal Integrity**
: Assurance that the timeline of achievements cannot be manipulated or backdated.

## Cryptography

**Commitment Scheme**
: A cryptographic primitive that allows one to commit to a value while keeping it hidden, with the ability to reveal the value later.

**Hash Function**
: A mathematical function that maps data of arbitrary size to fixed-size values, used in Merkle trees for data integrity.

**Inner Product Argument (IPA)**
: A type of cryptographic proof used in Verkle trees to verify that data belongs to a committed polynomial.

**Merkle Tree**
: A tree data structure where each leaf node contains a data hash and each non-leaf node contains the hash of its children. Used for efficient data verification.

**Polynomial Commitment**
: A cryptographic commitment to a polynomial, used in Verkle trees as an alternative to traditional hashing.

**Verkle Tree**
: A cryptographic tree structure using polynomial commitments instead of hashes, enabling constant-size membership proofs.

**Zero-Knowledge Proof**
: A cryptographic method to prove possession of information without revealing the information itself.

## Blockchain

**Gas**
: Transaction fee on Ethereum, measured in wei/gwei. Compensates validators for computational resources.

**Smart Contract**
: Self-executing code deployed on a blockchain that automatically enforces agreements.

**Testnet**
: A blockchain network used for testing, with no real economic value (e.g., Sepolia).

**Transaction**
: A signed data package storing a message to be sent from an externally owned account to another account on the blockchain.

**Wei/Gwei**
: Smallest unit of Ether (1 ETH = 10^18 wei = 10^9 gwei).

## System Components

**Issuer**
: The educational institution that creates and publishes academic credentials.

**Receipt**
: A JSON file containing course details and cryptographic proof, given to students for verification purposes.

**Root Commitment**
: The single 32-byte value representing an entire Verkle tree, published to the blockchain.

**Term**
: An academic period (e.g., Semester_1_2023) for which a single Verkle tree is constructed.

**Verifier**
: Any party (employer, institution, etc.) checking the authenticity of a credential.

## Protocols & Standards

**DID (Decentralized Identifier)**
: A W3C standard for decentralized digital identifiers not controlled by any central authority.

**EIP (Ethereum Improvement Proposal)**
: A design document providing information to the Ethereum community, including new features or processes.

**JSON-LD**
: JSON for Linked Data, a method of encoding linked data using JSON.

**Verifiable Credential (VC)**
: W3C standard for expressing credentials on the web in a cryptographically secure, privacy-respecting manner.

## Acronyms

- **API**: Application Programming Interface
- **CORS**: Cross-Origin Resource Sharing
- **EVM**: Ethereum Virtual Machine
- **IPA**: Inner Product Argument
- **LMS**: Learning Management System
- **REST**: Representational State Transfer
- **RPC**: Remote Procedure Call
- **W3C**: World Wide Web Consortium

---

**Helpful Resource**: For deeper technical explanations, see the [Technical Documentation](technical-docs.md).
