import DashboardLayout from "../components/DashboardLayout";
import { Mail } from "lucide-react";

export default function RevokePage() {
  return (
    <DashboardLayout activeSection="revoke">
      <div className="h-full flex items-center justify-center">
        <div className="glass-effect rounded-xl p-12 text-center max-w-md">
          <div className="w-16 h-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center mx-auto mb-6">
            <Mail className="w-8 h-8 text-white" />
          </div>
          <h2 className="text-2xl font-bold text-white font-space-grotesk mb-4">
            Contact Issuer for Revocation
          </h2>
          <p className="text-purple-200 font-inter mb-6">
            This feature is currently unavailable in this portal. Please contact the issuer directly to request certificate revocation.
          </p>
          <p className="text-white/60 font-inter text-sm">
            Certificate revocation requires institutional authorization.
          </p>
        </div>
      </div>
    </DashboardLayout>
  );
}
