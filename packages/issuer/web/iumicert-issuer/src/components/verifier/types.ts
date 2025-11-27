// Type definitions for Receipt Verifier

export interface Course {
  course_id?: string;
  course_name?: string;
  grade: string;
  credits: number;
}

export interface TermReceipt {
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

export interface JourneyReceipt {
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

export interface VerificationResult {
  verified?: boolean;
  status?: string;
  courses_verified?: number;
  courses_failed?: number;
  ipa_verified?: boolean;
  blockchain_verified?: boolean;
  blockchain_anchored?: boolean;
  blockchain_tx_hash?: string;
  blockchain_block?: number;
  blockchain_published_at?: string;
  details?: string;
  error?: string;
}
