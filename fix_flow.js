// Activar flujo por defecto
const { MongoClient } = require('mongodb');

const uri = "mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority";

async function fixFlow() {
    const client = new MongoClient(uri);
    
    try {
        await client.connect();
        console.log("✅ Conectado a MongoDB");
        
        const db = client.db('whatsapp_bot');
        const flows = db.collection('flows');
        
        // Desactivar todos
        await flows.updateMany({}, {
            $set: {
                is_active: false,
                is_default: false
            }
        });
        
        // Activar el primero como por defecto
        const result = await flows.updateOne(
            { _id: 'flow_1761169428062_3nuye99dx' },
            {
                $set: {
                    is_active: true,
                    is_default: true,
                    instance_id: 'default',
                    tenant_id: 'default'
                }
            }
        );
        
        console.log("\n✅ Flujo activado como por defecto");
        console.log(`   - Modified: ${result.modifiedCount}`);
        
        // Verificar
        const flow = await flows.findOne({ _id: 'flow_1761169428062_3nuye99dx' });
        console.log("\n📋 Flujo configurado:");
        console.log(`   - ID: ${flow._id}`);
        console.log(`   - Name: ${flow.name}`);
        console.log(`   - Active: ${flow.is_active}`);
        console.log(`   - Default: ${flow.is_default}`);
        
    } catch (error) {
        console.error("❌ Error:", error.message);
    } finally {
        await client.close();
    }
}

fixFlow();

