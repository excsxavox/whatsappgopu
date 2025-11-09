# 🚀 Setup Completo - WhatsApp Business Cloud API + MongoDB

## ✅ Stack Tecnológico

- **Backend**: Go 1.24
- **Arquitectura**: Hexagonal (Ports & Adapters)
- **Base de Datos**: MongoDB
- **API**: WhatsApp Business Cloud API (Meta)
- **HTTP**: API REST + Webhooks

## 📋 Requisitos Previos

### 1. Instalar Go

```bash
# Descargar desde: https://go.dev/dl/
# Versión mínima: 1.21+
```

### 2. Instalar MongoDB

**Opción A: Docker (Recomendado)**
```bash
docker run -d \
  --name mongodb \
  -p 27017:27017 \
  -v mongodb_data:/data/db \
  mongo:latest
```

**Opción B: Instalación Local**
- Windows: https://www.mongodb.com/try/download/community
- Linux: `sudo apt install mongodb`
- macOS: `brew install mongodb-community`

**Opción C: MongoDB Atlas (Cloud)**
- Crear cuenta en: https://cloud.mongodb.com
- Tier gratis: 512MB
- Obtener connection string

### 3. Crear App en Meta for Developers

1. Ve a: https://developers.facebook.com/apps
2. Crear App → **Business**
3. Agregar producto: **WhatsApp**
4. Completar configuración inicial

## 🔑 Obtener Credenciales

### WhatsApp Business API

1. **Phone Number ID** (`WABA_PHONE_ID`)
   - Dashboard → WhatsApp → API Setup
   - Copiar el "Phone number ID"

2. **Access Token** (`WABA_TOKEN`)
   - Dashboard → WhatsApp → API Setup
   - Generar "Permanent token" (recomendado)
   - O usar el temporal para pruebas

3. **App Secret** (`WHATSAPP_APP_SECRET`)
   - Dashboard → Settings → Basic
   - Click en "Show" junto a "App secret"

4. **Verify Token** (`WHATSAPP_VERIFY_TOKEN`)
   - Crear un string seguro aleatorio
   - Ejemplo: `my_super_secret_token_123`

## ⚙️ Configuración

### 1. Clonar Variables de Entorno

```bash
cp config.env.example .env
```

### 2. Editar `.env`

```bash
# MongoDB
MONGO_URI=mongodb://localhost:27017
MONGO_DB=whatsapp_api

# WhatsApp Cloud API
WHATSAPP_VERIFY_TOKEN=my_super_secret_token_123
WHATSAPP_APP_SECRET=tu_app_secret_aqui
WABA_PHONE_ID=123456789012345
WABA_TOKEN=EAAxxxxxxxxxxxxxxxxxxxxx

# Opcional
API_PORT=8080
LOG_LEVEL=INFO
```

### 3. Instalar Dependencias

```bash
go mod tidy
```

## 🚀 Ejecución

### Desarrollo (local)

```bash
# Windows
.\RUN.bat

# Linux/Mac
go run cmd/server/main.go
```

### Producción (compilado)

```bash
# Compilar
go build -o whatsapp-cloud-api.exe cmd/server/main.go

# Ejecutar
.\whatsapp-cloud-api.exe
```

## 🔔 Configurar Webhook en Meta

### 1. Exponer Servidor Público

**Opción A: ngrok (para desarrollo)**
```bash
ngrok http 8080
# Copiar la URL: https://xxxx.ngrok.io
```

**Opción B: Servidor con dominio (producción)**
- Subir a servidor con dominio y SSL
- Ejemplo: https://tu-dominio.com

### 2. Configurar en Meta Dashboard

1. WhatsApp → Configuration
2. Edit webhook:
   - **Callback URL**: `https://tu-url.com/webhook`
   - **Verify token**: El mismo que pusiste en `WHATSAPP_VERIFY_TOKEN`
   - Click en "Verify and save"

3. Subscribe to webhooks:
   - Marcar: ✅ **messages**
   - Click en "Subscribe"

### 3. Probar Webhook

```bash
# Enviar mensaje de prueba al número de WhatsApp
# El servidor debe recibir y responder automáticamente
```

## 🧪 Pruebas

### Health Check

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

### Enviar Mensaje

```bash
curl -X POST http://localhost:8080/send \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "5215512345678",
    "message": "¡Hola desde la API!"
  }'
```

### Verificar MongoDB

```bash
# Conectar a MongoDB
mongosh

# Usar base de datos
use whatsapp_api

# Ver colecciones
show collections

# Ver mensajes
db.messages.find().pretty()
```

## 📊 Colecciones MongoDB

### `messages`
```javascript
{
  "_id": "wamid.HBgLMTUyMTI...",
  "from": "",
  "to": "5215512345678",
  "content": "Hola",
  "message_type": "text",
  "timestamp": ISODate("2024-01-01T10:00:00Z"),
  "status": "sent",
  "media_url": ""
}
```

### `sessions`
```javascript
{
  "_id": "session_123",
  "phone_number": "123456789",
  "is_active": true,
  "is_connected": true,
  "connected_at": ISODate("2024-01-01T10:00:00Z"),
  "last_seen": ISODate("2024-01-01T10:00:00Z")
}
```

## 🔐 Seguridad en Producción

### Checklist

- [ ] **MongoDB**
  - [ ] Habilitar autenticación
  - [ ] Crear usuario con permisos limitados
  - [ ] Firewall solo desde el servidor

- [ ] **HTTPS**
  - [ ] Usar SSL/TLS (Let's Encrypt, Caddy)
  - [ ] Nunca HTTP en producción

- [ ] **Tokens**
  - [ ] Usar tokens permanentes de producción
  - [ ] No compartir `.env`
  - [ ] Rotar tokens periódicamente

- [ ] **Aplicación**
  - [ ] Agregar autenticación en `/send`
  - [ ] Implementar rate limiting global
  - [ ] Logs estructurados
  - [ ] Monitoreo (Prometheus, Grafana)

### Ejemplo: MongoDB con Autenticación

```bash
# Crear usuario admin
mongosh
use admin
db.createUser({
  user: "waba_admin",
  pwd: "password_seguro",
  roles: [{ role: "readWrite", db: "whatsapp_api" }]
})

# Actualizar .env
MONGO_URI=mongodb://waba_admin:password_seguro@localhost:27017
```

## 🐛 Troubleshooting

### MongoDB no conecta

```bash
# Verificar que está corriendo
docker ps | grep mongodb
# O
mongosh --eval "db.version()"

# Iniciar MongoDB
docker start mongodb
# O
systemctl start mongodb
```

### Error: "Invalid signature"

- Verifica que `WHATSAPP_APP_SECRET` sea correcto
- Revisa los logs del servidor

### Mensajes no llegan

1. Verifica webhook configurado en Meta
2. Verifica que el servidor esté público (ngrok)
3. Revisa logs: `LOG_LEVEL=DEBUG`

### Rate limit 429

- Estás enviando muy rápido al mismo usuario
- Límite: 1 mensaje cada 6 segundos por usuario
- El servidor lo maneja automáticamente

## 📚 Recursos

- [Documentación WhatsApp Cloud API](https://developers.facebook.com/docs/whatsapp/cloud-api)
- [MongoDB Go Driver](https://www.mongodb.com/docs/drivers/go/current/)
- [Arquitectura Hexagonal](ARCHITECTURE.md)

## 🎯 Próximos Pasos

1. **Testing**
   - Implementar tests unitarios
   - Tests de integración con MongoDB
   - Mock de Cloud API

2. **Features**
   - Soporte para imágenes/multimedia
   - Templates de mensajes
   - Botones interactivos
   - Listas

3. **Infraestructura**
   - CI/CD (GitHub Actions)
   - Docker Compose
   - Kubernetes deployment

4. **Observabilidad**
   - Métricas (Prometheus)
   - Logs centralizados (ELK)
   - Tracing (OpenTelemetry)

---

**¡Tu aplicación está lista para producción!** 🎉

