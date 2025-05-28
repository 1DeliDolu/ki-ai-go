#!/bin/bash

echo "Setting up your downloaded models for Ollama..."

MODELS_DIR="/mnt/d/Praxis/KI/lokaleKI/go_mustAI/local-ai-project/models"

# Check if models directory exists
if [ ! -d "$MODELS_DIR" ]; then
    echo "Models directory not found: $MODELS_DIR"
    exit 1
fi

echo "Models directory: $MODELS_DIR"
echo "Available models:"
ls -la "$MODELS_DIR"/*.gguf 2>/dev/null || echo "No .gguf files found"

echo ""
echo "Creating Ollama modelfiles for your downloaded models..."

# Function to create Ollama model with proper templates
create_ollama_model() {
    local model_file="$1"
    local model_name="$2"
    local template="$3"
    local system_prompt="$4"
    
    if [ -f "$MODELS_DIR/$model_file" ]; then
        echo "Creating Ollama model: $model_name"
        echo "  Source file: $model_file"
        
        # Create modelfile with proper template and system prompt
        cat > "/tmp/${model_name}.modelfile" << EOF
FROM $MODELS_DIR/$model_file

$template

SYSTEM """$system_prompt"""

PARAMETER temperature 0.7
PARAMETER top_p 0.9
PARAMETER top_k 40
PARAMETER repeat_penalty 1.1
PARAMETER num_predict 512
PARAMETER stop "<|end|>"
PARAMETER stop "<|eot_id|>"
PARAMETER stop "</s>"
EOF

        # Create the model in Ollama
        echo "  Creating model in Ollama..."
        if ollama create "$model_name" -f "/tmp/${model_name}.modelfile" 2>/dev/null; then
            echo "✅ Successfully created: $model_name"
        else
            echo "❌ Failed to create: $model_name"
            echo "  Retrying with simpler configuration..."
            
            # Fallback with simpler template
            cat > "/tmp/${model_name}.modelfile" << EOF
FROM $MODELS_DIR/$model_file

TEMPLATE """{{ .Prompt }}"""

PARAMETER temperature 0.7
PARAMETER top_p 0.9
PARAMETER num_predict 512
EOF
            if ollama create "$model_name" -f "/tmp/${model_name}.modelfile" 2>/dev/null; then
                echo "✅ Created with fallback template: $model_name"
            else
                echo "❌ Still failed to create: $model_name"
            fi
        fi
        
        # Clean up
        rm -f "/tmp/${model_name}.modelfile"
    else
        echo "⚠️  Model file not found: $model_file"
    fi
}

# Templates for different model types
LLAMA3_TEMPLATE='TEMPLATE """<|begin_of_text|><|start_header_id|>system<|end_header_id|>

{{ .System }}<|eot_id|><|start_header_id|>user<|end_header_id|>

{{ .Prompt }}<|eot_id|><|start_header_id|>assistant<|end_header_id|>

{{ .Response }}<|eot_id|>"""'

LLAMA2_TEMPLATE='TEMPLATE """<s>[INST] <<SYS>>
{{ .System }}
<</SYS>>

{{ .Prompt }} [/INST] {{ .Response }} </s>"""'

NEURAL_CHAT_TEMPLATE='TEMPLATE """### System:
{{ .System }}

### User:
{{ .Prompt }}

### Assistant:
{{ .Response }}"""'

OPENCHAT_TEMPLATE='TEMPLATE """GPT4 Correct User: {{ .Prompt }}<|end_of_turn|>GPT4 Correct Assistant: {{ .Response }}<|end_of_turn|>"""'

PHI_TEMPLATE='TEMPLATE """<|system|>
{{ .System }}<|end|>
<|user|>
{{ .Prompt }}<|end|>
<|assistant|>
{{ .Response }}<|end|>"""'

# System prompts for different models
NEMOTRON_SYSTEM="You are a helpful AI assistant created by NVIDIA. You are knowledgeable, accurate, and concise in your responses."

NEURAL_SYSTEM="You are Neural Chat, an AI assistant optimized by Intel. You provide helpful, harmless, and honest responses."

OPENCHAT_SYSTEM="You are OpenChat, a helpful AI assistant. You follow instructions carefully and provide accurate information."

LLAMA_SYSTEM="You are Llama, a helpful AI assistant created by Meta. You are honest, helpful, and harmless."

PHI_SYSTEM="You are Phi, a helpful AI assistant by Microsoft. You provide concise and accurate responses."

echo "Setting up your specific model files..."
echo "======================================"

# Create models for your actual downloaded files
create_ollama_model "nvidia_Llama-3.1-Nemotron-Nano-4B-v1.1-bf16.gguf" "nemotron-nano" "$LLAMA3_TEMPLATE" "$NEMOTRON_SYSTEM"

create_ollama_model "neural-chat-7b-v3-1.Q5_0.gguf" "neural-chat" "$NEURAL_CHAT_TEMPLATE" "$NEURAL_SYSTEM"

create_ollama_model "openchat-3.5-0106.Q5_K_M.gguf" "openchat" "$OPENCHAT_TEMPLATE" "$OPENCHAT_SYSTEM"

create_ollama_model "llama-2-7b-chat.Q4_K_M.gguf" "llama2-chat" "$LLAMA2_TEMPLATE" "$LLAMA_SYSTEM"

create_ollama_model "phi-2.Q8_0.gguf" "phi2" "$PHI_TEMPLATE" "$PHI_SYSTEM"

echo ""
echo "Testing created models..."
echo "========================"

echo "Available models in Ollama:"
if command -v ollama >/dev/null 2>&1; then
    ollama list
else
    echo "❌ Ollama command not found. Please ensure Ollama is installed and running."
fi

echo ""
echo "✅ Setup complete! Your models are ready:"
echo "========================================"
echo "🚀 nemotron-nano  - NVIDIA Llama 3.1 Nemotron Nano 4B (High quality, fast)"
echo "🧠 neural-chat   - Intel Neural Chat 7B Q5_0 (Optimized performance)"
echo "💬 openchat      - OpenChat 3.5 Q5_K_M (High quality conversations)"
echo "🦙 llama2-chat   - Llama 2 7B Chat Q4_K_M (General purpose)"
echo "🔬 phi2          - Microsoft Phi-2 Q8_0 (Compact but powerful)"

echo ""
echo "📋 Model Information:"
echo "===================="

check_model_size() {
    local filename="$1"
    local model_name="$2"
    
    if [ -f "$MODELS_DIR/$filename" ]; then
        size=$(du -h "$MODELS_DIR/$filename" 2>/dev/null | cut -f1 || echo "unknown")
        echo "  📄 $model_name: $size"
    fi
}

check_model_size "nvidia_Llama-3.1-Nemotron-Nano-4B-v1.1-bf16.gguf" "Nemotron Nano"
check_model_size "neural-chat-7b-v3-1.Q5_0.gguf" "Neural Chat"
check_model_size "openchat-3.5-0106.Q5_K_M.gguf" "OpenChat 3.5"
check_model_size "llama-2-7b-chat.Q4_K_M.gguf" "Llama 2 Chat"
check_model_size "phi-2.Q8_0.gguf" "Phi-2"

echo ""
echo "🔧 Testing model functionality..."
echo "================================="

test_model() {
    local model_name="$1"
    echo "Testing $model_name..."
    
    if timeout 30s ollama run "$model_name" "Hello, please respond with just 'OK' to confirm you're working." 2>/dev/null | grep -q "OK\|Hi\|Hello"; then
        echo "✅ $model_name is working"
    else
        echo "⚠️  $model_name may need configuration adjustment"
    fi
}

if command -v ollama >/dev/null 2>&1; then
    test_model "nemotron-nano"
    test_model "neural-chat"
    test_model "openchat"
    test_model "llama2-chat"
    test_model "phi2"
fi

echo ""
echo "🎯 Next Steps:"
echo "============="
echo "1. Start your backend server: cd ../backend && go run cmd/server/main.go"
echo "2. Open your frontend: http://localhost:5174"
echo "3. Select one of these models in the Models tab"
echo "4. Start chatting!"

echo ""
echo "💡 Model Recommendations:"
echo "========================="
echo "• For fastest responses: nemotron-nano (4B parameters)"
echo "• For best quality: neural-chat or openchat (7B parameters)"
echo "• For general use: llama2-chat (proven and reliable)"
echo "• For minimal resources: phi2 (compact but capable)"

echo ""
echo "🚀 You can now select these models in your frontend!"
