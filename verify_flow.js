// Script para verificar el flow
const { MongoClient } = require('mongodb');

const uri = "mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0";

async function verifyFlow() {
    const client = new MongoClient(uri);
    
    try {
        await client.connect();
        console.log("✅ Conectado a MongoDB");
        
        const db = client.db('test');
        const flows = db.collection('flows');
        
        // Buscar flows por defecto
        console.log("\n📋 Buscando flows por defecto con instance_id: 804818756055720");
        const defaultFlows = await flows.find({
            instance_id: '804818756055720',
            is_default: true,
            _isActive: true
        }).toArray();
        
        console.log(`\nEncontrados: ${defaultFlows.length} flow(s)\n`);
        
        defaultFlows.forEach((flow, i) => {
            console.log(`${i + 1}. Flow ID: ${flow._id}`);
            console.log(`   Nombre: ${flow._name}`);
            console.log(`   instance_id: ${flow.instance_id}`);
            console.log(`   is_default: ${flow.is_default}`);
            console.log(`   _isActive: ${flow._isActive}`);
            console.log(`   tenant_id: ${flow.tenant_id || 'N/A'}`);
            console.log('');
        });
        
        // Buscar el flow específico
        console.log("\n📋 Buscando flow específico: flow_1762360447819_wf2jy6z6d");
        const specificFlow = await flows.findOne({ _id: 'flow_1762360447819_wf2jy6z6d' });
        
        if (specificFlow) {
            console.log(`\n✅ Flow encontrado:`);
            console.log(`   ID: ${specificFlow._id}`);
            console.log(`   Nombre: ${specificFlow._name}`);
            console.log(`   instance_id: ${specificFlow.instance_id}`);
            console.log(`   is_default: ${specificFlow.is_default}`);
            console.log(`   _isActive: ${specificFlow._isActive}`);
            console.log(`   tenant_id: ${specificFlow.tenant_id || 'N/A'}`);
        } else {
            console.log("❌ Flow NO encontrado");
        }
        
    } catch (error) {
        console.error('❌ Error:', error.message);
    } finally {
        await client.close();
    }
}

verifyFlow();

