# 🚂 Deploy en Railway (FÁCIL)

## ✅ Ventajas
- $5 gratis al mes (suficiente para tu app)
- HTTPS automático
- Deploy desde GitHub (automático)
- Dashboard super fácil

## 📋 Pasos

### 1. Subir a GitHub (si no lo has hecho)
```bash
git init
git add .
git commit -m "WhatsApp API ready for deployment"
git remote add origin https://github.com/TU_USUARIO/whatsapp-api-go.git
git push -u origin main
```

### 2. Ir a Railway
```
https://railway.app
```

### 3. Crear cuenta
- Login con GitHub

### 4. Crear proyecto
- Click "New Project"
- "Deploy from GitHub repo"
- Selecciona tu repositorio

### 5. Configurar variables
En el dashboard, click en tu servicio → Variables:

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

### 6. Railway detecta Dockerfile automáticamente y hace deploy

### 7. Obtener URL
En Settings → Public URL → Generate Domain

Tu URL: **https://tuapp.railway.app**

## 🔄 Deploy automático
Cada `git push` → Deploy automático ✅

## 💰 Precio
- $5 gratis/mes
- $0.000231/GB-hora después
- Tu app: ~$1-2/mes

## 🎯 Resultado
✅ URL pública permanente
✅ Auto-deploy desde GitHub
✅ HTTPS incluido

