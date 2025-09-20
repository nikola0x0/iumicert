// Shared interfaces for proof verification
export interface CourseCompletion {
  course_code: string;
  course_name: string;
  grade: string;
  completion_date: string;
  credits: number;
  instructor: string;
}

export interface TermData {
  term: string;
  courses: CourseCompletion[];
  term_gpa?: number;
  total_credits: number;
}

export interface IndividualTermProof {
  type: "individual_term" | "single_term";
  claim: {
    student_id: string;
    term: string;
    claimed_courses: string[];
  };
  verkle_proof: {
    claimed_values: CourseCompletion[];
  };
  blockchain_reference: {
    contract_address: string;
    block_number: number;
    tree_commitment: string;
  };
  institutional_verification: {
    institution: string;
    semester: string;
  };
  metadata: {
    total_students_in_term?: number;
    proof_generation_timestamp: string;
  };
}

export interface AggregatedJourneyProof {
  type: "aggregated_journey";
  student_info: {
    student_id: string;
    student_name: string;
    program: string;
    enrollment_date: string;
  };
  academic_terms: TermData[];
  journey_summary: {
    total_terms: number;
    total_courses: number;
    total_credits: number;
    cumulative_gpa: number;
    start_date: string;
    latest_term: string;
  };
  multi_tree_verification_chain: Array<{
    term: string;
    blockchain_deployment: {
      contract_address: string;
      block_number: number;
      tree_commitment: string;
    };
  }>;
  institutional_verification: {
    institution: string;
  };
}

export type ProofData = IndividualTermProof | AggregatedJourneyProof;

export interface VerificationResult {
  isValid: boolean;
  message: string;
  proofData?: ProofData;
}