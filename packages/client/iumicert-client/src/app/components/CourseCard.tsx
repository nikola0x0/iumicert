import { CourseCompletion } from "../types/proofs";

interface CourseCardProps {
  course: CourseCompletion;
  index?: number;
}

export default function CourseCard({ course }: CourseCardProps) {
  const getGradeColor = (grade: string) => {
    if (["A+", "A", "A-"].includes(grade)) return "bg-green-500/20 text-green-300";
    if (["B+", "B", "B-"].includes(grade)) return "bg-blue-500/20 text-blue-300";
    return "bg-yellow-500/20 text-yellow-300";
  };

  return (
    <div className="bg-white/5 rounded-lg p-4 border border-white/10 hover:border-white/20 transition-colors duration-200">
      {/* Header */}
      <div className="flex items-center justify-between mb-3">
        <span className="text-sm font-mono text-blue-300 bg-blue-500/10 px-2 py-1 rounded">
          {course.course_code}
        </span>
        <span className={`px-2 py-1 rounded-full text-xs font-bold ${getGradeColor(course.grade)}`}>
          {course.grade}
        </span>
      </div>
      
      {/* Title */}
      <h5 className="text-white font-medium text-sm mb-3 font-space-grotesk">
        {course.course_name}
      </h5>
      
      {/* Details */}
      <div className="text-xs text-purple-200 space-y-1 font-inter">
        <div className="flex items-center gap-2">
          <span>👨‍🏫</span>
          <span>{course.instructor}</span>
        </div>
        <div className="flex justify-between">
          <span>📅 {course.completion_date}</span>
          <span>🎓 {course.credits}cr</span>
        </div>
      </div>
    </div>
  );
}