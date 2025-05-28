@echo off
echo 🧪 Testing Local AI Backend...
echo.

echo 🔍 Testing server health...
curl -s http://localhost:8082/health >nul
if %ERRORLEVEL% == 0 (
    echo ✅ Server is responding
    echo Response:
    curl -s http://localhost:8082/health
) else (
    echo ❌ Server is not responding
    echo 💡 Make sure the server is running: start.bat
    exit /b 1
)

echo.
echo 🤖 Testing models endpoint...
echo Response:
curl -s http://localhost:8082/api/v1/models

echo.
echo 📄 Testing documents endpoint...
echo Response:
curl -s http://localhost:8082/api/v1/documents

echo.
echo 📋 Testing document types endpoint...
echo Response:
curl -s http://localhost:8082/api/v1/documents/types

echo.
echo 📁 Testing test documents endpoint...
curl -s http://localhost:8082/api/v1/documents/test

echo.
echo 🧹 Testing test documents cleanup...
curl -s -X POST http://localhost:8082/api/v1/cleanup/test-documents

echo.
echo ✅ API tests completed!
echo.
echo 💡 Frontend document endpoints:
echo     GET /api/v1/documents/test - List test documents
echo     POST /api/v1/cleanup/test-documents - Clean test documents
echo.
echo 🌐 You can also test in browser:
echo     http://localhost:8082/health
echo     http://localhost:8082/api/v1/models
echo     http://localhost:8082/api/v1/documents/types

pause
