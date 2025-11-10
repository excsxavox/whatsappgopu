# 🚀 GUÍA PASO A PASO: AZURE APP SERVICE

## ✅ **PASO 1: ENTRAR AL PORTAL DE AZURE**

1. Abre tu navegador
2. Ve a: **https://portal.azure.com**
3. Inicia sesión con tu cuenta de Microsoft/Azure

---

## 📦 **PASO 2: CREAR APP SERVICE**

### 2.1. Buscar App Service

1. En la página principal de Azure, busca la barra de búsqueda (arriba)
2. Escribe: **"App Services"**
3. Click en **"App Services"**

### 2.2. Crear Nuevo

1. Click en el botón **"+ Crear"** o **"+ Create"** (arriba a la izquierda)
2. Click en **"Aplicación web"** o **"Web App"**

### 2.3. Configuración Básica (Pestaña "Básico")

Completa los siguientes campos:

#### **Detalles del proyecto:**
- **Suscripción**: Selecciona tu suscripción activa
- **Grupo de recursos**: 
  - Click en **"Crear nuevo"**
  - Nombre: `whatsapp-rg`
  - Click **OK**

#### **Detalles de la instancia:**
- **Nombre**: `whatsapp-api-go` (o cualquier nombre único)
  - ⚠️ Este nombre debe ser único globalmente
  - Será tu URL: `whatsapp-api-go.azurewebsites.net`
  
- **Publicar**: Selecciona **"Contenedor Docker"** o **"Docker Container"**

- **Sistema operativo**: Selecciona **"Linux"**

- **Región**: Selecciona la más cercana a ti o tus usuarios
  - Recomendado: `East US`, `West Europe`, o `Brazil South`

#### **Plan de App Service:**
- Click en **"Crear nuevo"**
- Nombre del plan: `whatsapp-plan`
- Click en **"Cambiar tamaño"** o **"Change size"**

**Seleccionar Plan:**
- Para desarrollo/pruebas: **"B1 (Basic)"** (~$13/mes)
  - 1 CPU core
  - 1.75 GB RAM
  - Suficiente para empezar
  
- Para producción: **"P1V2 (Premium)"** (~$73/mes)
  - 1 CPU core
  - 3.5 GB RAM
  - Auto-scaling
  - Mejor performance

- Click **"Aplicar"**

### 2.4. Configurar Contenedor (Nueva pestaña)

Después de seleccionar "Contenedor Docker", verás una nueva pestaña **"Docker"**:

1. Click en **"Siguiente: Docker"** (abajo)

En la pestaña Docker:
- **Opciones**: Selecciona **"Contenedor único"** o **"Single Container"**
- **Origen de imagen**: Selecciona **"Docker Hub"** (por ahora)
- **Acceso**: Selecciona **"Público"**
- **Imagen y etiqueta**: Escribe `alpine:latest` (temporal, lo cambiaremos después)

2. Click en **"Siguiente: Redes"** → **"Siguiente: Supervisión"** → **"Revisar y crear"**

### 2.5. Revisar y Crear

1. Click en **"Crear"**
2. ⏳ Espera 2-3 minutos mientras se crea...
3. Verás un mensaje: **"La implementación se completó"**
4. Click en **"Ir al recurso"**

✅ **¡App Service creado!**

---

## 🔧 **PASO 3: CONFIGURAR VARIABLES DE ENTORNO**

### 3.1. Ir a Configuración

1. En tu App Service, en el menú de la izquierda busca:
2. **"Configuración"** o **"Configuration"** (sección "Configuración")
3. Click en **"Configuración"**

### 3.2. Agregar Variables

1. Verás la pestaña **"Configuración de la aplicación"** o **"Application settings"**
2. Click en **"+ Nueva configuración de la aplicación"** o **"+ New application setting"**

**Agrega CADA UNA de estas variables:**

#### Variable 1:
- **Nombre**: `MONGODB_URL`
- **Valor**: `mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0`
- Click **OK**

#### Variable 2:
- **Nombre**: `MONGO_DB`
- **Valor**: `whatsapp_api`
- Click **OK**

#### Variable 3:
- **Nombre**: `WHATSAPP_VERIFY_TOKEN`
- **Valor**: `mi_token_seguro_whatsapp_2024`
- Click **OK**

#### Variable 4:
- **Nombre**: `WHATSAPP_APP_SECRET`
- **Valor**: `451614ef9eb9b35571dc352af6b2110e`
- Click **OK**

#### Variable 5:
- **Nombre**: `WABA_PHONE_ID`
- **Valor**: `804818756055720`
- Click **OK**

#### Variable 6:
- **Nombre**: `WABA_TOKEN`
- **Valor**: `EACO8kt4CNU4BP5ZAjnyEZBsatkJmx2XvPvYfO9cCllcjZANi1UeTvuh9LBWQ2t3Rse4B0q4rij37Ml3vgiFQB6krHWYhdW6mUkfRZBrA6w3ZBOZBYL1AAgTZCyS1Ls5zB4OwZAqPB2Dgpcz8Ucn2TjnPzGbVD3zza6IKmlGlYsaLSC3SNBHXvWjNj4W1FRPXtiY7y2ksG7n7xDzZBNe6kYTM3p0OZCc5RivuQrwpb6v4D4lKRGD5Ut2R2ownJzp8NRJ2BgfoeLttq5Cw5FOup76vssYgZDZD`
- Click **OK**

#### Variable 7:
- **Nombre**: `API_PORT`
- **Valor**: `8080`
- Click **OK**

#### Variable 8:
- **Nombre**: `LOG_LEVEL`
- **Valor**: `INFO`
- Click **OK**

### 3.3. Guardar

1. Click en **"Guardar"** o **"Save"** (arriba)
2. Aparecerá un mensaje de confirmación
3. Click **"Continuar"** o **"Continue"**
4. ⏳ Espera unos segundos...

✅ **Variables configuradas!**

---

## 🔌 **PASO 4: CONFIGURAR COMANDO DE INICIO**

### 4.1. Ir a Configuración General

1. En el mismo menú **"Configuración"**
2. Click en la pestaña **"Configuración general"** o **"General settings"**

### 4.2. Comando de inicio

1. Busca el campo: **"Comando de inicio"** o **"Startup Command"**
2. Escribe exactamente: `./whatsapp-api-server`
3. Click en **"Guardar"** (arriba)

✅ **Comando configurado!**

---

## 🔗 **PASO 5: CONFIGURAR DEPLOYMENT CON DOCKER**

### 5.1. Ir a Centro de Implementación

1. En el menú de la izquierda, busca:
2. **"Centro de implementación"** o **"Deployment Center"**
3. Click ahí

### 5.2. Configurar Origen

1. En **"Origen"** o **"Source"**, selecciona: **"GitHub"**
2. Click en **"Autorizar"** o **"Authorize"**
3. Se abrirá una ventana de GitHub
4. Inicia sesión en GitHub si es necesario
5. Click en **"Authorize Azure App Service"**

### 5.3. Seleccionar Repositorio

Después de autorizar, verás estos campos:

- **Organización**: Selecciona tu usuario de GitHub (`excsxavox`)
- **Repositorio**: Selecciona `whatsappgopu`
- **Rama**: Selecciona `main`

### 5.4. Configuración de Build

Como usamos Docker, Azure necesita:
- **Tipo de compilación**: Selecciona **"GitHub Actions"**
- Azure detectará automáticamente el `Dockerfile` en tu repositorio

### 5.5. Configuración del Dockerfile

Azure preguntará:
- **Archivo de Docker**: `/Dockerfile` (debe detectarlo automáticamente)
- **Contexto**: `/` (raíz del repositorio)

### 5.6. Guardar

1. Click en **"Guardar"** o **"Save"** (arriba)
2. ⏳ Espera unos segundos...

✅ **GitHub conectado!**

Azure ahora:
- Creó un archivo `.github/workflows/main_whatsapp-api-go.yml` en tu repo
- Compilará la imagen Docker automáticamente
- Cada vez que hagas `git push`, se construirá y desplegará el contenedor

---

## 🚀 **PASO 6: INICIAR PRIMER DEPLOYMENT**

### 6.1. Ver Deployment en Progreso

1. Quédate en **"Centro de implementación"**
2. Verás en la sección **"Registros"** o **"Logs"**
3. Aparecerá una entrada nueva con estado **"En curso"** o **"In Progress"**

### 6.2. Esperar Deployment

⏳ Este proceso toma **5-10 minutos** la primera vez:
- Azure descarga el código de GitHub
- Instala Go
- Compila la aplicación
- Inicia el servidor

### 6.3. Ver Logs

Para ver qué está pasando:
1. Click en la entrada del deployment
2. Verás logs en tiempo real
3. Espera hasta ver: **"Deployment successful"** o **"Implementación correcta"**

✅ **¡App desplegada!**

---

## 🌐 **PASO 7: OBTENER TU URL**

### 7.1. Ir a Información General

1. En el menú de la izquierda
2. Click en **"Información general"** o **"Overview"**

### 7.2. Copiar URL

1. Busca el campo: **"Dominio predeterminado"** o **"Default domain"**
2. Verás algo como: `whatsapp-api-go.azurewebsites.net`
3. **Copia esta URL** (la necesitarás para Meta)

### 7.3. Probar tu API

Abre en el navegador:
```
https://whatsapp-api-go.azurewebsites.net/health
```

**Deberías ver:**
```json
{
  "status": "ok",
  "timestamp": "..."
}
```

✅ **¡API funcionando!**

---

## 📱 **PASO 8: CONFIGURAR WEBHOOK EN META**

### 8.1. Ir a Meta Developers

1. Abre: **https://developers.facebook.com/**
2. Inicia sesión
3. Ve a **"Mis aplicaciones"**
4. Selecciona tu app (ID: `10058963160806734`)

### 8.2. Configurar WhatsApp

1. En el menú lateral, click en **"WhatsApp"**
2. Click en **"Configuración"** o **"Configuration"**

### 8.3. Configurar Webhook

Busca la sección **"Webhook"**:

1. Click en **"Editar"** o **"Edit"**

**Configuración:**
- **URL de devolución de llamada**: 
  ```
  https://whatsapp-api-go.azurewebsites.net/webhook
  ```
  (Reemplaza `whatsapp-api-go` con el nombre de tu App Service)

- **Token de verificación**: 
  ```
  mi_token_seguro_whatsapp_2024
  ```

2. Click en **"Verificar y guardar"** o **"Verify and save"**

⏳ Espera unos segundos...

**Deberías ver:**
- ✅ **"Webhook verificado correctamente"**

### 8.4. Suscribirse a Campos

Más abajo, en **"Campos del webhook"** o **"Webhook fields"**:

1. Busca **"messages"** y actívalo (toggle ON)
2. Click en **"Guardar"** o **"Save"**

✅ **¡Webhook configurado!**

---

## 🧪 **PASO 9: PROBAR ENVIANDO MENSAJE**

### 9.1. Enviar Mensaje de WhatsApp

1. Desde tu WhatsApp, envía un mensaje al número de prueba
2. Escribe: **"Hola"**

### 9.2. Ver Logs en Azure

**Para ver qué está pasando:**

1. En tu App Service de Azure
2. Menú lateral → **"Supervisión"** → **"Secuencia de registro"** o **"Log stream"**
3. Verás logs en tiempo real

**Deberías ver:**
```
📨 Mensaje entrante from: 521234567890
✅ Mensaje guardado en MongoDB
🔄 Iniciando flujo por defecto
```

### 9.3. Recibir Respuesta

Tu WhatsApp debería recibir:
- Si hay un flujo configurado: El primer mensaje del flujo
- Si no hay flujo: Un mensaje de confirmación

✅ **¡Sistema funcionando completo!**

---

## 🔍 **PASO 10: VER LOGS Y MONITOREAR**

### 10.1. Ver Logs en Tiempo Real

**Opción 1: Log Stream (Recomendado)**
1. Azure Portal → Tu App Service
2. **"Supervisión"** → **"Secuencia de registro"**
3. Verás logs en vivo

**Opción 2: Descargar Logs**
1. **"Supervisión"** → **"Registros de diagnóstico"**
2. Activa **"Registro de aplicaciones (Sistema de archivos)"**
3. Click **"Guardar"**

### 10.2. Ver Métricas

1. **"Supervisión"** → **"Métricas"**
2. Puedes ver:
   - CPU usage
   - Memory usage
   - HTTP requests
   - Response time

### 10.3. Configurar Alertas (Opcional)

1. **"Supervisión"** → **"Alertas"**
2. **"+ Crear"** → **"Regla de alerta"**
3. Configura alertas para:
   - CPU > 80%
   - Memory > 80%
   - HTTP 5xx errors

---

## 🛠️ **TROUBLESHOOTING**

### ❌ Error: "Application failed to start"

**Solución:**
1. Ve a **"Log stream"** y revisa los logs
2. Verifica que todas las variables estén configuradas
3. Verifica que el **Comando de inicio** sea: `./whatsapp-api-server`

### ❌ Error: "Cannot connect to MongoDB"

**Solución:**
1. Ve a MongoDB Atlas → Network Access
2. Click **"Add IP Address"**
3. Agrega: `0.0.0.0/0` (permitir todas las IPs)
4. Click **"Confirm"**

### ❌ Webhook no se verifica

**Solución:**
1. Verifica que la URL sea correcta y accesible
2. Prueba en el navegador: `https://tu-app.azurewebsites.net/health`
3. Verifica que `WHATSAPP_VERIFY_TOKEN` esté correcto

### ❌ La app no responde mensajes

**Solución:**
1. Ve a **"Log stream"** y mira los logs
2. Verifica que `WABA_TOKEN` y `WABA_PHONE_ID` estén correctos
3. Verifica que el webhook esté suscrito a "messages"

---

## 💰 **COSTOS**

### Plan B1 (Lo que probablemente elegiste)
- **$13 USD/mes** (~$0.43/día)
- Incluye:
  - 1 vCPU
  - 1.75 GB RAM
  - 10 GB almacenamiento
  - Ancho de banda incluido

### Cómo Reducir Costos
1. Si solo es para pruebas, apaga el App Service cuando no lo uses
2. Puedes cambiar a un plan más barato después (F1 gratis, pero muy limitado)

---

## ✅ **RESUMEN: ¿QUÉ LOGRASTE?**

🎉 **¡Felicidades! Ahora tienes:**

✅ Una API de WhatsApp funcionando en la nube (Azure)
✅ Conectada a MongoDB Atlas
✅ Con sistema completo de flujos implementado
✅ Deployment automático desde GitHub
✅ Webhook configurado con Meta
✅ Sistema listo para procesar conversaciones

---

## 📞 **SIGUIENTE PASO: CREAR TU PRIMER FLUJO**

Ahora que está funcionando, el siguiente paso es crear un flujo en MongoDB.

¿Quieres que te ayude a crear un flujo de ejemplo? 🚀

