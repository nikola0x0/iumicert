# Problem & Motivation

## Current Challenges in Academic Credentials

Traditional academic credential systems face several limitations:

### 1. Whole-Degree Focus

Current blockchain credential solutions primarily verify **entire degrees** rather than individual competencies or courses. This creates gaps in:

- Skill-based hiring needs
- Lifelong learning verification
- Micro-credential recognition

### 2. Lack of Temporal Context

Credentials are treated as **isolated achievements** without preserving:

- The timeline of learning progression
- Dependencies between courses
- Academic milestones and their order

This enables potential **backdating fraud** - claiming achievements occurred earlier than they actually did.

### 3. Scalability Issues

Traditional Merkle tree approaches require:

- Large proof sizes (logarithmic in tree size)
- Full path disclosure for verification
- Inefficient updates for new credentials

### 4. Privacy vs. Verification Trade-off

Students must choose between:

- **Full disclosure**: Share entire transcript to verify one course
- **No verification**: Keep data private but unverifiable

## Why This Matters

In the modern education landscape:

- **Employers** need to verify specific skills, not just degrees
- **Students** pursue lifelong learning with courses from multiple institutions
- **Institutions** issue thousands of micro-credentials annually
- **Verification** must be efficient, private, and tamper-proof

## Our Approach

IU-MiCert addresses these challenges through:

1. **Course-level granularity**: Each course is independently verifiable
2. **Term-based provenance**: Immutable timeline prevents backdating
3. **Verkle tree efficiency**: 32-byte proofs regardless of transcript size
4. **Selective disclosure**: Verify specific courses without revealing full academic history

The next section explores existing solutions and their limitations.
