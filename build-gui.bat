@echo off
REM {{ AURA-X: Add - Windows GUI编译脚本. Confirmed via 寸止 }}
REM Movie Data Capture GUI Build Script for Windows

echo ========================================
echo Movie Data Capture - GUI Build
echo ========================================
echo.

REM 检查Wails CLI
where wails >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo [错误] Wails CLI 未安装
    echo 请运行: go install github.com/wailsapp/wails/v2/cmd/wails@latest
    pause
    exit /b 1
)

echo [1/3] 检查环境...
wails doctor
if %ERRORLEVEL% NEQ 0 (
    echo [错误] Wails 环境检查失败
    pause
    exit /b 1
)

echo.
echo [2/3] 安装前端依赖...
cd frontend
if not exist node_modules (
    call npm install
    if %ERRORLEVEL% NEQ 0 (
        echo [错误] 安装npm依赖失败
        cd ..
        pause
        exit /b 1
    )
)
cd ..

echo.
echo [3/3] 编译GUI应用...
wails build -tags gui -clean

if %ERRORLEVEL% EQ 0 (
    echo.
    echo ========================================
    echo 编译成功！
    echo 可执行文件位于: build\bin\mdc-gui.exe
    echo ========================================
) else (
    echo.
    echo ========================================
    echo 编译失败！请检查错误信息。
    echo ========================================
)

echo.
pause

