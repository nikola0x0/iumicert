package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the issuer
type Config struct {
	// Blockchain settings
	IssuerPrivateKey     string
	ContractAddress      string
	Network              string
	DefaultGasLimit      uint64
	MaxGasPrice          uint64
	
	// RPC URLs
	LocalhostRPCURL      string
	SepoliaRPCURL        string
	MainnetRPCURL        string
	
	// Application settings
	Debug                bool
	LogLevel             string
	
	// Test settings
	TestPrivateKey       string
	TestContractAddress  string
}

// LoadConfig loads configuration from .env file and environment variables
func LoadConfig() (*Config, error) {
	// Try to load .env file from current directory or parent directories
	if err := loadDotEnv(); err != nil {
		// .env file is optional, so we don't fail here
		fmt.Printf("Note: .env file not found, using environment variables only\n")
	}

	config := &Config{
		// Blockchain settings
		IssuerPrivateKey:    getEnv("ISSUER_PRIVATE_KEY", ""),
		ContractAddress:     getEnv("IUMICERT_CONTRACT_ADDRESS", ""),
		Network:             getEnv("NETWORK", "localhost"),
		DefaultGasLimit:     getEnvUint64("DEFAULT_GAS_LIMIT", 500000),
		MaxGasPrice:         getEnvUint64("MAX_GAS_PRICE", 20000000000),
		
		// RPC URLs
		LocalhostRPCURL:     getEnv("LOCALHOST_RPC_URL", "http://localhost:8545"),
		SepoliaRPCURL:       getEnv("SEPOLIA_RPC_URL", "https://sepolia.infura.io/v3/YOUR_INFURA_KEY"),
		MainnetRPCURL:       getEnv("MAINNET_RPC_URL", "https://mainnet.infura.io/v3/YOUR_INFURA_KEY"),
		
		// Application settings
		Debug:               getEnvBool("DEBUG", false),
		LogLevel:            getEnv("LOG_LEVEL", "info"),
		
		// Test settings (fallback for development)
		TestPrivateKey:      getEnv("TEST_PRIVATE_KEY", "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"),
		TestContractAddress: getEnv("TEST_CONTRACT_ADDRESS", "0x5FbDB2315678afecb367f032d93F642f64180aa3"),
	}

	return config, nil
}

// GetRPCURL returns the appropriate RPC URL for the given network
func (c *Config) GetRPCURL(network string) (string, error) {
	switch network {
	case "localhost", "local":
		return c.LocalhostRPCURL, nil
	case "sepolia":
		if c.SepoliaRPCURL == "https://sepolia.infura.io/v3/YOUR_INFURA_KEY" {
			return "", fmt.Errorf("please set SEPOLIA_RPC_URL with your Infura key")
		}
		return c.SepoliaRPCURL, nil
	case "mainnet":
		if c.MainnetRPCURL == "https://mainnet.infura.io/v3/YOUR_INFURA_KEY" {
			return "", fmt.Errorf("please set MAINNET_RPC_URL with your Infura key")
		}
		return c.MainnetRPCURL, nil
	default:
		return "", fmt.Errorf("unsupported network: %s", network)
	}
}

// GetPrivateKey returns the private key to use, with fallback to test key for localhost
func (c *Config) GetPrivateKey() string {
	if c.IssuerPrivateKey != "" {
		return c.IssuerPrivateKey
	}
	
	// For localhost/testing, use the test private key
	if c.Network == "localhost" || c.Network == "local" {
		fmt.Println("⚠️  Using test private key for localhost development")
		return c.TestPrivateKey
	}
	
	return ""
}

// GetContractAddress returns the contract address to use, with fallback to test address for localhost
func (c *Config) GetContractAddress() string {
	if c.ContractAddress != "" {
		return c.ContractAddress
	}
	
	// For localhost/testing, use the test contract address
	if c.Network == "localhost" || c.Network == "local" {
		fmt.Println("⚠️  Using test contract address for localhost development")
		return c.TestContractAddress
	}
	
	return ""
}

// Validate checks if all required configuration is present
func (c *Config) Validate() error {
	privateKey := c.GetPrivateKey()
	if privateKey == "" {
		return fmt.Errorf("private key is required. Set ISSUER_PRIVATE_KEY environment variable")
	}
	
	contractAddress := c.GetContractAddress()
	if contractAddress == "" {
		return fmt.Errorf("contract address is required. Set IUMICERT_CONTRACT_ADDRESS environment variable")
	}
	
	_, err := c.GetRPCURL(c.Network)
	if err != nil {
		return fmt.Errorf("invalid network configuration: %w", err)
	}
	
	return nil
}

// loadDotEnv tries to load .env file from current directory or parent directories
func loadDotEnv() error {
	// Try current directory first
	if err := godotenv.Load(); err == nil {
		return nil
	}
	
	// Try parent directories
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	
	for i := 0; i < 5; i++ { // Search up to 5 levels up
		envPath := filepath.Join(currentDir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			return godotenv.Load(envPath)
		}
		
		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			break // Reached root directory
		}
		currentDir = parent
	}
	
	return fmt.Errorf(".env file not found")
}

// Helper functions for environment variable parsing

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvUint64(key string, defaultValue uint64) uint64 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseUint(value, 10, 64); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// PrintConfig prints the current configuration (excluding sensitive data)
func (c *Config) PrintConfig() {
	fmt.Println("📋 Current Configuration:")
	fmt.Printf("  Network: %s\n", c.Network)
	fmt.Printf("  Contract Address: %s\n", maskSensitive(c.GetContractAddress()))
	fmt.Printf("  Private Key: %s\n", maskPrivateKey(c.GetPrivateKey()))
	fmt.Printf("  Default Gas Limit: %d\n", c.DefaultGasLimit)
	fmt.Printf("  Debug Mode: %v\n", c.Debug)
	fmt.Printf("  Log Level: %s\n", c.LogLevel)
}

func maskSensitive(value string) string {
	if len(value) <= 8 {
		return strings.Repeat("*", len(value))
	}
	return value[:4] + strings.Repeat("*", len(value)-8) + value[len(value)-4:]
}

func maskPrivateKey(value string) string {
	if value == "" {
		return "not set"
	}
	if len(value) <= 8 {
		return strings.Repeat("*", len(value))
	}
	return value[:4] + strings.Repeat("*", len(value)-8) + value[len(value)-4:]
}