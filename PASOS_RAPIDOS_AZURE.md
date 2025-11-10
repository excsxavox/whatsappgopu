# ⚡ PASOS RÁPIDOS: AZURE APP SERVICE

## 🎯 RESUMEN EN 5 PASOS

---

## 📍 **PASO 1: CREAR APP SERVICE** (5 min)

1. Ve a: **https://portal.azure.com**
2. Busca **"App Services"** → Click **"+ Crear"** → **"Aplicación web"**
3. Configura:
   - **Grupo de recursos**: Crear nuevo → `whatsapp-rg`
   - **Nombre**: `whatsapp-api-go` (tu URL será: `whatsapp-api-go.azurewebsites.net`)
   - **Publicar**: **Contenedor Docker** (⚠️ NO "Código")
   - **Sistema**: **Linux**
   - **Plan**: **B1 Basic** ($13/mes)
4. Click **"Siguiente: Docker"**:
   - **Opciones**: Contenedor único
   - **Origen**: Docker Hub
   - **Imagen**: `alpine:latest` (temporal)
5. Click **"Revisar y crear"** → **"Crear"**
6. Espera 2-3 min → Click **"Ir al recurso"**

---

## 🔧 **PASO 2: CONFIGURAR VARIABLES** (3 min)

1. Menú izquierdo → **"Configuración"**
2. Click **"+ Nueva configuración"** para cada variable:

```
MONGODB_URL = mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0
MONGO_DB = whatsapp_api
WHATSAPP_VERIFY_TOKEN = mi_token_seguro_whatsapp_2024
WHATSAPP_APP_SECRET = 451614ef9eb9b35571dc352af6b2110e
WABA_PHONE_ID = 804818756055720
WABA_TOKEN = EACO8kt4CNU4BP5ZAjnyEZBsatkJmx2XvPvYfO9cCllcjZANi1UeTvuh9LBWQ2t3Rse4B0q4rij37Ml3vgiFQB6krHWYhdW6mUkfRZBrA6w3ZBOZBYL1AAgTZCyS1Ls5zB4OwZAqPB2Dgpcz8Ucn2TjnPzGbVD3zza6IKmlGlYsaLSC3SNBHXvWjNj4W1FRPXtiY7y2ksG7n7xDzZBNe6kYTM3p0OZCc5RivuQrwpb6v4D4lKRGD5Ut2R2ownJzp8NRJ2BgfoeLttq5Cw5FOup76vssYgZDZD
API_PORT = 8080
LOG_LEVEL = INFO
```

3. Click **"Guardar"** (arriba)
4. Pestaña **"Configuración general"** → **Comando de inicio**: `./whatsapp-api-server`
5. Click **"Guardar"**

---

## 🔗 **PASO 3: CONECTAR GITHUB** (2 min)

1. Menú izquierdo → **"Centro de implementación"**
2. **Origen**: Selecciona **"GitHub"**
3. Click **"Autorizar"** → Inicia sesión en GitHub
4. Selecciona:
   - **Organización**: `excsxavox`
   - **Repositorio**: `whatsappgopu`
   - **Rama**: `main`
5. Click **"Guardar"**
6. ⏳ Espera 5-10 min (primer deployment)

---

## 🌐 **PASO 4: OBTENER URL** (30 seg)

1. Menú izquierdo → **"Información general"**
2. Copia el **"Dominio predeterminado"**: `whatsapp-api-go.azurewebsites.net`
3. Prueba en navegador: `https://whatsapp-api-go.azurewebsites.net/health`
   - Debe responder: `{"status":"ok"}`

---

## 📱 **PASO 5: CONFIGURAR META WEBHOOK** (2 min)

1. Ve a: **https://developers.facebook.com/**
2. Tu app → **WhatsApp** → **Configuración**
3. **Webhook** → **Editar**:
   - **URL**: `https://whatsapp-api-go.azurewebsites.net/webhook`
   - **Token**: `mi_token_seguro_whatsapp_2024`
4. Click **"Verificar y guardar"**
5. Activa el campo **"messages"**
6. Click **"Guardar"**

---

## ✅ **¡LISTO!**

**Envía un mensaje a tu WhatsApp de prueba y debería responder.**

---

## 📊 **VER LOGS EN TIEMPO REAL**

Azure Portal → Tu App Service → **"Supervisión"** → **"Secuencia de registro"**

---

## 🆘 **SI ALGO FALLA**

### No verifica el webhook:
- MongoDB Atlas → Network Access → Add IP: `0.0.0.0/0`

### No responde mensajes:
- Ve a "Secuencia de registro" y revisa los logs
- Verifica que todas las variables estén correctas

---

## 💰 **COSTO: ~$13/mes** (Plan B1 Basic)

---

## 📚 **DOCUMENTACIÓN COMPLETA**

Si necesitas más detalles, abre: `GUIA_AZURE_APP_SERVICE.md`

