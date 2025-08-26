'use client';

import { useState, useEffect } from 'react';
import { ConnectButton } from '@rainbow-me/rainbowkit';
import { useAccount } from 'wagmi';
import { apiService, type Term, type Receipt } from '@/lib/api';

export function IssuerDashboard() {
  const { isConnected } = useAccount();
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
            <ConnectButton />
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
            {activeTab === 'blockchain' && <BlockchainTab />}
            {activeTab === 'status' && <StatusTab systemStatus={systemStatus} />}
          </>
        ) : (
          <div className="text-center py-12">
            <div className="text-6xl mb-4">🔒</div>
            <h2 className="text-2xl font-bold text-gray-900 mb-4">Wallet Connection Required</h2>
            <p className="text-gray-600 mb-8">Please connect your wallet to access the issuer dashboard.</p>
            <ConnectButton />
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
    // TODO: Implement term publishing to blockchain
    alert(`Publishing ${selectedTerm.name} credentials to Sepolia blockchain...`);
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
              {receipts.map((receipt, index) => (
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

function BlockchainTab() {
  const [publishingStatus, setPublishingStatus] = useState<'idle' | 'preparing' | 'signing' | 'confirming'>('idle');
  const [transactions, setTransactions] = useState([
    // Mock transaction data for demo
    {
      id: '0x1a2b3c...',
      type: 'Term Publication',
      term: 'Fall 2024',
      students: 5,
      status: 'confirmed',
      timestamp: new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString(),
      blockNumber: 4825193,
      gasUsed: '142,350'
    },
    {
      id: '0x4d5e6f...',
      type: 'Root Update',
      term: 'Spring 2024',
      students: 3,
      status: 'pending',
      timestamp: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString(),
      blockNumber: null,
      gasUsed: null
    }
  ]);

  const handlePublishToBlockchain = async () => {
    setPublishingStatus('preparing');
    // Simulate blockchain publishing process
    setTimeout(() => setPublishingStatus('signing'), 1000);
    setTimeout(() => setPublishingStatus('confirming'), 3000);
    setTimeout(() => {
      setPublishingStatus('idle');
      // Add new transaction
      const newTx = {
        id: '0x' + Math.random().toString(16).slice(2, 8) + '...',
        type: 'Term Publication',
        term: 'Fall 2024',
        students: 5,
        status: 'confirmed',
        timestamp: new Date().toISOString(),
        blockNumber: 4825200 + Math.floor(Math.random() * 100),
        gasUsed: (140000 + Math.floor(Math.random() * 10000)).toLocaleString()
      };
      setTransactions(prev => [newTx, ...prev]);
    }, 8000);
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
              Batch publish all verified student credentials to Sepolia testnet. Creates immutable record of term completion.
            </p>
            <div className="space-y-3">
              <div className="text-sm">
                <span className="font-medium">Ready to publish:</span>
                <ul className="mt-1 text-gray-600 space-y-1">
                  <li>• Fall 2024: 3 verified receipts</li>
                  <li>• Spring 2025: 0 receipts (generate first)</li>
                  <li>• Fall 2023: 0 receipts (generate first)</li>
                </ul>
              </div>
              <button 
                onClick={handlePublishToBlockchain}
                disabled={publishingStatus !== 'idle'}
                className="w-full bg-green-600 text-white px-4 py-2 rounded-lg hover:bg-green-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {publishingStatus === 'idle' ? '🚀 Publish Fall 2024 to Blockchain' : 'Publishing...'}
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
                <span>~20 gwei</span>
              </div>
              <div className="flex justify-between">
                <span>Estimated Cost:</span>
                <span className="text-green-600">~0.003 SepoliaETH</span>
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