#!/bin/bash

echo "🧪 Testing Local AI Backend..."
echo ""

# Test if server is running
echo "🔍 Testing server health..."
if curl -s http://localhost:8082/health >/dev/null; then
    echo "✅ Server is responding"
    echo "Response:"
    curl -s http://localhost:8082/health | jq . 2>/dev/null || curl -s http://localhost:8082/health
else
    echo "❌ Server is not responding"
    echo "💡 Make sure the server is running: ./start.sh"
    exit 1
fi

echo ""
echo "🤖 Testing models endpoint..."
echo "Response:"
curl -s http://localhost:8082/api/v1/models | jq . 2>/dev/null || curl -s http://localhost:8082/api/v1/models
echo ""

echo ""
echo "📄 Testing documents endpoint..."
echo "Response:"
curl -s http://localhost:8082/api/v1/documents | jq . 2>/dev/null || curl -s http://localhost:8082/api/v1/documents
echo ""

echo ""
echo "📄 Testing document types endpoint..."
echo "Response:"
curl -s http://localhost:8082/api/v1/documents/types | jq . 2>/dev/null || curl -s http://localhost:8082/api/v1/documents/types
echo ""

echo ""
echo "🤖 Testing models initialization..."
echo "Response:"
curl -s -X POST http://localhost:8082/api/v1/models/initialize | jq . 2>/dev/null || curl -s -X POST http://localhost:8082/api/v1/models/initialize

echo ""
echo "📋 Testing model types endpoint..."
echo "Response:"
curl -s http://localhost:8082/api/v1/models/types | jq . 2>/dev/null || curl -s http://localhost:8082/api/v1/models/types
echo ""

echo ""
echo "🔍 Testing chat models..."
echo "Response:"
curl -s http://localhost:8082/api/v1/models/type/chat | jq . 2>/dev/null || curl -s http://localhost:8082/api/v1/models/type/chat
echo ""

echo ""
echo "📄 Testing specific model info..."
echo "Response:"
curl -s http://localhost:8082/api/v1/models/tinyllama-1.1b | jq . 2>/dev/null || curl -s http://localhost:8082/api/v1/models/tinyllama-1.1b
echo ""

echo ""
echo "📄 Testing document content extraction..."
# Create test files first
mkdir -p test_files

# Create test TXT file
echo "This is a test text file for document processing." > test_files/test.txt

# Create test JSON file
echo '{"name": "test", "type": "json", "data": [1,2,3]}' > test_files/test.json

# Create test CSV file
echo "name,age,city
John,25,NYC
Jane,30,LA" > test_files/test.csv

# Upload and test TXT file
echo "Testing TXT file upload..."
curl -s -X POST -F "file=@test_files/test.txt" http://localhost:8082/api/v1/documents/upload | jq . 2>/dev/null || curl -s -X POST -F "file=@test_files/test.txt" http://localhost:8082/api/v1/documents/upload

echo ""
echo "Testing JSON file upload..."
curl -s -X POST -F "file=@test_files/test.json" http://localhost:8082/api/v1/documents/upload | jq . 2>/dev/null || curl -s -X POST -F "file=@test_files/test.json" http://localhost:8082/api/v1/documents/upload

echo ""
echo "Testing CSV file upload..."
curl -s -X POST -F "file=@test_files/test.csv" http://localhost:8082/api/v1/documents/upload | jq . 2>/dev/null || curl -s -X POST -F "file=@test_files/test.csv" http://localhost:8082/api/v1/documents/upload

# Clean up test files
rm -rf test_files

echo ""
echo "📄 Testing frontend document upload simulation..."
# Simulate frontend document upload
echo "Testing frontend document upload..."
curl -s -X POST -F "file=@test_files/test.txt" http://localhost:8082/api/v1/documents/upload | jq . 2>/dev/null || curl -s -X POST -F "file=@test_files/test.txt" http://localhost:8082/api/v1/documents/upload

echo ""
echo "📁 Testing test documents endpoint..."
curl -s http://localhost:8082/api/v1/documents/test | jq . 2>/dev/null || curl -s http://localhost:8082/api/v1/documents/test

echo ""
echo "🧹 Testing test documents cleanup..."
curl -s -X POST http://localhost:8082/api/v1/cleanup/test-documents | jq . 2>/dev/null || curl -s -X POST http://localhost:8082/api/v1/cleanup/test-documents

echo ""
echo "✅ API tests completed!"
echo ""
echo "🌐 You can also test in browser:"
echo "    http://localhost:8082/health"
echo "    http://localhost:8082/api/v1/models"
echo ""

# Test Ollama connection
echo "🦙 Testing Ollama connection..."
if curl -s http://localhost:11434/api/tags >/dev/null; then
    echo "✅ Ollama is accessible"
    echo "Available models:"
    ollama list 2>/dev/null | head -10
else
    echo "⚠️  Ollama is not accessible"
    echo "💡 Start Ollama: ollama serve"
fi
