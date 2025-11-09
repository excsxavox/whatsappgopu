# 📨 Estructura de Mensajes - WhatsApp Cloud API

## 🗄️ Esquema Completo en MongoDB

Tu aplicación ahora guarda mensajes con una estructura **profesional y escalable** en la colección `messages`:

```javascript
{
  // === IDENTIFICACIÓN ===
  "_id": "wamid.HBgLMTUyMTI3...",        // wamid de Meta (único)

  // === MULTI-TENANT Y ROUTING ===
  "tenant_id": "acme",                    // Opcional: para multi-tenant
  "instance_id": "123456789012345",       // WABA_PHONE_ID
  "channel": "whatsapp",                  // Siempre "whatsapp"
  "provider": "meta",                     // Siempre "meta"

  // === DIRECCIÓN Y CONVERSACIÓN ===
  "direction": "in",                      // "in" (entrante) | "out" (saliente)
  "conversation_id": "5939XXXXXXX@123456789012345",  // phone@instance

  // Ventana de 24h de Meta (solo si aplica)
  "wa_conversation": {
    "id": "b908e7a3-abc123",              // Desde statuses
    "category": "service",                // marketing|utility|authentication|service
    "origin": "user_initiated",           // user_initiated|business_initiated
    "expires_at": { "$date": "2025-10-24T18:45:00Z" }
  },

  // === PARTICIPANTES ===
  "from": "5939XXXXXXX",                  // E.164 del cliente
  "to": "123456789012345",                // Tu WABA number
  "contact_id": "contacts/abc123",        // Referencia opcional a contacto

  // === CONTENIDO DEL MENSAJE ===
  "message": {
    "id": "wamid.HBgLMTUyMTI3...",        // wamid (mismo que _id)
    "type": "text",                       // Tipo del mensaje ⬇️

    // TEXTO (type: text)
    "text": {
      "body": "Hola 👋 ¿Cómo estás?"
    },

    // INTERACTIVO (type: interactive)
    "interactive": {
      "type": "button_reply",             // button_reply|list_reply|nfm_reply
      "button_reply": {
        "id": "BTN_AYUDA",
        "title": "Ayuda"
      }
    },

    // MEDIA (type: image|video|audio|document|sticker)
    "media": {
      "mime_type": "image/jpeg",
      "file_name": "foto.jpg",
      "sha256": "abc123...",
      "size": 123456,
      "caption": "Mi foto",
      // NO guardamos binario en Mongo - usamos storage externo
      "storage": {
        "provider": "s3",                 // s3|gcs|azure
        "bucket": "wa-media",
        "key": "2025/10/23/wamid.HBgM...jpg",
        "public_url": "https://s3.amazonaws.com/..."
      }
    },

    // UBICACIÓN (type: location)
    "location": {
      "latitude": -0.1807,
      "longitude": -78.4678,
      "name": "Quito",
      "address": "Ecuador"
    },

    // CONTEXTO (reply_to)
    "context": {
      "message_id": "wamid.original...",  // wamid del mensaje al que responde
      "from": "5939XXXXXXX"
    }
  },

  // === ESTADO Y TRACKING ===
  "status": "delivered",                  // Estado actual

  // Historial completo de estados
  "status_history": [
    { "status": "queued",    "ts": { "$date": "2025-10-23T16:00:10Z" } },
    { "status": "sent",      "ts": { "$date": "2025-10-23T16:00:11Z" }, "provider_id": "wamid..." },
    { "status": "delivered", "ts": { "$date": "2025-10-23T16:00:12Z" } },
    { "status": "read",      "ts": { "$date": "2025-10-23T16:01:02Z" } }
  ],

  // Error (solo si status: failed)
  "error": {
    "code": 131000,
    "title": "Rate limit hit",
    "details": "Pair rate limit: 1 mensaje cada 6s"
  },

  // === FLOW ENGINE (opcional) ===
  "flow_state": {
    "flow_id": "onboarding",
    "version": 7,
    "step": "ask_email",
    "context": { "lang": "es", "lead_id": "..." }
  },

  // === IDEMPOTENCIA ===
  "dedup_key": "123456789012345|wamid.HBgLMTUyMTI3...",  // instance_id|wamid

  // === TRAZABILIDAD (raw minimal) ===
  "raw_min": {
    "entry_id": "...",
    "change_field": "messages",
    "metadata": { "display_phone_number": "123456789012345" }
  },

  // === TIMESTAMPS ===
  "timestamps": {
    "created_at":   { "$date": "2025-10-23T16:00:10Z" },  // Cuando se guardó
    "received_at":  { "$date": "2025-10-23T16:00:10Z" },  // in
    "queued_at":    { "$date": "2025-10-23T16:00:10Z" },  // out
    "sent_at":      { "$date": "2025-10-23T16:00:11Z" },
    "delivered_at": { "$date": "2025-10-23T16:00:12Z" },
    "read_at":      { "$date": "2025-10-23T16:01:02Z" },
    "updated_at":   { "$date": "2025-10-23T16:01:02Z" }
  }
}
```

---

## 📊 Tipos de Mensajes Soportados

### 1. **text** - Texto simple
```javascript
{
  "message": {
    "type": "text",
    "text": { "body": "Hola 👋" }
  }
}
```

### 2. **image** - Imagen
```javascript
{
  "message": {
    "type": "image",
    "media": {
      "mime_type": "image/jpeg",
      "caption": "Mi foto",
      "storage": { "public_url": "https://..." }
    }
  }
}
```

### 3. **video** - Video
```javascript
{
  "message": {
    "type": "video",
    "media": {
      "mime_type": "video/mp4",
      "caption": "Mi video"
    }
  }
}
```

### 4. **audio** - Audio/Voz
```javascript
{
  "message": {
    "type": "audio",
    "media": { "mime_type": "audio/ogg" }
  }
}
```

### 5. **document** - Documento
```javascript
{
  "message": {
    "type": "document",
    "media": {
      "mime_type": "application/pdf",
      "file_name": "factura.pdf"
    }
  }
}
```

### 6. **location** - Ubicación
```javascript
{
  "message": {
    "type": "location",
    "location": {
      "latitude": -0.1807,
      "longitude": -78.4678,
      "name": "Quito"
    }
  }
}
```

### 7. **interactive** - Botones/Listas
```javascript
{
  "message": {
    "type": "interactive",
    "interactive": {
      "type": "button_reply",
      "button_reply": {
        "id": "BTN_AYUDA",
        "title": "Ayuda"
      }
    }
  }
}
```

---

## 🔄 Estados del Mensaje

| Estado | Descripción | Aplica a |
|--------|-------------|----------|
| `queued` | En cola para envío | Saliente |
| `sent` | Enviado a Meta | Saliente |
| `delivered` | Entregado al destinatario | Saliente |
| `read` | Leído por el destinatario | Saliente |
| `received` | Recibido de usuario | Entrante |
| `failed` | Falló (ver campo `error`) | Ambos |

---

## 📈 Índices en MongoDB

Tu aplicación crea automáticamente estos índices para rendimiento óptimo:

```javascript
// Índice por conversación (queries principales)
{ "conversation_id": 1, "timestamps.created_at": -1 }

// Índice por dedup_key (idempotencia) - ÚNICO
{ "dedup_key": 1 } // unique: true

// Índice por instance (multi-instance)
{ "instance_id": 1, "timestamps.created_at": -1 }

// Índice por tenant (multi-tenant)
{ "tenant_id": 1, "timestamps.created_at": -1 }

// Índice por remitente
{ "from": 1, "timestamps.created_at": -1 }

// Índice por estado (reporting)
{ "status": 1, "timestamps.created_at": -1 }
```

---

## 🔍 Queries Típicas

### Obtener conversación completa
```javascript
db.messages.find({
  "conversation_id": "5939XXXXXXX@123456789012345"
}).sort({ "timestamps.created_at": -1 }).limit(50)
```

### Buscar por remitente
```javascript
db.messages.find({
  "from": "5939XXXXXXX"
}).sort({ "timestamps.created_at": -1 })
```

### Mensajes fallidos
```javascript
db.messages.find({
  "status": "failed"
}).sort({ "timestamps.created_at": -1 })
```

### Mensajes por instancia
```javascript
db.messages.find({
  "instance_id": "123456789012345",
  "direction": "out"
}).sort({ "timestamps.created_at": -1 })
```

### Verificar duplicado (idempotencia)
```javascript
db.messages.findOne({
  "dedup_key": "123456789012345|wamid.HBgLMTUyMTI3..."
})
```

---

## 💡 Buenas Prácticas

### 1. **Media Storage**
❌ **NO** guardes archivos binarios en MongoDB:
```javascript
// MAL
"media": { "binary_data": Buffer(...) }  // ❌
```

✅ **SÍ** usa storage externo (S3, GCS, Azure):
```javascript
// BIEN
"media": {
  "storage": {
    "provider": "s3",
    "public_url": "https://s3.amazonaws.com/..."
  }
}  // ✅
```

### 2. **Idempotencia**
Siempre verifica `dedup_key` antes de procesar:
```javascript
// Antes de guardar mensaje entrante
const exists = await db.messages.findOne({
  dedup_key: `${instanceID}|${wamid}`
});

if (exists) {
  return; // Ya procesado
}
```

### 3. **Timestamps**
Usa timestamps específicos para cada estado:
```javascript
"timestamps": {
  "sent_at":      ISODate("..."),  // Cuando se envió
  "delivered_at": ISODate("..."),  // Cuando se entregó
  "read_at":      ISODate("..."),  // Cuando se leyó
}
```

### 4. **Conversaciones**
Usa `conversation_id` para agrupar mensajes:
```javascript
// conversation_id = phone@instance
"conversation_id": "5939XXXXXXX@123456789012345"
```

### 5. **Multi-tenant** (opcional)
Si tienes múltiples clientes, usa `tenant_id`:
```javascript
"tenant_id": "empresa_abc",
"instance_id": "123456789012345"
```

---

## 🔐 Privacidad

### ⚠️ NO guardes datos sensibles en `raw_min`
```javascript
// Solo metadata mínima
"raw_min": {
  "entry_id": "...",
  "change_field": "messages"
  // ❌ NO guardes el JSON completo de Meta aquí
}
```

### ✅ Campos seguros
- ✅ `wamid`, `timestamps`, `status`
- ✅ Contenido del mensaje (ya lo eligió el usuario)
- ✅ Metadata de conversación

### ❌ NO guardar
- ❌ Tokens, secrets
- ❌ Payload completo de Meta
- ❌ Headers HTTP
- ❌ Datos de autenticación

---

## 📊 Métricas y Reporting

### Mensajes por estado (últimas 24h)
```javascript
db.messages.aggregate([
  {
    $match: {
      "timestamps.created_at": {
        $gte: new Date(Date.now() - 24*60*60*1000)
      }
    }
  },
  {
    $group: {
      _id: "$status",
      count: { $sum: 1 }
    }
  }
])
```

### Tiempo promedio de entrega
```javascript
db.messages.aggregate([
  {
    $match: {
      "status": "delivered",
      "timestamps.sent_at": { $exists: true },
      "timestamps.delivered_at": { $exists: true }
    }
  },
  {
    $project: {
      delivery_time: {
        $subtract: ["$timestamps.delivered_at", "$timestamps.sent_at"]
      }
    }
  },
  {
    $group: {
      _id: null,
      avg_delivery_ms: { $avg: "$delivery_time" }
    }
  }
])
```

---

## 🎯 Resumen

✅ **Tu aplicación ahora guarda**:
- ✅ Estructura completa y profesional
- ✅ Historial de estados
- ✅ Soporte para todos los tipos de mensajes
- ✅ Idempotencia robusta
- ✅ Multi-tenant ready
- ✅ Índices optimizados
- ✅ Trazabilidad completa
- ✅ Flow engine ready

**¡Listo para producción!** 🚀

