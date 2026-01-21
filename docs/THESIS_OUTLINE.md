# IU-MiCert Thesis Outline

**Title:** A Transparent and Granular Blockchain System for Verifiable Academic Provenance

**Target Length:** 50-80 pages

---

## Report Writing Progress

| Chapter             | Status       | Notes                          |
| ------------------- | ------------ | ------------------------------ |
| Ch1. Introduction   | **DONE**     |                                |
| Ch2. Related Work   | **DONE**     |                                |
| Ch3. Methodology    | NEEDS UPDATE | Add revocation section (3.6)   |
| Ch4. Prototyping    | **DONE**     |                                |
| Ch5. Implementation | NEEDS UPDATE | System is done - needs rewrite |
| Ch6. Results        | IN PROGRESS  | Needs benchmarks               |
| Ch7. Discussion     | PENDING      | Not started                    |
| Ch8. Conclusion     | PENDING      | Not started                    |

## Implementation Progress

| Component                | Status   | Details                                                                   |
| ------------------------ | -------- | ------------------------------------------------------------------------- |
| CLI Tools (15+ commands) | **DONE** | generate-data, batch-process, verify-local, publish-roots, etc.           |
| Verkle Tree Engine       | **DONE** | go-verkle integration, IPA proofs                                         |
| Smart Contract           | **DONE** | Deployed on Sepolia: `0x2452F0063c600BcFc232cC9daFc48B7372472f79`         |
| Student Portal           | **DONE** | [https://iu-micert.vercel.app](https://iu-micert.vercel.app/)             |
| Issuer Dashboard         | **DONE** | [https://iumicert-issuer.vercel.app](https://iumicert-issuer.vercel.app/) |
| Revocation System        | **DONE** | 8 API endpoints, version-based supersession                               |
| Database Layer           | **DONE** | SQLite with GORM, migration support                                       |

---

# Chapter Outline

## Chapter 1: INTRODUCTION

- 1.1 Motivation
- 1.2 Problem Statement
- 1.3 Scope
- 1.4 Objectives
- 1.5 Thesis Organization

## Chapter 2: RELATED WORK

- 2.1 Current Advancements (BlockCerts, CVSS, EduCTX, IU-TransCert, IU-VecCert)
- 2.2 Research Gaps
- 2.3 Theoretical Background
  - 2.3.1 Blockchain Technology (Immutability, Consensus, Smart Contracts)
  - 2.3.2 Ethereum Platform and Gas Economics
  - 2.3.3 Pedersen Commitments and Inner Product Argument (IPA)
  - 2.3.4 Verkle Trees (Structure built on polynomial commitments)

## Chapter 3: METHODOLOGY

- 3.1 Requirement Analysis (Use Cases, System Overview)
- 3.2 System Architecture and Design
- 3.3 Verkle Tree Construction (NewTermVerkleTree, AddCourses)
- 3.4 Proof Generation and Verification (GenerateStudentReceipt, VerifyCourseProof, VerifyReceiptOffChain)
- 3.5 Smart Contract Integration - _add versioning_
- 3.6 Credential Revocation System - _add new section_
- 3.7 Complexity and Performance Analysis

## Chapter 4: PROTOTYPING

- 4.1 Technology Stack (Go, Next.js, Solidity)
- 4.2 Development Tools (go-verkle, Ethers.js, Foundry, GORM)
- 4.3 Blockchain Strategy (Sepolia deployment)
- 4.4 Data Pipeline Design

## Chapter 5: IMPLEMENTATION

- 5.1 System Architecture Overview
- 5.2 CLI Tool Implementation (15+ commands)
- 5.3 Smart Contract Implementation
- 5.4 Revocation System Implementation
- 5.5 Web Interface Implementation
- 5.6 Integration and Testing

## Chapter 6: RESULTS

- 6.1 Performance Benchmarks
- 6.2 Security Analysis
- 6.3 Cost Analysis
- 6.4 Comparison with Existing Solutions
- 6.5 Usability Evaluation

## Chapter 7: DISCUSSION

- 7.1 Interpretation of Results
- 7.2 Research Contributions
- 7.3 Limitations
- 7.4 Future Work

## Chapter 8: CONCLUSION

- 8.1 Summary of Contributions
- 8.2 Research Impact
- 8.3 Final Remarks

---

# Appendices (Planned)

- **A.** Smart Contract Source Code
- **B.** API Documentation
- **C.** Sample Receipt JSON
- **D.** User Testing Questionnaire
- **E.** Benchmark Raw Data

---

# Technical Notes

## Why go-verkle Uses IPA Instead of KZG

| Factor               | KZG             | IPA (Pedersen) | Winner |
| -------------------- | --------------- | -------------- | ------ |
| Trusted Setup        | Required        | Not needed     | IPA    |
| Proof Size           | O(1)            | O(log n)       | KZG    |
| Verification         | Fast (pairings) | Slower         | KZG    |
| SNARK Compatibility  | Limited         | Better (Halo)  | IPA    |
| Quantum Upgrade Path | Harder          | Easier         | IPA    |

### Key Reasons

1. **No Trusted Setup** - KZG requires ceremony where secret `s` must be destroyed. If compromised, proofs can be forged. IPA avoids this entirely.

2. **Ethereum's "The Verge" Roadmap** - Plans to SNARK-ify everything. IPA works with Halo-style techniques for constant-time verification without trusted setup.

3. **Design Trade-off** - KZG is simpler internally, but IPA's no-trusted-setup wins for decentralized systems.

4. **Tree Width** - IPA uses 256-width (2^8), KZG could use 2^24. Smaller width is fine for this use case.

### References

- [Vitalik: Verkle Trees (2021)](https://vitalik.eth.limo/general/2021/06/18/verkle.html)
- [Anatomy of a Verkle Proof](https://ihagopian.com/posts/anatomy-of-a-verkle-proof)
- [Verification of Verkle Tree](https://medium.com/@chaisomsri96/statelessness-series-part3-verification-of-verkle-tree-7b9207790c49)
- [Halo and SNARKs without pairings](https://vitalik.ca/general/2021/11/05/halo.html)

---

_Last Updated: 3 December 2025_
