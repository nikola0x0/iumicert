# IU-MiCert: Blockchain-Based Verifiable Academic Micro-Credential Provenance System

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Blockchain](https://img.shields.io/badge/Blockchain-Ethereum-3C3C3D.svg)](https://ethereum.org/)
[![Verkle Trees](https://img.shields.io/badge/Data_Structure-Verkle_Trees-green.svg)](#)
[![Research Paper](https://img.shields.io/badge/📄-Research_Paper-orange.svg)](#)

> **Enhancing credential verification through verifiable academic micro-credential provenance using Verkle trees**

## 📖 Overview

IU-MiCert is a blockchain-based system designed to address critical limitations in current academic credential verification by implementing verifiable micro-credential provenance. The system leverages **Verkle tree technology** as an improvement over traditional Merkle trees to provide:

- 🛡️ Enhanced anti-forgery mechanisms through temporal verification
- 📈 Efficient storage and verification of granular learning achievements
- 🔗 Seamless integration with existing credential management systems
- ⏳ Comprehensive audit trail of learning achievements with verifiable timestamps

## 🎯 Motivation

Traditional credential systems fail to capture the rich mosaic of skills, projects, and competencies that constitute modern education. IU-MiCert addresses this by:

- Providing cryptographic verification of micro-credential sequences
- Preventing timeline manipulation and backdating of achievements
- Supporting the shift toward lifelong learning and skills-based assessment
- Serving as an upgrade to existing blockchain credential systems

## 🔍 Problem Statement

Current blockchain credential systems:

- Focus primarily on whole degree verification
- Lack efficient mechanisms for micro-credential management
- Create opportunities for credential forgery through timeline gaps
- Use Merkle trees that become inefficient at scale

IU-MiCert implements **Verkle tree structures** to address these challenges while maintaining compatibility with existing systems.

## 🏗️ System Architecture

```mermaid
graph TD
    A[Learning Management Systems] --> B[IU-MiCert Core]
    B --> C[Verkle Tree Structure]
    C --> D[Blockchain Storage]
    D --> E[Verification Protocols]
    E --> F[Student Interface]
    E --> G[Employer Interface]
    E --> H[Institution Interface]
```

### Key Components:

- **Verkle Tree Manager**: Efficient storage/verification of micro-credentials
- **Smart Contracts**: Automated credential issuance with term-based cycles
- **Commitment Engine**: Minimizes on-chain storage while maintaining provenance
- **Verification Protocols**: Efficient proof validation at scale

## 🚀 Key Features

| Feature                    | Benefit                                            |
| -------------------------- | -------------------------------------------------- |
| Verkle Tree Implementation | 40-50% smaller proofs than Merkle trees            |
| Temporal Verification      | Detects backdating and timeline manipulation       |
| Micro-Credential Tracking  | Captures courses, projects, and skill achievements |
| Compatibility Layer        | Works alongside existing credential systems        |
| Intuitive Interfaces       | For students, employers, and institutions          |

## 📚 Thesis Chapters

1. **Introduction**: Motivation, problem statement, and objectives
2. **Literature Review**: Blockchain credentials, Verkle trees, academic provenance
3. **System Design**: Architecture, data structures, and protocols
4. **Implementation**: Smart contracts, Verkle tree manager, interfaces
5. **Evaluation**: Security analysis, performance benchmarks, case studies
6. **Conclusion**: Findings and future work

## 🛠️ Technical Stack

- **Blockchain**: Ethereum (with potential for other EVM chains)
- **Smart Contracts**: Solidity
- **Cryptography**: Verkle tree implementation using established libraries
- **Frontend**: Next.js with React (for demo interfaces)
- **Backend**: Node.js (for institutional integration)
- **Styling**: Tailwind CSS

## 📂 Repository Structure

```
iumicert/
├── contracts/                    # Smart contract source code
├── packages/
│   ├── client/                  # Client-side applications
│   │   └── iumicert-client/     # Next.js web application for verification
│   └── issuer/                  # Credential issuing system
│       ├── data/                # Test data and academic records
│       └── backup_scripts/      # Backup and maintenance scripts
├── LICENSE                      # MIT License
└── README.md                   # This file
```

## 🔬 Evaluation Metrics

- **Security**: Resistance to credential forgery attempts
- **Performance**: Proof generation/verification times vs Merkle trees
- **Storage Efficiency**: On-chain footprint comparison
- **Usability**: Stakeholder feedback on interfaces

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📝 Citation

If you use this work in your research, please cite:

```bibtex
@thesis{phat2026iumicert,
  title={IU-MiCert: Blockchain-Based Verifiable Academic Micro-Credential Provenance System},
  author={Le Tien Phat},
  year={2026},
  institution={International University - Vietnam National University HCM},
  type={Bachelor's Thesis}
}
```

## ✉️ Contact

For questions about this research, please contact:

**Le Tien Phat**  
📧 ltphat.developer@gmail.com  
🏫 International University - Vietnam National University HCM

- IUMiCertRegistry: 0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60
- ReceiptRevocationRegistry: 0x8814ae511d54Dc10C088143d86110B9036B3aa92
- Owner: 0xf16221da98b931409195A395b290333edA85f90F (your wallet)
