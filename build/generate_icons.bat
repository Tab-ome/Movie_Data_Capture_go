@echo off
REM {{ AURA-X: Create - Windows图标生成批处理脚本 }}
chcp 65001 >nul
echo.
echo 🎨 Movie Data Capture 图标生成工具
echo ====================================
echo.

REM 检查Python是否安装
python --version >nul 2>&1
if errorlevel 1 (
    echo ❌ 错误：未检测到Python，请先安装Python 3.7+
    echo    下载地址：https://www.python.org/downloads/
    pause
    exit /b 1
)

echo ✅ Python 已安装
echo.

REM 检查并安装依赖
echo 📦 检查依赖库...
python -c "import cairosvg, PIL" >nul 2>&1
if errorlevel 1 (
    echo ⚠️  缺少依赖库，正在安装...
    echo.
    pip install cairosvg pillow
    if errorlevel 1 (
        echo.
        echo ❌ 依赖安装失败，请手动执行：
        echo    pip install cairosvg pillow
        pause
        exit /b 1
    )
) else (
    echo ✅ 依赖库已就绪
)

echo.
echo 🚀 开始生成图标...
echo.

REM 执行Python脚本
python "%~dp0generate_icons.py"

if errorlevel 1 (
    echo.
    echo ❌ 图标生成失败
    pause
    exit /b 1
)

echo.
echo ✅ 图标生成成功！
echo.
echo 💡 提示：现在可以运行 build-gui.bat 重新编译应用
echo.
pause

