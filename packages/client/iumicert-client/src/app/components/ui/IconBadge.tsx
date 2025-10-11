import { LucideIcon } from "lucide-react";

interface IconBadgeProps {
  icon: LucideIcon;
  size?: "sm" | "md" | "lg";
  color?: "blue" | "green" | "purple" | "yellow" | "red" | "gray";
  variant?: "solid" | "outline";
}

const sizeClasses = {
  sm: { container: "w-6 h-6", icon: "w-3 h-3" },
  md: { container: "w-8 h-8", icon: "w-4 h-4" },
  lg: { container: "w-10 h-10", icon: "w-5 h-5" }
};

const colorClasses = {
  solid: {
    blue: "bg-blue-500/20 text-blue-400",
    green: "bg-green-500/20 text-green-400",
    purple: "bg-purple-500/20 text-purple-400", 
    yellow: "bg-yellow-500/20 text-yellow-400",
    red: "bg-red-500/20 text-red-400",
    gray: "bg-white/10 text-white/60"
  },
  outline: {
    blue: "border border-blue-400/30 text-blue-400",
    green: "border border-green-400/30 text-green-400",
    purple: "border border-purple-400/30 text-purple-400",
    yellow: "border border-yellow-400/30 text-yellow-400", 
    red: "border border-red-400/30 text-red-400",
    gray: "border border-white/20 text-white/60"
  }
};

export default function IconBadge({ 
  icon: Icon, 
  size = "md", 
  color = "blue", 
  variant = "solid" 
}: IconBadgeProps) {
  return (
    <div className={`
      ${sizeClasses[size].container} 
      ${colorClasses[variant][color]} 
      rounded-lg flex items-center justify-center
    `}>
      <Icon className={sizeClasses[size].icon} />
    </div>
  );
}