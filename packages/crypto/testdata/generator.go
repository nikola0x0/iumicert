package testdata

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
	
	"iumicert/crypto/verkle"
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
			"IT064IU":   {"IT064IU", "Introduction to Computing", 3, "Information Technology"},
			"IT116IU":   {"IT116IU", "C/C++ Programming", 4, "Information Technology"},
			"IT069IU":   {"IT069IU", "Object-Oriented Programming", 4, "Information Technology"},
			"IT013IU":   {"IT013IU", "Algorithms & Data Structures", 4, "Information Technology"},
			"IT076IU":   {"IT076IU", "Software Engineering", 4, "Information Technology"},
			"IT091IU":   {"IT091IU", "Computer Networks", 4, "Information Technology"},
			"IT117IU":   {"IT117IU", "System and Network Security", 4, "Information Technology"},
			"IT160IU":   {"IT160IU", "Data Mining", 4, "Information Technology"},
			"IT134IU":   {"IT134IU", "Internet of Things", 4, "Information Technology"},
			"IT133IU":   {"IT133IU", "Mobile Application Development", 4, "Information Technology"},
			"MA001IU":   {"MA001IU", "Calculus 1", 4, "Mathematics"},
			"MA003IU":   {"MA003IU", "Calculus 2", 4, "Mathematics"},
			"IT154IU":   {"IT154IU", "Linear Algebra", 3, "Mathematics"},
			"MA026IU":   {"MA026IU", "Probability, Statistic & Random Process", 3, "Mathematics"},
			"PH013IU":   {"PH013IU", "Physics 1", 2, "Physics"},
			"PH014IU":   {"PH014IU", "Physics 2", 2, "Physics"},
			"PH015IU":   {"PH015IU", "Physics 3", 3, "Physics"},
			"CH011IU":   {"CH011IU", "Chemistry for Engineers", 3, "Chemistry"},
			"EN007IU":   {"EN007IU", "Writing AE1", 2, "English"},
			"EN008IU":   {"EN008IU", "Listening AE1", 2, "English"},
			"PE008IU":   {"PE008IU", "Critical Thinking", 3, "Philosophy"},
			"IT153IU":   {"IT153IU", "Discrete Mathematics", 3, "Mathematics"},
			"IT067IU":   {"IT067IU", "Digital Logic Design", 3, "Information Technology"},
			"IT089IU":   {"IT089IU", "Computer Architecture", 4, "Information Technology"},
		},
		students: []StudentInfo{
			{"ITITIU00001", "did:iu:student:ITITIU00001", "Nguyen Van Minh", "Computer Science"},
			{"ITITIU00002", "did:iu:student:ITITIU00002", "Tran Thi Lan", "Computer Science"}, 
			{"ITITIU00003", "did:iu:student:ITITIU00003", "Le Hoang Nam", "Information Technology"},
			{"ITITIU00004", "did:iu:student:ITITIU00004", "Pham Thi Hoa", "Computer Science"},
			{"ITITIU00005", "did:iu:student:ITITIU00005", "Vo Van Duc", "Information Technology"},
		},
		terms: []string{
			"Semester_1_2023", "Semester_2_2023", "Summer_2023", 
			"Semester_1_2024", "Semester_2_2024", "Summer_2024",
			"Semester_1_2025", "Semester_2_2025",
		},
		grades: []string{"A+", "A", "A-", "B+", "B", "B-", "C+", "C", "C-", "D+", "D"},
		rand:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GenerateTermData generates course completions for a specific term
func (g *TestDataGenerator) GenerateTermData(termID string, numStudents, coursesPerStudent int) ([]verkle.CourseCompletion, error) {
	if numStudents > len(g.students) {
		numStudents = len(g.students)
	}
	
	var completions []verkle.CourseCompletion
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
func (g *TestDataGenerator) generateCourseCompletion(student StudentInfo, course CourseInfo, termID string, termStart, termEnd time.Time) verkle.CourseCompletion {
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
	
	return verkle.CourseCompletion{
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
	// First check predefined terms
	switch termID {
	case "Semester_1_2023":
		return time.Date(2023, 8, 15, 0, 0, 0, 0, time.UTC),
			   time.Date(2023, 12, 20, 23, 59, 59, 0, time.UTC), nil
	case "Semester_2_2023":
		return time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
			   time.Date(2024, 5, 15, 23, 59, 59, 0, time.UTC), nil
	case "Summer_2023":
		return time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
			   time.Date(2023, 7, 31, 23, 59, 59, 0, time.UTC), nil
	case "Semester_1_2024":
		return time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC),
			   time.Date(2024, 12, 20, 23, 59, 59, 0, time.UTC), nil
	case "Semester_2_2024":
		return time.Date(2025, 1, 8, 0, 0, 0, 0, time.UTC),
			   time.Date(2025, 5, 15, 23, 59, 59, 0, time.UTC), nil
	case "Summer_2024":
		return time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
			   time.Date(2024, 7, 31, 23, 59, 59, 0, time.UTC), nil
	case "Semester_1_2025":
		return time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
			   time.Date(2025, 12, 20, 23, 59, 59, 0, time.UTC), nil
	case "Semester_2_2025":
		return time.Date(2026, 1, 8, 0, 0, 0, 0, time.UTC),
			   time.Date(2026, 5, 15, 23, 59, 59, 0, time.UTC), nil
	}

	// Try to parse term ID dynamically for custom terms
	// Expected format: "Semester_1_YYYY", "Semester_2_YYYY", or "Summer_YYYY"
	return g.parseTermDates(termID)
}

// parseTermDates dynamically parses term dates from term ID
func (g *TestDataGenerator) parseTermDates(termID string) (time.Time, time.Time, error) {
	var year int

	// Try Semester_1_YYYY format
	if _, err := fmt.Sscanf(termID, "Semester_1_%d", &year); err == nil {
		return time.Date(year, 8, 15, 0, 0, 0, 0, time.UTC),
			   time.Date(year, 12, 20, 23, 59, 59, 0, time.UTC), nil
	}

	// Try Semester_2_YYYY format (runs in YYYY+1)
	if _, err := fmt.Sscanf(termID, "Semester_2_%d", &year); err == nil {
		return time.Date(year+1, 1, 8, 0, 0, 0, 0, time.UTC),
			   time.Date(year+1, 5, 15, 23, 59, 59, 0, time.UTC), nil
	}

	// Try Summer_YYYY format
	if _, err := fmt.Sscanf(termID, "Summer_%d", &year); err == nil {
		return time.Date(year, 6, 1, 0, 0, 0, 0, time.UTC),
			   time.Date(year, 7, 31, 23, 59, 59, 0, time.UTC), nil
	}

	return time.Time{}, time.Time{}, fmt.Errorf("unknown term format: %s (expected Semester_1_YYYY, Semester_2_YYYY, or Summer_YYYY)", termID)
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
func (g *TestDataGenerator) GenerateMultiTermJourney(studentIndex int, termIDs []string, coursesPerTerm int) (map[string][]verkle.CourseCompletion, error) {
	if studentIndex >= len(g.students) {
		return nil, fmt.Errorf("student index %d out of range", studentIndex)
	}
	
	student := g.students[studentIndex]
	journey := make(map[string][]verkle.CourseCompletion)
	
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
func (g *TestDataGenerator) generateStudentTermCompletions(student StudentInfo, termID string, courseCount int) ([]verkle.CourseCompletion, error) {
	termStart, termEnd, err := g.getTermDates(termID)
	if err != nil {
		return nil, err
	}
	
	courses := g.getRandomCourses(courseCount)
	var completions []verkle.CourseCompletion
	
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