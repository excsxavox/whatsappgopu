# 🚀 Deploy en Fly.io (GRATIS)

## ✅ Ventajas
- GRATIS hasta 3 aplicaciones
- HTTPS automático
- Dominio incluido: tu-app.fly.dev
- Deploy en 2 minutos
- Perfecto para Go y WhatsApp webhooks

## 📋 Pasos

### 1. Instalar Fly CLI
```powershell
powershell -Command "iwr https://fly.io/install.ps1 -useb | iex"
```

### 2. Crear cuenta (GRATIS)
```bash
fly auth signup
```

### 3. Inicializar app
```bash
fly launch
```

Responde:
- **App name:** whatsapp-api-tuempresa (o deja en blanco para auto)
- **Region:** Miami (mia) - más cercano a Latinoamérica
- **Deploy now?** NO (aún falta configurar variables)

### 4. Configurar variables de entorno
```bash
fly secrets set MONGODB_URL="mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
fly secrets set MONGO_DB="whatsapp_api"
fly secrets set WHATSAPP_VERIFY_TOKEN="mi_token_seguro_whatsapp_2024"
fly secrets set WHATSAPP_APP_SECRET="451614ef9eb9b35571dc352af6b2110e"
fly secrets set WABA_PHONE_ID="804818756055720"
fly secrets set WABA_TOKEN="EACO8kt4CNU4BP5ZAjnyEZBsatkJmx2XvPvYfO9cCllcjZANi1UeTvuh9LBWQ2t3Rse4B0q4rij37Ml3vgiFQB6krHWYhdW6mUkfRZBrA6w3ZBOZBYL1AAgTZCyS1Ls5zB4OwZAqPB2Dgpcz8Ucn2TjnPzGbVD3zza6IKmlGlYsaLSC3SNBHXvWjNj4W1FRPXtiY7y2ksG7n7xDzZBNe6kYTM3p0OZCc5RivuQrwpb6v4D4lKRGD5Ut2R2ownJzp8NRJ2BgfoeLttq5Cw5FOup76vssYgZDZD"
fly secrets set API_PORT="8080"
```

### 5. Deploy
```bash
fly deploy
```

### 6. Obtener URL
```bash
fly info
```

Tu URL será: **https://tu-app.fly.dev**

### 7. Configurar webhook en Meta
URL del webhook: `https://tu-app.fly.dev/webhook`
Verify Token: `mi_token_seguro_whatsapp_2024`

## 🔄 Comandos útiles

```bash
# Ver logs
fly logs

# Ver status
fly status

# Abrir dashboard
fly dashboard

# SSH a la app
fly ssh console

# Escalar (si necesitas más recursos)
fly scale vm shared-cpu-1x --memory 512
```

## 💰 Límites Plan Gratis
- 3 aplicaciones
- 160 GB transferencia/mes
- 256 MB RAM por app
- Suficiente para WhatsApp API

## 🎯 Resultado
✅ URL: https://tu-app.fly.dev
✅ HTTPS automático
✅ Webhook: https://tu-app.fly.dev/webhook
✅ Siempre online

