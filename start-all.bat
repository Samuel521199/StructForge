@echo off
setlocal enabledelayedexpansion
chcp 65001 >nul
echo ========================================
echo   启动 StructForge 全栈服务
echo   前端 + 后端
echo ========================================
echo.

REM 检查环境
echo [信息] 检查环境...
echo.

REM 检查 Go 环境
where go >nul 2>&1
if !errorlevel! neq 0 (
    echo [错误] 未找到 Go 环境，请先安装 Go 1.21+
    pause
    exit /b 1
)

REM 检查 Node.js 环境
where node >nul 2>&1
if !errorlevel! neq 0 (
    echo [错误] 未找到 Node.js 环境，请先安装 Node.js 18+
    pause
    exit /b 1
)

echo [信息] 环境检查通过
echo.

REM 检查必要的批处理文件是否存在
if not exist "%~dp0start-backend.bat" (
    echo [错误] 未找到 start-backend.bat 文件
    pause
    exit /b 1
)

if not exist "%~dp0start-frontend.bat" (
    echo [错误] 未找到 start-frontend.bat 文件
    pause
    exit /b 1
)

echo [信息] 正在启动后端服务...
start "StructForge - 后端服务" cmd /k "%~dp0start-backend.bat"

REM 等待后端启动
timeout /t 3 /nobreak >nul

echo [信息] 正在启动前端服务...
start "StructForge - 前端服务" cmd /k "%~dp0start-frontend.bat"

echo.
echo ========================================
echo   服务启动完成！
echo ========================================
echo.
echo [信息] 后端服务: Gateway (新窗口)
echo [信息] 前端服务: http://localhost:5173 (新窗口)
echo.
echo [提示] 两个服务已在独立窗口中启动
echo [提示] 关闭对应窗口即可停止服务
echo.
echo 按任意键退出此窗口（服务将继续运行）...
pause >nul
