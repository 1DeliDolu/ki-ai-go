#!/bin/bash

echo "🔨 Building Local AI Project Backend..."
echo ""

# Check if Go is installed
if ! command -v go >/dev/null 2>&1; then
    echo "❌ Go is not installed or not in PATH"
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi

echo "✅ Go version:"
go version
echo ""

# Set build environment for Linux (no CGO needed)
echo "🔧 Setting build environment..."
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

# Create bin directory if it doesn't exist
mkdir -p bin

# Clean previous builds
rm -f bin/server

echo "🔨 Building backend server..."
echo "   - Target: Linux x64"
echo "   - Output: bin/server"
echo "   - Mode: Release build"
echo ""

# Build with optimizations (no CGO dependencies)
if go build -ldflags "-s -w -X main.version=1.0.0" -trimpath -o bin/server ./cmd/server; then
    echo ""
    echo "✅ Build successful! (No CGO dependencies)"
    
    # Show file size
    if command -v du >/dev/null 2>&1; then
        size=$(du -h bin/server | cut -f1)
        echo "📏 Size: $size"
    fi
    
    # Make executable
    chmod +x bin/server
    
    echo ""
    echo "🎯 Next steps:"
    echo "1. Ensure Ollama is running: ollama serve"
    echo "2. Check your models: ollama list"
    echo "3. Start the server: ./start.sh"
    echo "4. Test API: curl http://localhost:8082/health"
    echo ""
    echo "💡 Your available models:"
    echo "    🚀 nemotron-nano    - NVIDIA Llama 3.1 Nemotron Nano 4B"
    echo "    🧠 neural-chat     - Intel Neural Chat 7B Q5_0"
    echo "    💬 openchat        - OpenChat 3.5 Q5_K_M"
    echo "    🦙 llama2-chat     - Llama 2 7B Chat Q4_K_M"
    echo "    🔬 phi2            - Microsoft Phi-2 Q8_0"
    echo ""
    echo "💡 Features:"
    echo "    ✅ In-memory database (no SQLite dependency)"
    echo "    ✅ Static binary (CGO_ENABLED=0)"
    echo "    ✅ All your downloaded models supported"
else
    echo ""
    echo "❌ Build failed!"
    echo ""
    echo "🔍 Troubleshooting:"
    echo "1. Check Go modules: go mod tidy"
    echo "2. Download dependencies: go mod download"
    echo "3. Verify Go version: go version"
    echo "4. Check for syntax errors: go vet ./..."
    echo ""
    exit 1
fi

echo "🏁 Build process completed."
