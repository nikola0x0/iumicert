import React, { useState } from 'react';
import { ArrowRight, Star, Shield, Zap } from 'lucide-react';

const AnimationSection = () => {
  const [triggerAnimation, setTriggerAnimation] = useState(false);

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
  };

  const CodeBlock = ({ code }: { code: string }) => (
    <div className="bg-black/40 rounded-lg p-4 mt-4 backdrop-blur-sm border border-white/10">
      <code className="text-green-300 text-sm font-mono break-all whitespace-pre-wrap">{code}</code>
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
      {/* Hover Animations */}
      <div className="glass-effect rounded-2xl p-8">
        <h3 className="text-2xl font-bold text-white font-space-grotesk mb-6">Hover Animations</h3>
        <div className="grid gap-6">
          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">Scale on Hover</h4>
            <div className="flex gap-4 mb-4 flex-wrap">
              <button className="bg-gradient-to-r from-blue-500 to-purple-600 px-8 py-3 rounded-full font-semibold text-white transition-all duration-300 hover:scale-110 hover:shadow-2xl font-inter">
                Hover to Scale
              </button>
              <div className="glass-effect rounded-lg p-4 hover:scale-105 transition-all duration-300 cursor-pointer">
                <Shield className="w-6 h-6 text-purple-400 mb-2" />
                <p className="text-white text-sm font-inter">Hover this card</p>
              </div>
            </div>
            <CodeBlock code={`className="transition-all duration-300 hover:scale-110"
// For subtle scaling
className="transition-all duration-300 hover:scale-105"`} />
          </div>

          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">Icon Movement</h4>
            <div className="flex gap-4 mb-4 flex-wrap">
              <button className="flex items-center gap-2 bg-white/10 hover:bg-white/20 px-6 py-3 rounded-full font-semibold text-white transition-all duration-300 group font-inter">
                <span>Continue</span>
                <ArrowRight className="w-5 h-5 group-hover:translate-x-1 transition-transform duration-300" />
              </button>
              <button className="p-3 rounded-full bg-white/20 hover:bg-white/30 transition-all duration-200 hover:rotate-12 group">
                <Star className="w-5 h-5 text-white group-hover:animate-bounce" />
              </button>
            </div>
            <CodeBlock code={`// Arrow slide
<ArrowRight className="w-5 h-5 group-hover:translate-x-1 transition-transform duration-300" />

// Rotation and bounce
className="hover:rotate-12 transition-all duration-200"
className="group-hover:animate-bounce"`} />
          </div>

          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">Color Transitions</h4>
            <div className="flex gap-4 mb-4 flex-wrap">
              <button className="bg-white/10 hover:bg-white/30 text-white/70 hover:text-white px-8 py-3 rounded-full font-semibold transition-all duration-300 font-inter">
                Hover for Color
              </button>
              <div className="w-12 h-12 bg-blue-500/30 hover:bg-blue-500 rounded-full transition-colors duration-500 cursor-pointer"></div>
            </div>
            <CodeBlock code={`className="text-white/70 hover:text-white transition-all duration-300"
className="bg-blue-500/30 hover:bg-blue-500 transition-colors duration-500"`} />
          </div>
        </div>
      </div>

      {/* CSS Animations */}
      <div className="glass-effect rounded-2xl p-8">
        <h3 className="text-2xl font-bold text-white font-space-grotesk mb-6">CSS Animations</h3>
        <div className="grid gap-6">
          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">Floating Animation</h4>
            <div className="flex gap-4 mb-4 flex-wrap">
              <div className="bg-gradient-to-r from-purple-500 to-pink-500 p-4 rounded-xl floating">
                <Shield className="w-8 h-8 text-white" />
              </div>
              <div className="text-white font-inter text-lg floating">
                Floating Text
              </div>
            </div>
            <CodeBlock code={`className="floating"

/* CSS */
.floating {
  animation: float 6s ease-in-out infinite;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px);
  }
  50% {
    transform: translateY(-2px);
  }
}`} />
          </div>

          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">Pulse Animation</h4>
            <div className="flex gap-4 mb-4 items-center flex-wrap">
              <div className="w-4 h-4 bg-green-400 rounded-full animate-pulse"></div>
              <div className="w-4 h-4 bg-blue-400 rounded-full animate-pulse"></div>
              <div className="w-4 h-4 bg-red-400 rounded-full animate-pulse"></div>
              <span className="text-white font-inter">Status indicators</span>
            </div>
            <CodeBlock code={`className="animate-pulse"

// Or custom pulse
className="animate-pulse"
// Tailwind includes: opacity 1 -> 0.5 -> 1`} />
          </div>

          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">Bounce Animation</h4>
            <div className="flex gap-4 mb-4 items-center flex-wrap">
              <Star className="w-6 h-6 text-yellow-400 animate-bounce" />
              <button 
                className="bg-gradient-to-r from-green-500 to-blue-500 px-6 py-3 rounded-full font-semibold text-white font-inter"
                onMouseEnter={() => setTriggerAnimation(true)}
                onMouseLeave={() => setTriggerAnimation(false)}
              >
                <span className={triggerAnimation ? 'animate-bounce' : ''}>Hover me!</span>
              </button>
            </div>
            <CodeBlock code={`className="animate-bounce"

// Conditional bounce
className={isTriggered ? 'animate-bounce' : ''}`} />
          </div>

          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">Spin Animation</h4>
            <div className="flex gap-4 mb-4 items-center flex-wrap">
              <div className="w-8 h-8 border-2 border-blue-400 border-t-transparent rounded-full animate-spin"></div>
              <Zap className="w-6 h-6 text-yellow-400 hover:animate-spin transition-all duration-300" />
              <span className="text-white font-inter">Loading spinner & hover spin</span>
            </div>
            <CodeBlock code={`// Loading spinner
<div className="w-8 h-8 border-2 border-blue-400 border-t-transparent rounded-full animate-spin"></div>

// Hover spin
className="hover:animate-spin transition-all duration-300"`} />
          </div>
        </div>
      </div>

      {/* Custom Slide Animations */}
      <div className="glass-effect rounded-2xl p-8">
        <h3 className="text-2xl font-bold text-white font-space-grotesk mb-6">Custom Slide Animations</h3>
        <div className="grid gap-6">
          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">Slide Transitions (from globals.css)</h4>
            <div className="flex gap-4 mb-4 flex-wrap">
              <button 
                className="bg-blue-500 px-6 py-3 rounded-lg text-white font-inter slide-enter-right"
                onClick={() => {/* Animation already applied via class */}}
              >
                Slide from Right
              </button>
              <button 
                className="bg-purple-500 px-6 py-3 rounded-lg text-white font-inter slide-enter-up"
                onClick={() => {/* Animation already applied via class */}}
              >
                Slide from Bottom
              </button>
            </div>
            <CodeBlock code={`// Available slide classes (from globals.css):
className="slide-enter-right"   // Slides in from right
className="slide-enter-left"    // Slides in from left  
className="slide-enter-up"      // Fades in from bottom
className="slide-enter-zoom"    // Zooms in

// Exit animations:
className="slide-exit-left"     // Slides out to left
className="slide-exit-right"    // Slides out to right
className="slide-exit-down"     // Fades out downward
className="slide-exit-zoom"     // Zooms out`} />
          </div>

          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">Particle Animations</h4>
            <div className="relative bg-gradient-to-r from-blue-900 to-purple-900 rounded-xl p-8 mb-4 overflow-hidden">
              <p className="text-white font-inter text-center relative z-10">Animated Background with Particles</p>
              {/* Sample particles */}
              {[1, 2, 3, 4, 5].map((i) => (
                <div
                  key={i}
                  className="absolute w-1 h-1 bg-blue-400/30 rounded-full"
                  style={{
                    left: `${20 * i}%`,
                    top: `${15 + (i * 10)}%`,
                    animation: `float ${3 + i * 0.5}s ease-in-out infinite`,
                    animationDelay: `${i * 0.2}s`,
                  }}
                />
              ))}
            </div>
            <CodeBlock code={`// Floating particles
{particles.map((particle, i) => (
  <div
    key={i}
    className="absolute w-1 h-1 bg-blue-400/30 rounded-full"
    style={{
      left: \`\${particle.left}%\`,
      top: \`\${particle.top}%\`,
      animation: \`float \${particle.duration}s ease-in-out infinite\`,
      animationDelay: \`\${particle.delay}s\`,
    }}
  />
))}`} />
          </div>
        </div>
      </div>

      {/* Animation Performance Tips */}
      <div className="glass-effect rounded-2xl p-8">
        <h3 className="text-2xl font-bold text-white font-space-grotesk mb-6">Animation Best Practices</h3>
        <div className="space-y-4">
          <div className="bg-blue-900/20 border border-blue-500/30 rounded-lg p-4">
            <h4 className="text-lg font-semibold text-blue-300 mb-2 font-space-grotesk">Performance Tips</h4>
            <ul className="text-gray-300 font-inter text-sm space-y-2">
              <li>• Use <code className="text-blue-300">transform</code> and <code className="text-blue-300">opacity</code> for smooth animations</li>
              <li>• Prefer CSS transitions over JavaScript animations for simple effects</li>
              <li>• Use <code className="text-blue-300">will-change</code> property sparingly for complex animations</li>
              <li>• Keep animation durations reasonable (200-500ms for interactions)</li>
            </ul>
          </div>
          
          <div className="bg-green-900/20 border border-green-500/30 rounded-lg p-4">
            <h4 className="text-lg font-semibold text-green-300 mb-2 font-space-grotesk">Timing Guidelines</h4>
            <ul className="text-gray-300 font-inter text-sm space-y-2">
              <li>• Micro-interactions: 100-300ms</li>
              <li>• UI state changes: 200-500ms</li>
              <li>• Page transitions: 300-800ms</li>
              <li>• Decorative animations: 2-6s (with infinite loop)</li>
            </ul>
          </div>

          <div className="bg-purple-900/20 border border-purple-500/30 rounded-lg p-4">
            <h4 className="text-lg font-semibold text-purple-300 mb-2 font-space-grotesk">Easing Functions</h4>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-3">
              <div>
                <p className="text-gray-300 font-inter text-sm mb-2">Common Tailwind Easings:</p>
                <ul className="text-gray-400 font-mono text-xs space-y-1">
                  <li>• ease-linear</li>
                  <li>• ease-in</li>
                  <li>• ease-out</li>
                  <li>• ease-in-out</li>
                </ul>
              </div>
              <div>
                <p className="text-gray-300 font-inter text-sm mb-2">Custom Easings (CSS):</p>
                <ul className="text-gray-400 font-mono text-xs space-y-1">
                  <li>• cubic-bezier(0.25, 0.46, 0.45, 0.94)</li>
                  <li>• cubic-bezier(0.68, -0.55, 0.265, 1.55)</li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AnimationSection;