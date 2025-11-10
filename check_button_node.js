const { MongoClient } = require('mongodb');

const uri = "mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0";
const dbName = 'test';
const flowId = 'flow_1761770353752_amurb6yeq';
const nodeId = 'node_6_menu_opciones';

async function checkButtonNode() {
    const client = new MongoClient(uri);

    try {
        await client.connect();
        console.log("✅ Conectado a MongoDB\n");

        const db = client.db(dbName);
        const flowsCollection = db.collection('flows');

        const flow = await flowsCollection.findOne({ _id: flowId });

        if (flow) {
            const node = flow._flowData.nodes.find(n => n.id === nodeId);
            
            if (node) {
                console.log("📋 Nodo de botones encontrado:\n");
                console.log(JSON.stringify(node, null, 2));
                
                console.log("\n📝 Estructura del config:\n");
                console.log(JSON.stringify(node.config, null, 2));
                
                if (node.config.action) {
                    console.log("\n🔍 Action encontrado:");
                    console.log(JSON.stringify(node.config.action, null, 2));
                }
                
                if (node.config.buttons) {
                    console.log("\n🔘 Botones directos:");
                    console.log(JSON.stringify(node.config.buttons, null, 2));
                }
            } else {
                console.log(`⚠️ Nodo ${nodeId} no encontrado en el flow`);
            }
        } else {
            console.log(`⚠️ Flow ${flowId} no encontrado`);
        }

    } catch (error) {
        console.error('❌ Error:', error.message);
    } finally {
        await client.close();
    }
}

checkButtonNode();

