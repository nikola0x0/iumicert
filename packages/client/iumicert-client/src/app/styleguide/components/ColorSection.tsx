import React from 'react';

const ColorSection = () => {
  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
  };

  const ColorSwatch = ({ color, name, className, description }: { 
    color: string; 
    name: string; 
    className: string;
    description?: string;
  }) => (
    <div 
      className={`${className} rounded-lg p-6 cursor-pointer transition-all duration-200 hover:scale-105`}
      onClick={() => copyToClipboard(className)}
    >
      <div className="text-white font-inter">
        <div className="font-semibold mb-1">{name}</div>
        <div className="text-sm opacity-80">{color}</div>
        {description && <div className="text-xs opacity-60 mt-1">{description}</div>}
      </div>
    </div>
  );

  const GradientSwatch = ({ name, className }: { 
    name: string; 
    className: string;
  }) => (
    <div 
      className={`${className} rounded-lg p-6 cursor-pointer transition-all duration-200 hover:scale-105`}
      onClick={() => copyToClipboard(className)}
    >
      <div className="text-white font-inter">
        <div className="font-semibold mb-1">{name}</div>
        <div className="text-xs opacity-80">Click to copy class</div>
      </div>
    </div>
  );

  return (
    <div className="space-y-12">
      {/* Brand Colors */}
      <div className="glass-effect rounded-2xl p-8">
        <h3 className="text-2xl font-bold text-white font-space-grotesk mb-6">Brand Colors</h3>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <ColorSwatch 
            color="#3B82F6" 
            name="Primary Blue" 
            className="bg-blue-500" 
            description="Main brand color"
          />
          <ColorSwatch 
            color="#8B5CF6" 
            name="Primary Purple" 
            className="bg-purple-500" 
            description="Secondary brand color"
          />
          <ColorSwatch 
            color="#EC4899" 
            name="Accent Pink" 
            className="bg-pink-500" 
            description="Accent color"
          />
          <ColorSwatch 
            color="#10B981" 
            name="Success Green" 
            className="bg-emerald-500" 
            description="Success states"
          />
        </div>
      </div>

      {/* Background Gradients */}
      <div className="glass-effect rounded-2xl p-8">
        <h3 className="text-2xl font-bold text-white font-space-grotesk mb-6">Background Gradients</h3>
        <div className="grid gap-4">
          <GradientSwatch 
            name="Hero Gradient" 
            className="bg-gradient-to-r from-blue-800 via-purple-900 to-slate-900 h-24"
          />
          <GradientSwatch 
            name="Problem Section" 
            className="bg-gradient-to-r from-red-800 via-orange-900 to-slate-900 h-24"
          />
          <GradientSwatch 
            name="Solution Section" 
            className="bg-gradient-to-r from-green-800 via-blue-900 to-slate-900 h-24"
          />
          <GradientSwatch 
            name="Features Section" 
            className="bg-gradient-to-r from-purple-800 via-indigo-900 to-slate-900 h-24"
          />
          <GradientSwatch 
            name="CTA Section" 
            className="bg-gradient-to-r from-indigo-800 via-pink-900 to-slate-900 h-24"
          />
          <GradientSwatch 
            name="Style Guide" 
            className="bg-gradient-to-r from-slate-900 via-purple-900 to-indigo-900 h-24"
          />
        </div>
      </div>

      {/* Button Gradients */}
      <div className="glass-effect rounded-2xl p-8">
        <h3 className="text-2xl font-bold text-white font-space-grotesk mb-6">Button Gradients</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <GradientSwatch 
            name="Primary Button" 
            className="bg-gradient-to-r from-blue-500 to-purple-600 h-16"
          />
          <GradientSwatch 
            name="CTA Button" 
            className="bg-gradient-to-r from-pink-500 to-purple-600 h-16"
          />
          <GradientSwatch 
            name="Primary Hover" 
            className="bg-gradient-to-r from-blue-600 to-purple-700 h-16"
          />
          <GradientSwatch 
            name="CTA Hover" 
            className="bg-gradient-to-r from-pink-600 to-purple-700 h-16"
          />
        </div>
      </div>

      {/* Text Colors */}
      <div className="glass-effect rounded-2xl p-8">
        <h3 className="text-2xl font-bold text-white font-space-grotesk mb-6">Text Colors</h3>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div className="bg-black/20 rounded-lg p-4">
            <div className="text-white font-semibold mb-2">Primary White</div>
            <div className="text-white text-sm">text-white</div>
            <div className="text-white/80 text-sm mt-1">text-white/80</div>
            <div className="text-white/70 text-sm">text-white/70</div>
          </div>
          <div className="bg-black/20 rounded-lg p-4">
            <div className="text-white font-semibold mb-2">Gray Scale</div>
            <div className="text-gray-300 text-sm">text-gray-300</div>
            <div className="text-gray-400 text-sm">text-gray-400</div>
            <div className="text-gray-500 text-sm">text-gray-500</div>
          </div>
          <div className="bg-black/20 rounded-lg p-4">
            <div className="text-white font-semibold mb-2">Blue Accents</div>
            <div className="text-blue-200 text-sm">text-blue-200</div>
            <div className="text-blue-300 text-sm">text-blue-300</div>
            <div className="text-blue-400 text-sm">text-blue-400</div>
          </div>
          <div className="bg-black/20 rounded-lg p-4">
            <div className="text-white font-semibold mb-2">Purple Accents</div>
            <div className="text-purple-200 text-sm">text-purple-200</div>
            <div className="text-purple-300 text-sm">text-purple-300</div>
            <div className="text-purple-400 text-sm">text-purple-400</div>
          </div>
        </div>
      </div>

      {/* Glass Effects */}
      <div className="glass-effect rounded-2xl p-8">
        <h3 className="text-2xl font-bold text-white font-space-grotesk mb-6">Glass Effects</h3>
        <div className="grid gap-6">
          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">Standard Glass Effect</h4>
            <div className="glass-effect rounded-xl p-6 mb-4">
              <p className="text-white font-inter">This uses the .glass-effect class</p>
            </div>
            <div className="bg-black/40 rounded-lg p-4 backdrop-blur-sm border border-white/10">
              <code className="text-green-300 text-sm font-mono">className=&quot;glass-effect&quot;</code>
              <button
                onClick={() => copyToClipboard('glass-effect')}
                className="ml-4 text-blue-300 hover:text-blue-200 text-sm"
              >
                Copy
              </button>
            </div>
          </div>
          
          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">Custom Glass Variations</h4>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="bg-white/10 backdrop-blur-sm border border-white/20 rounded-xl p-6">
                <p className="text-white text-sm font-inter">Light Glass</p>
                <code className="text-xs text-gray-300">bg-white/10</code>
              </div>
              <div className="bg-white/20 backdrop-blur-md border border-white/30 rounded-xl p-6">
                <p className="text-white text-sm font-inter">Medium Glass</p>
                <code className="text-xs text-gray-300">bg-white/20</code>
              </div>
              <div className="bg-white/30 backdrop-blur-lg border border-white/40 rounded-xl p-6">
                <p className="text-white text-sm font-inter">Heavy Glass</p>
                <code className="text-xs text-gray-300">bg-white/30</code>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ColorSection;