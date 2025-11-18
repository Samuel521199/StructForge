@echo off
chcp 65001 >nul
echo ========================================
echo   停止 StructForge 所有服务
echo ========================================
echo.

echo [信息] 正在查找并停止相关进程...
echo.

REM 停止 Node.js 进程（前端）
echo [信息] 停止前端服务...
taskkill /FI "WINDOWTITLE eq StructForge - 前端服务*" /T /F >nul 2>&1
taskkill /FI "IMAGENAME eq node.exe" /FI "WINDOWTITLE eq *Vite*" /T /F >nul 2>&1

REM 停止 Go 进程（后端）
echo [信息] 停止后端服务...
taskkill /FI "WINDOWTITLE eq StructForge - 后端服务*" /T /F >nul 2>&1
taskkill /FI "IMAGENAME eq go.exe" /FI "COMMANDLINE eq *gateway*" /T /F >nul 2>&1

echo.
echo [信息] 服务已停止
echo.
pause

