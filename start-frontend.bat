@echo off
setlocal enabledelayedexpansion
chcp 65001 >nul
echo ========================================
echo   启动 StructForge 前端服务
echo ========================================
echo.

REM 切换到前端目录
cd /d "%~dp0frontend"

REM 检查 Node.js 环境
where node >nul 2>&1
if %errorlevel% neq 0 (
    echo [错误] 未找到 Node.js 环境，请先安装 Node.js 18+
    echo 下载地址: https://nodejs.org/
    pause
    exit /b 1
)

REM 检查 npm 环境
where npm >nul 2>&1
if %errorlevel% neq 0 (
    echo [错误] 未找到 npm，请先安装 Node.js
    pause
    exit /b 1
)

echo [信息] 当前工作目录: %CD%
echo [信息] Node.js 版本:
node --version
echo [信息] npm 版本:
npm --version
echo.

REM 检查 package.json 是否存在
if not exist "package.json" (
    echo [错误] 未找到 package.json 文件
    echo [错误] 请确保在 frontend 目录下存在 package.json
    echo.
    pause
    exit /b 1
)

REM 检查 node_modules 是否存在
if not exist "node_modules" (
    echo [信息] 检测到 node_modules 不存在，正在安装依赖...
    echo.
    call npm install
    if %errorlevel% neq 0 (
        echo [错误] 依赖安装失败
        pause
        exit /b 1
    )
    echo.
    echo [信息] 依赖安装完成
    echo.
)

REM 检查是否安装了 @vue-flow/core
echo [信息] 检查 Vue Flow 依赖...
npm list @vue-flow/core >nul 2>&1
if %errorlevel% neq 0 (
    echo [警告] 未找到 @vue-flow/core，正在安装...
    call npm install @vue-flow/core
    if %errorlevel% neq 0 (
        echo [错误] Vue Flow 安装失败
        pause
        exit /b 1
    )
    echo [信息] Vue Flow 安装完成
    echo.
) else (
    echo [信息] Vue Flow 已安装
    echo.
)

echo [信息] 正在启动前端开发服务器...
echo [信息] 前端服务将在浏览器中自动打开
echo [信息] 默认地址: http://localhost:5173
echo.
echo ========================================
echo   按 Ctrl+C 停止服务
echo ========================================
echo.

REM 启动前端开发服务器
echo [调试] 准备执行: npm run dev
echo.
call npm run dev
set DEV_EXIT_CODE=!errorlevel!

REM 如果服务退出，显示错误信息
if !DEV_EXIT_CODE! neq 0 (
    echo.
    echo [错误] 前端服务启动失败 (退出代码: !DEV_EXIT_CODE!)
    echo 请检查:
    echo 1. Node.js 环境是否正确安装
    echo 2. 依赖是否已安装: npm install
    echo 3. package.json 是否存在
    echo 4. vite 是否正确安装
    echo.
    echo [提示] 可以尝试手动运行: npm run dev
    echo.
    pause
) else (
    echo.
    echo [信息] 前端服务已正常退出
    pause
)

