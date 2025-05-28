#!/bin/bash

echo "🔍 Testing Go dependencies for document processing..."
echo "=================================================="

cd "$(dirname "$0")"

echo ""
echo "📦 Testing PDF library (github.com/ledongthuc/pdf)..."
go get github.com/ledongthuc/pdf
if [ $? -eq 0 ]; then
    echo "✅ PDF library downloaded successfully"
else
    echo "❌ PDF library failed"
fi

echo ""
echo "📄 Testing DOCX library (github.com/nguyenthenguyen/docx)..."
go get github.com/nguyenthenguyen/docx
if [ $? -eq 0 ]; then
    echo "✅ DOCX library downloaded successfully"
else
    echo "❌ DOCX library failed"
fi

echo ""
echo "📝 Testing Markdown library (github.com/yuin/goldmark)..."
go get github.com/yuin/goldmark
if [ $? -eq 0 ]; then
    echo "✅ Markdown library downloaded successfully"
else
    echo "❌ Markdown library failed"
fi

echo ""
echo "🌐 Testing HTML libraries..."
go get golang.org/x/net/html
if [ $? -eq 0 ]; then
    echo "✅ HTML parser (golang.org/x/net/html) downloaded successfully"
else
    echo "❌ HTML parser failed"
fi

go get github.com/PuerkitoBio/goquery
if [ $? -eq 0 ]; then
    echo "✅ goquery library downloaded successfully"
else
    echo "❌ goquery library failed"
fi

echo ""
echo "🔍 Testing file type detection library..."
go get github.com/h2non/filetype
if [ $? -eq 0 ]; then
    echo "✅ File type detection library downloaded successfully"
else
    echo "❌ File type detection library failed"
fi

echo ""
echo "🧹 Running go mod tidy..."
go mod tidy

echo ""
echo "📋 Current dependencies:"
go list -m all | grep -E "(pdf|docx|goldmark|goquery|filetype|html)"

echo ""
echo "🧪 Testing simple imports..."
cat > test_imports.go << 'EOF'
package main

import (
    "fmt"
    
    "github.com/ledongthuc/pdf"
    "github.com/nguyenthenguyen/docx"  
    "github.com/yuin/goldmark"
    "github.com/PuerkitoBio/goquery"
    "github.com/h2non/filetype"
    "golang.org/x/net/html"
)

func main() {
    fmt.Println("All libraries imported successfully!")
}
EOF

echo "Attempting to build test imports..."
if go build -o test_imports test_imports.go; then
    echo "✅ All libraries can be imported successfully!"
    rm -f test_imports test_imports.go
else
    echo "❌ Some libraries have import issues"
    rm -f test_imports.go
fi

echo ""
echo "🎯 Dependency test completed!"
