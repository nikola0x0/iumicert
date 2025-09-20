import React from 'react';

interface StyleGuideSectionProps {
  title: string;
  description: string;
  children: React.ReactNode;
}

const StyleGuideSection: React.FC<StyleGuideSectionProps> = ({ title, description, children }) => {
  return (
    <div className="mb-12">
      <div className="text-center mb-8">
        <h2 className="text-4xl font-bold text-white font-space-grotesk mb-2">
          {title}
        </h2>
        <p className="text-lg text-purple-200 font-inter">
          {description}
        </p>
      </div>
      <div className="space-y-8">
        {children}
      </div>
    </div>
  );
};

export default StyleGuideSection;