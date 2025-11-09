# 🚂 Desplegar en Railway - Instrucciones Completas

## ✅ Estado Actual
- [x] Código listo
- [x] Git inicializado
- [x] Primer commit realizado
- [ ] Subir a GitHub
- [ ] Conectar Railway
- [ ] Configurar variables
- [ ] Deploy

---

## 📋 Paso 4: Crear repositorio en GitHub

### Opción A: Desde la Web (MÁS FÁCIL)

1. **Ve a:** https://github.com/new

2. **Configuración:**
   ```
   Repository name: whatsapp-api-go
   Description: WhatsApp Business Cloud API with Hexagonal Architecture
   Privacy: Private (recomendado) o Public
   ```

3. **NO marques:** "Add README", "Add .gitignore", ni "Choose a license"
   (Ya los tenemos)

4. **Click:** "Create repository"

5. **Copia los comandos** que aparecen en "…or push an existing repository from the command line"
   
   Se verán así:
   ```bash
   git remote add origin https://github.com/TU_USUARIO/whatsapp-api-go.git
   git branch -M main
   git push -u origin main
   ```

6. **Ejecuta esos comandos aquí en tu terminal**

### Opción B: Con GitHub CLI (si la tienes instalada)

```bash
gh repo create whatsapp-api-go --private --source=. --remote=origin --push
```

---

## 📋 Paso 5: Conectar Railway

1. **Ve a:** https://railway.app

2. **Login con GitHub**
   - Click "Login with GitHub"
   - Autoriza Railway

3. **Crear nuevo proyecto**
   - Click "New Project"
   - Selecciona "Deploy from GitHub repo"
   - Busca "whatsapp-api-go"
   - Click en tu repositorio

4. **Railway detectará automáticamente el Dockerfile**
   ✅ Railway empezará a hacer el build

---

## 📋 Paso 6: Configurar Variables de Entorno

1. **En Railway, click en tu servicio** (whatsapp-api-go)

2. **Click en la pestaña "Variables"**

3. **Agregar todas estas variables:**

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

4. **Click "Add" para cada variable**

   O usa el modo "RAW Editor" para pegar todas a la vez.

---

## 📋 Paso 7: Generar Dominio Público

1. **En Railway, ve a "Settings"**

2. **Scroll hasta "Networking"**

3. **Click "Generate Domain"**

4. **Copia tu URL** (será algo como: `whatsapp-api-go-production-xxxx.up.railway.app`)

---

## 📋 Paso 8: Configurar Webhook en Meta

1. **Ve a:** https://developers.facebook.com/apps/10058963160806734/whatsapp-business/wa-settings/

2. **Configuration → Edit webhook**

3. **Configuración:**
   ```
   Callback URL: https://tu-app.up.railway.app/webhook
   Verify Token: mi_token_seguro_whatsapp_2024
   ```

4. **Click "Verify and Save"**

5. **Subscribe to fields:**
   - ✅ messages

---

## 📋 Paso 9: Verificar Deploy

### Ver logs:
En Railway → Click en tu servicio → Pestaña "Deployments" → Click en el último deploy

### Probar endpoints:
```bash
# Health check (desde tu terminal)
curl https://tu-app.up.railway.app/webhook?hub.mode=subscribe&hub.verify_token=mi_token_seguro_whatsapp_2024&hub.challenge=test

# Debería devolver: test
```

---

## 🔄 Actualizaciones Futuras

Cada vez que hagas cambios:

```bash
git add .
git commit -m "Descripción del cambio"
git push
```

Railway detectará el push y hará **deploy automático**. 🚀

---

## 💰 Costos

- **$5 gratis al mes** (suficiente para desarrollo)
- Después: ~$1-3/mes dependiendo del uso
- Monitorea en Railway Dashboard → Usage

---

## 🆘 Troubleshooting

### Error: "Build failed"
- Revisa los logs en Railway
- Verifica que el Dockerfile esté correcto

### Error: "Application crashed"
- Revisa las variables de entorno
- Verifica que MONGODB_URL esté correcta
- Ve a "Logs" en Railway

### Webhook no responde
- Verifica que el dominio esté generado
- Prueba la URL manualmente
- Revisa que WHATSAPP_VERIFY_TOKEN coincida

---

## ✅ Checklist Final

- [ ] Repositorio creado en GitHub
- [ ] Código subido
- [ ] Railway conectado
- [ ] Variables configuradas
- [ ] Dominio generado
- [ ] Webhook configurado en Meta
- [ ] Prueba de envío exitosa

---

¡Listo! Tu aplicación estará disponible en Railway 24/7. 🎉

