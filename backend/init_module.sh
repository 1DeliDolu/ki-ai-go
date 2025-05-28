#!/bin/bash

echo "🚀 Initializing Go module for Local AI Project Backend"
echo "======================================================"

# Check if we're in the right directory
if [ ! -f "cmd/server/main.go" ]; then
    echo "❌ Error: Not in backend directory or main.go not found"
    echo "Please run this script from the backend directory"
    exit 1
fi

echo "📁 Current directory: $(pwd)"
echo ""

# Remove existing go.mod and go.sum if they exist
if [ -f "go.mod" ]; then
    echo "🗑️  Removing existing go.mod"
    rm -f go.mod
fi

if [ -f "go.sum" ]; then
    echo "🗑️  Removing existing go.sum"
    rm -f go.sum
fi

echo "📦 Initializing Go module..."
go mod init github.com/1DeliDolu/go_mustAI/local-ai-project/backend

if [ $? -ne 0 ]; then
    echo "❌ Failed to initialize Go module"
    exit 1
fi

echo "✅ Go module initialized"
echo ""

echo "📋 Adding required dependencies..."
echo "================================="

# Add main dependencies
echo "Adding Gin web framework..."
go get github.com/gin-gonic/gin@latest

echo "Adding CORS middleware..."
go get github.com/gin-contrib/cors@latest

# Add other useful dependencies
echo "Adding additional utilities..."
go get golang.org/x/crypto@latest
go get golang.org/x/net@latest

echo ""
echo "🔄 Running go mod tidy..."
go mod tidy

if [ $? -ne 0 ]; then
    echo "❌ go mod tidy failed"
    exit 1
fi

echo ""
echo "📊 Checking module status..."
echo "=========================="
echo "Go version: $(go version)"
echo "Module name: $(grep '^module' go.mod)"
echo "Dependencies:"
go list -m all | head -10

echo ""
echo "🧪 Testing build..."
echo "=================="
echo "Attempting to build the project..."

if go build -o /tmp/test_build ./cmd/server; then
    echo "✅ Build test successful!"
    rm -f /tmp/test_build
else
    echo "❌ Build test failed. Checking for common issues..."
    echo ""
    echo "🔍 Checking import paths..."
    
    # Check for incorrect import paths
    if grep -r "local-ai-project/backend" . --include="*.go"; then
        echo "⚠️  Found old import paths. They should use:"
        echo "   github.com/1DeliDolu/go_mustAI/local-ai-project/backend"
        echo ""
        echo "🔧 Auto-fixing import paths..."
        
        # Fix import paths in all Go files
        find . -name "*.go" -type f -exec sed -i 's|local-ai-project/backend|github.com/1DeliDolu/go_mustAI/local-ai-project/backend|g' {} \;
        
        echo "✅ Import paths fixed"
        echo ""
        echo "🔄 Running go mod tidy again..."
        go mod tidy
        
        echo "🧪 Testing build again..."
        if go build -o /tmp/test_build ./cmd/server; then
            echo "✅ Build successful after fixing imports!"
            rm -f /tmp/test_build
        else
            echo "❌ Build still failing. Manual intervention needed."
        fi
    fi
fi

echo ""
echo "📁 Project structure:"
echo "===================="
tree -I 'node_modules|dist|build|bin|tmp' -L 3 2>/dev/null || find . -type d -name ".*" -prune -o -type d -print | head -20

echo ""
echo "🎯 Next steps:"
echo "============="
echo "1. Build the project: ./build.sh"
echo "2. Start Ollama: ./start_ollama.sh"  
echo "3. Run the server: ./start.sh"
echo "4. Test the API: ./test.sh"

echo ""
echo "✅ Go module initialization complete!"
