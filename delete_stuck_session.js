const { MongoClient } = require('mongodb');

const uri = "mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0";
const dbName = 'test';
const sessionId = 'c7fa6561-edf2-418a-b3c7-fc21d83ef5c6'; // La sesión atascada en node_6_menu_opciones

async function deleteSession() {
    const client = new MongoClient(uri);

    try {
        await client.connect();
        console.log("✅ Conectado a MongoDB\n");

        const db = client.db(dbName);
        const sessionsCollection = db.collection('flow_sessions');

        // Eliminar la sesión atascada
        const result = await sessionsCollection.deleteOne({ _id: sessionId });
        console.log(`✅ Sesión eliminada: ${result.deletedCount} documento(s)\n`);

        // Mostrar sesiones restantes
        const remaining = await sessionsCollection.find({ conversation_id: '593992686734@804818756055720' }).toArray();
        console.log(`📋 Sesiones restantes para esta conversación: ${remaining.length}\n`);

    } catch (error) {
        console.error('❌ Error:', error.message);
    } finally {
        await client.close();
    }
}

deleteSession();

