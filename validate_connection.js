// Validar conexión completa a MongoDB
const { MongoClient } = require('mongodb');

const uri = "mongodb+srv://nexti:sL1Vr3NSs46rB0ZLU7wl3VC8GV@cluster0.acnpcls.mongodb.net/?retryWrites=true&w=majority";

async function validateConnection() {
    const client = new MongoClient(uri);
    
    try {
        console.log("🔌 Conectando a MongoDB...");
        await client.connect();
        console.log("✅ Conexión exitosa\n");
        
        // Listar TODAS las bases de datos
        const adminDb = client.db().admin();
        const { databases } = await adminDb.listDatabases();
        
        console.log("📂 BASES DE DATOS DISPONIBLES:");
        console.log("=".repeat(60));
        databases.forEach(db => {
            console.log(`   - ${db.name} (${(db.sizeOnDisk / 1024 / 1024).toFixed(2)} MB)`);
        });
        
        console.log("\n" + "=".repeat(60));
        console.log("🔍 ANALIZANDO CADA BASE DE DATOS:");
        console.log("=".repeat(60));
        
        // Analizar cada base de datos
        for (const database of databases) {
            if (database.name === 'admin' || database.name === 'local' || database.name === 'config') {
                continue; // Saltar bases de sistema
            }
            
            console.log(`\n📁 Base de datos: ${database.name}`);
            const db = client.db(database.name);
            const collections = await db.listCollections().toArray();
            
            console.log(`   Colecciones (${collections.length}):`);
            for (const col of collections) {
                const count = await db.collection(col.name).countDocuments();
                console.log(`      - ${col.name}: ${count} documentos`);
                
                // Si es la colección flows, mostrar los IDs
                if (col.name === 'flows') {
                    const flows = await db.collection('flows').find({}, { projection: { _id: 1, name: 1 } }).toArray();
                    if (flows.length > 0) {
                        console.log(`         FLOWS:`);
                        flows.forEach(f => {
                            console.log(`            - ${f._id}: ${f.name || 'Sin nombre'}`);
                        });
                    }
                }
            }
        }
        
        console.log("\n" + "=".repeat(60));
        console.log("🔍 BUSCANDO FLOW ESPECÍFICO: flow_1761770353752_amurb6yeq");
        console.log("=".repeat(60));
        
        // Buscar en TODAS las bases de datos
        for (const database of databases) {
            if (database.name === 'admin' || database.name === 'local' || database.name === 'config') {
                continue;
            }
            
            const db = client.db(database.name);
            const collections = await db.listCollections().toArray();
            
            for (const col of collections) {
                if (col.name === 'flows') {
                    const flow = await db.collection('flows').findOne({ _id: 'flow_1761770353752_amurb6yeq' });
                    if (flow) {
                        console.log(`\n✅ ENCONTRADO en: ${database.name}.flows`);
                        console.log(JSON.stringify(flow, null, 2));
                        return;
                    }
                }
            }
        }
        
        console.log("\n❌ Flow NO encontrado en ninguna base de datos");
        
    } catch (error) {
        console.error("❌ Error de conexión:", error.message);
        console.error(error.stack);
    } finally {
        await client.close();
        console.log("\n🔌 Conexión cerrada");
    }
}

validateConnection();

