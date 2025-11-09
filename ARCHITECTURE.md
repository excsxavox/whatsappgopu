# Arquitectura Hexagonal - WhatsApp API

## 🏗️ Principios de Arquitectura Hexagonal (Ports & Adapters)

Esta aplicación sigue estrictamente los principios de **Arquitectura Hexagonal** (también conocida como Ports & Adapters), donde:

- El **dominio** es el núcleo de la aplicación y NO tiene dependencias externas
- Los **puertos** definen interfaces para comunicarse con el exterior
- Los **adaptadores** implementan los puertos usando tecnologías específicas
- La **inyección de dependencias** se hace en el punto de entrada (main.go)

## 📂 Estructura del Proyecto

```
whatsapp-api-go/
├── cmd/
│   └── server/
│       └── main.go                      # 🎯 Punto de entrada - Inyección de dependencias
│
├── internal/
│   ├── domain/                          # 🧠 NÚCLEO - Sin dependencias externas
│   │   ├── entities/                    # Entidades del dominio
│   │   │   ├── message.go              # Message, MessageType, MessageStatus
│   │   │   ├── session.go              # Session
│   │   │   ├── connection.go           # Connection
│   │   │   └── errors.go               # Errores del dominio
│   │   │
│   │   ├── ports/                       # 🔌 Interfaces (Puertos)
│   │   │   ├── messaging.go            # MessagingService, MessageRepository
│   │   │   ├── session.go              # SessionRepository, SessionManager
│   │   │   ├── usecases.go             # Puertos de entrada (use cases)
│   │   │   └── logger.go               # Logger
│   │   │
│   │   └── services/                    # Servicios del dominio (lógica de negocio)
│   │
│   ├── application/                     # 🎬 CASOS DE USO - Orquestación
│   │   └── usecases/
│   │       ├── send_message.go         # Enviar mensaje
│   │       ├── get_connection_status.go # Obtener estado
│   │       ├── establish_connection.go  # Conectar
│   │       ├── disconnect.go           # Desconectar
│   │       └── handle_webhook.go       # Procesar webhooks
│   │
│   └── infrastructure/                  # 🔧 ADAPTADORES - Implementaciones concretas
│       ├── adapters/
│       │   ├── whatsapp/               # Adaptador de WhatsApp (whatsmeow)
│       │   │   └── adapter.go          # Implementa MessagingService
│       │   │
│       │   ├── http/                   # Adaptador HTTP REST
│       │   │   ├── handlers.go         # Handlers HTTP
│       │   │   └── server.go           # Servidor HTTP
│       │   │
│       │   └── storage/                # Adaptador de almacenamiento
│       │       ├── message_repository.go  # Implementa MessageRepository
│       │       └── session_repository.go  # Implementa SessionRepository
│       │
│       └── config/
│           └── config.go               # Configuración
│
└── pkg/                                 # 🛠️ UTILIDADES - Compartidas
    └── logger/
        └── logger.go                    # Implementaciones de Logger
```

## 🔄 Flujo de Datos (Hexagonal)

### Ejemplo: Enviar un Mensaje

```
1. HTTP Request POST /send
       ↓
2. [HTTP Adapter] → handlers.go
       ↓
3. [Use Case] → SendMessageUseCase.Execute()
       ↓
4. [Domain Logic] → Validación con entities.Message
       ↓
5. [Port] → MessagingService.SendMessage()
       ↓
6. [WhatsApp Adapter] → whatsmeow (implementación)
       ↓
7. WhatsApp Web API
```

### Flujo inverso (webhook):

```
1. Webhook POST /webhook
       ↓
2. [HTTP Adapter] → HandleWebhookHandler
       ↓
3. [Use Case] → HandleWebhookUseCase.Execute()
       ↓
4. Delega a → SendMessageUseCase
       ↓
5. [Domain Logic] → ...
```

## 🎯 Capas y Responsabilidades

### 1. Domain (Dominio)

**Responsabilidad**: Contener la lógica de negocio pura.

**Reglas**:
- ❌ NO puede depender de frameworks externos
- ❌ NO puede depender de bases de datos específicas
- ❌ NO puede depender de HTTP o protocolos de red
- ✅ Solo contiene lógica de negocio pura
- ✅ Define interfaces (ports) para comunicarse con el exterior

**Archivos**:
- `entities/`: Objetos de dominio (Message, Session, Connection)
- `ports/`: Interfaces que definen contratos
- `services/`: Lógica de negocio compleja

### 2. Application (Aplicación)

**Responsabilidad**: Orquestar los casos de uso.

**Reglas**:
- ✅ Usa las interfaces del dominio (ports)
- ✅ Coordina el flujo de datos entre adaptadores
- ❌ NO contiene lógica de negocio (eso va en domain)
- ❌ NO sabe de implementaciones concretas

**Archivos**:
- `usecases/`: Implementaciones de casos de uso

### 3. Infrastructure (Infraestructura)

**Responsabilidad**: Implementar las interfaces del dominio con tecnologías específicas.

**Reglas**:
- ✅ Implementa los puertos (interfaces) del dominio
- ✅ Puede usar frameworks y librerías externas
- ✅ Puede acceder a bases de datos, APIs, etc.
- ❌ NO debe contener lógica de negocio

**Adaptadores**:
- `whatsapp/`: Implementa MessagingService usando whatsmeow
- `http/`: Expone API REST para casos de uso
- `storage/`: Implementa repositorios (en memoria o DB)

### 4. Ports (Puertos)

**Tipos de puertos**:

#### Puertos de ENTRADA (Driven Ports)
Lo que el dominio OFRECE al exterior:
- `SendMessageUseCase`
- `GetConnectionStatusUseCase`
- `HandleWebhookUseCase`

#### Puertos de SALIDA (Driving Ports)
Lo que el dominio NECESITA del exterior:
- `MessagingService` (para enviar mensajes)
- `MessageRepository` (para persistir mensajes)
- `SessionRepository` (para gestionar sesiones)
- `Logger` (para logging)

## 🔌 Inyección de Dependencias

Todo se conecta en `cmd/server/main.go`:

```go
// 1. Inicializar adaptadores (infraestructura)
logger := logger.NewColorLogger()
whatsappAdapter := whatsapp.NewWhatsAppAdapter(logger)
messageRepo := storage.NewInMemoryMessageRepository()

// 2. Inyectar en casos de uso (aplicación)
sendMessageUseCase := usecases.NewSendMessageUseCase(
    whatsappAdapter,  // implementa MessagingService
    messageRepo,      // implementa MessageRepository
    logger,           // implementa Logger
)

// 3. Inyectar en adaptadores de entrada
httpAdapter := http.NewHTTPAdapter(
    sendMessageUseCase,
    getConnectionStatusUseCase,
    handleWebhookUseCase,
    logger,
)

// 4. Iniciar servidor
httpServer.Start()
```

## ✅ Beneficios de esta Arquitectura

### 1. **Independencia de Frameworks**
El dominio no depende de whatsmeow, HTTP, o cualquier framework. Podemos cambiarlos sin afectar la lógica de negocio.

### 2. **Testeable**
```go
// Fácil de testear con mocks
mockMessaging := &MockMessagingService{}
useCase := NewSendMessageUseCase(mockMessaging, mockRepo, mockLogger)
```

### 3. **Mantenible**
- Cada capa tiene responsabilidades claras
- Los cambios están aislados
- Fácil de entender y navegar

### 4. **Escalable**
Podemos agregar nuevos adaptadores sin modificar el dominio:
- Cambiar de SQLite a PostgreSQL
- Agregar adaptador gRPC además de HTTP
- Cambiar whatsmeow por otra librería

### 5. **Reusable**
Los casos de uso pueden ser llamados desde:
- HTTP REST API
- gRPC
- CLI
- WebSockets
- Cron jobs

## 🔄 Cómo Agregar Nueva Funcionalidad

### Ejemplo: Agregar "Enviar Imagen"

#### 1. Domain (si es necesario)
```go
// internal/domain/entities/message.go
const MessageTypeImage MessageType = "image"
```

#### 2. Port (si es necesario)
```go
// internal/domain/ports/messaging.go
type MessagingService interface {
    SendImage(ctx context.Context, to string, imageURL string) error
}
```

#### 3. Use Case
```go
// internal/application/usecases/send_image.go
type SendImageUseCase struct {
    messagingService ports.MessagingService
}

func (uc *SendImageUseCase) Execute(ctx context.Context, to, imageURL string) error {
    // lógica del caso de uso
}
```

#### 4. Adapter
```go
// internal/infrastructure/adapters/whatsapp/adapter.go
func (a *WhatsAppAdapter) SendImage(ctx context.Context, to, imageURL string) error {
    // implementación con whatsmeow
}
```

#### 5. HTTP Handler
```go
// internal/infrastructure/adapters/http/handlers.go
func (h *HTTPAdapter) SendImageHandler(w http.ResponseWriter, r *http.Request) {
    // llamar al use case
}
```

#### 6. Wire en main.go
```go
// cmd/server/main.go
sendImageUseCase := usecases.NewSendImageUseCase(whatsappAdapter, logger)
httpAdapter := http.NewHTTPAdapter(..., sendImageUseCase, ...)
```

## 📚 Referencias

- [Hexagonal Architecture (Alistair Cockburn)](https://alistair.cockburn.us/hexagonal-architecture/)
- [Clean Architecture (Robert C. Martin)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Ports and Adapters Pattern](https://herbertograca.com/2017/09/14/ports-adapters-architecture/)

## 🎓 Conclusión

Esta arquitectura garantiza:
- ✅ **Dominio puro**: Sin dependencias externas
- ✅ **Flexibilidad**: Cambiar tecnologías fácilmente
- ✅ **Testeable**: Mockear cualquier dependencia
- ✅ **Mantenible**: Código organizado y claro
- ✅ **Escalable**: Agregar funcionalidades sin romper nada

