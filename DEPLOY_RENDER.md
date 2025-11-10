# 🚀 DESPLEGAR EN RENDER.COM

## ⏱️ Tiempo estimado: 5 minutos

---

## 📋 PASO 1: CREAR CUENTA EN RENDER

1. Ve a: **https://render.com**
2. Click en **"Get Started"**
3. **Sign up with GitHub** (usa tu cuenta de GitHub)
4. Autoriza Render para acceder a tus repos

---

## 📦 PASO 2: CREAR WEB SERVICE

1. En el dashboard de Render, click en **"New +"**
2. Selecciona **"Web Service"**
3. Click en **"Connect a repository"**
4. Busca y selecciona: **`whatsappgo`**
5. Click en **"Connect"**

---

## ⚙️ PASO 3: CONFIGURAR EL SERVICIO

### Información básica:
- **Name:** `whatsapp-api-go` (o el que quieras)
- **Region:** Oregon (USA West) - el más cercano
- **Branch:** `main`
- **Runtime:** Detectará automáticamente **Docker** ✅

### Build & Deploy:
- ✅ **Dockerfile detectado automáticamente**
- ✅ No necesitas configurar Build Command ni Start Command

### Plan:
- Selecciona: **Free** ✅

---

## 🔐 PASO 4: AGREGAR VARIABLES DE ENTORNO

**Antes de hacer deploy**, scroll down hasta **"Environment Variables"** y agrega estas variables:

Click en **"Add Environment Variable"** para cada una:

```
MONGODB_URL=mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0

MONGO_DB=whatsapp_api

WHATSAPP_VERIFY_TOKEN=mi_token_seguro_whatsapp_2024

WHATSAPP_APP_SECRET=451614ef9eb9b35571dc352af6b2110e

WABA_PHONE_ID=804818756055720

WABA_TOKEN=EACO8kt4CNU4BP5ZAjnyEZBsatkJmx2XvPvYfO9cCllcjZANi1UeTvuh9LBWQ2t3Rse4B0q4rij37Ml3vgiFQB6krHWYhdW6mUkfRZBrA6w3ZBOZBYL1AAgTZCyS1Ls5zB4OwZAqPB2Dgpcz8Ucn2TjnPzGbVD3zza6IKmlGlYsaLSC3SNBHXvWjNj4W1FRPXtiY7y2ksG7n7xDzZBNe6kYTM3p0OZCc5RivuQrwpb6v4D4lKRGD5Ut2R2ownJzp8NRJ2BgfoeLttq5Cw5FOup76vssYgZDZD

API_PORT=8080

LOG_LEVEL=INFO
```

### ⚠️ Forma rápida:
Puedes copiar todas las variables de una vez usando el **"Add from .env"** y pegar el contenido de `VARIABLES_RAILWAY.txt`

---

## 🎯 PASO 5: DEPLOY

1. Una vez agregadas todas las variables, click en **"Create Web Service"**
2. Render comenzará a construir tu aplicación (toma ~3-5 minutos)
3. Verás los logs en tiempo real

### ✅ Indicadores de éxito:
```
==> Building...
==> Pushing...
==> Deploying...
==> Your service is live 🎉
```

---

## 🌐 PASO 6: OBTENER TU URL

1. Una vez deployado, verás tu URL en la parte superior:
   ```
   https://whatsapp-api-go-XXXX.onrender.com
   ```

2. **COPIA ESA URL** (la necesitarás para el webhook)

---

## 🔗 PASO 7: CONFIGURAR WEBHOOK EN META

1. Ve a: **https://developers.facebook.com/apps/10058963160806734**
2. En el menú izquierdo: **WhatsApp → Configuration**
3. Click en **"Edit"** en Webhook
4. Configura:
   - **Callback URL:** `https://TU-URL-RENDER.onrender.com/webhook`
   - **Verify Token:** `mi_token_seguro_whatsapp_2024`
5. Click en **"Verify and Save"**
6. Subscribe to: **messages** ✅

---

## 🧪 PASO 8: PROBAR

### Verificar que está funcionando:
1. Abre en tu navegador:
   ```
   https://TU-URL-RENDER.onrender.com/health
   ```
   Deberías ver:
   ```json
   {"status":"ok"}
   ```

### Probar webhook:
1. Envía un WhatsApp al número de prueba de tu app
2. Ve los logs en Render:
   - Click en tu servicio
   - Tab **"Logs"**
   - Deberías ver: `📥 Webhook recibido`

---

## ⚠️ IMPORTANTE: FREE TIER

### Limitaciones del plan gratuito:
- ✅ **750 horas/mes** (suficiente para 24/7)
- ⚠️ **Duerme después de 15 minutos sin actividad**
- ⚠️ **Tarda ~30 segundos en despertar** al recibir el primer request

### ¿Cómo afecta esto?
- El **primer mensaje** después de 15 min de inactividad puede tardar en procesarse
- Los mensajes subsecuentes son instantáneos
- Para producción, considera el plan pagado ($7/mes) que **no duerme**

### Mantener despierta la app (opcional):
Puedes usar un servicio de "ping" gratuito como:
- **UptimeRobot** (https://uptimerobot.com)
- **Cron-job.org** (https://cron-job.org)

Configura un ping cada 10 minutos a tu URL `/health`

---

## 🎉 ¡LISTO!

Tu aplicación está desplegada y funcionando en Render.

### URLs importantes:
- **App:** https://TU-URL-RENDER.onrender.com
- **Health:** https://TU-URL-RENDER.onrender.com/health
- **Webhook:** https://TU-URL-RENDER.onrender.com/webhook

### Dashboard Render:
- **Logs:** Para ver qué está pasando
- **Metrics:** CPU, Memoria, Requests
- **Environment:** Para cambiar variables

---

## 🆘 PROBLEMAS COMUNES

### 1. "Build failed"
- Verifica que todas las variables estén configuradas
- Revisa los logs del build

### 2. "Application failed to respond"
- La app puede estar durmiendo (espera 30s)
- Verifica el PORT (debe ser 8080)

### 3. "Webhook verification failed"
- Verifica que el VERIFY_TOKEN coincida exactamente
- La URL debe terminar en `/webhook` (sin `/` adicional)

---

## 📞 SOPORTE

¿Problemas? Dime en qué paso estás y te ayudo.
