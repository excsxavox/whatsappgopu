package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"whatsapp-api-go/internal/domain/entities"
	"whatsapp-api-go/internal/domain/ports"
)

// MongoContextoRepository implementa ports.ContextoRepository con MongoDB
type MongoContextoRepository struct {
	collection *mongo.Collection
}

// NewMongoContextoRepository crea un nuevo repositorio de contextos con MongoDB
func NewMongoContextoRepository(db *mongo.Database) ports.ContextoRepository {
	return &MongoContextoRepository{
		collection: db.Collection("contexto"),
	}
}

// Save guarda o actualiza un contexto
func (r *MongoContextoRepository) Save(ctx context.Context, contexto *entities.Contexto) error {
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": contexto.ID}
	update := bson.M{"$set": contexto}
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// FindByID busca un contexto por ID
func (r *MongoContextoRepository) FindByID(ctx context.Context, contextoID string) (*entities.Contexto, error) {
	var contexto entities.Contexto
	err := r.collection.FindOne(ctx, bson.M{"_id": contextoID}).Decode(&contexto)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &contexto, err
}

// FindByCompanyID busca un contexto por companyID
func (r *MongoContextoRepository) FindByCompanyID(ctx context.Context, companyID string) (*entities.Contexto, error) {
	var contexto entities.Contexto
	err := r.collection.FindOne(ctx, bson.M{"_companyId": companyID}).Decode(&contexto)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &contexto, err
}

// FindActiveByCompanyID busca un contexto activo por companyID
func (r *MongoContextoRepository) FindActiveByCompanyID(ctx context.Context, companyID string) (*entities.Contexto, error) {
	var contexto entities.Contexto
	filter := bson.M{
		"_companyId": companyID,
		"_isActive":  true,
	}
	err := r.collection.FindOne(ctx, filter).Decode(&contexto)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &contexto, err
}

// FindAll busca todos los contextos
func (r *MongoContextoRepository) FindAll(ctx context.Context) ([]*entities.Contexto, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var contextos []*entities.Contexto
	if err = cursor.All(ctx, &contextos); err != nil {
		return nil, err
	}
	return contextos, nil
}

// Delete elimina un contexto
func (r *MongoContextoRepository) Delete(ctx context.Context, contextoID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": contextoID})
	return err
}
