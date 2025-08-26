import {
  createPublicClient,
  createWalletClient,
  custom,
  http,
  parseAbi,
} from "viem";
import { sepolia } from "viem/chains";

// IUMiCertRegistry contract ABI (key functions)
const IUMICERT_ABI = parseAbi([
  "function publishTermRoot(bytes32 _verkleRoot, string _termId, uint256 _totalStudents) external",
  "function verifyReceiptAnchor(bytes32 _blockchainAnchor) external view returns (bool isValid, string termId, uint256 publishedAt)",
  "function getTermRootInfo(bytes32 _verkleRoot) external view returns (string termId, uint256 totalStudents, uint256 publishedAt, bool exists)",
  "function getPublishedRootsCount() external view returns (uint256)",
  "function owner() external view returns (address)",
  "event TermRootPublished(bytes32 indexed verkleRoot, string indexed termId, uint256 totalStudents, uint256 timestamp)",
]);

// Contract address - this should be set via environment variable
const CONTRACT_ADDRESS =
  (process.env.NEXT_PUBLIC_IUMICERT_CONTRACT_ADDRESS as `0x${string}`) ||
  "0x5FbDB2315678afecb367f032d93F642f64180aa3"; // Fallback for localhost

// Public client for reading from blockchain
const publicClient = createPublicClient({
  chain: sepolia,
  transport: http(),
});

export interface TermRootData {
  term_id: string;
  verkle_root: string;
  total_students: number;
}

export interface PublishResult {
  transactionHash: `0x${string}`;
  blockNumber: bigint;
  gasUsed: bigint;
  status: "success" | "failed";
}

/**
 * Publishes a term root to the IUMiCertRegistry contract via MetaMask
 * This will trigger a MetaMask popup for user signature
 */
export async function publishTermRoot(
  termRootData: TermRootData,
  connectedAddress?: `0x${string}`,
  connectedWalletClient?: any
): Promise<PublishResult> {
  console.log("🚀 Publishing term root to blockchain:", termRootData);
  console.log("📥 Received params:", { connectedAddress, hasConnectedWalletClient: !!connectedWalletClient });
  
  // Use provided wagmi account and wallet client if available
  let account = connectedAddress;
  let walletClient = connectedWalletClient;
  
  // If no wagmi client provided, fall back to direct MetaMask access
  if (!account || !walletClient) {
    console.log("⚠️  No wagmi params provided, falling back to direct MetaMask access");
    
    // Check if MetaMask is available
    if (!window.ethereum) {
      throw new Error("MetaMask not found. Please install MetaMask to continue.");
    }

    console.log("✅ MetaMask detected:", !!window.ethereum);

    // Create wallet client for transactions
    walletClient = createWalletClient({
      chain: sepolia,
      transport: custom(window.ethereum),
    });

    console.log("✅ Fallback wallet client created");

    // Get connected account - try multiple approaches
    try {
      console.log("📞 Trying walletClient.getAddresses()...");
      const addresses = await walletClient.getAddresses();
      account = addresses[0];
      console.log("✅ Got addresses from wallet client:", addresses);
    } catch (error) {
      console.warn("❌ Failed to get addresses from wallet client:", error);
      
      // Fallback: request account access
      try {
        console.log("📞 Trying window.ethereum.request('eth_requestAccounts')...");
        const accounts = await window.ethereum.request({ 
          method: 'eth_requestAccounts' 
        }) as string[];
        account = accounts[0] as `0x${string}`;
        console.log("✅ Got accounts from direct MetaMask:", accounts);
      } catch (requestError) {
        console.error("❌ Failed to get accounts from MetaMask:", requestError);
        throw new Error("Failed to connect wallet. Please make sure MetaMask is unlocked and connected to this site.");
      }
    }
  } else {
    console.log("✅ Using wagmi wallet client and address");
  }
  
  console.log("🔍 Final account:", account);
  
  if (!account) {
    throw new Error("No wallet account connected. Please connect your wallet and try again.");
  }

  // Parse verkle root to bytes32
  let verkleRootHex = termRootData.verkle_root;
  if (!verkleRootHex.startsWith("0x")) {
    verkleRootHex = "0x" + verkleRootHex;
  }

  if (verkleRootHex.length !== 66) {
    // 0x + 64 hex chars = 66 total
    throw new Error(
      `Invalid verkle root format. Expected 32 bytes (66 chars with 0x), got ${verkleRootHex.length}`
    );
  }

  // Simulate the transaction first to catch errors
  const { request } = await publicClient.simulateContract({
    account,
    address: CONTRACT_ADDRESS,
    abi: IUMICERT_ABI,
    functionName: "publishTermRoot",
    args: [
      verkleRootHex as `0x${string}`,
      termRootData.term_id,
      BigInt(termRootData.total_students),
    ],
  });

  // Execute the transaction (this will trigger MetaMask popup)
  const hash = await walletClient.writeContract(request);

  // Wait for transaction confirmation
  const receipt = await publicClient.waitForTransactionReceipt({ hash });

  return {
    transactionHash: hash,
    blockNumber: receipt.blockNumber,
    gasUsed: receipt.gasUsed,
    status: receipt.status === "success" ? "success" : "failed",
  };
}

/**
 * Waits for transaction confirmation with status updates
 */
export async function waitForTransactionConfirmation(
  hash: `0x${string}`,
  onStatusUpdate?: (status: "pending" | "confirmed" | "failed") => void
): Promise<{
  blockNumber: bigint;
  gasUsed: bigint;
  status: "success" | "reverted";
}> {
  onStatusUpdate?.("pending");

  try {
    const receipt = await publicClient.waitForTransactionReceipt({
      hash,
      timeout: 60_000, // 60 seconds timeout
    });

    const status = receipt.status === "success" ? "confirmed" : "failed";
    onStatusUpdate?.(status as any);

    return {
      blockNumber: receipt.blockNumber,
      gasUsed: receipt.gasUsed,
      status: receipt.status,
    };
  } catch (error) {
    onStatusUpdate?.("failed");
    throw error;
  }
}

/**
 * Fetches term root information from the blockchain
 */
export async function getTermRootInfo(verkleRootHex: string) {
  let rootHex = verkleRootHex;
  if (!rootHex.startsWith("0x")) {
    rootHex = "0x" + rootHex;
  }

  const result = await publicClient.readContract({
    address: CONTRACT_ADDRESS,
    abi: IUMICERT_ABI,
    functionName: "getTermRootInfo",
    args: [rootHex as `0x${string}`],
  });

  return {
    termId: result[0] as string,
    totalStudents: Number(result[1]),
    publishedAt: Number(result[2]),
    exists: result[3] as boolean,
  };
}

/**
 * Verifies if a receipt's blockchain anchor is valid
 */
export async function verifyReceiptAnchor(blockchainAnchorHex: string) {
  let anchorHex = blockchainAnchorHex;
  if (!anchorHex.startsWith("0x")) {
    anchorHex = "0x" + anchorHex;
  }

  const result = await publicClient.readContract({
    address: CONTRACT_ADDRESS,
    abi: IUMICERT_ABI,
    functionName: "verifyReceiptAnchor",
    args: [anchorHex as `0x${string}`],
  });

  return {
    isValid: result[0] as boolean,
    termId: result[1] as string,
    publishedAt: Number(result[2]),
  };
}

/**
 * Gets the total number of published roots
 */
export async function getPublishedRootsCount(): Promise<number> {
  const count = await publicClient.readContract({
    address: CONTRACT_ADDRESS,
    abi: IUMICERT_ABI,
    functionName: "getPublishedRootsCount",
  });

  return Number(count);
}

/**
 * Gets the contract owner address
 */
export async function getContractOwner(): Promise<`0x${string}`> {
  const owner = await publicClient.readContract({
    address: CONTRACT_ADDRESS,
    abi: IUMICERT_ABI,
    functionName: "owner",
  });

  return owner as `0x${string}`;
}

/**
 * Gets recent TermRootPublished events from the blockchain
 */
export async function getTermRootHistory(fromBlock?: bigint) {
  // Get current block number
  const currentBlock = await publicClient.getBlockNumber();
  // Look back 1000 blocks to avoid "ranges over 10000 blocks" error on free tier
  const blockOffset = BigInt(1000);
  const startBlock = fromBlock || (currentBlock > blockOffset ? currentBlock - blockOffset : BigInt(0));
  
  const logs = await publicClient.getContractEvents({
    address: CONTRACT_ADDRESS,
    abi: IUMICERT_ABI,
    eventName: "TermRootPublished",
    fromBlock: startBlock,
    toBlock: "latest",
  });

  return logs.map((log) => ({
    verkleRoot: log.args.verkleRoot as string,
    termId: log.args.termId as string,
    totalStudents: Number(log.args.totalStudents),
    timestamp: Number(log.args.timestamp),
    transactionHash: log.transactionHash,
    blockNumber: Number(log.blockNumber),
  }));
}

/**
 * Estimates gas cost for publishing a term root
 */
export async function estimatePublishGas(
  termRootData: TermRootData,
  connectedAddress?: `0x${string}`,
  connectedWalletClient?: any
): Promise<{
  gasLimit: bigint;
  gasPrice: bigint;
  estimatedCost: bigint;
}> {
  console.log("⛽ Starting gas estimation...");
  console.log("📥 Gas estimation received params:", { connectedAddress, hasConnectedWalletClient: !!connectedWalletClient });
  
  // Use provided wagmi account if available
  let account = connectedAddress;
  
  // If no wagmi account provided, fall back to direct MetaMask access
  if (!account) {
    console.log("⚠️  Gas estimation: no wagmi address provided, falling back to direct MetaMask access");
    
    // Check if MetaMask is available
    if (!window.ethereum) {
      throw new Error("MetaMask not found");
    }

    console.log("✅ MetaMask detected for gas estimation:", !!window.ethereum);

    const walletClient = createWalletClient({
      chain: sepolia,
      transport: custom(window.ethereum),
    });

    console.log("✅ Gas estimation wallet client created");

    // Get connected account
    try {
      console.log("📞 Gas estimation: trying walletClient.getAddresses()...");
      const addresses = await walletClient.getAddresses();
      account = addresses[0];
      console.log("✅ Gas estimation: got addresses:", addresses);
    } catch (error) {
      console.warn("❌ Gas estimation: failed to get addresses from wallet client:", error);
      
      // Fallback: request account access
      try {
        console.log("📞 Gas estimation: trying window.ethereum.request('eth_requestAccounts')...");
        const accounts = await window.ethereum.request({ 
          method: 'eth_requestAccounts' 
        }) as string[];
        account = accounts[0] as `0x${string}`;
        console.log("✅ Gas estimation: got accounts from direct MetaMask:", accounts);
      } catch (requestError) {
        console.error("❌ Gas estimation: failed to get accounts from MetaMask:", requestError);
        throw new Error("Failed to connect wallet for gas estimation");
      }
    }
  } else {
    console.log("✅ Gas estimation: using wagmi address");
  }
  
  console.log("🔍 Gas estimation: final account:", account);
  
  if (!account) {
    throw new Error("No wallet account connected");
  }

  let verkleRootHex = termRootData.verkle_root;
  if (!verkleRootHex.startsWith("0x")) {
    verkleRootHex = "0x" + verkleRootHex;
  }

  // Estimate gas
  const gasLimit = await publicClient.estimateContractGas({
    account,
    address: CONTRACT_ADDRESS,
    abi: IUMICERT_ABI,
    functionName: "publishTermRoot",
    args: [
      verkleRootHex as `0x${string}`,
      termRootData.term_id,
      BigInt(termRootData.total_students),
    ],
  });

  // Get current gas price
  const gasPrice = await publicClient.getGasPrice();

  const estimatedCost = gasLimit * gasPrice;

  return {
    gasLimit,
    gasPrice,
    estimatedCost,
  };
}

/**
 * Formats wei to ETH for display
 */
export function formatEther(wei: bigint): string {
  const eth = Number(wei) / 1e18;
  return eth.toFixed(6);
}

/**
 * Checks if the connected account is the contract owner
 */
export async function isContractOwner(
  account: `0x${string}`
): Promise<boolean> {
  try {
    const owner = await getContractOwner();
    return owner.toLowerCase() === account.toLowerCase();
  } catch {
    return false;
  }
}
