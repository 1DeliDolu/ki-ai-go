#!/bin/bash

echo "🔍 Debugging Model Detection"
echo "============================"

MODELS_DIR="/mnt/d/Praxis/KI/lokaleKI/go_mustAI/local-ai-project/models"

echo "📁 Models directory: $MODELS_DIR"
echo ""

if [ ! -d "$MODELS_DIR" ]; then
    echo "❌ Models directory not found!"
    exit 1
fi

echo "📄 All files in models directory:"
echo "================================="
ls -la "$MODELS_DIR"
echo ""

echo "🤖 GGUF files specifically:"
echo "=========================="
find "$MODELS_DIR" -name "*.gguf" -type f -exec ls -lh {} \; 2>/dev/null
echo ""

echo "🎯 Expected vs Found:"
echo "==================="

check_model() {
    local expected="$1"
    local model_name="$2"
    
    echo "Checking $model_name:"
    echo "  Expected: $expected"
    
    if [ -f "$MODELS_DIR/$expected" ]; then
        size=$(du -h "$MODELS_DIR/$expected" | cut -f1)
        echo "  ✅ Found: $expected ($size)"
    else
        echo "  ❌ Not found: $expected"
        
        # Look for similar files
        echo "  🔍 Looking for similar files:"
        case $model_name in
            "nemotron-nano")
                find "$MODELS_DIR" -iname "*nemotron*" -o -iname "*nvidia*" -o -iname "*nano*" | while read file; do
                    echo "    Possible: $(basename "$file")"
                done
                ;;
            "openchat")
                find "$MODELS_DIR" -iname "*openchat*" | while read file; do
                    echo "    Possible: $(basename "$file")"
                done
                ;;
            "phi2")
                find "$MODELS_DIR" -iname "*phi*" | while read file; do
                    echo "    Possible: $(basename "$file")"
                done
                ;;
        esac
    fi
    echo ""
}

check_model "nvidia_Llama-3.1-Nemotron-Nano-4B-v1.1-bf16.gguf" "nemotron-nano"
check_model "neural-chat-7b-v3-1.Q5_0.gguf" "neural-chat"
check_model "openchat-3.5-0106.Q5_K_M.gguf" "openchat"
check_model "llama-2-7b-chat.Q4_K_M.gguf" "llama2-chat"
check_model "phi-2.Q8_0.gguf" "phi2"

echo "🔄 Testing backend model detection:"
echo "==================================="
echo "Starting backend temporarily to check model detection..."

# Test API if backend is running
if curl -s http://localhost:8082/api/v1/models >/dev/null 2>&1; then
    echo "✅ Backend is running, testing models API:"
    curl -s http://localhost:8082/api/v1/models | jq . 2>/dev/null || curl -s http://localhost:8082/api/v1/models
else
    echo "⚠️  Backend is not running. Start it with:"
    echo "   ./build.sh && ./start.sh"
fi

echo ""
echo "💡 If models are not detected:"
echo "1. Check file permissions: chmod 644 $MODELS_DIR/*.gguf"
echo "2. Restart backend server"
echo "3. Check backend logs for errors"
