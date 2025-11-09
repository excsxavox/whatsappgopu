package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"whatsapp-api-go/internal/domain/entities"
	"whatsapp-api-go/internal/domain/ports"
)

// MongoMessageRepository implementa ports.MessageRepository con MongoDB
type MongoMessageRepository struct {
	collection *mongo.Collection
}

// NewMongoMessageRepository crea un nuevo repositorio de mensajes con MongoDB
func NewMongoMessageRepository(db *mongo.Database) ports.MessageRepository {
	return &MongoMessageRepository{
		collection: db.Collection("messages"),
	}
}

// Save guarda un mensaje en MongoDB
func (r *MongoMessageRepository) Save(ctx context.Context, message *entities.Message) error {
	// Upsert por _id (wamid para incoming, generar para outgoing)
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": message.ID}
	update := bson.M{"$set": message}

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// FindByID busca un mensaje por ID (wamid)
func (r *MongoMessageRepository) FindByID(ctx context.Context, messageID string) (*entities.Message, error) {
	var message entities.Message
	err := r.collection.FindOne(ctx, bson.M{"_id": messageID}).Decode(&message)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, entities.ErrMessageNotFound
		}
		return nil, err
	}
	return &message, nil
}

// FindByDedupKey busca un mensaje por dedup_key (idempotencia)
func (r *MongoMessageRepository) FindByDedupKey(ctx context.Context, dedupKey string) (*entities.Message, error) {
	var message entities.Message
	err := r.collection.FindOne(ctx, bson.M{"dedup_key": dedupKey}).Decode(&message)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No es error, simplemente no existe
		}
		return nil, err
	}
	return &message, nil
}

// FindByConversation busca mensajes por conversation_id
func (r *MongoMessageRepository) FindByConversation(ctx context.Context, conversationID string, limit int) ([]*entities.Message, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamps.created_at", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, bson.M{"conversation_id": conversationID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*entities.Message
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

// FindByRecipient busca mensajes por destinatario (para compatibilidad)
func (r *MongoMessageRepository) FindByRecipient(ctx context.Context, recipient string, limit int) ([]*entities.Message, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamps.created_at", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, bson.M{"to": recipient}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*entities.Message
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

// FindByInstance busca mensajes por instance_id
func (r *MongoMessageRepository) FindByInstance(ctx context.Context, instanceID string, limit int) ([]*entities.Message, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamps.created_at", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, bson.M{"instance_id": instanceID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*entities.Message
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

// FindByTenant busca mensajes por tenant_id (multi-tenant)
func (r *MongoMessageRepository) FindByTenant(ctx context.Context, tenantID string, limit int) ([]*entities.Message, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamps.created_at", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, bson.M{"tenant_id": tenantID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*entities.Message
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

// UpdateStatus actualiza el estado de un mensaje
func (r *MongoMessageRepository) UpdateStatus(ctx context.Context, messageID string, status string) error {
	// Nota: Este método es legacy, el nuevo modelo usa message.UpdateStatus() y luego Save()
	filter := bson.M{"_id": messageID}
	update := bson.M{"$set": bson.M{"status": status}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return entities.ErrMessageNotFound
	}

	return nil
}

// ExistsByDedupKey verifica si existe un mensaje con ese dedup_key
func (r *MongoMessageRepository) ExistsByDedupKey(ctx context.Context, dedupKey string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"dedup_key": dedupKey}, options.Count().SetLimit(1))
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
