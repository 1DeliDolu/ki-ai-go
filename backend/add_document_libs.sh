#!/bin/bash

echo "📚 Adding document processing libraries one by one..."
echo "===================================================="

cd "$(dirname "$0")"

# Function to test if a library can be added
test_library() {
    local lib=$1
    local name=$2
    
    echo ""
    echo "Testing $name ($lib)..."
    
    if go get "$lib" 2>/dev/null; then
        echo "✅ $name added successfully"
        return 0
    else
        echo "❌ $name failed"
        return 1
    fi
}

# Test basic libraries first
test_library "golang.org/x/net/html" "HTML parser"
test_library "github.com/h2non/filetype" "File type detection"

# Test markdown
test_library "github.com/yuin/goldmark" "Goldmark (Markdown)"

# Test HTML processing
test_library "github.com/PuerkitoBio/goquery" "goquery (HTML)"

# Alternative PDF libraries
echo ""
echo "🔍 Testing PDF libraries..."
test_library "github.com/unidoc/unipdf/v3" "UniPDF" || \
test_library "github.com/gen2brain/go-fitz" "go-fitz" || \
test_library "github.com/ledongthuc/pdf" "ledongthuc/pdf" || \
echo "⚠️  No PDF library available"

# Alternative DOCX libraries  
echo ""
echo "📄 Testing DOCX libraries..."
test_library "github.com/unidoc/unioffice" "UniOffice" || \
test_library "github.com/nguyenthenguyen/docx" "nguyenthenguyen/docx" || \
test_library "github.com/fumiama/go-docx" "go-docx" || \
echo "⚠️  No DOCX library available"

echo ""
echo "🔄 Running go mod tidy..."
go mod tidy

echo ""
echo "🧪 Testing build with new dependencies..."
if go build -o /tmp/test_build ./cmd/server 2>/dev/null; then
    echo "✅ Build successful with document libraries!"
    rm -f /tmp/test_build
else
    echo "❌ Build failed, checking errors..."
    go build ./cmd/server
fi

echo ""
echo "📋 Current dependencies:"
go list -m all | grep -v indirect

echo ""
echo "✅ Document library installation completed!"
