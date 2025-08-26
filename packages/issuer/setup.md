# IU-MiCert Issuer Setup Guide

## 🚀 Quick Setup

### 1. Environment Configuration

Copy the example environment file and customize it:

```bash
cd packages/issuer
cp .env.example .env
```

### 2. Edit `.env` file

Open `.env` in your editor and configure:

```bash
# For local development/testing - these values work out of the box
NETWORK=localhost
ISSUER_PRIVATE_KEY=ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
IUMICERT_CONTRACT_ADDRESS=0x5FbDB2315678afecb367f032d93F642f64180aa3

# For production networks
# NETWORK=sepolia
# ISSUER_PRIVATE_KEY=your_actual_private_key_here
# IUMICERT_CONTRACT_ADDRESS=your_deployed_contract_address
# SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/YOUR_INFURA_KEY
```

## 🔧 Network-Specific Setup

### Local Development (Hardhat/Anvil)

1. **Start local blockchain:**
   ```bash
   # Using Hardhat
   npx hardhat node
   
   # OR using Anvil
   anvil
   ```

2. **Use default test values:**
   ```bash
   NETWORK=localhost
   ISSUER_PRIVATE_KEY=ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
   IUMICERT_CONTRACT_ADDRESS=0x5FbDB2315678afecb367f032d93F642f64180aa3
   ```

### Sepolia Testnet

1. **Get Sepolia ETH:** [Sepolia Faucet](https://sepoliafaucet.com/)

2. **Configure .env:**
   ```bash
   NETWORK=sepolia
   ISSUER_PRIVATE_KEY=your_private_key_here
   IUMICERT_CONTRACT_ADDRESS=deployed_contract_address
   SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/YOUR_INFURA_KEY
   ```

### Mainnet

1. **Configure .env:**
   ```bash
   NETWORK=mainnet
   ISSUER_PRIVATE_KEY=your_private_key_here
   IUMICERT_CONTRACT_ADDRESS=deployed_contract_address
   MAINNET_RPC_URL=https://mainnet.infura.io/v3/YOUR_INFURA_KEY
   ```

## 🔑 Private Key Management

### Generate New Private Key

```bash
# Using OpenSSL
openssl rand -hex 32

# Using Go (create a small script)
go run -c 'package main
import ("crypto/rand"; "fmt"; "encoding/hex")
func main() {
    bytes := make([]byte, 32)
    rand.Read(bytes)
    fmt.Println(hex.EncodeToString(bytes))
}'
```

### Security Notes

- **Never commit `.env` files** (already in `.gitignore`)
- **Use different keys for different networks**
- **Store production keys securely** (consider hardware wallets)
- **Fund your issuer account** with ETH for gas fees

## 🏗️ Contract Deployment

Deploy the IUMiCertRegistry contract to your target network:

```bash
# Deploy to localhost (using Forge)
cd ../contracts
forge create --rpc-url http://localhost:8545 \
    --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
    src/IUMiCertRegistry.sol:IUMiCertRegistry \
    --constructor-args 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266

# Deploy to Sepolia
forge create --rpc-url $SEPOLIA_RPC_URL \
    --private-key $ISSUER_PRIVATE_KEY \
    src/IUMiCertRegistry.sol:IUMiCertRegistry \
    --constructor-args $ISSUER_ADDRESS
```

Update your `.env` file with the deployed contract address.

## ✅ Verify Setup

Test your configuration:

```bash
cd packages/issuer

# Check configuration
go run cmd/main.go --help

# Test blockchain connection (this will show config without publishing)
go run cmd/main.go publish-roots test-term --network localhost
```

## 🔍 Configuration Validation

The system will automatically validate your configuration and provide helpful error messages:

- **Missing private key:** Guides you to set `ISSUER_PRIVATE_KEY`
- **Missing contract address:** Guides you to set `IUMICERT_CONTRACT_ADDRESS` 
- **Invalid network:** Lists supported networks
- **Missing RPC URL:** Reminds you to set Infura keys for testnets/mainnet

## 📁 File Structure

After setup, your structure should look like:

```
packages/issuer/
├── .env                 # Your configuration (not committed)
├── .env.example        # Template for configuration
├── .gitignore          # Excludes .env and sensitive files
├── config/
│   └── env.go          # Configuration loading logic
├── blockchain/
│   ├── client.go       # Ethereum client wrapper
│   ├── contracts.go    # Generated contract bindings
│   └── integration.go  # High-level blockchain operations
└── cmd/
    └── main.go         # CLI application
```

## 🆘 Troubleshooting

### "Private key required"
- Copy `.env.example` to `.env`
- Set `ISSUER_PRIVATE_KEY` in `.env`

### "Contract address required"
- Deploy the contract to your target network
- Set `IUMICERT_CONTRACT_ADDRESS` in `.env`

### "Please set RPC URL with Infura key"
- Sign up at [Infura](https://infura.io)
- Create a project and get your Project ID
- Set `SEPOLIA_RPC_URL` or `MAINNET_RPC_URL` with your key

### "Connection refused" 
- Make sure your local blockchain is running (Hardhat/Anvil)
- Check the RPC URL matches your blockchain

### "Insufficient funds"
- Fund your issuer account with ETH
- Use faucets for testnets
- Buy ETH for mainnet

## 🎯 Next Steps

Once setup is complete:

1. **Generate test data:** `go run cmd/main.go add-term test-term1 path/to/data.json`
2. **Publish to blockchain:** `go run cmd/main.go publish-roots test-term1`  
3. **Generate receipts:** `go run cmd/main.go issue-receipts test-term1`
4. **Verify receipts:** Use the client application to verify issued receipts

The system is now ready for full blockchain integration! 🚀