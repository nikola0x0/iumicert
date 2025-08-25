import { ethers } from "hardhat";

async function main() {
  console.log("🚀 Deploying IU-MiCert Registry Contract...");
  
  const [deployer] = await ethers.getSigners();
  
  console.log("📍 Deploying contracts with account:", deployer.address);
  console.log("💰 Account balance:", (await deployer.getBalance()).toString());
  
  // Deploy IUMiCertRegistry
  const IUMiCertRegistry = await ethers.getContractFactory("IUMiCertRegistry");
  const registry = await IUMiCertRegistry.deploy();
  await registry.deployed();
  
  console.log("✅ IUMiCertRegistry deployed to:", registry.address);
  
  // Deploy IUMiCertVerifier with registry address
  const IUMiCertVerifier = await ethers.getContractFactory("IUMiCertVerifier");
  const verifier = await IUMiCertVerifier.deploy(registry.address);
  await verifier.deployed();
  
  console.log("✅ IUMiCertVerifier deployed to:", verifier.address);
  
  // Save deployment addresses
  const deploymentInfo = {
    network: await ethers.provider.getNetwork(),
    registry: registry.address,
    verifier: verifier.address,
    deployer: deployer.address,
    timestamp: new Date().toISOString(),
  };
  
  console.log("📝 Deployment Summary:", deploymentInfo);
  
  // Verify contracts on Etherscan (if not local)
  if (deploymentInfo.network.chainId !== 31337) {
    console.log("⏳ Waiting for block confirmations...");
    await registry.deployTransaction.wait(6);
    await verifier.deployTransaction.wait(6);
    
    console.log("🔍 Verifying contracts on Etherscan...");
    try {
      await run("verify:verify", {
        address: registry.address,
        constructorArguments: [],
      });
      
      await run("verify:verify", {
        address: verifier.address,
        constructorArguments: [registry.address],
      });
      
      console.log("✅ Contracts verified on Etherscan");
    } catch (error) {
      console.log("⚠️ Verification failed:", error);
    }
  }
  
  return deploymentInfo;
}

main()
  .then((deployment) => {
    console.log("🎉 Deployment completed successfully!");
    process.exit(0);
  })
  .catch((error) => {
    console.error("❌ Deployment failed:", error);
    process.exit(1);
  });