# IU-MiCert Thesis Defense Presentation

**Title:** A Transparent and Granular Blockchain System for Verifiable Academic Provenance
**Duration:** 15 minutes (presentation + demo) | 15 minutes Q&A
**Author:** Nikola - IU Vietnam

---

## TABLE OF CONTENTS

### Main Presentation Slides

- [Slide 1: Title Slide](#slide-1-title-slide)
- [Slide 2: Agenda](#slide-2-agenda)
- [Slide 3: The Problem](#slide-3-the-problem)
- [Slide 4: Research Question & Objectives](#slide-4-research-question--objectives)
- [Slide 5: Proposed Solution](#slide-5-proposed-solution)
- [Slide 6: What is a Verkle Tree?](#slide-6-what-is-a-verkle-tree)
- [Slide 7: Verkle Trees in IU-MiCert](#slide-7-verkle-trees-in-iu-micert)
- [Slide 8: System Architecture](#slide-8-system-architecture)
- [Slide 9: How It Works - Issuance Flow](#slide-9-how-it-works---issuance-flow)
- [Slide 10: How It Works - Verification Flow](#slide-10-how-it-works---verification-flow)
- [Slide 11: Demo - Issuer Dashboard](#slide-11-demo---issuer-dashboard)
- [Slide 12: Demo - Student/Verifier Portal](#slide-12-demo---studentverifier-portal)
- [Slide 13: Demo - Selective Disclosure](#slide-13-demo---selective-disclosure)
- [Slide 14: Demo - Credential Revocation](#slide-14-demo---credential-revocation)
- [Slide 15: Results - Performance](#slide-15-results---performance)
- [Slide 16: Comparison with Existing Systems](#slide-16-comparison-with-existing-systems)
- [Slide 17: Research Contributions](#slide-17-research-contributions)
- [Slide 18: Conclusion & Future Work](#slide-18-conclusion--future-work)

### Q&A Backup Slides

- [Q&A 1: IPA vs KZG Deep Dive](#qa-slide-1-ipa-vs-kzg-deep-dive)
- [Q&A 2: Revocation Mechanism](#qa-slide-2-revocation-mechanism)
- [Q&A 3: Security Guarantees](#qa-slide-3-security-guarantees)
- [Q&A 4: Gas Costs Detail](#qa-slide-4-gas-costs-detail)
- [Q&A 5: Scalability Analysis](#qa-slide-5-scalability-analysis)
- [Q&A 6: Privacy Considerations](#qa-slide-6-privacy-considerations)
- [Q&A 7: Comparison with ZK Solutions](#qa-slide-7-comparison-with-zk-solutions)
- [Q&A 8: Real-World Deployment](#qa-slide-8-real-world-deployment)
- [Q&A 9: Technical Architecture Details](#qa-slide-9-technical-architecture-details)
- [Q&A 10: Limitations Acknowledgment](#qa-slide-10-limitations-acknowledgment)

### Appendix

- [Presentation Tips](#presentation-tips)

---

## PRESENTATION OVERVIEW

| Section                             | Slides | Time    | Focus                    |
| ----------------------------------- | ------ | ------- | ------------------------ |
| Opening                             | 1-2    | 1 min   | Title & Agenda           |
| I. Research Background & Motivation | 3-4    | 2 min   | Problem & objectives     |
| II. Proposed Solution               | 5-7    | 2 min   | Verkle trees explanation |
| III. Methodology                    | 8-10   | 2.5 min | Architecture & flows     |
| IV. Research Results                | 11-16  | 6.5 min | Demo & benchmarks        |
| V. Conclusion & Discussions         | 17-18  | 2 min   | Contributions & future   |

**Total: 18 slides + 10 Q&A backup slides**

---

# MAIN PRESENTATION SLIDES

---

## SLIDE 1: Title Slide

**Time: 30 seconds**

### Content:

```
IU-MiCert: A Transparent and Granular Blockchain System
for Verifiable Academic Provenance

Presented by: [Your Name]
Supervisor: [Supervisor Name]

International University - Vietnam National University HCMC
School of Computer Science and Engineering

December 2025
```

### Visual Suggestions:

- IU Vietnam logo
- Professional academic background
- Subtle blockchain or certificate iconography

### Speaker Notes:

> "Good morning, distinguished committee members and guests. I am [Name], and I will be presenting my thesis entitled 'IU-MiCert: A Transparent and Granular Blockchain System for Verifiable Academic Provenance.'"

---

## SLIDE 2: Agenda

**Time: 30 seconds**

### Title: Presentation Outline

### Content:

```
I.   Research Background & Motivation

II.  Hypotheses Development

III. Methodology

IV.  Research Results

V.   Conclusion & Discussions
```

### Visual Suggestions:

- Professional numbered list with consistent iconography
- Section progress indicator (reusable throughout presentation)
- Linear timeline or roadmap visualization

### Speaker Notes:

> "The presentation is organized into five sections. I will begin with the research background and motivation, followed by the hypotheses development. Then I will explain the methodology and system design. The results section will include a system demonstration and performance evaluation. Finally, I will conclude with a discussion of contributions and future directions."

---

## SLIDE 3: The Problem

**Time: 1 minute**

### Title: Challenges in Credential Verification

### Content:

**Traditional Systems:**

- Verification requires direct contact with issuing institution
- Credentials vulnerable to forgery and fraud
- No mechanism for course-level verification
- **No verifiable provenance** - cannot prove when credentials were issued

**Existing Blockchain Solutions:**

- Limited to whole-certificate verification (BlockCerts, IU-TransCert)
- Cannot verify individual courses as micro-credentials
- Proof sizes grow with data volume - O(log n)
- No selective disclosure - full transcript must be revealed
- **Lack of temporal integrity** - credentials can potentially be backdated

### Visual Suggestions:

- Comparative diagram: Traditional vs. Blockchain approaches
- Visual markers indicating limitations
- Provenance gap highlighted

### Speaker Notes:

> "Traditional credential verification requires employers to contact institutions directly, which is time-consuming and costly. Critically, there is no verifiable provenance - no way to prove when credentials were actually issued. While blockchain solutions address fraud concerns, they still lack temporal integrity guarantees. Students cannot prove specific course completions, and proof sizes grow logarithmically with data."

---

## SLIDE 4: Research Question & Objectives

**Time: 1 minute**

### Title: Research Question & Objectives

### Content:

**Research Question:**
How can we design a blockchain-based credential system that addresses the limitations of existing solutions?

**Research Objectives:**

1. **Granular verification** - Enable verification of individual courses
2. **Compact proofs** - Achieve constant-size proofs regardless of data scale
3. **Temporal integrity** - Provide cryptographic proof of issuance time
4. **Selective disclosure** - Allow students to reveal only chosen credentials
5. **Cost-effective deployment** - Minimize on-chain storage costs

### Visual Suggestions:

- Research question prominently displayed
- Five objectives with corresponding icons
- Clear visual hierarchy

### Speaker Notes:

> "This research addresses a central question: How can we design a blockchain credential system that overcomes current limitations? We identified five key objectives: enabling course-level verification, achieving constant-size proofs, ensuring temporal integrity, supporting selective disclosure, and maintaining cost-effectiveness for large-scale deployment."

---

## SLIDE 5: Proposed Solution

**Time: 30 seconds**

### Title: Proposed Solution: Verkle Trees

### Content:

**Core Innovation:**

Apply Verkle tree technology - originally designed for Ethereum's stateless clients - to academic credential verification.

**Why Verkle Trees?**

- Constant-size proofs regardless of data volume
- Proven cryptographic foundation from Ethereum research
- Designed for exactly this type of problem: efficient state proofs

### Visual Suggestions:

- Transition slide connecting problem to solution
- Verkle tree visual introduction
- Ethereum research foundation reference

### Speaker Notes:

> "To address these challenges, we propose applying Verkle tree technology to academic credentials. Verkle trees were developed for Ethereum's stateless client upgrade and offer properties ideal for credential verification: constant-size proofs regardless of how many credentials exist in the system."

---

## SLIDE 6: What is a Verkle Tree?

**Time: 1 minute**

### Title: Verkle Tree Technology

### Content:

**Verkle Tree Overview:**

- Cryptographic data structure developed for Ethereum's stateless client upgrade
- Designed to reduce witness sizes for state verification (from ~1MB to ~150KB)
- Utilizes polynomial commitments (Pedersen + IPA) instead of hash-based proofs

**Verkle vs. Merkle Comparison (at 1 million items):**

| Feature      | Merkle Tree           | Verkle Tree                  |
| ------------ | --------------------- | ---------------------------- |
| Tree Height  | 20 levels (binary)    | **32 levels (256-ary width)** |
| Proof Size   | 32B × 20 = 640 bytes  | **~1.8 KB (constant)**       |
| Verification | O(20 hash operations) | **O(1) IPA verification**    |

**Key Property:** Proof size remains constant regardless of dataset scale

### Visual Suggestions:

- Side-by-side tree structure diagrams (Merkle vs Verkle)
- Proof size comparison chart
- 256-ary branching visualization

### Speaker Notes:

> "Verkle trees were developed for Ethereum's stateless client upgrade. The key innovation is the 256-ary width at each level - our implementation has 32 levels with 256 children per node. Unlike Merkle proofs that grow with tree depth, Verkle proofs remain constant at approximately 1.8 KB regardless of scale - this constancy is critical for scalable credential systems."

---

## SLIDE 7: Verkle Trees in IU-MiCert

**Time: 30 seconds**

### Title: IU-MiCert: Pedersen + IPA Variant

### Content:

**Our Choice: Pedersen Commitments with IPA Proofs**

- Uses Pedersen vector commitments for tree nodes
- Inner Product Argument (IPA) for proof generation
- No trusted setup ceremony required

**Why This Variant?**

| Factor           | KZG Variant                     | Pedersen + IPA (Our Choice)  |
| ---------------- | ------------------------------- | ---------------------------- |
| Trusted Setup    | Required (multi-party ceremony) | **Not needed**               |
| Proof Size       | O(1) - 48 bytes                 | O(log n) - ~1.8 KB           |
| If Compromised   | All proofs forgeable            | **No such risk**             |
| Long-term Safety | Depends on ceremony integrity   | **Independent of any setup** |

**Critical for Credentials:** Must remain valid for 50+ years without setup compromise risk

### Visual Suggestions:

- Pedersen + IPA commitment diagram
- Comparison with KZG variant
- Trust model visualization

### Speaker Notes:

> "IU-MiCert specifically uses the Pedersen commitment with IPA proof variant of Verkle trees. This choice is critical because it requires no trusted setup ceremony. For academic credentials that must remain valid for decades, we cannot rely on a setup ceremony that could be compromised. Each term has one tree with its root anchored to the blockchain."

---

## SLIDE 8: System Architecture

**Time: 1 minute**

### Title: System Architecture

### Content:

```
┌─────────────────────────────────────────┐
│           PRESENTATION LAYER            │
│  Issuer Dashboard | Student/Verifier    │
└─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│           APPLICATION LAYER             │
│         REST API + CLI Tools            │
└─────────────────────────────────────────┘
                    │
       ┌────────────┼────────────┐
       ▼            ▼            ▼
┌───────────┐ ┌──────────┐ ┌──────────────┐
│  Verkle   │ │ Database │ │  Ethereum    │
│  Engine   │ │ (SQLite) │ │  Blockchain  │
└───────────┘ └──────────┘ └──────────────┘
```

**Core Components:**

- **Backend (Go):** Verkle tree construction and proof generation
- **Smart Contract:** Term root storage on Ethereum Sepolia
- **Web Interfaces:** Next.js portals for issuers and verifiers

### Visual Suggestions:

- Three-tier architecture diagram
- Color-coded layers with technology labels
- Component interaction arrows

### Speaker Notes:

> "The system architecture has three main layers. At the user layer, we have issuers, students, and verifiers. The issuer platform includes the CLI and dashboard for credential management, with the Verkle engine handling cryptographic operations and PostgreSQL storing credential data. The client application provides the verification portal. At the bottom, the blockchain infrastructure hosts our smart contract storing 32-byte term roots for public verification."

---

## SLIDE 9: How It Works - Issuance Flow

**Time: 1 minute**

### Title: Credential Issuance Process

### Content:

```
Step 1: Export course completion data from LMS
         ↓
Step 2: Construct Verkle tree for the term
        (2,500 courses processed in ~250ms)
         ↓
Step 3: Compute 32-byte root commitment
         ↓
Step 4: Publish root to blockchain
        (Cost: ~$2 per term)
         ↓
Step 5: Generate student receipts with proofs
```

**Key Design Decision:**

- Single 32-byte root covers all students in a term
- Each student receives individual cryptographic proofs

### Visual Suggestions:

- Sequential flowchart with numbered steps
- Data flow: LMS → JSON → Tree → Blockchain → Receipts
- Cost efficiency highlight

### Speaker Notes:

> "The issuance process begins with exporting grade data from the LMS. A Verkle tree is constructed for each academic term, processing 2,500 courses in approximately 250 milliseconds. The key design decision is that a single 32-byte root commitment covers all students in a term, while each student receives individual proofs. The blockchain cost is approximately $2 per term, regardless of student count."

---

## SLIDE 10: How It Works - Verification Flow

**Time: 1 minute**

### Title: Credential Verification Process

### Content:

**Two-Phase Verification:**

**Phase 1: Blockchain Validation**

- Query smart contract for published term root
- Confirm receipt root matches on-chain record
- Validate timestamp for temporal integrity

**Phase 2: Cryptographic Verification**

- Recompute hashes from course data
- Verify IPA proof path to root commitment
- Confirm cryptographic binding of credentials

**Result:** Trustless verification without institutional contact

### Visual Suggestions:

- Two-phase verification diagram
- Checkmarks for each verification step
- Trustless verification emphasis

### Speaker Notes:

> "Verification operates in two phases. First, the blockchain is queried to confirm the term root exists and matches the receipt. Second, cryptographic verification validates each course proof using IPA verification. This enables trustless verification in 3-5 milliseconds per course, without requiring any contact with the issuing institution."

---

## SLIDE 11: Demo - Issuer Dashboard

**Time: 1.5 minutes**

### Title: System Demonstration: Issuer Dashboard

### Content:

**Demonstration Scope:**

1. Term management interface
2. Root publication to blockchain
3. Transaction confirmation
4. Revocation management

### Visual Suggestions:

- Use: `issuer_login.png`
- Use: `issuer_publish_a_term.png`
- Use: `issuer_published_result.png`
- Use: `issuer_revocation_page.png`

### Speaker Notes:

> "I will now demonstrate the issuer dashboard. This interface allows the university to manage credentials. When a term is finalized, the administrator publishes the root commitment to Ethereum with a single action. Transaction confirmation occurs within seconds. The system also supports credential revocation through a version-based supersession mechanism."

---

## SLIDE 12: Demo - Student/Verifier Portal

**Time: 1.5 minutes**

### Title: System Demonstration: Verification Portal

### Content:

**Demonstration Scope:**

1. Receipt upload interface
2. Cryptographic verification process
3. Verification results with credential details
4. **Provenance timeline** - proof of issuance history

### Visual Suggestions:

- Use: `verifier_verify.png`
- Use: `verifier_verify_result.png`
- Use: `verifier_verify_detailed_result_timeline.png`

### Speaker Notes:

> "The verification portal allows verifiers to upload a student's receipt. The system performs cryptographic verification of each credential. The results display verified courses with their details. Critically, the provenance timeline demonstrates when each credential was issued - this prevents backdating and establishes a tamper-proof academic history."

---

## SLIDE 13: Demo - Selective Disclosure

**Time: 1 minute**

### Title: System Demonstration: Selective Disclosure

### Content:

**Privacy-Preserving Verification:**

- Student possesses complete receipt (35 courses)
- Selects only relevant credentials to share (e.g., 3 programming courses)
- Removes unwanted courses from receipt
- Remaining proofs maintain cryptographic validity

**Technical Basis:**

- Each course has an independent proof
- Proofs are not interdependent
- Enabled by Verkle tree mathematical properties

### Visual Suggestions:

- Use: `verifer_selective disclosure_computer_courses.png`
- Use: `verifier_download_selective_receipt.png`
- Use: `verifier_verify_selective_receipt.png`

### Speaker Notes:

> "Selective disclosure enables privacy-preserving verification. A student with 35 courses can share only the relevant credentials - for example, 3 programming courses for a software position. The proofs remain valid because each course is independently provable. This preserves student privacy while maintaining provenance integrity."

---

## SLIDE 14: Demo - Credential Revocation

**Time: 1 minute**

### Title: System Demonstration: Credential Revocation

### Content:

**Revocation Scenario:**

1. Academic misconduct detected for specific course
2. Revocation request submitted through issuer dashboard
3. Administrator reviews and approves request
4. System rebuilds term tree excluding revoked credential
5. New root (version 2) published to blockchain
6. Old receipts show "superseded by newer version" warning

**Version-Based Supersession:**

- Original proofs remain cryptographically valid
- Blockchain tracks both v1 and v2 roots
- Verifiers see version currency status

### Visual Suggestions:

- Use: `issuer_revocation_page.png`
- Version comparison diagram
- Blockchain version tracking visualization

### Speaker Notes:

> "The system supports credential revocation through version-based supersession. When a course must be revoked due to academic misconduct, the issuer submits a revocation request. After approval, the system rebuilds the term tree without the revoked credential and publishes a new root as version 2. The blockchain maintains both versions, and verifiers are alerted when receipts reference superseded versions."

---

## SLIDE 15: Results - Performance

**Time: 1.5 minutes**

### Title: Performance Evaluation

### Content:

**Proof Size Analysis:**
| Metric | Target | Achieved |
|--------|--------|----------|
| Single course proof (IPA + state diff) | <1 KB | **~1.8 KB (constant)** |
| Complete receipt (21 courses) | - | **~11 KB** |
| On-chain storage per term | - | **32 bytes** |

**Processing Speed:**
| Operation | Measurement |
|-----------|-------------|
| Tree construction (22,500 courses) | 1.88 seconds |
| Proof generation (per course) | ~70 ms |
| Proof verification (per course) | **~5 ms** |
| Full receipt verification (21 courses) | **<100 ms** |

**Cost Analysis (Measured on Sepolia):**

- Gas per term: **204,359 gas** (~$12.26 at 20 gwei, $3000/ETH)
- Cost per student (5,000 scale): **$0.017**
- Traditional verification: $50-100 per request

### Visual Suggestions:

- Proof size consistency chart across scales
- Performance benchmark visualization
- Cost comparison with traditional methods

### Speaker Notes:

> "Performance evaluation shows proof sizes of 1.8 KB - slightly above target but critically constant across all scales tested (100 to 50,000 students). Verification completes in approximately 5 milliseconds per course, exceeding our 100ms target by 20x. At 5,000 students, the cost is $0.017 per student for a complete degree - compared to $50-100 for a single traditional verification."

---

## SLIDE 16: Comparison with Existing Systems

**Time: 1 minute**

### Title: Comparative Analysis with Related Work

### Content:

| Feature                  | BlockCerts        | IU-TransCert     | **IU-MiCert**          |
| ------------------------ | ----------------- | ---------------- | ---------------------- |
| Verification Granularity | Whole certificate | Whole transcript | **Per course**         |
| Proof Size               | ~2 KB             | ~1 KB (O(log n)) | **~1.8 KB (constant)** |
| Selective Disclosure     | Not supported     | Not supported    | **Supported**          |
| Academic Provenance      | Limited           | Not supported    | **Full support**       |
| On-chain Storage         | Full hash         | Full hash        | **Root only (32B)**    |
| Revocation Mechanism     | Separate registry | None             | **Version-based**      |

### Visual Suggestions:

- Feature comparison matrix
- Visual indicators for IU-MiCert advantages
- Related system references

### Speaker Notes:

> "Comparing with existing systems: BlockCerts and IU-TransCert are limited to whole-credential verification. IU-MiCert is the first system to enable per-course verification with selective disclosure. Notably, IU-MiCert provides full academic provenance support - verifiable proof of when each credential was issued - which existing systems lack."

---

## SLIDE 17: Research Contributions

**Time: 1 minute**

### Title: Research Contributions

### Content:

**Six Novel Contributions:**

1. **Micro-Credential Architecture**
   First practical implementation treating courses as independently verifiable units

2. **Verkle Trees for Academic Credentials**
   Novel application of Ethereum's state proof technology to credential domain

3. **Temporal Integrity Through Blockchain Anchoring**
   Term-based trees with blockchain timestamps prevent credential backdating

4. **Cryptographic Selective Disclosure**
   Students can remove courses from receipts while maintaining proof validity

5. **Version-Based Revocation Mechanism**
   Elegant tree supersession without complex revocation registries

6. **Complete System Implementation**
   CLI tools, REST API, smart contracts, and web interfaces for all stakeholders

**Deployed System:** Sepolia testnet | iu-micert.vercel.app | iumicert-issuer.vercel.app

### Visual Suggestions:

- Six contributions with icons
- QR codes to deployed applications
- System architecture summary

### Speaker Notes:

> "This research presents six contributions. First, a micro-credential architecture for per-course verification. Second, novel application of Verkle trees to credentials. Third, temporal integrity through blockchain anchoring. Fourth, cryptographic selective disclosure. Fifth, version-based revocation. Sixth, a complete working system with CLI, API, smart contracts, and web interfaces - all deployed and functional."

---

## SLIDE 18: Conclusion & Future Work

**Time: 1 minute**

### Title: Conclusion & Future Directions

### Content:

**Summary:**
IU-MiCert addresses six critical gaps in academic credential verification:

- Granular micro-credential verification (per course)
- Constant-size proofs (~1.8 KB regardless of scale)
- Verifiable academic provenance (blockchain timestamps)
- Privacy-preserving selective disclosure
- Version-based revocation mechanism
- Cost-effective deployment ($0.017 per student at 5K scale)

**Future Research Directions:**

- Mobile application development
- Multi-institution federation support
- LMS integration capabilities
- Post-quantum cryptography migration

**Thank You**

### Visual Suggestions:

- Summary with achievement indicators
- Future work roadmap
- Acknowledgments

### Speaker Notes:

> "In conclusion, IU-MiCert demonstrates that Verkle tree technology can effectively address the limitations of existing credential systems. The system achieves granular verification, verifiable academic provenance, and privacy-preserving disclosure while remaining cost-effective. Future work will explore mobile applications, multi-institution support, and post-quantum migration. Thank you for your attention. I am ready to address any questions."

---

# Q&A BACKUP SLIDES

---

## Q&A SLIDE 1: IPA vs KZG Deep Dive

**For question: "Why IPA instead of KZG commitments?"**

### Content:

| Factor               | KZG                       | IPA (Our Choice)                |
| -------------------- | ------------------------- | ------------------------------- |
| Trusted Setup        | **Required** (ceremony)   | **Not needed**                  |
| If setup compromised | All proofs forgeable      | No such risk                    |
| Proof Size           | O(1) - 48 bytes           | O(log n) - ~1.8 KB              |
| Verification Speed   | Fast (pairing operations) | Slower (scalar multiplications) |
| Security Basis       | Pairing assumptions       | Discrete log (Bandersnatch)     |

**Decision Rationale (from thesis):**
Academic credentials must remain valid for **50+ years**, potentially outliving institutions and administrative structures. A trusted setup ceremony creates a long-term security vulnerability: if setup parameters are ever compromised, all historical credentials become forgeable. IPA eliminates this risk entirely.

---

## Q&A SLIDE 2: Revocation Mechanism

**For question: "How does revocation work?"**

### Content:

**Version-Based Supersession:**

```
Original State (v1):
├─ Student A: Course1, Course2
├─ Student B: Course1, Course3  ← To be revoked
└─ Student C: Course1

After Revocation (v2):
├─ Student A: Course1, Course2
├─ Student B: Course3  ← Course1 removed
└─ Student C: Course1

Blockchain:
  termRoots["Semester_1_2023"][1] = 0xabc... (original)
  termRoots["Semester_1_2023"][2] = 0xdef... (new)
  currentVersion = 2
```

**Key Point:** Old proofs remain cryptographically valid but verifiers see "superseded by newer version" warning.

---

## Q&A SLIDE 3: Security Guarantees

**For question: "What prevents forgery/tampering?"**

### Content:

**Threat Model & Mitigations:**

| Threat              | Mitigation                             |
| ------------------- | -------------------------------------- |
| Forge credential    | IPA proof must bind to blockchain root |
| Modify grade        | Value hash mismatch detected           |
| Backdate credential | Blockchain timestamp immutable         |
| Replay attack       | Version tracking                       |
| Key compromise      | Only issuer can publish                |

**Cryptographic Assumptions:**

- Discrete log hardness (Bandersnatch curve)
- SHA-256 collision resistance
- Ethereum consensus security

---

## Q&A SLIDE 4: Gas Costs Detail

**For question: "What are the blockchain costs?"**

### Content:

**Measured On-Chain Operations (Sepolia Testnet):**

| Operation          | Gas (Measured) | Cost @ 20 gwei | Cost in USD ($3000/ETH) |
| ------------------ | -------------- | -------------- | ----------------------- |
| publishTermRoot    | **204,359**    | 0.00408 ETH    | **~$12.26**             |
| getTermRoot (read) | 0              | Free           | $0                      |

**Complete Degree Cost (7 terms):**

- Total: 7 × $12.26 = **$85.82** total blockchain cost
- Per student (5,000 scale): **$0.017**
- Per student (50,000 scale): **$0.0017**

**Cost Independence Validated:**
Gas variance of ±0.03% across 100 to 50,000 students - cost is independent of student count.

**Comparison:** Traditional verification $50-100 per request → Break-even after 0.017% of one verification

---

## Q&A SLIDE 5: Scalability Analysis

**For question: "How does this scale?"**

### Content:

**Multi-Scale Performance (Measured):**

| Students   | Courses     | Tree Build | ms/course | Receipt Size |
| ---------- | ----------- | ---------- | --------- | ------------ |
| 100        | 450         | 0.57s      | 1.27      | 11.2 KB      |
| 1,000      | 4,500       | 0.81s      | 0.18      | 11.3 KB      |
| 5,000      | 22,500      | 1.88s      | 0.08      | ~8 KB        |
| 10,000     | 45,000      | 2.73s      | 0.06      | ~8 KB        |
| **50,000** | **225,000** | **12.11s** | **0.05**  | **~8 KB**    |

**Key Findings:**

- Receipt size **constant** regardless of total students (O(1) validated)
- Per-course construction time **improves** at larger scale (1.27ms → 0.05ms)
- Tree has 32 levels with 256-ary width per level

---

## Q&A SLIDE 6: Privacy Considerations

**For question: "How is student privacy protected?"**

### Content:

**Privacy Features:**

1. **Selective Disclosure**

   - Students control what to share
   - Remove unwanted courses, proofs still valid

2. **No On-Chain Personal Data**

   - Only 32-byte roots stored on blockchain
   - Course data stays in student's receipt

3. **GDPR Considerations**
   - Students own their receipts
   - No personal data on public blockchain
   - "Right to be forgotten" via selective sharing

**Not Implemented (Future Work):**

- Zero-knowledge proofs for grade ranges
- Attribute-based credentials

---

## Q&A SLIDE 7: Comparison with ZK Solutions

**For question: "Why not use ZK proofs?"**

### Content:

**ZK-SNARKs vs Our Approach:**

| Aspect           | ZK-SNARKs        | IU-MiCert (Verkle+IPA)            |
| ---------------- | ---------------- | --------------------------------- |
| Proof generation | Slow (seconds)   | Fast (ms)                         |
| Trusted setup    | Usually required | **Not needed**                    |
| Complexity       | Very high        | Moderate                          |
| Maturity         | Newer            | Based on Ethereum production code |
| What's hidden    | Everything       | Course data (not existence)       |

**Our Design Choice:**
Verkle+IPA provides sufficient privacy (selective disclosure) with simpler implementation and no trusted setup. Full ZK is future enhancement.

---

## Q&A SLIDE 8: Real-World Deployment

**For question: "Is this production-ready?"**

### Content:

**Current Deployment Status:**

| Component        | Status    | Location                   |
| ---------------- | --------- | -------------------------- |
| Smart Contract   | Deployed  | Sepolia: `0x2452F0...2f79` |
| Student Portal   | Live      | iu-micert.vercel.app       |
| Issuer Dashboard | Live      | iumicert-issuer.vercel.app |
| CLI Tools        | Complete  | 15+ commands               |
| Test Data        | Generated | 5 students × 6 terms       |

**Production Path:**

1. Security audit
2. Mainnet deployment (or L2 for lower costs)
3. LMS integration
4. Pilot program with IU Vietnam

---

## Q&A SLIDE 9: Technical Architecture Details

**For question: "Can you explain the data model?"**

### Content:

**Course Key Generation:**

```
courseKey = "studentDID:termID:courseID"
         = "did:example:ITITIU00001:Semester_1_2023:IT013IU"

keyHash = SHA256(courseKey)  → 32 bytes
        = [stem: 31 bytes][suffix: 1 byte]

valueHash = SHA256(JSON(courseRecord))  → 32 bytes
```

**Tree Structure:**

- One tree per academic term
- 256-ary (256 children per node)
- Depth: ~3-4 levels for typical university

**Receipt JSON Structure:**

```json
{
  "student_id": "ITITIU00001",
  "terms": [{
    "term_id": "Semester_1_2023",
    "verkle_root": "0x1a72...",
    "version": 1,
    "courses": [...],
    "proofs": {...}
  }]
}
```

---

## Q&A SLIDE 10: Limitations Acknowledgment

**For question: "What are the limitations?"**

### Content:

**Known Limitations:**

| Limitation               | Impact                  | Mitigation/Future Work        |
| ------------------------ | ----------------------- | ----------------------------- |
| Not quantum-resistant    | Long-term concern       | Monitor post-quantum research |
| Requires Ethereum access | Verifiers need internet | Caching, light clients        |
| Single issuer design     | Centralized trust       | Multi-issuer federation       |
| Manual data pipeline     | No real-time LMS        | LMS plugins planned           |
| JSON receipt handling    | Technical users         | Mobile app with UI            |

**Honest Assessment:**
System is thesis-scope prototype. Production deployment needs additional security audit and infrastructure.

---

# PRESENTATION TIPS

## Timing Guide

- Practice to fit 14-15 minutes
- Demo can be live OR pre-recorded video
- Keep 1 minute buffer for transitions

## Key Points to Emphasize

1. **Novelty:** First per-course verification system with provenance
2. **Practical:** Fully deployed and functional on Sepolia
3. **Efficient:** $0.017/student (5K scale) vs $50-100 traditional
4. **Privacy:** Cryptographic selective disclosure
5. **Scalable:** Constant proof size validated from 100 to 50,000 students

## Potential Difficult Questions

1. "Why not just use a database?" → Trustless verification, no institution contact needed
2. "What if Ethereum fails?" → Contract on immutable blockchain, could redeploy elsewhere
3. "Is this better than W3C Verifiable Credentials?" → Complementary, we use VC-compatible receipts

## Demo Backup Plan

If live demo fails:

- Have screenshots ready
- Have pre-recorded video
- Walk through receipt JSON structure

---

**Document Created:** December 2025
**For:** IU-MiCert Thesis Defense
