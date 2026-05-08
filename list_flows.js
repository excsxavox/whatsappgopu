const { MongoClient } = require('mongodb');

const uri = "mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0";

async function listFlows() {
    const client = new MongoClient(uri);
    
    try {
        await client.connect();
        console.log("✅ Conectado a MongoDB");
        
        const db = client.db("whatsapp");
        
        // Listar colecciones
        const collections = await db.listCollections().toArray();
        console.log("\n📂 Colecciones disponibles:");
        collections.forEach(c => console.log(`   - ${c.name}`));
        
        // Buscar flow específico
        const flowId = "flow_1762360447819_wf2jy6z6d";
        console.log(`\n🔍 Buscando flow: ${flowId}\n`);
        
        const flowsCollection = db.collection("flows");
        
        // Intentar diferentes formatos de búsqueda
        let flow = await flowsCollection.findOne({ _id: flowId });
        if (flow) {
            console.log("✅ Encontrado por _id");
        } else {
            flow = await flowsCollection.findOne({ id: flowId });
            if (flow) console.log("✅ Encontrado por id");
        }
        
        if (flow) {
            console.log(`\n📊 Flow: ${flow.name || flow._id}`);
            console.log(`   Nodes: ${flow.definition?.nodes?.length || 0}`);
            console.log(`   Edges: ${flow.definition?.edges?.length || 0}`);
            
            if (flow.definition?.edges) {
                console.log("\n🔗 Edges desde Condition-1762872376949:");
                const conditionEdges = flow.definition.edges.filter(e => 
                    e.from === "Condition-1762872376949"
                );
                
                if (conditionEdges.length > 0) {
                    conditionEdges.forEach(edge => {
                        console.log(`   ${edge.from} -> ${edge.to}`);
                        console.log(`      sourceHandle: ${edge.sourceHandle}`);
                        console.log(`      condition: ${edge.condition || 'none'}`);
                    });
                } else {
                    console.log("   ❌ No hay edges desde este nodo");
                }
            }
        } else {
            console.log("❌ Flow no encontrado");
            
            // Mostrar todos los flows
            const allFlows = await flowsCollection.find({}).limit(5).toArray();
            console.log(`\n📋 Flows disponibles (primeros 5):`);
            allFlows.forEach(f => {
                console.log(`   - ${f._id || f.id}`);
            });
        }
        
    } catch (error) {
        console.error("❌ Error:", error.message);
    } finally {
        await client.close();
    }
}

listFlows();

