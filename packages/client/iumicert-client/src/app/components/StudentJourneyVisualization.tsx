import React, { useState, useEffect } from 'react';

const StudentJourneyVisualization = () => {
  const [students, setStudents] = useState([]);
  const [selectedStudentId, setSelectedStudentId] = useState('');
  const [studentJourney, setStudentJourney] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    loadStudents();
  }, []);

  const loadStudents = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/students');
      const result = await response.json();
      const studentList = result.success ? (result.data || []) : [];
      setStudents(studentList.slice(0, 10)); // Show first 10 for demo
      if (studentList.length > 0) {
        setSelectedStudentId(studentList[0].student_id);
        loadStudentJourney(studentList[0].student_id);
      }
    } catch (error) {
      console.error('Failed to load students:', error);
      setError('Failed to load students');
    }
  };

  const loadStudentJourney = async (studentId) => {
    if (!studentId) return;
    
    setLoading(true);
    setError('');
    
    try {
      // Load student journey data from API
      const response = await fetch(`http://localhost:8080/api/students/${studentId}/journey`);
      if (!response.ok) {
        throw new Error('Student journey not found');
      }
      const result = await response.json();
      if (result.success) {
        setStudentJourney(result.data);
      } else {
        throw new Error(result.error || 'Failed to load journey');
      }
    } catch (error) {
      console.error('Failed to load student journey:', error);
      setError('Failed to load student journey');
    } finally {
      setLoading(false);
    }
  };

  const handleStudentChange = (studentId) => {
    setSelectedStudentId(studentId);
    loadStudentJourney(studentId);
  };

  const getGradeColor = (grade) => {
    const gradeColors = {
      'A+': 'bg-green-500', 'A': 'bg-green-500', 'A-': 'bg-green-400',
      'B+': 'bg-blue-400', 'B': 'bg-blue-400', 'B-': 'bg-blue-300',
      'C+': 'bg-yellow-400', 'C': 'bg-yellow-400', 'C-': 'bg-yellow-300',
      'D+': 'bg-orange-400', 'D': 'bg-orange-400',
      'F': 'bg-red-400'
    };
    return gradeColors[grade] || 'bg-gray-400';
  };

  const calculateTermGPA = (courses) => {
    const gradePoints = {
      'A+': 4.0, 'A': 4.0, 'A-': 3.7,
      'B+': 3.3, 'B': 3.0, 'B-': 2.7,
      'C+': 2.3, 'C': 2.0, 'C-': 1.7,
      'D+': 1.3, 'D': 1.0, 'F': 0.0
    };
    
    const total = courses.reduce((sum, course) => sum + (gradePoints[course.grade] || 0), 0);
    return (total / courses.length).toFixed(2);
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('en-US', { 
      year: 'numeric', 
      month: 'short'
    });
  };

  if (loading) {
    return (
      <div className="bg-white rounded-lg shadow-md p-6">
        <div className="flex items-center justify-center py-8">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          <span className="ml-2 text-gray-600">Loading student journey...</span>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <div className="mb-6">
        <h2 className="text-2xl font-bold text-gray-900 mb-4">Student Academic Journey</h2>
        
        {/* Student Selection */}
        <div className="mb-4">
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Select Student
          </label>
          <select
            value={selectedStudentId}
            onChange={(e) => handleStudentChange(e.target.value)}
            className="block w-full max-w-md px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="">Select a student...</option>
            {students.map((student) => (
              <option key={student.student_id} value={student.student_id}>
                {student.student_id} - Student {student.student_id}
              </option>
            ))}
          </select>
        </div>

        {error && (
          <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-md">
            <p className="text-red-600">{error}</p>
          </div>
        )}
      </div>

      {studentJourney && (
        <div className="space-y-6">
          {/* Student Overview */}
          <div className="bg-blue-50 rounded-lg p-4">
            <h3 className="text-lg font-semibold text-blue-900 mb-2">
              Academic Overview - {selectedStudentId}
            </h3>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <span className="text-sm text-blue-700">Total Terms:</span>
                <div className="font-semibold text-blue-900">
                  {Object.keys(studentJourney.terms || {}).length}
                </div>
              </div>
              <div>
                <span className="text-sm text-blue-700">Total Courses:</span>
                <div className="font-semibold text-blue-900">
                  {Object.values(studentJourney.terms || {}).reduce((sum, term) => 
                    sum + (term.courses?.length || 0), 0
                  )}
                </div>
              </div>
              <div>
                <span className="text-sm text-blue-700">Overall GPA:</span>
                <div className="font-semibold text-blue-900">
                  {Object.values(studentJourney.terms || {}).length > 0 ? (
                    Object.values(studentJourney.terms || {})
                      .reduce((sum, term) => sum + parseFloat(term.gpa || 0), 0) / 
                      Object.values(studentJourney.terms || {}).length
                  ).toFixed(2) : '0.00'}
                </div>
              </div>
            </div>
          </div>

          {/* Academic Timeline */}
          <div className="space-y-4">
            <h3 className="text-xl font-semibold text-gray-900">Academic Timeline</h3>
            
            {Object.entries(studentJourney.terms || {})
              .sort(([a], [b]) => a.localeCompare(b))
              .map(([termId, termData]) => (
                <div key={termId} className="border border-gray-200 rounded-lg p-6">
                  {/* Term Header */}
                  <div className="flex justify-between items-center mb-4">
                    <div>
                      <h4 className="text-lg font-semibold text-gray-900">
                        {termId.replace(/_/g, ' ')}
                      </h4>
                      <p className="text-sm text-gray-600">
                        {termData.courses?.length || 0} courses completed
                      </p>
                    </div>
                    <div className="text-right">
                      <div className="text-sm text-gray-600">Term GPA</div>
                      <div className="text-lg font-semibold text-green-600">
                        {calculateTermGPA(termData.courses || [])}
                      </div>
                    </div>
                  </div>

                  {/* Course Grid */}
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {(termData.courses || []).map((course, index) => (
                      <div key={index} className="border border-gray-100 rounded-md p-4 hover:shadow-md transition-shadow">
                        <div className="flex justify-between items-start mb-2">
                          <div>
                            <div className="font-semibold text-gray-900">{course.course_id}</div>
                            <div className="text-sm text-gray-600">{course.course_name}</div>
                          </div>
                          <div className={`px-2 py-1 rounded text-white text-xs font-semibold ${getGradeColor(course.grade)}`}>
                            {course.grade}
                          </div>
                        </div>
                        <div className="text-xs text-gray-500 space-y-1">
                          <div>Credits: {course.credits}</div>
                          <div>Instructor: {course.instructor}</div>
                          <div>Completed: {formatDate(course.completed_at)}</div>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              ))}
          </div>

          {/* Academic Progress Chart */}
          <div className="bg-gray-50 rounded-lg p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">GPA Progression</h3>
            <div className="flex items-end space-x-2 h-32">
              {Object.entries(studentJourney.terms || {})
                .sort(([a], [b]) => a.localeCompare(b))
                .map(([termId, termData], index) => {
                  const gpa = parseFloat(calculateTermGPA(termData.courses || []));
                  const height = (gpa / 4.0) * 100;
                  return (
                    <div key={termId} className="flex flex-col items-center flex-1">
                      <div 
                        className="bg-blue-500 rounded-t w-full transition-all duration-300 hover:bg-blue-600"
                        style={{ height: `${height}%` }}
                        title={`${termId}: ${gpa} GPA`}
                      ></div>
                      <div className="text-xs text-gray-600 mt-1 text-center">
                        {termId.split('_')[0]}
                      </div>
                      <div className="text-xs font-semibold text-gray-800">
                        {gpa}
                      </div>
                    </div>
                  );
                })}
            </div>
          </div>
        </div>
      )}

      {!studentJourney && !loading && selectedStudentId && (
        <div className="text-center py-8 text-gray-500">
          No journey data available for selected student
        </div>
      )}
    </div>
  );
};

export default StudentJourneyVisualization;