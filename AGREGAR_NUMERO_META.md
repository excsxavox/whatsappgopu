# 📱 Cómo Agregar tu Número en Meta (GUÍA COMPLETA)

## 🎯 PASO 1: Ir a la Configuración Correcta

1. Abre: https://developers.facebook.com/apps/10058963160806734/whatsapp-business/wa-dev-console/
2. En el menú izquierdo, busca **"WhatsApp"**
3. Haz clic en **"API Setup"**

---

## 🎯 PASO 2: Agregar Número de Prueba

En la sección **"Phone Numbers"** o **"Test Phone Numbers"**, busca un botón que dice:
- **"+ Add phone number"**
- O **"Manage phone numbers"**
- O **"Add recipient"**

### Formato CORRECTO del número:

```
+593992686734
```

✅ **CON** el símbolo `+`  
✅ **SIN** espacios  
✅ **SIN** guiones  
✅ **CON** código de país (593 para Ecuador)

---

## 🎯 PASO 3: Verificar que Aparece en la Lista

Después de agregar, deberías ver tu número en una lista:

```
+593992686734    ✅ Verified
```

O

```
593992686734     ✅ Active
```

---

## 🎯 PASO 4: Esperar Propagación

Meta puede tardar **5-15 minutos** en propagar el cambio a todos sus servidores.

**Durante este tiempo:**
- ❌ NO envíes más mensajes
- ⏳ Espera al menos 10 minutos
- ☕ Tómate un café

---

## 🎯 PASO 5: Verificar Modo de la App

En la parte superior de la página, verifica que diga:

```
🟢 Production
```

O

```
🟢 Live
```

**NO** debe decir:
```
🟠 Development
```

---

## 🆘 SI NO ENCUENTRAS DÓNDE AGREGAR EL NÚMERO

Hay **3 lugares diferentes** donde podría estar:

### Opción A: API Setup → Phone Numbers
https://developers.facebook.com/apps/10058963160806734/whatsapp-business/wa-dev-console/

### Opción B: WhatsApp → Configuration → Phone Numbers
https://developers.facebook.com/apps/10058963160806734/whatsapp-business/wa-settings/

### Opción C: WhatsApp Business Account → Phone Numbers
https://business.facebook.com/latest/whatsapp_manager

---

## 📸 ¿QUÉ DEBO VER?

Comparte una captura de pantalla de:

1. La sección donde aparece tu número agregado
2. El toggle/botón que muestra "Production" o "Live"
3. Cualquier mensaje de error que veas

---

## 🔄 DESPUÉS DE AGREGAR

Una vez que hayas:
1. ✅ Agregado el número con formato correcto (+593992686734)
2. ✅ Verificado que aparece en la lista
3. ✅ Confirmado que está en modo Production
4. ⏳ Esperado 10 minutos

Entonces **envía un mensaje de prueba** y revisa los logs.

---

## 🎯 ALTERNATIVA: Usar Otro Número

Si no puedes agregar tu número personal, puedes:

1. **Usar el número de prueba de WhatsApp Business**:
   - Descargar WhatsApp Business en otro celular
   - Registrar un número diferente
   - Agregar ESE número como tester

2. **Pedir a Meta que active tu cuenta**:
   - En algunos casos, necesitas que Meta revise tu app manualmente
   - Ve a: https://developers.facebook.com/docs/whatsapp/get-started/
   - Solicita revisión de la app

---

## 📞 CONTACTO CON META

Si nada funciona, contacta al soporte de Meta:
- https://business.facebook.com/business/help
- Menciona el error: `#131030 Recipient phone number not in allowed list`
- Proporciona tu App ID: `10058963160806734`
- Proporciona tu Phone Number ID: `804818756055720`

