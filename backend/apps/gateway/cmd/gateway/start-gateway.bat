@echo off
setlocal enabledelayedexpansion
chcp 65001 >nul
echo ========================================
echo   启动 Gateway 服务
echo ========================================
echo.

REM 检查 Go 环境
where go >nul 2>&1
if %errorlevel% neq 0 (
    echo [错误] 未找到 Go 环境，请先安装 Go 1.21+
    echo 下载地址: https://golang.org/dl/
    pause
    exit /b 1
)

REM 检查必要的文件是否存在
if not exist "%~dp0main.go" (
    echo [错误] 未找到 main.go 文件
    pause
    exit /b 1
)

if not exist "%~dp0wire_gen.go" (
    echo [警告] 未找到 wire_gen.go 文件，正在生成...
    cd /d %~dp0
    wire
    if %errorlevel% neq 0 (
        echo [错误] wire 命令执行失败
        pause
        exit /b 1
    )
)

echo [信息] 正在启动 Gateway 服务...
echo [提示] Gateway 服务将运行在: http://localhost:8000
echo [提示] 按 Ctrl+C 停止服务
echo.

cd /d %~dp0
go run . -env local

pause

