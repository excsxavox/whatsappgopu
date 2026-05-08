const { MongoClient } = require('mongodb');

const MONGO_URI = 'mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0';

async function forceDropSessions() {
    const client = new MongoClient(MONGO_URI);
    
    try {
        await client.connect();
        console.log('✅ Conectado a MongoDB');
        
        // Probar múltiples bases de datos
        const dbNames = ['whatsapp', 'test', 'whatsapp_api'];
        
        for (const dbName of dbNames) {
            console.log(`\n🔍 Revisando base de datos: ${dbName}`);
            const db = client.db(dbName);
            
            try {
                const flowSessions = db.collection('flow_sessions');
                const count = await flowSessions.countDocuments();
                
                if (count > 0) {
                    console.log(`   📊 Encontradas ${count} sesiones`);
                    
                    // Eliminar TODAS
                    const result = await flowSessions.deleteMany({});
                    console.log(`   🗑️  ${result.deletedCount} sesiones eliminadas`);
                }
            } catch (err) {
                console.log(`   ⚠️  No se pudo acceder a la colección`);
            }
        }
        
        console.log('\n✅ Limpieza completada');
        
    } catch (error) {
        console.error('❌ Error:', error.message);
    } finally {
        await client.close();
    }
}

forceDropSessions();

