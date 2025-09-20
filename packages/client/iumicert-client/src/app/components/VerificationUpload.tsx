import { useState } from "react";
import FileUploaderWrapper from "./FileUploaderWrapper";
import ProofTypesHelpModal from "./ProofTypesHelpModal";
import { Shield, Upload, FileText, HelpCircle } from "lucide-react";

interface VerificationUploadProps {
  credentialData: string;
  isLoading: boolean;
  onCredentialChange: (data: string) => void;
  onVerify: () => void;
  onFileChange: (file: File) => void;
  onTypeError: () => void;
}

export default function VerificationUpload({
  credentialData,
  isLoading,
  onCredentialChange,
  onVerify,
  onFileChange,
  onTypeError
}: VerificationUploadProps) {
  const fileTypes = ["JSON"];
  const [isHelpModalOpen, setIsHelpModalOpen] = useState(false);

  return (
    <>
      {/* Centered Upload Panel */}
      <div className="glass-effect rounded-xl p-8 max-w-lg w-full">
        {/* Header */}
        <div className="text-center mb-8">
          <div className="w-16 h-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center mx-auto mb-4">
            <Upload className="w-8 h-8 text-white" />
          </div>
          <h1 className="text-2xl font-bold text-white font-space-grotesk mb-2">
            Upload Proof Package
          </h1>
          <p className="text-purple-200 font-inter text-sm">
            Upload or paste your JSON proof file for verification
          </p>
          <button
            onClick={() => setIsHelpModalOpen(true)}
            className="inline-flex items-center gap-1 mt-2 text-xs text-blue-300 hover:text-blue-200 transition-colors duration-200"
            title="View supported proof types"
          >
            <HelpCircle className="w-3 h-3" />
            View supported formats
          </button>
        </div>

        {/* File Upload Area */}
        <div className="mb-4">
          <FileUploaderWrapper
            handleChange={onFileChange}
            name="proofFile"
            types={fileTypes}
            onTypeError={onTypeError}
            maxSize={10}
            classes="w-full file-uploader-custom hover:cursor-pointer"
          >
            <div className="w-full rounded-xl border-2 border-dashed border-white/30 bg-white/5 hover:border-blue-400 hover:bg-blue-500/20 transition-all duration-300 p-6 text-center">
              <div className="mb-3">
                <FileText className="mx-auto h-10 w-10 text-white/60" />
              </div>
              <p className="text-white/90 font-medium text-sm mb-1 font-space-grotesk">
                Drag & drop JSON file
              </p>
              <p className="text-white/60 text-xs font-inter">
                or click to browse files
              </p>
            </div>
          </FileUploaderWrapper>
        </div>

        <div className="text-center text-xs text-white/60 mb-4 font-inter">
          or paste JSON content below
        </div>

        {/* Textarea */}
        <textarea
          rows={6}
          className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-xl 
                   text-white placeholder-white/50 backdrop-blur-sm font-mono text-sm
                   focus:outline-none focus:ring-2 focus:ring-blue-400/50 focus:border-transparent 
                   transition duration-300 resize-none flex-1 mb-4"
          placeholder='Paste your proof package JSON here...'
          value={credentialData}
          onChange={(e) => onCredentialChange(e.target.value)}
        />

        {/* Verify Button */}
        <button
          onClick={onVerify}
          disabled={!credentialData.trim() || isLoading}
          className="w-full bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700
                   disabled:from-gray-600 disabled:to-gray-700 disabled:cursor-not-allowed
                   text-white font-bold py-3 px-6 rounded-xl transition-all duration-300
                   hover:scale-105 hover:shadow-xl flex items-center justify-center gap-2 font-space-grotesk"
        >
          {isLoading ? (
            <>
              <div className="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
              Verifying...
            </>
          ) : (
            <>
              <Shield className="w-5 h-5" />
              Verify Credentials
            </>
          )}
        </button>
      </div>

      {/* Help Modal */}
      <ProofTypesHelpModal
        isOpen={isHelpModalOpen}
        onClose={() => setIsHelpModalOpen(false)}
      />
    </>
  );
}