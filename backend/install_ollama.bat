@echo off
echo Installing Ollama for Windows...

:: Download Ollama installer
echo Downloading Ollama installer...
curl -fsSL https://ollama.com/install.sh | sh

if %ERRORLEVEL% == 0 (
    echo Ollama installed successfully!
    echo Starting Ollama service...
    ollama serve
) else (
    echo Failed to install Ollama. Please visit https://ollama.com/download for manual installation.
    echo.
    echo Manual installation steps:
    echo 1. Go to https://ollama.com/download
    echo 2. Download Ollama for Windows
    echo 3. Run the installer
    echo 4. Open a new command prompt and run: ollama serve
)

pause
