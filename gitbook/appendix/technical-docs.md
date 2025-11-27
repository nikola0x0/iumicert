# Technical Documentation

This GitBook provides a high-level, thesis-focused overview of IU-MiCert. For developers and technical implementers, detailed documentation is available in the project repository.

## Primary Technical Resources

### 1. Issuer System Documentation

**Location**: `packages/issuer/README.md`

**Contents**:
- Complete CLI command reference (15+ commands)
- API endpoint documentation
- Database schema and migrations
- Deployment instructions
- Development setup guide

**Audience**: Developers implementing or extending the issuer system

### 2. Cryptographic Implementation

**Location**: `docs/implementation-guide.md`

**Contents**:
- IPA verification implementation details
- go-verkle and go-ipa library usage
- Proof generation step-by-step
- Code examples with file references
- Security implementation details

**Audience**: Cryptography researchers, security auditors

### 3. Mathematical Foundation

**Location**: `docs/mathematical-foundation.md`

**Contents**:
- Verkle tree mathematical theory
- Polynomial commitment schemes
- Inner Product Argument proofs
- Security analysis and proofs
- Performance characteristics

**Audience**: Researchers, academics, theoretical cryptographers

### 4. Data Flow Documentation

**Location**: `packages/issuer/docs/data-flow.md`

**Contents**:
- Complete pipeline from LMS to blockchain
- Data format specifications
- File structure and organization
- Processing workflows
- Error handling

**Audience**: System integrators, data engineers

### 5. Smart Contract Documentation

**Location**: `packages/contracts/README.md`

**Contents**:
- Solidity contract source code
- Deployment scripts
- Testing framework (Foundry)
- Gas optimization details
- Security considerations

**Audience**: Smart contract developers, auditors

### 6. Client Portal Documentation

**Location**: `packages/client/README.md`

**Contents**:
- Next.js application structure
- Verification workflow implementation
- API integration
- UI/UX design decisions
- Deployment guide

**Audience**: Frontend developers, UI/UX designers

## Code Repository

**GitHub**: [Your Repository URL]

**Structure**:
```
iumicert/
├── docs/                    # Academic documentation (this GitBook source)
├── packages/
│   ├── issuer/             # Go backend + CLI
│   ├── client/             # Next.js verification portal
│   └── contracts/          # Solidity smart contracts
└── README.md               # Project overview
```

## Live Deployments

**Smart Contracts**:
- IUMiCertRegistry: [`0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60`](https://sepolia.etherscan.io/address/0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60)
- Network: Ethereum Sepolia Testnet

**Web Applications**:
- Student Portal: [https://iu-micert.vercel.app](https://iu-micert.vercel.app)
- Issuer Dashboard: [https://iumicert-issuer.vercel.app](https://iumicert-issuer.vercel.app)

## Development Setup

For developers wanting to run IU-MiCert locally:

1. **Clone repository**
2. **Install dependencies**:
   - Go 1.21+
   - Node.js 18+
   - PostgreSQL 15+
   - Docker (optional)
3. **Follow setup guides** in respective package READMEs
4. **See** `packages/issuer/README.md` for detailed instructions

## Technical Support

For technical questions, implementation help, or collaboration:

- **GitHub Issues**: [Repository Issues]
- **Documentation**: Package-specific READMEs
- **Academic Inquiries**: [Your Contact Email]

## Research Papers

Detailed research findings and academic contributions:

1. **Thesis Document**: [Link to full thesis PDF]
2. **Conference Papers**: [Links if published]
3. **Technical Reports**: Available in `docs/` directory

## API Reference

**Issuer API**: `http://localhost:8080` (development)

**Endpoints**:
- `POST /api/demo/generate-full` - Generate test data
- `POST /api/demo/reset` - Reset system
- `GET /api/receipts/:studentId` - Get student receipts
- `POST /api/verify` - Verify receipt

**Full API docs**: See `packages/issuer/README.md`

## Contributing

Contributions welcome! See:
- `CONTRIBUTING.md` (if exists)
- Individual package READMEs
- GitHub Issues for feature requests

---

**Note**: This GitBook focuses on research contributions and high-level concepts. Technical implementation details are maintained in the GitHub repository for easier updates and version control.
