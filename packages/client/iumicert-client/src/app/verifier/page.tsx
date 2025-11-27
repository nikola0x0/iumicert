"use client";

import DashboardLayout from "../components/DashboardLayout";
import AcademicJourneyVerifier from "../components/AcademicJourneyVerifier";

export default function VerifierDashboard() {
  return (
    <DashboardLayout activeSection="verify">
      <AcademicJourneyVerifier />
    </DashboardLayout>
  );
}
