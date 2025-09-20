import {
  ProofData,
  AggregatedJourneyProof,
  IndividualTermProof,
} from "../types/proofs";

export const detectProofType = (jsonData: unknown): ProofData | null => {
  try {
    if (typeof jsonData !== "object" || jsonData === null) {
      return null;
    }

    const data = jsonData as Record<string, unknown>;

    if (
      data.type === "individual_term" ||
      data.type === "single_term" ||
      (data.claim &&
        data.verkle_proof &&
        typeof data.claim === "object" &&
        data.claim !== null &&
        "term" in data.claim)
    ) {
      return {
        type: "individual_term",
        ...data,
      } as IndividualTermProof;
    }

    if (
      data.student_info &&
      data.academic_terms &&
      Array.isArray(data.academic_terms)
    ) {
      return {
        type: "aggregated_journey",
        ...data,
      } as AggregatedJourneyProof;
    }

    return null;
  } catch {
    return null;
  }
};