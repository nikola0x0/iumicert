import { Course } from "./types";

interface CourseCardProps {
  course: Course;
  courseId: string;
  termId: string;
  isExpanded: boolean;
  isSelectiveMode: boolean;
  isTermSelected: boolean;
  isSelected: boolean;
  proof: any;
  proofType: string;
  verificationPath: string;
  onToggle: () => void;
  onToggleSelection: () => void;
}

export function CourseCard({
  course,
  courseId,
  isExpanded,
  isSelectiveMode,
  isTermSelected,
  isSelected,
  proof,
  proofType,
  verificationPath,
  onToggle,
  onToggleSelection,
}: CourseCardProps) {
  const getGradeColor = (grade: string) => {
    if (grade === "A" || grade === "A+")
      return "bg-green-100 text-green-800 border-green-200";
    if (grade === "B" || grade === "B+")
      return "bg-blue-100 text-blue-800 border-blue-200";
    if (grade === "C" || grade === "C+")
      return "bg-amber-100 text-amber-800 border-amber-200";
    return "bg-gray-100 text-gray-800 border-gray-200";
  };

  return (
    <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
      {/* Course Header */}
      <div
        className="p-4 cursor-pointer hover:bg-gray-50/50 transition-colors"
        onClick={onToggle}
      >
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            {isSelectiveMode && isTermSelected && (
              <input
                type="checkbox"
                checked={isSelected}
                onChange={(e) => {
                  e.stopPropagation();
                  onToggleSelection();
                }}
                onClick={(e) => e.stopPropagation()}
                className="w-4 h-4 rounded border-gray-300 text-purple-600 focus:ring-purple-500 cursor-pointer"
              />
            )}
            <div
              className={`w-8 h-8 rounded-lg flex items-center justify-center ${
                isExpanded ? "bg-blue-100" : "bg-gray-100"
              }`}
            >
              <svg
                className={`w-4 h-4 ${
                  isExpanded ? "text-blue-600" : "text-gray-600"
                }`}
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"
                />
              </svg>
            </div>
            <div>
              <p className="font-semibold text-gray-800">
                {course.course_id || `Course ${courseId}`}
              </p>
              <p className="text-sm text-gray-600">
                {course.course_name || "N/A"}
              </p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            <span
              className={`px-3 py-1 rounded-full text-sm font-semibold border ${getGradeColor(
                course.grade
              )}`}
            >
              Grade: {course.grade}
            </span>
            <span className="text-sm text-gray-600">
              {course.credits} credits
            </span>
            <svg
              className="w-4 h-4 text-gray-400"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d={isExpanded ? "M19 9l-7 7-7-7" : "M9 5l7 7-7 7"}
              />
            </svg>
          </div>
        </div>
      </div>

      {/* Course Details (Proof) */}
      {isExpanded && (
        <div className="border-t border-gray-200 bg-gray-50 p-4">
          <h5 className="text-sm font-semibold text-slate-900 mb-2 flex items-center gap-2">
            <svg
              className="w-4 h-4 text-blue-600"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
              />
            </svg>
            Cryptographic Proof
          </h5>
          <div className="bg-white p-3 rounded-lg border border-gray-200">
            <pre className="text-xs font-mono text-gray-700 overflow-x-auto">
              {JSON.stringify(proof, null, 2)}
            </pre>
          </div>
          <div className="mt-3 flex flex-wrap items-center gap-2">
            <span className="text-xs text-gray-600">Proof Type:</span>
            <span className="px-2 py-1 bg-indigo-100 text-indigo-800 text-xs rounded-lg font-medium">
              {proofType}
            </span>
            <span className="text-xs text-gray-600">Verification Path:</span>
            <span className="px-2 py-1 bg-purple-100 text-purple-800 text-xs rounded-lg font-medium">
              {verificationPath}
            </span>
          </div>
        </div>
      )}
    </div>
  );
}
