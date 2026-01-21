Verkle Tree Implementation Details

---

Tree Structure

| Property     | Value                                          |
| ------------ | ---------------------------------------------- |
| Library      | go-verkle (Ethereum's official implementation) |
| Width        | 256 children per node                          |
| Depth        | 32 levels                                      |
| Key size     | 32 bytes (256 bits)                            |
| Value size   | 32 bytes (SHA-256 hash of course data)         |
| Commitment   | Pedersen commitment (Bandersnatch curve)       |
| Proof system | IPA (Inner Product Argument)                   |

---

Key Generation (Deterministic)

courseKey = "did:example:{studentID}:{termID}:{courseID}"

Example: "did:example:ITITIU00001:Semester_1_2023:IT013IU"

keyHash = SHA256(courseKey) → 32 bytes
= [stem: 31 bytes] + [suffix: 1 byte]

- Stem: First 31 bytes → path through tree
- Suffix: Last 1 byte → position within leaf node (0-255)

---

Value Storage

courseData = {
"course_id": "IT013IU",
"course_name": "Database Systems",
"grade": "A",
"credits": 3
}

valueHash = SHA256(JSON(courseData)) → 32 bytes

Stored in tree: tree.Insert(keyHash, valueHash)

---

One Tree Per Term

Semester_1_2023/
└── verkle_tree.json (serialized tree)
└── verkle_root.json (32-byte root commitment)

Semester_2_2023/
└── verkle_tree.json
└── verkle_root.json

Why per-term?

- Natural academic boundary
- Independent revocation per term
- Manageable tree size
- Clear temporal anchoring

---

Proof Generation

proof, _, _, \_ := verkle.MakeVerkleMultiProof(
tree, // The Verkle tree
keyList, // Keys to prove
keyPresenceMap, // Which keys exist
)

serialized := verkle.SerializeProof(proof)
// → ~1.8 KB regardless of tree size

Proof contains:

- IPA proof (commitments + challenges)
- Path commitments
- Extension status

---

Proof Verification

err := verkle.VerifyCourseProof(
courseKey, // "did:example:ITITIU00001:Sem1:IT013IU"
courseData, // {course_id, grade, credits...}
proofBytes, // Serialized IPA proof
verkleRoot, // 32-byte root from blockchain
)
// Returns nil if valid, error if tampered

Verification steps:

1. Regenerate keyHash from courseKey
2. Regenerate valueHash from courseData
3. Deserialize proof
4. Run IPA verification against root
5. Confirm value commitment matches

---

Why IPA (not KZG)?

| Aspect          | KZG                   | IPA                    |
| --------------- | --------------------- | ---------------------- |
| Trusted setup   | Required (ceremony)   | Not required           |
| Proof size      | 48 bytes              | ~1.8 KB                |
| Verification    | Pairing-based         | Inner product          |
| Ethereum status | Future (Danksharding) | Production (go-verkle) |

Our choice: IPA because:

- No trusted setup = no ceremony = simpler deployment
- Uses Ethereum's battle-tested code
- Acceptable proof size for credentials

---

Pedersen Commitment (How it works)

Commitment = g₁^v₁ · g₂^v₂ · ... · gₙ^vₙ

Where:

- gᵢ = generator points on Bandersnatch curve
- vᵢ = values being committed

Properties:

- Binding: Cannot open to different value
- Hiding: Reveals nothing about value
- Homomorphic: Can combine commitments

---

Bandersnatch Curve

- Type: Elliptic curve (embedded in BLS12-381)
- Security: ~128 bits
- Why chosen: Efficient for IPA, Ethereum-compatible

---

Receipt Structure

{
"student_id": "ITITIU00001",
"term_receipts": {
"Semester_1_2023": {
"verkle_root": "0x1a72f3...",
"receipt": {
"revealed_courses": [...],
"course_proofs": {
"IT013IU": {
"proof": "base64...",
"key_hash": "0xabc...",
"value_hash": "0xdef..."
}
}
}
}
}
}

---

Scalability Numbers

| Students | Courses | Tree Build | Proof Size |
| -------- | ------- | ---------- | ---------- |
| 100      | 450     | 0.57s      | ~1.8 KB    |
| 1,000    | 4,500   | 0.81s      | ~1.8 KB    |
| 5,000    | 22,500  | 1.88s      | ~1.8 KB    |
| 50,000   | 225,000 | 12.11s     | ~1.8 KB    |

Key insight: Proof size stays constant!

---

Code Location

packages/issuer/
├── crypto/
│ └── verkle/
│ ├── tree.go # Tree operations
│ ├── proof.go # Proof generation
│ └── verification.go # IPA verification

---

Need more detail on any specific part?
