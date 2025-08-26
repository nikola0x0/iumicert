package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"iumicert/crypto/merkle"
	"iumicert/crypto/testdata"
	"iumicert/crypto/verkle"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run comprehensive testing of the hybrid credential system",
	Long:  `Test the complete IU-MiCert system using the modular data generation and crypto integration programs`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🧪 Starting IU-MiCert Comprehensive System Test")
		fmt.Println("=" + string(make([]byte, 60)))
		
		useModular, _ := cmd.Flags().GetBool("modular")
		if useModular {
			if err := runModularSystemTest(); err != nil {
				log.Fatalf("❌ Modular system test failed: %v", err)
			}
		} else {
			if err := runCompleteSystemTest(); err != nil {
				log.Fatalf("❌ Legacy system test failed: %v", err)
			}
		}
		
		fmt.Println("\n🎉 All tests passed successfully!")
	},
}

func runModularSystemTest() error {
	fmt.Println("\n🔧 Running Modular System Test (using separate programs)")
	
	// Step 1: Run data generation program
	fmt.Println("\n📊 Step 1: Running data generation program...")
	
	dataCmd := exec.Command("./micert", "generate-data")
	dataCmd.Stdout = os.Stdout
	dataCmd.Stderr = os.Stderr
	
	if err := dataCmd.Run(); err != nil {
		return fmt.Errorf("data generation program failed: %w", err)
	}
	
	fmt.Println("✅ Data generation completed successfully!")
	
	// Step 2: Run crypto integration program
	fmt.Println("\n🌳 Step 2: Running crypto integration program...")
	
	cryptoCmd := exec.Command("./micert", "test-crypto")
	cryptoCmd.Stdout = os.Stdout  
	cryptoCmd.Stderr = os.Stderr
	
	if err := cryptoCmd.Run(); err != nil {
		return fmt.Errorf("crypto integration program failed: %w", err)
	}
	
	fmt.Println("✅ Crypto integration completed successfully!")
	
	// Step 3: Validate results
	fmt.Println("\n✅ Step 3: Validating modular test results...")
	
	// Check if expected output files exist
	expectedFiles := []string{
		"data/generated_student_data/system_summary.json",
		"data/generated_student_data/integration_results/verkle_summary_Fall_2023.json",
		"data/generated_student_data/integration_results/verkle_summary_Spring_2024.json",
		"data/generated_student_data/integration_results/verkle_summary_Fall_2024.json", 
		"data/generated_student_data/integration_results/verkle_summary_Spring_2025.json",
	}
	
	for _, file := range expectedFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return fmt.Errorf("expected output file not found: %s", file)
		}
		fmt.Printf("  ✓ Found expected file: %s\n", file)
	}
	
	fmt.Println("\n🎯 Modular test summary:")
	fmt.Println("  ✅ Data generation: SUCCESS")
	fmt.Println("  ✅ Crypto integration: SUCCESS") 
	fmt.Println("  ✅ Output validation: SUCCESS")
	
	return nil
}

func runCompleteSystemTest() error {
	// Step 1: Generate test data
	fmt.Println("\n📊 Step 1: Generating test data...")
	generator := testdata.NewTestDataGenerator()
	
	termID := "Fall_2024"
	completions, err := generator.GenerateTermData(termID, 3, 4) // 3 students, 4 courses each
	if err != nil {
		return fmt.Errorf("failed to generate test data: %w", err)
	}
	
	fmt.Printf("✅ Generated %d course completions for term %s\n", len(completions), termID)
	
	// Step 2: Organize completions by student
	fmt.Println("\n🗂️  Step 2: Organizing data by student...")
	studentCompletions := make(map[string][]merkle.CourseCompletion)
	
	for _, completion := range completions {
		studentDID := fmt.Sprintf("did:example:%s", completion.StudentID)
		studentCompletions[studentDID] = append(studentCompletions[studentDID], completion)
	}
	
	fmt.Printf("✅ Organized completions for %d students\n", len(studentCompletions))
	
	// Step 3: Test Student-Term Merkle Trees
	fmt.Println("\n🌲 Step 3: Testing Student-Term Merkle Trees...")
	studentTrees := make(map[string]*merkle.StudentTermMerkle)
	
	for studentDID, courses := range studentCompletions {
		tree, err := merkle.NewStudentTermMerkle(studentDID, termID, courses)
		if err != nil {
			return fmt.Errorf("failed to create merkle tree for %s: %w", studentDID, err)
		}
		
		studentTrees[studentDID] = tree
		fmt.Printf("  ✓ Student %s: %d courses, root: %x\n", 
			studentDID, len(courses), tree.Root[:8])
		
		// Test Merkle proof verification
		for _, course := range courses {
			valid, err := tree.VerifyProofForCourse(course.CourseID, tree.Root)
			if err != nil {
				return fmt.Errorf("failed to verify merkle proof for course %s: %w", course.CourseID, err)
			}
			if !valid {
				return fmt.Errorf("merkle proof verification failed for course %s", course.CourseID)
			}
		}
		fmt.Printf("  ✓ All Merkle proofs verified for %s\n", studentDID)
	}
	
	// Step 4: Test Term-Level Verkle Tree
	fmt.Println("\n🌳 Step 4: Testing Term-Level Verkle Tree...")
	termTree := verkle.NewTermVerkleTree(termID)
	
	for studentDID, courses := range studentCompletions {
		err := termTree.AddStudent(studentDID, courses)
		if err != nil {
			return fmt.Errorf("failed to add student %s to verkle tree: %w", studentDID, err)
		}
	}
	
	// Publish the term
	err = termTree.PublishTerm()
	if err != nil {
		return fmt.Errorf("failed to publish term: %w", err)
	}
	
	fmt.Printf("✅ Term published with Verkle root: %x\n", termTree.VerkleRoot[:8])
	
	// Step 5: Test Receipt Generation and Verification
	fmt.Println("\n📄 Step 5: Testing Receipt Generation and Verification...")
	
	for studentDID := range studentCompletions {
		// Generate full receipt (all courses)
		fullReceipt, err := termTree.GenerateVerificationReceipt(studentDID, nil)
		if err != nil {
			return fmt.Errorf("failed to generate full receipt for %s: %w", studentDID, err)
		}
		
		fmt.Printf("  ✓ Generated full receipt for %s (%d courses)\n", 
			studentDID, len(fullReceipt.RevealedCourses))
		
		// Verify receipt off-chain
		result, err := verkle.VerifyReceiptOffChain(fullReceipt, termTree.VerkleRoot)
		if err != nil {
			return fmt.Errorf("failed to verify receipt for %s: %w", studentDID, err)
		}
		
		if !result.Valid {
			return fmt.Errorf("receipt verification failed for %s: %v", studentDID, result.Errors)
		}
		
		fmt.Printf("  ✅ Receipt verified successfully for %s\n", studentDID)
		
		// Test selective disclosure (first 2 courses only)
		if len(studentCompletions[studentDID]) >= 2 {
			courseIDs := []string{
				studentCompletions[studentDID][0].CourseID,
				studentCompletions[studentDID][1].CourseID,
			}
			
			selectiveReceipt, err := termTree.GenerateVerificationReceipt(studentDID, courseIDs)
			if err != nil {
				return fmt.Errorf("failed to generate selective receipt for %s: %w", studentDID, err)
			}
			
			fmt.Printf("  ✓ Generated selective receipt for %s (%d/%d courses)\n", 
				studentDID, len(selectiveReceipt.RevealedCourses), len(studentCompletions[studentDID]))
			
			// Verify selective receipt
			result, err := verkle.VerifyReceiptOffChain(selectiveReceipt, termTree.VerkleRoot)
			if err != nil {
				return fmt.Errorf("failed to verify selective receipt for %s: %w", studentDID, err)
			}
			
			if !result.Valid {
				return fmt.Errorf("selective receipt verification failed for %s: %v", studentDID, result.Errors)
			}
			
			fmt.Printf("  ✅ Selective receipt verified successfully for %s\n", studentDID)
		}
	}
	
	// Step 6: Test Error Cases
	fmt.Println("\n🚨 Step 6: Testing Error Cases...")
	
	// Test invalid student
	_, err = termTree.GenerateVerificationReceipt("did:example:nonexistent", nil)
	if err == nil {
		return fmt.Errorf("expected error for nonexistent student, but got none")
	}
	fmt.Printf("  ✓ Correctly rejected nonexistent student\n")
	
	// Test invalid course ID
	firstStudentDID := getFirstKey(studentCompletions)
	_, err = termTree.GenerateVerificationReceipt(firstStudentDID, []string{"INVALID_COURSE"})
	if err != nil {
		fmt.Printf("  ✓ Correctly handled invalid course ID\n")
	}
	
	// Test receipt with wrong Verkle root
	firstStudentDID = getFirstKey(studentCompletions)
	receipt, err := termTree.GenerateVerificationReceipt(firstStudentDID, nil)
	if err != nil {
		return fmt.Errorf("failed to generate test receipt: %w", err)
	}
	
	// Create wrong root
	wrongRoot := termTree.VerkleRoot
	wrongRoot[0] = wrongRoot[0] ^ 0xFF // Flip some bits
	
	result, err := verkle.VerifyReceiptOffChain(receipt, wrongRoot)
	if err != nil {
		return fmt.Errorf("error during wrong root test: %w", err)
	}
	
	if result.Valid {
		return fmt.Errorf("expected verification to fail with wrong root, but it passed")
	}
	fmt.Printf("  ✓ Correctly rejected receipt with wrong Verkle root\n")
	
	return nil
}

func getFirstKey(m map[string][]merkle.CourseCompletion) string {
	for k := range m {
		return k
	}
	return ""
}

func init() {
	testCmd.Flags().Bool("modular", false, "Use modular test approach (run separate programs)")
	rootCmd.AddCommand(testCmd)
}