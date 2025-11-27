import { useState } from "react";
import { apiService } from "@/lib/api";
import { JourneyReceipt, TermReceipt, VerificationResult } from "./types";

export function useReceiptVerifier() {
  const [receipt, setReceipt] = useState<JourneyReceipt | null>(null);
  const [expandedTerms, setExpandedTerms] = useState<Set<string>>(new Set());
  const [expandedCourses, setExpandedCourses] = useState<Set<string>>(new Set());
  const [verificationResults, setVerificationResults] = useState<Record<string, VerificationResult>>({});
  const [isVerifying, setIsVerifying] = useState(false);
  const [error, setError] = useState<string>("");

  // Selective disclosure state
  const [selectedTerms, setSelectedTerms] = useState<Set<string>>(new Set());
  const [selectedCourses, setSelectedCourses] = useState<Record<string, Set<string>>>({});
  const [isSelectiveMode, setIsSelectiveMode] = useState(false);

  const handleFileUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (e) => {
      try {
        const json = JSON.parse(e.target?.result as string);
        setReceipt(json);
        setError("");
        setExpandedTerms(new Set(Object.keys(json.term_receipts || {})));
      } catch (err) {
        setError("Invalid JSON file. Please upload a valid receipt.");
      }
    };
    reader.readAsText(file);
  };

  const toggleTerm = (termId: string) => {
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

  const toggleCourse = (termId: string, courseId: string) => {
    const key = `${termId}-${courseId}`;
    setExpandedCourses((prev) => {
      const next = new Set(prev);
      if (next.has(key)) {
        next.delete(key);
      } else {
        next.add(key);
      }
      return next;
    });
  };

  const initializeSelectiveMode = () => {
    if (!receipt) return;

    const allTerms = new Set(Object.keys(receipt.term_receipts));
    setSelectedTerms(allTerms);

    const allCourses: Record<string, Set<string>> = {};
    Object.entries(receipt.term_receipts).forEach(([termId, termData]) => {
      const courseIds = new Set(
        termData.receipt.revealed_courses.map((c) => c.course_id || "")
      );
      allCourses[termId] = courseIds;
    });
    setSelectedCourses(allCourses);
    setIsSelectiveMode(true);
  };

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

  const generateFilteredReceipt = (): JourneyReceipt | null => {
    if (!receipt) return null;

    const filteredTermReceipts: Record<string, TermReceipt> = {};

    selectedTerms.forEach((termId) => {
      const termData = receipt.term_receipts[termId];
      const selectedCoursesForTerm = selectedCourses[termId] || new Set();

      const filteredCourses = termData.receipt.revealed_courses.filter((course) =>
        selectedCoursesForTerm.has(course.course_id || "")
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

  const getSelectedCounts = () => {
    const totalTerms = selectedTerms.size;
    const totalCourses = Array.from(selectedTerms).reduce((sum, termId) => {
      return sum + (selectedCourses[termId]?.size || 0);
    }, 0);
    return { totalTerms, totalCourses };
  };

  const verifyJourney = async () => {
    if (!receipt) return;

    setIsVerifying(true);
    setError("");

    try {
      const result = await apiService.verifyReceiptIPA(receipt);

      setVerificationResults({
        overall: {
          verified: result.status === "success",
          ipa_verified: result.verified_courses === result.total_courses,
          blockchain_anchored: true,
          details: `${result.verified_courses}/${result.total_courses} courses verified`,
        },
        ...result.term_results,
      });
    } catch (err: any) {
      setError(err.message || "Verification failed");
      setVerificationResults({
        overall: {
          verified: false,
          error: err.message,
        },
      });
    } finally {
      setIsVerifying(false);
    }
  };

  const resetReceipt = () => {
    setReceipt(null);
    setVerificationResults({});
    setSelectedTerms(new Set());
    setSelectedCourses({});
    setIsSelectiveMode(false);
  };

  return {
    receipt,
    expandedTerms,
    expandedCourses,
    verificationResults,
    isVerifying,
    error,
    selectedTerms,
    selectedCourses,
    isSelectiveMode,
    handleFileUpload,
    toggleTerm,
    toggleCourse,
    initializeSelectiveMode,
    toggleTermSelection,
    toggleCourseSelection,
    downloadFilteredReceipt,
    getSelectedCounts,
    verifyJourney,
    resetReceipt,
    setIsSelectiveMode,
  };
}
