const { MongoClient } = require('mongodb');

const uri = "mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0";

async function checkSession() {
    const client = new MongoClient(uri);

    try {
        await client.connect();
        console.log("✅ Conectado a MongoDB\n");

        const db = client.db('test');
        const sessionsCollection = db.collection('flow_sessions');

        // Buscar la sesión específica
        const session = await sessionsCollection.findOne({ 
            _id: "e69b5913-125e-4856-9f8f-c343bb9a31b2" 
        });

        if (session) {
            console.log("📋 Sesión encontrada:");
            console.log(`   ID: ${session._id}`);
            console.log(`   Conversation ID: ${session.conversation_id}`);
            console.log(`   Flow ID: ${session.flow_id}`);
            console.log(`   Current Node: ${session.current_node_id}`);
            console.log(`   Status: ${session.status}`);
            console.log(`   Waiting for Response: ${session.waiting_for_response}`);
            console.log(`   Waiting for Variable: ${session.waiting_for_variable}`);
            console.log(`   Created: ${session.created_at}`);
            console.log(`   Updated: ${session.updated_at}`);
            console.log(`   Variables:`, JSON.stringify(session.variables, null, 2));
            console.log('');
        } else {
            console.log("⚠️ Sesión no encontrada");
        }

        // Buscar TODAS las sesiones activas para este conversationID
        console.log("\n📋 Todas las sesiones para esta conversación:");
        const allSessions = await sessionsCollection.find({ 
            conversation_id: "593992686734" 
        }).toArray();
        
        allSessions.forEach((s, i) => {
            console.log(`\n${i + 1}. Sesión ${s._id}:`);
            console.log(`   Status: ${s.status}`);
            console.log(`   Current Node: ${s.current_node_id}`);
            console.log(`   Waiting: ${s.waiting_for_response}`);
            console.log(`   Created: ${s.created_at}`);
        });

    } catch (error) {
        console.error('❌ Error:', error.message);
    } finally {
        await client.close();
    }
}

checkSession();

