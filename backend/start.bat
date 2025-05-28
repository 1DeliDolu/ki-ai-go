@echo off
title Local AI Backend Server
echo 🚀 Starting Local AI Backend Server...
echo.

:: Check if Ollama is running
echo 🔍 Checking Ollama service...
curl -s http://localhost:11434/api/tags >nul 2>&1
if %ERRORLEVEL% == 0 (
    echo ✅ Ollama is running on http://localhost:11434
) else (
    echo ❌ Ollama is not running
    echo 💡 Please start Ollama first: ollama serve
    echo.
    pause
    exit /b 1
)

:: Check if server executable exists
if not exist "bin\server.exe" (
    echo ❌ Server executable not found
    echo 🔨 Building server first...
    call build.bat
    if %ERRORLEVEL% NEQ 0 (
        echo ❌ Build failed, cannot start server
        pause
        exit /b 1
    )
)

echo.
echo 🌐 Starting backend server...
echo    API: http://localhost:8082/api/v1
echo    Health: http://localhost:8082/health
echo    Models: http://localhost:8082/api/v1/models
echo.
echo 📋 Available models in Ollama:
ollama list 2>nul | findstr /C:"nemotron-nano" /C:"neural-chat" /C:"openchat" /C:"llama2-chat" /C:"phi2"
echo.
echo 💡 Press Ctrl+C to stop the server
echo.

:: Start the server
bin\server.exe

echo.
echo 🛑 Server stopped.
pause
