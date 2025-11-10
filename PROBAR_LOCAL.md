# 🏠 Probar Flows Localmente con ngrok

## Paso 1: Crear archivo .env local

Crea el archivo `.env` en la raíz del proyecto:

```env
# API Server
API_PORT=8080

# WhatsApp Cloud API
WHATSAPP_VERIFY_TOKEN=mi_token_secreto_123
WHATSAPP_APP_SECRET=451614ef9eb9b35571dc352af6b2110e
WABA_PHONE_ID=804818756055720
WABA_TOKEN=EACO8kt4CNU4BP4xgZAS2jwSIsUNwOS5ggJPnYz5WvPZCprIjoP5PfSE8JYD59lvwzBBAeTTKwQiVFdkGhLQoq7aaPpU1ZCYKV6mSZCmw7973W3q305S8B36ZAe19P75ZCsUqxpwWJaom0UebC0A10R3aNrfl7Tc2ItyrslOHKR7RR2SmrXABDG4lRdO0R3HjJZCXrwPYHVguWxLUFXQGU6yNFGVlqTtm77X0aYfOp5fAbgc98VYhshNdzGhkGZADNA5Qx0sxRAUVITxPOI7tYxMa
WABA_API_VERSION=v20.0

# MongoDB
MONGO_URI=mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0
MONGO_DB=test

# Logs
LOG_LEVEL=INFO
```

## Paso 2: Compilar y ejecutar

```powershell
# Compilar
go build -o whatsapp-api.exe cmd/api/main.go

# Ejecutar
.\whatsapp-api.exe
```

## Paso 3: Instalar y ejecutar ngrok

```powershell
# Descargar ngrok si no lo tienes: https://ngrok.com/download

# Exponer el puerto 8080
ngrok http 8080
```

Te dará una URL pública como: `https://abc123.ngrok.io`

## Paso 4: Configurar webhook en Meta

1. Ve a Meta Developer Console: https://developers.facebook.com/apps/10058963160806734
2. WhatsApp → Configuration
3. Webhook → Edit
4. Callback URL: `https://TU-URL-NGROK.ngrok.io/webhook`
5. Verify Token: `mi_token_secreto_123`
6. Subscribe to: `messages`
7. Verify and Save

## Paso 5: Probar

Ahora cuando envíes mensajes a WhatsApp, verás los logs en tiempo real en tu terminal donde ejecutaste `whatsapp-api.exe`.

---

## 🎯 Ventajas de probar local:

✅ Logs en tiempo real  
✅ Debugging más fácil  
✅ Cambios instantáneos (sin rebuild de Docker)  
✅ Ver exactamente qué está pasando con Meta  

---

## ⚠️ Recuerda:

Después de probar, vuelve a configurar el webhook de Meta apuntando a Azure:
```
https://whatsapp-api-go-dpb5cgbnaec2gdf2.eastus-01.azurewebsites.net/webhook
```

