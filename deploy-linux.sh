#!/bin/bash
# Waterlogger Linux Deployment Script
# This script automates the deployment of Waterlogger as a systemd service

set -e  # Exit on error

echo "=================================="
echo "Waterlogger Deployment Script"
echo "=================================="
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then
   echo "❌ This script must be run as root (use sudo)"
   exit 1
fi

# Check if waterlogger binary exists
if [ ! -f "waterlogger" ]; then
    echo "❌ waterlogger binary not found in current directory"
    echo "Please build the application first using: ./build.sh"
    exit 1
fi

# Check if config.yaml exists
if [ ! -f "config.yaml" ]; then
    echo "❌ config.yaml not found in current directory"
    exit 1
fi

# Check if web directory exists
if [ ! -d "web" ]; then
    echo "⚠️  web directory not found - templates and static files may be missing"
    echo "   (This is only required if the binary doesn't have embedded assets)"
fi

echo "✅ Found waterlogger binary and config.yaml"
echo ""

# Set installation directory
INSTALL_DIR="/opt/waterlogger"

# Create application directory structure
echo "📁 Creating directory structure at $INSTALL_DIR..."
mkdir -p "$INSTALL_DIR"
mkdir -p "$INSTALL_DIR/logs"
mkdir -p "$INSTALL_DIR/backups"

# Copy files
echo "📋 Copying files..."
cp waterlogger "$INSTALL_DIR/"
cp config.yaml "$INSTALL_DIR/"

# Copy web directory if it exists (contains templates and static assets)
if [ -d "web" ]; then
    echo "📋 Copying web assets (templates and static files)..."
    cp -r web "$INSTALL_DIR/"
    echo "✅ Web assets copied"
else
    echo "⚠️  Web directory not copied (binary must have embedded assets)"
fi

# Make executable
chmod +x "$INSTALL_DIR/waterlogger"

# Secure config file
chmod 600 "$INSTALL_DIR/config.yaml"

echo "✅ Files copied successfully"
echo ""

# Ask about production config
read -p "Configure for production? (recommended: json logging, file output) [y/N]: " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "⚙️  Updating config.yaml for production..."
    sed -i 's/format: console/format: json/g' "$INSTALL_DIR/config.yaml"
    sed -i 's/output: both/output: file/g' "$INSTALL_DIR/config.yaml"
    sed -i 's/level: debug/level: info/g' "$INSTALL_DIR/config.yaml"
    echo "✅ Config updated for production"
fi
echo ""

# Create dedicated user
echo "👤 Creating dedicated waterlogger user..."
if id "waterlogger" &>/dev/null; then
    echo "   User 'waterlogger' already exists, skipping..."
else
    useradd -r -s /bin/false waterlogger
    echo "✅ User 'waterlogger' created"
fi

# Set ownership
chown -R waterlogger:waterlogger "$INSTALL_DIR"
echo "✅ Ownership set to waterlogger:waterlogger"
echo ""

# Create systemd service file
echo "📝 Creating systemd service file..."
cat > /etc/systemd/system/waterlogger.service <<'EOF'
[Unit]
Description=Waterlogger - Pool and Hot Tub Water Management System
After=network.target
Documentation=https://github.com/johnzastrow/waterlogger

[Service]
Type=simple
User=waterlogger
Group=waterlogger
WorkingDirectory=/opt/waterlogger

# Environment
Environment=GIN_MODE=release

# Execution
ExecStart=/opt/waterlogger/waterlogger -config /opt/waterlogger/config.yaml

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=waterlogger

# Restart policy
Restart=always
RestartSec=10

# Security hardening
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/waterlogger

[Install]
WantedBy=multi-user.target
EOF

echo "✅ Service file created at /etc/systemd/system/waterlogger.service"
echo ""

# Reload systemd
echo "🔄 Reloading systemd daemon..."
systemctl daemon-reload

# Enable service
echo "⚙️  Enabling service to start on boot..."
systemctl enable waterlogger

echo ""
echo "=================================="
echo "✅ Installation Complete!"
echo "=================================="
echo ""
echo "Next steps:"
echo ""
echo "1. Start the service:"
echo "   sudo systemctl start waterlogger"
echo ""
echo "2. Check service status:"
echo "   sudo systemctl status waterlogger"
echo ""
echo "3. View logs:"
echo "   sudo journalctl -u waterlogger -f"
echo ""
echo "4. Access the application:"
echo "   http://your-server-ip:2342"
echo ""
echo "5. Default credentials (first run):"
echo "   Username: admin"
echo "   Password: admin"
echo "   ⚠️  IMPORTANT: Change this password immediately!"
echo ""
echo "To start now, run: sudo systemctl start waterlogger"
echo ""
