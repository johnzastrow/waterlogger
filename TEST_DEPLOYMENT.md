# Waterlogger Deployment Test Plan

## Current Status
- ✅ Binary built successfully: `/home/jcz/Github/waterlogger/waterlogger` (22 MB)
- ✅ Updated `deploy-linux.sh` to copy web directory
- ✅ Modified `cmd/waterlogger/main.go` to load assets from multiple paths

## Pre-Deployment Checklist

### 1. Verify Binary Works Locally (Dev Mode)
```bash
cd /home/jcz/Github/waterlogger
./waterlogger -version
```
Expected output: `Waterlogger v1.5.0`

### 2. Test Binary with Dev Config (Optional)
```bash
cd /home/jcz/Github/waterlogger
# Create a test database in a temp location
mkdir -p /tmp/waterlogger-test
cd /tmp/waterlogger-test
cp /home/jcz/Github/waterlogger/config.example.yaml ./config.yaml
# Run for a few seconds then stop with Ctrl+C
/home/jcz/Github/waterlogger/waterlogger -config ./config.yaml
```
Expected: Should start successfully and load templates from source directory

## Deployment Steps (Run as root/sudo)

### Option A: Use Updated Deploy Script (Recommended)
```bash
cd /home/jcz/Github/waterlogger
sudo ./deploy-linux.sh
```
This will:
- Copy binary to `/opt/waterlogger/waterlogger`
- Copy config to `/opt/waterlogger/config.yaml`
- **Copy web directory to `/opt/waterlogger/web`** ← KEY FIX
- Create systemd service
- Enable service
- Create dedicated waterlogger user

### Option B: Manual Deployment (Step by Step)
```bash
# 1. Stop current service (if running)
sudo systemctl stop waterlogger || true
sudo systemctl disable waterlogger || true

# 2. Create directory structure
sudo mkdir -p /opt/waterlogger/logs /opt/waterlogger/backups

# 3. Copy files
sudo cp /home/jcz/Github/waterlogger/waterlogger /opt/waterlogger/
sudo cp /home/jcz/Github/waterlogger/config.yaml /opt/waterlogger/
sudo cp -r /home/jcz/Github/waterlogger/web /opt/waterlogger/

# 4. Set permissions
sudo chmod +x /opt/waterlogger/waterlogger
sudo chmod 600 /opt/waterlogger/config.yaml
sudo chown -R waterlogger:waterlogger /opt/waterlogger

# 5. Update config for production (if desired)
sudo sed -i 's/format: console/format: json/g' /opt/waterlogger/config.yaml
sudo sed -i 's/output: both/output: file/g' /opt/waterlogger/config.yaml
sudo sed -i 's/level: debug/level: info/g' /opt/waterlogger/config.yaml

# 6. Create service file
sudo tee /etc/systemd/system/waterlogger.service > /dev/null <<'EOF'
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

# 7. Enable service
sudo systemctl daemon-reload
sudo systemctl enable waterlogger
```

## Post-Deployment Testing

### 1. Start the Service
```bash
sudo systemctl start waterlogger
```

### 2. Check Service Status
```bash
sudo systemctl status waterlogger
```
Should show:
```
● waterlogger.service - Waterlogger - Pool and Hot Tub Water Management System
     Loaded: loaded (/etc/systemd/system/waterlogger.service; enabled; preset: enabled)
     Active: active (running) since ...
```

### 3. Check Recent Logs
```bash
sudo journalctl -u waterlogger -n 50
```
Look for these SUCCESS indicators:
```
Waterlogger starting...
Build information
Connecting to SQLite database
Successfully connected to SQLite database
Running database migrations
Database migrations completed successfully
Loaded HTML templates ← KEY: This means templates were found!
Serving static files    ← KEY: This means static files were found!
Starting Waterlogger server
```

### 4. Follow Live Logs
```bash
sudo journalctl -u waterlogger -f
```
This will show new log entries in real-time. Leave it running and proceed to next step.

### 5. Test Web Interface (in another terminal)
```bash
curl -I http://localhost:2342/
```
Expected response:
```
HTTP/1.1 302 Found
Location: /setup
```

Or open your browser:
```
http://your-server-ip:2342
```
You should see:
- ✅ Setup wizard loads (or login page if already set up)
- ✅ CSS styling appears (proves static files loaded)
- ✅ No 404 errors for templates

### 6. Verify File Locations
```bash
# Check deployment has all files
ls -la /opt/waterlogger/
ls -la /opt/waterlogger/web/
ls -la /opt/waterlogger/web/templates/
ls -la /opt/waterlogger/web/static/
```

Expected output:
```
/opt/waterlogger/:
total 23M
-rwxr-xr-x 1 waterlogger waterlogger  22M waterlogger
-rw-r--r-- 1 waterlogger waterlogger 5.7K config.yaml
drwxr-xr-x 2 waterlogger waterlogger 4.0K logs
drwxr-xr-x 2 waterlogger waterlogger 4.0K backups
drwxr-xr-x 3 waterlogger waterlogger 4.0K web  ← KEY: This should exist now!

/opt/waterlogger/web/templates/:
-rw-r--r-- 1 waterlogger waterlogger 404.html
-rw-r--r-- 1 waterlogger waterlogger adjustments.html
... (11 HTML files total)

/opt/waterlogger/web/static/:
-rw-r--r-- 1 waterlogger waterlogger favicon.ico
drwxr-xr-x 2 waterlogger waterlogger css
drwxr-xr-x 2 waterlogger waterlogger js
```

## Troubleshooting

### Issue: "pattern matches no files: web/templates/*.html"
**Solution:**
```bash
# Verify web directory was copied
sudo ls -la /opt/waterlogger/web/templates/

# If missing, copy it manually:
sudo cp -r /home/jcz/Github/waterlogger/web /opt/waterlogger/
sudo chown -R waterlogger:waterlogger /opt/waterlogger/web

# Restart service
sudo systemctl restart waterlogger
```

### Issue: Service crashes immediately
```bash
# Check detailed error logs
sudo journalctl -u waterlogger -n 100

# Check if web files are readable by waterlogger user
sudo sudo -u waterlogger ls /opt/waterlogger/web/templates/

# Check config file is valid
sudo /opt/waterlogger/waterlogger -config /opt/waterlogger/config.yaml
```

### Issue: CSS/JS not loading (404 errors)
```bash
# Check static files exist and are readable
sudo ls -la /opt/waterlogger/web/static/
sudo ls -la /opt/waterlogger/web/static/css/
sudo ls -la /opt/waterlogger/web/static/js/

# Check permissions
sudo sudo -u waterlogger cat /opt/waterlogger/web/static/css/style.css | wc -l
```

## Success Criteria

✅ All of the following should be true:
1. Service status shows `active (running)`
2. Logs show "Loaded HTML templates" and "Serving static files"
3. Curl to `http://localhost:2342/` returns HTTP 302 (redirect, not 500)
4. Opening browser to `http://your-server-ip:2342` shows properly styled page
5. Setup wizard or login page appears with full styling
6. Web directory exists at `/opt/waterlogger/web/`
7. No "pattern matches no files" errors in logs

## Next Steps After Success

1. Complete setup wizard if first run
2. Create test admin account
3. Test basic functionality (add pool, add sample, etc.)
4. Configure production logging if desired
5. Monitor logs for any warnings

## Questions?

If deployment fails, capture:
1. Output from `sudo systemctl status waterlogger`
2. Last 50 lines of `sudo journalctl -u waterlogger`
3. Output of `ls -la /opt/waterlogger/`
4. Output of `ls -la /opt/waterlogger/web/`

This will help diagnose the issue quickly.
