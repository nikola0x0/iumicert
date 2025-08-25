package merkle

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/cbergoon/merkletree"
)

// CourseCompletion represents a single course completion with all metadata
type CourseCompletion struct {
	IssuerID     string    `json:"issuer_id"`
	StudentID    string    `json:"student_id"`
	TermID       string    `json:"term_id"`
	CourseID     string    `json:"course_id"`
	CourseName   string    `json:"course_name"`
	AttemptNo    uint8     `json:"attempt_no"`
	StartedAt    time.Time `json:"started_at"`
	CompletedAt  time.Time `json:"completed_at"`
	AssessedAt   time.Time `json:"assessed_at"`
	IssuedAt     time.Time `json:"issued_at"`
	Grade        string    `json:"grade"`
	Credits      uint8     `json:"credits"`
	Instructor   string    `json:"instructor"`
}

// CourseLeaf implements merkletree.Content interface for course completions
type CourseLeaf struct {
	Course       CourseCompletion `json:"course"`
	Timestamp32  [32]byte         `json:"timestamp32"`  // 4×u64 packed into 32 bytes
	CourseHash32 [32]byte         `json:"course_hash32"` // H256(issuer_id || student_id || term_id || course_id || attempt_no)
}

// CalculateHash computes H256(ts32 || courseHash256) for the Merkle leaf
func (cl CourseLeaf) CalculateHash() ([]byte, error) {
	h := sha256.New()
	// Write timestamp32 first, then courseHash32
	h.Write(cl.Timestamp32[:])
	h.Write(cl.CourseHash32[:])
	return h.Sum(nil), nil
}

// Equals checks equality of two CourseLeaf items
func (cl CourseLeaf) Equals(other merkletree.Content) (bool, error) {
	otherCL, ok := other.(CourseLeaf)
	if !ok {
		return false, fmt.Errorf("value is not of type CourseLeaf")
	}
	return cl.Timestamp32 == otherCL.Timestamp32 && cl.CourseHash32 == otherCL.CourseHash32, nil
}

// StudentTermMerkle manages a Merkle tree for a student's term courses
type StudentTermMerkle struct {
	StudentID   string                   `json:"student_id"`
	TermID      string                   `json:"term_id"`
	Courses     []CourseCompletion       `json:"courses"`
	Tree        *merkletree.MerkleTree   `json:"-"` // Don't serialize the tree object
	Root        [32]byte                 `json:"root"`
	Leaves      []CourseLeaf            `json:"leaves"`
	Proofs      map[string][]string     `json:"proofs"` // courseID -> merkle path (hex strings)
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

// sortCoursesByCanonicalOrder sorts courses by (course_id ASC, attempt_no ASC, completed_at ASC)
func sortCoursesByCanonicalOrder(courses []CourseCompletion) {
	sort.Slice(courses, func(i, j int) bool {
		if courses[i].CourseID != courses[j].CourseID {
			return courses[i].CourseID < courses[j].CourseID
		}
		if courses[i].AttemptNo != courses[j].AttemptNo {
			return courses[i].AttemptNo < courses[j].AttemptNo
		}
		return courses[i].CompletedAt.Before(courses[j].CompletedAt)
	})
}

// NewStudentTermMerkle creates a new student-term Merkle tree
func NewStudentTermMerkle(studentID, termID string, courses []CourseCompletion) (*StudentTermMerkle, error) {
	if len(courses) == 0 {
		return nil, fmt.Errorf("no courses provided for student %s in term %s", studentID, termID)
	}
	
	// Sort courses for canonical order
	sortedCourses := make([]CourseCompletion, len(courses))
	copy(sortedCourses, courses)
	sortCoursesByCanonicalOrder(sortedCourses)
	
	// Validate timestamps for each course
	for _, course := range sortedCourses {
		if err := validateCourseTimestamps(course); err != nil {
			return nil, fmt.Errorf("invalid timestamps for course %s: %w", course.CourseID, err)
		}
	}
	
	// Create course leaves
	var leaves []CourseLeaf
	var leafContents []merkletree.Content
	
	for _, course := range sortedCourses {
		ts32 := packTimestamps(course.StartedAt, course.CompletedAt, course.AssessedAt, course.IssuedAt)
		courseHash32 := computeCourseHash(course)
		
		leaf := CourseLeaf{
			Course:       course,
			Timestamp32:  ts32,
			CourseHash32: courseHash32,
		}
		
		leaves = append(leaves, leaf)
		leafContents = append(leafContents, leaf)
	}
	
	// Build Merkle tree
	tree, err := merkletree.NewTree(leafContents)
	if err != nil {
		return nil, fmt.Errorf("failed to create merkle tree: %w", err)
	}
	
	// Get root as [32]byte
	root := tree.MerkleRoot()
	var root32 [32]byte
	copy(root32[:], root)
	
	// Generate proofs for all courses
	proofs := make(map[string][]string)
	for i, course := range sortedCourses {
		proof, _, err := tree.GetMerklePath(leafContents[i])
		if err != nil {
			return nil, fmt.Errorf("failed to generate proof for course %s: %w", course.CourseID, err)
		}
		
		// Convert proof to hex strings for JSON serialization
		var proofHex []string
		for _, hash := range proof {
			proofHex = append(proofHex, hex.EncodeToString(hash))
		}
		proofs[course.CourseID] = proofHex
	}
	
	return &StudentTermMerkle{
		StudentID: studentID,
		TermID:    termID,
		Courses:   sortedCourses,
		Tree:      tree,
		Root:      root32,
		Leaves:    leaves,
		Proofs:    proofs,
	}, nil
}

// validateCourseTimestamps ensures timestamps are in correct order
func validateCourseTimestamps(course CourseCompletion) error {
	if course.StartedAt.After(course.CompletedAt) {
		return fmt.Errorf("started_at (%s) is after completed_at (%s)", 
			course.StartedAt, course.CompletedAt)
	}
	if course.CompletedAt.After(course.AssessedAt) {
		return fmt.Errorf("completed_at (%s) is after assessed_at (%s)", 
			course.CompletedAt, course.AssessedAt)
	}
	if course.AssessedAt.After(course.IssuedAt) {
		return fmt.Errorf("assessed_at (%s) is after issued_at (%s)", 
			course.AssessedAt, course.IssuedAt)
	}
	return nil
}

// GetProofForCourse returns the Merkle proof for a specific course
func (stm *StudentTermMerkle) GetProofForCourse(courseID string) ([]string, error) {
	proof, exists := stm.Proofs[courseID]
	if !exists {
		return nil, fmt.Errorf("no proof found for course %s", courseID)
	}
	return proof, nil
}

// VerifyProofForCourse verifies a Merkle proof for a specific course
func (stm *StudentTermMerkle) VerifyProofForCourse(courseID string, expectedRoot [32]byte) (bool, error) {
	// Find the course leaf
	var targetLeaf CourseLeaf
	found := false
	for _, leaf := range stm.Leaves {
		if leaf.Course.CourseID == courseID {
			targetLeaf = leaf
			found = true
			break
		}
	}
	
	if !found {
		return false, fmt.Errorf("course %s not found in tree", courseID)
	}
	
	// Get proof
	proofHex, err := stm.GetProofForCourse(courseID)
	if err != nil {
		return false, err
	}
	
	// Convert hex proof back to bytes
	var proof [][]byte
	for _, hexStr := range proofHex {
		hashBytes, err := hex.DecodeString(hexStr)
		if err != nil {
			return false, fmt.Errorf("failed to decode proof hash: %w", err)
		}
		proof = append(proof, hashBytes)
	}
	
	// Verify using the tree
	if stm.Tree == nil {
		return false, fmt.Errorf("merkle tree not available for verification")
	}
	
	return stm.Tree.VerifyContent(targetLeaf)
}

// SerializeToJSON serializes the student term data to JSON
func (stm *StudentTermMerkle) SerializeToJSON() ([]byte, error) {
	return json.MarshalIndent(stm, "", "  ")
}

// DeserializeFromJSON creates a StudentTermMerkle from JSON data
func DeserializeFromJSON(data []byte) (*StudentTermMerkle, error) {
	var stm StudentTermMerkle
	if err := json.Unmarshal(data, &stm); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	
	// Rebuild the Merkle tree from leaves
	var leafContents []merkletree.Content
	for _, leaf := range stm.Leaves {
		leafContents = append(leafContents, leaf)
	}
	
	tree, err := merkletree.NewTree(leafContents)
	if err != nil {
		return nil, fmt.Errorf("failed to rebuild merkle tree: %w", err)
	}
	
	stm.Tree = tree
	return &stm, nil
}

// GetCourseByID retrieves a specific course by ID
func (stm *StudentTermMerkle) GetCourseByID(courseID string) (*CourseCompletion, error) {
	for _, course := range stm.Courses {
		if course.CourseID == courseID {
			return &course, nil
		}
	}
	return nil, fmt.Errorf("course %s not found", courseID)
}

// ListCourses returns all course IDs in this term
func (stm *StudentTermMerkle) ListCourses() []string {
	var courseIDs []string
	for _, course := range stm.Courses {
		courseIDs = append(courseIDs, course.CourseID)
	}
	return courseIDs
}