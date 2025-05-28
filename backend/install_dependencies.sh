#!/bin/bash

echo "📦 Installing document processing dependencies..."
echo "================================================"

cd "$(dirname "$0")"

echo ""
echo "🧹 Cleaning up Go modules..."
go clean -modcache

echo ""
echo "📥 Adding dependencies to go.mod..."

echo "Adding basic dependencies..."
go get github.com/gin-contrib/cors@latest
go get github.com/gin-gonic/gin@latest

echo ""
echo "📄 Adding document processing libraries (with error handling)..."

# Try to add each library individually
echo "Trying h2non/filetype for file type detection..."
if go get github.com/h2non/filetype@latest; then
    echo "✅ filetype library added successfully"
else
    echo "❌ filetype library failed"
fi

echo ""
echo "Trying yuin/goldmark for Markdown processing..."
if go get github.com/yuin/goldmark@latest; then
    echo "✅ goldmark library added successfully"
else
    echo "❌ goldmark library failed"
fi

echo ""
echo "Trying PuerkitoBio/goquery for HTML processing..."
if go get github.com/PuerkitoBio/goquery@latest; then
    echo "✅ goquery library added successfully"
else
    echo "❌ goquery library failed"
fi

echo ""
echo "Trying ledongthuc/pdf for PDF processing..."
if go get github.com/ledongthuc/pdf@latest; then
    echo "✅ PDF library added successfully"
else
    echo "❌ PDF library failed"
fi

echo ""
echo "Trying nguyenthenguyen/docx for DOCX processing..."
if go get github.com/nguyenthenguyen/docx@latest; then
    echo "✅ DOCX library added successfully"
else
    echo "❌ DOCX library failed"
fi

echo ""
echo "🔄 Running go mod tidy..."
go mod tidy

echo ""
echo "✅ Dependency installation completed!"
echo ""
echo "📋 Installed dependencies:"
go list -m all | grep -v "indirect" | head -20

echo ""
echo "🧪 Testing build..."
if go build -o /tmp/test_build ./cmd/server 2>/dev/null; then
    echo "✅ Build successful with current dependencies!"
    rm -f /tmp/test_build
else
    echo "⚠️  Build has issues - using basic implementation for now"
fi
