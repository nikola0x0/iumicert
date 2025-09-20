import DashboardLayout from "../components/DashboardLayout";
import { Search } from "lucide-react";

export default function SearchPage() {
  return (
    <DashboardLayout activeSection="search">
      <div className="h-full flex items-center justify-center">
        <div className="glass-effect rounded-xl p-12 text-center max-w-md">
          <div className="w-16 h-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center mx-auto mb-6">
            <Search className="w-8 h-8 text-white" />
          </div>
          <h2 className="text-2xl font-bold text-white font-space-grotesk mb-4">
            Certificate Search
          </h2>
          <p className="text-purple-200 font-inter mb-6">
            Search through verified certificates and academic records.
          </p>
          <p className="text-white/60 font-inter text-sm">
            Coming soon...
          </p>
        </div>
      </div>
    </DashboardLayout>
  );
}