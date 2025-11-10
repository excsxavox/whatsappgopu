package storage

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"whatsapp-api-go/internal/domain/entities"
	"whatsapp-api-go/internal/domain/ports"
)

type mongoFlowRepository struct {
	collection *mongo.Collection
}

// NewMongoFlowRepository crea un nuevo repositorio de flujos
func NewMongoFlowRepository(client *MongoClient) ports.FlowRepository {
	return &mongoFlowRepository{
		collection: client.Database.Collection("flows"),
	}
}

func (r *mongoFlowRepository) Save(ctx context.Context, flow *entities.Flow) error {
	if flow.ID == "" {
		return fmt.Errorf("flow ID is required")
	}

	flow.CreatedAt = time.Now()
	flow.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, flow)
	if err != nil {
		return fmt.Errorf("error saving flow: %w", err)
	}

	return nil
}

func (r *mongoFlowRepository) FindByID(ctx context.Context, flowID string) (*entities.Flow, error) {
	var flow entities.Flow

	err := r.collection.FindOne(ctx, bson.M{"_id": flowID}).Decode(&flow)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("flow not found: %s", flowID)
		}
		return nil, fmt.Errorf("error finding flow: %w", err)
	}

	return &flow, nil
}

func (r *mongoFlowRepository) FindDefault(ctx context.Context, instanceID string) (*entities.Flow, error) {
	var flow entities.Flow

	filter := bson.M{
		"instance_id": instanceID,
		"is_default":  true,
		"_isActive":   true,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&flow)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no default flow found for instance: %s", instanceID)
		}
		return nil, fmt.Errorf("error finding default flow: %w", err)
	}

	return &flow, nil
}

func (r *mongoFlowRepository) FindByTenant(ctx context.Context, tenantID string) ([]*entities.Flow, error) {
	filter := bson.M{"tenant_id": tenantID}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error finding flows: %w", err)
	}
	defer cursor.Close(ctx)

	var flows []*entities.Flow
	if err := cursor.All(ctx, &flows); err != nil {
		return nil, fmt.Errorf("error decoding flows: %w", err)
	}

	return flows, nil
}

func (r *mongoFlowRepository) Update(ctx context.Context, flow *entities.Flow) error {
	flow.UpdatedAt = time.Now()

	filter := bson.M{"_id": flow.ID}
	update := bson.M{"$set": flow}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error updating flow: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("flow not found: %s", flow.ID)
	}

	return nil
}

func (r *mongoFlowRepository) Delete(ctx context.Context, flowID string) error {
	filter := bson.M{"_id": flowID}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting flow: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("flow not found: %s", flowID)
	}

	return nil
}


