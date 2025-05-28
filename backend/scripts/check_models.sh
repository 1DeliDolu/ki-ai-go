#!/bin/bash

echo "🔍 Checking available model files..."
echo ""

# WSL Debian path - adjust this to your actual models directory
MODELS_DIR="/mnt/d/Praxis/KI/lokaleKI/go_mustAI/local-ai-project/models"

# Alternative paths to check if the main one doesn't exist
ALT_PATHS=(
    "$(pwd)/../models"
    "$(pwd)/../../models"
    "$HOME/.local-ai-project/models"
    "/tmp/models"
)

# Find the correct models directory
if [ ! -d "$MODELS_DIR" ]; then
    echo "⚠️  Primary models directory not found: $MODELS_DIR"
    echo "🔍 Searching for alternative locations..."
    
    for alt_path in "${ALT_PATHS[@]}"; do
        if [ -d "$alt_path" ]; then
            MODELS_DIR="$alt_path"
            echo "✅ Found models directory at: $MODELS_DIR"
            break
        fi
    done
    
    if [ ! -d "$MODELS_DIR" ]; then
        echo "❌ No models directory found in any of the expected locations:"
        for alt_path in "${ALT_PATHS[@]}"; do
            echo "   - $alt_path"
        done
        echo ""
        echo "💡 Please create the models directory and place your .gguf files there:"
        echo "   mkdir -p /mnt/d/Praxis/KI/lokaleKI/go_mustAI/local-ai-project/models"
        exit 1
    fi
fi

echo "📁 Models directory: $MODELS_DIR"
echo ""

echo "🗂️ All files in models directory:"
if ls -la "$MODELS_DIR"/ 2>/dev/null; then
    echo ""
else
    echo "❌ Cannot list directory contents or directory is empty"
    echo ""
fi

echo "🤖 Model files (*.gguf):"
gguf_files=($(find "$MODELS_DIR" -name "*.gguf" -type f 2>/dev/null))

if [ ${#gguf_files[@]} -eq 0 ]; then
    echo "❌ No .gguf files found in $MODELS_DIR"
    echo ""
    echo "💡 Expected model files:"
    echo "   - nvidia_Llama-3.1-Nemotron-Nano-4B-v1.1-bf16.gguf"
    echo "   - neural-chat-7b-v3-1.Q5_0.gguf"
    echo "   - openchat-3.5-0106.Q5_K_M.gguf"
    echo "   - llama-2-7b-chat.Q4_K_M.gguf"
    echo "   - phi-2.Q8_0.gguf"
else
    for file in "${gguf_files[@]}"; do
        filename=$(basename "$file")
        if command -v stat >/dev/null 2>&1; then
            # Try different stat commands for different systems
            if stat -c%s "$file" >/dev/null 2>&1; then
                size=$(stat -c%s "$file")
            elif stat -f%z "$file" >/dev/null 2>&1; then
                size=$(stat -f%z "$file")
            else
                size="unknown"
            fi
            
            if [[ "$size" != "unknown" ]]; then
                size_gb=$(echo "scale=1; $size / 1024 / 1024 / 1024" | bc 2>/dev/null || echo "?")
                echo "  📄 $filename (${size_gb} GB)"
            else
                echo "  📄 $filename (size unknown)"
            fi
        else
            echo "  📄 $filename"
        fi
    done
fi

echo ""
echo "🔗 Model mapping check:"
echo "Expected models and their files:"

echo "  🚀 Nemotron Nano  → nvidia_Llama-3.1-Nemotron-Nano-4B-v1.1-bf16.gguf"
echo "  🧠 Neural Chat    → neural-chat-7b-v3-1.Q5_0.gguf"  
echo "  💬 OpenChat 3.5   → openchat-3.5-0106.Q5_K_M.gguf"
echo "  🦙 Llama 2 Chat   → llama-2-7b-chat.Q4_K_M.gguf"
echo "  🔬 Phi-2          → phi-2.Q8_0.gguf"

echo ""
echo "📋 Availability status:"

check_model_size() {
    local model_file="$1"
    local model_name="$2"
    local alternatives="$3"
    
    if [ -f "$MODELS_DIR/$model_file" ]; then
        # Get file size
        if command -v stat >/dev/null 2>&1; then
            if stat -c%s "$MODELS_DIR/$model_file" >/dev/null 2>&1; then
                size=$(stat -c%s "$MODELS_DIR/$model_file")
            elif stat -f%z "$MODELS_DIR/$model_file" >/dev/null 2>&1; then
                size=$(stat -f%z "$MODELS_DIR/$model_file")
            else
                size="unknown"
            fi
            
            if [[ "$size" != "unknown" ]]; then
                size_gb=$(echo "scale=1; $size / 1024 / 1024 / 1024" | bc 2>/dev/null || echo "?")
                echo "  ✅ $model_name - Found: $model_file (${size_gb} GB)"
            else
                echo "  ✅ $model_name - Found: $model_file (size unknown)"
            fi
        else
            echo "  ✅ $model_name - Found: $model_file"
        fi
    else
        echo "  ❌ $model_name - Missing: $model_file"
        
        if [ -n "$alternatives" ]; then
            echo "     🔍 Checking alternatives:"
            IFS=',' read -ra ALT_ARRAY <<< "$alternatives"
            for alt in "${ALT_ARRAY[@]}"; do
                alt=$(echo "$alt" | xargs) # trim whitespace
                if [ -f "$MODELS_DIR/$alt" ]; then
                    echo "     ✅ Alternative found: $alt"
                else
                    echo "     ❌ Alternative missing: $alt"
                fi
            done
        fi
        
        # Check for pattern matches
        echo "     🔍 Pattern matching:"
        case $model_name in
            "Nemotron Nano")
                pattern_files=($(find "$MODELS_DIR" -name "*nemotron*" -o -name "*Nemotron*" -type f 2>/dev/null | head -3))
                if [ ${#pattern_files[@]} -gt 0 ]; then
                    for match in "${pattern_files[@]}"; do
                        echo "     🎯 Possible match: $(basename "$match")"
                    done
                else
                    echo "     ❌ No files matching *nemotron* pattern"
                fi
                ;;
            "Neural Chat")
                pattern_files=($(find "$MODELS_DIR" -name "*neural*chat*" -type f 2>/dev/null | head -3))
                if [ ${#pattern_files[@]} -gt 0 ]; then
                    for match in "${pattern_files[@]}"; do
                        echo "     🎯 Possible match: $(basename "$match")"
                    done
                else
                    echo "     ❌ No files matching *neural*chat* pattern"
                fi
                ;;
            "OpenChat 3.5")
                pattern_files=($(find "$MODELS_DIR" -name "*openchat*" -type f 2>/dev/null | head -3))
                if [ ${#pattern_files[@]} -gt 0 ]; then
                    for match in "${pattern_files[@]}"; do
                        echo "     🎯 Possible match: $(basename "$match")"
                    done
                else
                    echo "     ❌ No files matching *openchat* pattern"
                fi
                ;;
            "Llama 2 Chat")
                pattern_files=($(find "$MODELS_DIR" -name "*llama*2*chat*" -o -name "*llama-2*" -type f 2>/dev/null | head -3))
                if [ ${#pattern_files[@]} -gt 0 ]; then
                    for match in "${pattern_files[@]}"; do
                        echo "     🎯 Possible match: $(basename "$match")"
                    done
                else
                    echo "     ❌ No files matching *llama*2* pattern"
                fi
                ;;
            "Phi-2")
                pattern_files=($(find "$MODELS_DIR" -name "*phi*" -type f 2>/dev/null | head -3))
                if [ ${#pattern_files[@]} -gt 0 ]; then
                    for match in "${pattern_files[@]}"; do
                        echo "     🎯 Possible match: $(basename "$match")"
                    done
                else
                    echo "     ❌ No files matching *phi* pattern"
                fi
                ;;
        esac
    fi
}

check_model_size "nvidia_Llama-3.1-Nemotron-Nano-4B-v1.1-bf16.gguf" "Nemotron Nano"
check_model_size "neural-chat-7b-v3-1.Q5_0.gguf" "Neural Chat"
check_model_size "openchat-3.5-0106.Q5_K_M.gguf" "OpenChat 3.5"
check_model_size "llama-2-7b-chat.Q4_K_M.gguf" "Llama 2 Chat"
check_model_size "phi-2.Q8_0.gguf" "Phi-2"

echo ""
echo "💡 If models are missing, you can:"
echo "   1. Download them from Hugging Face"
echo "   2. Copy existing model files to the expected names"
echo "   3. Use the fix_model_names.sh script to rename existing files"
echo "   4. Check if files are in a different location"

echo ""
echo "🔄 To fix naming issues, run:"
echo "   chmod +x scripts/fix_model_names.sh"
echo "   ./scripts/fix_model_names.sh"
