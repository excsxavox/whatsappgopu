# WhatsApp Business Cloud API - Hexagonal Architecture 🎯

> Servidor API REST para **WhatsApp Business Cloud API** (Meta) con Arquitectura Hexagonal pura en Go + MongoDB

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev)
[![Architecture](https://img.shields.io/badge/Architecture-Hexagonal-blue?style=flat)](ARCHITECTURE.md)
[![Meta Cloud API](https://img.shields.io/badge/Meta-Cloud_API-0668E1?style=flat)](https://developers.facebook.com/docs/whatsapp/cloud-api)
[![MongoDB](https://img.shields.io/badge/MongoDB-Database-47A248?style=flat&logo=mongodb)](https://www.mongodb.com)

## 🌟 Características

- ✅ **WhatsApp Business Cloud API** (Meta)
- ✅ **Arquitectura Hexagonal pura** (Ports & Adapters)
- ✅ **MongoDB** como base de datos
- ✅ **Webhooks de Meta** con validación de firma
- ✅ **Idempotencia** con wamid (sin duplicados)
- ✅ **Rate limiting** automático (pair rate limit)
- ✅ **API REST** para envío de mensajes
- ✅ **Sin sesión local** (100% Cloud)

## 🚀 Inicio Rápido

### 1. Instalar MongoDB

```bash
# Opción 1: Docker (recomendado)
docker run -d -p 27017:27017 --name mongodb mongo:latest

# Opción 2: Instalación local
# Descarga desde: https://www.mongodb.com/try/download/community

# Opción 3: MongoDB Atlas (cloud)
# https://cloud.mongodb.com (gratis hasta 512MB)
```

### 2. Configurar Variables de Entorno

Copia `config.env.example` a `.env` y configura:

```bash
# MongoDB - REQUERIDO
MONGO_URI=mongodb://localhost:27017
MONGO_DB=whatsapp_api

# WhatsApp Cloud API - REQUERIDOS
WHATSAPP_VERIFY_TOKEN=mi_token_seguro_12345
WHATSAPP_APP_SECRET=tu_app_secret_de_meta
WABA_PHONE_ID=tu_phone_number_id
WABA_TOKEN=tu_access_token_permanente
```

### 3. Obtener Credenciales de Meta

1. Ve a [Meta for Developers](https://developers.facebook.com/apps)
2. Crea una app y agrega el producto "WhatsApp"
3. Obtén:
   - **WABA_PHONE_ID**: En "WhatsApp > API Setup"
   - **WABA_TOKEN**: Token de acceso permanente
   - **WHATSAPP_APP_SECRET**: En "Settings > Basic > App Secret"
   - **WHATSAPP_VERIFY_TOKEN**: Créalo tú (cualquier string seguro)

### 4. Ejecutar el Servidor

```bash
# Instalar dependencias
go mod tidy

# Windows
RUN.bat

# O directamente con Go
go run cmd/server/main.go
```

### 5. Configurar Webhook en Meta

1. Ve a tu app en Meta for Developers
2. "WhatsApp > Configuration"
3. Configura el webhook:
   - **URL**: `https://tu-dominio.com/webhook`
   - **Verify Token**: El mismo que pusiste en `WHATSAPP_VERIFY_TOKEN`
   - **Subscribirse a**: `messages`

## 📡 API Endpoints

### REST API

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/status` | Estado de conexión |
| POST | `/send` | Enviar mensaje |

### Webhooks Meta

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | `/webhook` | Verificación de webhook (Meta) |
| POST | `/webhook` | Recepción de eventos (Meta) |

## 💬 Ejemplos de Uso

### Enviar Mensaje

```bash
curl -X POST http://localhost:8080/send \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "5215512345678",
    "message": "¡Hola desde Cloud API!"
  }'
```

Respuesta:
```json
{
  "status": "success",
  "message": "Mensaje enviado correctamente",
  "message_id": "wamid.HBgLMTUyMTI..."
}
```

### Verificar Estado

```bash
curl http://localhost:8080/status
```

Respuesta:
```json
{
  "connected": true,
  "logged_in": true
}
```

## 🔔 Webhooks de Meta

### Evento de Mensaje Entrante

Cuando un usuario te envía un mensaje, Meta envía:

```json
{
  "object": "whatsapp_business_account",
  "entry": [{
    "changes": [{
      "value": {
        "messages": [{
          "from": "5215512345678",
          "id": "wamid.HBgLMTUyMTI...",
          "timestamp": "1623456789",
          "type": "text",
          "text": {
            "body": "Hola"
          }
        }]
      }
    }]
  }]
}
```

El servidor:
- ✅ Valida la firma `X-Hub-Signature-256`
- ✅ Deduplica por `wamid` (idempotencia)
- ✅ Responde automáticamente
- ✅ Ignora `statuses` (sin loops)

## 🏗️ Arquitectura Hexagonal

```
whatsapp-api-go/
├── cmd/server/              # 🎯 Entry Point + DI
│   └── main.go
│
├── internal/
│   ├── domain/             # 🧠 DOMINIO (lógica pura)
│   │   ├── entities/       # Message, Session, Connection
│   │   └── ports/          # Interfaces (contratos)
│   │
│   ├── application/        # 🎬 CASOS DE USO
│   │   └── usecases/       # Orquestación
│   │
│   └── infrastructure/     # 🔧 ADAPTADORES
│       ├── adapters/
│       │   ├── whatsapp/   # Cloud API (HTTP client)
│       │   ├── http/       # REST API + Webhooks
│       │   └── storage/    # Persistencia
│       └── config/
│
└── pkg/                    # 🛠️ UTILIDADES
    └── logger/
```

### Flujo de Webhook

```
Meta Cloud API
     ↓
POST /webhook (validación de firma)
     ↓
WebhookHandler (idempotencia wamid)
     ↓
SendMessageUseCase
     ↓
CloudAPIAdapter
     ↓
graph.facebook.com (envío)
```

## ⚡ Características Avanzadas

### 1. **Validación de Firma HMAC-SHA256**

Todas las peticiones de Meta incluyen `X-Hub-Signature-256`. El servidor valida:

```go
sha256=<hmac_hex_del_body>
```

### 2. **Idempotencia con wamid**

Cada mensaje tiene un `wamid` único (ej: `wamid.HBgLMTUyMTI...`). El servidor:
- Guarda wamids vistos (TTL 1 hora)
- Ignora duplicados automáticamente

### 3. **Pair Rate Limiting**

Meta recomienda **1 mensaje cada 6 segundos** por usuario. El servidor:
- Throttling automático por destinatario
- No excede el límite

### 4. **Sin Loops**

El servidor:
- ✅ Responde solo a `messages[]`
- ❌ Ignora `statuses[]` (enviados por nosotros)

## 🔧 Configuración

### Variables de Entorno

| Variable | Descripción | Requerida |
|----------|-------------|-----------|
| `MONGO_URI` | URI de conexión a MongoDB | ✅ |
| `MONGO_DB` | Nombre de la base de datos (default: whatsapp_api) | ❌ |
| `WHATSAPP_VERIFY_TOKEN` | Token para verificar webhook | ✅ |
| `WHATSAPP_APP_SECRET` | App Secret de Meta | ✅ |
| `WABA_PHONE_ID` | ID del número de teléfono | ✅ |
| `WABA_TOKEN` | Access token permanente | ✅ |
| `WABA_API_VERSION` | Versión de API (default: v20.0) | ❌ |
| `API_PORT` | Puerto del servidor (default: 8080) | ❌ |
| `LOG_LEVEL` | Nivel de logs (default: INFO) | ❌ |

### Configurar en PowerShell

```powershell
$env:MONGO_URI = "mongodb://localhost:27017"
$env:WHATSAPP_VERIFY_TOKEN = "mi_token_123"
$env:WHATSAPP_APP_SECRET = "app_secret_meta"
$env:WABA_PHONE_ID = "123456789"
$env:WABA_TOKEN = "EAAx..."
go run cmd/server/main.go
```

## 📊 Límites de WhatsApp

### Messaging Limits

Tu número tiene límites diarios de conversaciones:

- **Tier 1**: 1,000 conversaciones únicas / 24h
- **Tier 2**: 10,000 conversaciones únicas / 24h
- **Tier 3**: 100,000 conversaciones únicas / 24h
- **Tier 4+**: Ilimitado

Revisa tu tier en: WhatsApp Manager > Insights

### Pair Rate Limit

- **1 mensaje cada 6 segundos** por destinatario
- El servidor lo maneja automáticamente

## 🔐 Seguridad en Producción

### ⚠️ Checklist

- [ ] Usa HTTPS (Caddy, nginx con SSL)
- [ ] No expongas directamente el servidor
- [ ] Valida siempre la firma `X-Hub-Signature-256`
- [ ] Usa tokens de producción (no de prueba)
- [ ] Implementa rate limiting global
- [ ] Implementa autenticación en endpoints `/send`
- [ ] Monitorea logs y métricas
- [ ] Usa Redis para idempotencia en producción

### Reverse Proxy (nginx)

```nginx
server {
    listen 443 ssl;
    server_name tu-dominio.com;
    
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

## 🧪 Testing

```bash
# Verificar que el servidor está corriendo
.\test_server.ps1
```

## 🛠️ Compilar

```bash
# Compilar binario
go build -o whatsapp-cloud-api.exe cmd/server/main.go

# Ejecutar
.\whatsapp-cloud-api.exe
```

## 🐛 Troubleshooting

### Error: "Variable de entorno requerida: MONGO_URI"

Verifica que MongoDB esté instalado y corriendo:
```bash
# Verificar MongoDB local
mongosh --eval "db.version()"

# O iniciar con Docker
docker start mongodb
```

### Error: "Variable de entorno requerida: WABA_TOKEN"

Configura todas las variables requeridas en `.env` o como variables de entorno.

### Webhook no recibe eventos

1. Verifica que la URL esté configurada en Meta
2. Verifica que el `Verify Token` coincida
3. Revisa los logs del servidor
4. Prueba con `curl` local primero

### Error 401 en webhook

La firma `X-Hub-Signature-256` no coincide. Verifica que `WHATSAPP_APP_SECRET` sea correcto.

### Rate limit 429

Estás enviando más de 1 mensaje cada 6 segundos al mismo usuario. El servidor lo maneja automáticamente.

## 📖 Recursos

- [WhatsApp Cloud API - Documentación Oficial](https://developers.facebook.com/docs/whatsapp/cloud-api)
- [Webhooks - Meta Docs](https://developers.facebook.com/docs/graph-api/webhooks)
- [Messaging Limits](https://developers.facebook.com/docs/whatsapp/messaging-limits)
- [Hexagonal Architecture](ARCHITECTURE.md)

## 🤝 Integración

Este servidor puede integrarse con:

- ✅ Aplicaciones web/móviles
- ✅ Sistemas CRM (Salesforce, HubSpot)
- ✅ Automatización (n8n, Zapier, Make)
- ✅ Bots con IA (OpenAI, Anthropic)
- ✅ E-commerce (notificaciones de pedidos)

## ✨ Ventajas vs whatsmeow

| Característica | whatsmeow | Cloud API |
|---------------|-----------|-----------|
| Tipo | WhatsApp Web | Oficial Meta |
| Sesión | Local (QR) | Sin sesión |
| Estabilidad | Variable | Alta |
| Límites | ∞ | Tier-based |
| Costo | Gratis | Por conversación |
| Soporte | Comunidad | Meta oficial |
| Producción | No recomendado | ✅ Recomendado |

## 📝 Licencia

MIT License

---

**Desarrollado con ❤️ usando WhatsApp Business Cloud API + Arquitectura Hexagonal**
