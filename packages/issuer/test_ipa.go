package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	
	"github.com/iumicert/issuer/crypto/verkle"
)

func main() {
	fmt.Println("🔬 Testing FULL IPA Cryptographic Verification")
	fmt.Println("==============================================")
	
	// Read a receipt
	data, err := os.ReadFile("publish_ready/receipts/ITITIU00001_journey.json")
	if err != nil {
		log.Fatal(err)
	}
	
	var receipt map[string]interface{}
	json.Unmarshal(data, &receipt)
	
	termReceipts := receipt["term_receipts"].(map[string]interface{})
	
	// Test Semester_1_2023
	termData := termReceipts["Semester_1_2023"].(map[string]interface{})
	verkleRootHex := termData["verkle_root"].(string)
	
	rootBytes, _ := hex.DecodeString(verkleRootHex)
	var verkleRoot [32]byte
	copy(verkleRoot[:], rootBytes)
	
	receiptData := termData["receipt"].(map[string]interface{})
	courseProofs := receiptData["course_proofs"].(map[string]interface{})
	revealedCourses := receiptData["revealed_courses"].([]interface{})
	
	fmt.Printf("Testing verification for %d courses\n", len(courseProofs))
	fmt.Printf("Verkle Root: %x\n\n", verkleRoot[:16])
	
	for courseID, proofStr := range courseProofs {
		// Find course data
		var course verkle.CourseCompletion
		for _, c := range revealedCourses {
			cMap := c.(map[string]interface{})
			if cMap["course_id"] == courseID {
				course.CourseID = courseID
				course.StudentID = "did:example:ITITIU00001"
				course.TermID = "Semester_1_2023"
				course.CourseName = cMap["course_name"].(string)
				course.Grade = cMap["grade"].(string)
				course.Credits = uint8(cMap["credits"].(float64))
				break
			}
		}
		
		courseKey := fmt.Sprintf("did:example:ITITIU00001:Semester_1_2023:%s", courseID)
		proofData := []byte(proofStr.(string))
		
		fmt.Printf("Verifying %s... ", courseID)
		err := verkle.VerifyCourseProof(courseKey, course, proofData, verkleRoot)
		if err != nil {
			fmt.Printf("❌ FAILED: %v\n", err)
		} else {
			fmt.Printf("✅ SUCCESS (IPA verified!)\n")
		}
	}
}