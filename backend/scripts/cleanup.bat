@echo off
echo 🧹 Manual Cleanup Script for Local AI Project
echo ==============================================

set API_URL=http://localhost:8082/api/v1
set UPLOADS_DIR=%USERPROFILE%\.local-ai-project\uploads

echo.
echo Available cleanup options:
echo 1. Clean all files (uploads + database via API)
echo 2. Clean documents only  
echo 3. Clean uploads directory manually
echo 4. Reset PostgreSQL database (truncate tables)
echo 5. Exit
echo.

set /p choice="Choose an option (1-5): "

if "%choice%"=="1" (
    echo 🗑️  Cleaning all files via API...
    curl -X POST "%API_URL%/cleanup/all" -H "Content-Type: application/json"
) else if "%choice%"=="2" (
    echo 📄 Cleaning documents only via API...
    curl -X POST "%API_URL%/cleanup/documents" -H "Content-Type: application/json"
) else if "%choice%"=="3" (
    echo 🗂️  Manually cleaning uploads directory...
    if exist "%UPLOADS_DIR%" (
        del /q "%UPLOADS_DIR%\*.*" 2>nul
        for /d %%x in ("%UPLOADS_DIR%\*") do rd /s /q "%%x" 2>nul
        echo ✅ Uploads directory cleaned: %UPLOADS_DIR%
    ) else (
        echo 📁 Uploads directory not found: %UPLOADS_DIR%
    )
) else if "%choice%"=="4" (
    echo 🗄️  Resetting PostgreSQL database...
    echo 💡 This will truncate all tables in the local_ai database
    curl -X POST "%API_URL%/cleanup/database" -H "Content-Type: application/json"
) else if "%choice%"=="5" (
    echo 👋 Exiting...
    exit /b 0
) else (
    echo ❌ Invalid option
    exit /b 1
)

echo.
echo ✅ Cleanup operation completed!
pause
