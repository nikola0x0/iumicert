"use client";

import { UploadSection } from "./verifier/UploadSection";
import { ReceiptHeader } from "./verifier/ReceiptHeader";
import { JourneySummary } from "./verifier/JourneySummary";
import { TermCard } from "./verifier/TermCard";
import { useReceiptVerifier } from "./verifier/useReceiptVerifier";

export default function ReceiptVerifier() {
  const {
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
  } = useReceiptVerifier();

  // Show upload section if no receipt
  if (!receipt) {
    return <UploadSection onFileUpload={handleFileUpload} error={error} />;
  }

  const overallResult = verificationResults.overall;
  const selectedCounts = getSelectedCounts();

  // Sort terms chronologically
  const sortedTermEntries = Object.entries(receipt.term_receipts).sort(
    ([termIdA], [termIdB]) => {
      const extractOrder = (id: string) => {
        const yearMatch = id.match(/(\d{4})/);
        const year = yearMatch ? parseInt(yearMatch[1]) : 0;
        const semester = id.includes("Semester_1")
          ? 1
          : id.includes("Semester_2")
          ? 2
          : 3;
        return year * 10 + semester;
      };
      return extractOrder(termIdA) - extractOrder(termIdB);
    }
  );

  return (
    <div className="space-y-6">
      <ReceiptHeader
        receipt={receipt}
        overallResult={overallResult}
        isVerifying={isVerifying}
        isSelectiveMode={isSelectiveMode}
        selectedCourseCount={selectedCounts.totalCourses}
        error={error}
        onUploadNew={resetReceipt}
        onVerify={verifyJourney}
        onInitializeSelective={initializeSelectiveMode}
        onDownloadSelective={downloadFilteredReceipt}
        onCancelSelective={() => setIsSelectiveMode(false)}
      />

      <JourneySummary receipt={receipt} />

      {/* Terms List */}
      <div className="space-y-4">
        {sortedTermEntries.map(([termId, termData]) => (
          <TermCard
            key={termId}
            termId={termId}
            termData={termData}
            termResult={verificationResults[termId]}
            isExpanded={expandedTerms.has(termId)}
            isSelectiveMode={isSelectiveMode}
            isSelected={selectedTerms.has(termId)}
            selectedCourses={selectedCourses[termId] || new Set()}
            expandedCourses={expandedCourses}
            onToggle={() => toggleTerm(termId)}
            onToggleSelection={() => toggleTermSelection(termId)}
            onToggleCourse={(courseId) => toggleCourse(termId, courseId)}
            onToggleCourseSelection={(courseId) =>
              toggleCourseSelection(termId, courseId)
            }
          />
        ))}
      </div>
    </div>
  );
}
