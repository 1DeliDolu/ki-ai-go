# Project Directory Navigation

cd /mnt/d/Praxis/KI/lokaleKI/go_mustAI/local-ai-project/frontend
cd /mnt/d/Praxis/KI/lokaleKI/go_mustAI/local-ai-project/backend
cd /mnt/d/Praxis/KI/lokaleKI/go_mustAI/local-ai-project/local-rag-system

# NPM Commands for Frontend Development

cd /mnt/d/Praxis/KI/lokaleKI/go_mustAI/local-ai-project/frontend

# Install dependencies

npm install

# Development server

npm run dev

# Build for production

npm run build

# Preview production build

npm run preview

# Install additional packages if needed

npm install @mui/material @emotion/react @emotion/styled
npm install @mui/icons-material
npm install axios
npm install react-router-dom
npm install typescript @types/react @types/react-dom

# GO Commands for Backend Development

cd /mnt/d/Praxis/KI/lokaleKI/go_mustAI/local-ai-project/backend

# Initialize Go module (if not already done)

go mod init local-ai-backend

# Install/update dependencies

go mod tidy

# Download specific Go packages

go get github.com/gin-gonic/gin
go get github.com/gorilla/websocket
go get github.com/rs/cors
go get gorm.io/gorm
go get gorm.io/driver/sqlite
go get github.com/joho/godotenv

# Run Go backend

go run main.go

# Build Go binary

go build -o local-ai-backend

# Run built binary

./local-ai-backend

# GO Commands for RAG System

cd /mnt/d/Praxis/KI/lokaleKI/go_mustAI/local-ai-project/local-rag-system

# Initialize Go module for RAG system

go mod init local-rag-system

# Install RAG-specific dependencies

go get github.com/tmc/langchaingo
go get github.com/chroma-core/chroma
go get github.com/pdfcpu/pdfcpu
go get github.com/ledongthuc/pdf
go get github.com/go-resty/resty/v2

# Run RAG system

go run main.go

# Build RAG system

go build -o local-rag-system

# Development Workflow Commands

# Terminal 1 - Frontend (React + Vite)

cd /mnt/d/Praxis/KI/lokaleKI/go_mustAI/local-ai-project/frontend && npm run dev

# Terminal 2 - Backend (Go)

cd /mnt/d/Praxis/KI/lokaleKI/go_mustAI/local-ai-project/backend && go run main.go

# Terminal 3 - RAG System (Go)

cd /mnt/d/Praxis/KI/lokaleKI/go_mustAI/local-ai-project/local-rag-system && go run main.go

# Terminal 4 - Model Downloads

cd /mnt/d/Praxis/KI/lokaleKI/go_mustAI/local-ai-project/models

# Package Management Commands

# Update all npm packages

cd /mnt/d/Praxis/KI/lokaleKI/go_mustAI/local-ai-project/frontend
npm update

# Check for outdated packages

npm outdated

# Update Go modules

cd /mnt/d/Praxis/KI/lokaleKI/go_mustAI/local-ai-project/backend
go get -u ./...
go mod tidy

cd /mnt/d/Praxis/KI/lokaleKI/go_mustAI/local-ai-project/local-rag-system
go get -u ./...
go mod tidy

# Model Download Commands

# Create models directory

mkdir -p /mnt/d/Praxis/KI/lokaleKI/go_mustAI/local-ai-project/models
cd /mnt/d/Praxis/KI/lokaleKI/go_mustAI/local-ai-project/models

# Download Llama 2 7B Chat model (4.1GB) - Choose one method:

# Method 1: Using wget

wget https://huggingface.co/TheBloke/Llama-2-7B-Chat-GGUF/resolve/main/llama-2-7b-chat.Q4_K_M.gguf

# Method 2: Using curl

curl -L -o llama-2-7b-chat.Q4_K_M.gguf https://huggingface.co/TheBloke/Llama-2-7B-Chat-GGUF/resolve/main/llama-2-7b-chat.Q4_K_M.gguf

# Method 3: Using curl with progress bar and resume capability

curl -L -C - --progress-bar -o llama-2-7b-chat.Q4_K_M.gguf https://huggingface.co/TheBloke/Llama-2-7B-Chat-GGUF/resolve/main/llama-2-7b-chat.Q4_K_M.gguf

# Alternative: Direct browser download link

# https://huggingface.co/TheBloke/Llama-2-7B-Chat-GGUF/resolve/main/llama-2-7b-chat.Q4_K_M.gguf

# Verify download

ls -lh llama-2-7b-chat.Q4_K_M.gguf

# Check file size (should be around 4.1GB)

du -h llama-2-7b-chat.Q4_K_M.gguf

# Other Popular Models to Download:

# 1. Mistral 7B Instruct (4.4GB) - Very good performance

wget https://huggingface.co/TheBloke/Mistral-7B-Instruct-v0.1-GGUF/resolve/main/mistral-7b-instruct-v0.1.Q4_K_M.gguf

# 2. Code Llama 7B Instruct (4.0GB) - Best for coding

wget https://huggingface.co/TheBloke/CodeLlama-7B-Instruct-GGUF/resolve/main/codellama-7b-instruct.Q4_K_M.gguf

# 3. Neural Chat 7B (4.1GB) - Intel optimized

wget https://huggingface.co/TheBloke/neural-chat-7B-v3-1-GGUF/resolve/main/neural-chat-7b-v3-1.Q4_K_M.gguf

# 4. Zephyr 7B Beta (4.1GB) - Great chat model

wget https://huggingface.co/TheBloke/zephyr-7B-beta-GGUF/resolve/main/zephyr-7b-beta.Q4_K_M.gguf

# 5. OpenChat 3.5 (4.1GB) - High quality conversations

wget https://huggingface.co/TheBloke/openchat-3.5-1210-GGUF/resolve/main/openchat-3.5-1210.Q4_K_M.gguf

# 6. Phi-2 (1.6GB) - Small but powerful Microsoft model

wget https://huggingface.co/TheBloke/phi-2-GGUF/resolve/main/phi-2.Q4_K_M.gguf

# 7. TinyLlama 1.1B (0.6GB) - Very fast, good for testing

wget https://huggingface.co/TheBloke/TinyLlama-1.1B-Chat-v1.0-GGUF/resolve/main/tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf

# Larger Models (if you have more resources):

# 8. Llama 2 13B Chat (7.4GB) - Better quality than 7B

wget https://huggingface.co/TheBloke/Llama-2-13B-Chat-GGUF/resolve/main/llama-2-13b-chat.Q4_K_M.gguf

# 9. Mixtral 8x7B Instruct (26GB) - Very powerful mixture of experts

# wget https://huggingface.co/TheBloke/Mixtral-8x7B-Instruct-v0.1-GGUF/resolve/main/mixtral-8x7b-instruct-v0.1.Q4_K_M.gguf

# Turkish Language Models:

# 10. Turkish Llama (4.1GB) - Turkish language support

wget https://huggingface.co/malhajar/llama2-turkish-7b-v1-GGUF/resolve/main/llama2-turkish-7b-v1.q4_k_m.gguf

# Specialized Models:

# 11. WizardCoder (4.1GB) - Excellent for programming

wget https://huggingface.co/TheBloke/WizardCoder-Python-7B-V1.0-GGUF/resolve/main/wizardcoder-python-7b-v1.0.Q4_K_M.gguf

# 12. Vicuna 7B (4.1GB) - Good general purpose

wget https://huggingface.co/TheBloke/vicuna-7B-v1.5-GGUF/resolve/main/vicuna-7b-v1.5.Q4_K_M.gguf

# Download with curl (alternative method):

# curl -L -C - --progress-bar -o [filename] [url]

# Check all downloaded models:

ls -lh *.gguf
du -h *.gguf

# OPEN SOURCE AI MODELS - All models below are completely free and open source:

# Meta's Open Source Models:

# 1. Llama 2 7B Chat (4.1GB) - Meta's open source model

wget https://huggingface.co/TheBloke/Llama-2-7B-Chat-GGUF/resolve/main/llama-2-7b-chat.Q4_K_M.gguf

# 2. Code Llama 7B (4.0GB) - Meta's open source coding model

wget https://huggingface.co/TheBloke/CodeLlama-7B-Instruct-GGUF/resolve/main/codellama-7b-instruct.Q4_K_M.gguf

# Mistral AI Open Source Models:

# 3. Mistral 7B Instruct (4.4GB) - Mistral AI's open source model

wget https://huggingface.co/TheBloke/Mistral-7B-Instruct-v0.1-GGUF/resolve/main/mistral-7b-instruct-v0.1.Q4_K_M.gguf

# 4. Mixtral 8x7B (26GB) - Mistral's open source mixture of experts

# wget https://huggingface.co/TheBloke/Mixtral-8x7B-Instruct-v0.1-GGUF/resolve/main/mixtral-8x7b-instruct-v0.1.Q4_K_M.gguf

# Microsoft Open Source Models:

# 5. Phi-2 (1.6GB) - Microsoft's small but powerful open source model

wget https://huggingface.co/TheBloke/phi-2-GGUF/resolve/main/phi-2.Q4_K_M.gguf

# 6. Phi-3 Mini (2.3GB) - Microsoft's latest small model

wget https://huggingface.co/microsoft/Phi-3-mini-4k-instruct-gguf/resolve/main/Phi-3-mini-4k-instruct-q4.gguf

# Intel Open Source Models:

# 7. Neural Chat 7B (4.1GB) - Intel's optimized open source model

wget https://huggingface.co/TheBloke/neural-chat-7B-v3-1-GGUF/resolve/main/neural-chat-7b-v3-1.Q4_K_M.gguf

# Community Open Source Models:

# 8. TinyLlama 1.1B (0.6GB) - Community developed, very fast

wget https://huggingface.co/TheBloke/TinyLlama-1.1B-Chat-v1.0-GGUF/resolve/main/tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf

# 9. Zephyr 7B Beta (4.1GB) - HuggingFace's open source chat model

wget https://huggingface.co/TheBloke/zephyr-7B-beta-GGUF/resolve/main/zephyr-7b-beta.Q4_K_M.gguf

# 10. OpenChat 3.5 (4.1GB) - Open source conversational AI

wget https://huggingface.co/TheBloke/openchat-3.5-1210-GGUF/resolve/main/openchat-3.5-1210.Q4_K_M.gguf

# 11. Vicuna 7B (4.1GB) - UC Berkeley's open source model

wget https://huggingface.co/TheBloke/vicuna-7B-v1.5-GGUF/resolve/main/vicuna-7b-v1.5.Q4_K_M.gguf

# Programming Specialized Open Source:

# 12. WizardCoder (4.1GB) - Open source coding assistant

wget https://huggingface.co/TheBloke/WizardCoder-Python-7B-V1.0-GGUF/resolve/main/wizardcoder-python-7b-v1.0.Q4_K_M.gguf

# 13. StarCoder (2.2GB) - BigCode's open source programming model

wget https://huggingface.co/TheBloke/starcoder-GGUF/resolve/main/starcoder.q4_k_m.gguf

# 14. DeepSeek Coder (4.1GB) - Open source coding model

wget https://huggingface.co/TheBloke/deepseek-coder-6.7B-instruct-GGUF/resolve/main/deepseek-coder-6.7b-instruct.Q4_K_M.gguf

# Multilingual Open Source Models:

# 15. Aya 8B (4.6GB) - Cohere's multilingual open source model

wget https://huggingface.co/CohereForAI/aya-23-8B-GGUF/resolve/main/aya-23-8b.Q4_K_M.gguf

# 16. Turkish Llama (4.1GB) - Open source Turkish language model

wget https://huggingface.co/malhajar/llama2-turkish-7b-v1-GGUF/resolve/main/llama2-turkish-7b-v1.q4_k_m.gguf

# Latest Open Source Models (2024):

# 17. Qwen2 7B (4.1GB) - Alibaba's open source model

wget https://huggingface.co/Qwen/Qwen2-7B-Instruct-GGUF/resolve/main/qwen2-7b-instruct-q4_k_m.gguf

# 18. Gemma 7B (4.2GB) - Google's open source model

wget https://huggingface.co/google/gemma-7b-it-gguf/resolve/main/gemma-7b-it.q4_k_m.gguf

# Open Source Model Verification:

# All these models have open source licenses:

# - Llama 2: Custom Open Source License (commercial use allowed)

# - Mistral: Apache 2.0 License

# - Phi: MIT License

# - Others: Various open source licenses (Apache, MIT, etc.)

# Check downloaded open source models:

ls -lh *.gguf
echo "All downloaded models are open source and free to use!"
