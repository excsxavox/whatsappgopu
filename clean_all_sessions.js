const { MongoClient } = require('mongodb');

const uri = "mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0";
const dbName = 'test';
const conversationID = '593992686734@804818756055720';

async function cleanAllSessions() {
    const client = new MongoClient(uri);

    try {
        await client.connect();
        console.log("✅ Conectado a MongoDB\n");

        const db = client.db(dbName);
        const sessionsCollection = db.collection('flow_sessions');

        // Eliminar TODAS las sesiones de esta conversación
        const result = await sessionsCollection.deleteMany({ conversation_id: conversationID });
        console.log(`✅ ${result.deletedCount} sesiones eliminadas\n`);

        // Verificar que no queden sesiones
        const remaining = await sessionsCollection.find({ conversation_id: conversationID }).toArray();
        console.log(`📋 Sesiones restantes: ${remaining.length}\n`);

    } catch (error) {
        console.error('❌ Error:', error.message);
    } finally {
        await client.close();
    }
}

cleanAllSessions();

