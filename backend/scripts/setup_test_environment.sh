#!/bin/bash

echo "🔧 Setting up test environment..."
echo "================================"

cd "$(dirname "$0")/.."

# Make scripts executable
chmod +x scripts/*.sh

# Create test documents
./scripts/create_test_documents.sh

# Install dependencies
echo ""
echo "📦 Installing Go dependencies..."
go get github.com/nguyenthenguyen/docx
go get github.com/ledongthuc/pdf
go get github.com/PuerkitoBio/goquery

# Run go mod tidy
go mod tidy

echo ""
echo "🧪 Building project..."
go build -o bin/test_server ./cmd/server

echo ""
echo "✅ Test environment ready!"
echo ""
echo "🚀 Next steps:"
echo "1. Start server: ./start.sh"
echo "2. Test uploads: ./scripts/test_document_upload.sh"
echo "3. Check http://localhost:8082/api/v1/documents/types"
