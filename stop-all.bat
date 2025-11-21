@echo off
setlocal enabledelayedexpansion
chcp 65001 >nul
echo ========================================
echo   停止 StructForge 全栈服务
echo   前端 + 后端
echo ========================================
echo.

REM 停止前端服务（Node.js/Vite）
echo [信息] 正在停止前端服务...

REM 通过端口 5173 查找并停止进程
netstat -ano | findstr :5173 >nul 2>&1
if %errorlevel% equ 0 (
    echo [信息] 发现端口 5173 被占用，正在停止相关进程...
    for /f "tokens=5" %%a in ('netstat -ano ^| findstr :5173') do (
        echo [信息] 终止占用端口 5173 的进程: %%a
        taskkill /PID %%a /T /F >nul 2>&1
    )
) else (
    echo [信息] 端口 5173 未被占用
)

REM 停止前端相关窗口和进程
taskkill /FI "WINDOWTITLE eq StructForge - 前端服务*" /T /F >nul 2>&1
taskkill /FI "IMAGENAME eq node.exe" /FI "WINDOWTITLE eq *StructForge*" /T /F >nul 2>&1
taskkill /FI "IMAGENAME eq node.exe" /FI "COMMANDLINE eq *vite*" /T /F >nul 2>&1
taskkill /FI "IMAGENAME eq node.exe" /FI "COMMANDLINE eq *npm*" /FI "COMMANDLINE eq *dev*" /T /F >nul 2>&1

REM 停止 Gateway 服务
echo [信息] 正在停止 Gateway 服务...
taskkill /FI "WINDOWTITLE eq StructForge - Gateway 服务*" /T /F >nul 2>&1
taskkill /FI "WINDOWTITLE eq StructForge - 后端服务*" /T /F >nul 2>&1
taskkill /FI "IMAGENAME eq go.exe" /FI "COMMANDLINE eq *gateway*" /T /F >nul 2>&1

REM 停止 User 服务
echo [信息] 正在停止 User 服务...
taskkill /FI "WINDOWTITLE eq StructForge - User 服务*" /T /F >nul 2>&1
taskkill /FI "IMAGENAME eq go.exe" /FI "COMMANDLINE eq *user*" /T /F >nul 2>&1

REM 停止所有相关的 Go 进程（排除 Nacos）
echo [信息] 正在清理 Go 进程...
for /f "tokens=2" %%i in ('tasklist /FI "IMAGENAME eq go.exe" /FO LIST 2^>nul ^| findstr /I "PID"') do (
    set PID=%%i
    set FOUND=0
    for /f "tokens=*" %%j in ('wmic process where "ProcessId=!PID!" get CommandLine /format:list 2^>nul ^| findstr /I "CommandLine"') do (
        set CMDLINE=%%j
        echo !CMDLINE! | findstr /I "nacos" >nul 2>&1
        if !errorlevel! equ 0 (
            set FOUND=1
        )
        if !FOUND! equ 0 (
            echo !CMDLINE! | findstr /I "gateway" >nul 2>&1
            if !errorlevel! equ 0 (
                echo [信息] 终止 Gateway 进程: !PID!
                taskkill /PID !PID! /T /F >nul 2>&1
                set FOUND=1
            )
        )
        if !FOUND! equ 0 (
            echo !CMDLINE! | findstr /I "user" >nul 2>&1
            if !errorlevel! equ 0 (
                echo !CMDLINE! | findstr /I "gateway" >nul 2>&1
                if !errorlevel! neq 0 (
                    echo !CMDLINE! | findstr /I "nacos" >nul 2>&1
                    if !errorlevel! neq 0 (
                        echo [信息] 终止 User 进程: !PID!
                        taskkill /PID !PID! /T /F >nul 2>&1
                        set FOUND=1
                    )
                )
            )
        )
    )
)

REM 停止所有相关的 Node.js 进程
echo [信息] 正在清理 Node.js 进程...
for /f "tokens=2" %%i in ('tasklist /FI "IMAGENAME eq node.exe" /FO LIST 2^>nul ^| findstr /I "PID"') do (
    set PID=%%i
    for /f "tokens=*" %%j in ('wmic process where "ProcessId=!PID!" get CommandLine /format:list 2^>nul ^| findstr /I "CommandLine"') do (
        set CMDLINE=%%j
        echo !CMDLINE! | findstr /I "vite" >nul 2>&1
        if !errorlevel! equ 0 (
            echo [信息] 终止 Vite 进程: !PID!
            taskkill /PID !PID! /T /F >nul 2>&1
        )
        echo !CMDLINE! | findstr /I "frontend" >nul 2>&1
        if !errorlevel! equ 0 (
            echo [信息] 终止前端进程: !PID!
            taskkill /PID !PID! /T /F >nul 2>&1
        )
    )
)

REM 关闭所有相关的命令窗口（排除 Nacos）
echo [信息] 正在关闭相关命令窗口...
for /f "tokens=2" %%i in ('tasklist /FI "IMAGENAME eq cmd.exe" /FO LIST 2^>nul ^| findstr /I "PID"') do (
    set PID=%%i
    set FOUND=0
    for /f "tokens=*" %%j in ('wmic process where "ProcessId=!PID!" get CommandLine /format:list 2^>nul ^| findstr /I "CommandLine"') do (
        set CMDLINE=%%j
        REM 排除 Nacos 相关窗口
        echo !CMDLINE! | findstr /I "nacos" >nul 2>&1
        if !errorlevel! equ 0 (
            set FOUND=1
        )
        if !FOUND! equ 0 (
            echo !CMDLINE! | findstr /I "StructForge" >nul 2>&1
            if !errorlevel! equ 0 (
                echo [信息] 关闭 StructForge 窗口: !PID!
                taskkill /PID !PID! /T /F >nul 2>&1
                set FOUND=1
            )
        )
        if !FOUND! equ 0 (
            echo !CMDLINE! | findstr /I "start-backend" >nul 2>&1
            if !errorlevel! equ 0 (
                echo [信息] 关闭后端窗口: !PID!
                taskkill /PID !PID! /T /F >nul 2>&1
                set FOUND=1
            )
        )
        if !FOUND! equ 0 (
            echo !CMDLINE! | findstr /I "start-frontend" >nul 2>&1
            if !errorlevel! equ 0 (
                echo [信息] 关闭前端窗口: !PID!
                taskkill /PID !PID! /T /F >nul 2>&1
                set FOUND=1
            )
        )
        if !FOUND! equ 0 (
            echo !CMDLINE! | findstr /I "gateway" >nul 2>&1
            if !errorlevel! equ 0 (
                echo !CMDLINE! | findstr /I "nacos" >nul 2>&1
                if !errorlevel! neq 0 (
                    echo [信息] 关闭 Gateway 窗口: !PID!
                    taskkill /PID !PID! /T /F >nul 2>&1
                    set FOUND=1
                )
            )
        )
        if !FOUND! equ 0 (
            echo !CMDLINE! | findstr /I "user" >nul 2>&1
            if !errorlevel! equ 0 (
                echo !CMDLINE! | findstr /I "gateway" >nul 2>&1
                if !errorlevel! neq 0 (
                    echo !CMDLINE! | findstr /I "nacos" >nul 2>&1
                    if !errorlevel! neq 0 (
                        echo [信息] 关闭 User 窗口: !PID!
                        taskkill /PID !PID! /T /F >nul 2>&1
                        set FOUND=1
                    )
                )
            )
        )
    )
)

REM 等待进程完全关闭
echo [信息] 等待进程完全关闭...
timeout /t 2 /nobreak >nul 2>&1

REM 再次检查端口 5173 是否已释放
netstat -ano | findstr :5173 >nul 2>&1
if %errorlevel% equ 0 (
    echo [警告] 端口 5173 仍被占用，尝试强制释放...
    for /f "tokens=5" %%a in ('netstat -ano ^| findstr :5173') do (
        taskkill /PID %%a /F >nul 2>&1
    )
    timeout /t 1 /nobreak >nul 2>&1
) else (
    echo [成功] 端口 5173 已释放
)

REM 最终检查并强制关闭
echo [信息] 进行最终清理...
taskkill /FI "WINDOWTITLE eq *StructForge*" /T /F >nul 2>&1
taskkill /FI "IMAGENAME eq go.exe" /FI "COMMANDLINE eq *gateway*" /T /F >nul 2>&1
taskkill /FI "IMAGENAME eq go.exe" /FI "COMMANDLINE eq *user*" /T /F >nul 2>&1
taskkill /FI "IMAGENAME eq node.exe" /FI "COMMANDLINE eq *vite*" /T /F >nul 2>&1
taskkill /FI "IMAGENAME eq node.exe" /FI "COMMANDLINE eq *frontend*" /T /F >nul 2>&1

echo.
echo ========================================
echo   服务停止完成！
echo ========================================
echo.
echo [信息] 所有 StructForge 相关进程已停止
echo [信息] 相关命令窗口已关闭
echo.
echo 按任意键退出...
pause >nul
