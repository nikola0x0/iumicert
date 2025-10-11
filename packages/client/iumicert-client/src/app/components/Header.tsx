import Image from "next/image";
import { Link } from "next-view-transitions";

export default function Header() {
  return (
    <header className="fixed top-16 left-1/2 transform -translate-x-1/2 z-50 w-[90%] max-w-4xl">
      <div className="relative bg-white/10 ring-1 ring-white/30 backdrop-blur-sm rounded-3xl shadow-[0_0_40px_rgba(255,255,255,0.05)] px-6 py-4 overflow-hidden">
        {/* Glow effect */}
        <div className="absolute inset-0 rounded-3xl ring-1 ring-white/20 blur-md opacity-30 pointer-events-none mix-blend-overlay"></div>

        <div className="relative flex items-center justify-between">
          <Link href="/" className="flex items-center space-x-4">
            <div className="flex hover:cursor-pointer items-center space-x-4">
              <Image
                src="/logo.svg"
                alt="IU-MiCert Logo"
                width={80}
                height={80}
                className="object-contain transition-all duration-300 ease-in-out hover:scale-110 hover:rotate-3"
                style={{ viewTransitionName: "logo" }}
              />

              <div className="flex gap-6 justify-center align-middle items-center">
                <h1
                  className="text-4xl font-bold text-white font-crimson"
                  style={{
                    fontFamily: "var(--font-crimson), serif",
                    viewTransitionName: "main-title",
                  }}
                >
                  IU-MiCert
                </h1>
              </div>
            </div>
          </Link>

          {/* Navigation links */}
          <nav className="hidden md:flex items-center space-x-6">
            <Link
              href="/verifier"
              className="text-white/80 hover:text-white font-medium transition duration-300 text-sm font-inter"
              style={{
                fontFamily: "var(--font-inter), sans-serif",
                viewTransitionName: "nav-verifier",
              }}
            >
              Verifier Dashboard
            </Link>

            <a
              href="https://nikolaempire.gitbook.io/iu-micert/"
              target="_blank"
              rel="noopener noreferrer"
              className="text-white/80 hover:text-white font-medium flex items-center justify-center gap-[4px] align-middle transition duration-300 text-sm font-inter"
              style={{
                fontFamily: "var(--font-inter), sans-serif",
              }}
            >
              About & Docs
              <svg
                className="w-4 h-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.746 0 3.332.477 4.5 1.253v13C20.168 18.477 18.582 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"
                />
              </svg>
            </a>
          </nav>

          {/* Mobile menu button */}
          <button className="md:hidden text-white/80 hover:text-white">
            <svg
              className="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M4 6h16M4 12h16M4 18h16"
              />
            </svg>
          </button>
        </div>
      </div>
    </header>
  );
}
