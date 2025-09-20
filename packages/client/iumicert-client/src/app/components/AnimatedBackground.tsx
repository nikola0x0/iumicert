"use client";

import { useState, useEffect } from "react";

interface Particle {
  id: number;
  left: number;
  top: number;
  delay: number;
  duration: number;
  size: number;
  opacity: number;
  animationType: "float" | "drift";
}

interface AnimatedBackgroundProps {
  gradient?: string;
  className?: string;
}

export default function AnimatedBackground({
  gradient = "from-blue-900 via-purple-900 to-indigo-900",
  className = "",
}: AnimatedBackgroundProps) {
  const [particles, setParticles] = useState<Particle[]>([]);
  const [isClient, setIsClient] = useState(false);

  // Initialize particles on client side only
  useEffect(() => {
    setIsClient(true);
    const generateParticles = () => {
      const animationTypes: ("float" | "drift")[] = ["float", "drift"];
      return [...Array(30)].map((_, i) => ({
        id: i,
        left: Math.random() * 100,
        top: Math.random() * 100,
        delay: Math.random() * 8,
        duration: 6 + Math.random() * 3,
        size: 2 + Math.random() * 3,
        opacity: 0.1 + Math.random() * 0.4,
        animationType:
          animationTypes[Math.floor(Math.random() * animationTypes.length)],
      }));
    };

    setParticles(generateParticles());

    // Reduced refresh frequency from every few seconds to every 30 seconds
    const refreshInterval = setInterval(() => {
      setParticles(generateParticles());
    }, 30000); // Refresh every 30 seconds

    return () => clearInterval(refreshInterval);
  }, []);

  return (
    <div
      className={`fixed inset-0 bg-gradient-to-br ${gradient} ${className} noise`}
    >
      <div className="absolute inset-0 bg-black/20"></div>
      {/* Floating particles */}
      <div className="absolute inset-0">
        {isClient &&
          particles.map((particle) => (
            <div
              key={particle.id}
              className={`absolute rounded-full animate-pulse ${
                particle.animationType === "float"
                  ? "bg-white/30"
                  : "bg-blue-300/40"
              }`}
              style={{
                left: `${particle.left}%`,
                top: `${particle.top}%`,
                width: `${particle.size}px`,
                height: `${particle.size}px`,
                opacity: particle.opacity,
                animation: `${particle.animationType} ${
                  particle.duration
                }s ease-in-out infinite, glow-pulse ${
                  particle.duration * 0.8
                }s ease-in-out infinite`,
                animationDelay: `${particle.delay}s`,
              }}
            ></div>
          ))}
      </div>

      {/* Noise overlay using SVG */}
      <style jsx>{`
        .noise:before {
          content: "";
          position: absolute;
          width: 100%;
          height: 100%;
          background: url("data:image/svg+xml,%0A%3Csvg xmlns='http://www.w3.org/2000/svg' width='500' height='500'%3E%3Cfilter id='noise' x='0' y='0'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.65' numOctaves='3' stitchTiles='stitch'/%3E%3CfeBlend mode='multiply'/%3E%3C/filter%3E%3Crect width='500' height='500' filter='url(%23noise)' opacity='0.3'/%3E%3C/svg%3E");
          mix-blend-mode: overlay;
          pointer-events: none;
          z-index: 1;
        }
      `}</style>
    </div>
  );
}
