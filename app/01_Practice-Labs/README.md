# Practice Labs

This repository contains distributed systems practice labs that require **Go version 1.23.x** or later.

## Go Setup Instructions

### Prerequisites
- **Required Go Version**: Go 1.23.x or later
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
brew install go@1.23
```

**Option 3: Version Manager (g)**
```bash
# Install g
curl -sSL https://git.io/g-install | sh -s
# Install Go 1.23
g install 1.23.4
```

#### Linux

**Option 1: Official Binary (Recommended)**
```bash
# Download and extract
wget https://go.dev/dl/go1.23.4.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz

# Add to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH=$PATH:/usr/local/go/bin
```

**Option 2: Package Manager**
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install golang-go

# CentOS/RHEL/Fedora
sudo dnf install golang
# or
sudo yum install golang
```

**Option 3: Version Manager (gvm)**
```bash
# Install gvm
bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
# Install Go 1.23
gvm install go1.23.4
gvm use go1.23.4 --default
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
# Should output: go version go1.23.x linux/amd64 (or your platform)
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

These labs were originally designed for Go 1.15 but have been updated to work with Go 1.23.x. The codebase uses:

- Standard Go libraries (no external dependencies)
- RPC communication
- File I/O operations
- Concurrent programming with goroutines and channels

### Troubleshooting

**Common Issues:**

1. **"go: command not found"**
   - Ensure Go is in your PATH
   - Restart your terminal after installation

2. **Version mismatch**
   - Check with `go version`
   - Update to Go 1.23.x if using an older version

3. **Permission issues (Linux/macOS)**
   - Use `sudo` for system-wide installation
   - Or install to user directory and update PATH

4. **Module vs GOPATH conflicts**
   - Use `go mod init` for modern development
   - Clear GOPATH if using modules

### Getting Started with Labs

1. Ensure Go 1.23.x is installed and working
2. Navigate to the lab directory:
   ```bash
   cd app/01_Practice-Labs/src/main
   ```
3. Run the first lab:
   ```bash
   go run wc.go master kjv12.txt sequential
   ```

For more information about Go, visit [go.dev](https://go.dev/).
