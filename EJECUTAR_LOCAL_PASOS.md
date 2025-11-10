# 🏠 Ejecutar Aplicación Local - Pasos

## ✅ Ya está listo:
- ✅ Archivo .env creado
- ✅ Aplicación compilada (whatsapp-api.exe)

---

## 📋 PASO 1: Ejecutar la aplicación

Abre una terminal PowerShell en esta carpeta y ejecuta:

```powershell
.\whatsapp-api.exe
```

Deberías ver:
```
[INFO] ✅ Conectado a MongoDB
[INFO] 🚀 Servidor iniciado en :8080
```

**⚠️ DEJA ESTA TERMINAL ABIERTA** - Aquí verás los logs en tiempo real.

---

## 📋 PASO 2: Ejecutar ngrok

Abre **OTRA terminal** PowerShell y ejecuta:

```powershell
ngrok http 8080
```

Verás algo como:
```
Forwarding   https://abc123def456.ngrok.io -> http://localhost:8080
```

**Copia la URL** que empieza con `https://` (ejemplo: `https://abc123def456.ngrok.io`)

---

## 📋 PASO 3: Configurar webhook en Meta

1. Ve a: https://developers.facebook.com/apps/10058963160806734
2. WhatsApp → **Configuration** 
3. **Webhook** → Click en **Edit**
4. Pega tu URL de ngrok + `/webhook`:
   ```
   https://TU-URL-NGROK.ngrok.io/webhook
   ```
5. Verify Token: `mi_token_secreto_123`
6. Click en **Verify and Save**
7. Subscribe to: **messages**

---

## 📋 PASO 4: ¡PROBAR!

Ahora envía un mensaje a tu WhatsApp.

**En la primera terminal** (donde ejecutaste `whatsapp-api.exe`) verás los logs en tiempo real:

```
[INFO] 📨 Mensaje entrante [from 593992686734 wamid xxx type text]
[INFO] Processing message in flow: 593992686734@804818756055720
[INFO] Starting flow for conversation: 593992686734@804818756055720
[INFO] Processing node node_1_bienvenida (type: TEXT)
[INFO] ✅ Mensaje enviado exitosamente
```

---

## 🎯 ¿Qué buscar en los logs?

✅ **SI VES ESTO** → Todo está bien:
```
[INFO] Processing TEXT node: node_1_bienvenida
[INFO] Mensaje enviado exitosamente
```

❌ **SI VES ESTO** → Hay un problema:
```
[ERROR] (#131030) Recipient phone number not in allowed list
```

---

## 🔄 DESPUÉS DE PROBAR

Cuando termines de probar localmente:

1. **Detén** whatsapp-api.exe (Ctrl+C)
2. **Detén** ngrok (Ctrl+C)
3. **Vuelve a configurar el webhook en Meta** apuntando a Azure:
   ```
   https://whatsapp-api-go-dpb5cgbnaec2gdf2.eastus-01.azurewebsites.net/webhook
   ```

---

## 🆘 ¿Necesitas ayuda?

Si algo sale mal, comparte los logs que veas en la terminal y te ayudo a debuggear.

