#!/bin/bash

echo "🧹 Manual Cleanup Script for Local AI Project"
echo "=============================================="

# Configuration
API_URL="http://localhost:8082/api/v1"
UPLOADS_DIR="$HOME/.local-ai-project/uploads"
DATABASE_PATH="$HOME/.local-ai-project/data/app.db"

echo ""
echo "Available cleanup options:"
echo "1. Clean all files (uploads + database)"
echo "2. Clean documents only"
echo "3. Clean uploads directory manually"
echo "4. Reset database manually"
echo "5. Exit"
echo ""

read -p "Choose an option (1-5): " choice

case $choice in
    1)
        echo "🗑️  Cleaning all files via API..."
        curl -X POST "$API_URL/cleanup/all" \
             -H "Content-Type: application/json" \
             -w "\nHTTP Status: %{http_code}\n"
        ;;
    2)
        echo "📄 Cleaning documents only via API..."
        curl -X POST "$API_URL/cleanup/documents" \
             -H "Content-Type: application/json" \
             -w "\nHTTP Status: %{http_code}\n"
        ;;
    3)
        echo "🗂️  Manually cleaning uploads directory..."
        if [ -d "$UPLOADS_DIR" ]; then
            rm -rf "$UPLOADS_DIR"/*
            echo "✅ Uploads directory cleaned: $UPLOADS_DIR"
        else
            echo "📁 Uploads directory not found: $UPLOADS_DIR"
        fi
        ;;
    4)
        echo "🗄️  Manually resetting database..."
        if [ -f "$DATABASE_PATH" ]; then
            rm -f "$DATABASE_PATH"
            echo "✅ Database file deleted: $DATABASE_PATH"
            echo "💡 Database will be recreated on next server start"
        else
            echo "📄 Database file not found: $DATABASE_PATH"
        fi
        ;;
    5)
        echo "👋 Exiting..."
        exit 0
        ;;
    *)
        echo "❌ Invalid option"
        exit 1
        ;;
esac

echo ""
echo "✅ Cleanup operation completed!"
