# Micro-Credentials

## What are Micro-Credentials?

In IU-MiCert, **every course is a micro-credential** - an independently verifiable unit of learning achievement.

Unlike traditional transcripts that bundle all courses together, micro-credentials enable:

- **Granular verification**: Prove specific skills or knowledge
- **Selective sharing**: Show only relevant courses to employers
- **Stackable credentials**: Combine micro-credentials from multiple institutions

## Course as a Credential

Each course credential contains:

```json
{
  "student_id": "did:example:ITITIU00001",
  "course_id": "CS201",
  "course_name": "Data Structures & Algorithms",
  "grade": "A",
  "credits": 3,
  "term": "Semester_1_2023",
  "issued_at": "2023-12-15T10:00:00Z"
}
```

Plus a **32-byte cryptographic proof** linking it to the term's Verkle tree.

## Use Cases

### 1. Job Applications

**Scenario**: Software engineering position requires data structures knowledge

**Traditional**: Submit entire transcript (reveals all courses, all grades)

**With IU-MiCert**: Share only "Data Structures" and "Algorithms" receipts

**Benefit**: Privacy + targeted verification

### 2. Continuing Education

**Scenario**: Graduate program requires prerequisite courses

**Traditional**: Request official transcript (slow, costs money, whole document)

**With IU-MiCert**: Instantly share prerequisite course receipts

**Benefit**: Speed + cost savings

### 3. Skills-Based Hiring

**Scenario**: Employer wants to verify specific technical competencies

**Traditional**: Degree verification doesn't show individual skills

**With IU-MiCert**: Portfolio of verified micro-credentials per skill

**Benefit**: Evidence-based hiring decisions

## Micro-Credential Features

### Independent Verification

Each course receipt is self-contained:

- No need to request full transcript
- Verifiable without contacting institution
- Blockchain-anchored timestamp proves authenticity

### Composability

Students can:

- Combine credentials from multiple institutions
- Build learning portfolios across years
- Demonstrate continuous skill development

### Privacy-Preserving

Selective disclosure means:

- Share only what's relevant
- Hide sensitive information (e.g., poor grades in unrelated courses)
- Control your academic narrative

## Comparison

| Aspect | Traditional Transcript | IU-MiCert Micro-Credentials |
|--------|----------------------|---------------------------|
| Granularity | Whole document | Per-course |
| Sharing | All or nothing | Selective |
| Verification | Contact institution | Instant, cryptographic |
| Privacy | Reveals everything | Minimal disclosure |
| Cost | $10-50 per transcript | Free verification |
| Speed | Days to weeks | Instant |

## Academic Integrity

Micro-credentials don't compromise academic integrity:

- Each course is cryptographically bound to its term
- Timeline provenance prevents backdating
- Blockchain anchoring ensures authenticity
- Verification is independent of issuer

---

**Next**: Understanding academic provenance and timeline integrity.
