import CourseCard from "./CourseCard";
import GlassCard from "./ui/GlassCard";
import StatCard from "./ui/StatCard";
import SectionHeader from "./ui/SectionHeader";
import { Shield, FileText } from "lucide-react";
import {
  ProofData,
  TermData,
  AggregatedJourneyProof,
} from "../types/proofs";

interface SingleTermViewProps {
  proofData: ProofData;
  selectedTerm: string | null;
  onBackToOverview: () => void;
}

export default function SingleTermView({
  proofData,
  selectedTerm,
  onBackToOverview,
}: SingleTermViewProps) {
  let termData: TermData;
  let studentInfo: { student_id?: string; student_name?: string };
  let blockchainInfo: {
    contract_address?: string;
    block_number?: number;
    tree_commitment?: string;
  };
  let institution: string;

  if (
    proofData.type === "individual_term" ||
    proofData.type === "single_term"
  ) {
    termData = {
      term: proofData.claim.term,
      courses: proofData.verkle_proof.claimed_values,
      total_credits: proofData.verkle_proof.claimed_values.reduce(
        (sum, course) => sum + course.credits,
        0
      ),
    };
    studentInfo = { student_id: proofData.claim.student_id };
    blockchainInfo = proofData.blockchain_reference;
    institution = proofData.institutional_verification.institution;
  } else {
    const selectedTermData = (
      proofData as AggregatedJourneyProof
    ).academic_terms.find((t) => t.term === selectedTerm);
    if (!selectedTermData) return null;

    termData = selectedTermData;
    studentInfo = (proofData as AggregatedJourneyProof).student_info;

    const blockchainEntry = (
      proofData as AggregatedJourneyProof
    ).multi_tree_verification_chain.find((v) => v.term === selectedTerm);
    blockchainInfo = blockchainEntry
      ? blockchainEntry.blockchain_deployment
      : {};
    institution = proofData.institutional_verification.institution;
  }

  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 h-full">
      {/* Left Column - Term Summary */}
      <div className="space-y-4">
        <div className="glass-effect rounded-xl p-6">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h3 className="text-xl font-bold text-white font-space-grotesk mb-1">
                {termData.term}
              </h3>
              <p className="text-purple-200 text-sm font-inter">
                {studentInfo.student_name || studentInfo.student_id} •{" "}
                {institution}
              </p>
            </div>
            {proofData.type === "aggregated_journey" && (
              <button
                onClick={onBackToOverview}
                className="px-3 py-1 bg-white/10 hover:bg-white/20 text-white border border-white/20 rounded-lg transition duration-300 text-sm font-inter"
              >
                ← Back
              </button>
            )}
          </div>

          {/* Quick Stats */}
          <div className="grid grid-cols-2 gap-3">
            <StatCard
              label="Courses"
              value={termData.courses.length}
              color="blue"
            />
            <StatCard
              label="Credits"
              value={termData.total_credits}
              color="green"
            />
          </div>

          {termData.term_gpa && (
            <StatCard
              label="Term GPA"
              value={termData.term_gpa}
              color="purple"
              className="mt-3"
            />
          )}
        </div>

        {/* Blockchain Verification */}
        <GlassCard padding="sm">
          <div className="flex items-center gap-3 mb-3">
            <div className="w-8 h-8 bg-green-500/20 rounded-lg flex items-center justify-center">
              <Shield className="w-4 h-4 text-green-400" />
            </div>
            <div>
              <div className="text-white font-semibold text-sm font-space-grotesk">
                Blockchain Verified
              </div>
              <div className="text-green-300 text-xs font-inter">
                Sepolia Testnet
              </div>
            </div>
          </div>
          <div className="text-xs text-white/70 space-y-1 font-mono">
            <div>
              Contract: {blockchainInfo?.contract_address?.substring(0, 10)}
              ...
            </div>
            <div>Block: #{blockchainInfo?.block_number}</div>
            <div>
              Commitment: {blockchainInfo?.tree_commitment?.substring(0, 10)}
              ...
            </div>
          </div>
        </GlassCard>
      </div>

      {/* Right Columns - Course List */}
      <div className="lg:col-span-2">
        <GlassCard className="h-full">
          <SectionHeader
            title="Course Completions"
            icon={FileText}
            className="mb-4"
          />
          <div className="grid grid-cols-1 xl:grid-cols-2 gap-4 max-h-96 overflow-y-auto">
            {termData.courses.map((course, index) => (
              <CourseCard key={index} course={course} index={index} />
            ))}
          </div>
        </GlassCard>
      </div>
    </div>
  );
}