#!/bin/bash

echo "📥 Manual Model Download Guide"
echo "=============================="
echo ""

MODELS_DIR="/mnt/d/Praxis/KI/lokaleKI/go_mustAI/local-ai-project/models"

echo "📁 Models will be saved to: $MODELS_DIR"
echo ""

# Create models directory if it doesn't exist
mkdir -p "$MODELS_DIR"

echo "🔍 Checking current model status..."
echo "=================================="

check_model_status() {
    local filename="$1"
    local model_name="$2"
    
    if [ -f "$MODELS_DIR/$filename" ]; then
        size=$(du -h "$MODELS_DIR/$filename" 2>/dev/null | cut -f1 || echo "unknown")
        echo "✅ $model_name: Available ($size)"
        return 0
    else
        echo "❌ $model_name: Missing"
        return 1
    fi
}

# Check all models
llama_ok=false
neural_ok=false
openchat_ok=false
phi2_ok=false

if check_model_status "llama-2-7b-chat.Q4_K_M.gguf" "Llama 2 Chat"; then
    llama_ok=true
fi

if check_model_status "neural-chat-7b-v3-1.Q4_K_M.gguf" "Neural Chat"; then
    neural_ok=true
fi

if check_model_status "openchat-3.5-1210.Q4_K_M.gguf" "OpenChat 3.5"; then
    openchat_ok=true
fi

if check_model_status "phi-2.Q4_K_M.gguf" "Phi-2"; then
    phi2_ok=true
fi

echo ""
echo "📋 Manual Download Instructions"
echo "==============================="
echo ""

if [ "$openchat_ok" = false ]; then
    echo "🔽 OpenChat 3.5 Download:"
    echo "========================="
    echo "File needed: openchat-3.5-1210.Q4_K_M.gguf"
    echo "Size: ~4.1 GB"
    echo ""
    echo "Download options:"
    echo "1. Direct download (recommended):"
    echo "   https://huggingface.co/TheBloke/openchat_3.5-GGUF/resolve/main/openchat_3.5.q4_k_m.gguf"
    echo ""
    echo "2. Alternative sources:"
    echo "   https://huggingface.co/openchat/openchat_3.5/resolve/main/openchat_3.5.q4_k_m.gguf"
    echo ""
    echo "3. Browser download page:"
    echo "   https://huggingface.co/TheBloke/openchat_3.5-GGUF/tree/main"
    echo ""
    echo "📥 Manual download commands:"
    echo "cd $MODELS_DIR"
    echo "wget 'https://huggingface.co/TheBloke/openchat_3.5-GGUF/resolve/main/openchat_3.5.q4_k_m.gguf' -O openchat-3.5-1210.Q4_K_M.gguf"
    echo ""
    echo "Or with curl:"
    echo "curl -L 'https://huggingface.co/TheBloke/openchat_3.5-GGUF/resolve/main/openchat_3.5.q4_k_m.gguf' -o openchat-3.5-1210.Q4_K_M.gguf"
    echo ""
    echo "─────────────────────────────────────────────────────────────"
    echo ""
fi

if [ "$phi2_ok" = false ]; then
    echo "🔽 Microsoft Phi-2 Download:"
    echo "============================"
    echo "File needed: phi-2.Q4_K_M.gguf"
    echo "Size: ~1.6 GB"
    echo ""
    echo "Download options:"
    echo "1. Direct download (recommended):"
    echo "   https://huggingface.co/TheBloke/phi-2-GGUF/resolve/main/phi-2.q4_k_m.gguf"
    echo ""
    echo "2. Alternative sources:"
    echo "   https://huggingface.co/microsoft/phi-2-gguf/resolve/main/phi-2.q4_k_m.gguf"
    echo ""
    echo "3. Browser download page:"
    echo "   https://huggingface.co/TheBloke/phi-2-GGUF/tree/main"
    echo ""
    echo "📥 Manual download commands:"
    echo "cd $MODELS_DIR"
    echo "wget 'https://huggingface.co/TheBloke/phi-2-GGUF/resolve/main/phi-2.q4_k_m.gguf' -O phi-2.Q4_K_M.gguf"
    echo ""
    echo "Or with curl:"
    echo "curl -L 'https://huggingface.co/TheBloke/phi-2-GGUF/resolve/main/phi-2.q4_k_m.gguf' -o phi-2.Q4_K_M.gguf"
    echo ""
    echo "─────────────────────────────────────────────────────────────"
    echo ""
fi

if [ "$llama_ok" = false ]; then
    echo "🔽 Llama 2 Chat Download:"
    echo "========================="
    echo "File needed: llama-2-7b-chat.Q4_K_M.gguf"
    echo "Size: ~4.1 GB"
    echo ""
    echo "Download options:"
    echo "1. Direct download:"
    echo "   https://huggingface.co/TheBloke/Llama-2-7B-Chat-GGUF/resolve/main/llama-2-7b-chat.Q4_K_M.gguf"
    echo ""
    echo "2. Browser download page:"
    echo "   https://huggingface.co/TheBloke/Llama-2-7B-Chat-GGUF/tree/main"
    echo ""
    echo "📥 Manual download commands:"
    echo "cd $MODELS_DIR"
    echo "wget 'https://huggingface.co/TheBloke/Llama-2-7B-Chat-GGUF/resolve/main/llama-2-7b-chat.Q4_K_M.gguf'"
    echo ""
    echo "─────────────────────────────────────────────────────────────"
    echo ""
fi

if [ "$neural_ok" = false ]; then
    echo "🔽 Neural Chat Download:"
    echo "======================="
    echo "File needed: neural-chat-7b-v3-1.Q4_K_M.gguf"
    echo "Size: ~4.1 GB"
    echo ""
    echo "Download options:"
    echo "1. Direct download:"
    echo "   https://huggingface.co/TheBloke/neural-chat-7B-v3-1-GGUF/resolve/main/neural-chat-7b-v3-1.Q4_K_M.gguf"
    echo ""
    echo "2. Browser download page:"
    echo "   https://huggingface.co/TheBloke/neural-chat-7B-v3-1-GGUF/tree/main"
    echo ""
    echo "📥 Manual download commands:"
    echo "cd $MODELS_DIR"
    echo "wget 'https://huggingface.co/TheBloke/neural-chat-7B-v3-1-GGUF/resolve/main/neural-chat-7b-v3-1.Q4_K_M.gguf'"
    echo ""
    echo "─────────────────────────────────────────────────────────────"
    echo ""
fi

echo "💡 General Instructions:"
echo "========================"
echo ""
echo "1. Open your web browser"
echo "2. Copy and paste the download URLs"
echo "3. Save files to: $MODELS_DIR"
echo "4. Make sure filenames match exactly (case sensitive)"
echo "5. Restart your backend server after downloading"
echo ""

echo "🔧 Alternative: Download with WSL commands"
echo "=========================================="
echo ""
echo "Run these commands in WSL terminal:"
echo ""
echo "# Go to models directory"
echo "cd $MODELS_DIR"
echo ""

if [ "$openchat_ok" = false ]; then
    echo "# Download OpenChat 3.5"
    echo "wget 'https://huggingface.co/TheBloke/openchat_3.5-GGUF/resolve/main/openchat_3.5.q4_k_m.gguf' -O openchat-3.5-1210.Q4_K_M.gguf"
    echo ""
fi

if [ "$phi2_ok" = false ]; then
    echo "# Download Phi-2"
    echo "wget 'https://huggingface.co/TheBloke/phi-2-GGUF/resolve/main/phi-2.q4_k_m.gguf' -O phi-2.Q4_K_M.gguf"
    echo ""
fi

echo "🎯 After downloading:"
echo "===================="
echo "1. Check files: ls -la $MODELS_DIR/*.gguf"
echo "2. Run: ./scripts/check_models.sh"
echo "3. Start backend: go run cmd/server/main.go"
echo "4. Check frontend model list"
echo ""

# Create a simple download script
echo "📝 Creating download script..."
cat > "$MODELS_DIR/download_models.sh" << 'EOF'
#!/bin/bash
echo "🚀 Downloading missing models..."

# OpenChat 3.5
if [ ! -f "openchat-3.5-1210.Q4_K_M.gguf" ]; then
    echo "📥 Downloading OpenChat 3.5..."
    wget 'https://huggingface.co/TheBloke/openchat_3.5-GGUF/resolve/main/openchat_3.5.q4_k_m.gguf' -O openchat-3.5-1210.Q4_K_M.gguf
fi

# Phi-2
if [ ! -f "phi-2.Q4_K_M.gguf" ]; then
    echo "📥 Downloading Phi-2..."
    wget 'https://huggingface.co/TheBloke/phi-2-GGUF/resolve/main/phi-2.q4_k_m.gguf' -O phi-2.Q4_K_M.gguf
fi

echo "✅ Download complete!"
ls -la *.gguf
EOF

chmod +x "$MODELS_DIR/download_models.sh"
echo "✅ Created download script: $MODELS_DIR/download_models.sh"
echo ""

echo "⚡ Quick download option:"
echo "cd $MODELS_DIR && ./download_models.sh"
echo ""

if [ "$openchat_ok" = true ] && [ "$phi2_ok" = true ] && [ "$llama_ok" = true ] && [ "$neural_ok" = true ]; then
    echo "🎉 All models are already available!"
    echo "🚀 You can start your backend server now."
fi
