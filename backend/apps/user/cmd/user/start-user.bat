@echo off
setlocal enabledelayedexpansion
chcp 65001 >nul

echo ======================================== 启动 StructForge 用户服务 ========================================

REM 检查 Go 环境
where go >nul 2>&1
if %errorlevel% neq 0 (
    echo [错误] 未找到 Go 环境，请先安装 Go 1.21+
    pause
    exit /b 1
)

REM 切换到用户服务目录
cd /d "%~dp0"
echo [信息] 当前工作目录: %CD%

REM 检查必要文件
if not exist "main.go" (
    echo [错误] 未找到 main.go 文件
    pause
    exit /b 1
)

REM 检查 wire_gen.go 是否存在
if not exist "wire_gen.go" (
    echo [警告] 未找到 wire_gen.go 文件，正在生成...
    echo [信息] 正在运行 wire 命令生成依赖注入代码...
    wire
    if !errorlevel! neq 0 (
        echo [错误] Wire 代码生成失败
        echo [提示] 请确保已安装 wire: go install github.com/google/wire/cmd/wire@latest
        pause
        exit /b 1
    )
    echo [信息] Wire 代码生成成功
)

REM 启动用户服务
echo [信息] 正在启动用户服务...
echo.
go run main.go wire_gen.go

set EXIT_CODE=!errorlevel!

if !EXIT_CODE! neq 0 (
    echo.
    echo [错误] 用户服务启动失败 (退出代码: !EXIT_CODE!)
    echo 请检查:
    echo 1. 数据库连接配置是否正确
    echo 2. 配置文件是否存在: ../../../../configs/local/user.yaml
    echo 3. Protobuf 代码是否已生成: make proto
    echo.
    pause
) else (
    echo.
    echo [信息] 用户服务已正常退出
    pause
)

