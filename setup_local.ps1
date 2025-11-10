# Script para configurar entorno local

Write-Host "🔧 Configurando entorno local..." -ForegroundColor Cyan

# Crear archivo .env
@'
API_PORT=8080
WHATSAPP_VERIFY_TOKEN=mi_token_secreto_123
WHATSAPP_APP_SECRET=451614ef9eb9b35571dc352af6b2110e
WABA_PHONE_ID=804818756055720
WABA_TOKEN=EACO8kt4CNU4BP4xgZAS2jwSIsUNwOS5ggJPnYz5WvPZCprIjoP5PfSE8JYD59lvwzBBAeTTKwQiVFdkGhLQoq7aaPpU1ZCYKV6mSZCmw7973W3q305S8B36ZAe19P75ZCsUqxpwWJaom0UebC0A10R3aNrfl7Tc2ItyrslOHKR7RR2SmrXABDG4lRdO0R3HjJZCXrwPYHVguWxLUFXQGU6yNFGVlqTtm77X0aYfOp5fAbgc98VYhshNdzGhkGZADNA5Qx0sxRAUVITxPOI7tYxMa
WABA_API_VERSION=v20.0
MONGO_URI=mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0
MONGO_DB=test
LOG_LEVEL=INFO
'@ | Out-File -FilePath ".env" -Encoding UTF8
Write-Host "✅ Archivo .env creado" -ForegroundColor Green

# Compilar
Write-Host "`n🔨 Compilando aplicación..." -ForegroundColor Cyan
go build -o whatsapp-api.exe cmd/server/main.go

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Compilación exitosa" -ForegroundColor Green
    Write-Host "`n📋 PRÓXIMOS PASOS:" -ForegroundColor Yellow
    Write-Host "1. Ejecuta: .\whatsapp-api.exe" -ForegroundColor White
    Write-Host "2. En otra terminal ejecuta: ngrok http 8080" -ForegroundColor White
    Write-Host "3. Configura el webhook en Meta con la URL de ngrok" -ForegroundColor White
    Write-Host "4. Envía un mensaje a WhatsApp y observa los logs" -ForegroundColor White
} else {
    Write-Host "❌ Error en la compilación" -ForegroundColor Red
}

