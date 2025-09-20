import React from 'react';

const TypographySection = () => {
  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
  };

  const CodeBlock = ({ code }: { code: string }) => (
    <div className="bg-black/40 rounded-lg p-4 mt-4 backdrop-blur-sm border border-white/10">
      <code className="text-green-300 text-sm font-mono">{code}</code>
      <button
        onClick={() => copyToClipboard(code)}
        className="ml-4 text-blue-300 hover:text-blue-200 text-sm"
      >
        Copy
      </button>
    </div>
  );

  return (
    <div className="space-y-12">
      {/* Font Families */}
      <div className="glass-effect rounded-2xl p-8">
        <h3 className="text-2xl font-bold text-white font-space-grotesk mb-6">Font Families</h3>
        <div className="grid gap-6">
          <div>
            <h4 className="text-lg font-semibold text-white mb-2 font-space-grotesk">Space Grotesk (Headings)</h4>
            <p className="text-4xl font-space-grotesk text-white mb-2">The quick brown fox</p>
            <CodeBlock code='className="font-space-grotesk" or style={{ fontFamily: "var(--font-space-grotesk), sans-serif" }}' />
          </div>
          <div>
            <h4 className="text-lg font-semibold text-white mb-2 font-space-grotesk">Inter (Body Text)</h4>
            <p className="text-xl font-inter text-white mb-2">The quick brown fox jumps over the lazy dog</p>
            <CodeBlock code='className="font-inter" or style={{ fontFamily: "var(--font-inter), sans-serif" }}' />
          </div>
          <div>
            <h4 className="text-lg font-semibold text-white mb-2 font-space-grotesk">Crimson Text (Brand)</h4>
            <p className="text-4xl font-crimson text-white mb-2">IU-MiCert</p>
            <CodeBlock code='className="font-crimson" or style={{ fontFamily: "var(--font-crimson), serif" }}' />
          </div>
        </div>
      </div>

      {/* Heading Hierarchy */}
      <div className="glass-effect rounded-2xl p-8">
        <h3 className="text-2xl font-bold text-white font-space-grotesk mb-6">Heading Hierarchy</h3>
        <div className="space-y-4">
          <div>
            <h1 className="text-6xl font-bold text-white font-space-grotesk">Heading 1</h1>
            <CodeBlock code='className="text-6xl font-bold text-white font-space-grotesk"' />
          </div>
          <div>
            <h2 className="text-4xl font-bold text-white font-space-grotesk">Heading 2</h2>
            <CodeBlock code='className="text-4xl font-bold text-white font-space-grotesk"' />
          </div>
          <div>
            <h3 className="text-2xl font-bold text-white font-space-grotesk">Heading 3</h3>
            <CodeBlock code='className="text-2xl font-bold text-white font-space-grotesk"' />
          </div>
          <div>
            <h4 className="text-xl font-semibold text-white font-space-grotesk">Heading 4</h4>
            <CodeBlock code='className="text-xl font-semibold text-white font-space-grotesk"' />
          </div>
        </div>
      </div>

      {/* Body Text Sizes */}
      <div className="glass-effect rounded-2xl p-8">
        <h3 className="text-2xl font-bold text-white font-space-grotesk mb-6">Body Text Sizes</h3>
        <div className="space-y-4">
          <div>
            <p className="text-xl text-white font-inter">Large body text (text-xl)</p>
            <CodeBlock code='className="text-xl text-white font-inter"' />
          </div>
          <div>
            <p className="text-lg text-white font-inter">Medium body text (text-lg)</p>
            <CodeBlock code='className="text-lg text-white font-inter"' />
          </div>
          <div>
            <p className="text-base text-white font-inter">Regular body text (text-base)</p>
            <CodeBlock code='className="text-base text-white font-inter"' />
          </div>
          <div>
            <p className="text-sm text-white font-inter">Small text (text-sm)</p>
            <CodeBlock code='className="text-sm text-white font-inter"' />
          </div>
        </div>
      </div>

      {/* Text Colors */}
      <div className="glass-effect rounded-2xl p-8">
        <h3 className="text-2xl font-bold text-white font-space-grotesk mb-6">Text Colors</h3>
        <div className="space-y-4">
          <div>
            <p className="text-white font-inter text-lg">Primary text - Pure white</p>
            <CodeBlock code='className="text-white"' />
          </div>
          <div>
            <p className="text-white/80 font-inter text-lg">Secondary text - White 80%</p>
            <CodeBlock code='className="text-white/80"' />
          </div>
          <div>
            <p className="text-white/70 font-inter text-lg">Tertiary text - White 70%</p>
            <CodeBlock code='className="text-white/70"' />
          </div>
          <div>
            <p className="text-gray-300 font-inter text-lg">Body text - Gray 300</p>
            <CodeBlock code='className="text-gray-300"' />
          </div>
          <div>
            <p className="text-gray-400 font-inter text-lg">Muted text - Gray 400</p>
            <CodeBlock code='className="text-gray-400"' />
          </div>
          <div>
            <p className="text-blue-200 font-inter text-lg">Accent blue - Blue 200</p>
            <CodeBlock code='className="text-blue-200"' />
          </div>
          <div>
            <p className="text-purple-200 font-inter text-lg">Accent purple - Purple 200</p>
            <CodeBlock code='className="text-purple-200"' />
          </div>
        </div>
      </div>

      {/* Special Text Effects */}
      <div className="glass-effect rounded-2xl p-8">
        <h3 className="text-2xl font-bold text-white font-space-grotesk mb-6">Special Effects</h3>
        <div className="space-y-6">
          <div>
            <h4 className="text-4xl font-bold bg-gradient-to-r from-white to-pink-200 bg-clip-text text-transparent font-space-grotesk">
              Gradient Text
            </h4>
            <CodeBlock code='className="bg-gradient-to-r from-white to-pink-200 bg-clip-text text-transparent"' />
          </div>
          <div>
            <h4 className="text-4xl font-bold text-white font-space-grotesk" 
                style={{ textShadow: "0 0 30px rgba(255, 255, 255, 0.3), 0 4px 8px rgba(0, 0, 0, 0.8)" }}>
              Text with Glow
            </h4>
            <CodeBlock code='style={{ textShadow: "0 0 30px rgba(255, 255, 255, 0.3), 0 4px 8px rgba(0, 0, 0, 0.8)" }}' />
          </div>
        </div>
      </div>
    </div>
  );
};

export default TypographySection;