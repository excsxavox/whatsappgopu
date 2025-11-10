const { MongoClient } = require('mongodb');

const uri = "mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0";

async function checkFlow() {
    const client = new MongoClient(uri);

    try {
        await client.connect();
        console.log("✅ Conectado a MongoDB\n");

        const db = client.db('test');
        const flowsCollection = db.collection('flows');

        // Buscar el flow específico
        const flow = await flowsCollection.findOne({ _id: "flow_1761770353752_amurb6yeq" });

        if (flow) {
            console.log("📋 Flow encontrado:");
            console.log(`   ID: ${flow._id}`);
            console.log(`   Nombre: ${flow._name}`);
            console.log(`   Activo: ${flow._isActive}`);
            console.log(`   Default: ${flow.is_default}`);
            console.log(`   Instance ID: ${flow.instance_id}`);
            console.log(`   Tenant ID: ${flow.tenant_id}`);
            console.log(`   Entry Node: ${flow._flowData?.entryNodeId}`);
            console.log(`   Nodos: ${flow._flowData?.nodes?.length || 0}`);
            console.log(`   Edges: ${flow._flowData?.edges?.length || 0}`);
            
            if (flow._flowData?.nodes) {
                console.log("\n📝 Nodos del flow:");
                flow._flowData.nodes.forEach((node, i) => {
                    console.log(`   ${i + 1}. ${node.id} (${node.type})`);
                    if (node.type === 'TEXT') {
                        console.log(`      Texto: ${node.config.content || node.config.bodyText}`);
                    } else if (node.type === 'AUDIO') {
                        console.log(`      Has audio: ${node.config.hasRecordedAudio}`);
                        console.log(`      Wait response: ${node.config.waitForVoiceResponse}`);
                    }
                });
            }
            
            console.log("\n✅ Configurando como flow por defecto...");
            
            // Desactivar otros flows por defecto
            await flowsCollection.updateMany(
                { instance_id: "804818756055720", is_default: true },
                { $set: { is_default: false } }
            );
            
            // Activar este flow como default
            await flowsCollection.updateOne(
                { _id: "flow_1761770353752_amurb6yeq" },
                { 
                    $set: { 
                        is_default: true, 
                        _isActive: true,
                        instance_id: "804818756055720",
                        tenant_id: "default"
                    } 
                }
            );
            
            console.log("✅ Flow configurado como por defecto\n");
            
        } else {
            console.log("❌ Flow no encontrado");
            console.log("\n📋 Flows disponibles:");
            const allFlows = await flowsCollection.find({}).toArray();
            allFlows.forEach((f, i) => {
                console.log(`   ${i + 1}. ${f._id} - ${f._name}`);
            });
        }

    } catch (error) {
        console.error('❌ Error:', error.message);
    } finally {
        await client.close();
    }
}

checkFlow();

