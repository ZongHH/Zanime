@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

:: 设置窗口标题
title 项目停止脚本

:: 关闭Go服务进程
taskkill /F /IM go.exe 2>nul
taskkill /F /IM main.exe 2>nul

:: 关闭Python进程
taskkill /F /IM python.exe 2>nul

:: 关闭Node.js进程
taskkill /F /IM node.exe 2>nul

:: 关闭Windows Terminal
taskkill /F /IM WindowsTerminal.exe 2>nul

echo 所有服务已停止
pause 