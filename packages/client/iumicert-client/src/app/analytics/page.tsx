import DashboardLayout from "../components/DashboardLayout";
import { BarChart3 } from "lucide-react";

export default function AnalyticsPage() {
  return (
    <DashboardLayout activeSection="analytics">
      <div className="h-full flex items-center justify-center">
        <div className="glass-effect rounded-xl p-12 text-center max-w-md">
          <div className="w-16 h-16 bg-gradient-to-br from-green-500 to-blue-600 rounded-full flex items-center justify-center mx-auto mb-6">
            <BarChart3 className="w-8 h-8 text-white" />
          </div>
          <h2 className="text-2xl font-bold text-white font-space-grotesk mb-4">
            Analytics & Statistics
          </h2>
          <p className="text-purple-200 font-inter mb-6">
            View verification statistics, usage metrics, and system analytics.
          </p>
          <p className="text-white/60 font-inter text-sm">
            Coming soon...
          </p>
        </div>
      </div>
    </DashboardLayout>
  );
}