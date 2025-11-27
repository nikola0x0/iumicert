"use client";

import { useState } from "react";
import DashboardLayout from "../components/DashboardLayout";
import FileUploaderWrapper from "../components/FileUploaderWrapper";
import { Lock, FileText, Download, ChevronDown, ChevronUp } from "lucide-react";

interface Course {
  course_id?: string;
  course_name?: string;
  grade: string;
  credits: number;
}

interface TermReceipt {
  term_id: string;
  student_id: string;
  receipt: {
    student_id: string;
    term_id: string;
    revealed_courses: Course[];
    verkle_root: string;
    course_proofs: Record<string, any>;
    proof_type: string;
    selective_disclosure: boolean;
    verification_path: string;
    timestamp: string;
  };
  verkle_root: string;
  revealed_courses: number;
  total_courses: number;
  generated_at: string;
}

interface JourneyReceipt {
  student_id: string;
  receipt_type: {
    selective_disclosure: boolean;
    specific_courses: boolean;
    specific_terms: boolean;
  };
  generation_timestamp: string;
  terms_included: string[];
  courses_filter: string[];
  term_receipts: Record<string, TermReceipt>;
  blockchain_ready: boolean;
}

export default function SelectiveDisclosurePage() {
  const [receipt, setReceipt] = useState<JourneyReceipt | null>(null);
  const [selectedTerms, setSelectedTerms] = useState<Set<string>>(new Set());
  const [selectedCourses, setSelectedCourses] = useState<Record<string, Set<string>>>({});
  const [expandedTerms, setExpandedTerms] = useState<Set<string>>(new Set());
  const [error, setError] = useState<string>("");

  // Handle file upload
  const handleFileChange = async (file: File) => {
    try {
      const text = await file.text();
      const json = JSON.parse(text);
      setReceipt(json);
      setError("");

      // Initialize with all terms and courses selected
      const allTerms = new Set(Object.keys(json.term_receipts || {}));
      setSelectedTerms(allTerms);
      setExpandedTerms(new Set());

      const allCourses: Record<string, Set<string>> = {};
      Object.entries(json.term_receipts).forEach(([termId, termData]: [string, any]) => {
        const courseIds = new Set<string>(
          termData.receipt.revealed_courses.map((c: Course) => c.course_id || "")
        );
        allCourses[termId] = courseIds;
      });
      setSelectedCourses(allCourses);
    } catch (err) {
      setError("Invalid JSON file. Please upload a valid receipt.");
      setReceipt(null);
    }
  };

  const handleTypeError = () => {
    setError("Please upload a JSON file containing the receipt.");
  };

  // Toggle term selection
  const toggleTermSelection = (termId: string) => {
    setSelectedTerms((prev) => {
      const next = new Set(prev);
      if (next.has(termId)) {
        next.delete(termId);
        setSelectedCourses((prevCourses) => {
          const nextCourses = { ...prevCourses };
          delete nextCourses[termId];
          return nextCourses;
        });
      } else {
        next.add(termId);
        if (receipt) {
          const termData = receipt.term_receipts[termId];
          setSelectedCourses((prevCourses) => ({
            ...prevCourses,
            [termId]: new Set(
              termData.receipt.revealed_courses.map((c) => c.course_id || "")
            ),
          }));
        }
      }
      return next;
    });
  };

  // Toggle course selection
  const toggleCourseSelection = (termId: string, courseId: string) => {
    setSelectedCourses((prev) => {
      const termCourses = new Set(prev[termId] || []);
      if (termCourses.has(courseId)) {
        termCourses.delete(courseId);
      } else {
        termCourses.add(courseId);
      }
      return { ...prev, [termId]: termCourses };
    });
  };

  // Toggle term expansion
  const toggleTermExpansion = (termId: string) => {
    setExpandedTerms((prev) => {
      const next = new Set(prev);
      if (next.has(termId)) {
        next.delete(termId);
      } else {
        next.add(termId);
      }
      return next;
    });
  };

  // Generate filtered receipt
  const generateFilteredReceipt = () => {
    if (!receipt) return null;

    const filteredTermReceipts: Record<string, TermReceipt> = {};

    selectedTerms.forEach((termId) => {
      const termData = receipt.term_receipts[termId];
      const selectedCoursesForTerm = selectedCourses[termId] || new Set();

      const filteredCourses = termData.receipt.revealed_courses.filter(
        (course) => selectedCoursesForTerm.has(course.course_id || "")
      );

      const filteredProofs: Record<string, any> = {};
      selectedCoursesForTerm.forEach((courseId) => {
        if (termData.receipt.course_proofs[courseId]) {
          filteredProofs[courseId] = termData.receipt.course_proofs[courseId];
        }
      });

      filteredTermReceipts[termId] = {
        ...termData,
        revealed_courses: filteredCourses.length,
        receipt: {
          ...termData.receipt,
          revealed_courses: filteredCourses,
          course_proofs: filteredProofs,
          selective_disclosure: true,
        },
      };
    });

    return {
      ...receipt,
      receipt_type: {
        selective_disclosure: true,
        specific_courses: true,
        specific_terms: true,
      },
      terms_included: Array.from(selectedTerms),
      term_receipts: filteredTermReceipts,
      generation_timestamp: new Date().toISOString(),
    };
  };

  // Download filtered receipt
  const downloadFilteredReceipt = () => {
    const filteredReceipt = generateFilteredReceipt();
    if (!filteredReceipt) return;

    const blob = new Blob([JSON.stringify(filteredReceipt, null, 2)], {
      type: "application/json",
    });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `${receipt?.student_id}_selective_receipt.json`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  // Get counts
  const getSelectedCounts = () => {
    const totalTerms = selectedTerms.size;
    const totalCourses = Array.from(selectedTerms).reduce((sum, termId) => {
      return sum + (selectedCourses[termId]?.size || 0);
    }, 0);
    return { totalTerms, totalCourses };
  };

  const { totalTerms, totalCourses } = getSelectedCounts();

  return (
    <DashboardLayout activeSection="selective">
      <div className="h-full flex flex-col">
        {!receipt ? (
          /* Upload View */
          <div className="h-full flex items-center justify-center p-4">
            <div className="glass-effect rounded-xl p-8 max-w-lg w-full">
              {/* Header */}
              <div className="text-center mb-8">
                <div className="w-16 h-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center mx-auto mb-4">
                  <Lock className="w-8 h-8 text-white" />
                </div>
                <h1 className="text-2xl font-bold text-white font-space-grotesk mb-2">
                  Selective Disclosure
                </h1>
                <p className="text-purple-200 font-inter text-sm">
                  Create a filtered receipt showing only selected courses
                </p>
              </div>

              {/* File Upload Area */}
              <div className="mb-4">
                <FileUploaderWrapper
                  handleChange={handleFileChange}
                  name="receipt"
                  types={["JSON"]}
                  onTypeError={handleTypeError}
                  maxSize={10}
                  classes="w-full file-uploader-custom hover:cursor-pointer"
                >
                  <div className="w-full rounded-xl border-2 border-dashed border-white/30 bg-white/5 hover:border-blue-400 hover:bg-blue-500/20 transition-all duration-300 p-6 text-center">
                    <div className="mb-3">
                      <FileText className="mx-auto h-10 w-10 text-white/60" />
                    </div>
                    <p className="text-white/90 font-medium text-sm mb-1 font-space-grotesk">
                      Drag & drop your full receipt
                    </p>
                    <p className="text-white/60 text-xs font-inter">
                      or click to browse files
                    </p>
                  </div>
                </FileUploaderWrapper>
              </div>

              {error && (
                <div className="mt-4 p-3 bg-red-500/20 border border-red-500/30 rounded-xl text-red-200 text-sm font-inter">
                  {error}
                </div>
              )}

              {/* Info */}
              <div className="mt-6 p-4 bg-blue-500/10 border border-blue-500/20 rounded-xl">
                <p className="text-sm text-blue-200 font-inter mb-2 font-semibold">How it works:</p>
                <ul className="text-xs text-blue-200/80 font-inter space-y-1 list-disc list-inside">
                  <li>Upload your complete academic receipt</li>
                  <li>Select which terms and courses to reveal</li>
                  <li>Download filtered receipt with proofs</li>
                  <li>Still verifies cryptographically!</li>
                </ul>
              </div>
            </div>
          </div>
        ) : (
          /* Selection View */
          <div className="h-full flex flex-col p-6 overflow-auto">
            {/* Header Card */}
            <div className="glass-effect rounded-xl p-6 mb-6">
              <div className="flex items-center justify-between mb-4">
                <div>
                  <h2 className="text-xl font-bold text-white font-space-grotesk">
                    Student: <span className="text-blue-300">{receipt.student_id}</span>
                  </h2>
                  <p className="text-purple-200 font-inter text-sm mt-1">
                    Select terms and courses to include in filtered receipt
                  </p>
                </div>
                <button
                  onClick={() => {
                    setReceipt(null);
                    setError("");
                  }}
                  className="px-4 py-2 bg-white/10 hover:bg-white/20 text-white rounded-lg transition-colors duration-200 font-inter text-sm"
                >
                  Upload New
                </button>
              </div>

              {/* Stats */}
              <div className="grid grid-cols-3 gap-4 mb-4">
                <div className="bg-blue-500/20 p-3 rounded-lg border border-blue-500/30">
                  <div className="text-2xl font-bold text-blue-300 font-space-grotesk">{totalTerms}</div>
                  <div className="text-xs text-white/70 font-inter">Terms Selected</div>
                </div>
                <div className="bg-purple-500/20 p-3 rounded-lg border border-purple-500/30">
                  <div className="text-2xl font-bold text-purple-300 font-space-grotesk">{totalCourses}</div>
                  <div className="text-xs text-white/70 font-inter">Courses Selected</div>
                </div>
                <div className="bg-white/5 p-3 rounded-lg border border-white/20">
                  <div className="text-sm font-bold text-white/90 font-space-grotesk">
                    {Object.keys(receipt.term_receipts).length} Terms
                  </div>
                  <div className="text-xs text-white/70 font-inter">
                    {Object.values(receipt.term_receipts).reduce((sum, term) => sum + term.total_courses, 0)} Courses Total
                  </div>
                </div>
              </div>

              {/* Download Button */}
              <button
                onClick={downloadFilteredReceipt}
                disabled={totalCourses === 0}
                className="w-full bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700
                         disabled:from-gray-600 disabled:to-gray-700 disabled:cursor-not-allowed
                         text-white font-bold py-3 px-6 rounded-xl transition-all duration-300
                         hover:scale-105 hover:shadow-xl flex items-center justify-center gap-2 font-space-grotesk"
              >
                <Download className="w-5 h-5" />
                Download Filtered Receipt ({totalCourses} courses)
              </button>
            </div>

            {/* Terms List */}
            <div className="space-y-3">
              {Object.entries(receipt.term_receipts)
                .sort(([termIdA], [termIdB]) => {
                  const extractOrder = (id: string) => {
                    const yearMatch = id.match(/(\d{4})/);
                    const year = yearMatch ? parseInt(yearMatch[1]) : 0;
                    const semester = id.includes("Semester_1") ? 1 : id.includes("Semester_2") ? 2 : 3;
                    return year * 10 + semester;
                  };
                  return extractOrder(termIdA) - extractOrder(termIdB);
                })
                .map(([termId, termData]) => {
                  const isExpanded = expandedTerms.has(termId);
                  const isSelected = selectedTerms.has(termId);
                  const selectedCount = selectedCourses[termId]?.size || 0;

                  return (
                    <div key={termId} className="glass-effect rounded-xl overflow-hidden">
                      {/* Term Header */}
                      <div className="p-4">
                        <div className="flex items-center justify-between">
                          <div className="flex items-center gap-3">
                            <input
                              type="checkbox"
                              checked={isSelected}
                              onChange={() => toggleTermSelection(termId)}
                              className="w-5 h-5 rounded border-white/30 bg-white/10 text-blue-500 focus:ring-blue-400 cursor-pointer"
                            />
                            <div>
                              <h3 className="text-lg font-bold text-white font-space-grotesk">
                                {termId.replace(/_/g, " ")}
                              </h3>
                              <p className="text-xs text-white/60 font-inter">
                                {isSelected ? `${selectedCount} selected / ` : ""}
                                {termData.total_courses} courses
                              </p>
                            </div>
                          </div>
                          <button
                            onClick={() => toggleTermExpansion(termId)}
                            disabled={!isSelected}
                            className={`px-3 py-1.5 rounded-lg font-inter text-sm transition-colors ${
                              !isSelected
                                ? "bg-white/5 text-white/30 cursor-not-allowed"
                                : "bg-blue-500/20 text-blue-300 hover:bg-blue-500/30"
                            }`}
                          >
                            {isExpanded ? (
                              <ChevronUp className="w-4 h-4" />
                            ) : (
                              <ChevronDown className="w-4 h-4" />
                            )}
                          </button>
                        </div>
                      </div>

                      {/* Courses List */}
                      {isExpanded && isSelected && (
                        <div className="border-t border-white/10 bg-white/5 p-4">
                          <div className="space-y-2">
                            {termData.receipt.revealed_courses.map((course, idx) => (
                              <div
                                key={course.course_id || idx}
                                className="bg-white/5 rounded-lg border border-white/10 p-3 hover:bg-white/10 transition-colors"
                              >
                                <div className="flex items-center justify-between">
                                  <div className="flex items-center gap-3">
                                    <input
                                      type="checkbox"
                                      checked={selectedCourses[termId]?.has(course.course_id || "")}
                                      onChange={() => toggleCourseSelection(termId, course.course_id || "")}
                                      className="w-4 h-4 rounded border-white/30 bg-white/10 text-blue-500 focus:ring-blue-400 cursor-pointer"
                                    />
                                    <div>
                                      <p className="font-semibold text-white font-space-grotesk text-sm">
                                        {course.course_id || `Course ${idx + 1}`}
                                      </p>
                                      <p className="text-xs text-white/60 font-inter">
                                        {course.course_name || "N/A"}
                                      </p>
                                    </div>
                                  </div>
                                  <div className="flex items-center gap-2">
                                    <span className={`px-2 py-1 rounded-full text-xs font-semibold font-inter ${
                                      course.grade === "A" || course.grade === "A+"
                                        ? "bg-green-500/20 text-green-300 border border-green-500/30"
                                        : course.grade === "B" || course.grade === "B+"
                                        ? "bg-blue-500/20 text-blue-300 border border-blue-500/30"
                                        : "bg-amber-500/20 text-amber-300 border border-amber-500/30"
                                    }`}>
                                      {course.grade}
                                    </span>
                                    <span className="text-xs text-white/60 font-inter">
                                      {course.credits} cr
                                    </span>
                                  </div>
                                </div>
                              </div>
                            ))}
                          </div>
                        </div>
                      )}
                    </div>
                  );
                })}
            </div>
          </div>
        )}
      </div>
    </DashboardLayout>
  );
}
