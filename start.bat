@echo off
setlocal enabledelayedexpansion

echo ========================================
echo   Starting Agent OS...
echo ========================================

cd /d "%~dp0"

set ROOT_DIR=%cd%
set BACKEND_DIR=%ROOT_DIR%
set FRONTEND_DIR=%ROOT_DIR%\dashboard

echo [1/3] Starting backend server...
start "Agent OS - Backend" cmd /k "cd /d %BACKEND_DIR% && go run cmd/server/main.go"

echo [2/3] Waiting for backend...
timeout /t 5 /nobreak >nul

echo [3/3] Starting frontend...
start "Agent OS - Frontend" cmd /k "cd /d %FRONTEND_DIR% && npm run dev"

echo.
echo ========================================
echo   Agent OS is running!
echo ========================================
echo   Backend:  http://localhost:8080
echo   Frontend: http://localhost:3000
echo.
echo   Press Ctrl+C in the windows to stop
echo ========================================
echo.

endlocal
