#!/bin/bash

# IU-MiCert Issuer - GCP VM Initial Setup Script
# Run this once on your VM to prepare it for automated deployments

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_step() {
    echo -e "${BLUE}==>${NC} $1"
}

log_step "Starting IU-MiCert Issuer VM Setup"

# Update system
log_info "Updating system packages..."
sudo apt-get update
sudo apt-get upgrade -y

# Install Docker
if ! command -v docker &> /dev/null; then
    log_info "Installing Docker..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    rm get-docker.sh

    # Add current user to docker group
    sudo usermod -aG docker $USER
    log_warn "You need to log out and log back in for docker group changes to take effect"
else
    log_info "Docker already installed ($(docker --version))"
fi

# Install Docker Compose
if ! docker compose version &> /dev/null; then
    log_info "Installing Docker Compose..."
    sudo apt-get install -y docker-compose-plugin
else
    log_info "Docker Compose already installed ($(docker compose version))"
fi

# Install Git
if ! command -v git &> /dev/null; then
    log_info "Installing Git..."
    sudo apt-get install -y git
else
    log_info "Git already installed ($(git --version))"
fi

# Install additional utilities
log_info "Installing utilities..."
sudo apt-get install -y curl wget htop nano vim

# Create project directory structure
log_info "Creating project directory structure..."
mkdir -p ~/iumicert
cd ~/iumicert

# Clone repository (optional - can be done by CI/CD)
read -p "Do you want to clone the repository now? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    read -p "Enter repository URL: " REPO_URL
    if [ -n "$REPO_URL" ]; then
        git clone "$REPO_URL" .
        log_info "Repository cloned successfully"
    fi
fi

# Configure firewall (using ufw)
log_info "Configuring firewall..."
if ! command -v ufw &> /dev/null; then
    sudo apt-get install -y ufw
fi

# Allow SSH
sudo ufw allow 22/tcp

# Allow API port
sudo ufw allow 8080/tcp

# Enable firewall (only if not already enabled)
if ! sudo ufw status | grep -q "Status: active"; then
    log_warn "Firewall will be enabled. Make sure SSH access is configured!"
    read -p "Enable firewall now? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "y" | sudo ufw enable
    fi
else
    log_info "Firewall already active"
fi

# Set up log rotation
log_info "Setting up log rotation..."
sudo tee /etc/logrotate.d/iumicert > /dev/null << 'EOF'
/home/*/iumicert/packages/issuer/logs/*.log {
    daily
    rotate 14
    compress
    delaycompress
    notifempty
    missingok
    create 0640 $USER $USER
}
EOF

# Create systemd service for auto-restart (optional)
log_step "Creating systemd service for auto-restart..."
sudo tee /etc/systemd/system/iumicert-issuer.service > /dev/null << EOF
[Unit]
Description=IU-MiCert Issuer Service
After=docker.service
Requires=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=$HOME/iumicert/packages/issuer
ExecStart=/usr/bin/docker compose -f docker-compose.prod.yml up -d
ExecStop=/usr/bin/docker compose -f docker-compose.prod.yml down
User=$USER

[Install]
WantedBy=multi-user.target
EOF

log_info "Systemd service created (not enabled by default)"
log_info "To enable on boot: sudo systemctl enable iumicert-issuer"

# Set up monitoring script
log_info "Creating monitoring script..."
cat > ~/check-health.sh << 'EOF'
#!/bin/bash
# Quick health check script

echo "=== IU-MiCert Issuer Health Check ==="
echo ""

# Check if containers are running
echo "Container Status:"
cd ~/iumicert/packages/issuer
docker compose -f docker-compose.prod.yml ps

echo ""
echo "API Health Check:"
if curl -f http://localhost:8080/api/health 2>/dev/null; then
    echo "✅ API is healthy"
else
    echo "❌ API is not responding"
fi

echo ""
echo "Disk Usage:"
df -h / | tail -1

echo ""
echo "Memory Usage:"
free -h | grep Mem

echo ""
echo "Recent Logs (last 10 lines):"
docker compose -f docker-compose.prod.yml logs --tail=10 issuer-backend
EOF

chmod +x ~/check-health.sh

# Print summary
log_step "Setup Complete!"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📋 Summary"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "✅ Docker installed"
echo "✅ Docker Compose installed"
echo "✅ Git installed"
echo "✅ Firewall configured (ports 22, 8080)"
echo "✅ Log rotation configured"
echo "✅ Systemd service created"
echo "✅ Health check script created"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📝 Next Steps"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "1. Log out and log back in for Docker group changes:"
echo "   exit"
echo ""
echo "2. Configure GitHub Actions secrets in your repository:"
echo "   - GCP_SA_KEY (Service account JSON)"
echo "   - GCP_VM_NAME (Your VM name)"
echo "   - GCP_VM_ZONE (Your VM zone)"
echo "   - ISSUER_PRIVATE_KEY (Ethereum private key)"
echo "   - SEPOLIA_RPC_URL (Infura URL)"
echo "   - POSTGRES_PASSWORD (Database password)"
echo "   - IUMICERT_CONTRACT_ADDRESS (Contract address)"
echo ""
echo "3. Push to your repository to trigger automatic deployment"
echo ""
echo "4. Check deployment health:"
echo "   ~/check-health.sh"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🔧 Useful Commands"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Check health:        ~/check-health.sh"
echo "View logs:           cd ~/iumicert/packages/issuer && docker compose logs -f"
echo "Restart services:    cd ~/iumicert/packages/issuer && docker compose restart"
echo "Manual deploy:       cd ~/iumicert/packages/issuer && ./deploy.sh"
echo ""
