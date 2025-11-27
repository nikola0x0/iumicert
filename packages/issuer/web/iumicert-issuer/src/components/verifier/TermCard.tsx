import { TermReceipt, VerificationResult, Course } from "./types";
import { CourseCard } from "./CourseCard";

interface TermCardProps {
  termId: string;
  termData: TermReceipt;
  termResult?: VerificationResult;
  isExpanded: boolean;
  isSelectiveMode: boolean;
  isSelected: boolean;
  selectedCourses: Set<string>;
  expandedCourses: Set<string>;
  onToggle: () => void;
  onToggleSelection: () => void;
  onToggleCourse: (courseId: string) => void;
  onToggleCourseSelection: (courseId: string) => void;
}

export function TermCard({
  termId,
  termData,
  termResult,
  isExpanded,
  isSelectiveMode,
  isSelected,
  selectedCourses,
  expandedCourses,
  onToggle,
  onToggleSelection,
  onToggleCourse,
  onToggleCourseSelection,
}: TermCardProps) {
  return (
    <div className="bg-white rounded-2xl shadow-lg border border-gray-100 overflow-hidden">
      {/* Term Header */}
      <div
        className="p-6 cursor-pointer hover:bg-gray-50/50 transition-colors"
        onClick={onToggle}
      >
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            {isSelectiveMode && (
              <input
                type="checkbox"
                checked={isSelected}
                onChange={(e) => {
                  e.stopPropagation();
                  onToggleSelection();
                }}
                onClick={(e) => e.stopPropagation()}
                className="w-5 h-5 rounded border-gray-300 text-purple-600 focus:ring-purple-500 cursor-pointer"
              />
            )}
            <div
              className={`w-10 h-10 rounded-xl flex items-center justify-center ${
                isExpanded ? "bg-blue-100" : "bg-gray-100"
              }`}
            >
              <svg
                className={`w-5 h-5 ${
                  isExpanded ? "text-blue-600" : "text-gray-600"
                }`}
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                {isExpanded ? (
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M5 19a2 2 0 01-2-2V7a2 2 0 012-2h4l2 2h4a2 2 0 012 2v1M5 19h14a2 2 0 002-2v-5a2 2 0 00-2-2H9a2 2 0 00-2 2v5a2 2 0 01-2 2z"
                  />
                ) : (
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"
                  />
                )}
              </svg>
            </div>
            <div>
              <h3 className="text-xl font-bold text-slate-900">
                {termId.replace(/_/g, " ")}
              </h3>
              <p className="text-sm text-gray-600">
                {isSelectiveMode && selectedCourses.size > 0
                  ? `${selectedCourses.size} selected / `
                  : ""}
                {termData.revealed_courses} of {termData.total_courses} courses
              </p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            {termResult && (
              <div
                className={`px-4 py-1.5 rounded-full text-sm font-semibold border ${
                  termResult.status === "completed" &&
                  (termResult.courses_verified || 0) > 0
                    ? "bg-green-100 text-green-800 border-green-200"
                    : "bg-red-100 text-red-800 border-red-200"
                }`}
              >
                {termResult.status === "completed" &&
                (termResult.courses_verified || 0) > 0
                  ? `Verified (${termResult.courses_verified || 0} courses)`
                  : "Failed"}
              </div>
            )}
            <svg
              className="w-5 h-5 text-gray-400"
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

      {/* Term Content */}
      {isExpanded && (
        <div className="border-t border-gray-200 bg-gray-50 p-6">
          {/* Blockchain Verification Status */}
          {termResult?.blockchain_verified && (
            <div className="mb-4 p-4 bg-green-50 rounded-xl border border-green-200">
              <div className="flex items-start gap-3">
                <div className="w-10 h-10 bg-green-100 rounded-xl flex items-center justify-center flex-shrink-0">
                  <svg
                    className="w-5 h-5 text-green-600"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M13 10V3L4 14h7v7l9-11h-7z"
                    />
                  </svg>
                </div>
                <div className="flex-1">
                  <h4 className="font-semibold text-green-800 mb-2">
                    Blockchain Verified
                  </h4>
                  <div className="space-y-2 text-sm">
                    <div>
                      <span className="text-green-700">Transaction:</span>
                      <a
                        href={`https://sepolia.etherscan.io/tx/${termResult.blockchain_tx_hash}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="ml-2 font-mono text-xs text-blue-600 hover:text-blue-800 underline break-all"
                      >
                        {termResult.blockchain_tx_hash}
                      </a>
                    </div>
                    <div>
                      <span className="text-green-700">Published:</span>
                      <span className="ml-2 text-green-800">
                        {termResult.blockchain_published_at
                          ? new Date(
                              termResult.blockchain_published_at
                            ).toLocaleString()
                          : "N/A"}
                      </span>
                    </div>
                    <div>
                      <span className="text-green-700">Block:</span>
                      <span className="ml-2 font-mono text-green-800">
                        {termResult.blockchain_block || "N/A"}
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          )}

          {/* Term Metadata */}
          <div className="mb-4 p-4 bg-white rounded-xl border border-gray-200">
            <div className="grid grid-cols-2 gap-4 text-sm">
              <div>
                <span className="text-gray-600 font-medium">Verkle Root:</span>
                <p className="font-mono text-xs mt-1 break-all text-slate-900 bg-gray-50 p-2 rounded">
                  {termData.verkle_root}
                </p>
              </div>
              <div>
                <span className="text-gray-600 font-medium">Generated:</span>
                <p className="text-slate-900 mt-1">
                  {new Date(termData.generated_at).toLocaleString()}
                </p>
              </div>
            </div>
          </div>

          {/* Courses */}
          <div className="space-y-2">
            <h4 className="font-semibold text-slate-900 mb-3 flex items-center gap-2">
              <svg
                className="w-5 h-5 text-blue-600"
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
              Courses
            </h4>
            {termData.receipt.revealed_courses.map((course: Course, idx: number) => {
              const courseId = course.course_id || idx.toString();
              const courseKey = `${termId}-${courseId}`;

              return (
                <CourseCard
                  key={courseKey}
                  course={course}
                  courseId={courseId}
                  termId={termId}
                  isExpanded={expandedCourses.has(courseKey)}
                  isSelectiveMode={isSelectiveMode}
                  isTermSelected={isSelected}
                  isSelected={selectedCourses.has(courseId)}
                  proof={termData.receipt.course_proofs[courseId]}
                  proofType={termData.receipt.proof_type}
                  verificationPath={termData.receipt.verification_path}
                  onToggle={() => onToggleCourse(courseId)}
                  onToggleSelection={() => onToggleCourseSelection(courseId)}
                />
              );
            })}
          </div>
        </div>
      )}
    </div>
  );
}
