# IU-MiCert Issuer Documentation

Documentation for the IU-MiCert issuer system (backend API, CLI, and admin dashboard).

---

## 📚 Core Documentation

### 1. **[setup.md](./setup.md)**

Initial setup guide for local development.

**Topics**:

- Prerequisites and system requirements
- Installing dependencies (Go, Docker, PostgreSQL)
- Environment configuration
- Building the CLI binary
- Running locally

**Read this first** for getting started with development.

---

### 2. **[data-flow.md](./data-flow.md)**

Complete system architecture and data pipeline.

**Topics**:

- System architecture overview
- Data flow from LMS to blockchain
- Verkle tree generation process
- Receipt creation workflow
- Component interactions

**Read this** for understanding how the system works end-to-end.

---

### 3. **[deployment.md](./deployment.md)**

Production deployment guide for GCP.

**Topics**:

- GCP VM setup
- Docker deployment
- Firewall configuration
- Environment variables
- Security best practices
- Monitoring and maintenance

**Read this** for deploying to production.

---

### 4. **[demo-guide.md](./demo-guide.md)**

Demonstration script for showcasing the system.

**Topics**:

- Complete demo flow
- CLI command examples
- What to highlight at each step
- Expected outputs

**Read this** for preparing demos or presentations.

---

## 🗄️ Archived Documents

The `archive/` directory contains historical and outdated documentation:

- **authentication-setup.md** - Auth configuration (now covered in deployment.md)
- **ci-cd-setup.md** - GitHub Actions CI/CD (advanced deployment topic)
- **thesis-demo-flow.md** - Thesis-specific demo (superseded by demo-guide.md)
- **ipa-verification-debugging.md** - Old debugging notes (issue resolved)

These are kept for historical reference but should not be used for current development.

---

## 🚀 Quick Start

### For Development

1. Read [setup.md](./setup.md)
2. Configure your environment
3. Run `./dev.sh` to start local server
4. Read [data-flow.md](./data-flow.md) to understand the system

### For Deployment

1. Read [deployment.md](./deployment.md)
2. Set up GCP VM or use Docker locally
3. Configure production environment variables
4. Deploy and monitor

### For Demos

1. Read [demo-guide.md](./demo-guide.md)
2. Ensure test data is generated (`./generate.sh`)
3. Follow the demo script
4. Highlight key features

---

## 🔗 Related Documentation

**Main System:**

- **Issuer README**: `../README.md` - Complete system overview and CLI commands
- **Project README**: `../../../README.md` - Thesis overview and research contributions
- **Technical Docs**: `../../../docs/` - Cryptographic implementation details

**Web Dashboard:**

- **Dashboard README**: `../web/iumicert-issuer/README.md` - Admin dashboard documentation
- **Dashboard Docs**: `../web/iumicert-issuer/docs/` - Design system and deployment

**Other Components:**

- **Client Portal**: `../../client/README.md` - Student/verifier portal
- **Smart Contracts**: `../../contracts/README.md` - Solidity contracts

---

## 📝 Document Maintenance

### When to Update

- **setup.md**: When dependencies or setup steps change
- **data-flow.md**: When system architecture or data pipeline changes
- **deployment.md**: When deployment process or infrastructure changes
- **demo-guide.md**: When adding new demo scenarios or updating commands
- **This README**: When documents are added, renamed, or archived

### Adding New Documentation

1. Create document with clear, descriptive name (use kebab-case)
2. Add entry to this README with purpose and topics covered
3. Link from relevant sections in main issuer README
4. Archive outdated docs to `archive/` directory

---

**Active Docs**: 4 core documents + archived references
