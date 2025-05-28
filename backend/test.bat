@echo off
echo 🧪 Testing Local AI Backend...
echo.

:: Test if server is running
echo 🔍 Testing server health...
curl -s http://localhost:8082/health
if %ERRORLEVEL__ == 0 (
    echo ✅ Server is responding
) else (
    echo ❌ Server is not responding
    echo 💡 Make sure the server is running: start.bat
    pause
    exit /b 1
)

echo.
echo 🤖 Testing models endpoint...
curl -s http://localhost:8082/api/v1/models
echo.

echo.
echo 📄 Testing documents endpoint...
curl -s http://localhost:8082/api/v1/documents
echo.

echo.
echo ✅ API tests completed!
echo.
echo 🌐 You can also test in browser:
echo    http://localhost:8082/health
echo    http://localhost:8082/api/v1/models
echo.
pause
