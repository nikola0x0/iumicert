import React from 'react';
import { 
  Shield, 
  Zap, 
  Users, 
  Award, 
  Eye, 
  ArrowRight, 
  CheckCircle, 
  Download,
  Upload,
  Settings,
  Home,
  Search,
  Bell,
  Heart,
  Star,
  Lock,
  Unlock,
  Plus,
  Minus,
  X,
  Check,
  AlertTriangle,
  Info,
  HelpCircle,
  Menu,
  ChevronUp,
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  ExternalLink,
  Copy,
  Share,
  Trash2,
  Edit,
  Save,
  File,
  Folder,
  Image,
  Video,
  Music,
  Calendar,
  Clock,
  Mail,
  Phone
} from 'lucide-react';

const IconSection = () => {
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

  const IconGrid = ({ icons, title }: { icons: Array<{icon: React.ComponentType<{className?: string}>, name: string}>, title: string }) => (
    <div className="glass-effect rounded-xl p-6 mb-6">
      <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">{title}</h4>
      <div className="grid grid-cols-4 md:grid-cols-8 gap-4">
        {icons.map(({ icon: Icon, name }, index) => (
          <div 
            key={index}
            className="flex flex-col items-center p-3 rounded-lg bg-white/5 hover:bg-white/10 transition-all duration-200 cursor-pointer group"
            onClick={() => copyToClipboard(`<${name} className="w-6 h-6" />`)}
          >
            <Icon className="w-6 h-6 text-white mb-2 group-hover:scale-110 transition-transform duration-200" />
            <span className="text-xs text-gray-300 text-center font-inter">{name}</span>
          </div>
        ))}
      </div>
    </div>
  );

  const brandIcons = [
    { icon: Shield, name: 'Shield' },
    { icon: Zap, name: 'Zap' },
    { icon: Users, name: 'Users' },
    { icon: Award, name: 'Award' },
    { icon: Eye, name: 'Eye' },
    { icon: CheckCircle, name: 'CheckCircle' },
    { icon: Lock, name: 'Lock' },
    { icon: Unlock, name: 'Unlock' },
  ];

  const navigationIcons = [
    { icon: Home, name: 'Home' },
    { icon: Search, name: 'Search' },
    { icon: Menu, name: 'Menu' },
    { icon: ArrowRight, name: 'ArrowRight' },
    { icon: ChevronUp, name: 'ChevronUp' },
    { icon: ChevronDown, name: 'ChevronDown' },
    { icon: ChevronLeft, name: 'ChevronLeft' },
    { icon: ChevronRight, name: 'ChevronRight' },
  ];

  const actionIcons = [
    { icon: Download, name: 'Download' },
    { icon: Upload, name: 'Upload' },
    { icon: Plus, name: 'Plus' },
    { icon: Minus, name: 'Minus' },
    { icon: X, name: 'X' },
    { icon: Check, name: 'Check' },
    { icon: Copy, name: 'Copy' },
    { icon: Share, name: 'Share' },
    { icon: Trash2, name: 'Trash2' },
    { icon: Edit, name: 'Edit' },
    { icon: Save, name: 'Save' },
    { icon: ExternalLink, name: 'ExternalLink' },
  ];

  const statusIcons = [
    { icon: AlertTriangle, name: 'AlertTriangle' },
    { icon: Info, name: 'Info' },
    { icon: HelpCircle, name: 'HelpCircle' },
    { icon: Bell, name: 'Bell' },
    { icon: Heart, name: 'Heart' },
    { icon: Star, name: 'Star' },
    { icon: Settings, name: 'Settings' },
  ];

  const fileIcons = [
    { icon: File, name: 'File' },
    { icon: Folder, name: 'Folder' },
    { icon: Image, name: 'Image' },
    { icon: Video, name: 'Video' },
    { icon: Music, name: 'Music' },
  ];

  const utilityIcons = [
    { icon: Calendar, name: 'Calendar' },
    { icon: Clock, name: 'Clock' },
    { icon: Mail, name: 'Mail' },
    { icon: Phone, name: 'Phone' },
  ];

  return (
    <div className="space-y-12">
      {/* Icon Library */}
      <div>
        <h3 className="text-2xl font-bold text-white font-space-grotesk mb-6">Icon Library (Lucide React)</h3>
        <p className="text-gray-300 font-inter mb-6">Click on any icon to copy its usage code. All icons use Lucide React.</p>
        
        <IconGrid icons={brandIcons} title="Brand & Security Icons" />
        <IconGrid icons={navigationIcons} title="Navigation Icons" />
        <IconGrid icons={actionIcons} title="Action Icons" />
        <IconGrid icons={statusIcons} title="Status & Feedback Icons" />
        <IconGrid icons={fileIcons} title="File Type Icons" />
        <IconGrid icons={utilityIcons} title="Utility Icons" />
      </div>

      {/* Icon Sizes */}
      <div className="glass-effect rounded-2xl p-8">
        <h3 className="text-2xl font-bold text-white font-space-grotesk mb-6">Icon Sizes</h3>
        <div className="grid gap-6">
          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">Standard Sizes</h4>
            <div className="flex items-center gap-6 mb-4 flex-wrap">
              <div className="flex flex-col items-center gap-2">
                <Shield className="w-4 h-4 text-white" />
                <span className="text-xs text-gray-300 font-inter">16px (w-4 h-4)</span>
              </div>
              <div className="flex flex-col items-center gap-2">
                <Shield className="w-5 h-5 text-white" />
                <span className="text-xs text-gray-300 font-inter">20px (w-5 h-5)</span>
              </div>
              <div className="flex flex-col items-center gap-2">
                <Shield className="w-6 h-6 text-white" />
                <span className="text-xs text-gray-300 font-inter">24px (w-6 h-6)</span>
              </div>
              <div className="flex flex-col items-center gap-2">
                <Shield className="w-8 h-8 text-white" />
                <span className="text-xs text-gray-300 font-inter">32px (w-8 h-8)</span>
              </div>
              <div className="flex flex-col items-center gap-2">
                <Shield className="w-12 h-12 text-white" />
                <span className="text-xs text-gray-300 font-inter">48px (w-12 h-12)</span>
              </div>
            </div>
            <CodeBlock code={`<Shield className="w-6 h-6 text-white" /> // Standard size
<Shield className="w-8 h-8 text-white" /> // Larger size`} />
          </div>
        </div>
      </div>

      {/* Icon Colors */}
      <div className="glass-effect rounded-2xl p-8">
        <h3 className="text-2xl font-bold text-white font-space-grotesk mb-6">Icon Colors</h3>
        <div className="grid gap-6">
          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">Color Variations</h4>
            <div className="flex items-center gap-6 mb-4 flex-wrap">
              <div className="flex flex-col items-center gap-2">
                <Shield className="w-8 h-8 text-white" />
                <span className="text-xs text-gray-300 font-inter">text-white</span>
              </div>
              <div className="flex flex-col items-center gap-2">
                <Shield className="w-8 h-8 text-white/80" />
                <span className="text-xs text-gray-300 font-inter">text-white/80</span>
              </div>
              <div className="flex flex-col items-center gap-2">
                <Shield className="w-8 h-8 text-blue-400" />
                <span className="text-xs text-gray-300 font-inter">text-blue-400</span>
              </div>
              <div className="flex flex-col items-center gap-2">
                <Shield className="w-8 h-8 text-purple-400" />
                <span className="text-xs text-gray-300 font-inter">text-purple-400</span>
              </div>
              <div className="flex flex-col items-center gap-2">
                <Shield className="w-8 h-8 text-green-400" />
                <span className="text-xs text-gray-300 font-inter">text-green-400</span>
              </div>
              <div className="flex flex-col items-center gap-2">
                <Shield className="w-8 h-8 text-red-400" />
                <span className="text-xs text-gray-300 font-inter">text-red-400</span>
              </div>
            </div>
            <CodeBlock code={`<Shield className="w-8 h-8 text-white" />
<Shield className="w-8 h-8 text-blue-400" />
<Shield className="w-8 h-8 text-purple-400" />`} />
          </div>
        </div>
      </div>

      {/* Icon in Context */}
      <div className="glass-effect rounded-2xl p-8">
        <h3 className="text-2xl font-bold text-white font-space-grotesk mb-6">Icons in Context</h3>
        <div className="grid gap-6">
          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">With Text</h4>
            <div className="space-y-3 mb-4">
              <div className="flex items-center gap-3">
                <Shield className="w-5 h-5 text-blue-400" />
                <span className="text-white font-inter">Security Features</span>
              </div>
              <div className="flex items-center gap-3">
                <Zap className="w-5 h-5 text-yellow-400" />
                <span className="text-white font-inter">Fast Performance</span>
              </div>
              <div className="flex items-center gap-3">
                <Users className="w-5 h-5 text-green-400" />
                <span className="text-white font-inter">Multi-User Support</span>
              </div>
            </div>
            <CodeBlock code={`<div className="flex items-center gap-3">
  <Shield className="w-5 h-5 text-blue-400" />
  <span className="text-white font-inter">Security Features</span>
</div>`} />
          </div>

          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">In Buttons</h4>
            <div className="flex gap-4 mb-4 flex-wrap">
              <button className="bg-gradient-to-r from-blue-500 to-purple-600 px-6 py-3 rounded-full font-semibold text-white transition-all duration-300 hover:scale-105 flex items-center gap-2 font-inter">
                <Eye className="w-5 h-5" />
                Verify
              </button>
              <button className="bg-white/10 hover:bg-white/20 backdrop-blur-sm border border-white/20 px-6 py-3 rounded-full font-semibold text-white transition-all duration-300 hover:scale-105 flex items-center gap-2 font-inter">
                <Download className="w-5 h-5" />
                Download
              </button>
              <button className="p-3 rounded-full bg-white/20 hover:bg-white/30 transition-all duration-200 hover:scale-110">
                <Settings className="w-5 h-5 text-white" />
              </button>
            </div>
            <CodeBlock code={`<button className="bg-gradient-to-r from-blue-500 to-purple-600 px-6 py-3 rounded-full font-semibold text-white transition-all duration-300 hover:scale-105 flex items-center gap-2 font-inter">
  <Eye className="w-5 h-5" />
  Verify
</button>`} />
          </div>

          <div>
            <h4 className="text-lg font-semibold text-white mb-4 font-space-grotesk">With Animation</h4>
            <div className="flex gap-4 mb-4 flex-wrap">
              <div className="flex items-center gap-2 text-white font-inter">
                <div className="w-3 h-3 bg-green-400 rounded-full animate-pulse"></div>
                <span>Online Status</span>
              </div>
              <button className="flex items-center gap-2 text-blue-300 hover:text-blue-200 transition-colors duration-200 group font-inter">
                <span>Learn More</span>
                <ArrowRight className="w-4 h-4 group-hover:translate-x-1 transition-transform duration-200" />
              </button>
            </div>
            <CodeBlock code={`// Pulsing dot
<div className="w-3 h-3 bg-green-400 rounded-full animate-pulse"></div>

// Arrow with hover animation  
<ArrowRight className="w-4 h-4 group-hover:translate-x-1 transition-transform duration-200" />`} />
          </div>
        </div>
      </div>
    </div>
  );
};

export default IconSection;