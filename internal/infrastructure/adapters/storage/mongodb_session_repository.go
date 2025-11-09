package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"whatsapp-api-go/internal/domain/entities"
	"whatsapp-api-go/internal/domain/ports"
)

// MongoSessionRepository implementa ports.SessionRepository con MongoDB
type MongoSessionRepository struct {
	collection *mongo.Collection
}

// NewMongoSessionRepository crea un nuevo repositorio de sesiones con MongoDB
func NewMongoSessionRepository(db *mongo.Database) ports.SessionRepository {
	return &MongoSessionRepository{
		collection: db.Collection("sessions"),
	}
}

// sessionDoc es el documento de MongoDB para sesiones
type sessionDoc struct {
	ID             string     `bson:"_id"`
	PhoneNumber    string     `bson:"phone_number"`
	DeviceName     string     `bson:"device_name,omitempty"`
	IsActive       bool       `bson:"is_active"`
	IsConnected    bool       `bson:"is_connected"`
	ConnectedAt    *time.Time `bson:"connected_at,omitempty"`
	LastSeen       *time.Time `bson:"last_seen,omitempty"`
	DisconnectedAt *time.Time `bson:"disconnected_at,omitempty"`
}

// Save guarda o actualiza una sesión
func (r *MongoSessionRepository) Save(ctx context.Context, session *entities.Session) error {
	doc := sessionDoc{
		ID:             session.ID,
		PhoneNumber:    session.PhoneNumber,
		DeviceName:     session.DeviceName,
		IsActive:       session.IsActive,
		IsConnected:    session.IsConnected,
		ConnectedAt:    session.ConnectedAt,
		LastSeen:       session.LastSeen,
		DisconnectedAt: session.DisconnectedAt,
	}

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": session.ID}
	update := bson.M{"$set": doc}

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// FindByID busca una sesión por ID
func (r *MongoSessionRepository) FindByID(ctx context.Context, sessionID string) (*entities.Session, error) {
	var doc sessionDoc
	err := r.collection.FindOne(ctx, bson.M{"_id": sessionID}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, entities.ErrSessionNotFound
		}
		return nil, err
	}

	return &entities.Session{
		ID:             doc.ID,
		PhoneNumber:    doc.PhoneNumber,
		DeviceName:     doc.DeviceName,
		IsActive:       doc.IsActive,
		IsConnected:    doc.IsConnected,
		ConnectedAt:    doc.ConnectedAt,
		LastSeen:       doc.LastSeen,
		DisconnectedAt: doc.DisconnectedAt,
	}, nil
}

// FindActive busca la sesión activa
func (r *MongoSessionRepository) FindActive(ctx context.Context) (*entities.Session, error) {
	var doc sessionDoc
	err := r.collection.FindOne(ctx, bson.M{"is_active": true}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, entities.ErrSessionNotFound
		}
		return nil, err
	}

	return &entities.Session{
		ID:             doc.ID,
		PhoneNumber:    doc.PhoneNumber,
		DeviceName:     doc.DeviceName,
		IsActive:       doc.IsActive,
		IsConnected:    doc.IsConnected,
		ConnectedAt:    doc.ConnectedAt,
		LastSeen:       doc.LastSeen,
		DisconnectedAt: doc.DisconnectedAt,
	}, nil
}

// Delete elimina una sesión
func (r *MongoSessionRepository) Delete(ctx context.Context, sessionID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": sessionID})
	return err
}

// MarkAsInactive marca una sesión como inactiva
func (r *MongoSessionRepository) MarkAsInactive(ctx context.Context, sessionID string) error {
	filter := bson.M{"_id": sessionID}
	update := bson.M{"$set": bson.M{"is_active": false}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return entities.ErrSessionNotFound
	}

	return nil
}
