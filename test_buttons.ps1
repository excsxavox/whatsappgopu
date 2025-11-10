# Script para probar envío de botones directamente a WhatsApp Cloud API

$token = "EACO8kt4CNU4BPzzUrgdmDLZBDyxarAfaiBET0SBulEFP4le6D8gJ4eCjuTbZAKfiyY34WNTWOINifv8YEP8tZC4JFnGWAfW7atE332Epks6wKzJe4wMqnuoqKCf0AZBfnrfkTgKK2pjER9yoxxT69P5b1TXA2PE3ObeVmq4NSR2qefcDZBHMh8C43dz0ZAZAPXZC7yGUZAXINVuQFGYdqeoJzX5BWwCQs1jGx4s01C9GSBshG9gVMuCLUVWN9doEKYCuJcZCFV27EaWiHfm1eHNGGxegZDZD"
$phoneNumberId = "804818756055720"
$to = "593992686734"

# Formato de botones según documentación de WhatsApp
$body = @{
    messaging_product = "whatsapp"
    recipient_type = "individual"
    to = $to
    type = "interactive"
    interactive = @{
        type = "button"
        header = @{
            type = "text"
            text = "Prueba de Botones"
        }
        body = @{
            text = "Hola edison! Selecciona una opción:"
        }
        footer = @{
            text = "Powered by WhatsApp API"
        }
        action = @{
            buttons = @(
                @{
                    type = "reply"
                    reply = @{
                        id = "test_1"
                        title = "Opción 1"
                    }
                },
                @{
                    type = "reply"
                    reply = @{
                        id = "test_2"
                        title = "Opción 2"
                    }
                },
                @{
                    type = "reply"
                    reply = @{
                        id = "test_3"
                        title = "Opción 3"
                    }
                }
            )
        }
    }
} | ConvertTo-Json -Depth 10

Write-Host "Enviando botones directamente a WhatsApp Cloud API..."
Write-Host "Payload:"
Write-Host $body
Write-Host ""

try {
    $response = Invoke-WebRequest `
        -Uri "https://graph.facebook.com/v20.0/$phoneNumberId/messages" `
        -Method POST `
        -Headers @{
            "Authorization" = "Bearer $token"
            "Content-Type" = "application/json"
        } `
        -Body $body `
        -UseBasicParsing

    Write-Host "✅ Respuesta exitosa:"
    Write-Host $response.Content
} catch {
    Write-Host "❌ Error:"
    Write-Host $_.Exception.Message
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $responseBody = $reader.ReadToEnd()
        Write-Host "Body de error:"
        Write-Host $responseBody
    }
}

