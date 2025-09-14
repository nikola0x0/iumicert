# IU-MiCert Issuer System - Project Overview

## Purpose
Academic credential management system for universities using Verkle trees with blockchain integration and zero-knowledge proofs. Enables privacy-preserving verification of student achievements.

## Technology Stack
- **Backend**: Go 1.21+ with Cobra CLI framework
- **Frontend**: Next.js 15 with React 18, TypeScript, TailwindCSS
- **Blockchain**: wagmi v2, viem v2, ConnectKit for Web3 integration
- **Cryptography**: Ethereum's go-verkle library for 32-byte proofs
- **Blockchain Network**: Ethereum Sepolia testnet
- **API**: REST API server on port 8080 with CORS support

## Architecture
- Single Verkle tree per academic term containing all student course completions
- CLI application (`micert` binary) with comprehensive commands
- REST API backend serving Next.js frontend
- Smart contracts deployed: IUMiCertRegistry (0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60)

## Current State
- Backend: Fully implemented CLI with all commands working
- API Server: Running on port 8080, healthy status confirmed
- Frontend: Next.js app with Web3 integration (wagmi/viem)
- Data: Test data exists, system operational