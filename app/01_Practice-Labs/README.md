# Practice Labs

This repository contains distributed systems practice labs that require **Go version 1.25.1** or later.

## What’s new in this folder
- Lab guides refined with clearer steps and hints:
  - `lab4a.md`: Shardmaster with step-by-step checklist and pitfalls
  - `lab4b.md`: Sharded KV with reconfiguration protocol and dedup transfer details
  - `lab5.md`: Persistence with atomic file writes, recovery plan, and disk layout
- Source updates for Go 1.25.1 compatibility in `src/shardkv` and `src/diskv` with improved inline documentation.

## Quick links
- Lab 1: `lab1.md`
- Lab 2A: `lab2a.md`, Lab 2B: `lab2b.md`
- Lab 3A: `lab3a.md`, Lab 3B: `lab3b.md`
- Lab 4A (Shardmaster): `lab4a.md`
- Lab 4B (Sharded KV): `lab4b.md`
- Lab 5 (Persistence): `lab5.md`

## Run tests and examples
- Shardmaster (Lab 4A):
  ```bash
  cd app/01_Practice-Labs/src/shardmaster
  go test
  ```
- Sharded KV (Lab 4B):
  ```bash
  cd app/01_Practice-Labs/src/shardkv
  go test
  ```
- Persistent KV (Lab 5):
  ```bash
  cd app/01_Practice-Labs/src/diskv
  # run lab4-compatible subset
  go test -run Test4
  # run full lab5 suite
  go test
  ```
- MapReduce examples:
  ```bash
  cd app/01_Practice-Labs/src/main
  go run wc.go master kjv12.txt sequential
  ```

## Go Setup Instructions

### Prerequisites
- **Required Go Version**: Go 1.25.1 or later
- **Minimum Requirements**: 
  - macOS 11 Big Sur or later
  - Linux kernel 3.2 or later
  - Windows 10 or later

### Installation by Platform

#### macOS

**Option 1: Official Installer (Recommended)**
1. Download the macOS installer from [go.dev/dl](https://go.dev/dl/)
2. Select the appropriate architecture (Intel or Apple Silicon)
3. Run the `.pkg` installer
4. Verify installation:
   ```bash
   go version
   ```

**Option 2: Homebrew**
```bash
brew install go@1.25
```

**Option 3: Version Manager (g)**
```bash
# Install g
curl -sSL https://git.io/g-install | sh -s
# Install Go 1.25.1
g install 1.25.1
```

#### Linux

**Option 1: Official Binary (Recommended)**
```bash
# Download and extract
wget https://go.dev/dl/go1.25.1.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.1.linux-amd64.tar.gz

# Add to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH=$PATH:/usr/local/go/bin
```

**Option 2: Package Manager**
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install golang

# CentOS/RHEL/Fedora
sudo dnf install golang
# or
sudo yum install golang
```

**Option 3: Version Manager (gvm)**
```bash
# Install gvm
bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
# Install Go 1.25.1
gvm install go1.25.1
gvm use go1.25.1 --default
```

#### Windows

**Option 1: Official Installer (Recommended)**
1. Download the Windows installer from [go.dev/dl](https://go.dev/dl/)
2. Run the `.msi` installer
3. Follow the installation wizard
4. Verify installation:
   ```cmd
   go version
   ```

**Option 2: Chocolatey**
```cmd
choco install golang
```

**Option 3: Scoop**
```cmd
scoop install go
```

### Verification

After installation, verify your Go setup:

```bash
go version
# Should output: go version go1.25.1 linux/amd64 (or your platform)
```

### Environment Setup

**For Modern Go Development (Recommended):**
```bash
# Initialize Go modules in your project
cd /path/to/your/project
go mod init your-project-name

# Verify Go environment
go env GOPATH
go env GOROOT
```

**For Legacy GOPATH Setup (if needed):**
```bash
# Set GOPATH (not recommended for new projects)
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

### Lab-Specific Setup

These labs were originally designed for older Go releases but have been updated to work with **Go 1.25.1**. The codebase uses:

- Standard Go libraries (no external dependencies)
- RPC communication
- File I/O operations
- Concurrent programming with goroutines

### Troubleshooting

**Common Issues:**

1. **"go: command not found"**
   - Ensure Go is in your PATH
   - Restart your terminal after installation

2. **Version mismatch**
   - Check with `go version`
   - Update to Go 1.25.1 if using an older version

3. **Permission issues (Linux/macOS)**
   - Use `sudo` for system-wide installation
   - Or install to user directory and update PATH

4. **Module vs GOPATH conflicts**
   - Use `go mod init` for modern development
   - Clear GOPATH if using modules
