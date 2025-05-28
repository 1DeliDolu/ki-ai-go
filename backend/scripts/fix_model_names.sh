#!/bin/bash

echo "🔧 Fixing model filenames for WSL Debian..."
echo ""

# WSL Debian path - adjust this to your actual models directory
MODELS_DIR="/mnt/d/Praxis/KI/lokaleKI/go_mustAI/local-ai-project/models"

# Alternative paths to check
ALT_PATHS=(
    "$(pwd)/../models"
    "$(pwd)/../../models"
    "$HOME/.local-ai-project/models"
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
        echo "❌ No models directory found. Creating it..."
        mkdir -p "/mnt/d/Praxis/KI/lokaleKI/go_mustAI/local-ai-project/models"
        MODELS_DIR="/mnt/d/Praxis/KI/lokaleKI/go_mustAI/local-ai-project/models"
        echo "✅ Created: $MODELS_DIR"
    fi
fi

cd "$MODELS_DIR" || exit 1

echo "📁 Working in: $MODELS_DIR"
echo ""

# Function to rename file if it exists
rename_if_exists() {
    local pattern="$1"
    local target="$2"
    local description="$3"
    
    echo "🔍 Looking for $description..."
    
    # Find files matching the pattern (case insensitive)
    found_files=($(find . -maxdepth 1 -iname "*${pattern}*" -type f 2>/dev/null | head -5))
    
    if [ ${#found_files[@]} -gt 0 ]; then
        echo "  📄 Found ${#found_files[@]} potential matches:"
        for i in "${!found_files[@]}"; do
            echo "    $((i+1)). $(basename "${found_files[$i]}")"
        done
        
        # Use the first match if target doesn't exist
        source_file=$(basename "${found_files[0]}")
        
        if [ "$source_file" != "$target" ]; then
            echo "  ➡️  Selected: $source_file"
            echo "  ➡️  Renaming to: $target"
            
            if [ -f "$target" ]; then
                echo "  ⚠️  Target file already exists, creating backup..."
                mv "$target" "${target}.backup.$(date +%s)"
            fi
            
            if mv "$source_file" "$target" 2>/dev/null; then
                echo "  ✅ Renamed successfully!"
            else
                echo "  ❌ Failed to rename file (check permissions)"
            fi
        else
            echo "  ✅ Already correctly named: $target"
        fi
    else
        echo "  ❌ No files found matching pattern: *${pattern}*"
        
        # Show all .gguf files for reference
        all_gguf=($(find . -maxdepth 1 -name "*.gguf" -type f 2>/dev/null))
        if [ ${#all_gguf[@]} -gt 0 ]; then
            echo "  📁 Available .gguf files:"
            for file in "${all_gguf[@]}"; do
                echo "    - $(basename "$file")"
            done
        fi
    fi
    echo ""
}

# Create symlinks as an alternative approach
create_symlink() {
    local pattern="$1"
    local target="$2"
    local description="$3"
    
    if [ ! -f "$target" ]; then
        found_files=($(find . -maxdepth 1 -iname "*${pattern}*" -type f 2>/dev/null | head -1))
        
        if [ ${#found_files[@]} -gt 0 ]; then
            source_file=$(basename "${found_files[0]}")
            echo "🔗 Creating symlink: $target → $source_file"
            if ln -sf "$source_file" "$target" 2>/dev/null; then
                echo "  ✅ Symlink created successfully!"
            else
                echo "  ❌ Failed to create symlink"
            fi
        fi
    fi
}

echo "🔄 Method 1: Renaming files to expected names"
echo "============================================="

# Rename openchat variants
rename_if_exists "openchat" "openchat-3.5-1210.Q4_K_M.gguf" "OpenChat model"

# Rename phi variants  
rename_if_exists "phi" "phi-2.Q4_K_M.gguf" "Phi-2 model"

echo "🔗 Method 2: Creating symlinks for missing files"
echo "================================================"

create_symlink "openchat" "openchat-3.5-1210.Q4_K_M.gguf" "OpenChat model"
create_symlink "phi" "phi-2.Q4_K_M.gguf" "Phi-2 model"

echo "🎯 Final status check:"
echo "======================"

if [ -f "llama-2-7b-chat.Q4_K_M.gguf" ]; then
    echo "✅ llama2-chat: llama-2-7b-chat.Q4_K_M.gguf"
else
    echo "❌ llama2-chat: Missing"
fi

if [ -f "neural-chat-7b-v3-1.Q4_K_M.gguf" ]; then
    echo "✅ neural-chat: neural-chat-7b-v3-1.Q4_K_M.gguf"
else
    echo "❌ neural-chat: Missing"
fi

if [ -f "openchat-3.5-1210.Q4_K_M.gguf" ]; then
    echo "✅ openchat: openchat-3.5-1210.Q4_K_M.gguf"
else
    echo "❌ openchat: Missing"
fi

if [ -f "phi-2.Q4_K_M.gguf" ]; then
    echo "✅ phi2: phi-2.Q4_K_M.gguf"
else
    echo "❌ phi2: Missing"
fi

echo ""
echo "🔄 Next steps:"
echo "1. Restart your backend server to detect the changes"
echo "2. Check the backend logs for model availability"
echo "3. If issues persist, run: ./scripts/check_models.sh"

echo ""
echo "💡 Manual download commands if files are missing:"
echo "================================="
echo "# For OpenChat:"
echo "wget https://huggingface.co/TheBloke/openchat_3.5-GGUF/resolve/main/openchat_3.5.q4_k_m.gguf -O openchat-3.5-1210.Q4_K_M.gguf"
echo ""
echo "# For Phi-2:"
echo "wget https://huggingface.co/microsoft/phi-2/resolve/main/pytorch_model.bin -O phi-2.Q4_K_M.gguf"
