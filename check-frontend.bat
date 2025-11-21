@echo off
setlocal enabledelayedexpansion
chcp 65001 >nul
echo ========================================
echo   检查前端服务状态
echo ========================================
echo.

REM 检查端口 5173
echo [检查] 端口 5173 占用情况:
netstat -ano | findstr :5173
if %errorlevel% equ 0 (
    echo [结果] 端口 5173 已被占用
) else (
    echo [结果] 端口 5173 未被占用（前端服务可能未运行）
)
echo.

REM 检查 node.exe 进程
echo [检查] Node.js 进程:
tasklist /FI "IMAGENAME eq node.exe" /FO TABLE 2>nul
if %errorlevel% equ 0 (
    echo [结果] 发现 Node.js 进程
) else (
    echo [结果] 未发现 Node.js 进程
)
echo.

REM 检查是否有 Vite 相关的进程
echo [检查] Vite 相关进程:
for /f "tokens=2" %%i in ('tasklist /FI "IMAGENAME eq node.exe" /FO LIST 2^>nul ^| findstr /I "PID"') do (
    set PID=%%i
    for /f "tokens=*" %%j in ('wmic process where "ProcessId=!PID!" get CommandLine /format:list 2^>nul ^| findstr /I "CommandLine"') do (
        set CMDLINE=%%j
        echo !CMDLINE! | findstr /I "vite" >nul 2>&1
        if !errorlevel! equ 0 (
            echo [发现] Vite 进程 PID: !PID!
            echo [命令] !CMDLINE!
        )
    )
)
echo.

REM 检查前端目录
echo [检查] 前端目录:
if exist "frontend\package.json" (
    echo [结果] frontend\package.json 存在
    cd /d "%~dp0frontend"
    if exist "node_modules" (
        echo [结果] node_modules 存在
    ) else (
        echo [警告] node_modules 不存在，需要运行 npm install
    )
    cd /d "%~dp0"
) else (
    echo [错误] frontend\package.json 不存在
)
echo.

REM 尝试连接 localhost:5173
echo [检查] 尝试连接 http://localhost:5173...
powershell -Command "try { $response = Invoke-WebRequest -Uri 'http://localhost:5173' -TimeoutSec 2 -UseBasicParsing; Write-Host '[结果] 连接成功，状态码:' $response.StatusCode } catch { Write-Host '[结果] 连接失败:' $_.Exception.Message }"
echo.

echo ========================================
echo   诊断完成
echo ========================================
pause

