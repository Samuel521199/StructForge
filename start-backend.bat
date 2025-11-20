@echo off
setlocal enabledelayedexpansion
chcp 65001 >nul
echo ========================================
echo   启动 StructForge 后端服务
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
if not exist "%~dp0backend\apps\gateway\cmd\gateway\main.go" (
    echo [错误] 未找到 Gateway 服务文件
    pause
    exit /b 1
)

if not exist "%~dp0backend\apps\user\cmd\user\main.go" (
    echo [错误] 未找到 User 服务文件
    pause
    exit /b 1
)

echo [信息] 正在启动 Gateway 服务...
start "StructForge - Gateway 服务" cmd /k "cd /d %~dp0backend\apps\gateway\cmd\gateway && echo [Gateway] 工作目录: %CD% && echo [Gateway] 正在启动... && go run . -env local && pause"

REM 等待 Gateway 启动
echo [信息] 等待 Gateway 服务启动...
timeout /t 3 /nobreak >nul

echo [信息] 正在启动 User 服务...
start "StructForge - User 服务" cmd /k "cd /d %~dp0backend\apps\user\cmd\user && echo [User] 工作目录: %CD% && echo [User] 正在启动... && go run main.go wire_gen.go && pause"

echo.
echo [信息] 后端服务启动完成
echo [提示] Gateway 服务运行在: http://localhost:8000
echo [提示] User 服务运行在: http://localhost:8001 (HTTP), :9001 (gRPC)
echo [提示] 关闭服务窗口即可停止对应服务
echo.
pause
