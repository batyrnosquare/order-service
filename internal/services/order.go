package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
	"order/internal/models"
	"time"
)

type OrderService struct {
	log           *slog.Logger
	orderProvider OrderRepository
	tokenTTL      time.Duration
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *models.Order) (primitive.ObjectID, error)
	OrderByID(ctx context.Context, id primitive.ObjectID) (*models.Order, error)
	DeleteOrder(ctx context.Context, id primitive.ObjectID) error
}

func New(log *slog.Logger, orderProvider OrderRepository, tokenTTL time.Duration) *OrderService {
	return &OrderService{
		log:           log,
		orderProvider: orderProvider,
		tokenTTL:      tokenTTL,
	}

}

func (o *OrderService) CreateOrder(ctx context.Context, order *models.Order) (primitive.ObjectID, error) {
	const op = "OrderService.CreateOrder"

	log := o.log.With(slog.String("op", op))

	log.Info("creating order")

	orderID, err := o.orderProvider.CreateOrder(ctx, order)
	if err != nil {
		log.Error("failed to create order", slog.Error)
		return primitive.NilObjectID, err
	}
	return orderID, nil
}

func (o *OrderService) OrderByID(ctx context.Context, id primitive.ObjectID) (*models.Order, error) {
	const op = "OrderService.OrderByID"

	log := o.log.With(slog.String("op", op))

	log.Info("getting order by id")

	order, err := o.orderProvider.OrderByID(ctx, id)
	if err != nil {
		log.Error("failed to get order by id", slog.Error)
		return nil, err
	}
	return order, nil
}

func (o *OrderService) DeleteOrder(ctx context.Context, id primitive.ObjectID) error {
	const op = "OrderService.DeleteOrder"

	log := o.log.With(slog.String("op", op))

	log.Info("deleting order")

	err := o.orderProvider.DeleteOrder(ctx, id)
	if err != nil {
		log.Error("failed to delete order", slog.Error)
		return err
	}
	return nil
}
