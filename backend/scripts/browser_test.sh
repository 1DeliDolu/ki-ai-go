#!/bin/bash

echo "🌐 Browser Test Helper"
echo "====================="

echo "📋 Available API endpoints for browser testing:"
echo ""
echo "🔹 Health Check:"
echo "   http://localhost:8082/health"
echo ""
echo "🔹 API v1 Health:"
echo "   http://localhost:8082/api/v1/health"
echo ""
echo "🔹 List Documents:"
echo "   http://localhost:8082/api/v1/documents"
echo ""
echo "🔹 Test Documents:"
echo "   http://localhost:8082/api/v1/documents/test"
echo ""
echo "🔹 Document Types:"
echo "   http://localhost:8082/api/v1/documents/types"
echo ""
echo "🔹 Models:"
echo "   http://localhost:8082/api/v1/models"
echo ""

echo "💡 Copy and paste these URLs into your browser to test"
echo "💡 If browser works but curl doesn't, it's a firewall issue"

# Try to open browser automatically
if command -v xdg-open >/dev/null 2>&1; then
    echo ""
    echo "🚀 Opening browser automatically..."
    xdg-open "http://localhost:8082/health" &
elif command -v open >/dev/null 2>&1; then
    echo ""
    echo "🚀 Opening browser automatically..."
    open "http://localhost:8082/health" &
else
    echo ""
    echo "📱 Please open your browser manually and test the URLs above"
fi

echo ""
echo "✅ Browser test helper ready!"
