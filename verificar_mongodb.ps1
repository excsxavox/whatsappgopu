# ==================================
# Script para verificar MongoDB
# ==================================

Write-Host "`n============================================" -ForegroundColor Cyan
Write-Host " Verificación de MongoDB" -ForegroundColor Cyan
Write-Host "============================================`n" -ForegroundColor Cyan

# MongoDB Atlas URI
$mongoUri = "mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
$dbName = "whatsapp_api"

Write-Host "📊 Conectando a MongoDB Atlas..." -ForegroundColor Yellow
Write-Host "   Base de datos: $dbName`n" -ForegroundColor Gray

# Crear archivo temporal con comandos de MongoDB
$mongoCommands = @"
use $dbName
print('\\n=== COLECCIONES ===')
db.getCollectionNames().forEach(function(col) {
    var count = db[col].countDocuments()
    print(col + ': ' + count + ' documentos')
})

print('\\n=== ÍNDICES ===')
db.getCollectionNames().forEach(function(col) {
    print('\\nColección: ' + col)
    db[col].getIndexes().forEach(function(idx) {
        print('  - ' + idx.name + ': ' + JSON.stringify(idx.key))
    })
})

print('\\n=== ÚLTIMO MENSAJE (si existe) ===')
if (db.messages.countDocuments() > 0) {
    var lastMsg = db.messages.findOne({}, {sort: {'timestamps.created_at': -1}})
    print('ID: ' + lastMsg._id)
    print('Direction: ' + lastMsg.direction)
    print('From: ' + lastMsg.from)
    print('Status: ' + lastMsg.status)
    if (lastMsg.message && lastMsg.message.text) {
        print('Text: ' + lastMsg.message.text.body)
    }
} else {
    print('No hay mensajes aún')
}

print('\\n=== EMPRESAS (si existen) ===')
if (db.companies.countDocuments() > 0) {
    db.companies.find().forEach(function(c) {
        print(c.code + ' - ' + c.name + ' (activa: ' + c.is_active + ')')
    })
} else {
    print('No hay empresas aún')
}
"@

$tempFile = [System.IO.Path]::GetTempFileName()
$mongoCommands | Out-File -FilePath $tempFile -Encoding utf8

Write-Host "🔍 Consultando estado..." -ForegroundColor Yellow

# Ejecutar mongosh (si está instalado)
$mongoshPath = "mongosh"
$mongoshInstalled = $null -ne (Get-Command mongosh -ErrorAction SilentlyContinue)

if ($mongoshInstalled) {
    & $mongoshPath $mongoUri --quiet --file $tempFile
    Remove-Item $tempFile
} else {
    Write-Host "`n⚠️  mongosh no está instalado localmente" -ForegroundColor Yellow
    Write-Host "`n📝 Para verificar manualmente:" -ForegroundColor Cyan
    Write-Host "   1. Ve a: https://cloud.mongodb.com" -ForegroundColor White
    Write-Host "   2. Click en tu cluster" -ForegroundColor White
    Write-Host "   3. Click en 'Browse Collections'" -ForegroundColor White
    Write-Host "   4. Selecciona la base de datos: $dbName" -ForegroundColor White
    Write-Host "`n   O instala mongosh:" -ForegroundColor Cyan
    Write-Host "   https://www.mongodb.com/try/download/shell" -ForegroundColor Gray
    Remove-Item $tempFile
}

Write-Host "`n============================================" -ForegroundColor Cyan
Write-Host " Nota: Las colecciones se crean automáticamente" -ForegroundColor Yellow
Write-Host " cuando ejecutes la aplicación por primera vez." -ForegroundColor Yellow
Write-Host "============================================`n" -ForegroundColor Cyan

