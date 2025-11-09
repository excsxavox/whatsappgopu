package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"whatsapp-api-go/internal/domain/entities"
	"whatsapp-api-go/internal/domain/ports"
)

// MongoCompanyRepository implementa ports.CompanyRepository con MongoDB
type MongoCompanyRepository struct {
	collection *mongo.Collection
}

// NewMongoCompanyRepository crea un nuevo repositorio de empresas con MongoDB
func NewMongoCompanyRepository(db *mongo.Database) ports.CompanyRepository {
	return &MongoCompanyRepository{
		collection: db.Collection("companies"),
	}
}

// Save guarda o actualiza una empresa
func (r *MongoCompanyRepository) Save(ctx context.Context, company *entities.Company) error {
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": company.ID}
	update := bson.M{"$set": company}
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// FindByID busca una empresa por ID
func (r *MongoCompanyRepository) FindByID(ctx context.Context, companyID string) (*entities.Company, error) {
	var company entities.Company
	err := r.collection.FindOne(ctx, bson.M{"_id": companyID}).Decode(&company)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &company, err
}

// FindByCode busca una empresa por código
func (r *MongoCompanyRepository) FindByCode(ctx context.Context, code string) (*entities.Company, error) {
	var company entities.Company
	err := r.collection.FindOne(ctx, bson.M{"code": code}).Decode(&company)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &company, err
}

// FindAll busca todas las empresas
func (r *MongoCompanyRepository) FindAll(ctx context.Context) ([]*entities.Company, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var companies []*entities.Company
	if err = cursor.All(ctx, &companies); err != nil {
		return nil, err
	}
	return companies, nil
}

// FindActive busca todas las empresas activas
func (r *MongoCompanyRepository) FindActive(ctx context.Context) ([]*entities.Company, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"is_active": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var companies []*entities.Company
	if err = cursor.All(ctx, &companies); err != nil {
		return nil, err
	}
	return companies, nil
}

// FindInactive busca todas las empresas inactivas
func (r *MongoCompanyRepository) FindInactive(ctx context.Context) ([]*entities.Company, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"is_active": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var companies []*entities.Company
	if err = cursor.All(ctx, &companies); err != nil {
		return nil, err
	}
	return companies, nil
}

// Delete elimina una empresa
func (r *MongoCompanyRepository) Delete(ctx context.Context, companyID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": companyID})
	return err
}
