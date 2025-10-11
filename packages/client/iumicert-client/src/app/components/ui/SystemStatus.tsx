import { BarChart3, AlertCircle, Clock } from "lucide-react";

export type SystemStatusType = "online" | "development" | "down";

interface SystemStatusProps {
  status: SystemStatusType;
}

const statusConfig = {
  online: {
    background: "bg-green-500/10",
    border: "border-green-500/20", 
    iconBg: "bg-green-500/20",
    iconColor: "text-green-400",
    titleColor: "text-green-300",
    textColor: "text-green-200",
    dotColor: "bg-green-400",
    label: "Online",
    icon: BarChart3
  },
  development: {
    background: "bg-yellow-500/10",
    border: "border-yellow-500/20",
    iconBg: "bg-yellow-500/20", 
    iconColor: "text-yellow-400",
    titleColor: "text-yellow-300",
    textColor: "text-yellow-200",
    dotColor: "bg-yellow-400",
    label: "Development",
    icon: Clock
  },
  down: {
    background: "bg-red-500/10",
    border: "border-red-500/20",
    iconBg: "bg-red-500/20",
    iconColor: "text-red-400", 
    titleColor: "text-red-300",
    textColor: "text-red-200",
    dotColor: "bg-red-400",
    label: "System Down",
    icon: AlertCircle
  }
};

export default function SystemStatus({ status }: SystemStatusProps) {
  const config = statusConfig[status];
  const IconComponent = config.icon;
  
  return (
    <div className={`${config.background} border ${config.border} rounded-xl p-4`}>
      <div className="flex items-center gap-3">
        <div className={`w-6 h-6 ${config.iconBg} rounded-lg flex items-center justify-center`}>
          <IconComponent className={`w-3 h-3 ${config.iconColor}`} />
        </div>
        <div className="flex-1 min-w-0">
          <div className={`text-xs font-medium ${config.titleColor} font-inter`}>
            System Status
          </div>
          <div className="flex items-center gap-2">
            <div className={`w-1.5 h-1.5 ${config.dotColor} rounded-full ${status !== 'down' ? 'animate-pulse' : ''}`}></div>
            <span className={`text-xs ${config.textColor} font-inter`}>
              {config.label}
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}