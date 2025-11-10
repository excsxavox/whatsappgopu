const { MongoClient } = require('mongodb');

const uri = "mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0";

async function cleanErrorSessions() {
    const client = new MongoClient(uri);

    try {
        await client.connect();
        console.log("✅ Conectado a MongoDB\n");

        const db = client.db('test');
        const sessionsCollection = db.collection('flow_sessions');

        // Eliminar todas las sesiones con error
        const result = await sessionsCollection.deleteMany({ status: "error" });
        console.log(`✅ ${result.deletedCount} sesiones con error eliminadas\n`);

        // Mostrar sesiones restantes
        const allSessions = await sessionsCollection.find({}).toArray();
        console.log(`📋 Sesiones restantes: ${allSessions.length}`);
        allSessions.forEach((s, i) => {
            console.log(`   ${i + 1}. ${s._id} - conversation_id: ${s.conversation_id} - status: ${s.status}`);
        });

    } catch (error) {
        console.error('❌ Error:', error.message);
    } finally {
        await client.close();
    }
}

cleanErrorSessions();

