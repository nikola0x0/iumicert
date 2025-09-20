import DashboardLayout from "../components/DashboardLayout";
import { Settings } from "lucide-react";

export default function SettingsPage() {
  return (
    <DashboardLayout activeSection="settings">
      <div className="h-full flex items-center justify-center">
        <div className="glass-effect rounded-xl p-12 text-center max-w-md">
          <div className="w-16 h-16 bg-gradient-to-br from-purple-500 to-pink-600 rounded-full flex items-center justify-center mx-auto mb-6">
            <Settings className="w-8 h-8 text-white" />
          </div>
          <h2 className="text-2xl font-bold text-white font-space-grotesk mb-4">
            Settings
          </h2>
          <p className="text-purple-200 font-inter mb-6">
            Configure system preferences, notifications, and user settings.
          </p>
          <p className="text-white/60 font-inter text-sm">
            Coming soon...
          </p>
        </div>
      </div>
    </DashboardLayout>
  );
}