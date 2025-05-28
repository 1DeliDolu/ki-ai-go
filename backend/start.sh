#!/bin/bash

echo "🚀 Starting Local AI Backend Server..."
echo ""

# Check if Ollama is running
echo "🔍 Checking Ollama service..."
if curl -s http://localhost:11434/api/tags >/dev/null 2>&1; then
    echo "✅ Ollama is running on http://localhost:11434"
else
    echo "❌ Ollama is not running"
    echo "💡 Please start Ollama first:"
    echo "   In WSL: ollama serve"
    echo "   Or run: ./start_ollama.sh"
    echo ""
    exit 1
fi

# Check if server executable exists
if [ ! -f "bin/server" ]; then
    echo "❌ Server executable not found"
    echo "🔨 Building server first..."
    if ! ./build.sh; then
        echo "❌ Build failed, cannot start server"
        exit 1
    fi
fi

echo ""
echo "🌐 Starting backend server..."
echo "    API: http://localhost:8082/api/v1"
echo "    Health: http://localhost:8082/health"
echo "    Models: http://localhost:8082/api/v1/models"
echo ""
echo "📋 Available models in Ollama:"
ollama list 2>/dev/null | grep -E "(nemotron-nano|neural-chat|openchat|llama2-chat|phi2)" || echo "No custom models found"
echo ""
echo "💡 Press Ctrl+C to stop the server"
echo ""

# Start the server
./bin/server
echo ""
echo "🛑 Server stopped."
