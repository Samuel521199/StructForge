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

REM 停止现有前端进程（确保干净启动）
echo [信息] 停止现有前端进程...
taskkill /FI "IMAGENAME eq node.exe" /FI "COMMANDLINE eq *vite*" /T /F >nul 2>&1
taskkill /FI "IMAGENAME eq node.exe" /FI "COMMANDLINE eq *npm*" /FI "COMMANDLINE eq *dev*" /T /F >nul 2>&1
timeout /t 2 /nobreak >nul 2>&1

REM 检查并清理前端端口
echo [信息] 检查前端端口 5173...
netstat -ano | findstr :5173 >nul 2>&1
if %errorlevel% equ 0 (
    echo [警告] 端口 5173 仍被占用，强制释放...
    for /f "tokens=5" %%a in ('netstat -ano ^| findstr :5173') do (
        taskkill /PID %%a /F >nul 2>&1
    )
    timeout /t 2 /nobreak >nul 2>&1
    echo [信息] 端口已清理
) else (
    echo [信息] 端口 5173 可用
)
echo.

echo [信息] 正在启动前端服务...
REM 使用快速启动脚本（减少检查步骤，加快启动速度）
if exist "%~dp0start-frontend-quick.bat" (
    start "StructForge - 前端服务" cmd /k "%~dp0start-frontend-quick.bat"
) else (
    start "StructForge - 前端服务" cmd /k "%~dp0start-frontend.bat"
)

REM 等待前端启动（增加等待时间，确保 Vite 完全启动）
echo [信息] 等待前端服务启动（最多等待 15 秒）...
timeout /t 5 /nobreak >nul 2>&1

REM 循环检查前端服务是否启动（最多等待 15 秒）
set CHECK_COUNT=0
:CHECK_FRONTEND
set /a CHECK_COUNT+=1
netstat -ano | findstr :5173 >nul 2>&1
if %errorlevel% equ 0 (
    echo [成功] 前端服务已启动，端口 5173 正在监听
    REM 再等待 2 秒确保 Vite 完全就绪
    timeout /t 2 /nobreak >nul 2>&1
    goto :FRONTEND_READY
)
if %CHECK_COUNT% geq 15 (
    echo [警告] 前端服务启动超时（已等待 15 秒），请检查前端服务窗口
    echo [提示] 如果前端无法访问，请等待几秒后刷新浏览器，或运行 restart-frontend.bat
    goto :FRONTEND_READY
)
timeout /t 1 /nobreak >nul 2>&1
goto :CHECK_FRONTEND

:FRONTEND_READY

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
echo [提示] 或运行 stop-all.bat 停止所有服务
echo.
echo [提示] 如果前端无法访问，请检查：
echo   1. 前端服务窗口是否显示 "VITE v5.x.x  ready"
echo   2. 端口 5173 是否被其他程序占用
echo   3. 防火墙是否阻止了端口访问
echo.
echo 按任意键退出此窗口（服务将继续运行）...
pause >nul
