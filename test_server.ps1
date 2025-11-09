# Test WhatsApp API Server

Write-Host "Probando servidor WhatsApp API..." -ForegroundColor Cyan
Write-Host ""

try {
    $health = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method Get -ErrorAction Stop
    Write-Host "[OK] Health: $($health.status)" -ForegroundColor Green
    
    $status = Invoke-RestMethod -Uri "http://localhost:8080/status" -Method Get -ErrorAction Stop
    Write-Host "[OK] Conectado: $($status.connected)" -ForegroundColor Green
    Write-Host "[OK] Autenticado: $($status.logged_in)" -ForegroundColor Green
}
catch {
    Write-Host "[ERROR] Servidor no responde" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Yellow
}
