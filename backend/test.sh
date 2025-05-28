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
