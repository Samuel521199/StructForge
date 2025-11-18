@echo off
chcp 65001 >nul
echo ========================================
echo   启动 StructForge 后端服务 (Gateway)
echo ========================================
echo.

REM 切换到后端目录
cd /d "%~dp0backend\apps\gateway\cmd\gateway"

REM 检查 Go 环境
where go >nul 2>&1
if %errorlevel% neq 0 (
    echo [错误] 未找到 Go 环境，请先安装 Go 1.21+
    echo 下载地址: https://golang.org/dl/
    pause
    exit /b 1
)

echo [信息] 当前工作目录: %CD%
echo [信息] 正在启动 Gateway 服务...
echo.

REM 运行 Gateway 服务
REM 使用 local 环境，配置文件路径: ../../../../configs/local/gateway.yaml
go run main.go wire_gen.go app.go logger.go -env local

REM 如果服务退出，显示错误信息
if %errorlevel% neq 0 (
    echo.
    echo [错误] Gateway 服务启动失败
    echo 请检查:
    echo 1. Go 环境是否正确安装
    echo 2. 配置文件是否存在: backend\configs\local\gateway.yaml
    echo 3. 依赖是否已安装: cd backend ^&^& go mod tidy
    echo.
    pause
)

