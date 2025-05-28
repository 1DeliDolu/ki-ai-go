#!/bin/bash

echo "🔍 Checking Ollama status..."
echo ""

# Check if Ollama process is running
if pgrep -x "ollama" > /dev/null; then
    echo "✅ Ollama process is running (PID: $(pgrep -x ollama))"
    
    # Check if API is responding
    if curl -s http://localhost:11434/api/tags > /dev/null 2>&1; then
        echo "✅ Ollama API is responding on http://localhost:11434"
        
        echo ""
        echo "📋 Available models in Ollama:"
        ollama list
        
        echo ""
        echo "📁 Your downloaded model files:"
        ls -la /mnt/d/Praxis/KI/lokaleKI/go_mustAI/local-ai-project/models/*.gguf 2>/dev/null || echo "No .gguf files found"
        
        echo ""
        echo "🚀 Backend can now connect to Ollama!"
        
    else
        echo "❌ Ollama process running but API not responding"
        echo "Try restarting Ollama with: pkill ollama && ollama serve"
    fi
else
    echo "❌ Ollama is not running"
    echo "Start it with: ./start_ollama.sh"
fi

echo ""
echo "🌐 Test Ollama API:"
echo "curl http://localhost:11434/api/tags"
