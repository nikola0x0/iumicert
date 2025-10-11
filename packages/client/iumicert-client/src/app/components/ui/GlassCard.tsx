import { ReactNode } from "react";

interface GlassCardProps {
  children: ReactNode;
  className?: string;
  padding?: "sm" | "md" | "lg";
}

const paddingClasses = {
  sm: "p-4",
  md: "p-6", 
  lg: "p-8"
};

export default function GlassCard({ 
  children, 
  className = "", 
  padding = "md" 
}: GlassCardProps) {
  return (
    <div className={`glass-effect rounded-xl ${paddingClasses[padding]} ${className}`}>
      {children}
    </div>
  );
}