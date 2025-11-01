#!/bin/bash
# {{ AURA-X: Add - Linux/macOS GUI编译脚本. Confirmed via 寸止 }}
# Movie Data Capture GUI Build Script for Linux/macOS

echo "========================================"
echo "Movie Data Capture - GUI Build"
echo "========================================"
echo ""

# 检查Wails CLI
if ! command -v wails &> /dev/null; then
    echo "[错误] Wails CLI 未安装"
    echo "请运行: go install github.com/wailsapp/wails/v2/cmd/wails@latest"
    exit 1
fi

echo "[1/3] 检查环境..."
wails doctor
if [ $? -ne 0 ]; then
    echo "[错误] Wails 环境检查失败"
    exit 1
fi

echo ""
echo "[2/3] 安装前端依赖..."
cd frontend
if [ ! -d "node_modules" ]; then
    npm install
    if [ $? -ne 0 ]; then
        echo "[错误] 安装npm依赖失败"
        cd ..
        exit 1
    fi
fi
cd ..

echo ""
echo "[3/3] 编译GUI应用..."
wails build -tags gui -clean

if [ $? -eq 0 ]; then
    echo ""
    echo "========================================"
    echo "编译成功！"
    echo "可执行文件位于: build/bin/mdc-gui"
    echo "========================================"
else
    echo ""
    echo "========================================"
    echo "编译失败！请检查错误信息。"
    echo "========================================"
    exit 1
fi

echo ""

