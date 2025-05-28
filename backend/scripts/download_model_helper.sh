#!/bin/bash

echo "🤖 Model Download Helper for Local AI Project"
echo "============================================="

cd "$(dirname "$0")/.."

echo "📍 Current location: $(pwd)"
echo "📁 Models will be saved to: ./models/"

echo ""
echo "📋 Available models:"
go run scripts/download_models.go

echo ""
echo "💡 Usage examples:"
echo "   go run scripts/download_models.go 1  # Download Llama 2 7B Chat"
echo "   go run scripts/download_models.go 2  # Download Mistral 7B"
echo "   go run scripts/download_models.go 3  # Download Phi-2"
echo "   go run scripts/download_models.go 4  # Download OpenChat 3.5"

echo ""
echo "⚠️  Note: Large models may take time to download"
echo "📊 Ensure you have enough disk space"
