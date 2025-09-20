import { Github } from "lucide-react";

export default function Footer() {
  return (
    <footer className="fixed bottom-16 left-1/2 transform -translate-x-1/2 z-50 w-[90%] max-w-4xl">
      <div className="relative bg-white/10 ring-1 ring-white/30 backdrop-blur-sm rounded-3xl shadow-[0_0_40px_rgba(255,255,255,0.05)] px-6 py-4 overflow-hidden">
        {/* Glow effect */}
        <div className="absolute inset-0 rounded-3xl ring-1 ring-white/20 blur-md opacity-30 pointer-events-none mix-blend-overlay"></div>

        <div className="relative z-10 flex flex-col md:flex-row justify-between items-center space-y-2 md:space-y-0">
          <div
            className="text-sm text-white/80 flex items-center space-x-4"
            style={{ fontFamily: "var(--font-space-grotesk), sans-serif" }}
          >
            <span className="text-lg">© 2025 IU-MiCert</span>

            <div className="w-px h-4 bg-white/30"></div>

            <span className="text-sm bg-gradient-to-r from-blue-300 to-purple-300 bg-clip-text text-transparent font-inter">
              By nikola0x0 (Sepolia Testnet)
            </span>
          </div>

          <div className="flex items-center space-x-6">
            <a
              href="/privacy"
              className="text-sm text-white/60 hover:text-white/80 transition duration-300 font-inter"
            >
              Privacy
            </a>
            <a
              href="/terms"
              className="text-sm text-white/60 hover:text-white/80 transition duration-300 font-inter"
            >
              Terms
            </a>
            <a
              href="https://github.com/Niko1444/iumicert"
              target="_blank"
              rel="noopener noreferrer"
              className="text-white/60 hover:text-white/80 transition duration-300 p-1"
              aria-label="GitHub Repository"
            >
              <Github size={20} />
            </a>
          </div>
        </div>
      </div>
    </footer>
  );
}
