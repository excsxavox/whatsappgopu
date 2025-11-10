// Verificar estructura real de flows
const { MongoClient } = require('mongodb');

const uri = "mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority";

async function checkFlowsStructure() {
    const client = new MongoClient(uri);
    
    try {
        await client.connect();
        console.log("✅ Conectado a MongoDB");
        
        const db = client.db('whatsapp_bot');
        const flows = db.collection('flows');
        
        // Buscar TODOS los flows
        const allFlows = await flows.find({}).toArray();
        
        console.log(`\n📊 Total de flows: ${allFlows.length}\n`);
        
        if (allFlows.length > 0) {
            console.log("📋 ESTRUCTURA COMPLETA DE FLOWS:\n");
            allFlows.forEach((flow, index) => {
                console.log(`\n${'='.repeat(60)}`);
                console.log(`FLOW ${index + 1}:`);
                console.log('='.repeat(60));
                console.log(JSON.stringify(flow, null, 2));
            });
        } else {
            console.log("❌ No hay flows en la base de datos");
        }
        
        // Buscar el flow específico por si tiene nombre diferente
        console.log(`\n\n${'='.repeat(60)}`);
        console.log("🔍 Buscando flow: flow_1761770353752_amurb6yeq");
        console.log('='.repeat(60));
        
        const specificFlow = await flows.findOne({ _id: 'flow_1761770353752_amurb6yeq' });
        if (specificFlow) {
            console.log("✅ Encontrado:");
            console.log(JSON.stringify(specificFlow, null, 2));
        } else {
            console.log("❌ No encontrado con _id");
            
            // Buscar por otros campos posibles
            const byName = await flows.findOne({ _name: /flow_1761770353752_amurb6yeq/i });
            if (byName) {
                console.log("✅ Encontrado por _name:");
                console.log(JSON.stringify(byName, null, 2));
            }
        }
        
    } catch (error) {
        console.error("❌ Error:", error.message);
        console.error(error.stack);
    } finally {
        await client.close();
    }
}

checkFlowsStructure();

