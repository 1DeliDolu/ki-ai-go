#!/bin/bash

echo "🔍 Checking Document Uploads Locally"
echo "===================================="

# Check different possible upload locations
LOCATIONS=(
    "./test_documents"
    "$HOME/.local-ai-project/test_documents"
    "$HOME/.local-ai-project/uploads"
    "/tmp/local-ai-uploads"
    "./uploads"
)

echo "📁 Checking upload locations..."
for location in "${LOCATIONS[@]}"; do
    echo "Checking: $location"
    if [ -d "$location" ]; then
        echo "✅ Directory exists"
        file_count=$(ls -1 "$location" 2>/dev/null | wc -l)
        echo "   Files: $file_count"
        if [ $file_count -gt 0 ]; then
            echo "   Contents:"
            ls -la "$location" | head -10
        fi
    else
        echo "❌ Directory not found"
    fi
    echo "----------------------------------------"
done

echo ""
echo "🔧 Process and port check..."
echo "Server processes:"
ps aux | grep -E "(server|8082)" | grep -v grep

echo ""
echo "Port 8082 usage:"
netstat -tlnp 2>/dev/null | grep :8082 || ss -tlnp 2>/dev/null | grep :8082

echo ""
echo "📋 Config check..."
if [ -f "./internal/config/config.go" ]; then
    echo "✅ Config file found"
    echo "Upload paths in config:"
    grep -E "(UploadsPath|TestDocumentsPath)" ./internal/config/config.go 2>/dev/null || echo "Not found in config"
else
    echo "❌ Config file not found"
fi

echo ""
echo "🧪 Manual API test (without firewall issues)..."
echo "Using wget instead of curl:"

if command -v wget >/dev/null 2>&1; then
    echo "Testing with wget..."
    wget -q -O - "http://localhost:8082/health" 2>/dev/null | head -20
else
    echo "wget not available, trying nc..."
    if command -v nc >/dev/null 2>&1; then
        echo -e "GET /health HTTP/1.0\r\n\r\n" | nc localhost 8082 | head -20
    else
        echo "No alternative tools available"
    fi
fi

echo ""
echo "✅ Upload check completed!"
