#!/bin/bash

echo "📥 Downloading missing model files..."
echo "===================================="

MODELS_DIR="/mnt/d/Praxis/KI/lokaleKI/go_mustAI/local-ai-project/models"

# Ensure models directory exists
mkdir -p "$MODELS_DIR"
cd "$MODELS_DIR" || exit 1

echo "📁 Working in: $MODELS_DIR"
echo ""

# Function to download with progress
download_model() {
    local name="$1"
    local url="$2"
    local filename="$3"
    local description="$4"
    
    echo "🔍 Checking $description..."
    
    if [ -f "$filename" ]; then
        echo "✅ $filename already exists"
        return 0
    fi
    
    echo "📥 Downloading $description..."
    echo "   URL: $url"
    echo "   File: $filename"
    echo ""
    
    # Download with progress bar
    if wget --progress=bar:force:noscroll "$url" -O "$filename.tmp" 2>&1; then
        mv "$filename.tmp" "$filename"
        echo ""
        echo "✅ Successfully downloaded: $filename"
        
        # Show file size
        if command -v du >/dev/null 2>&1; then
            size=$(du -h "$filename" | cut -f1)
            echo "   Size: $size"
        fi
    else
        echo ""
        echo "❌ Failed to download: $filename"
        rm -f "$filename.tmp"
        return 1
    fi
    echo ""
}

# Download OpenChat 3.5 (GGUF format)
download_model \
    "openchat" \
    "https://huggingface.co/TheBloke/openchat_3.5-GGUF/resolve/main/openchat_3.5.q4_k_m.gguf" \
    "openchat-3.5-1210.Q4_K_M.gguf" \
    "OpenChat 3.5 (4-bit quantized)"

# Download Microsoft Phi-2 (GGUF format)
download_model \
    "phi2" \
    "https://huggingface.co/TheBloke/phi-2-GGUF/resolve/main/phi-2.q4_k_m.gguf" \
    "phi-2.Q4_K_M.gguf" \
    "Microsoft Phi-2 (4-bit quantized)"

# Alternative: Try different Phi-2 sources if the first fails
if [ ! -f "phi-2.Q4_K_M.gguf" ]; then
    echo "🔄 Trying alternative Phi-2 source..."
    download_model \
        "phi2-alt" \
        "https://huggingface.co/microsoft/phi-2-gguf/resolve/main/phi-2.q4_k_m.gguf" \
        "phi-2.Q4_K_M.gguf" \
        "Microsoft Phi-2 (alternative source)"
fi

# Alternative: Try different OpenChat sources if the first fails
if [ ! -f "openchat-3.5-1210.Q4_K_M.gguf" ]; then
    echo "🔄 Trying alternative OpenChat source..."
    download_model \
        "openchat-alt" \
        "https://huggingface.co/openchat/openchat_3.5/resolve/main/openchat_3.5.q4_k_m.gguf" \
        "openchat-3.5-1210.Q4_K_M.gguf" \
        "OpenChat 3.5 (alternative source)"
fi

echo "🎯 Download Summary:"
echo "==================="

check_file() {
    local filename="$1"
    local model_name="$2"
    
    if [ -f "$filename" ]; then
        size=$(du -h "$filename" 2>/dev/null | cut -f1 || echo "unknown")
        echo "✅ $model_name: $filename ($size)"
    else
        echo "❌ $model_name: Missing"
    fi
}

check_file "llama-2-7b-chat.Q4_K_M.gguf" "Llama 2 Chat"
check_file "neural-chat-7b-v3-1.Q4_K_M.gguf" "Neural Chat"
check_file "openchat-3.5-1210.Q4_K_M.gguf" "OpenChat 3.5"
check_file "phi-2.Q4_K_M.gguf" "Phi-2"

echo ""
echo "💡 Manual download options if automatic download fails:"
echo "======================================================"

if [ ! -f "openchat-3.5-1210.Q4_K_M.gguf" ]; then
    echo ""
    echo "🔗 OpenChat 3.5 alternatives:"
    echo "   1. https://huggingface.co/TheBloke/openchat_3.5-GGUF/blob/main/openchat_3.5.q4_k_m.gguf"
    echo "   2. https://huggingface.co/openchat/openchat_3.5-GGUF/blob/main/openchat_3.5.q4_k_m.gguf"
    echo ""
    echo "   Download command:"
    echo "   wget 'https://huggingface.co/TheBloke/openchat_3.5-GGUF/resolve/main/openchat_3.5.q4_k_m.gguf' -O openchat-3.5-1210.Q4_K_M.gguf"
fi

if [ ! -f "phi-2.Q4_K_M.gguf" ]; then
    echo ""
    echo "🔗 Phi-2 alternatives:"
    echo "   1. https://huggingface.co/TheBloke/phi-2-GGUF/blob/main/phi-2.q4_k_m.gguf"
    echo "   2. https://huggingface.co/microsoft/phi-2-gguf/blob/main/phi-2.q4_k_m.gguf"
    echo ""
    echo "   Download command:"
    echo "   wget 'https://huggingface.co/TheBloke/phi-2-GGUF/resolve/main/phi-2.q4_k_m.gguf' -O phi-2.Q4_K_M.gguf"
fi

echo ""
echo "🚀 After downloading, restart your backend server to detect the new models!"
