// Script para verificar mensajes en MongoDB
const { MongoClient } = require('mongodb');

const uri = "mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0";

async function checkMessages() {
    const client = new MongoClient(uri);
    
    try {
        await client.connect();
        console.log("✅ Conectado a MongoDB");
        
        const db = client.db('test');
        const messages = db.collection('messages');
        
        // Buscar los últimos 10 mensajes
        const lastMessages = await messages.find({})
            .sort({ 'timestamps.created_at': -1 })
            .limit(10)
            .toArray();
        
        console.log(`\n📊 Últimos ${lastMessages.length} mensajes:\n`);
        
        lastMessages.forEach((msg, i) => {
            console.log(`${i + 1}. ${msg.direction === 'in' ? '📨 ENTRANTE' : '📤 SALIENTE'}`);
            console.log(`   De: ${msg.from || 'N/A'}`);
            console.log(`   Para: ${msg.to || 'N/A'}`);
            console.log(`   Tipo: ${msg.message?.type || msg.type || 'N/A'}`);
            console.log(`   Texto: ${msg.message?.text?.body || msg.text?.body || 'N/A'}`);
            console.log(`   Fecha: ${msg.timestamps?.created_at || 'N/A'}`);
            console.log(`   Estado: ${msg.status}`);
            console.log('');
        });
        
        // Contar mensajes totales
        const total = await messages.countDocuments({});
        console.log(`📝 Total de mensajes en DB: ${total}`);
        
    } catch (error) {
        console.error('❌ Error:', error.message);
    } finally {
        await client.close();
    }
}

checkMessages();

