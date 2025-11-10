================================================================================
GUÍA DE PROCESAMIENTO DE FLUJOS DE WHATSAPP
Cómo debe funcionar el sistema para procesar nodos y mantener conversaciones
================================================================================

OBJETIVO PRINCIPAL:
-------------------
La aplicación debe procesar cada tipo de nodo en un flujo de WhatsApp y mantener
la conversación dentro del flujo, sin que el usuario se salga del flujo hasta
completarlo o abandonarlo.

================================================================================
1. GESTIÓN DE SESIONES DE CONVERSACIÓN
================================================================================

¿QUÉ ES UNA SESIÓN?
-------------------
Una sesión es el estado actual de una conversación dentro de un flujo específico.
Cada conversación (identificada por conversationId) debe tener una sesión activa
que rastree:

- ID del flujo actual
- ID del nodo actual donde está el usuario
- Variables capturadas durante la conversación (nombre_usuario, imagen_cedula, etc.)
- Si está esperando una respuesta del usuario (waitingForResponse)
- Qué variable está esperando (waitingForVariable)
- Historial de nodos ejecutados

CUANDO CREAR UNA SESIÓN:
------------------------
- Cuando un usuario inicia una conversación y se asigna un flujo
- Cuando un usuario envía un mensaje y no hay sesión activa (iniciar flujo por defecto)
- Cuando se activa un flujo manualmente para un usuario

CUANDO ACTUALIZAR UNA SESIÓN:
-----------------------------
- Cada vez que se procesa un nodo
- Cada vez que se captura una variable del usuario
- Cada vez que se avanza al siguiente nodo
- Cuando se completa o abandona el flujo

CUANDO COMPLETAR UNA SESIÓN:
----------------------------
- Cuando se llega al final del flujo (no hay más nodos)
- Cuando se ejecuta un nodo que marca el flujo como completado
- Cuando el usuario completa exitosamente todos los pasos requeridos

CUANDO ABANDONAR UNA SESIÓN:
----------------------------
- Cuando el usuario no responde por un tiempo determinado (timeout)
- Cuando el usuario envía un comando para salir del flujo
- Cuando ocurre un error crítico que impide continuar

================================================================================
2. PROCESAMIENTO DE TIPOS DE NODOS
================================================================================

2.1. NODO TIPO: TEXT
--------------------
PROPÓSITO: Enviar un mensaje de texto al usuario.

CONFIGURACIÓN TÍPICA:
{
  "id": "node_1_bienvenida",
  "type": "TEXT",
  "config": {
    "content": "¡Hola! 👋 Bienvenido...",
    "bodyText": "¡Hola! 👋 Bienvenido...",
    "waitForResponse": false,  // Si espera respuesta del usuario
    "responseVariableName": "nombre_usuario",  // Variable donde guardar respuesta
    "responseType": "text",  // Tipo de respuesta esperada: text, image, audio
    "validation": {
      "required": true,
      "minLength": 3,
      "maxLength": 50
    }
  }
}

QUÉ DEBE HACER LA APLICACIÓN:
-----------------------------
1. Reemplazar variables en el contenido usando valores de la sesión:
   - {nombre_usuario} → valor de session.variables.nombre_usuario
   - [imagen_cedula] → valor de session.variables.imagen_cedula

2. Enviar el mensaje al usuario vía WhatsApp API

3. Si waitForResponse = true:
   - Actualizar sesión: waitingForResponse = true
   - Guardar: waitingForVariable = responseVariableName
   - NO avanzar al siguiente nodo, esperar respuesta del usuario
   - El siguiente mensaje del usuario se procesará como respuesta a esta variable

4. Si waitForResponse = false:
   - Avanzar automáticamente al siguiente nodo según los edges

EJEMPLO:
--------
Nodo: "¿Cuál es tu nombre?"
- waitForResponse: true
- responseVariableName: "nombre_usuario"
- Acción: Enviar mensaje, esperar respuesta, guardar en variables["nombre_usuario"]

================================================================================

2.2. NODO TIPO: BUTTONS
-----------------------
PROPÓSITO: Enviar botones interactivos al usuario.

CONFIGURACIÓN TÍPICA:
{
  "id": "node_3_menu_botones",
  "type": "BUTTONS",
  "config": {
    "content": "Selecciona una opción:",
    "buttons": [
      {
        "id": "btn_productos",
        "type": "reply",
        "title": "Productos"
      },
      {
        "id": "btn_soporte",
        "type": "reply",
        "title": "Soporte"
      }
    ],
    "responseVariableName": "button_response"
  }
}

QUÉ DEBE HACER LA APLICACIÓN:
-----------------------------
1. Reemplazar variables en el contenido y títulos de botones

2. Enviar mensaje con botones interactivos vía WhatsApp API
   - Formato: interactive message con type: "button"

3. Siempre espera respuesta (implícito):
   - Actualizar sesión: waitingForResponse = true
   - Guardar: waitingForVariable = responseVariableName (ej: "button_response")
   - NO avanzar, esperar que usuario presione un botón

4. Cuando el usuario presiona un botón:
   - El mensaje recibido tendrá type: "interactive"
   - Extraer: message.interactive.button_reply.id (ej: "btn_productos")
   - Guardar en variables[responseVariableName] = "btn_productos"
   - Avanzar al siguiente nodo (generalmente un CONDITION)

EJEMPLO:
--------
Usuario presiona "Productos" → variables["button_response"] = "btn_productos"
Luego se procesa un nodo CONDITION que evalúa esta variable.

================================================================================

2.3. NODO TIPO: RESPONSE
------------------------
PROPÓSITO: Capturar y validar la respuesta del usuario.

CONFIGURACIÓN TÍPICA:
{
  "id": "node_3_response_nombre",
  "type": "RESPONSE",
  "config": {
    "variableName": "nombre_usuario",
    "responseType": "text",
    "validation": {
      "required": true,
      "minLength": 3,
      "maxLength": 50
    }
  }
}

QUÉ DEBE HACER LA APLICACIÓN:
-----------------------------
1. Este nodo se procesa DESPUÉS de que el usuario respondió a un nodo TEXT/BUTTONS

2. Validar la respuesta según las reglas de validation:
   - required: debe tener valor
   - minLength/maxLength: validar longitud
   - pattern: validar formato (si existe)

3. Si la validación falla:
   - Enviar mensaje de error al usuario
   - Volver al nodo anterior que solicitó la respuesta
   - Pedir nuevamente la información

4. Si la validación pasa:
   - La variable ya está guardada en la sesión (se guardó en handleUserResponse)
   - Continuar al siguiente nodo

NOTA: Este nodo es opcional. Si un nodo TEXT tiene waitForResponse=true,
la respuesta se guarda automáticamente. El nodo RESPONSE permite validación
adicional.

================================================================================

2.4. NODO TIPO: HTTP
---------------------
PROPÓSITO: Hacer una llamada HTTP a un endpoint externo.

CONFIGURACIÓN TÍPICA:
{
  "id": "node_6_http_validar_cedula",
  "type": "HTTP",
  "config": {
    "method": "POST",
    "url": "https://whatsapp-three-eta.vercel.app/api/whatsapp/ocr/validate-id",
    "headers": {
      "Content-Type": "application/json"
    },
    "body": {
      "image": "{imagen_cedula}"
    },
    "responseVariable": "response_validacion_cedula"
  }
}

QUÉ DEBE HACER LA APLICACIÓN:
-----------------------------
1. Reemplazar variables en URL, headers y body:
   - {imagen_cedula} → valor de session.variables.imagen_cedula
   - Si la variable es una URL de imagen, convertir a base64 si es necesario

2. Hacer la llamada HTTP con el método, URL, headers y body especificados

3. Procesar la respuesta:
   - Guardar la respuesta completa en variables[responseVariable]
   - Ejemplo: variables["response_validacion_cedula"] = { valid: true }

4. Manejar errores:
   - Si la llamada falla, guardar error en la variable
   - O lanzar excepción para que el flujo maneje el error

5. NO espera respuesta del usuario, avanza automáticamente al siguiente nodo

EJEMPLO:
--------
URL: /api/whatsapp/ocr/validate-id
Body: { "image": "base64_encoded_image_or_url" }
Respuesta: { "valid": true }
Guardado en: variables["response_validacion_cedula"] = { "valid": true }

================================================================================

2.5. NODO TIPO: CONDITION
--------------------------
PROPÓSITO: Evaluar una condición y seguir una rama u otra del flujo.

CONFIGURACIÓN TÍPICA:
{
  "id": "node_7_condition_cedula_valida",
  "type": "CONDITION",
  "config": {
    "conditions": [
      {
        "id": "cond_cedula_valida",
        "type": "si",
        "operator": "equals",
        "field": "response_validacion_cedula.valid",
        "value": true
      }
    ]
  }
}

QUÉ DEBE HACER LA APLICACIÓN:
-----------------------------
1. Evaluar cada condición usando las variables de la sesión:
   - field: ruta a la variable (ej: "response_validacion_cedula.valid")
   - operator: equals, not_equals, greater_than, less_than, contains, etc.
   - value: valor a comparar

2. Buscar edges que salen de este nodo:
   - Edge con condición "yes" o "si" → si la condición es verdadera
   - Edge con condición "no" → si la condición es falsa

3. Seguir el edge correspondiente:
   - Si condición verdadera → seguir edge "yes"
   - Si condición falsa → seguir edge "no"

4. Avanzar al nodo destino del edge seleccionado

EJEMPLO:
--------
Condición: response_validacion_cedula.valid == true
- Si es true → seguir edge "yes" → nodo "node_9_cedula_valida"
- Si es false → seguir edge "no" → nodo "node_8_cedula_invalida"

================================================================================

2.6. NODO TIPO: AUDIO
----------------------
PROPÓSITO: Enviar o recibir audio del usuario.

CONFIGURACIÓN TÍPICA (ENVIAR AUDIO):
{
  "id": "node_audio_1",
  "type": "AUDIO",
  "config": {
    "mediaType": "recorded",
    "hasRecordedAudio": true,
    "recordedAudio": "data:audio/webm;codecs=opus;base64,UklGRiQ...",
    "waitForVoiceResponse": false
  }
}

CONFIGURACIÓN TÍPICA (RECIBIR AUDIO):
{
  "id": "node_audio_2",
  "type": "AUDIO",
  "config": {
    "mediaType": "recorded",
    "waitForVoiceResponse": true,
    "responseVariableName": "audio_respuesta"
  }
}

QUÉ DEBE HACER LA APLICACIÓN:
-----------------------------
CASO 1: ENVIAR AUDIO (hasRecordedAudio = true)
1. Convertir el base64 del audio a un formato que WhatsApp acepte
2. Enviar el audio al usuario vía WhatsApp API
3. Si waitForVoiceResponse = false: avanzar al siguiente nodo
4. Si waitForVoiceResponse = true: esperar respuesta de audio del usuario

CASO 2: RECIBIR AUDIO (waitForVoiceResponse = true)
1. Enviar mensaje pidiendo al usuario que grabe un audio
2. Actualizar sesión: waitingForResponse = true
3. Guardar: waitingForVariable = responseVariableName
4. Cuando el usuario envía audio:
   - Guardar el ID del audio en variables[responseVariableName]
   - O descargar y convertir a base64 si es necesario
5. Avanzar al siguiente nodo

================================================================================
3. FLUJO DE EJECUCIÓN GENERAL
================================================================================

PASO 1: INICIO DE CONVERSACIÓN
-------------------------------
1. Usuario envía mensaje a WhatsApp
2. Webhook recibe el mensaje
3. Buscar sesión activa para conversationId:
   - Si existe → continuar con PASO 2
   - Si NO existe → iniciar flujo (PASO 1.1)

PASO 1.1: INICIAR FLUJO
------------------------
1. Determinar qué flujo usar:
   - Flujo por defecto del canal
   - Flujo basado en reglas de negocio
2. Crear sesión nueva:
   - flowId = ID del flujo
   - currentNodeId = entryNodeId del flujo
   - variables = {}
   - waitingForResponse = false
3. Procesar nodo de entrada (entryNodeId)

PASO 2: PROCESAR MENSAJE EN SESIÓN ACTIVA
------------------------------------------
1. Verificar estado de la sesión:
   - Si waitingForResponse = true:
     → Procesar como respuesta a la variable esperada (PASO 2.1)
   - Si waitingForResponse = false:
     → Procesar como nuevo mensaje o continuar flujo (PASO 2.2)

PASO 2.1: PROCESAR RESPUESTA DEL USUARIO
----------------------------------------
1. Extraer valor del mensaje según tipo:
   - text → message.text.body
   - image → message.image.id o URL
   - audio → message.audio.id
   - interactive → message.interactive.button_reply.id

2. Validar según el tipo esperado (responseType del nodo anterior)

3. Guardar en variables:
   - session.variables[waitingForVariable] = valor extraído

4. Actualizar sesión:
   - waitingForResponse = false
   - waitingForVariable = null

5. Buscar nodo RESPONSE correspondiente (si existe) y procesarlo

6. Avanzar al siguiente nodo según edges

PASO 2.2: CONTINUAR FLUJO
--------------------------
1. Si hay currentNodeId, procesar ese nodo
2. Si no hay currentNodeId, buscar flujo o iniciar uno nuevo

PASO 3: PROCESAR NODO
---------------------
1. Identificar tipo de nodo (TEXT, BUTTONS, HTTP, CONDITION, RESPONSE, AUDIO)

2. Llamar al procesador correspondiente:
   - TextNodeProcessor.process()
   - ButtonsNodeProcessor.process()
   - HttpNodeProcessor.process()
   - ConditionNodeProcessor.process()
   - ResponseNodeProcessor.process()
   - AudioNodeProcessor.process()

3. El procesador:
   - Reemplaza variables en la configuración
   - Ejecuta la acción del nodo (enviar mensaje, hacer HTTP, etc.)
   - Retorna: { waitingForResponse, waitingForVariable, stopFlow, etc. }

4. Actualizar sesión con el resultado

5. Si waitingForResponse = true:
   - Detener, esperar siguiente mensaje del usuario
6. Si waitingForResponse = false:
   - Avanzar al siguiente nodo (PASO 4)

PASO 4: AVANZAR AL SIGUIENTE NODO
----------------------------------
1. Buscar edges que salen del nodo actual:
   - edges.filter(e => e.from === currentNodeId)

2. Si no hay edges:
   - Completar sesión (status = "completed")
   - Finalizar flujo

3. Si hay un solo edge:
   - Obtener nodo destino: edges[0].to
   - Actualizar: session.currentNodeId = nodo destino
   - Procesar nodo destino (volver a PASO 3)

4. Si hay múltiples edges (nodo CONDITION):
   - El ConditionNodeProcessor ya seleccionó el edge correcto
   - Seguir el edge seleccionado
   - Procesar nodo destino

================================================================================
4. MANEJO DE EDGES (CONEXIONES ENTRE NODOS)
================================================================================

ESTRUCTURA DE UN EDGE:
----------------------
{
  "id": "edge_1_2",
  "from": "node_1_bienvenida",  // Nodo origen
  "to": "node_2_solicitar_nombre",  // Nodo destino
  "sourceHandle": "default",
  "targetHandle": "input",
  "delayMs": 0
}

REGLAS DE EDGES:
----------------
1. Un nodo TEXT sin waitForResponse debe tener UN SOLO edge saliente
2. Un nodo TEXT con waitForResponse NO debe tener edge directo
   - La respuesta del usuario activa el siguiente nodo
3. Un nodo BUTTONS siempre espera respuesta, no tiene edge directo
4. Un nodo HTTP no espera respuesta, tiene UN SOLO edge saliente
5. Un nodo CONDITION tiene DOS edges:
   - Uno con condición "yes" o "si"
   - Uno con condición "no"
6. Un nodo RESPONSE tiene UN SOLO edge saliente

FLUJO CON EDGES:
----------------
Nodo TEXT (waitForResponse=false) 
  → Edge único 
  → Siguiente nodo

Nodo TEXT (waitForResponse=true)
  → Espera respuesta usuario
  → Nodo RESPONSE (opcional)
  → Edge desde RESPONSE
  → Siguiente nodo

Nodo BUTTONS
  → Espera respuesta usuario (botón presionado)
  → Nodo CONDITION (generalmente)
  → Edge "yes" o "no" según botón
  → Siguiente nodo

Nodo HTTP
  → Ejecuta llamada
  → Edge único
  → Siguiente nodo (generalmente CONDITION)

Nodo CONDITION
  → Evalúa condición
  → Edge "yes" o "no"
  → Siguiente nodo según resultado

================================================================================
5. REEMPLAZO DE VARIABLES
================================================================================

SINTAXIS DE VARIABLES:
----------------------
- {nombre_variable} → Reemplazar con valor
- [nombre_variable] → Reemplazar con valor (alternativa)

DÓNDE SE REEMPLAZAN:
--------------------
1. En contenido de nodos TEXT:
   "content": "Hola {nombre_usuario}" → "Hola Juan Pérez"

2. En URLs de nodos HTTP:
   "url": "https://api.com/user/{user_id}" → "https://api.com/user/123"

3. En body de nodos HTTP:
   {
     "image": "{imagen_cedula}",
     "name": "{nombre_usuario}"
   }

4. En títulos de botones:
   "title": "Ver {producto}"

CÓMO REEMPLAZAR:
----------------
1. Obtener variables de session.variables
2. Buscar patrones {variable} o [variable] en strings
3. Reemplazar con session.variables[variable] si existe
4. Si no existe, dejar el patrón original o usar valor por defecto

EJEMPLO:
--------
Variables de sesión:
{
  "nombre_usuario": "Juan Pérez",
  "imagen_cedula": "https://example.com/image.jpg"
}

Texto: "Hola {nombre_usuario}, tu cédula {imagen_cedula} fue validada"
Resultado: "Hola Juan Pérez, tu cédula https://example.com/image.jpg fue validada"

================================================================================
6. VALIDACIONES Y ERRORES
================================================================================

VALIDACIONES DE RESPUESTAS:
----------------------------
Cuando un nodo TEXT tiene waitForResponse=true y validation:
- required: true → el valor no puede estar vacío
- minLength: 3 → el texto debe tener al menos 3 caracteres
- maxLength: 50 → el texto no debe exceder 50 caracteres
- pattern: regex → el texto debe cumplir el patrón

Si la validación falla:
1. Enviar mensaje de error al usuario
2. Volver al nodo que solicitó la respuesta
3. Pedir nuevamente la información
4. NO avanzar al siguiente nodo

MANEJO DE ERRORES HTTP:
-----------------------
Si un nodo HTTP falla:
1. Opción A: Guardar error en la variable de respuesta
   - variables["response_validacion"] = { error: "Connection failed" }
   - Continuar flujo, dejar que CONDITION maneje el error

2. Opción B: Lanzar excepción
   - Detener flujo
   - Marcar sesión como error
   - Enviar mensaje de error al usuario

RECOMENDACIÓN: Usar Opción A para mayor robustez.

TIMEOUTS:
---------
Si el usuario no responde después de X tiempo (ej: 30 minutos):
1. Marcar sesión como "abandoned"
2. Opcional: Enviar mensaje recordatorio
3. Si el usuario responde después, puede:
   - Reiniciar el flujo desde el principio
   - Continuar desde donde se quedó (si se mantiene la sesión)

================================================================================
7. EJEMPLO COMPLETO DE FLUJO
================================================================================

FLUJO: Validación de Cédula
----------------------------

NODO 1: TEXT (Bienvenida)
- content: "¡Hola! 👋 Bienvenido..."
- waitForResponse: false
- Edge → NODO 2

NODO 2: TEXT (Solicitar Nombre)
- content: "¿Cuál es tu nombre completo?"
- waitForResponse: true
- responseVariableName: "nombre_usuario"
- Espera respuesta...

[Usuario envía: "Juan Pérez"]
- Guardar: variables["nombre_usuario"] = "Juan Pérez"
- Edge → NODO 3

NODO 3: RESPONSE (Validar Nombre)
- variableName: "nombre_usuario"
- validation: minLength: 3
- Validación pasa
- Edge → NODO 4

NODO 4: TEXT (Solicitar Cédula)
- content: "Perfecto {nombre_usuario}, envía foto de tu cédula"
- Reemplazar: "Perfecto Juan Pérez, envía foto de tu cédula"
- waitForResponse: true
- responseVariableName: "imagen_cedula"
- responseType: "image"
- Espera respuesta...

[Usuario envía imagen]
- Guardar: variables["imagen_cedula"] = "image_id_123"
- Edge → NODO 5

NODO 5: RESPONSE (Validar Imagen)
- variableName: "imagen_cedula"
- Validación pasa
- Edge → NODO 6

NODO 6: HTTP (Validar Cédula OCR)
- method: POST
- url: "/api/whatsapp/ocr/validate-id"
- body: { "image": "{imagen_cedula}" }
- Reemplazar: { "image": "image_id_123" }
- responseVariable: "response_validacion_cedula"
- Ejecutar HTTP...
- Respuesta: { "valid": true }
- Guardar: variables["response_validacion_cedula"] = { "valid": true }
- Edge → NODO 7

NODO 7: CONDITION (¿Cédula Válida?)
- condition: response_validacion_cedula.valid == true
- Evaluar: true
- Edge "yes" → NODO 9
- Edge "no" → NODO 8

NODO 9: TEXT (Cédula Válida)
- content: "¡Excelente {nombre_usuario}! Tu cédula fue validada."
- Reemplazar: "¡Excelente Juan Pérez! Tu cédula fue validada."
- waitForResponse: false
- No hay más edges
- COMPLETAR SESIÓN

================================================================================
8. CHECKLIST DE IMPLEMENTACIÓN
================================================================================

□ Crear modelo de sesión (FlowSessionModel)
  - Almacenar: flowId, currentNodeId, variables, waitingForResponse, etc.

□ Crear repositorio de sesiones (FlowSessionRepository)
  - Métodos: createOrGetActiveSession, findActiveByConversation, update, save

□ Crear motor de flujos (FlowEngine)
  - Métodos: startFlow, processMessage, processNode, moveToNextNode

□ Crear procesadores de nodos:
  □ TextNodeProcessor
  □ ButtonsNodeProcessor
  □ HttpNodeProcessor
  □ ConditionNodeProcessor
  □ ResponseNodeProcessor
  □ AudioNodeProcessor

□ Integrar con webhook de WhatsApp:
  - Al recibir mensaje, buscar sesión activa
  - Si existe: procesar mensaje en contexto del flujo
  - Si no existe: iniciar flujo por defecto

□ Implementar reemplazo de variables:
  - Función replaceVariables() que busca {variable} y [variable]
  - Aplicar en: contenido de mensajes, URLs, bodies de HTTP, etc.

□ Manejar edges correctamente:
  - Buscar edges que salen del nodo actual
  - Seguir edge correcto según tipo de nodo
  - Manejar condiciones (yes/no)

□ Validaciones:
  - Validar respuestas según reglas del nodo
  - Manejar errores de validación
  - Manejar errores de HTTP

□ Timeouts y abandono:
  - Detectar sesiones inactivas
  - Marcar como abandonadas después de X tiempo

================================================================================
FIN DEL DOCUMENTO
================================================================================

