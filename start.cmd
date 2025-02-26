@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

:: 设置窗口标题
title 项目启动脚本

:: 定义服务路径
set GATE_SERVICE=gateService\cmd
set CRAWLER_SERVICE=crawler\cmd
set RECOMMEND_SERVICE=recommendSimilarly
set MONITOR_SERVICE=monitorService\cmd
set ZANIME_WEB=Zanime
set ZANIME_MONITOR=ZanimeMonitor

:: 使用Windows Terminal启动所有服务
wt ^
new-tab -p "Windows PowerShell" --title "Gate Service" -d "%CD%\%GATE_SERVICE%" cmd /k "go run main.go" ^; ^
new-tab -p "Windows PowerShell" --title "Crawler Service" -d "%CD%\%CRAWLER_SERVICE%" cmd /k "go run main.go" ^; ^
new-tab -p "Windows PowerShell" --title "Monitor Service" -d "%CD%\%MONITOR_SERVICE%" cmd /k "go run main.go" ^; ^
new-tab -p "Windows PowerShell" --title "Recommend Service" -d "%CD%\%RECOMMEND_SERVICE%" cmd /k ".\venv\Scripts\activate && python main.py" ^; ^
new-tab -p "Windows PowerShell" --title "Zanime Web" -d "%CD%\%ZANIME_WEB%" cmd /k "npm run dev" ^; ^
new-tab -p "Windows PowerShell" --title "Zanime Monitor" -d "%CD%\%ZANIME_MONITOR%" cmd /k "npm run dev"

echo 所有服务已在Windows Terminal中启动
pause
