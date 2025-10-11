import { Eye, BarChart3 } from "lucide-react";
import { AggregatedJourneyProof } from "../types/proofs";

interface AggregatedJourneyViewProps {
  proofData: AggregatedJourneyProof;
  onTermSelect: (term: string) => void;
}

export default function AggregatedJourneyView({
  proofData,
  onTermSelect,
}: AggregatedJourneyViewProps) {
  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 h-full">
      {/* Left Column - Student Info */}
      <div className="space-y-4">
        <div className="glass-effect rounded-xl p-6 text-center">
          <div className="w-16 h-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center mx-auto mb-4">
            <span className="text-white font-bold text-xl">🎓</span>
          </div>
          <h3 className="text-xl font-bold text-white font-space-grotesk mb-1">
            {proofData.student_info.student_name}
          </h3>
          <p className="text-purple-200 text-sm font-inter">
            {proofData.student_info.program}
          </p>
          <p className="text-purple-300 text-xs font-inter">
            {proofData.institutional_verification.institution}
          </p>
        </div>

        {/* Journey Summary */}
        <div className="glass-effect rounded-xl p-6">
          <h4 className="text-white font-bold text-sm mb-4 font-space-grotesk">
            Journey Overview
          </h4>
          <div className="grid grid-cols-2 gap-3">
            <div className="bg-blue-500/20 rounded-lg p-3 text-center">
              <div className="text-2xl font-bold text-blue-300">
                {proofData.journey_summary.total_terms}
              </div>
              <div className="text-blue-200 text-xs font-inter">Terms</div>
            </div>
            <div className="bg-green-500/20 rounded-lg p-3 text-center">
              <div className="text-2xl font-bold text-green-300">
                {proofData.journey_summary.total_courses}
              </div>
              <div className="text-green-200 text-xs font-inter">Courses</div>
            </div>
            <div className="bg-purple-500/20 rounded-lg p-3 text-center">
              <div className="text-2xl font-bold text-purple-300">
                {proofData.journey_summary.total_credits}
              </div>
              <div className="text-purple-200 text-xs font-inter">
                Credits
              </div>
            </div>
            <div className="bg-yellow-500/20 rounded-lg p-3 text-center">
              <div className="text-2xl font-bold text-yellow-300">
                {proofData.journey_summary.cumulative_gpa}
              </div>
              <div className="text-yellow-200 text-xs font-inter">CGPA</div>
            </div>
          </div>
        </div>
      </div>

      {/* Right Columns - Academic Timeline */}
      <div className="lg:col-span-2">
        <div className="glass-effect rounded-xl p-6 h-full">
          <h4 className="text-white font-bold text-lg mb-4 flex items-center font-space-grotesk">
            <BarChart3 className="w-5 h-5 mr-2" />
            Academic Timeline
          </h4>
          <div className="grid grid-cols-1 xl:grid-cols-2 gap-4 max-h-96 overflow-y-auto">
            {proofData.academic_terms.map((term, index) => (
              <div
                key={term.term}
                onClick={() => onTermSelect(term.term)}
                className="bg-white/5 hover:bg-white/10 rounded-lg p-4 border border-white/10 hover:border-blue-300/30 cursor-pointer transition-all duration-200 hover:scale-[1.02] group"
              >
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center gap-3">
                    <div className="w-8 h-8 bg-gradient-to-r from-blue-500 to-purple-500 rounded-full flex items-center justify-center text-white font-bold text-sm">
                      {index + 1}
                    </div>
                    <div>
                      <h5 className="text-white font-semibold text-sm font-space-grotesk">
                        {term.term}
                      </h5>
                      <p className="text-purple-200 text-xs font-inter">
                        {term.courses.length} courses • {term.total_credits}{" "}
                        credits
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    {term.term_gpa && (
                      <div className="bg-white/10 px-2 py-1 rounded-full">
                        <span className="text-white font-medium text-xs">
                          {term.term_gpa}
                        </span>
                      </div>
                    )}
                    <Eye className="w-4 h-4 text-white/60 group-hover:text-blue-300 transition-colors duration-200" />
                  </div>
                </div>

                <div className="flex flex-wrap gap-1">
                  {term.courses.slice(0, 4).map((course) => (
                    <div
                      key={course.course_code}
                      className="bg-white/10 px-2 py-1 rounded text-xs text-purple-200 font-mono"
                    >
                      {course.course_code}
                    </div>
                  ))}
                  {term.courses.length > 4 && (
                    <div className="bg-white/10 px-2 py-1 rounded text-xs text-purple-200">
                      +{term.courses.length - 4}
                    </div>
                  )}
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}