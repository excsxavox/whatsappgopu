@echo off
chcp 65001 >nul
cls
echo.
echo ============================================
echo   WhatsApp API Server Core
echo   Arquitectura Hexagonal
echo ============================================
echo.
echo Iniciando servidor...
echo.

"C:\Program Files\Go\bin\go.exe" run cmd/server/main.go

echo.
echo ============================================
echo El servidor se ha detenido
echo ============================================
pause

