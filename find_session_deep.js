const { MongoClient } = require('mongodb');

const MONGO_URI = 'mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0';

async function findSession() {
    const client = new MongoClient(MONGO_URI);
    
    try {
        await client.connect();
        console.log('✅ Conectado a MongoDB');
        
        const db = client.db('whatsapp');
        
        // Listar todas las colecciones
        const collections = await db.listCollections().toArray();
        console.log('\n📂 Colecciones en la base de datos "whatsapp":');
        
        for (const col of collections) {
            const collName = col.name;
            const count = await db.collection(collName).countDocuments();
            console.log(`   - ${collName}: ${count} documentos`);
            
            // Buscar la sesión en cada colección
            const session = await db.collection(collName).findOne({
                $or: [
                    { _id: '6a3110eb-aaf6-4d22-a639-0c51eb2a3574' },
                    { id: '6a3110eb-aaf6-4d22-a639-0c51eb2a3574' },
                    { sessionId: '6a3110eb-aaf6-4d22-a639-0c51eb2a3574' },
                    { conversationId: '593992686734@804818756055720' }
                ]
            });
            
            if (session) {
                console.log(`\n🎯 ¡Sesión encontrada en ${collName}!`);
                console.log('   _id:', session._id);
                console.log('   conversationId:', session.conversationId);
                console.log('   status:', session.status);
                console.log('   waitingForResponse:', session.waitingForResponse);
                
                // Eliminar la sesión
                await db.collection(collName).deleteOne({ _id: session._id });
                console.log('   🗑️  ¡Sesión eliminada!');
            }
        }
        
        console.log('\n✅ Búsqueda completada');
        
    } catch (error) {
        console.error('❌ Error:', error.message);
    } finally {
        await client.close();
    }
}

findSession();

