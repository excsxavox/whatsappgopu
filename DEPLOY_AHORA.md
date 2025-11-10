# 🚀 Deploy en Railway - Guía Rápida

## ✅ Ya tienes el código en GitHub
Repositorio: https://github.com/excsxavox/whatsappgo

---

## 📋 PASO 1: Conectar Railway (2 minutos)

### 1.1 Ir a Railway
🔗 **https://railway.app**

### 1.2 Login
- Click **"Login with GitHub"**
- Autoriza Railway

### 1.3 Crear Proyecto
- Click **"New Project"**
- Selecciona **"Deploy from GitHub repo"**
- Busca y click en **"whatsappgo"**

### 1.4 Railway empezará a compilar
✅ Verás logs de compilación (toma ~3-5 minutos)

---

## 📋 PASO 2: Configurar Variables (2 minutos)

Mientras compila, configura las variables:

### 2.1 Click en tu servicio
En el dashboard, click en el card que dice **"whatsappgo"**

### 2.2 Ir a Variables
Click en la pestaña **"Variables"** o **"Settings"** → **"Variables"**

### 2.3 Agregar Variables

**Opción A: Una por una**
Click "New Variable" y agrega cada una:

**Opción B: RAW Editor** (más rápido)
Click en "RAW Editor" y pega todo esto:

```env
MONGODB_URL=mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0
MONGO_DB=whatsapp_api
WHATSAPP_VERIFY_TOKEN=mi_token_seguro_whatsapp_2024
WHATSAPP_APP_SECRET=451614ef9eb9b35571dc352af6b2110e
WABA_PHONE_ID=804818756055720
WABA_TOKEN=EACO8kt4CNU4BP5ZAjnyEZBsatkJmx2XvPvYfO9cCllcjZANi1UeTvuh9LBWQ2t3Rse4B0q4rij37Ml3vgiFQB6krHWYhdW6mUkfRZBrA6w3ZBOZBYL1AAgTZCyS1Ls5zB4OwZAqPB2Dgpcz8Ucn2TjnPzGbVD3zza6IKmlGlYsaLSC3SNBHXvWjNj4W1FRPXtiY7y2ksG7n7xDzZBNe6kYTM3p0OZCc5RivuQrwpb6v4D4lKRGD5Ut2R2ownJzp8NRJ2BgfoeLttq5Cw5FOup76vssYgZDZD
WABA_API_VERSION=v20.0
API_PORT=8080
PORT=8080
LOG_LEVEL=INFO
```

### 2.4 Guardar
Click **"Add"** o **"Save Changes"**

Railway hará **redeploy automático** con las nuevas variables.

---

## 📋 PASO 3: Obtener URL Pública (1 minuto)

### 3.1 Ir a Settings
En tu servicio, click en **"Settings"**

### 3.2 Generar Dominio
- Scroll hasta **"Networking"** o **"Domains"**
- Click **"Generate Domain"**

### 3.3 Copiar URL
Railway generará una URL como:
```
whatsappgo-production-xxxx.up.railway.app
```

**COPIA ESTA URL** ✅

---

## 📋 PASO 4: Configurar Webhook en Meta (2 minutos)

### 4.1 Ir a Meta Dashboard
🔗 **https://developers.facebook.com/apps/10058963160806734/whatsapp-business/wa-settings/**

### 4.2 Edit Webhook
- Click en **"Configuration"**
- Click en **"Edit"** (junto a Callback URL)

### 4.3 Configurar
```
Callback URL: https://TU-URL-RAILWAY.up.railway.app/webhook
Verify Token: mi_token_seguro_whatsapp_2024
```

Ejemplo:
```
https://whatsappgo-production-a1b2.up.railway.app/webhook
```

### 4.4 Verify and Save
- Click **"Verify and Save"**
- Debería aparecer ✅ verde

### 4.5 Subscribe
Marca la casilla:
- ✅ **messages**

---

## 📋 PASO 5: Probar (1 minuto)

### 5.1 Verificar que la app esté corriendo
En Railway, ve a **"Deployments"** → Debería estar ✅ verde

### 5.2 Ver logs
Click en **"View Logs"** para ver si hay errores

### 5.3 Probar webhook
Desde tu terminal local:

```powershell
curl "https://TU-URL-RAILWAY.up.railway.app/webhook?hub.mode=subscribe&hub.verify_token=mi_token_seguro_whatsapp_2024&hub.challenge=test123"
```

Debería devolver: `test123`

### 5.4 Enviar mensaje de prueba
Envía un mensaje de WhatsApp al número configurado (+1 555 152 6940)

Deberías ver el mensaje en los logs de Railway.

---

## ✅ CHECKLIST FINAL

- [ ] Railway conectado
- [ ] Variables configuradas
- [ ] Dominio generado
- [ ] Webhook configurado en Meta
- [ ] Webhook verificado (✅ verde)
- [ ] Messages subscribed
- [ ] Prueba de webhook exitosa
- [ ] App corriendo sin errores

---

## 🎉 ¡LISTO!

Tu aplicación está desplegada en Railway 24/7.

**URL de tu app:**
```
https://TU-URL.up.railway.app
```

**Webhook para Meta:**
```
https://TU-URL.up.railway.app/webhook
```

**Endpoints disponibles:**
- `GET /webhook` - Verificación
- `POST /webhook` - Recibir mensajes
- `POST /send` - Enviar mensajes
- `GET /api/companies` - API empresas

---

## 🔄 Actualizaciones Futuras

Cada vez que hagas cambios:

```bash
git add .
git commit -m "Tu mensaje"
git push
```

Railway detectará el push y **redesplegará automáticamente**. 🚀

---

## 🆘 Troubleshooting

### Build failed
- Revisa logs en Railway
- Verifica que el Dockerfile esté correcto

### App crashed
- Verifica variables de entorno
- Revisa logs: busca errores de MongoDB
- Verifica que todas las variables estén configuradas

### Webhook no funciona
- Verifica que el token coincida
- Prueba la URL manualmente con curl
- Revisa logs de Railway cuando Meta envíe el webhook

---

## 💰 Costos

Railway:
- **$5 gratis al mes**
- Después: ~$1-3/mes para esta app
- Monitorea en: Railway Dashboard → Usage

---

**¿Preguntas?** Lee los logs de Railway, ahí está toda la info.


