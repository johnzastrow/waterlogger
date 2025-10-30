@echo off
REM Waterlogger Windows Deployment Script
REM This script automates the deployment of Waterlogger as a Windows service
REM
REM ⚠️ IMPORTANT: Run this script as Administrator

setlocal enabledelayedexpansion

echo ==================================
echo Waterlogger Deployment Script
echo ==================================
echo.

REM Check if running as Administrator
net session >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ This script must be run as Administrator
    echo Right-click this file and select "Run as administrator"
    pause
    exit /b 1
)

REM Check if waterlogger.exe exists
if not exist "waterlogger.exe" (
    echo ❌ waterlogger.exe not found in current directory
    echo Please build the application first using: build.bat
    pause
    exit /b 1
)

REM Check if config.yaml exists
if not exist "config.yaml" (
    echo ❌ config.yaml not found in current directory
    pause
    exit /b 1
)

echo ✅ Found waterlogger.exe and config.yaml
echo.

REM Set installation directory
set "INSTALL_DIR=C:\Program Files\Waterlogger"

REM Create directory structure
echo 📁 Creating directory structure at %INSTALL_DIR%...
mkdir "%INSTALL_DIR%" 2>nul
mkdir "%INSTALL_DIR%\logs" 2>nul
mkdir "%INSTALL_DIR%\backups" 2>nul

REM Copy files
echo 📋 Copying files...
copy /Y waterlogger.exe "%INSTALL_DIR%\" >nul
copy /Y config.yaml "%INSTALL_DIR%\" >nul

echo ✅ Files copied successfully
echo.

REM Ask about production config
set /p PROD_CONFIG="Configure for production? (recommended: json logging, file output) [y/N]: "
if /i "%PROD_CONFIG%"=="y" (
    echo ⚙️  Updating config.yaml for production...
    powershell -Command "(gc '%INSTALL_DIR%\config.yaml') -replace 'format: console', 'format: json' | Set-Content '%INSTALL_DIR%\config.yaml'"
    powershell -Command "(gc '%INSTALL_DIR%\config.yaml') -replace 'output: both', 'output: file' | Set-Content '%INSTALL_DIR%\config.yaml'"
    powershell -Command "(gc '%INSTALL_DIR%\config.yaml') -replace 'level: debug', 'level: info' | Set-Content '%INSTALL_DIR%\config.yaml'"
    echo ✅ Config updated for production
)
echo.

REM Check if service already exists
sc query Waterlogger >nul 2>&1
if %errorlevel% equ 0 (
    echo ⚠️  Service 'Waterlogger' already exists
    set /p DELETE_SERVICE="Delete existing service and recreate? [y/N]: "
    if /i "!DELETE_SERVICE!"=="y" (
        echo Stopping service...
        sc stop Waterlogger >nul 2>&1
        timeout /t 2 /nobreak >nul
        echo Deleting service...
        sc delete Waterlogger >nul
        timeout /t 2 /nobreak >nul
        echo ✅ Existing service removed
    ) else (
        echo Keeping existing service. Exiting...
        pause
        exit /b 0
    )
)
echo.

REM Create Windows service
echo 📝 Creating Windows service...
sc create Waterlogger binpath="\"%INSTALL_DIR%\waterlogger.exe\" -config \"%INSTALL_DIR%\config.yaml\"" start=auto DisplayName="Waterlogger" >nul
if %errorlevel% neq 0 (
    echo ❌ Failed to create service
    pause
    exit /b 1
)

REM Set service description
sc description Waterlogger "Pool and Hot Tub Water Management System" >nul

REM Set service to restart on failure
sc failure Waterlogger reset=86400 actions=restart/60000/restart/60000/restart/60000 >nul

echo ✅ Service created successfully
echo.

echo ==================================
echo ✅ Installation Complete!
echo ==================================
echo.
echo Next steps:
echo.
echo 1. Start the service:
echo    sc start Waterlogger
echo.
echo 2. Check service status:
echo    sc query Waterlogger
echo.
echo 3. View logs:
echo    type "%INSTALL_DIR%\logs\waterlogger.log"
echo.
echo 4. Access the application:
echo    http://localhost:2342
echo.
echo 5. Default credentials (first run):
echo    Username: admin
echo    Password: admin
echo    ⚠️  IMPORTANT: Change this password immediately!
echo.
echo Management commands:
echo    Start:   sc start Waterlogger
echo    Stop:    sc stop Waterlogger
echo    Delete:  sc delete Waterlogger
echo.

REM Ask if user wants to start now
set /p START_NOW="Start the service now? [Y/n]: "
if /i not "%START_NOW%"=="n" (
    echo.
    echo Starting service...
    sc start Waterlogger
    if %errorlevel% equ 0 (
        echo ✅ Service started successfully!
        echo.
        echo You can now access Waterlogger at: http://localhost:2342
    ) else (
        echo ❌ Failed to start service
        echo Check logs at: %INSTALL_DIR%\logs\waterlogger.log
    )
)

echo.
pause
