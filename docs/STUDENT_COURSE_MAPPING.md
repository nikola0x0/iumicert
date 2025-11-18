# Student & Course Data Mapping in Verkle Tree

## How Student IDs and Course Data Are Encoded

### 1. Key Construction

In the `AddCourses` function in `term_aggregation.go`:

```go
// Course key construction
courseKey := fmt.Sprintf("%s:%s:%s", studentDID, tvt.TermID, course.CourseID)
courseKeyHash := sha256.Sum256([]byte(courseKey))

// Example:
// studentDID = "did:example:ITITIU00001"
// termID = "Semester_1_2023" 
// courseID = "IT001IU"
// courseKey = "did:example:ITITIU00001:Semester_1_2023:IT001IU"
// courseKeyHash = SHA256("did:example:ITITIU00001:Semester_1_2023:IT001IU")
```

### 2. Value Construction

```go
// Course data serialization
courseData, err := json.Marshal(course)  // CourseCompletion struct
courseValueHash := sha256.Sum256(courseData)  // Hash of the serialized data
```

### 3. Tree Insertion

```go
// Insert into Verkle tree: key = H(student:term:course), value = H(course_data)
err = tvt.tree.Insert(courseKeyHash[:], courseValueHash[:], nil)
```

## Verkle Tree Structure with Student/Course Data

```
Verkle Tree for Term "Semester_1_2023":

┌─────────────────────────────────────────────────────────────┐
│ Root: H(commitment of entire tree)                         │
└─────────────────┬───────────────────────────────────────────┘
                  │
        ┌─────────┴─────────┐
        │  Stem: H(student:term) (first 31 bytes of key)     │
        └─────────┬─────────┘
                  │
      ┌───────────┼───────────┐
      │           │           │
  Suffix:0    Suffix:1    Suffix:2    ...
  (course1)   (course2)   (course3)   ...
      │           │           │
"H(IT001IU)" "H(IT153IU)" "H(PE008IU)"
      │           │           │
"H(course_data1)" "H(course_data2)" "H(course_data3)"
```

## Specific Example

For Student "ITITIU00001" in "Semester_1_2023":

```
Student: "did:example:ITITIU00001"
Term: "Semester_1_2023" 
Courses: ["IT001IU", "IT153IU", "PE008IU", "MA001IU"]

Key-Value Pairs:
1. Key: "did:example:ITITIU00001:Semester_1_2023:IT001IU" 
   → Hash: [0x12, 0x34, 0x56, ..., 0xAB] (32 bytes)
   → Stem: [0x12, 0x34, ..., 0xAB] (first 31 bytes)
   → Suffix: 0xAB (last byte)
   → Value: H(serialized IT001IU course data)

2. Key: "did:example:ITITIU00001:Semester_1_2023:IT153IU"
   → Hash: [0x23, 0x45, 0x67, ..., 0xBC] (32 bytes)
   → Stem: [0x23, 0x45, ..., 0xBC] (first 31 bytes) 
   → Suffix: 0xBC (last byte)
   → Value: H(serialized IT153IU course data)

3. Key: "did:example:ITITIU00001:Semester_1_2023:PE008IU"
   → Hash: [0x34, 0x56, 0x78, ..., 0xCD] (32 bytes)
   → Stem: [0x34, 0x56, ..., 0xCD] (first 31 bytes)
   → Suffix: 0xCD (last byte)
   → Value: H(serialized PE008IU course data)
```

## How Keys Are Used in Proofs

### Proof Generation:
1. **Course Key**: Used to generate proof for specific course
2. **Course Key Hash**: Used as input to `MakeVerkleMultiProof()`
3. **Student ID + Term**: Implicitly encoded in the key structure

### Proof Verification:
1. **Verification Receipt**: Contains StudentDID, TermID, RevealedCourses
2. **Course Key Reconstruction**: 
   ```go
   expectedKey := fmt.Sprintf("%s:%s:%s", receipt.StudentDID, receipt.TermID, course.CourseID)
   ```
3. **Hash Comparison**: 
   - Reconstruct key hash from expected course data
   - Verify it matches what's in the StateDiff
   - Verify the value hash matches the course data

## StateDiff Structure with Student/Course Data

```
StateDiff Example:
[
  {
    Stem: [0x12, 0x34, ..., 0xAB],  // H("did:example:ITITIU00001:Semester_1_2023")[0:31]
    SuffixDiffs: [
      {
        Suffix: 0x00,  // Represents "IT001IU" (course position in polynomial)
        CurrentValue: [0x78, 0x9A, ...]  // H(serialized course data)
      }
    ]
  }
]
```

## Verification Process with Student Data

```
1. Student presents: {StudentID: "ITITIU00001", Term: "Semester_1_2023", Course: "IT001IU"}

2. Receipt contains: {VerkleProof, StateDiff, CourseProofs}

3. Verification:
   a) Reconstruct expected key: "ITITIU00001:Semester_1_2023:IT001IU"
   b) Hash the key: H("ITITIU00001:Semester_1_2023:IT001IU")
   c) Extract stem (first 31 bytes) and suffix (last byte)
   d) Find matching stem in StateDiff
   e) Find matching suffix in SuffixDiffs
   f) Verify the CurrentValue matches H(serialized course data)
   g) Use VerkleProof to verify this key-value exists in tree with known root
```

## Security Features

1. **Student Isolation**: Different student stems are at different tree locations
2. **Term Isolation**: Different terms are separate trees (or separate sections)
3. **Course Identification**: Each course within a student-term is at a specific suffix position
4. **Data Integrity**: Course values are hashed, preventing content tampering

## Selective Disclosure

Students can hide certain courses by omitting them from the receipt:

```
Full Student Data: [CourseA, CourseB, CourseC, CourseD]
Selective Receipt: [CourseA, CourseC]  // Only reveal these courses

The hidden courses still exist in the Verkle tree (and blockchain root),
but are not included in the StateDiff, so aren't revealed to the verifier.
```

This structure ensures that:
- Each student's data is uniquely addressable
- Course data integrity is maintained through hashing
- Selective disclosure is possible without affecting proof validity
- Verification can tie course data back to specific students and terms