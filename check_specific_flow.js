// Buscar flujo específico
const { MongoClient } = require('mongodb');

const uri = "mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority";

async function checkFlow() {
    const client = new MongoClient(uri);
    
    try {
        await client.connect();
        console.log("✅ Conectado a MongoDB");
        
        const db = client.db('whatsapp_bot');
        const flows = db.collection('flows');
        
        // Buscar flujo específico
        const flow = await flows.findOne({ _id: 'flow_1761770353752_amurb6yeq' });
        
        if (flow) {
            console.log("\n✅ Flujo encontrado:");
            console.log(JSON.stringify(flow, null, 2));
        } else {
            console.log("\n❌ Flujo NO encontrado con ID: flow_1761770353752_amurb6yeq");
            
            // Listar todos los flujos disponibles
            const allFlows = await flows.find({}).toArray();
            console.log("\n📋 Flujos disponibles:");
            allFlows.forEach(f => {
                console.log(`   - ${f._id} (${f.name})`);
            });
        }
        
    } catch (error) {
        console.error("❌ Error:", error.message);
    } finally {
        await client.close();
    }
}

checkFlow();

