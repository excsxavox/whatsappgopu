# 🚀 Guía Rápida: Configurar WhatsApp Business Cloud API

## ✅ Lo que ya tienes configurado

- ✅ **MongoDB Atlas** - Conexión cloud configurada
- ✅ **Backend compilado** - `whatsapp-cloud-api.exe`
- ✅ **Arquitectura completa** - Hexagonal + CRUD empresas

## 🔧 Lo que necesitas configurar

### 1. Obtener Credenciales de Meta

#### Paso 1: Ir a Meta for Developers
🔗 https://developers.facebook.com/apps

#### Paso 2: Crear/Seleccionar App
1. Click en **"Create App"** o selecciona una existente
2. Tipo: **Business**
3. Nombre: `WhatsApp API - [Tu Nombre]`

#### Paso 3: Agregar WhatsApp
1. En el dashboard, click en **"Add Product"**
2. Busca **"WhatsApp"**
3. Click en **"Set Up"**

#### Paso 4: Obtener Credenciales

##### 📱 WABA_PHONE_ID
```
Dashboard → WhatsApp → API Setup
Copiar: "Phone number ID"
Ejemplo: 123456789012345
```

##### 🔑 WABA_TOKEN (Access Token)
```
Dashboard → WhatsApp → API Setup
Generar: "Permanent token" (recomendado)
O usar el "Temporary token" para pruebas (24 horas)
Ejemplo: EAAFj8xVZBZBs0BAxxxxxxxxxxxxxx
```

##### 🔐 WHATSAPP_APP_SECRET
```
Dashboard → Settings → Basic
Click en "Show" junto a "App secret"
Copiar el valor
Ejemplo: a1b2c3d4e5f6g7h8i9j0
```

##### ✅ WHATSAPP_VERIFY_TOKEN
```
Este lo CREAS TÚ (cualquier string seguro)
Ejemplo: mi_token_super_secreto_12345
Guárdalo, lo necesitarás para configurar el webhook
```

---

## ⚙️ Configurar Variables de Entorno

### Opción 1: Editar START.ps1

Abre `START.ps1` y reemplaza los valores:

```powershell
$env:WHATSAPP_VERIFY_TOKEN = "mi_token_super_secreto_12345"
$env:WHATSAPP_APP_SECRET = "a1b2c3d4e5f6g7h8i9j0"
$env:WABA_PHONE_ID = "123456789012345"
$env:WABA_TOKEN = "EAAFj8xVZBZBs0BAxxxxxxxxxxxxxx"
```

### Opción 2: Variables de Sistema (PowerShell)

```powershell
$env:WHATSAPP_VERIFY_TOKEN = "tu_token_aqui"
$env:WHATSAPP_APP_SECRET = "tu_secret_aqui"
$env:WABA_PHONE_ID = "tu_phone_id_aqui"
$env:WABA_TOKEN = "tu_access_token_aqui"
```

---

## 🚀 Iniciar el Servidor

### Con Script (Recomendado)
```powershell
.\START.ps1
```

### Manual
```powershell
# Configurar variables (ver arriba)
# Luego ejecutar:
go run cmd/server/main.go

# O usar el binario:
.\whatsapp-cloud-api.exe
```

---

## 🔔 Configurar Webhook en Meta

### 1. Exponer Servidor Público

#### Desarrollo (ngrok):
```bash
# Instalar ngrok: https://ngrok.com/download
ngrok http 8080

# Copiar la URL HTTPS que te da
# Ejemplo: https://abcd-1234-5678-9012.ngrok.io
```

#### Producción:
- Sube el servidor a un VPS/Cloud con dominio
- Ejemplo: https://whatsapp-api.tu-dominio.com

### 2. Configurar en Meta Dashboard

```
Dashboard → WhatsApp → Configuration
```

#### Editar Webhook:
- **Callback URL**: `https://tu-url.com/webhook`
- **Verify token**: El mismo que pusiste en `WHATSAPP_VERIFY_TOKEN`
- Click en **"Verify and save"**

#### Subscribe to webhooks:
Marcar:
- ✅ **messages** (REQUERIDO)

Click en **"Subscribe"**

---

## 🧪 Probar la Integración

### 1. Verificar que el servidor está corriendo

```bash
curl http://localhost:8080/health
```

Respuesta esperada:
```json
{
  "service": "whatsapp-api",
  "status": "ok"
}
```

### 2. Enviar mensaje de prueba

```bash
curl -X POST http://localhost:8080/send \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "5215512345678",
    "message": "¡Hola desde la API!"
  }'
```

### 3. Probar webhook

Envía un mensaje de WhatsApp al número configurado en Meta.
Deberías recibir una respuesta automática.

---

## 📊 Verificar MongoDB

Tu aplicación ya está guardando datos en MongoDB Atlas:

```javascript
// Base de datos: whatsapp_api
// Colecciones:
// - companies (empresas)
// - messages (mensajes)
// - sessions (sesiones)
```

Para ver los datos:
1. Ve a https://cloud.mongodb.com
2. Click en tu cluster
3. **Browse Collections**
4. Selecciona la base de datos `whatsapp_api`

---

## 🏢 API de Empresas

Ya está implementada y lista:

```http
GET    http://localhost:8080/api/companies
POST   http://localhost:8080/api/companies
GET    http://localhost:8080/api/companies/{id}
PUT    http://localhost:8080/api/companies/{id}
DELETE http://localhost:8080/api/companies/{id}
POST   http://localhost:8080/api/companies/{id}/activate
POST   http://localhost:8080/api/companies/{id}/deactivate
```

Consulta `API_EMPRESAS.md` para detalles.

---

## ⚠️ Troubleshooting

### Error: "Variable de entorno requerida: WABA_TOKEN"
**Solución**: Configura todas las variables de WhatsApp Cloud API

### Error: "Cloud API error: (status 401)"
**Solución**: Tu `WABA_TOKEN` es inválido o expiró. Genera uno nuevo.

### Error: "Cloud API error: (status 404)"
**Solución**: Tu `WABA_PHONE_ID` es incorrecto. Verifica en Meta Dashboard.

### Webhook no recibe eventos
**Solución**: 
1. Verifica que el webhook esté configurado en Meta
2. Verifica que el servidor esté público (ngrok)
3. Verifica que `WHATSAPP_VERIFY_TOKEN` coincida
4. Revisa los logs del servidor

### MongoDB no conecta
**Solución**: Ya está configurado con Atlas. Si hay error:
1. Verifica que la contraseña no tenga caracteres especiales sin encodear
2. Verifica el whitelist de IPs en MongoDB Atlas (permite 0.0.0.0/0)

---

## 📚 Documentación Adicional

- **README.md** - Documentación general
- **API_EMPRESAS.md** - API REST de empresas
- **SETUP.md** - Guía de instalación completa
- **ARCHITECTURE.md** - Arquitectura hexagonal

---

## ✅ Checklist Final

- [ ] Obtener credenciales de Meta (WABA_PHONE_ID, WABA_TOKEN, WHATSAPP_APP_SECRET)
- [ ] Configurar variables de entorno (editar START.ps1)
- [ ] Iniciar servidor (`.\START.ps1`)
- [ ] Exponer servidor público (ngrok o dominio)
- [ ] Configurar webhook en Meta Dashboard
- [ ] Enviar mensaje de prueba desde WhatsApp
- [ ] Verificar que recibe respuesta automática
- [ ] Verificar datos en MongoDB Atlas

---

**¡Listo!** Una vez completes estos pasos, tu sistema estará 100% funcional. 🎉

