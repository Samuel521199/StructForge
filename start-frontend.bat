@echo off
setlocal enabledelayedexpansion
chcp 65001 >nul 2>&1
echo ========================================
echo   Starting StructForge Frontend Service
echo ========================================
echo.

REM Switch to frontend directory
cd /d "%~dp0frontend"

REM Check Node.js environment
where node >nul 2>&1
if %errorlevel% neq 0 (
    echo [Error] Node.js not found, please install Node.js 18+
    echo Download: https://nodejs.org/
    pause
    exit /b 1
)

REM Check npm environment
where npm >nul 2>&1
if %errorlevel% neq 0 (
    echo [Error] npm not found, please install Node.js
    pause
    exit /b 1
)

REM Quick environment check (reduce output)
echo [Info] Working directory: %CD%
node --version >nul 2>&1
if %errorlevel% neq 0 (
    echo [Error] Node.js environment error
    pause
    exit /b 1
)
echo.

REM Check if package.json exists
if not exist "package.json" (
    echo [Error] package.json not found
    echo [Error] Please ensure package.json exists in frontend directory
    echo.
    pause
    exit /b 1
)

REM Check if node_modules exists
if not exist "node_modules" (
    echo [Info] node_modules not found, installing dependencies...
    echo.
    call npm install
    if %errorlevel% neq 0 (
        echo [Error] Dependency installation failed
        pause
        exit /b 1
    )
    echo.
    echo [Info] Dependencies installed
    echo.
)

REM Quick check for @vue-flow/core (check only, don't install)
echo [Info] Checking Vue Flow dependency...
npm list @vue-flow/core >nul 2>&1
if %errorlevel% neq 0 (
    echo [Warning] @vue-flow/core not found, installing in background...
    start /B npm install @vue-flow/core >nul 2>&1
) else (
    echo [Info] Vue Flow installed
)
echo.

echo [Info] Starting frontend development server...
echo [Info] Server will auto-open in browser
echo [Info] Default address: http://localhost:5173
echo.
echo ========================================
echo   Press Ctrl+C to stop
echo ========================================
echo.

REM Start frontend development server
echo [Info] Starting Vite development server...
echo.

REM Check if port is in use
netstat -ano | findstr :5173 >nul 2>&1
if %errorlevel% equ 0 (
    echo [Warning] Port 5173 is in use!
    echo [Info] Attempting to stop process using the port...
    for /f "tokens=5" %%a in ('netstat -ano ^| findstr :5173') do (
        taskkill /PID %%a /F >nul 2>&1
    )
    timeout /t 2 /nobreak >nul 2>&1
    REM Check port again
    netstat -ano | findstr :5173 >nul 2>&1
    if %errorlevel% equ 0 (
        echo [Error] Cannot free port 5173, please stop the process manually
        echo [Info] You can run stop-all.bat or restart-frontend.bat
        pause
        exit /b 1
    ) else (
        echo [Success] Port freed, continuing...
    )
    echo.
)

REM Start Vite (don't use call, so we can see real-time output)
npm run dev

REM Note: If we reach here, npm run dev has exited
set DEV_EXIT_CODE=!errorlevel!

REM If service exited, show error message
if !DEV_EXIT_CODE! neq 0 (
    echo.
    echo ----------------------------------------
    echo [Error] Frontend service failed (exit code: !DEV_EXIT_CODE!)
    echo ----------------------------------------
    echo.
    echo Please check:
    echo 1. Node.js environment is correctly installed
    echo 2. Dependencies are installed: npm install
    echo 3. package.json exists
    echo 4. vite is correctly installed
    echo 5. Port 5173 is not in use
    echo.
    echo [Info] You can try:
    echo   - Run restart-frontend.bat to restart
    echo   - Manually run: npm run dev
    echo   - Check browser console for errors
    echo.
    pause
) else (
    echo.
    echo [Info] Frontend service exited normally
    pause
)
