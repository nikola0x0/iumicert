# Academic Provenance

## What is Academic Provenance?

**Provenance** = the chronological record of ownership or existence of an object.

For academics, **provenance means the complete, verifiable timeline** of a student's learning journey - when each achievement occurred and in what context.

## The Provenance Problem

Traditional credential systems lack provenance:

- Courses appear as isolated achievements
- No verification of when learning occurred
- Possibility of backdating credentials
- No proof of progression through curriculum

### Real-World Scenario

**Without Provenance**:

A student could claim to have completed "Advanced Machine Learning" before taking "Introduction to Programming" - logically impossible but technically undetectable.

**With IU-MiCert Provenance**:

Each course is timestamped to a specific term, published immutably on blockchain. The timeline cannot be manipulated.

## How IU-MiCert Achieves Provenance

### Term-Based Architecture

```
Semester_1_2023 → Verkle Tree → Blockchain (Jan 2024)
Semester_2_2023 → Verkle Tree → Blockchain (Jun 2024)
Summer_2023     → Verkle Tree → Blockchain (Sep 2024)
```

**Key Insight**: Each term's credentials are published together, creating an immutable timeline.

### Timeline Integrity

Cannot backdate because:

1. **Blockchain timestamp**: When root was published
2. **Term structure**: Courses tied to specific academic periods
3. **Cryptographic binding**: Proofs only valid for correct term

### Verification Flow

To verify a course:

1. Check course is in claimed term's Verkle tree (cryptographic proof)
2. Check term root exists on blockchain (immutable record)
3. Check blockchain timestamp matches claimed term dates
4. ✅ Course provenance verified

## Benefits of Provenance

### For Students

- **Authentic timeline**: Proof of when learning occurred
- **Learning journey**: Demonstrate progression through curriculum
- **Credential integrity**: No possibility of backdating claims

### For Employers

- **Timeline verification**: Confirm sequence of skill acquisition
- **Red flag detection**: Identify impossible timelines
- **Trust**: Blockchain-anchored evidence of achievements

### For Institutions

- **Fraud prevention**: Eliminate backdating attacks
- **Reputation protection**: Verifiable issuance history
- **Compliance**: Audit trail for accreditation

## Comparison: Backdating Attack

### Without Provenance

1. Student claims "Graduated 2022" with fake credential
2. Backdates courses to appear earlier
3. No way to verify actual completion dates
4. ❌ Attack succeeds

### With IU-MiCert

1. Student attempts to backdate course to earlier term
2. Cryptographic proof fails (course not in earlier term's tree)
3. Blockchain shows when term was actually published
4. ✅ Attack prevented

## Academic Integrity

Provenance ensures:

- **Chronological authenticity**: Courses completed when claimed
- **Prerequisite verification**: Proper progression through curriculum
- **Accreditation compliance**: Audit trail for quality assurance
- **Institutional trust**: Verifiable credential issuance history

---

**Next**: Implementation details and technology stack.
