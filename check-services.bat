@echo off
setlocal enabledelayedexpansion
chcp 65001 >nul
echo ========================================
echo   检查 StructForge 服务状态
echo ========================================
echo.

echo [检查] Gateway 服务 (端口 8000)...
netstat -ano | findstr ":8000" >nul 2>&1
if %errorlevel% equ 0 (
    echo [✓] Gateway 服务正在运行
    netstat -ano | findstr ":8000"
) else (
    echo [✗] Gateway 服务未运行
    echo [提示] 请运行 start-backend.bat 或手动启动 Gateway 服务
)
echo.

echo [检查] User 服务 HTTP (端口 8001)...
netstat -ano | findstr ":8001" >nul 2>&1
if %errorlevel% equ 0 (
    echo [✓] User 服务 HTTP 正在运行
    netstat -ano | findstr ":8001"
) else (
    echo [✗] User 服务 HTTP 未运行
)
echo.

echo [检查] User 服务 gRPC (端口 9001)...
netstat -ano | findstr ":9001" >nul 2>&1
if %errorlevel% equ 0 (
    echo [✓] User 服务 gRPC 正在运行
    netstat -ano | findstr ":9001"
) else (
    echo [✗] User 服务 gRPC 未运行
)
echo.

echo [检查] 前端服务 (端口 5173)...
netstat -ano | findstr ":5173" >nul 2>&1
if %errorlevel% equ 0 (
    echo [✓] 前端服务正在运行
    netstat -ano | findstr ":5173"
) else (
    echo [✗] 前端服务未运行
)
echo.

echo ========================================
echo   服务状态检查完成
echo ========================================
echo.
echo [提示] 如果 Gateway 服务未运行，请执行：
echo   1. 运行 start-backend.bat
echo   2. 或手动启动 Gateway 服务
echo.
pause

