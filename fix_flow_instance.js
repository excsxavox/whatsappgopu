// Script para actualizar el flow con el instance_id correcto
const { MongoClient } = require('mongodb');

const uri = "mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0";

async function fixFlowInstance() {
    const client = new MongoClient(uri);
    
    try {
        await client.connect();
        console.log("✅ Conectado a MongoDB");
        
        const db = client.db('test');
        const flows = db.collection('flows');
        
        // Buscar el flow por ID
        const flowId = 'flow_1762360447819_wf2jy6z6d';
        const flow = await flows.findOne({ _id: flowId });
        
        if (!flow) {
            console.log(`❌ Flow ${flowId} no encontrado`);
            return;
        }
        
        console.log(`\n📋 Flow actual:`);
        console.log(`   ID: ${flow._id}`);
        console.log(`   Nombre: ${flow._name}`);
        console.log(`   instance_id: ${flow.instance_id || 'NO CONFIGURADO'}`);
        console.log(`   is_default: ${flow.is_default}`);
        console.log(`   is_active: ${flow._isActive}`);
        
        // Actualizar con el instance_id correcto
        const instanceId = '804818756055720';
        
        const result = await flows.updateOne(
            { _id: flowId },
            {
                $set: {
                    instance_id: instanceId,
                    is_default: true,
                    _isActive: true
                }
            }
        );
        
        console.log(`\n✅ Flow actualizado: ${result.modifiedCount} documento(s) modificado(s)`);
        
        // Verificar actualización
        const updatedFlow = await flows.findOne({ _id: flowId });
        console.log(`\n📋 Flow después de actualizar:`);
        console.log(`   ID: ${updatedFlow._id}`);
        console.log(`   instance_id: ${updatedFlow.instance_id}`);
        console.log(`   is_default: ${updatedFlow.is_default}`);
        console.log(`   is_active: ${updatedFlow._isActive}`);
        
    } catch (error) {
        console.error('❌ Error:', error.message);
    } finally {
        await client.close();
    }
}

fixFlowInstance();

