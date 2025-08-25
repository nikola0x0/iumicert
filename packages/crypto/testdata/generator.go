package testdata

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"iumicert/crypto/merkle"
)

// TestDataGenerator generates realistic academic data for testing
type TestDataGenerator struct {
	institutions []string
	courses      map[string]CourseInfo
	students     []StudentInfo
	terms        []string
	grades       []string
	rand         *rand.Rand
}

// CourseInfo represents course metadata
type CourseInfo struct {
	CourseID   string `json:"course_id"`
	CourseName string `json:"course_name"`
	Credits    uint8  `json:"credits"`
	Department string `json:"department"`
}

// StudentInfo represents student metadata
type StudentInfo struct {
	StudentID string `json:"student_id"`
	DID       string `json:"did"`
	Name      string `json:"name"`
	Program   string `json:"program"`
}

// NewTestDataGenerator creates a new test data generator
func NewTestDataGenerator() *TestDataGenerator {
	return &TestDataGenerator{
		institutions: []string{
			"IU-CS", "IU-MATH", "IU-PHYS", "IU-CHEM", "IU-BUS",
		},
		courses: map[string]CourseInfo{
			"CS101":   {"CS101", "Introduction to Programming", 3, "Computer Science"},
			"CS102":   {"CS102", "Data Structures", 3, "Computer Science"},
			"CS201":   {"CS201", "Algorithms", 3, "Computer Science"},
			"CS301":   {"CS301", "Database Systems", 3, "Computer Science"},
			"MATH101": {"MATH101", "Calculus I", 4, "Mathematics"},
			"MATH102": {"MATH102", "Calculus II", 4, "Mathematics"},
			"MATH201": {"MATH201", "Linear Algebra", 3, "Mathematics"},
			"PHYS101": {"PHYS101", "Physics I", 4, "Physics"},
			"PHYS102": {"PHYS102", "Physics II", 4, "Physics"},
			"CHEM101": {"CHEM101", "General Chemistry", 4, "Chemistry"},
			"ENG101":  {"ENG101", "English Composition", 3, "English"},
			"HIST101": {"HIST101", "World History", 3, "History"},
			"BUS101":  {"BUS101", "Business Fundamentals", 3, "Business"},
		},
		students: []StudentInfo{
			{"STU001", "did:example:student001", "Alice Johnson", "Computer Science"},
			{"STU002", "did:example:student002", "Bob Smith", "Computer Science"},
			{"STU003", "did:example:student003", "Carol Davis", "Mathematics"},
			{"STU004", "did:example:student004", "David Wilson", "Physics"},
			{"STU005", "did:example:student005", "Eva Brown", "Computer Science"},
		},
		terms: []string{
			"Fall_2021", "Spring_2022", "Fall_2022", "Spring_2023", 
			"Fall_2023", "Spring_2024", "Fall_2024", "Spring_2025",
		},
		grades: []string{"A+", "A", "A-", "B+", "B", "B-", "C+", "C", "C-", "D+", "D"},
		rand:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GenerateTermData generates course completions for a specific term
func (g *TestDataGenerator) GenerateTermData(termID string, numStudents, coursesPerStudent int) ([]merkle.CourseCompletion, error) {
	if numStudents > len(g.students) {
		numStudents = len(g.students)
	}
	
	var completions []merkle.CourseCompletion
	termStart, termEnd, err := g.getTermDates(termID)
	if err != nil {
		return nil, err
	}
	
	// Get random course list
	courseList := g.getRandomCourses(coursesPerStudent * numStudents)
	courseIndex := 0
	
	for i := 0; i < numStudents; i++ {
		student := g.students[i]
		
		for j := 0; j < coursesPerStudent; j++ {
			if courseIndex >= len(courseList) {
				break
			}
			
			course := courseList[courseIndex]
			courseIndex++
			
			completion := g.generateCourseCompletion(student, course, termID, termStart, termEnd)
			completions = append(completions, completion)
		}
	}
	
	return completions, nil
}

// generateCourseCompletion creates a single course completion with realistic timestamps
func (g *TestDataGenerator) generateCourseCompletion(student StudentInfo, course CourseInfo, termID string, termStart, termEnd time.Time) merkle.CourseCompletion {
	// Generate realistic timeline within term, ensuring all timestamps are valid
	duration := termEnd.Sub(termStart)
	
	// Reserve buffer time for processing (21 days max for assessment + issuance)
	bufferDays := 21
	bufferTime := time.Duration(bufferDays) * 24 * time.Hour
	effectiveTermEnd := termEnd.Add(-bufferTime)
	
	// Started: 0-20% into term
	startOffset := time.Duration(float64(duration) * g.rand.Float64() * 0.2)
	startedAt := termStart.Add(startOffset)
	
	// Completed: 60-80% into effective term (leaving room for processing)
	effectiveDuration := effectiveTermEnd.Sub(termStart)
	completedOffset := time.Duration(float64(effectiveDuration) * (0.6 + g.rand.Float64()*0.2))
	completedAt := termStart.Add(completedOffset)
	
	// Assessed: 1-7 days after completion
	assessedAt := completedAt.Add(time.Duration(1+g.rand.Intn(7)) * 24 * time.Hour)
	
	// Issued: 1-14 days after assessment  
	issuedAt := assessedAt.Add(time.Duration(1+g.rand.Intn(14)) * 24 * time.Hour)
	
	// Final safety check - if still after term end, compress timeline
	if issuedAt.After(termEnd) {
		// Work backwards from term end with safe margins
		issuedAt = termEnd.Add(-time.Duration(1+g.rand.Intn(3)) * 24 * time.Hour)
		assessedAt = issuedAt.Add(-time.Duration(1+g.rand.Intn(3)) * 24 * time.Hour)
		completedAt = assessedAt.Add(-time.Duration(1+g.rand.Intn(3)) * 24 * time.Hour)
		// Keep original startedAt if it's still valid
		if startedAt.After(completedAt) {
			startedAt = completedAt.Add(-time.Duration(7+g.rand.Intn(14)) * 24 * time.Hour)
		}
	}
	
	return merkle.CourseCompletion{
		IssuerID:    g.getRandomInstitution(),
		StudentID:   student.StudentID,
		TermID:      termID,
		CourseID:    course.CourseID,
		CourseName:  course.CourseName,
		AttemptNo:   1, // Assume first attempt
		StartedAt:   startedAt,
		CompletedAt: completedAt,
		AssessedAt:  assessedAt,
		IssuedAt:    issuedAt,
		Grade:       g.getRandomGrade(),
		Credits:     course.Credits,
		Instructor:  fmt.Sprintf("Prof. %s", g.generateRandomName()),
	}
}

// getTermDates returns start and end dates for a term
func (g *TestDataGenerator) getTermDates(termID string) (time.Time, time.Time, error) {
	switch {
	case termID == "Fall_2021":
		return time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC), 
			   time.Date(2021, 12, 31, 23, 59, 59, 0, time.UTC), nil
	case termID == "Spring_2022":
		return time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC), 
			   time.Date(2022, 5, 15, 23, 59, 59, 0, time.UTC), nil
	case termID == "Fall_2022":
		return time.Date(2022, 9, 1, 0, 0, 0, 0, time.UTC), 
			   time.Date(2022, 12, 31, 23, 59, 59, 0, time.UTC), nil
	case termID == "Spring_2023":
		return time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC), 
			   time.Date(2023, 5, 15, 23, 59, 59, 0, time.UTC), nil
	case termID == "Fall_2023":
		return time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC), 
			   time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC), nil
	case termID == "Spring_2024":
		return time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), 
			   time.Date(2024, 5, 15, 23, 59, 59, 0, time.UTC), nil
	case termID == "Fall_2024":
		return time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC), 
			   time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC), nil
	case termID == "Spring_2025":
		return time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC), 
			   time.Date(2025, 5, 15, 23, 59, 59, 0, time.UTC), nil
	default:
		return time.Time{}, time.Time{}, fmt.Errorf("unknown term: %s", termID)
	}
}

// Helper methods
func (g *TestDataGenerator) getRandomInstitution() string {
	return g.institutions[g.rand.Intn(len(g.institutions))]
}

func (g *TestDataGenerator) getRandomGrade() string {
	return g.grades[g.rand.Intn(len(g.grades))]
}

func (g *TestDataGenerator) getRandomCourses(count int) []CourseInfo {
	var courses []CourseInfo
	courseIDs := make([]string, 0, len(g.courses))
	for id := range g.courses {
		courseIDs = append(courseIDs, id)
	}
	
	for i := 0; i < count; i++ {
		courseID := courseIDs[g.rand.Intn(len(courseIDs))]
		courses = append(courses, g.courses[courseID])
	}
	
	return courses
}

func (g *TestDataGenerator) generateRandomName() string {
	names := []string{
		"Johnson", "Smith", "Davis", "Wilson", "Brown", "Miller", 
		"Anderson", "Taylor", "Thomas", "Jackson", "White", "Harris",
	}
	return names[g.rand.Intn(len(names))]
}

// GenerateMultiTermJourney generates a complete academic journey across multiple terms
func (g *TestDataGenerator) GenerateMultiTermJourney(studentIndex int, termIDs []string, coursesPerTerm int) (map[string][]merkle.CourseCompletion, error) {
	if studentIndex >= len(g.students) {
		return nil, fmt.Errorf("student index %d out of range", studentIndex)
	}
	
	student := g.students[studentIndex]
	journey := make(map[string][]merkle.CourseCompletion)
	
	for _, termID := range termIDs {
		completions, err := g.generateStudentTermCompletions(student, termID, coursesPerTerm)
		if err != nil {
			return nil, fmt.Errorf("failed to generate completions for %s in %s: %w", student.StudentID, termID, err)
		}
		journey[termID] = completions
	}
	
	return journey, nil
}

// generateStudentTermCompletions generates completions for one student in one term
func (g *TestDataGenerator) generateStudentTermCompletions(student StudentInfo, termID string, courseCount int) ([]merkle.CourseCompletion, error) {
	termStart, termEnd, err := g.getTermDates(termID)
	if err != nil {
		return nil, err
	}
	
	courses := g.getRandomCourses(courseCount)
	var completions []merkle.CourseCompletion
	
	for _, course := range courses {
		completion := g.generateCourseCompletion(student, course, termID, termStart, termEnd)
		completions = append(completions, completion)
	}
	
	return completions, nil
}

// SaveToJSON saves test data to JSON file
func (g *TestDataGenerator) SaveToJSON(data interface{}, filename string) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	
	// In a real implementation, this would write to a file
	fmt.Printf("Generated test data would be saved to: %s\n", filename)
	fmt.Printf("Data size: %d bytes\n", len(jsonData))
	return nil
}

// GetAvailableTerms returns all available terms
func (g *TestDataGenerator) GetAvailableTerms() []string {
	return g.terms
}

// GetAvailableStudents returns all available students
func (g *TestDataGenerator) GetAvailableStudents() []StudentInfo {
	return g.students
}

// GetAvailableCourses returns all available courses
func (g *TestDataGenerator) GetAvailableCourses() map[string]CourseInfo {
	return g.courses
}