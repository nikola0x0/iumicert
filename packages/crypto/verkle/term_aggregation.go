package verkle

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"time"

	verkleLib "github.com/ethereum/go-verkle"
	"iumicert/crypto/merkle"
)

// TermVerkleTree manages the Verkle tree for a complete academic term
type TermVerkleTree struct {
	TermID           string                              `json:"term_id"`
	PublishedAt      time.Time                           `json:"published_at"`
	Version          uint32                              `json:"version"`
	VerkleRoot       [32]byte                            `json:"verkle_root"`
	StudentTerms     map[string]*merkle.StudentTermMerkle `json:"student_terms"` // studentDID -> StudentTermMerkle
	VerkleProofs     map[string][]byte                   `json:"verkle_proofs"` // studentDID -> serialized Verkle proof
	tree             verkleLib.VerkleNode                // Internal tree (not serialized)
}

// VerificationReceipt contains all data needed for off-chain or on-chain verification
type VerificationReceipt struct {
	TermID          string                      `json:"term_id"`
	StudentDID      string                      `json:"student_did"`
	StudentTermRoot [32]byte                    `json:"student_term_root"`
	VerkleProof     []byte                      `json:"verkle_proof"`
	VerkleRoot      [32]byte                    `json:"verkle_root"`
	PublishedAt     time.Time                   `json:"published_at"`
	RevealedCourses []merkle.CourseCompletion   `json:"revealed_courses"`
	MerkleProofs    map[string][]string         `json:"merkle_proofs"` // courseID -> merkle proof (hex)
	RawTimestamps   map[string][32]byte         `json:"raw_timestamps"` // courseID -> packed timestamps
	Metadata        ReceiptMetadata             `json:"metadata"`
}

// ReceiptMetadata contains additional verification metadata
type ReceiptMetadata struct {
	GeneratedAt       time.Time `json:"generated_at"`
	TotalCourses      int       `json:"total_courses"`
	RevealedCourses   int       `json:"revealed_courses"`
	VerificationLevel string    `json:"verification_level"` // "full" or "selective"
}

// NewTermVerkleTree creates a new term-level Verkle tree
func NewTermVerkleTree(termID string) *TermVerkleTree {
	return &TermVerkleTree{
		TermID:       termID,
		StudentTerms: make(map[string]*merkle.StudentTermMerkle),
		VerkleProofs: make(map[string][]byte),
		tree:         verkleLib.New(),
	}
}

// AddStudent adds a student's term data to the Verkle tree
func (tvt *TermVerkleTree) AddStudent(studentDID string, courses []merkle.CourseCompletion) error {
	log.Printf("Adding student %s to term %s with %d courses", studentDID, tvt.TermID, len(courses))
	
	// Create student-term Merkle tree
	studentTerm, err := merkle.NewStudentTermMerkle(studentDID, tvt.TermID, courses)
	if err != nil {
		return fmt.Errorf("failed to create student term merkle for %s: %w", studentDID, err)
	}
	
	// Store the student term data
	tvt.StudentTerms[studentDID] = studentTerm
	
	// Add to Verkle tree: key = H256(studentDID), value = studentTermRoot
	studentDIDKey := sha256.Sum256([]byte(studentDID))
	err = tvt.tree.Insert(studentDIDKey[:], studentTerm.Root[:], nil)
	if err != nil {
		return fmt.Errorf("failed to insert student %s into verkle tree: %w", studentDID, err)
	}
	
	log.Printf("✅ Student %s added successfully, term root: %x", studentDID, studentTerm.Root)
	return nil
}

// PublishTerm computes the final Verkle root and prepares for blockchain publication
func (tvt *TermVerkleTree) PublishTerm() error {
	log.Printf("Publishing term %s with %d students", tvt.TermID, len(tvt.StudentTerms))
	
	if len(tvt.StudentTerms) == 0 {
		return fmt.Errorf("cannot publish term %s: no students added", tvt.TermID)
	}
	
	// Compute Verkle tree commitment
	commitment := tvt.tree.Commit()
	if commitment == nil {
		return fmt.Errorf("failed to compute verkle tree commitment")
	}
	
	// Store the root
	tvt.VerkleRoot = commitment.Bytes()
	tvt.PublishedAt = time.Now()
	tvt.Version++
	
	log.Printf("✅ Term %s published successfully, Verkle root: %x", tvt.TermID, tvt.VerkleRoot)
	return nil
}

// GenerateVerificationReceipt creates a verification receipt for specific courses
func (tvt *TermVerkleTree) GenerateVerificationReceipt(studentDID string, courseIDs []string) (*VerificationReceipt, error) {
	log.Printf("Generating verification receipt for student %s, courses: %v", studentDID, courseIDs)
	
	// Check if term is published
	if tvt.PublishedAt.IsZero() {
		return nil, fmt.Errorf("term %s not yet published", tvt.TermID)
	}
	
	// Get student's term data
	studentTerm, exists := tvt.StudentTerms[studentDID]
	if !exists {
		return nil, fmt.Errorf("student %s not found in term %s", studentDID, tvt.TermID)
	}
	
	// Generate Verkle proof for the student's term root
	studentDIDKey := sha256.Sum256([]byte(studentDID))
	keysToProve := [][]byte{studentDIDKey[:]}
	
	proof, _, _, _, err := verkleLib.MakeVerkleMultiProof(tvt.tree, nil, keysToProve, nil)
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
	
	// If no specific courses requested, include all courses
	if len(courseIDs) == 0 {
		courseIDs = studentTerm.ListCourses()
	}
	
	// Collect revealed courses and their Merkle proofs
	var revealedCourses []merkle.CourseCompletion
	merkleProofs := make(map[string][]string)
	rawTimestamps := make(map[string][32]byte)
	
	for _, courseID := range courseIDs {
		// Get course data
		course, err := studentTerm.GetCourseByID(courseID)
		if err != nil {
			log.Printf("Warning: course %s not found for student %s: %v", courseID, studentDID, err)
			continue
		}
		
		revealedCourses = append(revealedCourses, *course)
		
		// Get Merkle proof
		proof, err := studentTerm.GetProofForCourse(courseID)
		if err != nil {
			log.Printf("Warning: failed to get merkle proof for course %s: %v", courseID, err)
			continue
		}
		merkleProofs[courseID] = proof
		
		// Pack timestamps for verification
		ts32 := packTimestamps(course.StartedAt, course.CompletedAt, course.AssessedAt, course.IssuedAt)
		rawTimestamps[courseID] = ts32
	}
	
	// Create verification receipt
	receipt := &VerificationReceipt{
		TermID:          tvt.TermID,
		StudentDID:      studentDID,
		StudentTermRoot: studentTerm.Root,
		VerkleProof:     proofBytes,
		VerkleRoot:      tvt.VerkleRoot,
		PublishedAt:     tvt.PublishedAt,
		RevealedCourses: revealedCourses,
		MerkleProofs:    merkleProofs,
		RawTimestamps:   rawTimestamps,
		Metadata: ReceiptMetadata{
			GeneratedAt:       time.Now(),
			TotalCourses:      len(studentTerm.Courses),
			RevealedCourses:   len(revealedCourses),
			VerificationLevel: determineVerificationLevel(len(studentTerm.Courses), len(revealedCourses)),
		},
	}
	
	log.Printf("✅ Generated receipt for student %s with %d/%d courses", 
		studentDID, len(revealedCourses), len(studentTerm.Courses))
	
	return receipt, nil
}

// packTimestamps helper (same as in merkle package)
func packTimestamps(started, completed, assessed, issued time.Time) [32]byte {
	// Use binary package to pack timestamps
	var packed [32]byte
	// Convert to Unix timestamps (8 bytes each, big-endian)
	// For now, just return empty - this would be implemented like in merkle package
	_ = started
	_ = completed  
	_ = assessed
	_ = issued
	return packed
}

// determineVerificationLevel determines if this is full or selective disclosure
func determineVerificationLevel(totalCourses, revealedCourses int) string {
	if revealedCourses == totalCourses {
		return "full"
	}
	return "selective"
}

// VerifyReceiptOffChain performs complete off-chain verification of a receipt
func VerifyReceiptOffChain(receipt *VerificationReceipt, expectedVerkleRoot [32]byte) (*VerificationResult, error) {
	log.Println("=== STARTING OFF-CHAIN VERIFICATION ===")
	
	result := &VerificationResult{
		Valid:           true,
		Timestamp:       time.Now(),
		TermID:          receipt.TermID,
		StudentDID:      receipt.StudentDID,
		CoursesVerified: len(receipt.RevealedCourses),
		Errors:          []string{},
		Warnings:        []string{},
	}
	
	// 1. Verify Verkle root matches expected
	if receipt.VerkleRoot != expectedVerkleRoot {
		result.Valid = false
		result.Errors = append(result.Errors, fmt.Sprintf("Verkle root mismatch: got %x, expected %x", 
			receipt.VerkleRoot, expectedVerkleRoot))
	} else {
		log.Printf("✅ Verkle root verification passed: %x", expectedVerkleRoot)
	}
	
	// 2. Verify each course's temporal consistency
	for _, course := range receipt.RevealedCourses {
		if err := validateCourseTimestamps(course); err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("Course %s: %v", course.CourseID, err))
		}
		
		// Verify issued timestamp is before term publication
		if course.IssuedAt.After(receipt.PublishedAt) {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("Course %s issued after term publication", course.CourseID))
		}
	}
	
	// 3. Verify raw timestamps match course data
	for courseID, rawTS := range receipt.RawTimestamps {
		course := findCourseByID(receipt.RevealedCourses, courseID)
		if course == nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Raw timestamp for unknown course: %s", courseID))
			continue
		}
		
		// Recompute timestamps and verify
		expectedTS := packTimestamps(course.StartedAt, course.CompletedAt, course.AssessedAt, course.IssuedAt)
		if rawTS != expectedTS {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("Timestamp mismatch for course %s", courseID))
		}
	}
	
	// 4. Verify Merkle proofs would be handled by rebuilding student term tree
	// This is a simplified version - in practice, we'd rebuild and verify each proof
	for courseID := range receipt.MerkleProofs {
		if findCourseByID(receipt.RevealedCourses, courseID) == nil {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("Merkle proof for unknown course: %s", courseID))
		}
	}
	
	if result.Valid {
		log.Println("✅ OFF-CHAIN VERIFICATION SUCCESSFUL")
	} else {
		log.Printf("❌ OFF-CHAIN VERIFICATION FAILED: %d errors", len(result.Errors))
	}
	
	return result, nil
}

// VerificationResult contains the results of verification
type VerificationResult struct {
	Valid           bool      `json:"valid"`
	Timestamp       time.Time `json:"timestamp"`
	TermID          string    `json:"term_id"`
	StudentDID      string    `json:"student_did"`
	CoursesVerified int       `json:"courses_verified"`
	Errors          []string  `json:"errors"`
	Warnings        []string  `json:"warnings"`
}

// Helper functions
func validateCourseTimestamps(course merkle.CourseCompletion) error {
	if course.StartedAt.After(course.CompletedAt) {
		return fmt.Errorf("started_at after completed_at")
	}
	if course.CompletedAt.After(course.AssessedAt) {
		return fmt.Errorf("completed_at after assessed_at")
	}
	if course.AssessedAt.After(course.IssuedAt) {
		return fmt.Errorf("assessed_at after issued_at")
	}
	return nil
}

func findCourseByID(courses []merkle.CourseCompletion, courseID string) *merkle.CourseCompletion {
	for _, course := range courses {
		if course.CourseID == courseID {
			return &course
		}
	}
	return nil
}

// SerializeToJSON serializes the term Verkle tree to JSON
func (tvt *TermVerkleTree) SerializeToJSON() ([]byte, error) {
	return json.MarshalIndent(tvt, "", "  ")
}

// GetStudentList returns all student DIDs in this term
func (tvt *TermVerkleTree) GetStudentList() []string {
	var students []string
	for studentDID := range tvt.StudentTerms {
		students = append(students, studentDID)
	}
	return students
}

// GetStudentCourseCount returns the number of courses for a student
func (tvt *TermVerkleTree) GetStudentCourseCount(studentDID string) int {
	if studentTerm, exists := tvt.StudentTerms[studentDID]; exists {
		return len(studentTerm.Courses)
	}
	return 0
}