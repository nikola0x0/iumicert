package blockchain

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// BlockchainClient manages Ethereum blockchain interactions
type BlockchainClient struct {
	client     *ethclient.Client
	privateKey *ecdsa.PrivateKey
	chainID    *big.Int
	gasLimit   uint64
}

// NetworkConfig holds blockchain network configuration
type NetworkConfig struct {
	RPC_URL     string
	ChainID     *big.Int
	GasLimit    uint64
	PrivateKey  string
}

// GetNetworkConfig returns network configuration for supported networks
func GetNetworkConfig(network string) (*NetworkConfig, error) {
	switch network {
	case "sepolia":
		return &NetworkConfig{
			RPC_URL:    "https://sepolia.infura.io/v3/YOUR_INFURA_KEY", // User will need to set this
			ChainID:    big.NewInt(11155111), // Sepolia chain ID
			GasLimit:   500000,
		}, nil
	case "mainnet":
		return &NetworkConfig{
			RPC_URL:    "https://mainnet.infura.io/v3/YOUR_INFURA_KEY",
			ChainID:    big.NewInt(1), // Mainnet chain ID
			GasLimit:   500000,
		}, nil
	case "localhost":
		return &NetworkConfig{
			RPC_URL:    "http://localhost:8545",
			ChainID:    big.NewInt(1337), // Local development chain ID
			GasLimit:   500000,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported network: %s", network)
	}
}

// NewBlockchainClient creates a new blockchain client
func NewBlockchainClient(network, privateKeyHex string) (*BlockchainClient, error) {
	// Get network configuration
	config, err := GetNetworkConfig(network)
	if err != nil {
		return nil, fmt.Errorf("failed to get network config: %w", err)
	}

	// Connect to Ethereum client
	client, err := ethclient.Dial(config.RPC_URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	// Parse private key
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return &BlockchainClient{
		client:     client,
		privateKey: privateKey,
		chainID:    config.ChainID,
		gasLimit:   config.GasLimit,
	}, nil
}

// GetTransactOpts returns transaction options for contract calls
func (bc *BlockchainClient) GetTransactOpts(ctx context.Context) (*bind.TransactOpts, error) {
	// Get account nonce
	publicKey := bc.privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := bc.client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get account nonce: %w", err)
	}

	// Get suggested gas price
	gasPrice, err := bc.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas price: %w", err)
	}

	// Create transaction options
	auth, err := bind.NewKeyedTransactorWithChainID(bc.privateKey, bc.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create keyed transactor: %w", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = bc.gasLimit
	auth.GasPrice = gasPrice

	return auth, nil
}

// GetCallOpts returns call options for read-only contract calls
func (bc *BlockchainClient) GetCallOpts(ctx context.Context) *bind.CallOpts {
	return &bind.CallOpts{
		Context: ctx,
	}
}

// GetClient returns the underlying Ethereum client
func (bc *BlockchainClient) GetClient() *ethclient.Client {
	return bc.client
}

// Close closes the blockchain client connection
func (bc *BlockchainClient) Close() {
	bc.client.Close()
}

// WaitForTransaction waits for a transaction to be mined and returns the receipt
func (bc *BlockchainClient) WaitForTransaction(ctx context.Context, txHash common.Hash) error {
	// In newer versions of go-ethereum, we can directly use the client to wait for receipt
	_, err := bc.client.TransactionReceipt(ctx, txHash)
	if err != nil {
		return fmt.Errorf("failed to get transaction receipt: %w", err)
	}
	return nil
}