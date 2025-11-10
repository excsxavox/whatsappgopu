# 🚀 DESPLEGAR EN AZURE APP SERVICE

## 📋 **REQUISITOS PREVIOS**

1. ✅ Cuenta de Azure activa
2. ✅ Azure CLI instalado (opcional pero recomendado)
3. ✅ Código en GitHub
4. ✅ MongoDB Atlas configurado

---

## 🎯 **MÉTODO 1: PORTAL DE AZURE (MÁS FÁCIL)**

### **PASO 1: Crear App Service**

1. Ir a [Portal de Azure](https://portal.azure.com)
2. Click en **"Crear un recurso"**
3. Buscar **"App Service"** y click en **Crear**
4. Configurar:
   - **Suscripción**: Tu suscripción activa
   - **Grupo de recursos**: Crear nuevo o usar existente
   - **Nombre**: `whatsapp-api-go` (debe ser único globalmente)
   - **Publicar**: **Código**
   - **Pila del entorno en tiempo de ejecución**: **Go 1.21** (o la más reciente)
   - **Sistema operativo**: **Linux**
   - **Región**: Elegir la más cercana a tus usuarios
   - **Plan de App Service**: 
     - Para desarrollo: **B1 Basic** (~$13/mes)
     - Para producción: **P1V2 Premium** (~$73/mes)

5. Click en **"Revisar y crear"** → **Crear**

### **PASO 2: Configurar Variables de Entorno**

1. Ir a tu App Service
2. En el menú lateral: **Configuración** → **Configuración de la aplicación**
3. Click en **"Nueva configuración de la aplicación"**
4. Agregar todas estas variables:

```
MONGODB_URL=mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0
MONGO_DB=whatsapp_api
WHATSAPP_VERIFY_TOKEN=mi_token_seguro_whatsapp_2024
WHATSAPP_APP_SECRET=451614ef9eb9b35571dc352af6b2110e
WABA_PHONE_ID=804818756055720
WABA_TOKEN=EACO8kt4CNU4BP5ZAjnyEZBsatkJmx2XvPvYfO9cCllcjZANi1UeTvuh9LBWQ2t3Rse4B0q4rij37Ml3vgiFQB6krHWYhdW6mUkfRZBrA6w3ZBOZBYL1AAgTZCyS1Ls5zB4OwZAqPB2Dgpcz8Ucn2TjnPzGbVD3zza6IKmlGlYsaLSC3SNBHXvWjNj4W1FRPXtiY7y2ksG7n7xDzZBNe6kYTM3p0OZCc5RivuQrwpb6v4D4lKRGD5Ut2R2ownJzp8NRJ2BgfoeLttq5Cw5FOup76vssYgZDZD
API_PORT=8080
LOG_LEVEL=INFO
```

5. Click en **Guardar** (arriba)

### **PASO 3: Configurar Puerto**

1. En **Configuración** → **Configuración general**
2. En **"Comando de inicio"**, agregar:
   ```
   ./whatsapp-api-server
   ```
3. Guardar

### **PASO 4: Desplegar desde GitHub**

#### **Opción A: Centro de implementación (Recomendado)**

1. Ir a **Centro de implementación**
2. Seleccionar **GitHub**
3. Autorizar Azure con tu cuenta de GitHub
4. Seleccionar:
   - **Organización**: Tu usuario de GitHub
   - **Repositorio**: `whatsappgopu`
   - **Rama**: `main`
5. Click en **Guardar**

Azure creará automáticamente un workflow de GitHub Actions.

#### **Opción B: GitHub Actions Manual**

1. En tu App Service, ir a **Centro de implementación**
2. Click en **"Administrar perfil de publicación"**
3. Descargar el archivo `.publishsettings`
4. Ir a tu repositorio de GitHub
5. **Settings** → **Secrets and variables** → **Actions**
6. Click en **"New repository secret"**
7. Nombre: `AZURE_WEBAPP_PUBLISH_PROFILE`
8. Valor: Pegar todo el contenido del archivo `.publishsettings`
9. Copiar el archivo `azure-deploy.yml` a `.github/workflows/`
10. Hacer commit y push

---

## 🎯 **MÉTODO 2: AZURE CLI (AVANZADO)**

### **PASO 1: Instalar Azure CLI**

```bash
# Windows (con winget)
winget install Microsoft.AzureCLI

# O descargar de: https://aka.ms/installazurecliwindows
```

### **PASO 2: Login**

```bash
az login
```

### **PASO 3: Crear Grupo de Recursos**

```bash
az group create --name whatsapp-rg --location eastus
```

### **PASO 4: Crear Plan de App Service**

```bash
az appservice plan create \
  --name whatsapp-plan \
  --resource-group whatsapp-rg \
  --sku B1 \
  --is-linux
```

### **PASO 5: Crear Web App**

```bash
az webapp create \
  --resource-group whatsapp-rg \
  --plan whatsapp-plan \
  --name whatsapp-api-go-tu-nombre \
  --runtime "GO:1.21"
```

### **PASO 6: Configurar Variables de Entorno**

```bash
az webapp config appsettings set \
  --resource-group whatsapp-rg \
  --name whatsapp-api-go-tu-nombre \
  --settings \
    MONGODB_URL="mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0" \
    MONGO_DB="whatsapp_api" \
    WHATSAPP_VERIFY_TOKEN="mi_token_seguro_whatsapp_2024" \
    WHATSAPP_APP_SECRET="451614ef9eb9b35571dc352af6b2110e" \
    WABA_PHONE_ID="804818756055720" \
    WABA_TOKEN="EACO8kt4CNU4BP5ZAjnyEZBsatkJmx2XvPvYfO9cCllcjZANi1UeTvuh9LBWQ2t3Rse4B0q4rij37Ml3vgiFQB6krHWYhdW6mUkfRZBrA6w3ZBOZBYL1AAgTZCyS1Ls5zB4OwZAqPB2Dgpcz8Ucn2TjnPzGbVD3zza6IKmlGlYsaLSC3SNBHXvWjNj4W1FRPXtiY7y2ksG7n7xDzZBNe6kYTM3p0OZCc5RivuQrwpb6v4D4lKRGD5Ut2R2ownJzp8NRJ2BgfoeLttq5Cw5FOup76vssYgZDZD" \
    API_PORT="8080" \
    LOG_LEVEL="INFO"
```

### **PASO 7: Configurar Deployment desde GitHub**

```bash
az webapp deployment source config \
  --resource-group whatsapp-rg \
  --name whatsapp-api-go-tu-nombre \
  --repo-url https://github.com/excsxavox/whatsappgopu \
  --branch main \
  --manual-integration
```

---

## 🎯 **MÉTODO 3: DOCKER EN APP SERVICE (RECOMENDADO PARA PRODUCCIÓN)**

### **PASO 1: Crear Container Registry**

```bash
az acr create \
  --resource-group whatsapp-rg \
  --name whatsappregistry \
  --sku Basic
```

### **PASO 2: Build y Push de la Imagen**

```bash
# Login al registry
az acr login --name whatsappregistry

# Build y push
az acr build \
  --registry whatsappregistry \
  --image whatsapp-api:latest \
  .
```

### **PASO 3: Crear Web App con Docker**

```bash
az webapp create \
  --resource-group whatsapp-rg \
  --plan whatsapp-plan \
  --name whatsapp-api-go-tu-nombre \
  --deployment-container-image-name whatsappregistry.azurecr.io/whatsapp-api:latest
```

### **PASO 4: Configurar Variables de Entorno**

```bash
az webapp config appsettings set \
  --resource-group whatsapp-rg \
  --name whatsapp-api-go-tu-nombre \
  --settings \
    MONGODB_URL="..." \
    MONGO_DB="whatsapp_api" \
    # ... resto de variables
```

---

## 📡 **CONFIGURAR WEBHOOK EN META**

Una vez desplegado, tu URL será:
```
https://whatsapp-api-go-tu-nombre.azurewebsites.net
```

1. Ir a [Meta Developers](https://developers.facebook.com/)
2. Tu App → **WhatsApp** → **Configuración**
3. **URL de devolución de llamada**: 
   ```
   https://whatsapp-api-go-tu-nombre.azurewebsites.net/webhook
   ```
4. **Token de verificación**: `mi_token_seguro_whatsapp_2024`
5. Click en **Verificar y guardar**

---

## 🔍 **VERIFICAR DEPLOYMENT**

### **1. Ver logs en tiempo real**

```bash
az webapp log tail \
  --resource-group whatsapp-rg \
  --name whatsapp-api-go-tu-nombre
```

### **2. Probar endpoint**

```bash
curl https://whatsapp-api-go-tu-nombre.azurewebsites.net/health
```

### **3. Ver métricas**

En el Portal de Azure:
- **Supervisión** → **Métricas**
- **Supervisión** → **Registros**

---

## 💰 **COSTOS ESTIMADOS**

### **Plan B1 Basic (Desarrollo)**
- **Precio**: ~$13 USD/mes
- **Especificaciones**:
  - 1 vCPU
  - 1.75 GB RAM
  - 10 GB almacenamiento
- **Ideal para**: Pruebas, desarrollo, bajo tráfico

### **Plan P1V2 Premium (Producción)**
- **Precio**: ~$73 USD/mes
- **Especificaciones**:
  - 1 vCPU
  - 3.5 GB RAM
  - 250 GB almacenamiento
  - Auto-scaling
- **Ideal para**: Producción, alto tráfico

### **Gratis (F1)**
- **Precio**: $0
- **Limitaciones**:
  - 60 minutos CPU/día
  - 1 GB RAM
  - No custom domains
  - Solo para pruebas muy básicas

---

## 🔧 **TROUBLESHOOTING**

### **Error: "Application failed to start"**

1. Ver logs:
   ```bash
   az webapp log tail --name whatsapp-api-go-tu-nombre --resource-group whatsapp-rg
   ```

2. Verificar que `API_PORT=8080` esté configurado

3. Verificar comando de inicio: `./whatsapp-api-server`

### **Error: "Cannot connect to MongoDB"**

1. Verificar que `MONGODB_URL` esté correctamente configurado
2. Verificar que MongoDB Atlas permita conexiones desde Azure:
   - MongoDB Atlas → Network Access → Add IP Address
   - Agregar: `0.0.0.0/0` (permitir todas las IPs)

### **Webhook no funciona**

1. Verificar que la URL sea accesible:
   ```bash
   curl https://whatsapp-api-go-tu-nombre.azurewebsites.net/webhook
   ```

2. Verificar variables `WHATSAPP_VERIFY_TOKEN` y `WHATSAPP_APP_SECRET`

---

## 🎓 **RECOMENDACIONES**

### ✅ **HACER**
1. Usar **HTTPS** siempre (Azure lo proporciona gratis)
2. Configurar **Application Insights** para monitoreo
3. Activar **Auto-scaling** en producción
4. Usar **Deployment Slots** para staging/production
5. Configurar **Alertas** de CPU/RAM

### ❌ **NO HACER**
1. No usar plan F1 (Free) para producción
2. No exponer tokens en logs
3. No olvidar configurar IP whitelist en MongoDB
4. No usar la misma app para dev y prod

---

## 📚 **RECURSOS ADICIONALES**

- [Documentación Azure App Service](https://docs.microsoft.com/azure/app-service/)
- [Go en Azure](https://docs.microsoft.com/azure/app-service/quickstart-golang)
- [CI/CD con GitHub Actions](https://docs.microsoft.com/azure/app-service/deploy-github-actions)

---

## ✅ **SIGUIENTE PASO**

¿Qué método prefieres?

1. **Portal de Azure** → Más visual, recomendado para comenzar
2. **Azure CLI** → Más rápido si ya tienes experiencia
3. **Docker** → Más control y portable

Dime cuál eliges y te guío paso a paso! 🚀

