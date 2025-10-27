# Build Requirements for Waterlogger

## Important: CGO Requirement

**Waterlogger uses SQLite which requires CGO (C bindings) to be enabled.**

If you see this error:
```
Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
```

Follow the instructions below for your platform.

---

## Windows Build Requirements

### 1. Install Go
- Download from: https://go.dev/dl/
- Version required: Go 1.21 or higher

### 2. Install MSYS2 with MinGW-w64 GCC (Required for SQLite)

**⚠️ IMPORTANT: DO NOT USE CYGWIN**
Cygwin's GCC cannot build native Windows Go programs with CGO. Use MSYS2 instead.

**Step 1: Download and Install MSYS2**
1. Download from: **https://www.msys2.org/**
2. Run the installer (e.g., `msys2-x86_64-latest.exe`)
3. Follow the installation wizard (default: `C:\msys64`)
4. When finished, it will open an MSYS2 terminal

**Step 2: Update MSYS2**
In the MSYS2 terminal:
```bash
pacman -Syu
```
If prompted to close the terminal, close it and reopen "MSYS2 MSYS" from Start Menu, then:
```bash
pacman -Su
```

**Step 3: Install MinGW-w64 GCC**
```bash
pacman -S mingw-w64-x86_64-gcc
```

**Step 4: Add to Windows PATH**
1. Press `Win + X` → "System"
2. Click "Advanced system settings"
3. Click "Environment Variables"
4. Under "System variables", select "Path" → "Edit"
5. Click "New" → Add: `C:\msys64\mingw64\bin`
6. Click "OK" on all dialogs
7. **Restart your terminal/command prompt**

**Step 5: Verify Installation**
Open a **new Windows Command Prompt** (NOT MSYS2 terminal):
```cmd
gcc --version
```

Should output something like:
```
gcc.exe (Rev10, Built by MSYS2 project) 13.2.0
```

### 3. Build Waterlogger

**IMPORTANT: Build from Windows Command Prompt, NOT MSYS2 terminal**

Open Windows Command Prompt or PowerShell:
```cmd
cd C:\Users\YourUsername\path\to\waterlogger
build.bat
```

The build script automatically:
- Enables CGO (required for SQLite)
- Checks for GCC availability
- Sets build timestamp
- Creates `waterlogger.exe`

---

## Windows Subsystem for Linux (WSL2) - Alternative Approach

If you're using Windows with WSL2 enabled, you can build the **Linux version** directly without dealing with Windows C compilers:

### Advantages of Using WSL2:
- No need to install MinGW or deal with Windows toolchains
- Native Linux build tools (already available in WSL2)
- Simpler setup and fewer compatibility issues

### Build Steps in WSL2:
```bash
# Open WSL2 terminal (Ubuntu, Debian, etc.)
cd /mnt/c/Users/YourUsername/path/to/waterlogger

# Install build tools (if not already installed)
sudo apt update
sudo apt install build-essential golang-go

# Build the Linux binary
./build.sh

# Run the application
./waterlogger
```

**Note:** The Linux binary will run perfectly in WSL2 and can be accessed from Windows browsers at `http://localhost:2342`

---

## Linux Build Requirements

### 1. Install Go
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install golang-go

# Fedora/RHEL
sudo dnf install golang

# Or download from: https://go.dev/dl/
```

### 2. Install Build Tools (Usually Already Installed)

**Ubuntu/Debian:**
```bash
sudo apt install build-essential
```

**Fedora/RHEL:**
```bash
sudo dnf groupinstall "Development Tools"
```

**Verify:**
```bash
gcc --version
```

### 3. Build Waterlogger

**Using build script:**
```bash
chmod +x build.sh
./build.sh
```

**Using make:**
```bash
make build
```

**Manual build:**
```bash
export CGO_ENABLED=1
go build -o waterlogger ./cmd/waterlogger
```

---

## macOS Build Requirements

### 1. Install Xcode Command Line Tools
```bash
xcode-select --install
```

This includes GCC/Clang which is required for CGO.

### 2. Install Go
```bash
# Using Homebrew
brew install go

# Or download from: https://go.dev/dl/
```

### 3. Build Waterlogger

**Using build script:**
```bash
chmod +x build.sh
./build.sh
```

**Using make:**
```bash
make build
```

---

## Common Issues and Solutions

### Issue 1: "CGO_ENABLED=0" Error

**Symptom:**
```
Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work
```

**Solution:**
Ensure CGO is enabled when building:
```bash
# Linux/macOS
export CGO_ENABLED=1
go build ./cmd/waterlogger

# Windows
set CGO_ENABLED=1
go build ./cmd/waterlogger
```

---

### Issue 2: "gcc: command not found" (Windows)

**Symptom:**
```
gcc: executable file not found in %PATH%
```

**Solution:**
1. Make sure MSYS2 is installed (see Windows section above)
2. Verify GCC is installed in MSYS2:
   ```bash
   # In MSYS2 terminal:
   pacman -S mingw-w64-x86_64-gcc
   ```
3. Add `C:\msys64\mingw64\bin` to Windows PATH
4. **Restart your terminal/command prompt**
5. Verify: `gcc --version` in Windows Command Prompt

---

### Issue 3: "cannot find -lgcc_s" (Windows)

**Symptom:**
```
cannot find -lgcc_s
```

**Solution:**
This means your MSYS2 MinGW-w64 installation is incomplete or not in PATH.
1. Verify GCC is fully installed:
   ```bash
   # In MSYS2 terminal:
   pacman -S mingw-w64-x86_64-gcc
   ```
2. Ensure `C:\msys64\mingw64\bin` is in your Windows PATH
3. Restart terminal/command prompt

---

### Issue 4: Cygwin Compiler Error (Windows)

**Symptom:**
```
# runtime/cgo
gcc_libinit_windows.c:6:2: error: #error "don't use the cygwin compiler to build native Windows programs; use MinGW instead"
```

**Root Cause:**
Cygwin's GCC creates POSIX-emulated binaries, not native Windows executables. Go's CGO runtime explicitly blocks Cygwin for Windows builds.

**Solution:**
You **must** use MSYS2 with MinGW-w64 instead of Cygwin:

1. **Install MSYS2** (see Windows section above)
2. Install GCC in MSYS2:
   ```bash
   pacman -S mingw-w64-x86_64-gcc
   ```
3. Add `C:\msys64\mingw64\bin` to Windows PATH
4. **Remove Cygwin from PATH**:
   - Edit your system PATH environment variable
   - Remove any paths containing "cygwin"
   - Restart terminal after removing

---

### Issue 5: Cross-Compilation Issues

**Important:** Cross-compiling with CGO is complex and requires platform-specific C compilers.

**For Linux → Windows:**
```bash
# Install mingw-w64 cross-compiler
sudo apt install gcc-mingw-w64

# Build
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc \
  go build -o waterlogger.exe ./cmd/waterlogger
```

**For Windows → Linux:**
Not recommended. Build on target platform instead.

---

## Verifying Your Build

After building, verify CGO is enabled:

**Windows:**
```cmd
.\waterlogger.exe -version
```

**Linux/macOS:**
```bash
./waterlogger -version
```

Should output:
```
Waterlogger v1.3.0
```

If you see the SQLite error, CGO was not enabled during build.

---

## Alternative: Using MariaDB Instead of SQLite

If you cannot get CGO working, you can use MariaDB instead:

1. Install MariaDB server
2. Edit `config.yaml`:
   ```yaml
   database:
     type: mariadb
     mariadb:
       host: localhost
       port: 3306
       username: waterlogger
       password: your-password
       database: waterlogger
   ```
3. Build without CGO requirement? **NO** - Still need CGO for SQLite driver even if using MariaDB.

**Best solution:** Fix CGO setup using instructions above.

---

## Build Scripts Overview

### build.sh (Linux/macOS)
- Enables CGO automatically
- Sets build timestamp
- Produces `./waterlogger`

### build.bat (Windows)
- Enables CGO automatically
- Checks for GCC
- Sets build timestamp
- Produces `.\waterlogger.exe`

### Makefile
- Cross-platform
- Multiple targets (build, clean, test, etc.)
- Enables CGO for all builds

---

## Quick Start (After Requirements Installed)

**Windows:**
```cmd
git clone https://github.com/johnzastrow/waterlogger.git
cd waterlogger
go mod download
build.bat
```

**Linux/macOS:**
```bash
git clone https://github.com/johnzastrow/waterlogger.git
cd waterlogger
go mod download
chmod +x build.sh
./build.sh
```

---

## Testing Your Build Environment

**Test if CGO is available:**
```bash
# Linux/macOS
go env CGO_ENABLED

# Windows
go env CGO_ENABLED
```

Should output: `1`

**Test if GCC is available:**
```bash
gcc --version
```

Should output version information.

**Test building a simple CGO program:**
```bash
# Create test file
cat > test.go << EOF
package main

/*
#include <stdio.h>
void hello() {
    printf("CGO Works!\n");
}
*/
import "C"

func main() {
    C.hello()
}
EOF

# Try to build
CGO_ENABLED=1 go build test.go

# Run it
./test  # Linux/macOS
test.exe  # Windows

# Cleanup
rm test.go test  # Linux/macOS
del test.go test.exe  # Windows
```

If this works, your CGO setup is correct and you can build Waterlogger.

---

## Getting Help

If you still have build issues:

1. Check Go version: `go version` (need 1.21+)
2. Check GCC version: `gcc --version`
3. Check CGO status: `go env CGO_ENABLED`
4. Try the test program above
5. Open an issue: https://github.com/johnzastrow/waterlogger/issues

Include:
- Operating system and version
- Go version
- GCC version
- Full error message
- Build command used

---

**Last Updated:** 2025-10-26
**Version:** 1.3.0
