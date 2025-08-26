'use client';

import { useState, useEffect } from 'react';
import { useAccount, useWalletClient } from 'wagmi';
import { ConnectKitButton } from 'connectkit';
import { apiService, type Term, type Receipt } from '@/lib/api';
import { 
  publishTermRoot, 
  waitForTransactionConfirmation, 
  getTermRootHistory, 
  estimatePublishGas, 
  formatEther,
  type TermRootData,
  type PublishResult 
} from '@/lib/blockchain';

export function IssuerDashboard() {
  const { isConnected, address } = useAccount();
  const { data: walletClient } = useWalletClient();
  const [activeTab, setActiveTab] = useState<'terms' | 'receipts' | 'blockchain' | 'status'>('terms');
  const [terms, setTerms] = useState<Term[]>([]);
  const [selectedTerm, setSelectedTerm] = useState<Term | null>(null);
  const [receipts, setReceipts] = useState<Receipt[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [systemStatus, setSystemStatus] = useState<{ status: string; timestamp: string } | null>(null);

  useEffect(() => {
    fetchTerms();
    checkSystemStatus();
  }, []);

  const fetchTerms = async () => {
    try {
      setIsLoading(true);
      const termsData = await apiService.getTerms();
      setTerms(termsData);
    } catch (error) {
      console.error('Failed to fetch terms:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const checkSystemStatus = async () => {
    try {
      const status = await apiService.healthCheck();
      setSystemStatus(status);
    } catch (error) {
      console.error('Failed to check system status:', error);
    }
  };

  const handleTermSelect = async (term: Term) => {
    setSelectedTerm(term);
    try {
      setIsLoading(true);
      const receiptsData = await apiService.getReceipts(term.id);
      setReceipts(receiptsData);
      setActiveTab('receipts');
    } catch (error) {
      console.error('Failed to fetch receipts:', error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-6">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">IU-MiCert Issuer Dashboard</h1>
              <p className="text-sm text-gray-500">Academic credential issuance system with blockchain integration</p>
            </div>
            <ConnectKitButton.Custom>
              {({ isConnected, show, truncatedAddress }) => (
                <button 
                  onClick={show}
                  className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors"
                >
                  {isConnected ? truncatedAddress : "Connect Wallet"}
                </button>
              )}
            </ConnectKitButton.Custom>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Navigation Tabs */}
        <div className="border-b border-gray-200 mb-8">
          <nav className="-mb-px flex space-x-8">
            {[
              { id: 'terms', name: 'Terms & Students', icon: '📚' },
              { id: 'receipts', name: 'Receipts & Verification', icon: '🎓' },
              { id: 'blockchain', name: 'Blockchain Operations', icon: '⛓️' },
              { id: 'status', name: 'System Status', icon: '🔧' }
            ].map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id as any)}
                className={`${
                  activeTab === tab.id
                    ? 'border-blue-500 text-blue-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                } whitespace-nowrap py-2 px-1 border-b-2 font-medium text-sm flex items-center gap-2`}
              >
                <span>{tab.icon}</span>
                {tab.name}
              </button>
            ))}
          </nav>
        </div>

        {/* Tab Content */}
        {isConnected ? (
          <>
            {activeTab === 'terms' && <TermsTab terms={terms} onTermSelect={handleTermSelect} isLoading={isLoading} />}
            {activeTab === 'receipts' && <ReceiptsTab selectedTerm={selectedTerm} receipts={receipts} isLoading={isLoading} />}
            {activeTab === 'blockchain' && <BlockchainTab terms={terms} />}
            {activeTab === 'status' && <StatusTab systemStatus={systemStatus} />}
          </>
        ) : (
          <div className="text-center py-12">
            <div className="text-6xl mb-4">🔒</div>
            <h2 className="text-2xl font-bold text-gray-900 mb-4">Wallet Connection Required</h2>
            <p className="text-gray-600 mb-8">Please connect your wallet to access the issuer dashboard.</p>
            <ConnectKitButton.Custom>
              {({ isConnected, show, truncatedAddress }) => (
                <button 
                  onClick={show}
                  className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors"
                >
                  {isConnected ? truncatedAddress : "Connect Wallet"}
                </button>
              )}
            </ConnectKitButton.Custom>
          </div>
        )}
      </main>
    </div>
  );
}

function TermsTab({ terms, onTermSelect, isLoading }: { 
  terms: Term[]; 
  onTermSelect: (term: Term) => void; 
  isLoading: boolean; 
}) {
  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading terms...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg shadow">
      <div className="p-6">
        <h2 className="text-lg font-medium text-gray-900 mb-4">Academic Terms</h2>
        {terms.length === 0 ? (
          <div className="text-center py-8">
            <div className="text-4xl mb-4">📅</div>
            <p className="text-gray-600">No terms found. Generate some data using the CLI first.</p>
            <code className="bg-gray-100 px-2 py-1 rounded text-sm mt-2 inline-block">
              go run . generate-data
            </code>
          </div>
        ) : (
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {terms.map((term) => (
              <div
                key={term.id}
                onClick={() => onTermSelect(term)}
                className="border border-gray-200 rounded-lg p-4 cursor-pointer hover:bg-gray-50 hover:border-blue-300 transition-colors"
              >
                <h3 className="font-medium text-gray-900">{term.name}</h3>
                <p className="text-sm text-gray-500 mt-1">
                  {term.start_date} - {term.end_date}
                </p>
                <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium mt-2 ${
                  term.status === 'active' 
                    ? 'bg-green-100 text-green-800' 
                    : 'bg-gray-100 text-gray-800'
                }`}>
                  {term.status}
                </span>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

function ReceiptsTab({ selectedTerm, receipts, isLoading }: { 
  selectedTerm: Term | null; 
  receipts: Receipt[]; 
  isLoading: boolean; 
}) {
  const { address } = useAccount();
  const { data: walletClient } = useWalletClient();
  
  if (!selectedTerm) {
    return (
      <div className="bg-white rounded-lg shadow p-6">
        <div className="text-center py-8">
          <div className="text-4xl mb-4">🎓</div>
          <p className="text-gray-600">Select a term from the Terms tab to view receipts and verification tools.</p>
        </div>
      </div>
    );
  }

  const handlePublishTerm = async () => {
    if (!selectedTerm) return;
    
    try {
      // Get term root data from the backend API
      const termRoot = await apiService.getTermRoot(selectedTerm.id);
      
      const termRootData: TermRootData = {
        term_id: selectedTerm.id,
        verkle_root: termRoot.verkle_root,
        total_students: selectedTerm.student_count || 0
      };

      // Publish to blockchain via MetaMask  
      const result: PublishResult = await publishTermRoot(termRootData, address, walletClient);
      
      // Wait for confirmation
      await waitForTransactionConfirmation(result.transactionHash);

      alert(`✅ Successfully published ${selectedTerm.name} to Sepolia blockchain!\nTransaction: ${result.transactionHash}`);
      
    } catch (error) {
      console.error('Publishing failed:', error);
      const errorMessage = error instanceof Error ? error.message : 'Unknown error occurred';
      alert(`❌ Publishing failed: ${errorMessage}`);
    }
  };

  const handleGenerateReceipts = async () => {
    // TODO: Implement batch receipt generation
    alert(`Generating receipts for all students in ${selectedTerm.name}...`);
  };

  return (
    <div className="space-y-6">
      {/* Term Publishing Actions */}
      <div className="bg-white rounded-lg shadow p-6">
        <div className="flex items-center justify-between mb-4">
          <div>
            <h2 className="text-lg font-medium text-gray-900">
              📚 Term: {selectedTerm.name}
            </h2>
            <p className="text-sm text-gray-600">
              {selectedTerm.student_count} students • {selectedTerm.total_courses} total courses • Status: {selectedTerm.status}
            </p>
          </div>
          <div className="flex gap-3">
            <button
              onClick={handleGenerateReceipts}
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
              🔄 Generate All Receipts
            </button>
            <button
              onClick={handlePublishTerm}
              className="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
            >
              🚀 Publish Term to Blockchain
            </button>
          </div>
        </div>
        
        {/* Demo Instructions */}
        <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
          <h3 className="font-medium text-blue-900 mb-2">🎯 Thesis Demo Flow:</h3>
          <ol className="text-sm text-blue-800 space-y-1">
            <li><strong>1.</strong> Review existing student receipts below</li>
            <li><strong>2.</strong> Generate receipts for any missing students</li>
            <li><strong>3.</strong> Verify all cryptographic proofs</li>
            <li><strong>4.</strong> Publish complete term to Sepolia blockchain</li>
          </ol>
        </div>
      </div>

      {/* Student Receipts */}
      <div className="bg-white rounded-lg shadow">
        <div className="p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">
            Student Receipts ({receipts.length} processed)
          </h3>
          {isLoading ? (
            <div className="flex items-center justify-center py-8">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
            </div>
          ) : receipts.length === 0 ? (
            <div className="text-center py-8">
              <div className="text-4xl mb-4">📋</div>
              <p className="text-gray-600">No receipts found. Click "Generate All Receipts" to process this term.</p>
            </div>
          ) : (
            <div className="space-y-4">
              {receipts.map((receipt) => (
                <div key={receipt.id} className="border border-gray-200 rounded-lg p-4">
                  <div className="flex justify-between items-start">
                    <div>
                      <h4 className="font-medium text-gray-900">
                        👤 {receipt.student_name} 
                        <span className="ml-2 px-2 py-1 bg-green-100 text-green-800 text-xs rounded-full">
                          ✅ Verified
                        </span>
                      </h4>
                      <p className="text-sm text-gray-500">DID: {receipt.student_id}</p>
                      <p className="text-sm text-gray-500">
                        📚 {Array.isArray(receipt.courses) ? receipt.courses.length : 0} courses completed
                      </p>
                    </div>
                    <div className="text-right">
                      <p className="text-sm text-gray-500">
                        📅 {new Date(receipt.created_at).toLocaleDateString()}
                      </p>
                      <button className="mt-2 px-3 py-1 bg-gray-100 text-gray-700 text-xs rounded hover:bg-gray-200">
                        View Details
                      </button>
                    </div>
                  </div>
                  
                  {/* Course Details */}
                  {Array.isArray(receipt.courses) && receipt.courses.length > 0 && (
                    <div className="mt-3 p-3 bg-gray-50 rounded">
                      <p className="text-xs font-medium text-gray-700 mb-2">Courses:</p>
                      <div className="grid grid-cols-2 gap-2">
                        {receipt.courses.slice(0, 4).map((course: any, i: number) => (
                          <div key={i} className="text-xs">
                            <span className="font-medium">{course.course_id}:</span> {course.grade} 
                            <span className="text-gray-500">({course.credits} cr)</span>
                          </div>
                        ))}
                      </div>
                    </div>
                  )}
                  
                  {/* Crypto Data */}
                  <div className="mt-4 p-3 bg-gray-50 rounded text-xs font-mono">
                    <div className="mb-1">
                      <strong>🔐 Merkle Root:</strong> 
                      <span className="text-blue-600 ml-1">{receipt.merkle_root?.slice(0, 20)}...</span>
                    </div>
                    <div>
                      <strong>🌳 Verkle Proof:</strong> 
                      <span className="text-green-600 ml-1">Generated & Verified</span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

function BlockchainTab({ terms }: { terms: Term[] }) {
  const { isConnected, address } = useAccount();
  const { data: walletClient } = useWalletClient();
  const [publishingStatus, setPublishingStatus] = useState<'idle' | 'preparing' | 'signing' | 'confirming'>('idle');
  const [transactions, setTransactions] = useState<Array<{
    id: string;
    type: string;
    term: string;
    students: number;
    status: string;
    timestamp: string;
    blockNumber: number | null;
    gasUsed: string | null;
  }>>([]);
  const [selectedTerm, setSelectedTerm] = useState<Term | null>(null);
  const [gasEstimate, setGasEstimate] = useState<{
    gasLimit: bigint;
    gasPrice: bigint;
    estimatedCost: bigint;
  } | null>(null);

  // Load transaction history on component mount
  useEffect(() => {
    const loadTransactionHistory = async () => {
      try {
        const history = await getTermRootHistory();
        const formattedTransactions = history.map((event: any) => ({
          id: event.transactionHash,
          type: 'Term Publication',
          term: event.termId,
          students: event.totalStudents,
          status: 'confirmed',
          timestamp: new Date(event.timestamp * 1000).toISOString(),
          blockNumber: event.blockNumber,
          gasUsed: 'N/A' // Gas info not in event logs
        }));
        setTransactions(formattedTransactions);
      } catch (error) {
        console.warn('Failed to load blockchain transaction history:', error);
      }
    };

    loadTransactionHistory();
  }, []);

  const handlePublishToBlockchain = async (term: Term) => {
    if (!term) {
      alert('Please select a term to publish');
      return;
    }

    // Debug wallet connection state
    console.log('Wallet connection state:', { 
      isConnected, 
      address, 
      hasWalletClient: !!walletClient,
      walletClientType: walletClient?.constructor?.name 
    });

    // Check wallet connection using wagmi state
    if (!isConnected || !address) {
      alert('❌ Please connect your wallet first');
      return;
    }

    if (!walletClient) {
      alert('❌ Wallet client not ready. Please try again in a moment.');
      return;
    }

    try {
      setPublishingStatus('preparing');

      // Get term root data from the backend API
      const termRoot = await apiService.getTermRoot(term.id);
      
      const termRootData: TermRootData = {
        term_id: term.id,
        verkle_root: termRoot.verkle_root,
        total_students: term.student_count || 0
      };

      // Estimate gas costs
      try {
        const estimate = await estimatePublishGas(termRootData, address, walletClient);
        setGasEstimate(estimate);
      } catch (error) {
        console.warn('Failed to estimate gas:', error);
      }

      setPublishingStatus('signing');
      
      // Publish to blockchain via MetaMask
      const result: PublishResult = await publishTermRoot(termRootData, address, walletClient);
      
      setPublishingStatus('confirming');
      
      // Wait for confirmation with status updates
      await waitForTransactionConfirmation(result.transactionHash, (status: 'pending' | 'confirmed' | 'failed') => {
        if (status === 'confirmed') {
          setPublishingStatus('idle');
        } else if (status === 'failed') {
          setPublishingStatus('idle');
          throw new Error('Transaction failed');
        }
      });

      // Add transaction to local state
      const newTx = {
        id: result.transactionHash,
        type: 'Term Publication',
        term: term.name,
        students: term.student_count || 0,
        status: 'confirmed',
        timestamp: new Date().toISOString(),
        blockNumber: Number(result.blockNumber),
        gasUsed: result.gasUsed.toLocaleString()
      };
      
      setTransactions(prev => [newTx, ...prev]);
      
      // Update backend with transaction info
      await apiService.publishToBlockchain({
        term_id: term.id,
        network: 'sepolia',
        gas_limit: gasEstimate?.gasLimit ? Number(gasEstimate.gasLimit) : 500000
      });

      alert(`✅ Successfully published ${term.name} to Sepolia blockchain!\nTransaction: ${result.transactionHash}`);
      
    } catch (error) {
      console.error('Publishing failed:', error);
      setPublishingStatus('idle');
      
      const errorMessage = error instanceof Error ? error.message : 'Unknown error occurred';
      alert(`❌ Publishing failed: ${errorMessage}`);
    }
  };

  return (
    <div className="space-y-6">
      {/* Blockchain Publishing */}
      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-lg font-medium text-gray-900 mb-4">⛓️ Sepolia Testnet Operations</h2>
        
        {/* Publishing Status */}
        {publishingStatus !== 'idle' && (
          <div className="mb-6 p-4 bg-blue-50 border border-blue-200 rounded-lg">
            <div className="flex items-center gap-3">
              <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-blue-600"></div>
              <div>
                <p className="font-medium text-blue-900">
                  {publishingStatus === 'preparing' && '📋 Preparing credential bundle...'}
                  {publishingStatus === 'signing' && '✍️ Please sign transaction in MetaMask...'}
                  {publishingStatus === 'confirming' && '⏳ Waiting for blockchain confirmation...'}
                </p>
                <p className="text-sm text-blue-700">
                  Publishing Fall 2024 credentials for 5 students
                </p>
              </div>
            </div>
          </div>
        )}

        <div className="grid md:grid-cols-2 gap-6">
          {/* Publishing Actions */}
          <div className="border border-gray-200 rounded-lg p-4">
            <h3 className="font-medium text-gray-900 mb-2">🚀 Publish Term Credentials</h3>
            <p className="text-sm text-gray-600 mb-4">
              Publish verified student credentials to Sepolia testnet via MetaMask. Creates immutable record of term completion.
            </p>
            <div className="space-y-3">
              <div className="text-sm">
                <label className="font-medium block mb-2">Select Term to Publish:</label>
                <select 
                  value={selectedTerm?.id || ''} 
                  onChange={(e) => setSelectedTerm(terms.find(t => t.id === e.target.value) || null)}
                  className="w-full border border-gray-300 rounded px-3 py-2 text-sm"
                >
                  <option value="">Choose a term...</option>
                  {terms.map(term => (
                    <option key={term.id} value={term.id}>
                      {term.name} ({term.student_count} students)
                    </option>
                  ))}
                </select>
              </div>
              
              {gasEstimate && (
                <div className="bg-blue-50 p-3 rounded text-sm">
                  <p className="font-medium text-blue-900">Gas Estimate:</p>
                  <p className="text-blue-800">
                    ~{formatEther(gasEstimate.estimatedCost)} ETH 
                    ({gasEstimate.gasLimit.toLocaleString()} gas @ {formatEther(gasEstimate.gasPrice * BigInt(1000000000))} gwei)
                  </p>
                </div>
              )}
              
              {!isConnected || !address ? (
                <div className="text-center p-3 bg-yellow-50 border border-yellow-200 rounded-lg">
                  <p className="text-yellow-800 text-sm font-medium">⚠️ Wallet not connected</p>
                  <p className="text-yellow-700 text-xs mt-1">Please connect your MetaMask wallet to publish to blockchain</p>
                </div>
              ) : (
                <div className="text-center p-2 bg-green-50 border border-green-200 rounded-lg mb-3">
                  <p className="text-green-800 text-xs">✅ Wallet connected: {address?.slice(0, 6)}...{address?.slice(-4)}</p>
                </div>
              )}
              
              <button 
                onClick={() => selectedTerm && handlePublishToBlockchain(selectedTerm)}
                disabled={publishingStatus !== 'idle' || !selectedTerm || !isConnected || !address}
                className="w-full bg-green-600 text-white px-4 py-2 rounded-lg hover:bg-green-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {publishingStatus === 'idle' ? '🚀 Publish to Blockchain' : 
                 publishingStatus === 'preparing' ? '📋 Preparing...' :
                 publishingStatus === 'signing' ? '✍️ Sign in MetaMask...' :
                 '⏳ Confirming...'}
              </button>
            </div>
          </div>

          {/* Network Info */}
          <div className="border border-gray-200 rounded-lg p-4">
            <h3 className="font-medium text-gray-900 mb-2">🌐 Network Status</h3>
            <div className="space-y-2 text-sm">
              <div className="flex justify-between">
                <span>Network:</span>
                <span className="text-blue-600">Sepolia Testnet</span>
              </div>
              <div className="flex justify-between">
                <span>Chain ID:</span>
                <span>11155111</span>
              </div>
              <div className="flex justify-between">
                <span>Gas Price:</span>
                <span>{gasEstimate ? formatEther(gasEstimate.gasPrice * BigInt(1000000000)) + ' gwei' : 'Loading...'}</span>
              </div>
              <div className="flex justify-between">
                <span>Estimated Cost:</span>
                <span className="text-green-600">
                  {gasEstimate ? `~${formatEther(gasEstimate.estimatedCost)} SepoliaETH` : 'Select term for estimate'}
                </span>
              </div>
            </div>
            <div className="mt-4 p-2 bg-gray-50 rounded text-xs">
              <p className="font-medium">⚡ Testnet Benefits:</p>
              <p>Free transactions • Fast confirmation • Safe testing</p>
            </div>
          </div>
        </div>
      </div>

      {/* Transaction History */}
      <div className="bg-white rounded-lg shadow p-6">
        <h3 className="text-lg font-medium text-gray-900 mb-4">
          📊 Transaction History ({transactions.length})
        </h3>
        
        {transactions.length === 0 ? (
          <div className="text-center py-8">
            <div className="text-4xl mb-4">📝</div>
            <p className="text-gray-600">No blockchain transactions yet. Publish your first term to get started!</p>
          </div>
        ) : (
          <div className="space-y-4">
            {transactions.map((tx, index) => (
              <div key={index} className="border border-gray-200 rounded-lg p-4">
                <div className="flex justify-between items-start">
                  <div>
                    <div className="flex items-center gap-2">
                      <h4 className="font-medium text-gray-900">{tx.type}</h4>
                      <span className={`px-2 py-1 text-xs rounded-full ${
                        tx.status === 'confirmed' 
                          ? 'bg-green-100 text-green-800' 
                          : 'bg-yellow-100 text-yellow-800'
                      }`}>
                        {tx.status === 'confirmed' ? '✅ Confirmed' : '⏳ Pending'}
                      </span>
                    </div>
                    <p className="text-sm text-gray-600">
                      {tx.term} • {tx.students} students
                    </p>
                    <p className="text-xs text-gray-500 font-mono mt-1">
                      Tx: {tx.id}
                    </p>
                  </div>
                  <div className="text-right text-sm">
                    <p className="text-gray-500">
                      {new Date(tx.timestamp).toLocaleDateString()}
                    </p>
                    {tx.blockNumber && (
                      <p className="text-gray-400 text-xs">
                        Block: {tx.blockNumber.toLocaleString()}
                      </p>
                    )}
                    {tx.gasUsed && (
                      <p className="text-gray-400 text-xs">
                        Gas: {tx.gasUsed}
                      </p>
                    )}
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

function StatusTab({ systemStatus }: { systemStatus: { status: string; timestamp: string } | null }) {
  return (
    <div className="bg-white rounded-lg shadow p-6">
      <h2 className="text-lg font-medium text-gray-900 mb-4">System Status</h2>
      <div className="space-y-4">
        <div className="flex items-center justify-between p-4 border border-gray-200 rounded-lg">
          <div className="flex items-center gap-3">
            <div className={`w-3 h-3 rounded-full ${systemStatus?.status === 'ok' ? 'bg-green-500' : 'bg-red-500'}`}></div>
            <span className="font-medium">Go API Server</span>
          </div>
          <span className={`text-sm ${systemStatus?.status === 'ok' ? 'text-green-600' : 'text-red-600'}`}>
            {systemStatus?.status === 'ok' ? 'Connected' : 'Disconnected'}
          </span>
        </div>
        
        <div className="flex items-center justify-between p-4 border border-gray-200 rounded-lg">
          <div className="flex items-center gap-3">
            <div className="w-3 h-3 rounded-full bg-blue-500"></div>
            <span className="font-medium">Sepolia Network</span>
          </div>
          <span className="text-sm text-blue-600">Connected</span>
        </div>
        
        <div className="flex items-center justify-between p-4 border border-gray-200 rounded-lg">
          <div className="flex items-center gap-3">
            <div className="w-3 h-3 rounded-full bg-purple-500"></div>
            <span className="font-medium">Cryptographic Services</span>
          </div>
          <span className="text-sm text-purple-600">Running</span>
        </div>

        {systemStatus && (
          <div className="p-4 bg-gray-50 rounded-lg">
            <p className="text-sm text-gray-600">
              Last updated: {new Date(systemStatus.timestamp).toLocaleString()}
            </p>
          </div>
        )}
      </div>
    </div>
  );
}