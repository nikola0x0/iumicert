"use client";
import React, { useState, useEffect, useCallback, useRef } from "react";
import Link from "next/link";
import Image from "next/image";
import AnimatedBackground from "./components/AnimatedBackground";
import { gsap } from "gsap";
import { TextPlugin } from "gsap/TextPlugin";
import {
  ChevronUp,
  ChevronDown,
  Shield,
  Zap,
  Users,
  Award,
  Eye,
  ArrowRight,
  CheckCircle,
} from "lucide-react";

// Register GSAP plugins
gsap.registerPlugin(TextPlugin);

interface GridFeature {
  icon: React.ComponentType<{ className?: string }>;
  title: string;
  desc: string;
}

interface Slide {
  id: string;
  title: string;
  subtitle: string;
  description: string;
  cta?: string;
  theme: string;
  features?: string[];
  gridFeatures?: GridFeature[];
}

const LandingPage = () => {
  const [currentSlide, setCurrentSlide] = useState(0);
  const [isAnimating, setIsAnimating] = useState(false);
  const [animationPhase, setAnimationPhase] = useState<
    "idle" | "exit" | "enter"
  >("idle");

  // Refs for GSAP animations
  const heroTitleRef = useRef<HTMLHeadingElement>(null);
  const heroSubtitleRef = useRef<HTMLParagraphElement>(null);
  const heroDescRef = useRef<HTMLParagraphElement>(null);
  const heroButtonRef = useRef<HTMLDivElement>(null);
  const heroContainerRef = useRef<HTMLDivElement>(null);
  const heroTitleContainerRef = useRef<HTMLDivElement>(null);
  const timelineRef = useRef<gsap.core.Timeline | null>(null);
  const floatAnimationRef = useRef<gsap.core.Tween | null>(null);

  const slides = [
    {
      id: "hero",
      title: "IU-MiCert",
      subtitle: "Blockchain Certificate Provenance Verification",
      description:
        "An advanced academic micro-credential provenance tracking system, leveraging Verkle trees for efficient and secure verification",
      cta: "Start Verification",
      theme: "from-blue-800 via-purple-900 to-slate-900",
    },
    {
      id: "problem",
      title: "The Challenge",
      subtitle: "Current Limitations",
      description:
        "Traditional systems only verify complete degrees, missing granular learning achievements and creating opportunities for credential forgery",
      features: [
        "Limited granular tracking",
        "Vulnerable to timeline manipulation",
        "Inefficient storage systems",
        "Lack of provenance",
      ],
      theme: "from-red-800 via-orange-900 to-slate-900",
    },
    {
      id: "solution",
      title: "Our Solution",
      subtitle: "Verkle Tree Architecture",
      description:
        "IU-MiCert leverages advanced Verkle trees to provide efficient micro-credential tracking with verifiable academic provenance",
      features: [
        "Enhanced storage efficiency",
        "Temporal verification",
        "Anti-forgery mechanisms",
        "Seamless integration",
      ],
      theme: "from-green-800 via-blue-900 to-slate-900",
    },
    {
      id: "features",
      title: "Key Features",
      subtitle: "Advanced Capabilities",
      description:
        "Comprehensive blockchain solution for modern educational credentialing",
      gridFeatures: [
        {
          icon: Shield,
          title: "Anti-Forgery",
          desc: "Timeline verification prevents credential manipulation",
        },
        {
          icon: Zap,
          title: "Efficient Storage",
          desc: "Verkle trees optimize micro-credential management",
        },
        {
          icon: Users,
          title: "Multi-Stakeholder",
          desc: "Students, employers, and institutions unified",
        },
        {
          icon: Award,
          title: "Comprehensive",
          desc: "Track entire learning journey, not just degrees",
        },
      ],
      theme: "from-purple-800 via-indigo-900 to-slate-900",
    },
    {
      id: "cta",
      title: "Ready to Verify?",
      subtitle: "Join the Future of Education",
      description:
        "Experience the next generation of credential verification with IU-MiCert",
      cta: "Start Verification Now",
      theme: "from-indigo-800 via-pink-900 to-slate-900",
    },
  ];

  const nextSlide = useCallback(() => {
    if (!isAnimating) {
      setIsAnimating(true);
      setAnimationPhase("exit");

      setTimeout(() => {
        setCurrentSlide((prev) => (prev + 1) % slides.length);
        setAnimationPhase("enter");
      }, 300);

      setTimeout(() => {
        setAnimationPhase("idle");
        setIsAnimating(false);
      }, 900);
    }
  }, [isAnimating, slides.length]);

  const prevSlide = useCallback(() => {
    if (!isAnimating) {
      setIsAnimating(true);
      setAnimationPhase("exit");

      setTimeout(() => {
        setCurrentSlide((prev) => (prev - 1 + slides.length) % slides.length);
        setAnimationPhase("enter");
      }, 300);

      setTimeout(() => {
        setAnimationPhase("idle");
        setIsAnimating(false);
      }, 900);
    }
  }, [isAnimating, slides.length]);

  const goToSlide = useCallback(
    (index: number) => {
      if (!isAnimating && index !== currentSlide) {
        setIsAnimating(true);
        setAnimationPhase("exit");

        setTimeout(() => {
          setCurrentSlide(index);
          setAnimationPhase("enter");
        }, 300);

        setTimeout(() => {
          setAnimationPhase("idle");
          setIsAnimating(false);
        }, 900);
      }
    },
    [isAnimating, currentSlide]
  );

  // Auto-advance slides with proper cleanup
  useEffect(() => {
    const interval = setInterval(() => {
      if (!isAnimating && document.visibilityState === 'visible') {
        nextSlide();
      }
    }, 300000);
    
    // Pause animations when page is not visible
    const handleVisibilityChange = () => {
      if (document.visibilityState === 'hidden') {
        if (timelineRef.current) {
          timelineRef.current.pause();
        }
        if (floatAnimationRef.current) {
          floatAnimationRef.current.pause();
        }
      } else {
        if (timelineRef.current) {
          timelineRef.current.resume();
        }
        if (floatAnimationRef.current) {
          floatAnimationRef.current.resume();
        }
      }
    };

    document.addEventListener('visibilitychange', handleVisibilityChange);
    
    return () => {
      clearInterval(interval);
      document.removeEventListener('visibilitychange', handleVisibilityChange);
    };
  }, [isAnimating, nextSlide]);

  // Cleanup animations on unmount
  useEffect(() => {
    return () => {
      if (timelineRef.current) {
        timelineRef.current.kill();
      }
      if (floatAnimationRef.current) {
        floatAnimationRef.current.kill();
      }
    };
  }, []);

  // GSAP Hero Animation
  useEffect(() => {
    // Cleanup previous animations
    if (timelineRef.current) {
      timelineRef.current.kill();
    }
    if (floatAnimationRef.current) {
      floatAnimationRef.current.kill();
    }

    // Always set initial states first to prevent flash
    if (currentSlide === 0) {
      gsap.set(
        [
          heroTitleContainerRef.current,
          heroSubtitleRef.current,
          heroDescRef.current,
          heroButtonRef.current,
        ],
        {
          opacity: 0,
          y: 50,
        }
      );
    }

    if (currentSlide === 0 && animationPhase === "idle") {
      // Reset and animate hero elements
      timelineRef.current = gsap.timeline();

      // Animate in sequence - simplified for performance
      timelineRef.current.to(heroTitleContainerRef.current, {
        opacity: 1,
        y: 0,
        duration: 1.2,
        ease: "power3.out",
      })
        .to(
          heroSubtitleRef.current,
          {
            opacity: 1,
            y: 0,
            duration: 0.8,
            ease: "power2.out",
          },
          "-=0.6"
        )
        .to(
          heroDescRef.current,
          {
            opacity: 1,
            y: 0,
            duration: 0.8,
            ease: "power2.out",
          },
          "-=0.4"
        )
        .to(
          heroButtonRef.current,
          {
            opacity: 1,
            y: 0,
            duration: 0.6,
            ease: "back.out(1.7)",
          },
          "-=0.2"
        );

      // Add subtle floating animation to the entire title container (includes backdrop)
      floatAnimationRef.current = gsap.to(heroTitleContainerRef.current, {
        y: -8,
        duration: 3,
        repeat: -1,
        yoyo: true,
        ease: "sine.inOut",
        delay: 2,
      });
    }
  }, [currentSlide, animationPhase]);

  // Get animation class based on current state
  const getAnimationClass = () => {
    if (animationPhase === "idle") return "";

    // Use zoom animation for all slides (same as "Our Solution" slide)
    if (animationPhase === "exit") {
      return "slide-exit-zoom";
    } else if (animationPhase === "enter") {
      return "slide-enter-zoom";
    }

    return "";
  };

  const renderSlideContent = (slide: Slide) => {
    switch (slide.id) {
      case "hero":
        return (
          <div
            ref={heroContainerRef}
            className="text-center space-y-8 relative"
          >
            {/* Optimized particles background - only show on hero slide */}
            {currentSlide === 0 && (
              <div className="absolute inset-0 pointer-events-none">
                {[
                  { left: 22, top: 15, duration: 4.2, delay: 0.5 },
                  { left: 78, top: 35, duration: 5.1, delay: 1.2 },
                  { left: 12, top: 75, duration: 3.8, delay: 0.8 },
                  { left: 65, top: 25, duration: 4.7, delay: 1.8 },
                  { left: 35, top: 85, duration: 3.9, delay: 0.3 },
                  { left: 88, top: 55, duration: 5.3, delay: 1.5 },
                  { left: 8, top: 45, duration: 4.1, delay: 0.9 },
                  { left: 55, top: 15, duration: 4.8, delay: 1.1 },
                  { left: 75, top: 65, duration: 3.7, delay: 1.7 },
                  { left: 25, top: 95, duration: 5.0, delay: 0.4 },
                ].map((particle, i) => (
                  <div
                    key={i}
                    className="absolute w-1 h-1 bg-blue-400/30 rounded-full will-change-transform"
                    style={{
                      left: `${particle.left}%`,
                      top: `${particle.top}%`,
                      animation: `float ${particle.duration}s ease-in-out infinite`,
                      animationDelay: `${particle.delay}s`,
                    }}
                  />
                ))}
              </div>
            )}

            <div className="space-y-6 relative z-10">
              <div
                ref={heroTitleContainerRef}
                className="relative flex justify-center"
                style={{ minHeight: "120px" }} // Prevent layout shift
              >
                {/* Single container that includes both backdrop and title */}
                <div className="relative inline-block">
                  {/* Background overlay for better contrast - positioned relative to title */}
                  <div className="absolute inset-0 bg-black/30 backdrop-blur-sm transform -rotate-1 rounded-2xl -mx-4 -my-2"></div>
                  <h1
                    ref={heroTitleRef}
                    className="text-7xl md:text-8xl font-bold text-white font-space-grotesk tracking-tight relative z-10 py-4 px-8"
                    style={{
                      textShadow:
                        "0 0 30px rgba(255, 255, 255, 0.3), 0 4px 8px rgba(0, 0, 0, 0.8)",
                    }}
                  >
                    {slide.title}
                  </h1>

                  {/* Interactive Certificate Easter Egg */}
                  <div className="absolute -top-10 -right-12 transform rotate-12 hover:rotate-[20deg] hover:scale-110 transition-all duration-500 cursor-pointer group">
                    <Image
                      src="/horizontal-certificate.svg"
                      alt="Secret Certificate"
                      width={100}
                      height={80}
                      className="certificate-easter-egg opacity-80 hover:opacity-100 drop-shadow-lg group-hover:drop-shadow-xl transition-all duration-300"
                    />
                    {/* Star effects on hover */}
                    <div className="absolute -top-2 -right-2 opacity-0 group-hover:opacity-100 transition-opacity duration-300">
                      <Image
                        src="/star.svg"
                        alt="Star"
                        width={16}
                        height={16}
                        className="group-hover:animate-bounce"
                        style={{ animationDelay: "0.1s" }}
                      />
                    </div>
                    <div className="absolute -top-3 right-2 opacity-0 group-hover:opacity-100 transition-opacity duration-300">
                      <Image
                        src="/star.svg"
                        alt="Star"
                        width={12}
                        height={12}
                        className="group-hover:animate-bounce"
                        style={{ animationDelay: "0.3s" }}
                      />
                    </div>
                    <div className="absolute -bottom-2 -left-3 opacity-0 group-hover:opacity-100 transition-opacity duration-300">
                      <Image
                        src="/star.svg"
                        alt="Star"
                        width={14}
                        height={14}
                        className="group-hover:animate-bounce"
                        style={{ animationDelay: "0.5s" }}
                      />
                    </div>
                    <div className="absolute top-1 -left-4 opacity-0 group-hover:opacity-100 transition-opacity duration-300">
                      <Image
                        src="/star.svg"
                        alt="Star"
                        width={10}
                        height={10}
                        className="group-hover:animate-bounce"
                        style={{ animationDelay: "0.2s" }}
                      />
                    </div>
                  </div>
                </div>
              </div>

              <p
                ref={heroSubtitleRef}
                className="text-2xl pt-6 text-blue-200 font-medium"
                style={{ fontFamily: "var(--font-space-grotesk), sans-serif" }}
              >
                {slide.subtitle}
              </p>
            </div>

            <p
              ref={heroDescRef}
              className="text-md text-gray-300 max-w-3xl mx-auto leading-relaxed font-inter"
            >
              {slide.description}
            </p>

            <div ref={heroButtonRef} className="flex justify-center">
              <Link href="/verifier" passHref>
                <button className="hover:cursor-pointer group relative bg-gradient-to-r from-blue-500 to-purple-600 hover:from-blue-600 hover:to-purple-700 px-10 py-4 rounded-full font-semibold text-white transition-all duration-300 hover:scale-110 hover:shadow-2xl flex items-center gap-3 font-inter overflow-hidden">
                  <span className="absolute inset-0 bg-gradient-to-r from-blue-400 to-purple-500 opacity-0 group-hover:opacity-20 transition-opacity duration-300"></span>
                  <span className="relative z-10">{slide.cta}</span>
                  <ArrowRight className="w-6 h-6 relative z-10 group-hover:translate-x-1 transition-transform duration-300" />
                </button>
              </Link>
            </div>

            {/* Custom animations CSS */}
            <style jsx>{`
              @keyframes gradient-shift {
                0%,
                100% {
                  background-position: 0% 50%;
                }
                50% {
                  background-position: 100% 50%;
                }
              }

              @keyframes float {
                0%,
                100% {
                  transform: translateY(0px) rotate(0deg);
                }
                50% {
                  transform: translateY(-20px) rotate(180deg);
                }
              }

              /* Optimize animations for better performance */
              .will-change-transform {
                will-change: transform;
              }

              /* Performance optimization for slide transitions */
              .slide-exit-zoom {
                animation: slideExitZoom 0.3s ease-in-out forwards;
                transform-origin: center;
              }

              .slide-enter-zoom {
                animation: slideEnterZoom 0.6s ease-out forwards;
                transform-origin: center;
              }

              @keyframes slideExitZoom {
                from {
                  opacity: 1;
                  transform: scale(1);
                }
                to {
                  opacity: 0;
                  transform: scale(0.95);
                }
              }

              @keyframes slideEnterZoom {
                from {
                  opacity: 0;
                  transform: scale(1.05);
                }
                to {
                  opacity: 1;
                  transform: scale(1);
                }
              }
            `}</style>
          </div>
        );

      case "problem":
        return (
          <div className="text-center space-y-6">
            <div className="space-y-2">
              <h2 className="text-4xl font-bold text-white font-space-grotesk">
                {slide.title}
              </h2>
              <p className="text-lg text-red-200 font-medium font-inter">
                {slide.subtitle}
              </p>
            </div>
            <p className="text-base text-gray-300 max-w-2xl mx-auto leading-relaxed font-inter">
              {slide.description}
            </p>
            <div className="grid grid-cols-2 gap-4 max-w-2xl mx-auto">
              {slide.features?.map((feature, index) => (
                <div
                  key={index}
                  className="bg-red-900/30 border border-red-500/30 rounded-lg p-4 backdrop-blur-sm hover:bg-red-800/40 transition-all duration-300"
                >
                  <div className="w-3 h-3 bg-red-400 rounded-full animate-pulse mx-auto mb-2"></div>
                  <h3 className="text-base font-semibold text-white mb-1 font-space-grotesk">
                    {feature
                      .split(" ")
                      .map(
                        (word) => word.charAt(0).toUpperCase() + word.slice(1)
                      )
                      .join(" ")}
                  </h3>
                  <p className="text-gray-300 text-xs font-inter">
                    Current limitation in credential systems
                  </p>
                </div>
              ))}
            </div>
          </div>
        );

      case "solution":
        return (
          <div className="text-center space-y-6">
            <div className="space-y-2">
              <h2 className="text-4xl font-bold text-white font-space-grotesk">
                {slide.title}
              </h2>
              <p className="text-lg text-green-200 font-medium font-inter">
                {slide.subtitle}
              </p>
            </div>
            <p className="text-base text-gray-300 max-w-2xl mx-auto leading-relaxed font-inter">
              {slide.description}
            </p>
            <div className="grid grid-cols-2 gap-4 max-w-2xl mx-auto">
              {slide.features?.map((feature, index) => (
                <div
                  key={index}
                  className="bg-green-900/30 border border-green-500/30 rounded-lg p-4 backdrop-blur-sm hover:bg-green-800/40 transition-all duration-300"
                >
                  <CheckCircle className="w-6 h-6 text-green-400 mx-auto mb-2" />
                  <h3 className="text-base font-semibold text-white mb-1 font-space-grotesk">
                    {feature
                      .split(" ")
                      .map(
                        (word) => word.charAt(0).toUpperCase() + word.slice(1)
                      )
                      .join(" ")}
                  </h3>
                  <p className="text-gray-300 text-xs font-inter">
                    Advanced solution capability
                  </p>
                </div>
              ))}
            </div>
          </div>
        );

      case "features":
        return (
          <div className="text-center space-y-6">
            <div className="space-y-2">
              <h2 className="text-4xl font-bold text-white font-space-grotesk">
                {slide.title}
              </h2>
              <p className="text-lg text-purple-200 font-medium font-inter">
                {slide.subtitle}
              </p>
            </div>
            <p className="text-base text-gray-300 max-w-2xl mx-auto leading-relaxed font-inter">
              {slide.description}
            </p>
            <div className="grid grid-cols-2 gap-4 max-w-2xl mx-auto">
              {slide.gridFeatures?.map((feature, index) => (
                <div
                  key={index}
                  className="bg-purple-900/30 border border-purple-500/30 rounded-lg p-4 backdrop-blur-sm hover:bg-purple-800/40 transition-all duration-300 hover:scale-105"
                >
                  <feature.icon className="w-6 h-6 text-purple-400 mx-auto mb-2" />
                  <h3 className="text-base font-semibold text-white mb-1 font-space-grotesk">
                    {feature.title}
                  </h3>
                  <p className="text-gray-300 text-xs font-inter">
                    {feature.desc}
                  </p>
                </div>
              ))}
            </div>
          </div>
        );

      case "cta":
        return (
          <div className="text-center space-y-8">
            <div className="space-y-4">
              <h2 className="text-5xl font-bold bg-gradient-to-r from-white to-pink-200 bg-clip-text text-transparent font-space-grotesk">
                {slide.title}
              </h2>
              <p className="text-xl text-pink-200 font-medium font-inter">
                {slide.subtitle}
              </p>
            </div>
            <p className="text-lg text-gray-300 max-w-2xl mx-auto leading-relaxed font-inter">
              {slide.description}
            </p>
            <div className="flex flex-col items-center gap-4">
              <Link href="/verifier" passHref>
                <button className="bg-gradient-to-r hover:cursor-pointer from-pink-500 to-purple-600 hover:from-pink-600 hover:to-purple-700 px-10 py-4 rounded-full font-bold text-white text-lg transition-all duration-300 hover:scale-110 hover:shadow-xl flex items-center gap-3 font-inter">
                  <Eye className="w-6 h-6" />
                  {slide.cta}
                  <ArrowRight className="w-6 h-6" />
                </button>
              </Link>
              <p className="text-sm text-gray-400 font-inter">
                No registration required • Instant verification
              </p>
            </div>
          </div>
        );

      default:
        return null;
    }
  };

  return (
    <div className="h-full w-full relative overflow-hidden">
      {/* Animated Background */}
      <AnimatedBackground
        gradient={slides[currentSlide].theme}
        className="transition-all duration-1000"
      />

      {/* Main Content */}
      <div className="relative z-10 h-full flex flex-col justify-center px-4 md:px-8 md:pr-20">
        <div className={`${getAnimationClass()}`}>
          {renderSlideContent(slides[currentSlide])}
        </div>
      </div>

      {/* Navigation Controls */}
      <div className="absolute right-8 top-1/2 transform -translate-y-1/2 z-20 hidden md:flex">
        <div className="flex flex-col items-center gap-6 bg-black/30 backdrop-blur-md rounded-full px-3 py-6">
          {/* Arrow Controls */}
          <div className="flex flex-col items-center gap-3">
            <button
              onClick={prevSlide}
              className="p-2 rounded-full bg-white/20 hover:bg-white/30 transition-all duration-200 hover:scale-110"
            >
              <ChevronUp className="w-5 h-5 text-white" />
            </button>
            <button
              onClick={nextSlide}
              className="p-2 rounded-full bg-white/20 hover:bg-white/30 transition-all duration-200 hover:scale-110"
            >
              <ChevronDown className="w-5 h-5 text-white" />
            </button>
          </div>

          {/* Slide Indicators */}
          <div className="flex flex-col gap-3">
            {slides.map((_, index) => (
              <button
                key={index}
                onClick={() => goToSlide(index)}
                className={`w-3 h-3 rounded-full transition-all duration-300 ${
                  index === currentSlide
                    ? "bg-white shadow-lg"
                    : "bg-white/40 hover:bg-white/60"
                }`}
              />
            ))}
          </div>
        </div>
      </div>

    </div>
  );
};

export default LandingPage;
