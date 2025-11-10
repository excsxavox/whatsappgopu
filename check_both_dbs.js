// Script para verificar mensajes en AMBAS bases de datos
const { MongoClient } = require('mongodb');

const uri = "mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0";

async function checkMessages() {
    const client = new MongoClient(uri);
    
    try {
        await client.connect();
        console.log("✅ Conectado a MongoDB\n");
        
        // Revisar base de datos TEST
        console.log("========================================");
        console.log("📊 BASE DE DATOS: TEST");
        console.log("========================================");
        const dbTest = client.db('test');
        const messagesTest = dbTest.collection('messages');
        
        const lastMessagesTest = await messagesTest.find({})
            .sort({ 'timestamps.created_at': -1 })
            .limit(5)
            .toArray();
        
        console.log(`\nÚltimos ${lastMessagesTest.length} mensajes:\n`);
        
        lastMessagesTest.forEach((msg, i) => {
            console.log(`${i + 1}. ${msg.direction === 'in' ? '📨 ENTRANTE' : '📤 SALIENTE'}`);
            console.log(`   De: ${msg.from || 'N/A'}`);
            console.log(`   Para: ${msg.to || 'N/A'}`);
            console.log(`   Tipo: ${msg.message?.type || msg.type || 'N/A'}`);
            console.log(`   Texto: ${msg.message?.text?.body || msg.text?.body || 'N/A'}`);
            console.log(`   Fecha: ${msg.timestamps?.created_at || 'N/A'}`);
            console.log(`   Estado: ${msg.status}`);
            console.log('');
        });
        
        const totalTest = await messagesTest.countDocuments({});
        console.log(`📝 Total de mensajes en TEST: ${totalTest}\n`);
        
        // Revisar base de datos WHATSAPP_BOT
        console.log("========================================");
        console.log("📊 BASE DE DATOS: WHATSAPP_BOT");
        console.log("========================================");
        const dbWhatsapp = client.db('whatsapp_bot');
        const messagesWhatsapp = dbWhatsapp.collection('messages');
        
        const lastMessagesWhatsapp = await messagesWhatsapp.find({})
            .sort({ 'timestamps.created_at': -1 })
            .limit(5)
            .toArray();
        
        console.log(`\nÚltimos ${lastMessagesWhatsapp.length} mensajes:\n`);
        
        lastMessagesWhatsapp.forEach((msg, i) => {
            console.log(`${i + 1}. ${msg.direction === 'in' ? '📨 ENTRANTE' : '📤 SALIENTE'}`);
            console.log(`   De: ${msg.from || 'N/A'}`);
            console.log(`   Para: ${msg.to || 'N/A'}`);
            console.log(`   Tipo: ${msg.message?.type || msg.type || 'N/A'}`);
            console.log(`   Texto: ${msg.message?.text?.body || msg.text?.body || 'N/A'}`);
            console.log(`   Fecha: ${msg.timestamps?.created_at || 'N/A'}`);
            console.log(`   Estado: ${msg.status}`);
            console.log('');
        });
        
        const totalWhatsapp = await messagesWhatsapp.countDocuments({});
        console.log(`📝 Total de mensajes en WHATSAPP_BOT: ${totalWhatsapp}\n`);
        
        // Revisar base de datos WHATSAPP_API
        console.log("========================================");
        console.log("📊 BASE DE DATOS: WHATSAPP_API");
        console.log("========================================");
        const dbWhatsappApi = client.db('whatsapp_api');
        const messagesWhatsappApi = dbWhatsappApi.collection('messages');
        
        const lastMessagesWhatsappApi = await messagesWhatsappApi.find({})
            .sort({ 'timestamps.created_at': -1 })
            .limit(5)
            .toArray();
        
        console.log(`\nÚltimos ${lastMessagesWhatsappApi.length} mensajes:\n`);
        
        lastMessagesWhatsappApi.forEach((msg, i) => {
            console.log(`${i + 1}. ${msg.direction === 'in' ? '📨 ENTRANTE' : '📤 SALIENTE'}`);
            console.log(`   De: ${msg.from || 'N/A'}`);
            console.log(`   Para: ${msg.to || 'N/A'}`);
            console.log(`   Tipo: ${msg.message?.type || msg.type || 'N/A'}`);
            console.log(`   Texto: ${msg.message?.text?.body || msg.text?.body || 'N/A'}`);
            console.log(`   Fecha: ${msg.timestamps?.created_at || 'N/A'}`);
            console.log(`   Estado: ${msg.status}`);
            console.log('');
        });
        
        const totalWhatsappApi = await messagesWhatsappApi.countDocuments({});
        console.log(`📝 Total de mensajes en WHATSAPP_API: ${totalWhatsappApi}\n`);
        
    } catch (error) {
        console.error('❌ Error:', error.message);
    } finally {
        await client.close();
    }
}

checkMessages();
