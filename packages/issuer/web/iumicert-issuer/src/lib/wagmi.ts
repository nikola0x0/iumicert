import { getDefaultConfig } from "connectkit";
import { http, createConfig } from "wagmi";
import { mainnet, sepolia } from "wagmi/chains";

export const config = createConfig(
  getDefaultConfig({
    // Chains our dApp supports
    chains: [mainnet, sepolia],

    // RPC transport mapping
    transports: {
      [mainnet.id]: http(),
      [sepolia.id]: http(),
    },

    // Your actual WalletConnect project ID
    walletConnectProjectId: "afe25c4d6b70033b081c93e3cb146426",

    // Required dApp metadata
    appName: "IU-MiCert Issuer",
    appDescription:
      "Academic credential issuance system with blockchain integration",
    appUrl: "http://localhost:3001",
    appIcon: "/next.svg",
  })
);
