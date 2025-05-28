#!/bin/bash

echo "⚡ Quick Model Setup for WSL Debian"
echo "=================================="

MODELS_DIR="/mnt/d/Praxis/KI/lokaleKI/go_mustAI/local-ai-project/models"

# Ensure we're in the right directory
mkdir -p "$MODELS_DIR"
cd "$MODELS_DIR" || exit 1

echo "📁 Working in: $MODELS_DIR"
echo ""

# Check what we have and what we need
echo "🔍 Current model status:"
echo "======================="

models_needed=()

if [ -f "llama-2-7b-chat.Q4_K_M.gguf" ]; then
    echo "✅ Llama 2 Chat: Available"
else
    echo "❌ Llama 2 Chat: Missing"
    models_needed+=("llama2")
fi

if [ -f "neural-chat-7b-v3-1.Q4_K_M.gguf" ]; then
    echo "✅ Neural Chat: Available"
else
    echo "❌ Neural Chat: Missing"
    models_needed+=("neural")
fi

if [ -f "openchat-3.5-1210.Q4_K_M.gguf" ]; then
    echo "✅ OpenChat 3.5: Available"
else
    echo "❌ OpenChat 3.5: Missing"
    models_needed+=("openchat")
fi

if [ -f "phi-2.Q4_K_M.gguf" ]; then
    echo "✅ Phi-2: Available"
else
    echo "❌ Phi-2: Missing"  
    models_needed+=("phi2")
fi

echo ""

if [ ${#models_needed[@]} -eq 0 ]; then
    echo "🎉 All models are available!"
    echo "🔄 Run your backend server to use them."
    exit 0
fi

echo "📥 Need to download ${#models_needed[@]} models"
echo ""

# Quick download function with better error handling
quick_download() {
    local url="$1"
    local filename="$2"
    local description="$3"
    
    echo "📥 Downloading $description..."
    echo "   🔗 URL: $url"
    echo "   📄 File: $filename"
    
    # Use curl if available, otherwise wget
    if command -v curl >/dev/null 2>&1; then
        if curl -L --progress-bar "$url" -o "$filename.tmp"; then
            mv "$filename.tmp" "$filename"
            echo "   ✅ Downloaded successfully!"
            return 0
        else
            echo "   ❌ Download failed with curl"
            rm -f "$filename.tmp"
            return 1
        fi
    elif command -v wget >/dev/null 2>&1; then
        if wget --progress=bar:force:noscroll "$url" -O "$filename.tmp"; then
            mv "$filename.tmp" "$filename"
            echo "   ✅ Downloaded successfully!"
            return 0
        else
            echo "   ❌ Download failed with wget"
            rm -f "$filename.tmp"
            return 1
        fi
    else
        echo "   ❌ Neither curl nor wget available"
        return 1
    fi
}

# Download missing models
for model in "${models_needed[@]}"; do
    case $model in
        "openchat")
            echo "🔽 Downloading OpenChat 3.5..."
            if ! quick_download \
                "https://huggingface.co/TheBloke/openchat_3.5-GGUF/resolve/main/openchat_3.5.q4_k_m.gguf" \
                "openchat-3.5-1210.Q4_K_M.gguf" \
                "OpenChat 3.5"; then
                echo "⚠️  Primary download failed, trying alternative..."
                quick_download \
                    "https://huggingface.co/openchat/openchat_3.5/resolve/main/openchat_3.5.q4_k_m.gguf" \
                    "openchat-3.5-1210.Q4_K_M.gguf" \
                    "OpenChat 3.5 (alternative)"
            fi
            ;;
        "phi2")
            echo "🔽 Downloading Phi-2..."
            if ! quick_download \
                "https://huggingface.co/TheBloke/phi-2-GGUF/resolve/main/phi-2.q4_k_m.gguf" \
                "phi-2.Q4_K_M.gguf" \
                "Microsoft Phi-2"; then
                echo "⚠️  Primary download failed, trying alternative..."
                quick_download \
                    "https://huggingface.co/microsoft/phi-2-gguf/resolve/main/phi-2.q4_k_m.gguf" \
                    "phi-2.Q4_K_M.gguf" \
                    "Microsoft Phi-2 (alternative)"
            fi
            ;;
    esac
    echo ""
done

echo "🎯 Final Status Check:"
echo "====================="

all_available=true

for model_file in "llama-2-7b-chat.Q4_K_M.gguf" "neural-chat-7b-v3-1.Q4_K_M.gguf" "openchat-3.5-1210.Q4_K_M.gguf" "phi-2.Q4_K_M.gguf"; do
    if [ -f "$model_file" ]; then
        size=$(du -h "$model_file" 2>/dev/null | cut -f1 || echo "unknown")
        echo "✅ $model_file ($size)"
    else
        echo "❌ $model_file (missing)"
        all_available=false
    fi
done

echo ""
if [ "$all_available" = true ]; then
    echo "🎉 All models are now available!"
    echo "🚀 Start your backend server with: go run cmd/server/main.go"
else
    echo "⚠️  Some models are still missing."
    echo "💡 Try manual download or check your internet connection."
fi

echo ""
echo "📋 Next steps:"
echo "1. Go to backend directory: cd ../backend"
echo "2. Start the server: go run cmd/server/main.go"
echo "3. Check model availability in the frontend"
