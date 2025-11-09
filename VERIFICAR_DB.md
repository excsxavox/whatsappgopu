# 🔍 Cómo Verificar MongoDB

## 📊 Estado de las Colecciones

### ⚠️ Importante
Las colecciones en MongoDB se crean **automáticamente** cuando:
1. Ejecutas la aplicación por primera vez
2. Se inserta el primer documento

**Los índices se crean al iniciar la aplicación.**

---

## 🎯 Verificación Paso a Paso

### Opción 1: MongoDB Atlas Web (Más Fácil)

1. **Ve a MongoDB Atlas**
   ```
   https://cloud.mongodb.com
   ```

2. **Inicia sesión** con tus credenciales

3. **Selecciona tu Cluster**
   - Click en **"Cluster0"** (o el nombre de tu cluster)

4. **Browse Collections**
   - Click en el botón **"Browse Collections"**

5. **Verifica la base de datos**
   - Busca: `whatsapp_api`
   - Deberías ver 3 colecciones:
     - ✅ `messages` (estructura nueva)
     - ✅ `companies` (empresas)
     - ✅ `sessions` (sesiones)

6. **Ver documentos**
   - Click en cada colección para ver los documentos guardados

---

### Opción 2: Script PowerShell

```powershell
.\verificar_mongodb.ps1
```

**Requisitos:**
- Tener `mongosh` instalado (MongoDB Shell)
- Descarga: https://www.mongodb.com/try/download/shell

---

### Opción 3: mongosh (Terminal)

```bash
# Conectar a MongoDB Atlas
mongosh "mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/"

# Usar la base de datos
use whatsapp_api

# Ver colecciones
show collections

# Contar documentos
db.messages.countDocuments()
db.companies.countDocuments()
db.sessions.countDocuments()

# Ver índices de messages
db.messages.getIndexes()

# Ver último mensaje
db.messages.findOne({}, {sort: {'timestamps.created_at': -1}})

# Ver todas las empresas
db.companies.find().pretty()
```

---

## 📋 Qué Esperar

### Antes de Ejecutar la Aplicación
```javascript
// Base de datos: whatsapp_api NO EXISTE aún
// O existe pero SIN colecciones
```

### Después de Ejecutar la Aplicación (Primera Vez)
```javascript
// Base de datos: whatsapp_api ✅
// Colecciones:
//   - messages (0 documentos, 7 índices) ✅
//   - companies (0 documentos, 2 índices) ✅
//   - sessions (0 documentos, 1 índice) ✅
```

### Después de Recibir Primer Mensaje
```javascript
// messages: 1 documento ✅
{
  "_id": "wamid.HBgLMTUyMTI3...",
  "direction": "in",
  "from": "5939XXXXXXX",
  "message": {
    "type": "text",
    "text": { "body": "Hola" }
  },
  "status": "received",
  "timestamps": { ... }
}
```

### Después de Crear Primera Empresa (API)
```javascript
// companies: 1 documento ✅
{
  "_id": "uuid-123",
  "code": "EMPRESA_001",
  "name": "Mi Primera Empresa",
  "is_active": true
}
```

---

## 🔎 Índices Creados Automáticamente

### Colección: `messages`
```javascript
// 1. _id (por defecto)
{ "_id": 1 }

// 2. conversation_id + timestamps
{ "conversation_id": 1, "timestamps.created_at": -1 }

// 3. dedup_key (ÚNICO) - Idempotencia
{ "dedup_key": 1 } // unique: true

// 4. instance_id + timestamps
{ "instance_id": 1, "timestamps.created_at": -1 }

// 5. tenant_id + timestamps (sparse)
{ "tenant_id": 1, "timestamps.created_at": -1 }

// 6. from + timestamps
{ "from": 1, "timestamps.created_at": -1 }

// 7. status + timestamps
{ "status": 1, "timestamps.created_at": -1 }
```

### Colección: `companies`
```javascript
// 1. _id (por defecto)
{ "_id": 1 }

// 2. code (ÚNICO)
{ "code": 1 } // unique: true

// 3. is_active
{ "is_active": 1 }
```

### Colección: `sessions`
```javascript
// 1. _id (por defecto)
{ "_id": 1 }

// 2. is_active
{ "is_active": 1 }
```

---

## ✅ Checklist de Verificación

### Al Iniciar la Aplicación
- [ ] MongoDB Atlas accesible
- [ ] Base de datos `whatsapp_api` creada
- [ ] 3 colecciones creadas (messages, companies, sessions)
- [ ] Todos los índices creados (ver logs: "✅ Índices creados")

### Logs Esperados
```
📊 Conectando a MongoDB...
✅ Conectado a MongoDB
📑 Creando índices...
✅ Índices creados
```

### Al Recibir Primer Webhook
```
📨 Mensaje entrante from=... wamid=... type=text
✅ Mensaje guardado en MongoDB wamid=...
```

### Al Crear Primera Empresa
```
🏢 POST /api/companies
✅ Empresa creada id=... code=...
```

---

## 🐛 Troubleshooting

### No se crean las colecciones
**Causa**: La aplicación no se ha ejecutado aún.
**Solución**: Ejecuta `.\START.ps1` o `go run cmd/server/main.go`

### Error: "Índices no se pueden crear"
**Causa**: Permisos insuficientes en MongoDB.
**Solución**: Verifica que el usuario tenga permisos de escritura en Atlas.

### Colecciones vacías
**Causa**: No se ha recibido ningún mensaje ni creado ninguna empresa.
**Solución**: Normal! Espera el primer webhook o crea una empresa vía API.

### No puedo ver las colecciones en Atlas
**Causa**: La base de datos no existe hasta que se inserta el primer documento.
**Solución**: Ejecuta la aplicación y espera el primer evento.

---

## 📝 Queries Útiles

### Ver todos los mensajes de una conversación
```javascript
db.messages.find({
  "conversation_id": "5939XXXXXXX@123456789012345"
}).sort({ "timestamps.created_at": -1 })
```

### Ver mensajes fallidos
```javascript
db.messages.find({ "status": "failed" })
```

### Contar mensajes por tipo
```javascript
db.messages.aggregate([
  { $group: { _id: "$message.type", count: { $sum: 1 } } }
])
```

### Ver empresas activas
```javascript
db.companies.find({ "is_active": true })
```

---

**¡Tu MongoDB está listo!** Solo falta ejecutar la aplicación. 🚀

