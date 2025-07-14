# Waterlogger Deployment Guide

This guide covers various deployment options for Waterlogger.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Standalone Deployment](#standalone-deployment)
- [Docker Deployment](#docker-deployment)
- [Service Deployment](#service-deployment)
- [Reverse Proxy Setup](#reverse-proxy-setup)
- [Database Setup](#database-setup)
- [SSL/TLS Configuration](#ssltls-configuration)
- [Monitoring and Logging](#monitoring-and-logging)
- [Backup and Recovery](#backup-and-recovery)
- [Troubleshooting](#troubleshooting)

## Prerequisites

- Linux server (Ubuntu 20.04+ recommended) or Windows Server
- Minimum 512MB RAM, 1GB recommended
- 10GB disk space (more for large datasets)
- Network access on desired port (default: 2341)

## Standalone Deployment

### Linux

1. **Download and extract the binary:**
```bash
wget https://github.com/your-org/waterlogger/releases/latest/download/waterlogger-linux-amd64.tar.gz
tar -xzf waterlogger-linux-amd64.tar.gz
cd waterlogger-linux-amd64
```

2. **Create application directory:**
```bash
sudo mkdir -p /opt/waterlogger
sudo cp waterlogger /opt/waterlogger/
sudo cp config.example.yaml /opt/waterlogger/config.yaml
sudo chmod +x /opt/waterlogger/waterlogger
```

3. **Create user and set permissions:**
```bash
sudo useradd -r -s /bin/false waterlogger
sudo chown -R waterlogger:waterlogger /opt/waterlogger
```

4. **Configure the application:**
```bash
sudo nano /opt/waterlogger/config.yaml
```

5. **Run the application:**
```bash
sudo -u waterlogger /opt/waterlogger/waterlogger -config /opt/waterlogger/config.yaml
```

### Windows

1. **Download and extract the binary:**
   - Download `waterlogger-windows-amd64.zip` from releases
   - Extract to `C:\Program Files\Waterlogger\`

2. **Create configuration:**
   - Copy `config.example.yaml` to `config.yaml`
   - Edit configuration as needed

3. **Run the application:**
```cmd
cd "C:\Program Files\Waterlogger"
waterlogger.exe -config config.yaml
```

## Docker Deployment

### Using Docker Compose (Recommended)

1. **Create docker-compose.yml:**
```yaml
version: '3.8'

services:
  waterlogger:
    image: waterlogger:latest
    ports:
      - "2341:2341"
    environment:
      - GIN_MODE=release
    volumes:
      - ./data:/app/data
      - ./config.yaml:/app/config.yaml
    depends_on:
      - db
    restart: unless-stopped

  db:
    image: mariadb:10.9
    environment:
      - MYSQL_ROOT_PASSWORD=rootpassword
      - MYSQL_DATABASE=waterlogger
      - MYSQL_USER=waterlogger
      - MYSQL_PASSWORD=waterlogger
    volumes:
      - db_data:/var/lib/mysql
    restart: unless-stopped

volumes:
  db_data:
```

2. **Deploy:**
```bash
docker-compose up -d
```

### Using Docker

1. **Build the image:**
```bash
docker build -t waterlogger:latest .
```

2. **Run the container:**
```bash
docker run -d \
  --name waterlogger \
  -p 2341:2341 \
  -v /path/to/data:/app/data \
  -v /path/to/config.yaml:/app/config.yaml \
  waterlogger:latest
```

## Service Deployment

### Linux (systemd)

1. **Create service file:**
```bash
sudo tee /etc/systemd/system/waterlogger.service > /dev/null <<EOF
[Unit]
Description=Waterlogger - Pool and Hot Tub Water Management System
After=network.target
Wants=network.target

[Service]
Type=simple
User=waterlogger
Group=waterlogger
WorkingDirectory=/opt/waterlogger
ExecStart=/opt/waterlogger/waterlogger -config /opt/waterlogger/config.yaml
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=waterlogger

# Security settings
NoNewPrivileges=yes
PrivateTmp=yes
ProtectSystem=strict
ReadWritePaths=/opt/waterlogger

[Install]
WantedBy=multi-user.target
EOF
```

2. **Enable and start service:**
```bash
sudo systemctl daemon-reload
sudo systemctl enable waterlogger
sudo systemctl start waterlogger
```

3. **Check status:**
```bash
sudo systemctl status waterlogger
sudo journalctl -u waterlogger -f
```

### Windows (Windows Service)

1. **Install as service using NSSM:**
```cmd
# Download NSSM from https://nssm.cc/
nssm install Waterlogger "C:\Program Files\Waterlogger\waterlogger.exe"
nssm set Waterlogger AppParameters "-config C:\Program Files\Waterlogger\config.yaml"
nssm set Waterlogger AppDirectory "C:\Program Files\Waterlogger"
nssm set Waterlogger Description "Pool and Hot Tub Water Management System"
nssm set Waterlogger Start SERVICE_AUTO_START
nssm start Waterlogger
```

2. **Using PowerShell (alternative):**
```powershell
New-Service -Name "Waterlogger" -BinaryPathName "C:\Program Files\Waterlogger\waterlogger.exe -config C:\Program Files\Waterlogger\config.yaml" -DisplayName "Waterlogger" -StartupType Automatic
Start-Service -Name "Waterlogger"
```

## Reverse Proxy Setup

### Nginx

1. **Install Nginx:**
```bash
sudo apt update
sudo apt install nginx
```

2. **Create configuration:**
```bash
sudo tee /etc/nginx/sites-available/waterlogger > /dev/null <<EOF
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        proxy_pass http://localhost:2341;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        
        # WebSocket support (if needed)
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
EOF
```

3. **Enable and restart:**
```bash
sudo ln -s /etc/nginx/sites-available/waterlogger /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### Apache

1. **Install Apache:**
```bash
sudo apt update
sudo apt install apache2
sudo a2enmod proxy proxy_http
```

2. **Create configuration:**
```bash
sudo tee /etc/apache2/sites-available/waterlogger.conf > /dev/null <<EOF
<VirtualHost *:80>
    ServerName your-domain.com
    
    ProxyPreserveHost On
    ProxyRequests Off
    ProxyPass / http://localhost:2341/
    ProxyPassReverse / http://localhost:2341/
    
    ProxyPassReverse / http://localhost:2341/
    ProxyPassReverseAdjustAddress On
</VirtualHost>
EOF
```

3. **Enable and restart:**
```bash
sudo a2ensite waterlogger
sudo systemctl reload apache2
```

## Database Setup

### SQLite (Default)

SQLite requires no additional setup. The database file is created automatically.

**Configuration:**
```yaml
database:
  type: "sqlite"
  sqlite:
    path: "/opt/waterlogger/waterlogger.db"
```

### MariaDB

1. **Install MariaDB:**
```bash
sudo apt update
sudo apt install mariadb-server
sudo mysql_secure_installation
```

2. **Create database and user:**
```sql
CREATE DATABASE waterlogger;
CREATE USER 'waterlogger'@'localhost' IDENTIFIED BY 'secure_password';
GRANT ALL PRIVILEGES ON waterlogger.* TO 'waterlogger'@'localhost';
FLUSH PRIVILEGES;
```

3. **Configure application:**
```yaml
database:
  type: "mariadb"
  mariadb:
    host: "localhost"
    port: 3306
    username: "waterlogger"
    password: "secure_password"
    database: "waterlogger"
```

### PostgreSQL (Future Support)

PostgreSQL support is planned for future releases.

## SSL/TLS Configuration

### Using Let's Encrypt with Nginx

1. **Install Certbot:**
```bash
sudo apt install certbot python3-certbot-nginx
```

2. **Obtain certificate:**
```bash
sudo certbot --nginx -d your-domain.com
```

3. **Auto-renewal:**
```bash
sudo crontab -e
# Add this line:
0 12 * * * /usr/bin/certbot renew --quiet
```

### Using Custom SSL Certificate

1. **Update Nginx configuration:**
```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;
    
    ssl_certificate /path/to/certificate.crt;
    ssl_certificate_key /path/to/private.key;
    
    # SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    
    location / {
        proxy_pass http://localhost:2341;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# Redirect HTTP to HTTPS
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}
```

## Monitoring and Logging

### Systemd Logging

```bash
# View logs
sudo journalctl -u waterlogger -f

# View logs with timestamps
sudo journalctl -u waterlogger --since "2024-07-14 10:00:00"

# Export logs
sudo journalctl -u waterlogger --since "2024-07-14" > waterlogger.log
```

### Log Rotation

```bash
sudo tee /etc/logrotate.d/waterlogger > /dev/null <<EOF
/var/log/waterlogger/*.log {
    daily
    missingok
    rotate 14
    compress
    delaycompress
    notifempty
    create 0644 waterlogger waterlogger
    postrotate
        systemctl reload waterlogger
    endscript
}
EOF
```

### Monitoring with Prometheus (Optional)

1. **Add metrics endpoint to application**
2. **Configure Prometheus to scrape metrics**
3. **Set up Grafana dashboards**

## Backup and Recovery

### SQLite Backup

```bash
#!/bin/bash
# backup-sqlite.sh

BACKUP_DIR="/opt/waterlogger/backups"
DB_FILE="/opt/waterlogger/waterlogger.db"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p "$BACKUP_DIR"
sqlite3 "$DB_FILE" ".backup $BACKUP_DIR/waterlogger_$DATE.db"
gzip "$BACKUP_DIR/waterlogger_$DATE.db"

# Keep only last 30 days
find "$BACKUP_DIR" -name "waterlogger_*.db.gz" -mtime +30 -delete
```

### MariaDB Backup

```bash
#!/bin/bash
# backup-mariadb.sh

BACKUP_DIR="/opt/waterlogger/backups"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p "$BACKUP_DIR"
mysqldump -u waterlogger -p waterlogger > "$BACKUP_DIR/waterlogger_$DATE.sql"
gzip "$BACKUP_DIR/waterlogger_$DATE.sql"

# Keep only last 30 days
find "$BACKUP_DIR" -name "waterlogger_*.sql.gz" -mtime +30 -delete
```

### Automated Backups

```bash
# Add to crontab
crontab -e

# Daily backup at 2 AM
0 2 * * * /opt/waterlogger/backup-sqlite.sh
```

## Troubleshooting

### Common Issues

#### Port Already in Use

```bash
# Check what's using the port
sudo netstat -tlnp | grep :2341
sudo ss -tlnp | grep :2341

# Change port in config.yaml
server:
  port: 3000
```

#### Database Connection Issues

```bash
# Check database status
sudo systemctl status mariadb

# Test connection
mysql -u waterlogger -p waterlogger

# Check logs
sudo journalctl -u waterlogger -f
```

#### Permission Issues

```bash
# Fix file permissions
sudo chown -R waterlogger:waterlogger /opt/waterlogger
sudo chmod +x /opt/waterlogger/waterlogger
```

#### Service Won't Start

```bash
# Check service status
sudo systemctl status waterlogger

# Check logs
sudo journalctl -u waterlogger -f

# Test manually
sudo -u waterlogger /opt/waterlogger/waterlogger -config /opt/waterlogger/config.yaml
```

### Performance Tuning

#### Database Optimization

**SQLite:**
```yaml
database:
  sqlite:
    path: "/opt/waterlogger/waterlogger.db"
    # Add pragma settings for performance
    options: "?cache=shared&mode=rwc&_journal_mode=WAL"
```

**MariaDB:**
```sql
-- Optimize for read-heavy workloads
SET GLOBAL innodb_buffer_pool_size = 256M;
SET GLOBAL innodb_log_file_size = 64M;
SET GLOBAL max_connections = 100;
```

#### Application Tuning

```yaml
server:
  port: 2341
  host: "0.0.0.0"
  # Add performance settings
  read_timeout: 30
  write_timeout: 30
  max_header_bytes: 1048576
```

### Security Considerations

1. **Change default ports**
2. **Use strong passwords**
3. **Enable SSL/TLS**
4. **Configure firewall**
5. **Regular security updates**
6. **Monitor access logs**
7. **Use non-root user**
8. **Restrict file permissions**

### Health Checks

```bash
#!/bin/bash
# health-check.sh

URL="http://localhost:2341"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$URL")

if [ "$STATUS" -eq 200 ]; then
    echo "Waterlogger is healthy"
    exit 0
else
    echo "Waterlogger is unhealthy (HTTP $STATUS)"
    exit 1
fi
```

## Production Checklist

- [ ] Application deployed with proper user permissions
- [ ] Database configured and optimized
- [ ] SSL/TLS certificate installed
- [ ] Reverse proxy configured
- [ ] Firewall configured
- [ ] Monitoring and logging set up
- [ ] Backup system configured
- [ ] Health checks implemented
- [ ] Documentation updated
- [ ] Testing completed

## Support

For deployment issues:
- Check the [Troubleshooting](#troubleshooting) section
- Review application logs
- Consult the [GitHub Issues](https://github.com/your-org/waterlogger/issues)
- Join the community discussions