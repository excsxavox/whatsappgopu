# 🎨 Deploy en Render (GRATIS)

## ✅ Ventajas
- Plan gratis disponible
- HTTPS automático
- Deploy desde GitHub

## ⚠️ Limitaciones Plan Gratis
- La app se "duerme" después de 15 min sin uso
- Tarda ~30 seg en despertar
- 750 horas/mes gratis

## 📋 Pasos

### 1. Subir a GitHub
```bash
git init
git add .
git commit -m "WhatsApp API"
git push
```

### 2. Ir a Render
```
https://render.com
```

### 3. Crear Web Service
- New → Web Service
- Conecta GitHub
- Selecciona tu repo

### 4. Configuración
```
Name: whatsapp-api
Region: Oregon (US West)
Branch: main
Runtime: Docker
```

### 5. Variables de entorno
```
MONGODB_URL=mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0
MONGO_DB=whatsapp_api
WHATSAPP_VERIFY_TOKEN=mi_token_seguro_whatsapp_2024
WHATSAPP_APP_SECRET=451614ef9eb9b35571dc352af6b2110e
WABA_PHONE_ID=804818756055720
WABA_TOKEN=EACO8kt4CNU4BP5ZAjnyEZBsatkJmx2XvPvYfO9cCllcjZANi1UeTvuh9LBWQ2t3Rse4B0q4rij37Ml3vgiFQB6krHWYhdW6mUkfRZBrA6w3ZBOZBYL1AAgTZCyS1Ls5zB4OwZAqPB2Dgpcz8Ucn2TjnPzGbVD3zza6IKmlGlYsaLSC3SNBHXvWjNj4W1FRPXtiY7y2ksG7n7xDzZBNe6kYTM3p0OZCc5RivuQrwpb6v4D4lKRGD5Ut2R2ownJzp8NRJ2BgfoeLttq5Cw5FOup76vssYgZDZD
API_PORT=8080
PORT=8080
```

### 6. Deploy
Click "Create Web Service"

### 7. URL
Tu URL: **https://whatsapp-api-xxxx.onrender.com**

## 💰 Plan Paid ($7/mes)
- App siempre activa
- 512 MB RAM
- Recomendado para WhatsApp (necesita estar 24/7)

## 🎯 Resultado
✅ Gratis pero se duerme
⚠️ Para WhatsApp mejor plan $7/mes

