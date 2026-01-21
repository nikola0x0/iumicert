# Inner Product Arguments (IPAs)

> A summary of Dankrad Feist's article on Inner Product Arguments, a core component of zero-knowledge proof systems like Bulletproofs.

---

## Table of Contents

1. [Introduction](#1-introduction)
2. [Pedersen Commitments](#2-pedersen-commitments)
3. [The IPA Protocol](#3-the-ipa-protocol)
4. [Optimizations and Extensions](#4-optimizations-and-extensions)
5. [Visual Summary](#5-visual-summary)

---

## 1. Introduction

**IPA (Inner Product Argument)** is a zero-knowledge proof mechanism that allows a **Prover** to convince a **Verifier** about the correctness of an **inner product** between two vectors, with logarithmic proof size and computation cost $O(\log n)$.

### 1.1 Inner Product Definition

The inner product of two vectors $\vec{a}$ and $\vec{b}$ is:

$$\vec{a} \cdot \vec{b} = \sum_{i=0}^{n-1} a_i \cdot b_i$$

**Example:**

```
a = [3, 4, 5]
b = [1, 2, 3]

a · b = (3×1) + (4×2) + (5×3) = 3 + 8 + 15 = 26
```

### 1.2 Connection to Polynomial Evaluation

If we set:

- $\vec{b} = (1, z, z^2, \ldots, z^{n-1})$ — powers of evaluation point $z$
- $\vec{a} = (a_0, a_1, \ldots, a_{n-1})$ — coefficients of polynomial $f(X) = \sum a_i X^i$

Then the inner product equals the **polynomial evaluation at $z$**:

$$\vec{a} \cdot \vec{b} = f(z)$$

**This is the key insight**: Proving an inner product is equivalent to proving a polynomial evaluation!

### 1.3 IPA vs KZG Comparison

| Aspect                  | Pedersen + IPA (Bulletproofs)                                                | KZG Commitments                                         |
| ----------------------- | ---------------------------------------------------------------------------- | ------------------------------------------------------- |
| **Security Assumption** | Discrete Logarithm                                                           | Bilinear Groups (Pairings)                              |
| **Trusted Setup**       | ✅ **NOT required**                                                          | ❌ **Required**                                         |
| **Proof Size**          | $O(\log n)$ group elements                                                   | $O(1)$ — single group element                           |
| **Verification Cost**   | $O(n)$ group operations                                                      | $O(1)$ — single pairing                                 |
| **Best Use Case**       | When no trusted setup is acceptable, or for proof aggregation (amortization) | When succinct proofs and fast verification are critical |

---

## 2. Pedersen Commitments

IPA is built on **Pedersen Commitments**, which use elliptic curve points.

### 2.1 Basic Commitment

To commit to scalars $(a_0, a_1, \ldots, a_{n-1})$ using generators $(G_0, G_1, \ldots, G_{n-1})$:

$$C = a_0 \cdot G_0 + a_1 \cdot G_1 + \ldots + a_{n-1} \cdot G_{n-1} = \sum_{i=0}^{n-1} a_i \cdot G_i$$

**In code (conceptually):**

```go
func PedersenCommit(values []Scalar, generators []Point) Point {
    C := IdentityPoint
    for i := 0; i < len(values); i++ {
        C = C.Add(generators[i].ScalarMul(values[i]))
    }
    return C
}
```

### 2.2 Key Properties

| Property                   | Description                                                                                                                                                                             |
| -------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Binding**                | Cannot find different values $(b_0, b_1, \ldots) \neq (a_0, a_1, \ldots)$ that produce the same commitment $C$ (unless you can solve Discrete Log, which is computationally infeasible) |
| **Hiding**                 | The commitment $C$ reveals nothing about the values (information-theoretically hiding if using a blinding factor)                                                                       |
| **Additively Homomorphic** | $C_1 + C_2$ is a valid commitment to the element-wise sum of the original vectors                                                                                                       |

### 2.3 Homomorphic Property Example

```
Commit(a) = a₀·G₀ + a₁·G₁ = C_a
Commit(b) = b₀·G₀ + b₁·G₁ = C_b

C_a + C_b = (a₀+b₀)·G₀ + (a₁+b₁)·G₁ = Commit(a + b)
```

This property is crucial for the IPA folding technique.

---

## 3. The IPA Protocol

The core strategy is **"Divide and Conquer"** — reduce the problem size by half in each of $\log n$ rounds.

### 3.1 The Statement to Prove

The Prover wants to prove that commitment $C$ satisfies the **Inner Product Property**:

$$C = \langle\vec{a}, \vec{G}\rangle + \langle\vec{b}, \vec{H}\rangle + (\vec{a} \cdot \vec{b}) \cdot Q$$

Where:

- $\vec{a}, \vec{b}$ — secret vectors (Prover knows these)
- $\vec{G}, \vec{H}$ — public generator vectors
- $Q$ — public generator for the inner product term
- $C$ — the commitment (public)

### 3.2 The Reduction Step (Folding)

In each round, vectors $\vec{a}, \vec{b}, \vec{G}, \vec{H}$ are split into left (L) and right (R) halves.

```
┌─────────────────────────────────────────────────────────────────┐
│                     SINGLE IPA ROUND                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Input: vectors of size n                                        │
│         a = [a_L | a_R],  b = [b_L | b_R]                       │
│         G = [G_L | G_R],  H = [H_L | H_R]                       │
│                                                                  │
│  Step 1: Prover computes and sends                              │
│  ────────────────────────────────                               │
│    Cross products:                                               │
│      z_L = ⟨a_R, b_L⟩  (right a with left b)                    │
│      z_R = ⟨a_L, b_R⟩  (left a with right b)                    │
│                                                                  │
│    Auxiliary commitments:                                        │
│      L = ⟨a_R, G_L⟩ + ⟨b_L, H_R⟩ + z_L·Q                       │
│      R = ⟨a_L, G_R⟩ + ⟨b_R, H_L⟩ + z_R·Q                       │
│                                                                  │
│  Step 2: Verifier sends random challenge x                      │
│  ─────────────────────────────────────────                      │
│    x ← random (or Fiat-Shamir hash of transcript)               │
│                                                                  │
│  Step 3: Both parties compute folded values                     │
│  ──────────────────────────────────────────                     │
│    a' = a_L + x·a_R           (size n/2)                        │
│    b' = b_L + x⁻¹·b_R         (size n/2)                        │
│    G' = G_L + x⁻¹·G_R         (size n/2)                        │
│    H' = H_L + x·H_R           (size n/2)                        │
│    C' = x·L + C + x⁻¹·R       (new commitment)                  │
│                                                                  │
│  Output: vectors of size n/2                                     │
│          Continue with (a', b', G', H', C')                     │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 3.3 Why Folding Works

The magic: if $C'$ satisfies the inner product property with $(a', b', G', H')$, then the original $C$ must have satisfied it with $(a, b, G, H)$.

**Algebraically:**

```
C' = x·L + C + x⁻¹·R

Expanding and simplifying shows:
C' = ⟨a', G'⟩ + ⟨b', H'⟩ + (a'·b')·Q

Where: a'·b' = a_L·b_L + x·(a_R·b_L) + x⁻¹·(a_L·b_R) + (a_R·b_R)
             = a_L·b_L + a_R·b_R + x·z_L + x⁻¹·z_R
```

### 3.4 Final Step

After $\log_2 n$ rounds, vectors have size 1. The Prover reveals final scalars $a$ and $b$.

**Verifier checks directly:**
$$C_{\text{final}} = a \cdot G_{\text{final}} + b \cdot H_{\text{final}} + (a \cdot b) \cdot Q$$

### 3.5 Complete Protocol Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                    COMPLETE IPA PROTOCOL                         │
│                      (n = 256 values)                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Round 1: n=256 → n=128                                         │
│    P → V: L₁, R₁                                                │
│    V → P: x₁ (challenge)                                        │
│                                                                  │
│  Round 2: n=128 → n=64                                          │
│    P → V: L₂, R₂                                                │
│    V → P: x₂ (challenge)                                        │
│                                                                  │
│  Round 3: n=64 → n=32                                           │
│    P → V: L₃, R₃                                                │
│    V → P: x₃ (challenge)                                        │
│                                                                  │
│  Round 4: n=32 → n=16                                           │
│    P → V: L₄, R₄                                                │
│    V → P: x₄ (challenge)                                        │
│                                                                  │
│  Round 5: n=16 → n=8                                            │
│    P → V: L₅, R₅                                                │
│    V → P: x₅ (challenge)                                        │
│                                                                  │
│  Round 6: n=8 → n=4                                             │
│    P → V: L₆, R₆                                                │
│    V → P: x₆ (challenge)                                        │
│                                                                  │
│  Round 7: n=4 → n=2                                             │
│    P → V: L₇, R₇                                                │
│    V → P: x₇ (challenge)                                        │
│                                                                  │
│  Round 8: n=2 → n=1                                             │
│    P → V: L₈, R₈                                                │
│    V → P: x₈ (challenge)                                        │
│                                                                  │
│  Final: P reveals (a, b)                                        │
│         V checks: C_final = a·G_final + b·H_final + (ab)·Q      │
│                                                                  │
│  Proof size: 8×(L,R) + (a,b) = 16 points + 2 scalars           │
│            = 16×32 + 2×32 = 576 bytes                          │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 3.6 Soundness (Security)

If a cheating Prover tries to prove a false statement:

- They must find values that satisfy the folded equation for a **random** challenge $x$
- This requires solving a quadratic equation in $x$
- Since $x$ is random (or derived via Fiat-Shamir), the probability of success is negligible (~$2/p$ where $p$ is the field size)

---

## 4. Optimizations and Extensions

### 4.1 Deferred Generator Computation

The Verifier can delay computing $G'$ and $H'$ until the final step by tracking all challenges $(x_1, x_2, \ldots, x_k)$.

**Final generators via Multi-Scalar Multiplication (MSM):**

$$G_{\text{final}} = \sum_{i=0}^{n-1} s_i \cdot G_i$$

Where $s_i$ is computed from the challenges:

```
s_i = ∏_{j where bit j of i is 1} x_j × ∏_{j where bit j of i is 0} x_j⁻¹
```

This replaces $k$ rounds of generator folding with a single MSM at the end.

### 4.2 Polynomial Evaluation Optimization

When proving polynomial evaluation $f(z) = \vec{a} \cdot \vec{b}$ where $\vec{b} = (1, z, z^2, \ldots)$:

- **The Verifier already knows $\vec{b}$!**
- Verifier can compute $b_{\text{final}}$ from initial $\vec{b}$ and challenges
- This eliminates $\vec{H}$ generators entirely
- Proof becomes more compact

### 4.3 Evaluation Form Polynomials

Instead of committing to coefficients $\vec{a} = (a_0, a_1, \ldots)$, commit to **evaluations**:

$$\vec{a} = (f(\omega^0), f(\omega^1), \ldots, f(\omega^{n-1}))$$

Where $\omega$ is a primitive $n$-th root of unity.

**Advantage:** Using the **Barycentric formula**, we can compute the inner product vector $\vec{b}$ for any evaluation point $z$, allowing IPA without $O(n \log n)$ FFT conversions.

---

## 5. Visual Summary

### 5.1 Proof Size Breakdown

```
┌─────────────────────────────────────────────────────────────────┐
│                    IPA PROOF COMPONENTS                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  For n = 256 values:                                            │
│                                                                  │
│  ┌──────────────────┬──────────────┬──────────────────────────┐ │
│  │ Component        │ Count        │ Size                     │ │
│  ├──────────────────┼──────────────┼──────────────────────────┤ │
│  │ L points         │ 8            │ 8 × 32 = 256 bytes       │ │
│  │ R points         │ 8            │ 8 × 32 = 256 bytes       │ │
│  │ Final scalar a   │ 1            │ 32 bytes                 │ │
│  │ Final scalar b   │ 1            │ 32 bytes (often omitted) │ │
│  ├──────────────────┼──────────────┼──────────────────────────┤ │
│  │ TOTAL            │              │ ~544-576 bytes           │ │
│  └──────────────────┴──────────────┴──────────────────────────┘ │
│                                                                  │
│  Compare to naive proof: 256 × 32 = 8,192 bytes                 │
│  Compression ratio: ~14x                                         │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 5.2 Complexity Summary

| Operation        | Prover                     | Verifier                        |
| ---------------- | -------------------------- | ------------------------------- |
| Rounds           | $\log_2 n$                 | $\log_2 n$                      |
| Group operations | $O(n)$                     | $O(n)$ (dominated by final MSM) |
| Field operations | $O(n)$                     | $O(n)$                          |
| Communication    | $O(\log n)$ group elements | —                               |

### 5.3 Security Assumptions

```
┌─────────────────────────────────────────────────────────────────┐
│                    SECURITY FOUNDATION                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  IPA security relies on:                                        │
│                                                                  │
│  1. Discrete Logarithm Problem (DLP)                            │
│     Given G and aG, finding a is computationally infeasible     │
│                                                                  │
│  2. Random Oracle Model (for Fiat-Shamir)                       │
│     Hash function behaves as random oracle                      │
│                                                                  │
│  Result:                                                         │
│  - Computational soundness: cheating Prover succeeds with       │
│    negligible probability                                        │
│  - Perfect completeness: honest Prover always convinces         │
│  - (Optional) Zero-knowledge: can be added with blinding        │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## References

1. **Dankrad Feist's IPA Explanation**

   - [dankradfeist.de/ethereum/2021/07/27/inner-product-arguments.html](https://dankradfeist.de/ethereum/2021/07/27/inner-product-arguments.html)

2. **Bulletproofs Paper** (Bünz et al., 2018)

   - [eprint.iacr.org/2017/1066.pdf](https://eprint.iacr.org/2017/1066.pdf)

3. **Halo Paper** (Bowe et al., 2019)

   - Recursive proof composition without trusted setup

4. **go-ipa Implementation**
   - [github.com/crate-crypto/go-ipa](https://github.com/crate-crypto/go-ipa)

---

_This document summarizes the Inner Product Arguments protocol as used in Verkle trees and Bulletproofs-style proof systems._
