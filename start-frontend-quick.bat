@echo off
setlocal enabledelayedexpansion
chcp 65001 >nul 2>&1

REM Quick start script for start-all.bat
REM Minimal checks, fast startup

REM Switch to frontend directory
cd /d "%~dp0frontend"

REM Quick check: package.json exists
if not exist "package.json" (
    echo [Error] package.json not found
    pause
    exit /b 1
)

REM Check if node_modules exists (only if missing)
if not exist "node_modules" (
    echo [Info] Installing dependencies...
    call npm install
    if %errorlevel% neq 0 (
        echo [Error] Dependency installation failed
        pause
        exit /b 1
    )
)

REM Start Vite immediately (no other checks)
npm run dev

