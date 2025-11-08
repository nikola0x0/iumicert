# IU-MiCert: Verifiable Academic Micro-Credential Provenance System

> Enhancing credential verification through verifiable academic micro-credential provenance using Verkle trees

## Overview

IU-MiCert is a blockchain-based system that transforms how academic credentials are verified by treating **each course as an independent micro-credential** with complete temporal provenance. Unlike traditional systems that verify entire degrees, IU-MiCert enables term-by-term verification while maintaining a tamper-proof timeline of academic achievements.

### Key Innovation

**Academic Provenance Timeline**: Every course completion is cryptographically anchored to a specific term, creating an immutable learning journey that prevents backdating or timeline manipulation.

## The Problem

Current academic credential systems:

- Treat courses as isolated achievements without temporal context
- Focus on whole-degree verification, not individual competencies
- Cannot efficiently verify specific skills or course completions
- Use Merkle trees that become inefficient at scale
- Lack mechanisms to prevent credential timeline manipulation

## Our Solution

IU-MiCert introduces:

1. **Micro-Credential Architecture**: Each course = one verifiable credential
2. **Verkle Tree Technology**: Compact cryptographic proofs (32 bytes vs. traditional approaches)
3. **Term-Based Provenance**: Immutable academic timeline preventing backdating
4. **Selective Disclosure**: Share specific courses without revealing full transcript

## System Architecture

```
LMS Data → Verkle Trees (per term) → Blockchain Anchoring → Verifiable Receipts
```

Each academic term becomes one Verkle tree containing all course completions for that period. Tree roots are published to Ethereum, enabling independent verification.

## Live Demonstration

- **Student/Verifier Portal**: [https://iu-micert.vercel.app](https://iu-micert.vercel.app)
- **Issuer Dashboard**: [https://iumicert-issuer.vercel.app](https://iumicert-issuer.vercel.app)
- **Smart Contracts**: [View on Sepolia Testnet](https://sepolia.etherscan.io/address/0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60)

## Research Contributions

1. **Verkle-based credential system** with term-level granularity
2. **Academic provenance model** preventing timeline manipulation
3. **Efficient verification** with 32-byte proofs per course
4. **Production deployment** demonstrating real-world feasibility

---

**Technology Stack**: Ethereum, Verkle Trees, Go, Next.js
**Institution**: International University - Vietnam National University HCMC
**License**: MIT
