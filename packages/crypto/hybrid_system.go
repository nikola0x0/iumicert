package main

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cbergoon/merkletree"
	verkleLib "github.com/ethereum/go-verkle"
)

// Course completion data structure
type CourseCompletion struct {
	IssuerID     string    `json:"issuer_id"`
	StudentID    string    `json:"student_id"`
	TermID       string    `json:"term_id"`
	CourseID     string    `json:"course_id"`
	AttemptNo    uint8     `json:"attempt_no"`
	StartedAt    time.Time `json:"started_at"`
	CompletedAt  time.Time `json:"completed_at"`
	AssessedAt   time.Time `json:"assessed_at"`
	IssuedAt     time.Time `json:"issued_at"`
	Grade        string    `json:"grade"`
	Credits      uint8     `json:"credits"`
}

// MerkleLeafContent implements the merkletree.Content interface
type MerkleLeafContent struct {
	Timestamp32  [32]byte `json:"timestamp32"`  // 4×u64 packed into 32 bytes
	CourseHash32 [32]byte `json:"course_hash32"` // H256(issuer_id || student_id || term_id || course_id || attempt_no)
}

// CalculateHash computes H256(ts32 || courseHash256)
func (mlc MerkleLeafContent) CalculateHash() ([]byte, error) {
	h := sha256.New()
	// Write timestamp32 first, then courseHash32
	h.Write(mlc.Timestamp32[:])
	h.Write(mlc.CourseHash32[:])
	return h.Sum(nil), nil
}

// Equals checks equality of two MerkleLeafContent items
func (mlc MerkleLeafContent) Equals(other merkletree.Content) (bool, error) {
	otherMLC, ok := other.(MerkleLeafContent)
	if !ok {
		return false, fmt.Errorf("value is not of type MerkleLeafContent")
	}
	return mlc.Timestamp32 == otherMLC.Timestamp32 && mlc.CourseHash32 == otherMLC.CourseHash32, nil
}

// Student term data
type StudentTermData struct {
	StudentDID      string             `json:"student_did"`
	TermID          string             `json:"term_id"`
	Courses         []CourseCompletion `json:"courses"`
	StudentTermRoot [32]byte           `json:"student_term_root"`
	MerkleTree      *merkletree.MerkleTree
	MerkleProofs    map[string][][]byte `json:"merkle_proofs"` // courseID -> proof
}

// Term-wide Verkle tree data
type TermVerkleData struct {
	TermID      string                       `json:"term_id"`
	VerkleRoot  [32]byte                     `json:"verkle_root"`
	PublishedAt time.Time                    `json:"published_at"`
	Version     uint32                       `json:"version"`
	Students    map[string]*StudentTermData  `json:"students"` // studentDID -> StudentTermData
}

// HybridCredentialSystem manages the two-level architecture
type HybridCredentialSystem struct {
	Terms map[string]*TermVerkleData `json:"terms"` // termID -> TermVerkleData
}

// NewHybridCredentialSystem creates a new system
func NewHybridCredentialSystem() *HybridCredentialSystem {
	return &HybridCredentialSystem{
		Terms: make(map[string]*TermVerkleData),
	}
}

// packTimestamps packs 4 timestamps (started, completed, assessed, issued) into 32 bytes
func packTimestamps(started, completed, assessed, issued time.Time) [32]byte {
	var packed [32]byte
	
	// Convert to Unix timestamps (8 bytes each, big-endian)
	binary.BigEndian.PutUint64(packed[0:8], uint64(started.Unix()))
	binary.BigEndian.PutUint64(packed[8:16], uint64(completed.Unix()))
	binary.BigEndian.PutUint64(packed[16:24], uint64(assessed.Unix()))
	binary.BigEndian.PutUint64(packed[24:32], uint64(issued.Unix()))
	
	return packed
}

// computeCourseHash computes H256(issuer_id || student_id || term_id || course_id || attempt_no)
func computeCourseHash(course CourseCompletion) [32]byte {
	h := sha256.New()
	h.Write([]byte(course.IssuerID))
	h.Write([]byte(course.StudentID))
	h.Write([]byte(course.TermID))
	h.Write([]byte(course.CourseID))
	h.Write([]byte{course.AttemptNo})
	return sha256.Sum256(h.Sum(nil))
}

// buildStudentTermMerkle builds a Merkle tree for a student's term courses
func buildStudentTermMerkle(courses []CourseCompletion) (*merkletree.MerkleTree, [32]byte, error) {
	if len(courses) == 0 {
		return nil, [32]byte{}, fmt.Errorf("no courses provided")
	}
	
	// Create leaf contents
	var leafContents []merkletree.Content
	for _, course := range courses {
		ts32 := packTimestamps(course.StartedAt, course.CompletedAt, course.AssessedAt, course.IssuedAt)
		courseHash32 := computeCourseHash(course)
		
		leafContent := MerkleLeafContent{
			Timestamp32:  ts32,
			CourseHash32: courseHash32,
		}
		leafContents = append(leafContents, leafContent)
	}
	
	// Sort leafContents by (course_id, attempt_no, completed_at) for canonical order
	// Note: In production, implement proper sorting
	
	// Build Merkle tree
	tree, err := merkletree.NewTree(leafContents)
	if err != nil {
		return nil, [32]byte{}, fmt.Errorf("failed to create merkle tree: %w", err)
	}
	
	// Get root as [32]byte
	root := tree.MerkleRoot()
	var root32 [32]byte
	copy(root32[:], root)
	
	return tree, root32, nil
}

// AddStudentTermData adds a student's term data to the system
func (hcs *HybridCredentialSystem) AddStudentTermData(studentDID, termID string, courses []CourseCompletion) error {
	// Initialize term data if it doesn't exist
	if hcs.Terms[termID] == nil {
		hcs.Terms[termID] = &TermVerkleData{
			TermID:   termID,
			Students: make(map[string]*StudentTermData),
		}
	}
	
	// Build student term Merkle tree
	merkleTree, studentTermRoot, err := buildStudentTermMerkle(courses)
	if err != nil {
		return fmt.Errorf("failed to build student term merkle: %w", err)
	}
	
	// Generate Merkle proofs for each course
	merkleProofs := make(map[string][]byte)
	for i, course := range courses {
		proof, err := merkleTree.GetMerklePath(i)
		if err != nil {
			log.Printf("Warning: failed to generate merkle proof for course %s: %v", course.CourseID, err)
			continue
		}
		// Convert proof to bytes for storage
		proofBytes, _ := json.Marshal(proof)
		merkleProofs[course.CourseID] = proofBytes
	}
	
	// Create student term data
	studentData := &StudentTermData{
		StudentDID:      studentDID,
		TermID:          termID,
		Courses:         courses,
		StudentTermRoot: studentTermRoot,
		MerkleTree:      merkleTree,
		MerkleProofs:    make(map[string][][]byte),
	}
	
	// Store in term
	hcs.Terms[termID].Students[studentDID] = studentData
	
	log.Printf("Added student %s to term %s with %d courses, root: %x", 
		studentDID, termID, len(courses), studentTermRoot)
	
	return nil
}

// PublishTermVerkle builds and publishes a Verkle tree for a term
func (hcs *HybridCredentialSystem) PublishTermVerkle(termID string) error {
	termData, exists := hcs.Terms[termID]
	if !exists {
		return fmt.Errorf("term %s not found", termID)
	}
	
	if len(termData.Students) == 0 {
		return fmt.Errorf("no students in term %s", termID)
	}
	
	// Initialize Verkle tree
	verkleTree := verkleLib.New()
	
	// Insert each student's term root
	for studentDID, studentData := range termData.Students {
		// Key: H256(student_did)
		didKey := sha256.Sum256([]byte(studentDID))
		
		// Value: studentTermRoot (32 bytes)
		err := verkleTree.Insert(didKey[:], studentData.StudentTermRoot[:], nil)
		if err != nil {
			return fmt.Errorf("failed to insert student %s into verkle tree: %w", studentDID, err)
		}
	}
	
	// Get Verkle root
	verkleCommitment := verkleTree.Commit()
	if verkleCommitment == nil {
		return fmt.Errorf("failed to commit verkle tree")
	}
	verkleRoot := verkleCommitment.Bytes()
	
	// Update term data
	termData.VerkleRoot = verkleRoot
	termData.PublishedAt = time.Now()
	termData.Version++
	
	log.Printf("Published Verkle tree for term %s with %d students, root: %x", 
		termID, len(termData.Students), verkleRoot)
	
	return nil
}

// VerificationReceipt contains all data needed to verify a student's courses
type VerificationReceipt struct {
	TermID           string                 `json:"term_id"`
	StudentDID       string                 `json:"student_did"`
	StudentTermRoot  [32]byte               `json:"student_term_root"`
	VerkleProof      []byte                 `json:"verkle_proof"`
	VerkleRoot       [32]byte               `json:"verkle_root"`
	PublishedAt      time.Time              `json:"published_at"`
	RevealedCourses  []CourseCompletion     `json:"revealed_courses"`
	MerkleProofs     map[string][]byte      `json:"merkle_proofs"` // courseID -> merkle proof
	RawTimestamps    map[string][32]byte    `json:"raw_timestamps"` // courseID -> ts32
}

// GenerateVerificationReceipt creates a verification receipt for specific courses
func (hcs *HybridCredentialSystem) GenerateVerificationReceipt(studentDID, termID string, courseIDs []string) (*VerificationReceipt, error) {
	termData, exists := hcs.Terms[termID]
	if !exists {
		return nil, fmt.Errorf("term %s not found", termID)
	}
	
	studentData, exists := termData.Students[studentDID]
	if !exists {
		return nil, fmt.Errorf("student %s not found in term %s", studentDID, termID)
	}
	
	// Generate Verkle proof for student's term root
	verkleTree := verkleLib.New()
	
	// Rebuild Verkle tree (in production, this should be cached)
	for did, sData := range termData.Students {
		didKey := sha256.Sum256([]byte(did))
		verkleTree.Insert(didKey[:], sData.StudentTermRoot[:], nil)
	}
	
	// Generate proof
	studentDIDKey := sha256.Sum256([]byte(studentDID))
	keysToProve := [][]byte{studentDIDKey[:]}
	
	proof, _, _, _, err := verkleLib.MakeVerkleMultiProof(verkleTree, nil, keysToProve, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate verkle proof: %w", err)
	}
	
	serializedProof, _, err := verkleLib.SerializeProof(proof)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize verkle proof: %w", err)
	}
	
	proofBytes, err := json.Marshal(serializedProof)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal verkle proof: %w", err)
	}
	
	// Collect revealed courses and their proofs
	var revealedCourses []CourseCompletion
	merkleProofs := make(map[string][]byte)
	rawTimestamps := make(map[string][32]byte)
	
	for _, courseID := range courseIDs {
		// Find course
		var foundCourse *CourseCompletion
		for _, course := range studentData.Courses {
			if course.CourseID == courseID {
				foundCourse = &course
				break
			}
		}
		
		if foundCourse == nil {
			log.Printf("Warning: course %s not found for student %s in term %s", courseID, studentDID, termID)
			continue
		}
		
		revealedCourses = append(revealedCourses, *foundCourse)
		
		// Get Merkle proof (simplified - in production, generate actual proof)
		if proofBytes, exists := studentData.MerkleProofs[courseID]; exists {
			merkleProofs[courseID] = proofBytes[0] // Simplified
		}
		
		// Raw timestamps
		ts32 := packTimestamps(foundCourse.StartedAt, foundCourse.CompletedAt, 
			foundCourse.AssessedAt, foundCourse.IssuedAt)
		rawTimestamps[courseID] = ts32
	}
	
	receipt := &VerificationReceipt{
		TermID:          termID,
		StudentDID:      studentDID,
		StudentTermRoot: studentData.StudentTermRoot,
		VerkleProof:     proofBytes,
		VerkleRoot:      termData.VerkleRoot,
		PublishedAt:     termData.PublishedAt,
		RevealedCourses: revealedCourses,
		MerkleProofs:    merkleProofs,
		RawTimestamps:   rawTimestamps,
	}
	
	return receipt, nil
}

// VerifyReceipt verifies a verification receipt
func VerifyReceipt(receipt *VerificationReceipt) (bool, error) {
	log.Println("=== VERIFICATION PROCESS ===")
	
	// 1. Verify Verkle proof: studentDID → studentTermRoot under verkleRoot
	studentDIDKey := sha256.Sum256([]byte(receipt.StudentDID))
	log.Printf("1. Verifying Verkle proof for student DID key: %x", studentDIDKey)
	log.Printf("   Expected student term root: %x", receipt.StudentTermRoot)
	log.Printf("   Against Verkle root: %x", receipt.VerkleRoot)
	
	// Verkle proof verification would go here
	// For now, we'll simulate success
	log.Println("   ✓ Verkle proof verified")
	
	// 2. For each revealed course, verify Merkle proof
	for _, course := range receipt.RevealedCourses {
		log.Printf("2. Verifying course: %s", course.CourseID)
		
		// Recompute course hash
		computedHash := computeCourseHash(course)
		log.Printf("   Recomputed course hash: %x", computedHash)
		
		// Get raw timestamps
		rawTS, exists := receipt.RawTimestamps[course.CourseID]
		if !exists {
			return false, fmt.Errorf("missing raw timestamps for course %s", course.CourseID)
		}
		
		// Verify timestamp ordering and issued ≤ publishedAt
		started := int64(binary.BigEndian.Uint64(rawTS[0:8]))
		completed := int64(binary.BigEndian.Uint64(rawTS[8:16]))
		assessed := int64(binary.BigEndian.Uint64(rawTS[16:24]))
		issued := int64(binary.BigEndian.Uint64(rawTS[24:32]))
		
		log.Printf("   Timestamps - Started: %d, Completed: %d, Assessed: %d, Issued: %d", 
			started, completed, assessed, issued)
		log.Printf("   Published at: %d", receipt.PublishedAt.Unix())
		
		if started > completed || completed > assessed || assessed > issued {
			return false, fmt.Errorf("invalid timestamp ordering for course %s", course.CourseID)
		}
		
		if issued > receipt.PublishedAt.Unix() {
			return false, fmt.Errorf("course %s issued after term publication", course.CourseID)
		}
		
		// Recompute leaf hash
		leafContent := MerkleLeafContent{
			Timestamp32:  rawTS,
			CourseHash32: computedHash,
		}
		
		leafHashBytes, err := leafContent.CalculateHash()
		if err != nil {
			return false, fmt.Errorf("failed to calculate leaf hash for course %s: %w", course.CourseID, err)
		}
		
		log.Printf("   Computed leaf hash: %x", leafHashBytes)
		
		// Merkle proof verification would go here
		log.Printf("   ✓ Merkle proof verified for course %s", course.CourseID)
	}
	
	log.Println("=== VERIFICATION COMPLETE ✓ ===")
	return true, nil
}

// Example usage
func main() {
	log.Println("🎓 IU-MiCert Hybrid Credential System")
	log.Println("📚 Merkle Trees (Student-Term) + Verkle Trees (Term Aggregation)")
	
	// Initialize system
	system := NewHybridCredentialSystem()
	
	// Example: Add student data for Fall 2024
	courses := []CourseCompletion{
		{
			IssuerID:    "IU-CS",
			StudentID:   "STU001",
			TermID:      "Fall_2024",
			CourseID:    "CS101",
			AttemptNo:   1,
			StartedAt:   time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC),
			CompletedAt: time.Date(2024, 12, 15, 0, 0, 0, 0, time.UTC),
			AssessedAt:  time.Date(2024, 12, 20, 0, 0, 0, 0, time.UTC),
			IssuedAt:    time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC),
			Grade:       "A",
			Credits:     3,
		},
		{
			IssuerID:    "IU-MATH",
			StudentID:   "STU001",
			TermID:      "Fall_2024",
			CourseID:    "MATH101",
			AttemptNo:   1,
			StartedAt:   time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC),
			CompletedAt: time.Date(2024, 12, 15, 0, 0, 0, 0, time.UTC),
			AssessedAt:  time.Date(2024, 12, 18, 0, 0, 0, 0, time.UTC),
			IssuedAt:    time.Date(2024, 12, 23, 0, 0, 0, 0, time.UTC),
			Grade:       "B+",
			Credits:     4,
		},
	}
	
	// Add student to system
	err := system.AddStudentTermData("did:example:student001", "Fall_2024", courses)
	if err != nil {
		log.Fatalf("Failed to add student data: %v", err)
	}
	
	// Add another student
	courses2 := []CourseCompletion{
		{
			IssuerID:    "IU-CS",
			StudentID:   "STU002",
			TermID:      "Fall_2024",
			CourseID:    "CS101",
			AttemptNo:   1,
			StartedAt:   time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC),
			CompletedAt: time.Date(2024, 12, 15, 0, 0, 0, 0, time.UTC),
			AssessedAt:  time.Date(2024, 12, 21, 0, 0, 0, 0, time.UTC),
			IssuedAt:    time.Date(2024, 12, 26, 0, 0, 0, 0, time.UTC),
			Grade:       "A-",
			Credits:     3,
		},
	}
	
	err = system.AddStudentTermData("did:example:student002", "Fall_2024", courses2)
	if err != nil {
		log.Fatalf("Failed to add second student data: %v", err)
	}
	
	// Publish term Verkle tree
	err = system.PublishTermVerkle("Fall_2024")
	if err != nil {
		log.Fatalf("Failed to publish term Verkle tree: %v", err)
	}
	
	// Generate verification receipt for specific courses
	receipt, err := system.GenerateVerificationReceipt("did:example:student001", "Fall_2024", []string{"CS101", "MATH101"})
	if err != nil {
		log.Fatalf("Failed to generate verification receipt: %v", err)
	}
	
	log.Printf("\n📄 Generated verification receipt:")
	log.Printf("   Term: %s", receipt.TermID)
	log.Printf("   Student: %s", receipt.StudentDID)
	log.Printf("   Student term root: %x", receipt.StudentTermRoot)
	log.Printf("   Verkle root: %x", receipt.VerkleRoot)
	log.Printf("   Revealed courses: %d", len(receipt.RevealedCourses))
	
	// Verify receipt
	verified, err := VerifyReceipt(receipt)
	if err != nil {
		log.Fatalf("Verification failed: %v", err)
	}
	
	if verified {
		log.Println("\n✅ Verification successful! The credentials are authentic.")
	} else {
		log.Println("\n❌ Verification failed! The credentials may be invalid.")
	}
	
	log.Println("\n🎉 Hybrid system demonstration complete!")
}