# IU-MiCert Issuer - GCP Deployment Guide

This guide walks you through deploying the IU-MiCert issuer backend and PostgreSQL database on a Google Cloud Platform VM using Docker.

## Deployment Options

1. **[Automated CI/CD with GitHub Actions](./CI-CD-SETUP.md)** (Recommended) - Automatic deployments on git push
2. **Manual Docker Deployment** (This guide) - One-time setup or manual control

## Prerequisites

- A GCP VM instance (Ubuntu 20.04 or later recommended, **AMD64 architecture**)
- Docker and Docker Compose installed on the VM
- Ethereum Sepolia testnet private key
- Infura API key (or other Ethereum RPC provider)

## Step 1: Prepare Your GCP VM

### Create a VM Instance

1. Go to GCP Console → Compute Engine → VM Instances
2. Create a new instance with these specs:
   - **Machine type**: e2-medium (2 vCPU, 4 GB memory) or higher
   - **Boot disk**: Ubuntu 22.04 LTS, 20 GB
   - **Firewall**: Allow HTTP (80) and HTTPS (443) traffic
3. Note your VM's external IP address

### SSH into Your VM

```bash
gcloud compute ssh YOUR_VM_NAME --zone=YOUR_ZONE
```

## Step 2: Install Docker and Docker Compose

```bash
# Update package list
sudo apt-get update

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Add your user to docker group (to run without sudo)
sudo usermod -aG docker $USER

# Install Docker Compose
sudo apt-get install -y docker-compose-plugin

# Verify installation
docker --version
docker compose version

# Logout and login again for group changes to take effect
exit
```

SSH back into your VM after logging out.

## Step 3: Clone and Setup the Project

```bash
# Clone the repository
git clone <your-repo-url>
cd iumicert/packages/issuer

# Or if uploading files manually:
# mkdir -p ~/iumicert/packages/issuer
# Upload your files using scp or gcloud compute scp
```

## Step 4: Configure Environment Variables

```bash
# Copy the production environment template
cp .env.production .env

# Edit the environment file
nano .env
```

Configure the following required values:

```bash
# Strong password for PostgreSQL
POSTGRES_PASSWORD=your_strong_password_here

# Your Sepolia private key (WITHOUT 0x prefix)
ISSUER_PRIVATE_KEY=your_private_key_here

# Your Infura API key
SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/your_infura_api_key

# Optional: Change default contract if needed
IUMICERT_CONTRACT_ADDRESS=0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60
```

Save and exit (Ctrl+X, then Y, then Enter).

## Step 5: Configure Firewall Rules

Allow traffic to your API port (8080):

```bash
# Using GCP Console
gcloud compute firewall-rules create allow-issuer-api \
    --allow tcp:8080 \
    --source-ranges 0.0.0.0/0 \
    --target-tags http-server

# Or using gcloud CLI
gcloud compute instances add-tags YOUR_VM_NAME \
    --tags http-server \
    --zone YOUR_ZONE
```

## Step 6: Deploy the Application

```bash
# Run the deployment script
./deploy.sh
```

The script will:
1. Validate your configuration
2. Build Docker images
3. Start PostgreSQL and the issuer backend
4. Check service health
5. Display logs

## Step 7: Verify Deployment

### Check Running Containers

```bash
docker compose ps
```

You should see:
- `iumicert-postgres` - healthy
- `iumicert-issuer` - running

### Test the API

```bash
# Health check
curl http://localhost:8080/api/health

# List terms
curl http://localhost:8080/api/terms

# From external network (use your VM's external IP)
curl http://YOUR_VM_EXTERNAL_IP:8080/api/health
```

### View Logs

```bash
# All services
docker compose logs -f

# Just the issuer backend
docker compose logs -f issuer-backend

# Just PostgreSQL
docker compose logs -f postgres
```

## Management Commands

### Start Services

```bash
docker compose up -d
```

### Stop Services

```bash
docker compose down
```

### Restart Services

```bash
docker compose restart
```

### Update and Redeploy

```bash
# Pull latest changes
git pull

# Redeploy
./deploy.sh
```

### Access PostgreSQL Database

```bash
# Using docker exec
docker compose exec postgres psql -U iumicert -d iumicert

# Or enable Adminer (optional - lightweight web UI)
docker compose --profile admin up -d
# Access at http://YOUR_VM_IP:8081
# Server: postgres, Username: iumicert, Password: <your password>, Database: iumicert
```

## Data Management

### Backup Data

```bash
# Backup PostgreSQL
docker compose exec postgres pg_dump -U iumicert iumicert > backup.sql

# Backup application data
tar -czf data_backup.tar.gz data/ publish_ready/
```

### Restore Data

```bash
# Restore PostgreSQL
docker compose exec -T postgres psql -U iumicert iumicert < backup.sql

# Restore application data
tar -xzf data_backup.tar.gz
```

### Generate Test Data

```bash
# Enter the container
docker compose exec issuer-backend sh

# Run generation commands
./micert generate-data
./micert batch-process
./micert generate-all-receipts

# Exit container
exit
```

## Troubleshooting

### Container Won't Start

```bash
# Check logs
docker compose logs issuer-backend

# Check if port is already in use
sudo netstat -tlnp | grep 8080

# Rebuild from scratch
docker compose down -v
docker compose up -d --build
```

### Database Connection Issues

```bash
# Check PostgreSQL is running
docker compose ps postgres

# Check database logs
docker compose logs postgres

# Test database connection
docker compose exec postgres psql -U iumicert -d iumicert -c "SELECT 1;"
```

### Out of Memory

```bash
# Check memory usage
free -h
docker stats

# Consider upgrading to a larger VM instance
```

### Can't Access API from Outside

```bash
# Check firewall rules
gcloud compute firewall-rules list

# Check if service is listening
docker compose exec issuer-backend netstat -tlnp

# Verify external IP
curl ifconfig.me
```

## Security Best Practices

1. **Never commit `.env` file** - Keep private keys secure
2. **Use strong passwords** - For PostgreSQL and pgAdmin
3. **Restrict firewall rules** - Only allow necessary ports
4. **Regular updates** - Keep Docker images and system updated
5. **Use HTTPS** - Set up a reverse proxy with SSL (Nginx/Caddy)
6. **Monitor logs** - Watch for suspicious activity
7. **Backup regularly** - Automate database and data backups

## Production Hardening (Optional)

### Set up Nginx Reverse Proxy with SSL

```bash
sudo apt-get install -y nginx certbot python3-certbot-nginx

# Configure Nginx
sudo nano /etc/nginx/sites-available/iumicert

# Add:
server {
    listen 80;
    server_name yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # Optional: Adminer admin interface
    location /admin {
        proxy_pass http://localhost:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}

# Enable site
sudo ln -s /etc/nginx/sites-available/iumicert /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx

# Get SSL certificate
sudo certbot --nginx -d yourdomain.com
```

### Set up Log Rotation

```bash
sudo nano /etc/logrotate.d/iumicert
```

Add:
```
/home/youruser/iumicert/packages/issuer/logs/*.log {
    daily
    rotate 14
    compress
    delaycompress
    notifempty
    create 0640 youruser youruser
}
```

## Monitoring

### Set up Health Check Monitoring

Use GCP monitoring or external services like UptimeRobot to monitor:
- `http://YOUR_VM_IP:8080/api/health`

### View Resource Usage

```bash
# Real-time container stats
docker stats

# VM resource usage
htop  # Install with: sudo apt-get install htop
```

## Cost Optimization

- Use **Preemptible/Spot VMs** for development (much cheaper)
- **Stop VM** when not in use: `gcloud compute instances stop YOUR_VM_NAME`
- Use **committed use discounts** for production
- Monitor with **GCP Cost Management**

## Support

For issues or questions:
- Check logs: `docker compose logs -f`
- Review this documentation
- Check project README: `../../README.md`
- Check project issues on GitHub

---

