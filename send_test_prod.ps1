# Enviar mensaje de prueba en modo producción
$body = @{
    phone = "593992686734"
    message = "APP EN PRODUCCION - Responde para activar el flow"
} | ConvertTo-Json

Invoke-WebRequest `
    -Uri "https://whatsapp-api-go-dpb5cgbnaec2gdf2.eastus-01.azurewebsites.net/send" `
    -Method POST `
    -Body $body `
    -ContentType "application/json" `
    -UseBasicParsing

