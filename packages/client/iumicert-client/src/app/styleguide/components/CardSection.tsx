import React from 'react';
import { Shield, Zap, CheckCircle, ArrowRight } from 'lucide-react';

const CardSection = () => {
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
      {/* Basic Glass Cards */}
      <div className="glass-effect rounded-2xl p-8">
        <h3 className="text-2xl font-bold text-white font-space-grotesk mb-6">Basic Glass Cards</h3>
        <div className="grid gap-6">
          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">Standard Glass Card</h4>
            <div className="glass-effect rounded-xl p-6 mb-4">
              <h5 className="text-xl font-bold text-white font-space-grotesk mb-2">Card Title</h5>
              <p className="text-gray-300 font-inter">This is a standard glass effect card with some content inside.</p>
            </div>
            <CodeBlock code={`<div className="glass-effect rounded-xl p-6">
  <h5 className="text-xl font-bold text-white font-space-grotesk mb-2">Card Title</h5>
  <p className="text-gray-300 font-inter">Card content goes here.</p>
</div>`} />
          </div>

          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">Hoverable Glass Card</h4>
            <div className="glass-effect rounded-xl p-6 mb-4 hover:bg-white/15 transition-all duration-300 hover:scale-105 cursor-pointer">
              <h5 className="text-xl font-bold text-white font-space-grotesk mb-2">Hoverable Card</h5>
              <p className="text-gray-300 font-inter">This card has hover effects and scales on hover.</p>
            </div>
            <CodeBlock code={`<div className="glass-effect rounded-xl p-6 hover:bg-white/15 transition-all duration-300 hover:scale-105 cursor-pointer">
  <h5 className="text-xl font-bold text-white font-space-grotesk mb-2">Hoverable Card</h5>
  <p className="text-gray-300 font-inter">This card has hover effects.</p>
</div>`} />
          </div>
        </div>
      </div>

      {/* Feature Cards */}
      <div className="glass-effect rounded-2xl p-8">
        <h3 className="text-2xl font-bold text-white font-space-grotesk mb-6">Feature Cards</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="bg-purple-900/30 border border-purple-500/30 rounded-lg p-6 backdrop-blur-sm hover:bg-purple-800/40 transition-all duration-300 hover:scale-105">
            <Shield className="w-8 h-8 text-purple-400 mb-4" />
            <h4 className="text-lg font-semibold text-white mb-2 font-space-grotesk">Anti-Forgery</h4>
            <p className="text-gray-300 text-sm font-inter">Timeline verification prevents credential manipulation</p>
          </div>
          <div className="bg-blue-900/30 border border-blue-500/30 rounded-lg p-6 backdrop-blur-sm hover:bg-blue-800/40 transition-all duration-300 hover:scale-105">
            <Zap className="w-8 h-8 text-blue-400 mb-4" />
            <h4 className="text-lg font-semibold text-white mb-2 font-space-grotesk">Efficient Storage</h4>
            <p className="text-gray-300 text-sm font-inter">Verkle trees optimize micro-credential management</p>
          </div>
        </div>
        <CodeBlock code={`<div className="bg-purple-900/30 border border-purple-500/30 rounded-lg p-6 backdrop-blur-sm hover:bg-purple-800/40 transition-all duration-300 hover:scale-105">
  <Shield className="w-8 h-8 text-purple-400 mb-4" />
  <h4 className="text-lg font-semibold text-white mb-2 font-space-grotesk">Feature Title</h4>
  <p className="text-gray-300 text-sm font-inter">Feature description</p>
</div>`} />
      </div>

      {/* Status Cards */}
      <div className="glass-effect rounded-2xl p-8">
        <h3 className="text-2xl font-bold text-white font-space-grotesk mb-6">Status Cards</h3>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="bg-green-900/30 border border-green-500/30 rounded-lg p-4 backdrop-blur-sm">
            <CheckCircle className="w-6 h-6 text-green-400 mb-2" />
            <h5 className="text-base font-semibold text-white mb-1 font-space-grotesk">Success State</h5>
            <p className="text-gray-300 text-xs font-inter">Operation completed successfully</p>
          </div>
          <div className="bg-yellow-900/30 border border-yellow-500/30 rounded-lg p-4 backdrop-blur-sm">
            <div className="w-3 h-3 bg-yellow-400 rounded-full animate-pulse mb-2"></div>
            <h5 className="text-base font-semibold text-white mb-1 font-space-grotesk">Warning State</h5>
            <p className="text-gray-300 text-xs font-inter">Attention required</p>
          </div>
          <div className="bg-red-900/30 border border-red-500/30 rounded-lg p-4 backdrop-blur-sm">
            <div className="w-3 h-3 bg-red-400 rounded-full animate-pulse mb-2"></div>
            <h5 className="text-base font-semibold text-white mb-1 font-space-grotesk">Error State</h5>
            <p className="text-gray-300 text-xs font-inter">Operation failed</p>
          </div>
        </div>
        <CodeBlock code={`<div className="bg-green-900/30 border border-green-500/30 rounded-lg p-4 backdrop-blur-sm">
  <CheckCircle className="w-6 h-6 text-green-400 mb-2" />
  <h5 className="text-base font-semibold text-white mb-1 font-space-grotesk">Success State</h5>
  <p className="text-gray-300 text-xs font-inter">Operation completed successfully</p>
</div>`} />
      </div>

      {/* Interactive Cards */}
      <div className="glass-effect rounded-2xl p-8">
        <h3 className="text-2xl font-bold text-white font-space-grotesk mb-6">Interactive Cards</h3>
        <div className="grid gap-6">
          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">Clickable Card with CTA</h4>
            <div className="glass-effect rounded-xl p-6 mb-4 hover:bg-white/15 transition-all duration-300 hover:scale-105 cursor-pointer group">
              <div className="flex justify-between items-start mb-4">
                <div>
                  <h5 className="text-xl font-bold text-white font-space-grotesk mb-2">Verify Certificate</h5>
                  <p className="text-gray-300 font-inter">Upload your certificate to verify its authenticity</p>
                </div>
                <ArrowRight className="w-6 h-6 text-white/60 group-hover:text-white group-hover:translate-x-1 transition-all duration-300" />
              </div>
              <div className="flex items-center gap-2 text-blue-300 font-inter text-sm">
                <span>Click to continue</span>
              </div>
            </div>
            <CodeBlock code={`<div className="glass-effect rounded-xl p-6 hover:bg-white/15 transition-all duration-300 hover:scale-105 cursor-pointer group">
  <div className="flex justify-between items-start mb-4">
    <div>
      <h5 className="text-xl font-bold text-white font-space-grotesk mb-2">Card Title</h5>
      <p className="text-gray-300 font-inter">Card description</p>
    </div>
    <ArrowRight className="w-6 h-6 text-white/60 group-hover:text-white group-hover:translate-x-1 transition-all duration-300" />
  </div>
  <div className="flex items-center gap-2 text-blue-300 font-inter text-sm">
    <span>Click to continue</span>
  </div>
</div>`} />
          </div>

          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">Card with Actions</h4>
            <div className="glass-effect rounded-xl p-6 mb-4">
              <h5 className="text-xl font-bold text-white font-space-grotesk mb-2">Certificate Status</h5>
              <p className="text-gray-300 font-inter mb-4">Your certificate has been successfully verified.</p>
              <div className="flex gap-3">
                <button className="bg-gradient-to-r from-blue-500 to-purple-600 px-4 py-2 rounded-full font-semibold text-white text-sm transition-all duration-300 hover:scale-105 font-inter">
                  Download
                </button>
                <button className="bg-white/10 hover:bg-white/20 backdrop-blur-sm border border-white/20 px-4 py-2 rounded-full font-semibold text-white text-sm transition-all duration-300 hover:scale-105 font-inter">
                  Share
                </button>
              </div>
            </div>
            <CodeBlock code={`<div className="glass-effect rounded-xl p-6">
  <h5 className="text-xl font-bold text-white font-space-grotesk mb-2">Card Title</h5>
  <p className="text-gray-300 font-inter mb-4">Card content</p>
  <div className="flex gap-3">
    <button className="bg-gradient-to-r from-blue-500 to-purple-600 px-4 py-2 rounded-full font-semibold text-white text-sm transition-all duration-300 hover:scale-105 font-inter">
      Primary Action
    </button>
    <button className="bg-white/10 hover:bg-white/20 backdrop-blur-sm border border-white/20 px-4 py-2 rounded-full font-semibold text-white text-sm transition-all duration-300 hover:scale-105 font-inter">
      Secondary Action
    </button>
  </div>
</div>`} />
          </div>
        </div>
      </div>

      {/* Card Layouts */}
      <div className="glass-effect rounded-2xl p-8">
        <h3 className="text-2xl font-bold text-white font-space-grotesk mb-6">Card Layouts</h3>
        <div className="grid gap-6">
          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">Grid Layout</h4>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
              {[1, 2, 3].map((num) => (
                <div key={num} className="glass-effect rounded-lg p-4">
                  <h5 className="text-lg font-bold text-white font-space-grotesk mb-2">Card {num}</h5>
                  <p className="text-gray-300 font-inter text-sm">Content for card {num}</p>
                </div>
              ))}
            </div>
            <CodeBlock code={`<div className="grid grid-cols-1 md:grid-cols-3 gap-4">
  <div className="glass-effect rounded-lg p-4">
    <h5 className="text-lg font-bold text-white font-space-grotesk mb-2">Card Title</h5>
    <p className="text-gray-300 font-inter text-sm">Card content</p>
  </div>
</div>`} />
          </div>
        </div>
      </div>
    </div>
  );
};

export default CardSection;