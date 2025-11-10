# Test webhook con firma HMAC válida
$appSecret = "451614ef9eb9b35571dc352af6b2110e"

$body = @"
{
  "object": "whatsapp_business_account",
  "entry": [{
    "id": "123456",
    "changes": [{
      "value": {
        "messaging_product": "whatsapp",
        "metadata": {
          "display_phone_number": "16505551111",
          "phone_number_id": "804818756055720"
        },
        "contacts": [{
          "profile": {
            "name": "Test User"
          },
          "wa_id": "593992686734"
        }],
        "messages": [{
          "from": "593992686734",
          "id": "wamid.TEST_$(Get-Date -Format 'yyyyMMddHHmmss')",
          "timestamp": "$([DateTimeOffset]::UtcNow.ToUnixTimeSeconds())",
          "type": "text",
          "text": {
            "body": "Mensaje de prueba local $(Get-Date -Format 'HH:mm:ss')"
          }
        }]
      },
      "field": "messages"
    }]
  }]
}
"@

# Calcular firma HMAC-SHA256
$hmac = New-Object System.Security.Cryptography.HMACSHA256
$hmac.Key = [Text.Encoding]::UTF8.GetBytes($appSecret)
$hash = $hmac.ComputeHash([Text.Encoding]::UTF8.GetBytes($body))
$signature = "sha256=" + [BitConverter]::ToString($hash).Replace("-", "").ToLower()

Write-Host "Signature: $signature"
Write-Host "Enviando webhook..."

# Enviar webhook
$headers = @{
    "X-Hub-Signature-256" = $signature
    "Content-Type" = "application/json"
}

Invoke-WebRequest `
    -Uri "https://whatsapp-api-go-dpb5cgbnaec2gdf2.eastus-01.azurewebsites.net/webhook" `
    -Method POST `
    -Body $body `
    -Headers $headers `
    -UseBasicParsing

