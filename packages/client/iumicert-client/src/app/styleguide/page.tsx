"use client";
import React, { useState } from "react";
import AnimatedBackground from "../components/AnimatedBackground";
import StyleGuideSection from "./components/StyleGuideSection";
import TypographySection from "./components/TypographySection";
import ButtonSection from "./components/ButtonSection";
import ColorSection from "./components/ColorSection";
import CardSection from "./components/CardSection";
import IconSection from "./components/IconSection";
import AnimationSection from "./components/AnimationSection";

type SectionType = 'typography' | 'buttons' | 'colors' | 'cards' | 'icons' | 'animations';

const StyleGuidePage = () => {
  const [activeSection, setActiveSection] = useState<SectionType>('typography');

  const sections: { id: SectionType; title: string; description: string }[] = [
    { id: 'typography', title: 'Typography', description: 'Fonts, sizes, and text styles' },
    { id: 'buttons', title: 'Buttons', description: 'Interactive elements and CTAs' },
    { id: 'colors', title: 'Colors', description: 'Color palette and gradients' },
    { id: 'cards', title: 'Cards', description: 'Glass effects and containers' },
    { id: 'icons', title: 'Icons', description: 'Icon system and usage' },
    { id: 'animations', title: 'Animations', description: 'Motion and transitions' },
  ];

  const renderSection = () => {
    switch (activeSection) {
      case 'typography':
        return <TypographySection />;
      case 'buttons':
        return <ButtonSection />;
      case 'colors':
        return <ColorSection />;
      case 'cards':
        return <CardSection />;
      case 'icons':
        return <IconSection />;
      case 'animations':
        return <AnimationSection />;
      default:
        return <TypographySection />;
    }
  };

  return (
    <div className="h-full w-full relative overflow-auto">
      {/* Animated Background */}
      <AnimatedBackground
        gradient="from-slate-900 via-purple-900 to-indigo-900"
        className="transition-all duration-1000"
      />

      {/* Main Content */}
      <div className="relative z-10 min-h-full pt-32 pb-32 px-4 md:px-8">
        {/* Header */}
        <div className="text-center mb-12">
          <h1 className="text-6xl font-bold text-white font-space-grotesk mb-4">
            Style Guide
          </h1>
          <p className="text-xl text-purple-200 font-inter max-w-2xl mx-auto">
            Design system and UI components for IU-MiCert applications
          </p>
        </div>

        {/* Navigation */}
        <div className="flex flex-wrap justify-center gap-4 mb-12">
          {sections.map((section) => (
            <button
              key={section.id}
              onClick={() => setActiveSection(section.id)}
              className={`px-6 py-3 rounded-full font-medium transition-all duration-300 backdrop-blur-sm ${
                activeSection === section.id
                  ? 'bg-white/20 text-white ring-2 ring-white/30'
                  : 'bg-white/10 text-white/70 hover:bg-white/15 hover:text-white'
              }`}
            >
              {section.title}
            </button>
          ))}
        </div>

        {/* Active Section */}
        <div className="max-w-6xl mx-auto">
          <StyleGuideSection
            title={sections.find(s => s.id === activeSection)?.title || ''}
            description={sections.find(s => s.id === activeSection)?.description || ''}
          >
            {renderSection()}
          </StyleGuideSection>
        </div>
      </div>
    </div>
  );
};

export default StyleGuidePage;