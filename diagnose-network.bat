@echo off
setlocal enabledelayedexpansion
chcp 65001 >nul
echo ========================================
echo   网络连接诊断工具
echo ========================================
echo.

REM 检查前端服务
echo [1] 检查前端服务 (端口 5173)...
netstat -ano | findstr :5173 >nul 2>&1
if !errorlevel! equ 0 (
    echo     [✓] 前端服务正在运行
) else (
    echo     [✗] 前端服务未运行
    echo     [提示] 运行 start-frontend.bat 启动前端服务
)
echo.

REM 检查 Gateway 服务
echo [2] 检查 Gateway 服务 (端口 8000)...
netstat -ano | findstr :8000 >nul 2>&1
if !errorlevel! equ 0 (
    echo     [✓] Gateway 服务正在运行
) else (
    echo     [✗] Gateway 服务未运行
    echo     [提示] 运行 start-backend.bat 启动后端服务
)
echo.

REM 检查 User 服务
echo [3] 检查 User 服务 (端口 8001)...
netstat -ano | findstr :8001 >nul 2>&1
if !errorlevel! equ 0 (
    echo     [✓] User 服务正在运行
) else (
    echo     [✗] User 服务未运行
    echo     [提示] 运行 start-backend.bat 启动后端服务
)
echo.

REM 测试 Gateway 连接
echo [4] 测试 Gateway 连接...
curl -s -o nul -w "HTTP状态码: %%{http_code}\n" http://localhost:8000/health 2>nul
if !errorlevel! equ 0 (
    echo     [✓] Gateway 健康检查通过
) else (
    echo     [✗] 无法连接到 Gateway
    echo     [提示] 请检查 Gateway 服务是否正常运行
)
echo.

echo ========================================
echo   诊断完成
echo ========================================
echo.
echo [解决方案]
echo 1. 如果前端服务未运行: start-frontend.bat
echo 2. 如果后端服务未运行: start-backend.bat
echo 3. 如果所有服务未运行: start-all.bat
echo 4. 如果服务运行但无法连接，请检查防火墙设置
echo.
pause

