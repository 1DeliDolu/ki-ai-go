#!/bin/bash

echo "🔧 Fixing import paths in Go files"
echo "=================================="

MODULE_NAME="github.com/1DeliDolu/go_mustAI/local-ai-project/backend"
OLD_IMPORT="local-ai-project/backend"

echo "📁 Scanning for files with incorrect imports..."

# Find all Go files with old import paths
FILES_TO_FIX=$(grep -r "$OLD_IMPORT" . --include="*.go" -l 2>/dev/null)

if [ -z "$FILES_TO_FIX" ]; then
    echo "✅ No files need fixing - all imports are correct"
    exit 0
fi

echo "📄 Files that need fixing:"
echo "$FILES_TO_FIX"
echo ""

echo "🔄 Fixing import paths..."

# Fix each file
while IFS= read -r file; do
    if [ -f "$file" ]; then
        echo "  Fixing: $file"
        # Use sed to replace the import path
        sed -i "s|$OLD_IMPORT|$MODULE_NAME|g" "$file"
    fi
done <<< "$FILES_TO_FIX"

echo ""
echo "🧹 Cleaning up Go modules..."
go mod tidy

echo ""
echo "🧪 Testing build..."
if go build -o /tmp/test_build ./cmd/server 2>/dev/null; then
    echo "✅ Build successful!"
    rm -f /tmp/test_build
else
    echo "❌ Build failed. Checking for remaining issues..."
    go build ./cmd/server
fi

echo ""
echo "✅ Import path fixing complete!"
