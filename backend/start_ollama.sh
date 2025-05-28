#!/bin/bash

echo "Checking Ollama service..."
echo ""

# Check if Ollama is already running
if pgrep -x "ollama" > /dev/null; then
    echo "✅ Ollama is already running on http://localhost:11434"
    echo ""
    
    # Test if Ollama is responding
    if curl -s http://localhost:11434/api/tags > /dev/null 2>&1; then
        echo "🌐 Ollama API is responding correctly"
        echo ""
        
        echo "Current models in Ollama:"
        ollama list
        echo ""
        
        echo "Your downloaded model files:"
        echo "- nvidia_Llama-3.1-Nemotron-Nano-4B-v1.1-bf16.gguf (Nemotron Nano)"
        echo "- neural-chat-7b-v3-1.Q5_0.gguf (Neural Chat)"
        echo "- openchat-3.5-0106.Q5_K_M.gguf (OpenChat 3.5)"
        echo "- llama-2-7b-chat.Q4_K_M.gguf (Llama 2 Chat)"
        echo "- phi-2.Q8_0.gguf (Phi-2)"
        echo ""
        
        echo "Options:"
        echo "1. Continue with current Ollama instance"
        echo "2. Restart Ollama service"
        echo "3. Stop Ollama service"
        echo "4. Setup your downloaded models"
        echo "5. Check model files status"
        echo ""
        read -p "Choose an option (1-5): " choice
        
        case $choice in
            1)
                echo "Using existing Ollama instance."
                ;;
            2)
                echo "Restarting Ollama..."
                pkill ollama
                sleep 2
                ollama serve &
                echo "Ollama restarted."
                ;;
            3)
                echo "Stopping Ollama..."
                pkill ollama
                echo "Ollama stopped."
                exit 0
                ;;
            4)
                echo "Setting up downloaded models..."
                ./setup_models.sh
                ;;
            5)
                echo "Checking model files status..."
                MODELS_DIR="/mnt/d/Praxis/KI/lokaleKI/go_mustAI/local-ai-project/models"
                echo "📁 Models directory: $MODELS_DIR"
                echo ""
                
                check_model_file() {
                    local filename="$1"
                    local display_name="$2"
                    
                    if [ -f "$MODELS_DIR/$filename" ]; then
                        size=$(du -h "$MODELS_DIR/$filename" 2>/dev/null | cut -f1 || echo "unknown")
                        echo "✅ $display_name: $filename ($size)"
                    else
                        echo "❌ $display_name: Missing ($filename)"
                    fi
                }
                
                check_model_file "nvidia_Llama-3.1-Nemotron-Nano-4B-v1.1-bf16.gguf" "Nemotron Nano"
                check_model_file "neural-chat-7b-v3-1.Q5_0.gguf" "Neural Chat"
                check_model_file "openchat-3.5-0106.Q5_K_M.gguf" "OpenChat 3.5"
                check_model_file "llama-2-7b-chat.Q4_K_M.gguf" "Llama 2 Chat"
                check_model_file "phi-2.Q8_0.gguf" "Phi-2"
                ;;
            *)
                echo "Invalid option. Using existing instance."
                ;;
        esac
    else
        echo "⚠️  Ollama process found but API not responding. Restarting..."
        pkill ollama
        sleep 2
        ollama serve
    fi
else
    echo "Starting Ollama service..."
    echo ""
    echo "Ollama will run on http://localhost:11434"
    echo ""
    echo "Your available model files:"
    echo "🚀 nvidia_Llama-3.1-Nemotron-Nano-4B-v1.1-bf16.gguf (Nemotron Nano - Fast & Efficient)"
    echo "🧠 neural-chat-7b-v3-1.Q5_0.gguf (Neural Chat - Intel Optimized)"
    echo "💬 openchat-3.5-0106.Q5_K_M.gguf (OpenChat 3.5 - High Quality)"
    echo "🦙 llama-2-7b-chat.Q4_K_M.gguf (Llama 2 Chat - Reliable)"
    echo "🔬 phi-2.Q8_0.gguf (Phi-2 - Compact & Powerful)"
    echo ""
    echo "💡 After Ollama starts, run ./setup_models.sh to configure these models"
    echo ""
    echo "Press Ctrl+C to stop"
    echo ""

    # Check if ollama is installed
    if ! command -v ollama &> /dev/null; then
        echo "Ollama not found. Installing..."
        curl -fsSL https://ollama.com/install.sh | sh
        echo "Ollama installed successfully!"
        echo ""
    fi

    # Start ollama service
    ollama serve
fi
