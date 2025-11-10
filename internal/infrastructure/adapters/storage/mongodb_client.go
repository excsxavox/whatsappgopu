package storage

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient encapsula la conexión a MongoDB
type MongoClient struct {
	client   *mongo.Client
	db       *mongo.Database
	Database *mongo.Database // Exportada para uso en repositorios
}

// NewMongoClient crea una nueva conexión a MongoDB
func NewMongoClient(ctx context.Context, mongoURI, dbName string) (*MongoClient, error) {
	// Configurar opciones de cliente
	clientOpts := options.Client().
		ApplyURI(mongoURI).
		SetConnectTimeout(10 * time.Second).
		SetServerSelectionTimeout(5 * time.Second)

	// Conectar a MongoDB
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("error conectando a MongoDB: %w", err)
	}

	// Verificar conexión
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("error al hacer ping a MongoDB: %w", err)
	}

	db := client.Database(dbName)
	return &MongoClient{
		client:   client,
		db:       db,
		Database: db, // Exportar para repositorios
	}, nil
}

// GetDatabase retorna la base de datos
func (m *MongoClient) GetDatabase() *mongo.Database {
	return m.db
}

// Close cierra la conexión a MongoDB
func (m *MongoClient) Close(ctx context.Context) error {
	if m.client != nil {
		return m.client.Disconnect(ctx)
	}
	return nil
}

// CreateIndexes crea los índices necesarios para las colecciones
func (m *MongoClient) CreateIndexes(ctx context.Context) error {
	// Índices para messages
	messagesCollection := m.db.Collection("messages")

	// Índice por conversation_id + timestamp (queries principales)
	_, err := messagesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{
			"conversation_id":       1,
			"timestamps.created_at": -1,
		},
	})
	if err != nil {
		return fmt.Errorf("error creando índice conversation_id: %w", err)
	}

	// Índice por dedup_key (idempotencia) - ÚNICO
	_, err = messagesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "dedup_key", Value: 1}},
		Options: options.Index().SetUnique(true).SetSparse(true),
	})
	if err != nil {
		return fmt.Errorf("error creando índice dedup_key: %w", err)
	}

	// Índice por instance_id + timestamps (multi-instance)
	_, err = messagesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{
			"instance_id":           1,
			"timestamps.created_at": -1,
		},
	})
	if err != nil {
		return fmt.Errorf("error creando índice instance_id: %w", err)
	}

	// Índice por tenant_id (multi-tenant)
	_, err = messagesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "tenant_id", Value: 1},
			{Key: "timestamps.created_at", Value: -1},
		},
		Options: options.Index().SetSparse(true), // sparse porque tenant_id es opcional
	})
	if err != nil {
		return fmt.Errorf("error creando índice tenant_id: %w", err)
	}

	// Índice por from (búsquedas por remitente)
	_, err = messagesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "from", Value: 1},
			{Key: "timestamps.created_at", Value: -1},
		},
	})
	if err != nil {
		return fmt.Errorf("error creando índice from: %w", err)
	}

	// Índice por status (reporting)
	_, err = messagesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "status", Value: 1},
			{Key: "timestamps.created_at", Value: -1},
		},
	})
	if err != nil {
		return fmt.Errorf("error creando índice status: %w", err)
	}

	// Índices para sessions
	sessionsCollection := m.db.Collection("sessions")
	_, err = sessionsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "is_active", Value: 1}},
	})
	if err != nil {
		return fmt.Errorf("error creando índice sessions: %w", err)
	}

	// Índices para companies
	companiesCollection := m.db.Collection("companies")

	// Índice único por code
	_, err = companiesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "code", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("error creando índice companies.code: %w", err)
	}

	// Índice por is_active
	_, err = companiesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "is_active", Value: 1}},
	})
	if err != nil {
		return fmt.Errorf("error creando índice companies.is_active: %w", err)
	}

	// Índices para flows
	flowsCollection := m.db.Collection("flows")

	// Índice por instance_id + is_active
	_, err = flowsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "instance_id", Value: 1},
			{Key: "is_active", Value: 1},
		},
	})
	if err != nil {
		return fmt.Errorf("error creando índice flows.instance_id: %w", err)
	}

	// Índice por tenant_id
	_, err = flowsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "tenant_id", Value: 1}},
	})
	if err != nil {
		return fmt.Errorf("error creando índice flows.tenant_id: %w", err)
	}

	// Índice por is_default + instance_id (para buscar flujo por defecto)
	_, err = flowsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "is_default", Value: 1},
			{Key: "instance_id", Value: 1},
		},
	})
	if err != nil {
		return fmt.Errorf("error creando índice flows.is_default: %w", err)
	}

	// Índices para flow_sessions
	flowSessionsCollection := m.db.Collection("flow_sessions")

	// Índice por conversation_id + status (buscar sesión activa)
	_, err = flowSessionsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "conversation_id", Value: 1},
			{Key: "status", Value: 1},
		},
	})
	if err != nil {
		return fmt.Errorf("error creando índice flow_sessions.conversation_id: %w", err)
	}

	// Índice por flow_id
	_, err = flowSessionsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "flow_id", Value: 1}},
	})
	if err != nil {
		return fmt.Errorf("error creando índice flow_sessions.flow_id: %w", err)
	}

	// Índice por last_activity_at + status (para cleanup de sesiones inactivas)
	_, err = flowSessionsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "status", Value: 1},
			{Key: "last_activity_at", Value: 1},
		},
	})
	if err != nil {
		return fmt.Errorf("error creando índice flow_sessions.last_activity_at: %w", err)
	}

	// Índice por instance_id
	_, err = flowSessionsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "instance_id", Value: 1}},
	})
	if err != nil {
		return fmt.Errorf("error creando índice flow_sessions.instance_id: %w", err)
	}

	return nil
}
