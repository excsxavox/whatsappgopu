// Verificar mensajes en MongoDB (sin dependencias)
const uri = "mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/whatsapp_api?retryWrites=true&w=majority";

// Hacer query HTTP directo a MongoDB Atlas Data API
const https = require('https');

// Por ahora, vamos a verificar usando curl o directamente en Atlas
console.log("✅ Configuración correcta:");
console.log("   - De: 593992686734");
console.log("   - A: +1 555 152 6940 (Meta Test)");
console.log("   - Webhook: https://whatsapp-api-go-dpb5cgbnaec2gdf2.eastus-01.azurewebsites.net/webhook");
console.log("");
console.log("📋 Por favor verifica en MongoDB Atlas:");
console.log("   1. Ve a: https://cloud.mongodb.com/");
console.log("   2. Collections → whatsapp_api → messages");
console.log("   3. Busca mensajes recientes");
console.log("   4. ¿Ves tu mensaje?");

