import Modal from "./ui/Modal";
import { Shield, FileText, BarChart3 } from "lucide-react";

interface ProofTypesHelpModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export default function ProofTypesHelpModal({ isOpen, onClose }: ProofTypesHelpModalProps) {
  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title="Supported Proof Types"
      size="lg"
    >
      <div className="space-y-6">
        <p className="text-purple-200 text-sm font-inter">
          Learn about the different types of verification formats supported by IU-MiCert.
        </p>

        <div className="space-y-4">
          <div className="bg-white/5 rounded-xl p-4 border border-white/10">
            <div className="flex items-start gap-3">
              <div className="w-6 h-6 bg-blue-500/20 rounded-lg flex items-center justify-center flex-shrink-0 mt-1">
                <FileText className="w-3 h-3 text-blue-400" />
              </div>
              <div>
                <h3 className="text-white font-bold mb-2 font-space-grotesk">
                  Individual Term Proof
                </h3>
                <p className="text-purple-200 text-sm font-inter leading-relaxed mb-3">
                  Verify courses completed in a specific semester with cryptographic proof anchored on blockchain. Perfect for selective disclosure of academic achievements.
                </p>
                <div className="flex flex-wrap gap-2">
                  <span className="bg-blue-500/20 text-blue-300 text-xs px-2 py-1 rounded font-mono">single_term</span>
                  <span className="bg-blue-500/20 text-blue-300 text-xs px-2 py-1 rounded font-mono">individual_term</span>
                </div>
              </div>
            </div>
          </div>

          <div className="bg-white/5 rounded-xl p-4 border border-white/10">
            <div className="flex items-start gap-3">
              <div className="w-6 h-6 bg-purple-500/20 rounded-lg flex items-center justify-center flex-shrink-0 mt-1">
                <BarChart3 className="w-3 h-3 text-purple-400" />
              </div>
              <div>
                <h3 className="text-white font-bold mb-2 font-space-grotesk">
                  Aggregated Journey
                </h3>
                <p className="text-purple-200 text-sm font-inter leading-relaxed mb-3">
                  Complete academic timeline across multiple terms with comprehensive verification chain. Shows full educational progression and achievements.
                </p>
                <div className="flex flex-wrap gap-2">
                  <span className="bg-purple-500/20 text-purple-300 text-xs px-2 py-1 rounded font-mono">aggregated_journey</span>
                </div>
              </div>
            </div>
          </div>

          <div className="bg-green-500/10 rounded-xl p-4 border border-green-500/20">
            <div className="flex items-start gap-3">
              <div className="w-6 h-6 bg-green-500/20 rounded-lg flex items-center justify-center flex-shrink-0 mt-1">
                <Shield className="w-3 h-3 text-green-400" />
              </div>
              <div>
                <h3 className="text-green-300 font-bold mb-2 font-space-grotesk">
                  Security Features
                </h3>
                <ul className="text-green-200 text-sm space-y-1 font-inter">
                  <li>• Verkle tree cryptographic proofs</li>
                  <li>• Blockchain anchoring on Sepolia</li>
                  <li>• Tamper-proof verification</li>
                  <li>• No institutional dependency</li>
                </ul>
              </div>
            </div>
          </div>
        </div>

        <div className="bg-blue-500/10 rounded-xl p-4 border border-blue-500/20">
          <h4 className="text-blue-300 font-bold mb-2 font-space-grotesk">
            How to Use
          </h4>
          <ol className="text-blue-200 text-sm space-y-1 font-inter list-decimal list-inside">
            <li>Upload or paste your JSON proof file in the upload area</li>
            <li>Click "Verify Credentials" to start the verification process</li>
            <li>View detailed results including course completions and blockchain verification</li>
            <li>Use "Copy Data" to export verification results for sharing</li>
          </ol>
        </div>
      </div>
    </Modal>
  );
}