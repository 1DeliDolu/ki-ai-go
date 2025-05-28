#!/bin/bash

echo "🧪 Testing Document Upload API (Local)"
echo "======================================"

API_URL="http://localhost:8082/api/v1"

# Function to test API with error handling
test_api() {
    local endpoint="$1"
    local description="$2"
    
    echo "Testing: $description"
    echo "URL: $API_URL$endpoint"
    
    response=$(curl -s -w "%{http_code}" "$API_URL$endpoint" 2>/dev/null)
    http_code="${response: -3}"
    content="${response%???}"
    
    echo "HTTP Code: $http_code"
    if [ "$http_code" = "200" ]; then
        echo "✅ Success"
        if command -v jq >/dev/null 2>&1; then
            echo "$content" | jq . 2>/dev/null || echo "$content"
        else
            echo "$content"
        fi
    else
        echo "❌ Failed (HTTP $http_code)"
        echo "Response: $content"
    fi
    echo "----------------------------------------"
}

# Check if server is running
echo "🔍 Testing server health..."
if curl -s http://localhost:8082/health >/dev/null 2>&1; then
    echo "✅ Server is responding"
else
    echo "❌ Server is not responding!"
    echo "💡 Make sure the server is running: ./start.sh"
    exit 1
fi

# Create test documents if they don't exist
if [ ! -d "test_documents" ]; then
    echo "📄 Creating test documents..."
    mkdir -p test_documents
    
    # Create simple test files
    echo "This is a test TXT file for document upload testing." > test_documents/test.txt
    echo '{"name": "test", "type": "json"}' > test_documents/test.json
    echo "name,value\ntest,123" > test_documents/test.csv
    echo "# Test Markdown\nThis is a **test** file." > test_documents/test.md
    echo "<html><body><h1>Test HTML</h1></body></html>" > test_documents/test.html
fi

echo ""
echo "📤 Testing document uploads..."

# Test uploads with better error handling
upload_file() {
    local file="$1"
    local description="$2"
    
    echo "$description"
    if [ -f "test_documents/$file" ]; then
        response=$(curl -s -w "%{http_code}" -X POST -F "file=@test_documents/$file" "$API_URL/documents/upload" 2>/dev/null)
        http_code="${response: -3}"
        content="${response%???}"
        
        echo "HTTP Code: $http_code"
        if [ "$http_code" = "200" ]; then
            echo "✅ Upload successful"
            if command -v jq >/dev/null 2>&1; then
                echo "$content" | jq '.document.name // .message' 2>/dev/null
            fi
        else
            echo "❌ Upload failed"
            echo "Response: $content"
        fi
    else
        echo "❌ File not found: test_documents/$file"
    fi
    echo "----------------------------------------"
}

# Test each file type
upload_file "test.txt" "1. Testing TXT file upload..."
upload_file "test.md" "2. Testing Markdown file upload..."
upload_file "test.html" "3. Testing HTML file upload..."
upload_file "test.json" "4. Testing JSON file upload..."
upload_file "test.csv" "5. Testing CSV file upload..."

echo ""
echo "📋 Testing API endpoints..."

# Test various endpoints
test_api "/health" "Health check"
test_api "/documents" "List all documents"
test_api "/documents/test" "List test documents"
test_api "/documents/types" "Get document types"

echo ""
echo "🔧 Alternative local test..."

# Direct file check
echo "📁 Checking uploaded files locally..."
if [ -d "$HOME/.local-ai-project/test_documents" ]; then
    echo "✅ Test documents directory exists:"
    ls -la "$HOME/.local-ai-project/test_documents/" 2>/dev/null || echo "Directory is empty"
else
    echo "❌ Test documents directory not found"
fi

echo ""
echo "📊 Server process check..."
if pgrep -f "server" >/dev/null; then
    echo "✅ Server process is running"
    echo "Process info:"
    pgrep -fl "server" | head -3
else
    echo "❌ No server process found"
fi

echo ""
echo "🌐 Direct localhost test..."
echo "Testing direct connection:"
curl -s "http://127.0.0.1:8082/health" | head -20

echo ""
echo "✅ Local document upload tests completed!"
echo ""
echo "💡 If uploads work but API responses are blocked:"
echo "   - Uploads are successful (files are being processed)"
echo "   - Network firewall is blocking API responses"
echo "   - Check your local test_documents directory"
echo "   - Use browser to test: http://localhost:8082/health"
