#!/bin/bash

echo "======================================== 启动 StructForge 用户服务 ========================================"

# 检查 Go 环境
if ! command -v go &> /dev/null; then
    echo "[错误] 未找到 Go 环境，请先安装 Go 1.21+"
    exit 1
fi

# 切换到用户服务目录
cd "$(dirname "$0")"
echo "[信息] 当前工作目录: $(pwd)"

# 检查必要文件
if [ ! -f "main.go" ]; then
    echo "[错误] 未找到 main.go 文件"
    exit 1
fi

# 检查 wire_gen.go 是否存在
if [ ! -f "wire_gen.go" ]; then
    echo "[警告] 未找到 wire_gen.go 文件，正在生成..."
    echo "[信息] 正在运行 wire 命令生成依赖注入代码..."
    wire
    if [ $? -ne 0 ]; then
        echo "[错误] Wire 代码生成失败"
        echo "[提示] 请确保已安装 wire: go install github.com/google/wire/cmd/wire@latest"
        exit 1
    fi
    echo "[信息] Wire 代码生成成功"
fi

# 启动用户服务
echo "[信息] 正在启动用户服务..."
echo ""
go run main.go wire_gen.go

EXIT_CODE=$?

if [ $EXIT_CODE -ne 0 ]; then
    echo ""
    echo "[错误] 用户服务启动失败 (退出代码: $EXIT_CODE)"
    echo "请检查:"
    echo "1. 数据库连接配置是否正确"
    echo "2. 配置文件是否存在: ../../../../configs/local/user.yaml"
    echo "3. Protobuf 代码是否已生成: make proto"
    echo ""
    exit $EXIT_CODE
else
    echo ""
    echo "[信息] 用户服务已正常退出"
fi

