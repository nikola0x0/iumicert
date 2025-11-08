# Technology Stack

## Overview

IU-MiCert uses a modern, production-ready technology stack combining blockchain, cryptography, and web technologies.

## Core Technologies

### Blockchain Layer

**Ethereum (Sepolia Testnet)**

- **Why**: Proven security, wide adoption, smart contract support
- **Usage**: Anchor term roots, permanent timestamp record
- **Contract**: IUMiCertRegistry at `0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60`

### Cryptography

**Verkle Trees (go-verkle library)**

- **Implementation**: Ethereum Foundation's reference implementation
- **Proofs**: Inner Product Arguments (IPA)
- **Security**: Same guarantees as Merkle trees, constant-size proofs

### Backend System

**Go 1.21+**

- **Why**: Performance, strong cryptography libraries, CLI tooling
- **Components**:
  - Verkle tree management
  - Receipt generation
  - REST API server
  - CLI commands (15+ operations)

**PostgreSQL**

- **Usage**: Store student journeys, receipts, verification logs
- **Why**: Reliable, ACID compliance, good for relational academic data

### Frontend Applications

**Next.js + React**

- **Issuer Dashboard**: Institution admin interface
- **Student/Verifier Portal**: Public verification interface
- **Deployment**: Vercel (serverless)

**Styling**: TailwindCSS + shadcn/ui components

### Smart Contracts

**Solidity**

- **Development**: Foundry framework
- **Deployment**: Ethereum Sepolia testnet
- **Features**:
  - Root publication
  - Term management
  - Revocation registry (future)

## Architecture Diagram

```
┌─────────────┐
│  LMS Data   │
└──────┬──────┘
       │
       v
┌─────────────────────────────┐
│   Issuer System (Go)        │
│  - Verkle Tree Builder      │
│  - Receipt Generator        │
│  - CLI Interface            │
│  - REST API                 │
└──────┬──────────────────────┘
       │
       v
┌─────────────────────────────┐
│  Smart Contracts (Solidity) │
│  - IUMiCertRegistry         │
│  - RevocationRegistry       │
└──────┬──────────────────────┘
       │
       v
┌─────────────────────────────┐
│   Ethereum Sepolia          │
│   (Blockchain Anchor)       │
└──────┬──────────────────────┘
       │
       v
┌─────────────────────────────┐
│  Verification Portals       │
│  - Student Portal (Next.js) │
│  - Employer Portal          │
└─────────────────────────────┘
```

## Development Tools

- **Version Control**: Git + GitHub
- **Package Management**: Go modules, npm
- **Testing**: Go testing framework, Foundry
- **Deployment**: Docker, Vercel, Ethereum testnet
- **Monitoring**: Etherscan, API logs

## Production Considerations

### Scalability

- **Database**: Connection pooling, indexed queries
- **API**: CORS-enabled, rate limiting ready
- **Blockchain**: Gas-efficient smart contracts

### Security

- **Private keys**: Environment variable management
- **API**: Input validation, error handling
- **Smart contracts**: Audited patterns, access control

### Reliability

- **Error handling**: Graceful degradation
- **Data validation**: Schema enforcement
- **Backup**: Database migrations, blockchain redundancy

## Why These Choices?

| Technology | Reason |
|------------|--------|
| Go | Performance, strong typing, excellent crypto libraries |
| Verkle Trees | Constant-size proofs, Ethereum-aligned |
| PostgreSQL | ACID compliance, relational data model |
| Next.js | SEO-friendly, serverless deployment, React ecosystem |
| Ethereum | Decentralized, proven security, smart contract platform |
| Foundry | Modern Solidity tooling, fast testing |

---

**Next**: Smart contract implementation details.
