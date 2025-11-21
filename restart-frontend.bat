@echo off
setlocal enabledelayedexpansion
chcp 65001 >nul
echo ========================================
echo   重启 StructForge 前端服务
echo ========================================
echo.

REM 步骤 1: 停止所有前端相关进程
echo [步骤 1/3] 正在停止现有前端进程...
taskkill /FI "IMAGENAME eq node.exe" /FI "COMMANDLINE eq *vite*" /T /F >nul 2>&1
taskkill /FI "IMAGENAME eq node.exe" /FI "COMMANDLINE eq *npm*" /FI "COMMANDLINE eq *dev*" /T /F >nul 2>&1

REM 等待进程完全退出
timeout /t 2 /nobreak >nul 2>&1

REM 步骤 2: 检查端口是否释放
echo [步骤 2/3] 检查端口 5173...
netstat -ano | findstr :5173 >nul 2>&1
if %errorlevel% equ 0 (
    echo [警告] 端口 5173 仍被占用，尝试强制释放...
    for /f "tokens=5" %%a in ('netstat -ano ^| findstr :5173') do (
        echo [信息] 终止占用端口 5173 的进程: %%a
        taskkill /PID %%a /F >nul 2>&1
    )
    timeout /t 1 /nobreak >nul 2>&1
) else (
    echo [成功] 端口 5173 已释放
)
echo.

REM 步骤 3: 切换到前端目录并启动服务
cd /d "%~dp0frontend"

REM 检查必要文件
if not exist "package.json" (
    echo [错误] 未找到 package.json 文件
    pause
    exit /b 1
)

echo [步骤 3/3] 正在启动前端服务...
echo.
echo ========================================
echo   前端服务输出（按 Ctrl+C 停止）
echo ========================================
echo.

REM 启动 Vite 开发服务器
npm run dev

REM 如果执行到这里，说明服务已退出
echo.
echo [信息] 前端服务已退出
pause
