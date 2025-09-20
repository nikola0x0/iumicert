import { ReactNode } from "react";

interface StatCardProps {
  label: string;
  value: string | number;
  color?: "blue" | "green" | "purple" | "yellow" | "red";
  className?: string;
}

const colorClasses = {
  blue: "bg-blue-500/20 text-blue-200",
  green: "bg-green-500/20 text-green-200", 
  purple: "bg-purple-500/20 text-purple-200",
  yellow: "bg-yellow-500/20 text-yellow-200",
  red: "bg-red-500/20 text-red-200"
};

const valueColorClasses = {
  blue: "text-blue-300",
  green: "text-green-300",
  purple: "text-purple-300", 
  yellow: "text-yellow-300",
  red: "text-red-300"
};

export default function StatCard({ 
  label, 
  value, 
  color = "blue", 
  className = "" 
}: StatCardProps) {
  return (
    <div className={`${colorClasses[color]} rounded-lg p-3 ${className}`}>
      <div className="text-xs font-medium font-inter">{label}</div>
      <div className={`text-xl font-bold ${valueColorClasses[color]}`}>
        {value}
      </div>
    </div>
  );
}