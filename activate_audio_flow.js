// Activar flow de audio
const { MongoClient } = require('mongodb');

const uri = "mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority";

async function activateAudioFlow() {
    const client = new MongoClient(uri);
    
    try {
        await client.connect();
        console.log("✅ Conectado a MongoDB");
        
        const db = client.db('test');
        const flows = db.collection('flows');
        
        // Desactivar todos los flows
        await flows.updateMany({}, {
            $set: {
                _isActive: false,
                is_default: false
            }
        });
        
        console.log("✅ Todos los flows desactivados");
        
        // Activar el flow de audio como por defecto
        const result = await flows.updateOne(
            { _id: 'flow_1762360447819_wf2jy6z6d' },
            {
                $set: {
                    _isActive: true,
                    is_default: true,
                    instance_id: 'default',
                    tenant_id: 'default'
                }
            }
        );
        
        console.log(`\n✅ Flow de audio activado como por defecto`);
        console.log(`   - Modified: ${result.modifiedCount}`);
        
        // Verificar
        const flow = await flows.findOne({ _id: 'flow_1762360447819_wf2jy6z6d' });
        console.log("\n📋 Flow configurado:");
        console.log(`   - ID: ${flow._id}`);
        console.log(`   - Name: ${flow._name}`);
        console.log(`   - Active: ${flow._isActive}`);
        console.log(`   - Default: ${flow.is_default}`);
        console.log(`   - Nodos: ${flow._flowData.nodes.length}`);
        console.log(`   - Tipo primer nodo: ${flow._flowData.nodes[0].type}`);
        
    } catch (error) {
        console.error("❌ Error:", error.message);
    } finally {
        await client.close();
    }
}

activateAudioFlow();

