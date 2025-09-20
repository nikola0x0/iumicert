import React from 'react';
import { ArrowRight, Download, Eye, Shield, Plus, Check, X } from 'lucide-react';

const ButtonSection = () => {
  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
  };

  const CodeBlock = ({ code }: { code: string }) => (
    <div className="bg-black/40 rounded-lg p-4 mt-4 backdrop-blur-sm border border-white/10">
      <code className="text-green-300 text-sm font-mono break-all">{code}</code>
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
      {/* Primary Buttons */}
      <div className="glass-effect rounded-2xl p-8">
        <h3 className="text-2xl font-bold text-white font-space-grotesk mb-6">Primary Buttons</h3>
        <div className="grid gap-6">
          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">Primary CTA Button</h4>
            <div className="flex flex-wrap gap-4 mb-4">
              <button className="group relative bg-gradient-to-r from-blue-500 to-purple-600 hover:from-blue-600 hover:to-purple-700 px-10 py-4 rounded-full font-semibold text-white transition-all duration-300 hover:scale-110 hover:shadow-2xl flex items-center gap-3 font-inter overflow-hidden">
                <span className="absolute inset-0 bg-gradient-to-r from-blue-400 to-purple-500 opacity-0 group-hover:opacity-20 transition-opacity duration-300"></span>
                <span className="relative z-10">Start Verification</span>
                <ArrowRight className="w-6 h-6 relative z-10 group-hover:translate-x-1 transition-transform duration-300" />
              </button>
            </div>
            <CodeBlock code={`<button className="group relative bg-gradient-to-r from-blue-500 to-purple-600 hover:from-blue-600 hover:to-purple-700 px-10 py-4 rounded-full font-semibold text-white transition-all duration-300 hover:scale-110 hover:shadow-2xl flex items-center gap-3 font-inter overflow-hidden">
  <span className="absolute inset-0 bg-gradient-to-r from-blue-400 to-purple-500 opacity-0 group-hover:opacity-20 transition-opacity duration-300"></span>
  <span className="relative z-10">Start Verification</span>
  <ArrowRight className="w-6 h-6 relative z-10 group-hover:translate-x-1 transition-transform duration-300" />
</button>`} />
          </div>

          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">Alternative Primary Button</h4>
            <div className="flex flex-wrap gap-4 mb-4">
              <button className="bg-gradient-to-r from-pink-500 to-purple-600 hover:from-pink-600 hover:to-purple-700 px-10 py-4 rounded-full font-bold text-white text-lg transition-all duration-300 hover:scale-110 hover:shadow-xl flex items-center gap-3 font-inter">
                <Eye className="w-6 h-6" />
                Verify Now
                <ArrowRight className="w-6 h-6" />
              </button>
            </div>
            <CodeBlock code={`<button className="bg-gradient-to-r from-pink-500 to-purple-600 hover:from-pink-600 hover:to-purple-700 px-10 py-4 rounded-full font-bold text-white text-lg transition-all duration-300 hover:scale-110 hover:shadow-xl flex items-center gap-3 font-inter">
  <Eye className="w-6 h-6" />
  Verify Now
  <ArrowRight className="w-6 h-6" />
</button>`} />
          </div>
        </div>
      </div>

      {/* Secondary Buttons */}
      <div className="glass-effect rounded-2xl p-8">
        <h3 className="text-2xl font-bold text-white font-space-grotesk mb-6">Secondary Buttons</h3>
        <div className="grid gap-6">
          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">Glass Effect Button</h4>
            <div className="flex flex-wrap gap-4 mb-4">
              <button className="bg-white/10 hover:bg-white/20 backdrop-blur-sm border border-white/20 px-8 py-3 rounded-full font-semibold text-white transition-all duration-300 hover:scale-105 flex items-center gap-2 font-inter">
                <Download className="w-5 h-5" />
                Download
              </button>
            </div>
            <CodeBlock code={`<button className="bg-white/10 hover:bg-white/20 backdrop-blur-sm border border-white/20 px-8 py-3 rounded-full font-semibold text-white transition-all duration-300 hover:scale-105 flex items-center gap-2 font-inter">
  <Download className="w-5 h-5" />
  Download
</button>`} />
          </div>

          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">Outline Button</h4>
            <div className="flex flex-wrap gap-4 mb-4">
              <button className="border-2 border-white/40 hover:border-white/60 bg-transparent hover:bg-white/10 px-8 py-3 rounded-full font-semibold text-white transition-all duration-300 hover:scale-105 flex items-center gap-2 font-inter">
                <Shield className="w-5 h-5" />
                Learn More
              </button>
            </div>
            <CodeBlock code={`<button className="border-2 border-white/40 hover:border-white/60 bg-transparent hover:bg-white/10 px-8 py-3 rounded-full font-semibold text-white transition-all duration-300 hover:scale-105 flex items-center gap-2 font-inter">
  <Shield className="w-5 h-5" />
  Learn More
</button>`} />
          </div>
        </div>
      </div>

      {/* Button Sizes */}
      <div className="glass-effect rounded-2xl p-8">
        <h3 className="text-2xl font-bold text-white font-space-grotesk mb-6">Button Sizes</h3>
        <div className="grid gap-6">
          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">Large Button</h4>
            <div className="flex flex-wrap gap-4 mb-4">
              <button className="bg-gradient-to-r from-blue-500 to-purple-600 px-12 py-5 rounded-full font-bold text-white text-lg transition-all duration-300 hover:scale-105 font-inter">
                Large Button
              </button>
            </div>
            <CodeBlock code='className="px-12 py-5 text-lg"' />
          </div>

          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">Medium Button</h4>
            <div className="flex flex-wrap gap-4 mb-4">
              <button className="bg-gradient-to-r from-blue-500 to-purple-600 px-8 py-3 rounded-full font-semibold text-white transition-all duration-300 hover:scale-105 font-inter">
                Medium Button
              </button>
            </div>
            <CodeBlock code='className="px-8 py-3 font-semibold"' />
          </div>

          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">Small Button</h4>
            <div className="flex flex-wrap gap-4 mb-4">
              <button className="bg-gradient-to-r from-blue-500 to-purple-600 px-6 py-2 rounded-full font-medium text-white text-sm transition-all duration-300 hover:scale-105 font-inter">
                Small Button
              </button>
            </div>
            <CodeBlock code='className="px-6 py-2 text-sm font-medium"' />
          </div>
        </div>
      </div>

      {/* Icon Buttons */}
      <div className="glass-effect rounded-2xl p-8">
        <h3 className="text-2xl font-bold text-white font-space-grotesk mb-6">Icon Buttons</h3>
        <div className="grid gap-6">
          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">Circle Icon Buttons</h4>
            <div className="flex flex-wrap gap-4 mb-4">
              <button className="p-3 rounded-full bg-white/20 hover:bg-white/30 transition-all duration-200 hover:scale-110">
                <Plus className="w-5 h-5 text-white" />
              </button>
              <button className="p-3 rounded-full bg-green-500/20 hover:bg-green-500/30 transition-all duration-200 hover:scale-110">
                <Check className="w-5 h-5 text-green-300" />
              </button>
              <button className="p-3 rounded-full bg-red-500/20 hover:bg-red-500/30 transition-all duration-200 hover:scale-110">
                <X className="w-5 h-5 text-red-300" />
              </button>
            </div>
            <CodeBlock code={`<button className="p-3 rounded-full bg-white/20 hover:bg-white/30 transition-all duration-200 hover:scale-110">
  <Plus className="w-5 h-5 text-white" />
</button>`} />
          </div>

          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">Square Icon Buttons</h4>
            <div className="flex flex-wrap gap-4 mb-4">
              <button className="p-3 rounded-lg bg-white/10 hover:bg-white/20 backdrop-blur-sm border border-white/20 transition-all duration-200 hover:scale-105">
                <Download className="w-5 h-5 text-white" />
              </button>
              <button className="p-3 rounded-lg bg-blue-500/20 hover:bg-blue-500/30 backdrop-blur-sm border border-blue-500/30 transition-all duration-200 hover:scale-105">
                <Eye className="w-5 h-5 text-blue-300" />
              </button>
            </div>
            <CodeBlock code={`<button className="p-3 rounded-lg bg-white/10 hover:bg-white/20 backdrop-blur-sm border border-white/20 transition-all duration-200 hover:scale-105">
  <Download className="w-5 h-5 text-white" />
</button>`} />
          </div>
        </div>
      </div>

      {/* Button States */}
      <div className="glass-effect rounded-2xl p-8">
        <h3 className="text-2xl font-bold text-white font-space-grotesk mb-6">Button States</h3>
        <div className="grid gap-6">
          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">Normal, Hover, and Disabled</h4>
            <div className="flex flex-wrap gap-4 mb-4">
              <button className="bg-gradient-to-r from-blue-500 to-purple-600 px-8 py-3 rounded-full font-semibold text-white transition-all duration-300 hover:scale-105 font-inter">
                Normal
              </button>
              <button className="bg-gradient-to-r from-blue-600 to-purple-700 px-8 py-3 rounded-full font-semibold text-white scale-105 shadow-xl font-inter">
                Hovered
              </button>
              <button 
                disabled 
                className="bg-gray-600 px-8 py-3 rounded-full font-semibold text-gray-400 cursor-not-allowed font-inter opacity-50"
              >
                Disabled
              </button>
            </div>
            <CodeBlock code={`// Disabled state
<button 
  disabled 
  className="bg-gray-600 px-8 py-3 rounded-full font-semibold text-gray-400 cursor-not-allowed font-inter opacity-50"
>
  Disabled
</button>`} />
          </div>
        </div>
      </div>
    </div>
  );
};

export default ButtonSection;