#!/bin/bash

echo "🐧 Setting up Local AI Project for WSL Debian..."
echo ""

# Check if we're running in WSL
if ! grep -q Microsoft /proc/version 2>/dev/null; then
    echo "⚠️  This script is designed for WSL (Windows Subsystem for Linux)"
    echo "Current system: $(uname -a)"
    echo "Proceeding anyway..."
fi

echo "📋 System Information:"
echo "====================="
echo "OS: $(lsb_release -d 2>/dev/null | cut -f2 || echo "Unknown")"
echo "Kernel: $(uname -r)"
echo "Architecture: $(uname -m)"
echo "Working directory: $(pwd)"
echo ""

# Install required dependencies
echo "📦 Installing dependencies..."
echo "============================="

# Update package list
sudo apt update

# Install required packages
packages=(
    "curl"
    "wget" 
    "git"
    "build-essential"
    "bc"
    "jq"
    "unzip"
)

for package in "${packages[@]}"; do
    if ! command -v "$package" >/dev/null 2>&1; then
        echo "Installing $package..."
        sudo apt install -y "$package"
    else
        echo "✅ $package already installed"
    fi
done

# Check Go installation
echo ""
echo "🔍 Checking Go installation..."
if command -v go >/dev/null 2>&1; then
    echo "✅ Go is installed: $(go version)"
else
    echo "❌ Go is not installed"
    echo "💡 Installing Go..."
    
    # Download and install Go
    GO_VERSION="1.21.5"
    wget -q "https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz"
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf "go${GO_VERSION}.linux-amd64.tar.gz"
    rm "go${GO_VERSION}.linux-amd64.tar.gz"
    
    # Add Go to PATH
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    export PATH=$PATH:/usr/local/go/bin
    
    echo "✅ Go installed: $(go version)"
fi

# Setup directories
echo ""
echo "📁 Setting up directories..."
echo "============================"

MODELS_DIR="/mnt/d/Praxis/KI/lokaleKI/go_mustAI/local-ai-project/models"
DATA_DIR="$HOME/.local-ai-project"

# Create models directory
if [ ! -d "$MODELS_DIR" ]; then
    echo "Creating models directory: $MODELS_DIR"
    mkdir -p "$MODELS_DIR"
    echo "✅ Models directory created"
else
    echo "✅ Models directory exists: $MODELS_DIR"
fi

# Create local data directory
if [ ! -d "$DATA_DIR" ]; then
    echo "Creating local data directory: $DATA_DIR"
    mkdir -p "$DATA_DIR"/{data,uploads,models}
    echo "✅ Local data directory created"
else
    echo "✅ Local data directory exists: $DATA_DIR"
fi

# Make scripts executable
echo ""
echo "🔧 Making scripts executable..."
echo "==============================="

SCRIPT_DIR="$(pwd)/scripts"
if [ -d "$SCRIPT_DIR" ]; then
    chmod +x "$SCRIPT_DIR"/*.sh
    echo "✅ All scripts are now executable"
    
    echo ""
    echo "📄 Available scripts:"
    ls -la "$SCRIPT_DIR"/*.sh
else
    echo "⚠️  Scripts directory not found: $SCRIPT_DIR"
fi

# Check Ollama installation
echo ""
echo "🦙 Checking Ollama..."
echo "===================="

if command -v ollama >/dev/null 2>&1; then
    echo "✅ Ollama is installed: $(ollama --version 2>/dev/null || echo "version unknown")"
    
    # Check if Ollama is running
    if curl -s http://localhost:11434/api/tags >/dev/null 2>&1; then
        echo "✅ Ollama service is running"
    else
        echo "⚠️  Ollama is installed but not running"
        echo "💡 Start it with: ollama serve"
    fi
else
    echo "❌ Ollama is not installed"
    echo "💡 Install it with: curl -fsSL https://ollama.com/install.sh | sh"
fi

# Display next steps
echo ""
echo "🎯 Setup Complete!"
echo "=================="
echo ""
echo "🔄 Next steps:"
echo "1. Check your models: ./scripts/check_models.sh"
echo "2. Fix model names if needed: ./scripts/fix_model_names.sh" 
echo "3. Start Ollama: ollama serve"
echo "4. Build and run backend: go run cmd/server/main.go"
echo ""
echo "📁 Important paths:"
echo "   Models: $MODELS_DIR"
echo "   Data: $DATA_DIR"
echo "   Scripts: $SCRIPT_DIR"
echo ""
echo "🌐 Once running, the API will be available at:"
echo "   http://localhost:8082/api/v1"
