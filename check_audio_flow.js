// Verificar flow con audio
const { MongoClient } = require('mongodb');

const uri = "mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority";

async function checkAudioFlow() {
    const client = new MongoClient(uri);
    
    try {
        await client.connect();
        console.log("✅ Conectado a MongoDB");
        
        const db = client.db('test');
        const flows = db.collection('flows');
        
        // Buscar flow específico
        const flow = await flows.findOne({ _id: 'flow_1762360447819_wf2jy6z6d' });
        
        if (flow) {
            console.log("\n✅ FLOW ENCONTRADO: flow_1762360447819_wf2jy6z6d");
            console.log("=".repeat(80));
            console.log(`Nombre: ${flow._name || 'Sin nombre'}`);
            console.log(`Descripción: ${flow._description || 'Sin descripción'}`);
            console.log(`Estado: ${flow._isActive ? 'ACTIVO ✅' : 'INACTIVO ❌'}`);
            console.log(`Status: ${flow._status}`);
            
            if (flow._flowData) {
                const flowData = flow._flowData;
                console.log(`\nEntry Node: ${flowData.entryNodeId}`);
                console.log(`Total de nodos: ${flowData.nodes ? flowData.nodes.length : 0}`);
                console.log(`Total de edges: ${flowData.edges ? flowData.edges.length : 0}`);
                
                console.log("\n" + "=".repeat(80));
                console.log("📋 NODOS DEL FLOW:");
                console.log("=".repeat(80));
                
                if (flowData.nodes) {
                    flowData.nodes.forEach((node, index) => {
                        console.log(`\n${index + 1}. ${node.id}`);
                        console.log(`   Tipo: ${node.type}`);
                        console.log(`   Config:`);
                        console.log(JSON.stringify(node.config, null, 4));
                    });
                }
                
                console.log("\n" + "=".repeat(80));
                console.log("🔗 EDGES (CONEXIONES):");
                console.log("=".repeat(80));
                
                if (flowData.edges) {
                    flowData.edges.forEach((edge, index) => {
                        console.log(`${index + 1}. ${edge.from} → ${edge.to}`);
                    });
                }
                
                // Buscar nodos de AUDIO
                console.log("\n" + "=".repeat(80));
                console.log("🎵 NODOS DE AUDIO:");
                console.log("=".repeat(80));
                
                const audioNodes = flowData.nodes ? flowData.nodes.filter(n => n.type === 'AUDIO') : [];
                if (audioNodes.length > 0) {
                    audioNodes.forEach((node, index) => {
                        console.log(`\n✅ Audio ${index + 1}: ${node.id}`);
                        console.log(`   URL: ${node.config.audioUrl || node.config.url || 'N/A'}`);
                        console.log(`   Config completo:`);
                        console.log(JSON.stringify(node.config, null, 4));
                    });
                } else {
                    console.log("\n⚠️  No se encontraron nodos de tipo AUDIO");
                }
            }
            
        } else {
            console.log("\n❌ Flow NO encontrado: flow_1762360447819_wf2jy6z6d");
        }
        
    } catch (error) {
        console.error("❌ Error:", error.message);
        console.error(error.stack);
    } finally {
        await client.close();
    }
}

checkAudioFlow();

