import { LucideIcon } from "lucide-react";
import IconBadge from "./IconBadge";

interface SectionHeaderProps {
  title: string;
  subtitle?: string;
  icon?: LucideIcon;
  iconColor?: "blue" | "green" | "purple" | "yellow" | "red" | "gray";
  className?: string;
}

export default function SectionHeader({ 
  title, 
  subtitle, 
  icon, 
  iconColor = "blue",
  className = "" 
}: SectionHeaderProps) {
  return (
    <div className={`flex items-center gap-3 ${className}`}>
      {icon && <IconBadge icon={icon} color={iconColor} />}
      <div>
        <h2 className="text-xl font-bold text-white font-space-grotesk">
          {title}
        </h2>
        {subtitle && (
          <p className="text-purple-200 text-sm font-inter">
            {subtitle}
          </p>
        )}
      </div>
    </div>
  );
}