# Script de prueba local para Google TTS + WhatsApp
# Uso: .\test_tts_local.ps1

param(
    [Parameter(Mandatory=$true)]
    [string]$Token,
    
    [Parameter(Mandatory=$true)]
    [string]$PhoneNumberID,
    
    [Parameter(Mandatory=$true)]
    [string]$ToPhone,
    
    [string]$ApiVersion = "v21.0",
    
    [string]$GoogleTTSKey = $env:GOOGLE_TTS_API_KEY,
    [string]$OpenAIKey = $env:OPENAI_API_KEY,
    [string]$GoogleTTSVoice = $env:GOOGLE_TTS_VOICE,
    [string]$GoogleTTSModel = $env:GOOGLE_TTS_MODEL,
    [string]$GoogleTTSLanguage = $env:GOOGLE_TTS_LANGUAGE,
    [string]$GoogleTTSPrompt = $env:GOOGLE_TTS_PROMPT
)

Write-Host "============================================" -ForegroundColor Cyan
Write-Host "🧪 Prueba Local: Google TTS + WhatsApp" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan
Write-Host ""

# Configurar variables de entorno
$env:WABA_TOKEN = $Token
$env:WABA_PHONE_ID = $PhoneNumberID
$env:WABA_API_VERSION = $ApiVersion
$env:TO_PHONE = $ToPhone

if ($GoogleTTSKey) {
    $env:GOOGLE_TTS_API_KEY = $GoogleTTSKey
}
if ($OpenAIKey) {
    $env:OPENAI_API_KEY = $OpenAIKey
}
if ($GoogleTTSVoice) {
    $env:GOOGLE_TTS_VOICE = $GoogleTTSVoice
}
if ($GoogleTTSModel) {
    $env:GOOGLE_TTS_MODEL = $GoogleTTSModel
}
if ($GoogleTTSLanguage) {
    $env:GOOGLE_TTS_LANGUAGE = $GoogleTTSLanguage
}
if ($GoogleTTSPrompt) {
    $env:GOOGLE_TTS_PROMPT = $GoogleTTSPrompt
}

Write-Host "📋 Configuración:" -ForegroundColor Yellow
Write-Host "   Token: $($Token.Substring(0, 20))..." -ForegroundColor Gray
Write-Host "   Phone Number ID: $PhoneNumberID" -ForegroundColor Gray
Write-Host "   Destino: $ToPhone" -ForegroundColor Gray
Write-Host "   API Version: $ApiVersion" -ForegroundColor Gray
if ($GoogleTTSKey) {
    Write-Host "   Google TTS: ✅ Configurado" -ForegroundColor Green
} else {
    Write-Host "   Google TTS: ⚠️ No configurado (usará OpenAI)" -ForegroundColor Yellow
}
Write-Host ""

# Ejecutar prueba
Write-Host "🚀 Ejecutando prueba..." -ForegroundColor Cyan
Write-Host ""

cd $PSScriptRoot
go run cmd/test_tts/main.go $Token $PhoneNumberID $ToPhone $ApiVersion

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "✅ Prueba completada exitosamente!" -ForegroundColor Green
} else {
    Write-Host ""
    Write-Host "❌ Error en la prueba" -ForegroundColor Red
    exit $LASTEXITCODE
}

