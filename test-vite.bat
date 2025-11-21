@echo off
setlocal enabledelayedexpansion
chcp 65001 >nul
echo ========================================
echo   测试 Vite 开发服务器
echo ========================================
echo.

cd /d "%~dp0frontend"

echo [1] 检查当前目录...
echo     目录: %CD%
echo.

echo [2] 检查 package.json...
if exist "package.json" (
    echo     [✓] package.json 存在
) else (
    echo     [✗] package.json 不存在
    pause
    exit /b 1
)
echo.

echo [3] 检查 node_modules...
if exist "node_modules" (
    echo     [✓] node_modules 存在
) else (
    echo     [✗] node_modules 不存在，需要运行 npm install
    pause
    exit /b 1
)
echo.

echo [4] 检查 vite 是否安装...
npm list vite >nul 2>&1
if !errorlevel! equ 0 (
    echo     [✓] vite 已安装
    npm list vite
) else (
    echo     [✗] vite 未安装
    echo     [提示] 运行: npm install
    pause
    exit /b 1
)
echo.

echo [5] 检查端口 5173 是否被占用...
netstat -ano | findstr :5173 >nul 2>&1
if !errorlevel! equ 0 (
    echo     [警告] 端口 5173 已被占用
    netstat -ano | findstr :5173
    echo.
    echo     [提示] 请先停止占用该端口的进程，或修改 vite.config.ts 中的端口号
) else (
    echo     [✓] 端口 5173 可用
)
echo.

echo [6] 尝试启动 Vite（测试模式）...
echo     执行: npm run dev
echo.
echo ========================================
echo   如果看到以下输出，说明启动成功:
echo   - VITE v5.x.x  ready in xxx ms
echo   - ➜  Local:   http://localhost:5173/
echo ========================================
echo.
echo [提示] 按 Ctrl+C 停止服务器
echo.

REM 直接运行，不使用 call，确保能看到输出
npm run dev

pause

