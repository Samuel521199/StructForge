@echo off
setlocal enabledelayedexpansion
chcp 65001 >nul
echo ========================================
echo   检查 StructForge 服务状态
echo ========================================
echo.

REM 检查前端服务（端口 5173）
echo [检查] 前端服务 (端口 5173)...
netstat -ano | findstr :5173 >nul 2>&1
if !errorlevel! equ 0 (
    echo [状态] 前端服务正在运行
    netstat -ano | findstr :5173
) else (
    echo [状态] 前端服务未运行
    echo [提示] 请运行 start-frontend.bat 启动前端服务
)
echo.

REM 检查 Gateway 服务（端口 8000）
echo [检查] Gateway 服务 (端口 8000)...
netstat -ano | findstr :8000 >nul 2>&1
if !errorlevel! equ 0 (
    echo [状态] Gateway 服务正在运行
    netstat -ano | findstr :8000
) else (
    echo [状态] Gateway 服务未运行
    echo [提示] 请运行 start-backend.bat 启动后端服务
)
echo.

REM 检查 User 服务（端口 8001）
echo [检查] User 服务 (端口 8001)...
netstat -ano | findstr :8001 >nul 2>&1
if !errorlevel! equ 0 (
    echo [状态] User 服务正在运行
    netstat -ano | findstr :8001
) else (
    echo [状态] User 服务未运行
    echo [提示] 请运行 start-backend.bat 启动后端服务
)
echo.

REM 检查 User gRPC 服务（端口 9001）
echo [检查] User gRPC 服务 (端口 9001)...
netstat -ano | findstr :9001 >nul 2>&1
if !errorlevel! equ 0 (
    echo [状态] User gRPC 服务正在运行
    netstat -ano | findstr :9001
) else (
    echo [状态] User gRPC 服务未运行
)
echo.

echo ========================================
echo   服务状态检查完成
echo ========================================
echo.
echo [提示] 如果服务未运行，请使用以下命令启动：
echo   1. 启动所有服务: start-all.bat
echo   2. 仅启动前端: start-frontend.bat
echo   3. 仅启动后端: start-backend.bat
echo.
pause
