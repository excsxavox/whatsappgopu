// Verificar flujos en MongoDB
const { MongoClient } = require('mongodb');

const uri = "mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority";

async function checkFlows() {
    const client = new MongoClient(uri);
    
    try {
        await client.connect();
        console.log("✅ Conectado a MongoDB");
        
        const db = client.db('whatsapp_bot');
        
        // Verificar colecciones
        const collections = await db.listCollections().toArray();
        console.log("\n📂 Colecciones existentes:");
        collections.forEach(c => console.log(`   - ${c.name}`));
        
        // Verificar flows
        const flows = db.collection('flows');
        const flowCount = await flows.countDocuments();
        console.log(`\n🔄 Total de flujos: ${flowCount}`);
        
        if (flowCount > 0) {
            const allFlows = await flows.find({}).toArray();
            console.log("\n📋 Flujos:");
            allFlows.forEach(f => {
                console.log(`   - ID: ${f._id}`);
                console.log(`     Name: ${f.name || 'N/A'}`);
                console.log(`     Active: ${f.is_active}`);
                console.log(`     Default: ${f.is_default}`);
            });
        }
        
        // Verificar flow_sessions
        const sessions = db.collection('flow_sessions');
        const sessionCount = await sessions.countDocuments();
        console.log(`\n💬 Total de sesiones de flujo: ${sessionCount}`);
        
    } catch (error) {
        console.error("❌ Error:", error.message);
    } finally {
        await client.close();
    }
}

checkFlows();

