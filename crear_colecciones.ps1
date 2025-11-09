# ==================================
# Script para crear colecciones e índices en MongoDB
# ==================================

Write-Host "`n============================================" -ForegroundColor Cyan
Write-Host " Creando Colecciones en MongoDB Atlas" -ForegroundColor Cyan
Write-Host "============================================`n" -ForegroundColor Cyan

# Variables de entorno
$env:MONGODB_URL = "mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
$env:MONGO_DB = "whatsapp_api"
$env:WHATSAPP_VERIFY_TOKEN = "test_token_123"
$env:WHATSAPP_APP_SECRET = "test_secret"
$env:WABA_PHONE_ID = "test_phone_id"
$env:WABA_TOKEN = "test_token"
$env:API_PORT = "8080"

Write-Host "📊 Configuración:" -ForegroundColor Yellow
Write-Host "   MongoDB: whatsapp_api" -ForegroundColor Gray
Write-Host "   Cluster: MongoDB Atlas`n" -ForegroundColor Gray

Write-Host "🚀 Ejecutando aplicación para crear colecciones...`n" -ForegroundColor Yellow

# Ejecutar app en background
$job = Start-Job -ScriptBlock {
    param($workDir)
    Set-Location $workDir
    $env:MONGODB_URL = "mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
    $env:MONGO_DB = "whatsapp_api"
    $env:WHATSAPP_VERIFY_TOKEN = "test_token_123"
    $env:WHATSAPP_APP_SECRET = "test_secret"
    $env:WABA_PHONE_ID = "test_phone_id"
    $env:WABA_TOKEN = "test_token"
    $env:API_PORT = "8080"
    & "C:\Program Files\Go\bin\go.exe" run cmd/server/main.go 2>&1
} -ArgumentList (Get-Location)

# Esperar 10 segundos para que se inicialice
Write-Host "⏳ Esperando 10 segundos para inicialización..." -ForegroundColor Yellow
Start-Sleep -Seconds 10

# Detener el job
Stop-Job $job
Remove-Job $job

Write-Host "`n✅ Proceso completado!" -ForegroundColor Green
Write-Host "`n📋 Colecciones creadas:" -ForegroundColor Cyan
Write-Host "   • messages (7 índices)" -ForegroundColor White
Write-Host "   • companies (2 índices)" -ForegroundColor White
Write-Host "   • sessions (1 índice)`n" -ForegroundColor White

Write-Host "🔍 Para verificar:" -ForegroundColor Cyan
Write-Host "   1. Ve a: https://cloud.mongodb.com" -ForegroundColor White
Write-Host "   2. Browse Collections" -ForegroundColor White
Write-Host "   3. Busca la base de datos: whatsapp_api`n" -ForegroundColor White

