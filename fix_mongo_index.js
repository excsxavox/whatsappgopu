// Script para eliminar el índice problemático id_1 en test.messages
const { MongoClient } = require('mongodb');

const uri = "mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0";

async function fixIndex() {
    const client = new MongoClient(uri);
    
    try {
        await client.connect();
        console.log("✅ Conectado a MongoDB");
        
        const db = client.db('test');
        const messages = db.collection('messages');
        
        // Listar índices actuales
        console.log("\n📋 Índices actuales:");
        const indexes = await messages.indexes();
        indexes.forEach(idx => {
            console.log(`   - ${idx.name}:`, JSON.stringify(idx.key));
        });
        
        // Eliminar el índice problemático id_1
        try {
            await messages.dropIndex('id_1');
            console.log("\n✅ Índice 'id_1' eliminado exitosamente");
        } catch (err) {
            if (err.codeName === 'IndexNotFound') {
                console.log("\n⚠️  Índice 'id_1' no existe");
            } else {
                throw err;
            }
        }
        
        // Listar índices después
        console.log("\n📋 Índices después de la limpieza:");
        const indexesAfter = await messages.indexes();
        indexesAfter.forEach(idx => {
            console.log(`   - ${idx.name}:`, JSON.stringify(idx.key));
        });
        
    } catch (error) {
        console.error('❌ Error:', error.message);
    } finally {
        await client.close();
    }
}

fixIndex();

