#!/bin/bash

echo "🧹 Cleaning up Go modules and installing dependencies..."
echo "======================================================="

cd "$(dirname "$0")"

echo ""
echo "🗑️  Removing problematic cache..."
go clean -modcache

echo ""
echo "🔄 Resetting go.mod to clean state..."
rm -f go.sum

echo ""
echo "📦 Installing core dependencies only..."

# Install only the core dependencies that we know work
echo "Installing Gin framework..."
go get github.com/gin-gonic/gin@latest

echo "Installing CORS middleware..."
go get github.com/gin-contrib/cors@latest

echo "Installing PostgreSQL driver..."
go get github.com/lib/pq@latest

echo ""
echo "🧪 Testing basic build..."
if go build -o /tmp/test_build ./cmd/server 2>/dev/null; then
    echo "✅ Core build successful!"
    rm -f /tmp/test_build
else
    echo "❌ Build failed, let's see the error:"
    go build ./cmd/server
fi

echo ""
echo "🔄 Running go mod tidy..."
go mod tidy

echo ""
echo "📋 Current clean dependencies:"
go list -m all

echo ""
echo "✅ Cleanup and basic installation completed!"
