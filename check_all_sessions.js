const { MongoClient } = require('mongodb');

const uri = "mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0";
const dbName = 'test';

async function checkAllSessions() {
    const client = new MongoClient(uri);

    try {
        await client.connect();
        console.log("✅ Conectado a MongoDB\n");

        const db = client.db(dbName);
        const sessionsCollection = db.collection('flow_sessions');

        // Buscar TODAS las sesiones ordenadas por fecha
        const allSessions = await sessionsCollection.find({}).sort({ createdAt: -1 }).limit(10).toArray();
        
        console.log(`📋 Últimas ${allSessions.length} sesiones:\n`);
        
        allSessions.forEach((s, i) => {
            console.log(`${i + 1}. Sesión ${s._id}`);
            console.log(`   Conversation ID: ${s.conversation_id}`);
            console.log(`   Flow ID: ${s.flow_id}`);
            console.log(`   Current Node: ${s.current_node_id}`);
            console.log(`   Status: ${s.status}`);
            console.log(`   Waiting: ${s.waiting_for_response}`);
            console.log(`   Variable esperada: ${s.waiting_for_variable || 'N/A'}`);
            console.log(`   Variables: ${JSON.stringify(s.variables || {})}`);
            console.log(`   Creada: ${s.createdAt}`);
            console.log(`   Última actividad: ${s.last_activity_at}`);
            console.log(`   Ejecutados: ${JSON.stringify(s.executed_nodes || [])}`);
            console.log('');
        });

    } catch (error) {
        console.error('❌ Error:', error.message);
    } finally {
        await client.close();
    }
}

checkAllSessions();

