@echo off
echo Building Local AI Project Backend...
echo.

:: Check if Go is installed
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo ❌ Go is not installed or not in PATH
    echo Please install Go from https://golang.org/dl/
    pause
    exit /b 1
)

echo ✅ Go version:
go version
echo.

:: Set build environment
echo 🔧 Setting build environment...
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64

:: Create bin directory if it doesn't exist
if not exist "bin" mkdir bin

:: Clean previous builds
if exist "bin\server.exe" del "bin\server.exe"

echo 🔨 Building backend server...
echo   - Target: Windows x64
echo   - Output: bin\server.exe
echo   - Mode: Release build
echo.

:: Build with optimizations
go build -ldflags "-s -w -X main.version=1.0.0" -trimpath -o bin/server.exe ./cmd/server

if %ERRORLEVEL% == 0 (
    echo.
    echo ✅ Build successful!
    echo 📁 Executable: bin\server.exe
    
    :: Show file size
    for %%I in (bin\server.exe) do echo 📏 Size: %%~zI bytes
    
    echo.
    echo 🎯 Next steps:
    echo 1. Ensure Ollama is running: ollama serve
    echo 2. Check your models: ollama list
    echo 3. Start the server: bin\server.exe
    echo 4. Test API: http://localhost:8082/health
    echo.
    echo 💡 Your available models:
    echo    🚀 nemotron-nano    - NVIDIA Llama 3.1 Nemotron Nano 4B
    echo    🧠 neural-chat     - Intel Neural Chat 7B Q5_0  
    echo    💬 openchat        - OpenChat 3.5 Q5_K_M
    echo    🦙 llama2-chat     - Llama 2 7B Chat Q4_K_M
    echo    🔬 phi2            - Microsoft Phi-2 Q8_0
    echo.
) else (
    echo.
    echo ❌ Build failed!
    echo.
    echo 🔍 Troubleshooting:
    echo 1. Check Go modules: go mod tidy
    echo 2. Download dependencies: go mod download
    echo 3. Verify Go version: go version
    echo 4. Check for syntax errors: go vet ./...
    echo.
    pause
    exit /b 1
)

echo 🏁 Build process completed.
pause
