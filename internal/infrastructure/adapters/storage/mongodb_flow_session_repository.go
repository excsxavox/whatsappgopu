package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"whatsapp-api-go/internal/domain/entities"
	"whatsapp-api-go/internal/domain/ports"
)

type mongoFlowSessionRepository struct {
	collection *mongo.Collection
}

// NewMongoFlowSessionRepository crea un nuevo repositorio de sesiones
func NewMongoFlowSessionRepository(client *MongoClient) ports.FlowSessionRepository {
	return &mongoFlowSessionRepository{
		collection: client.Database.Collection("flow_sessions"),
	}
}

func (r *mongoFlowSessionRepository) Save(ctx context.Context, session *entities.FlowSession) error {
	if session.ID == "" {
		session.ID = uuid.New().String()
	}

	_, err := r.collection.InsertOne(ctx, session)
	if err != nil {
		return fmt.Errorf("error saving flow session: %w", err)
	}

	return nil
}

func (r *mongoFlowSessionRepository) FindByID(ctx context.Context, sessionID string) (*entities.FlowSession, error) {
	var session entities.FlowSession

	err := r.collection.FindOne(ctx, bson.M{"_id": sessionID}).Decode(&session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("flow session not found: %s", sessionID)
		}
		return nil, fmt.Errorf("error finding flow session: %w", err)
	}

	return &session, nil
}

func (r *mongoFlowSessionRepository) FindActiveByConversation(ctx context.Context, conversationID string) (*entities.FlowSession, error) {
	var session entities.FlowSession

	filter := bson.M{
		"conversation_id": conversationID,
		"status":          "active",
	}

	err := r.collection.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No hay sesión activa, devolver nil sin error
		}
		return nil, fmt.Errorf("error finding active session: %w", err)
	}

	return &session, nil
}

func (r *mongoFlowSessionRepository) Update(ctx context.Context, session *entities.FlowSession) error {
	session.UpdatedAt = time.Now()

	filter := bson.M{"_id": session.ID}
	update := bson.M{"$set": session}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error updating flow session: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("flow session not found: %s", session.ID)
	}

	return nil
}

func (r *mongoFlowSessionRepository) FindInactiveSessions(ctx context.Context, minutesInactive int) ([]*entities.FlowSession, error) {
	threshold := time.Now().Add(-time.Duration(minutesInactive) * time.Minute)

	filter := bson.M{
		"status":           "active",
		"last_activity_at": bson.M{"$lt": threshold},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error finding inactive sessions: %w", err)
	}
	defer cursor.Close(ctx)

	var sessions []*entities.FlowSession
	if err := cursor.All(ctx, &sessions); err != nil {
		return nil, fmt.Errorf("error decoding sessions: %w", err)
	}

	return sessions, nil
}


