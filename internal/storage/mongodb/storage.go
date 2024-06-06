package mongodb

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"order/internal/models"
	"order/internal/storage"
)

type Storage struct {
	DB *mongo.Client
}

var (
	ErrorNotFound = errors.New("not found")
)

func New(storagePath string) (*Storage, error) {
	const op = "storage.mongodb.New"

	db, err := mongo.Connect(context.Background(), options.Client().ApplyURI(storagePath))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{DB: db}, nil

}

func (s *Storage) CreateOrder(ctx context.Context, order *models.Order) (primitive.ObjectID, error) {
	const op = "storage.mongodb.CreateOrder"

	collection := s.DB.Database("pizzeria").Collection("orders")

	result, err := collection.InsertOne(ctx, order)
	if err != nil {
		var writeException mongo.WriteException
		if errors.As(err, &writeException) {
			for _, we := range writeException.WriteErrors {
				if we.Code == 11000 {
					return primitive.NilObjectID, fmt.Errorf("%s: %w", op, storage.ErrorOrderAlreadyExists)
				}
			}
		}
	}

	insertedID := result.InsertedID.(primitive.ObjectID)
	return insertedID, nil
}

func (s *Storage) OrderByID(ctx context.Context, id primitive.ObjectID) (*models.Order, error) {
	const op = "storage.mongodb.OrderByID"

	collection := s.DB.Database("pizzeria").Collection("orders")

	var order models.Order
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, ErrorNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &order, nil
}

func (s *Storage) DeleteOrder(ctx context.Context, id primitive.ObjectID) error {
	const op = "storage.mongodb.DeleteOrder"

	collection := s.DB.Database("pizzeria").Collection("orders")

	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
